package readers

import (
	logger "agentless/infra/log"
	model "agentless/infra/model/common"
	utils "agentless/infra/utils"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

const EC2 = "EC2"

type EC2Reader struct {
	Svc     ec2iface.EC2API
	Region  string
	updates chan Resource
}

func NewEC2Reader(sess *session.Session, region string, resource chan Resource) *EC2Reader {
	return &EC2Reader{
		Svc:     ec2.New(sess),
		Region:  region,
		updates: resource,
	}
}

func (r *EC2Reader) Read() {
	logger.Log.Infof("Reader Started: Type=EC2, region=%s", r.Region)
	var err error = nil

	input := &ec2.DescribeInstancesInput{}
	err = r.Svc.DescribeInstancesPages(input, func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				item, err := ToInventoryItemFrom(instance)
				if err != nil {
					logger.Log.Errorf("failed to discover ec2 %s", *instance.InstanceId)
					continue
				}
				r.updates <- Resource{ID: *instance.InstanceId, Region: r.Region, Type: "EC2", Item: item}
			}
		}
		return !lastPage
	})

	if err != nil {
		logger.Log.Fatalf("Failed read instances: %s", err)
	}

	logger.Log.Infof("Reader Completed: Type=EC2, region=%s", r.Region)
}

func ToInventoryItemFrom(instance *ec2.Instance) (*model.InventoryItem, error) {
	entityData := &model.InventoryItem_EntityData{}
	entityData.FromVMData(*ToVMDataFrom(instance))
	tags := AwsTagsToMap(instance.Tags)

	item := &model.InventoryItem{
		EntityCategory: utils.StrPtr("compute"),
		EntityData:     entityData,
		EntityName:     utils.StrPtr(GetValueOrDefault(tags, "Name", *instance.InstanceId)),
		EntityType:     utils.StrPtr("virtual machine"),
		ExternalIds:    utils.SlicePtr([]string{*instance.InstanceId}),
		ItemId:         instance.InstanceId,
		ItemType:       (*model.InventoryItemItemType)(utils.StrPtr(string(model.Asset))),
		Labels:         AwsLabelsListFrom(instance.Tags),
	}

	return item, nil
}

func ToVMDataFrom(instance *ec2.Instance) *model.VMData {
	nics := &[]model.NetworkInterfaceData{}

	for _, ec2Nic := range instance.NetworkInterfaces {
		nicData := ToNetwotkInterfaceDataFrom(ec2Nic)
		*nics = append(*nics, *nicData)
	}

	return &model.VMData{
		Type: model.VM,
		Nics: nics,
	}
}

func ToNetwotkInterfaceDataFrom(nic *ec2.InstanceNetworkInterface) *model.NetworkInterfaceData {
	publicIPs := []string{}
	privateIPAddresses := []string{}

	for _, ipEntry := range nic.PrivateIpAddresses {
		if ipEntry.Association != nil && ipEntry.Association.PublicIp != nil {
			publicIPs = append(publicIPs, *ipEntry.Association.PublicIp)
		}
		if ipEntry.PrivateIpAddress != nil {
			privateIPAddresses = append(privateIPAddresses, *ipEntry.PrivateIpAddress)
		}
	}

	networkInterfaceData := &model.NetworkInterfaceData{
		Id:                 nic.NetworkInterfaceId,
		MacAddress:         *nic.MacAddress,
		PrivateIpAddresses: &privateIPAddresses,
		PublicIpAddresses:  &publicIPs,
		Network:            *nic.VpcId,
		SubnetId:           *nic.SubnetId,
	}

	return networkInterfaceData
}
