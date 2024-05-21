package main

import (
	logger "agentless/infra/log"
	"agentless/inventory/readers"
	"errors"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWSFetcher struct {
	regions    []string
	inventory  map[string]readers.Resource
	updates    chan readers.Resource
	ticker     *time.Ticker
	sema       chan struct{} // semaphore allows X maximum readers running concurrently
	done       chan bool
	repository Repository
	cycle      int
	mu         sync.Mutex
	fetching   bool
}

func NewAWSFetcher(interval time.Duration, maxConcurrency int, regions []string, repository Repository) *AWSFetcher {
	return &AWSFetcher{
		regions:    regions,
		inventory:  make(map[string]readers.Resource),
		updates:    make(chan readers.Resource, 100),
		ticker:     time.NewTicker(interval),
		sema:       make(chan struct{}, maxConcurrency),
		done:       make(chan bool),
		repository: repository,
		cycle:      0,
	}
}

func (f *AWSFetcher) Start() {
	logger.Log.Infof("Start AWS Fetcher. regions=%v", f.regions)
	go f.scheduleFetchCycle()
}

func (f *AWSFetcher) Stop() {
	f.ticker.Stop()
	close(f.done)
	close(f.updates)
	logger.Log.Info("Stop Aws Fetcher")
}

func (f *AWSFetcher) scheduleFetchCycle() {
	if err := f.fetch(); err != nil {
		logger.Log.Errorf("Scheduled fetch error: %s", err)
	}

	for {
		select {
		case <-f.ticker.C:
			if err := f.fetch(); err != nil {
				logger.Log.Errorf("Scheduled fetch error: %s", err)
			}
		case <-f.done:
			return
		}
	}
}

func (f *AWSFetcher) fetch() error {
	if err := f.startFetch(); err != nil {
		return err
	}

	doneUpdating := make(chan bool)
	go f.handleUpdates(doneUpdating)

	f.read()

	f.endFetch()

	<-doneUpdating
	f.repository.Update(f.cycle, &f.inventory)
	return nil
}

func (f *AWSFetcher) startFetch() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fetching {
		return errors.New("fetch already in progress")
	}
	logger.Log.Infof("Start fetch cycle %d", f.cycle)
	f.fetching = true
	f.updates = make(chan readers.Resource, 100)
	f.cycle++
	f.inventory = make(map[string]readers.Resource)
	return nil
}

func (f *AWSFetcher) endFetch() {
	f.mu.Lock()
	defer f.mu.Unlock()
	close(f.updates)
	f.fetching = false
	logger.Log.Infof("End fetch cycle %d", f.cycle)
}

func (f *AWSFetcher) read() {
	logger.Log.Infof("AWSFetcher#%d Starting read", f.cycle)

	var wg sync.WaitGroup
	for _, region := range f.regions {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewSharedCredentials("", "default"),
		})

		if err != nil {
			logger.Log.Fatalf("Failed to create session for region %s: %s", region, err)
			continue
		}

		regionReaders := readers.GetRegionReaders(sess, region, f.updates)

		for _, reader := range regionReaders {
			f.sema <- struct{}{} // Acquire a slot
			wg.Add(1)
			go func(r readers.ResourceReader) {
				defer wg.Done()
				r.Read()
				<-f.sema
			}(reader)
		}
	}
	wg.Wait()
}

func (f *AWSFetcher) handleUpdates(doneUpdating chan<- bool) {
	for item := range f.updates {
		readers.AddDefaultLabels(item)
		f.inventory[item.ID] = item
		logger.Log.Debugf("AWSFetcher#%d item:%s, type:%s, region: %s fetched",
			f.cycle, item.ID, item.Type, item.Region)
	}
	doneUpdating <- true
}
