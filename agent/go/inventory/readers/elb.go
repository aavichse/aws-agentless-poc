package readers

import (
	logger "agentless/infra/log"
	model "agentless/infra/model/common"
	utils "agentless/infra/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

const (
	LBSvcType string = "load balancer"
)

type ELBReader struct {
	ELBV2Svc *elbv2.ELBV2 // for Application Load Balancers (ALBs) and Network Load Balancers (NLBs)
	//TODO use when supporting classic ELBSvc   *elb.ELB      for Classic Load Balancers (CLBs)
	EC2Svc  ec2iface.EC2API
	Region  string
	updates chan Resource
}

func NewELBReader(sess *session.Session, region string, resource chan Resource) *ELBReader {
	return &ELBReader{
		ELBV2Svc: elbv2.New(sess),
		//ELBSvc:  elb.New(sess),  TODO use when supporting classic
		EC2Svc:  ec2.New(sess),
		Region:  region,
		updates: resource,
	}
}

func (r *ELBReader) Read() {
	logger.Log.Infof("Reader Started: Type=ELB, region=%s", r.Region)
	var err error = nil

	// TODO use when supporting classic
	//err = r.ELBSvc.DescribeLoadBalancersPages(
	//	&elb.DescribeLoadBalancersInput{},
	//	func(page *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
	//		for _, instance := range page.LoadBalancerDescriptions {
	//			item, _ := r.toInventoryItemFromELB(instance)
	//			r.updates <- Resource{ID: *instance.DNSName /*?*/, Region: r.Region, Type: string(model.MS), Item: item}
	//		}
	//
	//		return !lastPage
	//	})
	//if err != nil {
	//	logger.Log.Fatalf("Failed read instances: %s", err)
	//}

	err = r.ELBV2Svc.DescribeLoadBalancersPages(
		&elbv2.DescribeLoadBalancersInput{},
		func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, instance := range page.LoadBalancers {
				item, _ := r.toInventoryItemFromELBV2(instance)
				r.updates <- Resource{ID: *instance.LoadBalancerArn, Region: r.Region, Type: string(model.MS), Item: item}
			}

			return !lastPage
		})
	if err != nil {
		logger.Log.Errorf("Failed read instances: %s", err)
	}

	logger.Log.Infof("Reader Completed: Type=ELB, region=%s", r.Region)
}

func (r *ELBReader) toInventoryItemFromELBV2(instance *elbv2.LoadBalancer) (*model.InventoryItem, error) {

	entityData := &model.InventoryItem_EntityData{}
	_ = entityData.FromManagedServiceData(r.toManagedServiceDataFromElb(instance))

	tagResult, err := r.ELBV2Svc.DescribeTags(&elbv2.DescribeTagsInput{
		ResourceArns: []*string{instance.LoadBalancerArn},
	})
	if err != nil {
		logger.Log.Errorf("failed to describe tags for load balancer %s, %v", *instance.LoadBalancerName, err)
	}
	tags := tagResult.TagDescriptions[0].Tags
	itemType := model.Asset

	item := &model.InventoryItem{
		EntityCategory: utils.StrPtr("compute"),
		EntityData:     entityData,
		EntityName:     instance.LoadBalancerName,
		EntityType:     utils.StrPtr(LBSvcType),
		ExternalIds:    &[]string{*instance.LoadBalancerArn, *instance.DNSName},
		ItemId:         instance.LoadBalancerArn,
		ItemType:       &itemType,
		Labels:         awsELBTagsToList(tags),
	}

	return item, nil
}

func awsELBTagsToList(tags []*elbv2.Tag) *[]model.Label {
	labels := make([]model.Label, 0, len(tags))
	for _, tag := range tags {
		labels = append(labels, *awsELBTagToLabel(tag))
	}
	return &labels
}

func awsELBTagToLabel(tag *elbv2.Tag) *model.Label {
	return &model.Label{
		Key:   *tag.Key,
		Value: *tag.Value,
	}
}

func (r *ELBReader) toManagedServiceDataFromElb(elb *elbv2.LoadBalancer) model.ManagedServiceData {

	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				// FIXME: Better match pattern considering '/' 
				// example from aws interface: Description: "ELB app/aws-poc-elb-example/25130940da3ebdcd",
				// the elb name is unique per region.
				Name:   aws.String("description"),
				Values: []*string{aws.String("*" + *elb.LoadBalancerName + "*")},
			},
			{
				// When the requester is an AWS service, such as the Elastic Load Balancing service,
				// the requester-id will be a special alias
				Name:   aws.String("requester-id"),
				Values: []*string{aws.String("amazon-elb")},
			},
		}}
	output, err := r.EC2Svc.DescribeNetworkInterfaces(input)
	if err != nil {
		logger.Log.Errorf("failed to describe network interfaces, %v", err)
	}

	return ToManagedServiceDataFromNIC(output.NetworkInterfaces)
}

// TODO use when supporting classic - check the InventoryItem fields are correct
//func (r *ELBReader) toInventoryItemFromELB(instance *elb.LoadBalancerDescription) (*model.InventoryItem, error) {
//	entityData := &model.InventoryItem_EntityData{}
//	_ = entityData.FromManagedServiceData(r.toManagedServiceDataFromELB(instance))
//	tagResult, err := r.ELBSvc.DescribeTags(&elb.DescribeTagsInput{
//		LoadBalancerNames: []*string{instance.LoadBalancerName},
//	})
//	if err != nil {
//		logger.Log.Fatalf("failed to describe tags for load balancer %s, %v", *instance.LoadBalancerName, err)
//	}
//	tags := tagResult.TagDescriptions[0].Tags
//
//	item := &model.InventoryItem{
//		EntityCategory: utils.StrPtr("compute"),
//		EntityData:     entityData,
//		EntityName:     utils.StrPtr(GetValueOrDefault(awsELBTagsToMap(tags), "Name", *instance.LoadBalancerName)),
//		EntityType:     utils.StrPtr(LBSvcType),
//		ExternalIds:    utils.SlicePtr([]string{*instance.DNSName}),
//		ItemId:         instance.DNSName,
//		ItemType:       (*model.InventoryItemItemType)(utils.StrPtr(string(model.Asset))),
//		Labels:         awsELBTagsToList(tags),
//	}
//
//	return item, nil
//}

// TODO use when supporting classic
//func (r *ELBReader) toManagedServiceDataFromELB(elb *elb.LoadBalancerDescription) model.ManagedServiceData {
//
//	var instanceIDs []*string
//	for _, instance := range elb.Instances {
//		instanceIDs = append(instanceIDs, instance.InstanceId)
//	}
//	ec2Result, err := r.EC2Svs.DescribeInstances(&ec2.DescribeInstancesInput{
//		InstanceIds: instanceIDs,
//	})
//	if err != nil {
//		logger.Log.Errorf("failed to describe instances, %v", err)
//	}
//
//	nics := &[]model.NetworkInterfaceData{}
//
//	for _, reservation := range ec2Result.Reservations {
//		for _, ec2Instance := range reservation.Instances {
//			for _, ec2Nic := range ec2Instance.NetworkInterfaces {
//				nicData := ToNetwotkInterfaceDataFrom(ec2Nic)
//				*nics = append(*nics, *nicData)
//			}
//		}
//	}
//
//	return model.ManagedServiceData{
//		Type: model.MS,
//		Nics: nics,
//	}
//}

// TODO use when supporting classic
//func awsELBTagsToMap(tags []*elb.Tag) *map[string]string {
//	tagMap := make(map[string]string)
//	for _, tag := range tags {
//		tagMap[*tag.Key] = *tag.Value
//	}
//	return &tagMap
//}

// TODO use when supporting classic
//func awsELBTagsToList(tags []*elb.Tag) *[]contracts.Label {
//	labels := make([]contracts.Label, 0, len(tags))
//	for _, tag := range tags {
//		labels = append(labels, *awsELBTagToLabel(tag))
//	}
//	return &labels
//}

// TODO use when supporting classic
//func awsELBTagToLabel(tag *elb.Tag) *contracts.Label {
//	return &contracts.Label{
//		Key:   *tag.Key,
//		Value: *tag.Value,
//	}
//}
