package readers

import model "agentless/infra/model/common"

type Resource struct {
	ID     string
	Region string
	Type   string
	Item   *model.InventoryItem
}

type ResourceReader interface {
	Read()
}
