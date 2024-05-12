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
	// Mock response
	labels := []string{"label1", "label2", "label3"}
	c.JSON(http.StatusOK, gin.H{"labels": labels})
}

func (s *InventoryService) GetV1ProviderInventory(c *gin.Context, params GetV1ProviderInventoryParams) {
	items := s.listOfInventoryItems()
	c.JSON(http.StatusOK, items)
}

func (s *InventoryService) GetV1ProviderTopology(c *gin.Context, params GetV1ProviderTopologyParams) {
	// Mock response
	topology := map[string]interface{}{
		"nodes": []string{"node1", "node2", "node3"},
		"edges": []string{"edge1", "edge2", "edge3"},
	}
	c.JSON(http.StatusOK, gin.H{"topology": topology})
}

func (s *InventoryService) Update(revision int, fetchedInventory *map[string]readers.Resource) {
	logger.Log.Infof("Cycle %d: Total EC2=%d", revision, len(*fetchedInventory))

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
