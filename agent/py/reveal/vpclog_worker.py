import threading
import time
import boto3
from datetime import datetime, timedelta, timezone
from .vpclog_reader import FlowLogsS3Reader
from .vpclog_reveal import Reveal
from .ipmap_fetcher import fetch_ipmap
from common.logger import get_logger


LOG = get_logger(module_name=__name__)

class VPCLogWorker:

    def __init__(self, region: str, location: str, interval: int, 
                 reporting_entity_id: str) -> None:
        self.region = region
        self.location = location
        self.interval = interval
        self.reader_lock = threading.Lock()
        self.running = False
        self.worker_thread = None
        self.boto_client = boto3.client("s3")
        self.reporting_entity_id = reporting_entity_id

    def start(self): 
        if not self.running: 
            LOG.info(f'Starting vpclog worker for region {self.region}')
            self.start_time = datetime.now() - timedelta(hours=1)
            self.running = True
            self.worker_thread = threading.Thread(target=self.run)
            self.worker_thread.daemon = True
            self.worker_thread.start()
            self.worker_thread.join()

    def run(self):
        next_run_time = time.time() + self.interval
        while self.running:
            if not self.reader_lock.locked():
                start_time = time.time()
                self.reader_lock.acquire()
                try:
                    self.vpclog_reader()
                finally:
                    self.reader_lock.release()
                duration = time.time() - start_time
                sleep_time = max(0, next_run_time - time.time())
            else:
                LOG.warn("Reader is already working, skipping this cycle.")
                sleep_time = max(0, next_run_time - time.time())

            time.sleep(sleep_time)
            next_run_time += self.interval  # Schedule the next run

    def stop(self): 
        self.running = False
        if self.worker_thread: 
            self.worker_thread.join()
            LOG.info(f'Stop vpclog worker for region {self.region}')

    def vpclog_reader(self):
        ipmap = fetch_ipmap()
        reader = FlowLogsS3Reader(location=self.location, boto_client=self.boto_client)
        reader.memorize_previous_runs()

        reveal = Reveal(ipmap=ipmap, reporting_entity_id=self.reporting_entity_id)

        for record in reader: 
            reveal.send(record)