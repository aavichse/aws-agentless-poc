package main

import (
	logger "agentless/infra/log"
	model "agentless/infra/model/common"
	"agentless/inventory/readers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	Update(revision int, fetchedInventory *map[string]readers.Resource)
}

type InventoryService struct {
	inventory *map[string]readers.Resource
	revision  int
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		inventory: nil,
	}
}

func (s *InventoryService) GetV1ProviderImportLabels(c *gin.Context, params GetV1ProviderImportLabelsParams) {
	c.JSON(http.StatusOK, gin.H{})
}

func (s *InventoryService) GetV1ProviderInventory(c *gin.Context, params GetV1ProviderInventoryParams) {
	items := s.listOfInventoryItems()
	logger.Log.Infof("GET %s, Total items: %d", c.Request.URL.String(), len(items))
	c.JSON(http.StatusOK, items)
}

func (s *InventoryService) GetV1ProviderTopology(c *gin.Context, params GetV1ProviderTopologyParams) {
	// Mock response
	topology := map[string]interface{}{}
	c.JSON(http.StatusOK, gin.H{"topology": topology})
}

func (s *InventoryService) Update(revision int, fetchedInventory *map[string]readers.Resource) {
	logger.Log.Infof("Cycle %d: Total Inventory=%d", revision, len(*fetchedInventory))
	s.inventory = fetchedInventory
	s.revision = revision
}

func (s *InventoryService) listOfInventoryItems() []*model.InventoryItem {
	items := make([]*model.InventoryItem, 0, len(*s.inventory))
	for _, v := range *s.inventory {
		items = append(items, v.Item)
	}
	return items
}
