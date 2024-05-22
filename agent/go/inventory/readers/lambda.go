package readers

import (
	logger "agentless/infra/log"
	model "agentless/infra/model/common"
	utils "agentless/infra/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/lambda"
)

const (
	LambdaSvcType string = "lambda"
)

type LambdaReader struct {
	LambdaSvc *lambda.Lambda
	EC2Svc    ec2iface.EC2API
	Region    string
	updates   chan Resource
}

func NewLambdaReader(sess *session.Session, region string, resource chan Resource) *LambdaReader {
	return &LambdaReader{
		LambdaSvc: lambda.New(sess),
		EC2Svc:    ec2.New(sess),
		Region:    region,
		updates:   resource,
	}
}

func (r *LambdaReader) Read() {
	
	logger.Log.Infof("Reader Started: Type=Lambda, region=%s", r.Region)
	
	err := r.LambdaSvc.ListFunctionsPages(
		&lambda.ListFunctionsInput{},
		func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
			for _, instance := range page.Functions {
				item, err := r.toInventoryItemFromLambda(instance)
				if err != nil {
					logger.Log.Errorf("failed discover lambda %s", *instance.FunctionArn)
					continue
				}
				r.updates <- Resource{ID: *instance.FunctionArn, Region: r.Region, Type: string(model.MS), Item: item}
			}
			
			return !lastPage
		})
	if err != nil {
		logger.Log.Errorf("Failed read instances: %s", err)
	}
	
	logger.Log.Infof("Reader Completed: Type=Lambda, region=%s", r.Region)
}

func (r *LambdaReader) toInventoryItemFromLambda(instance *lambda.FunctionConfiguration) (*model.InventoryItem, error) {
	
	entityData := &model.InventoryItem_EntityData{}
	if instance.VpcConfig != nil {
		nic, err := r.toManagedServiceDataFromLambda(instance)
		if err != nil {
			return nil, err
		}
		err = entityData.FromManagedServiceData(*nic)
		if err != nil {
			return nil, err
		}
	}
	
	tagResult, err := r.LambdaSvc.ListTags(&lambda.ListTagsInput{Resource: instance.FunctionArn})
	if err != nil {
		logger.Log.Errorf("failed to describe tags for lambda %s, %v", *instance.FunctionName, err)
	}
	itemType := model.Asset
	
	item := &model.InventoryItem{
		EntityCategory: utils.StrPtr("compute"),
		EntityData:     entityData,
		EntityName:     instance.FunctionName,
		EntityType:     utils.StrPtr(LambdaSvcType),
		ExternalIds:    &[]string{*instance.FunctionArn},
		ItemId:         instance.FunctionArn,
		ItemType:       &itemType,
		Labels:         awsMapTagsToList(tagResult.Tags),
	}
	
	return item, nil
}

func awsMapTagsToList(tags map[string]*string) *[]model.Label {
	
	var list []model.Label
	for k, v := range tags {
		list = append(list,
			model.Label{
				Key:   k,
				Value: *v,
			})
	}
	return &list
}

func (r *LambdaReader) toManagedServiceDataFromLambda(result *lambda.FunctionConfiguration) (*model.ManagedServiceData,
	error) {
	
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("subnet-id"),
				Values: result.VpcConfig.SubnetIds,
			},
			{
				Name:   aws.String("group-id"),
				Values: result.VpcConfig.SecurityGroupIds,
			},
			{
				// example from aws interface:
				// Description: "AWS Lambda VPC ENI-poc-examples-lmabda-generate-trade-data-52370f8a-99ef-4b36-b14a-068b5e766ad2"
				Name:   aws.String("description"),
				Values: []*string{aws.String("*-" + *result.FunctionName + "-*")},
			},
		},
	}
	
	output, err := r.EC2Svc.DescribeNetworkInterfaces(input)
	if err != nil {
		logger.Log.Errorf("failed to describe network interfaces, %v", err)
		return nil, err
	}
	
	return ToManagedServiceDataFromNIC(output.NetworkInterfaces), nil
}
