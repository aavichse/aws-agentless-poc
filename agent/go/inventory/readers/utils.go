package readers

import (
	contracts "agentless/infra/model/common"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func AwsLabelFrom(tag *ec2.Tag) *contracts.Label {
	return &contracts.Label{
		Key:   *tag.Key,
		Value: *tag.Value,
	}
}

func AwsLabelsListFrom(tags []*ec2.Tag) *[]contracts.Label {
	labels := make([]contracts.Label, 0, len(tags))
	for _, tag := range tags {
		labels = append(labels, *AwsLabelFrom(tag))
	}
	return &labels
}

func AwsTagsToMap(tags []*ec2.Tag) *map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[*tag.Key] = *tag.Value
	}
	return &tagMap
}

func GetValueOrDefault(tags *map[string]string, key, defaultValue string) string {
	if value, ok := (*tags)[key]; ok {
		return value
	}
	return defaultValue
}

func AddDefaultLabels(item Resource) {

	if item.Item.Labels == nil {
		item.Item.Labels = &[]contracts.Label{}
	}
	*item.Item.Labels = append(*item.Item.Labels, contracts.Label{Key: "Region", Value: item.Region})
}

func ToManagedServiceDataFromNIC(ec2Nics []*ec2.NetworkInterface) *contracts.ManagedServiceData {

	nics := &[]contracts.NetworkInterfaceData{}

	for _, ec2Nic := range ec2Nics {
		nicData := toNetworkInterfaceDataFromNetworkInterfaces(ec2Nic)
		*nics = append(*nics, *nicData)
	}

	return &contracts.ManagedServiceData{
		Type: contracts.MS,
		Nics: nics,
	}
}
