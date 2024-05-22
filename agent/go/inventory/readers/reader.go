package readers

import (
	model "agentless/infra/model/common"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Resource struct {
	ID     string
	Region string
	Type   string
	Item   *model.InventoryItem
}

type ResourceReader interface {
	Read()
}

func GetRegionReaders(sess *session.Session, region string, updateChen chan Resource) []ResourceReader {
	return []ResourceReader{
		NewEC2Reader(sess, region, updateChen),
		NewVpcEndpointReader(sess, region, updateChen),
		NewELBReader(sess, region, updateChen),
		NewLambdaReader(sess, region, updateChen),
	}
}
