package readers

import (
	logger "agentless/infra/log"
	model "agentless/infra/model/common"
	utils "agentless/infra/utils"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"regexp"
)

type VpcEndpointReader struct {
	Svc     ec2iface.EC2API
	Region  string
	updates chan Resource
}

const (
	dynamodb string = "dynamodb"
	s3       string = "s3"
)

var supportedServices = map[string]*string{
	dynamodb: utils.StrPtr(dynamodb),
	s3:       utils.StrPtr(s3),
}

var serviceTypeToEntityCategory = map[string]*string{
	dynamodb: utils.StrPtr("database"),
	s3:       utils.StrPtr("storage"),
}

var svcNameRegex *regexp.Regexp

func init() {

	var err error
	svcNameRegex, err = regexp.Compile("com\\.amazonaws\\.[^.]+\\.(.+)$")

	if err != nil {
		logger.Log.Fatalf("Failed to compile service name regex", err)
	}
}

func NewVpcEndpointReader(sess *session.Session, region string, resource chan Resource) *VpcEndpointReader {
	return &VpcEndpointReader{
		Svc:     ec2.New(sess),
		Region:  region,
		updates: resource,
	}
}

func (r *VpcEndpointReader) Read() {
	logger.Log.Infof("Reader Started: Type=VpcEndpoint, region=%s", r.Region)
	var err error = nil

	input := &ec2.DescribeVpcEndpointsInput{}
	err = r.Svc.DescribeVpcEndpointsPages(input, func(page *ec2.DescribeVpcEndpointsOutput, lastPage bool) bool {
		for _, instance := range page.VpcEndpoints {
			svcType := getServiceTypeName(instance.ServiceName)
			if svcType == nil {
				continue
			}
			// TODO only get interface for now
			if *instance.VpcEndpointType != "Interface" {
				continue
			}
			item, _ := r.toInventoryItemFromVpcEndpoint(instance)
			item.EntityType = svcType
			item.EntityCategory = serviceTypeToEntityCategory[*svcType]

			r.updates <- Resource{ID: *instance.VpcEndpointId, Region: r.Region, Type: string(model.MS), Item: item}
		}
		return !lastPage
	})

	if err != nil {
		logger.Log.Fatalf("Failed read instances: %s", err)
	}

	logger.Log.Infof("Reader Completed: Type=VpcEndpoint, region=%s", r.Region)
}

func (r *VpcEndpointReader) toInventoryItemFromVpcEndpoint(instance *ec2.VpcEndpoint) (*model.InventoryItem, error) {
	entityData := &model.InventoryItem_EntityData{}
	_ = entityData.FromManagedServiceData(r.toManagedServiceDataFromVpcEndpoint(instance))
	itemType := model.Asset

	item := &model.InventoryItem{
		EntityData: entityData,
		//TODO this or the next line? EntityName:   utils.StrPtr(GetValueOrDefault(AwsTagsToMap(instance.Tags), "Name", *instance.VpcEndpointId)),
		EntityName:  instance.ServiceName,
		ExternalIds: &[]string{*instance.VpcEndpointId},
		ItemId:      instance.VpcEndpointId,
		ItemType:    &itemType,
		Labels:      AwsLabelsListFrom(instance.Tags),
	}

	return item, nil
}

func (r *VpcEndpointReader) toManagedServiceDataFromVpcEndpoint(instance *ec2.VpcEndpoint) model.ManagedServiceData {

	var nics []model.NetworkInterfaceData
	for _, interfaceId := range instance.NetworkInterfaceIds {
		if interfaceId == nil {
			continue
		}
		input := &ec2.DescribeNetworkInterfacesInput{
			NetworkInterfaceIds: []*string{interfaceId}}
		output, err := r.Svc.DescribeNetworkInterfaces(input)
		if err != nil {
			logger.Log.Errorf("failed to describe network interface ID:%s err:%s", *interfaceId, err.Error())
			continue
		}

		// Check if a network interface was found
		if len(output.NetworkInterfaces) == 0 {
			logger.Log.Errorf("No network interface found with ID:%s", *interfaceId)
			continue
		}

		for _, ec2Nic := range output.NetworkInterfaces {
			nicData := toNetworkInterfaceDataFromNetworkInterfaces(ec2Nic)
			nics = append(nics, *nicData)
		}
	}

	return model.ManagedServiceData{
		Type: model.MS,
		Nics: &nics,
	}
}

func toNetworkInterfaceDataFromNetworkInterfaces(nic *ec2.NetworkInterface) *model.NetworkInterfaceData {
	var publicIPs []string
	var privateIPAddresses []string

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

// extract service type from ServiceName using regex
// return nil if the type is not supported
func getServiceTypeName(serviceName *string) *string {

	if serviceName == nil {
		return nil
	}
	svcType := svcNameRegex.FindStringSubmatch(*serviceName)
	return supportedServices[svcType[1]]

	//TODO getting the type using strings package
	//svcTypeIndex := strings.LastIndex(*serviceName, ".") //obtain the extension after the dot
	//svcType := (*serviceName)[svcTypeIndex+1:]
	//return supportedServices[svcType]
}
