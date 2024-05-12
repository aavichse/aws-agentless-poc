// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/oapi-codegen/runtime"
)

// Defines values for IOTType.
const (
	IOTTypeIOT IOTType = "IOT"
)

// Defines values for InventoryItemItemType.
const (
	Asset InventoryItemItemType = "asset"
)

// Defines values for ManagedServiceDataType.
const (
	MS ManagedServiceDataType = "MS"
)

// Defines values for OnlineStatus.
const (
	Offline OnlineStatus = "offline"
	Online  OnlineStatus = "online"
)

// Defines values for PowerState.
const (
	Restarting PowerState = "restarting"
	Running    PowerState = "running"
	Stopped    PowerState = "stopped"
)

// Defines values for RegistrationStatus.
const (
	Registered   RegistrationStatus = "registered"
	Unregistered RegistrationStatus = "unregistered"
)

// Defines values for VMDataType.
const (
	VM VMDataType = "VM"
)

// AzureNetworkAccess defines model for AzureNetworkAccess.
type AzureNetworkAccess struct {
	PrivateAccessEnabled   *bool `json:"private-access-enabled,omitempty"`
	PublicAccessEnabled    *bool `json:"public-access-enabled,omitempty"`
	ServiceEndpointEnabled *bool `json:"service-endpoint-enabled,omitempty"`
}

// AzureNetworkTopology defines model for AzureNetworkTopology.
type AzureNetworkTopology struct {
	Region        *string `json:"region,omitempty"`
	ResourceGroup *string `json:"resource-group,omitempty"`
	ResourceId    *string `json:"resource-id,omitempty"`
	Subscription  *string `json:"subscription,omitempty"`
}

// HardwareInfo defines model for HardwareInfo.
type HardwareInfo struct {
	HwClass        *string `json:"hw-class,omitempty"`
	HwManufacturer *string `json:"hw-manufacturer,omitempty"`
	HwType         *string `json:"hw-type,omitempty"`
}

// IOT defines model for IOT.
type IOT struct {
	DeviceScore        *string                 `json:"device-score,omitempty"`
	HwDetails          HardwareInfo            `json:"hw-details"`
	Nics               *[]NetworkInterfaceData `json:"nics,omitempty"`
	OnlineStatus       *OnlineStatus           `json:"online-status,omitempty"`
	OsDetails          *map[string]interface{} `json:"os-details,omitempty"`
	PowerState         *PowerState             `json:"power-state,omitempty"`
	RegistrationStatus *IOTRegistrationStatus  `json:"registration-status,omitempty"`
	Type               IOTType                 `json:"type"`
}

// IOTType defines model for IOT.Type.
type IOTType string

// IOTRegistrationStatus defines model for IOTRegistrationStatus.
type IOTRegistrationStatus struct {
	RegDate   *time.Time          `json:"reg-date,omitempty"`
	RegStatus *RegistrationStatus `json:"reg-status,omitempty"`
	UnregDate *time.Time          `json:"unreg-date,omitempty"`
}

// InventoryItem defines model for InventoryItem.
type InventoryItem struct {
	// EntityCategory asset category such as Compute or Database
	EntityCategory *string `json:"entity-category,omitempty"`

	// EntityData additional entity data depending on its type
	EntityData *InventoryItem_EntityData `json:"entity-data,omitempty"`

	// EntityName asset name
	EntityName *string `json:"entity-name,omitempty"`

	// EntityType asset type such as Virtual Machine or Azure SQL Server
	EntityType *string `json:"entity-type,omitempty"`

	// ExternalIds list of external IDs that are managed by external asset management systems such as Azure resource management
	ExternalIds *[]string `json:"external-ids,omitempty"`

	// ItemId inventory item id assigned by Centra
	ItemId *string `json:"item-id,omitempty"`

	// ItemType inventory item type - currently asset only
	ItemType *InventoryItemItemType `json:"item-type,omitempty"`
	Labels   *[]Label               `json:"labels,omitempty"`
	union    json.RawMessage
}

// InventoryItemEntityData3 defines model for .
type InventoryItemEntityData3 = map[string]interface{}

// InventoryItem_EntityData additional entity data depending on its type
type InventoryItem_EntityData struct {
	union json.RawMessage
}

// InventoryItemItemType inventory item type - currently asset only
type InventoryItemItemType string

// InventoryItem0 defines model for .
type InventoryItem0 = interface{}

// InventoryItem1 defines model for .
type InventoryItem1 = interface{}

// Label defines model for Label.
type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ManagedServiceData defines model for ManagedServiceData.
type ManagedServiceData struct {
	NetworkAccess   *ManagedServiceData_NetworkAccess   `json:"network-access,omitempty"`
	NetworkTopology *ManagedServiceData_NetworkTopology `json:"network-topology,omitempty"`
	Nics            *[]NetworkInterfaceData             `json:"nics,omitempty"`
	Type            ManagedServiceDataType              `json:"type"`
}

// ManagedServiceDataNetworkAccess1 defines model for .
type ManagedServiceDataNetworkAccess1 = map[string]interface{}

// ManagedServiceData_NetworkAccess defines model for ManagedServiceData.NetworkAccess.
type ManagedServiceData_NetworkAccess struct {
	union json.RawMessage
}

// ManagedServiceDataNetworkTopology1 defines model for .
type ManagedServiceDataNetworkTopology1 = map[string]interface{}

// ManagedServiceData_NetworkTopology defines model for ManagedServiceData.NetworkTopology.
type ManagedServiceData_NetworkTopology struct {
	union json.RawMessage
}

// ManagedServiceDataType defines model for ManagedServiceData.Type.
type ManagedServiceDataType string

// NetworkInterfaceData defines model for NetworkInterfaceData.
type NetworkInterfaceData struct {
	// Id nic id, managed as resource id in Azure
	Id                 *string   `json:"id,omitempty"`
	MacAddress         string    `json:"mac-address"`
	Network            string    `json:"network"`
	PrivateIpAddresses *[]string `json:"private-ip-addresses,omitempty"`
	PublicIpAddresses  *[]string `json:"public-ip-addresses,omitempty"`
	SubnetId           string    `json:"subnet-id"`
}

// OnlineStatus defines model for OnlineStatus.
type OnlineStatus string

// PowerState defines model for PowerState.
type PowerState string

// RegistrationStatus defines model for RegistrationStatus.
type RegistrationStatus string

// VMData defines model for VMData.
type VMData struct {
	HwDetails       *HardwareInfo           `json:"hw-details,omitempty"`
	NetworkTopology *VMData_NetworkTopology `json:"network-topology,omitempty"`
	Nics            *[]NetworkInterfaceData `json:"nics,omitempty"`
	OsDetails       *map[string]interface{} `json:"os-details,omitempty"`
	PowerState      *PowerState             `json:"power-state,omitempty"`
	Type            VMDataType              `json:"type"`
}

// VMDataNetworkTopology1 defines model for .
type VMDataNetworkTopology1 = map[string]interface{}

// VMData_NetworkTopology defines model for VMData.NetworkTopology.
type VMData_NetworkTopology struct {
	union json.RawMessage
}

// VMDataType defines model for VMData.Type.
type VMDataType string

// AsInventoryItem0 returns the union data inside the InventoryItem as a InventoryItem0
func (t InventoryItem) AsInventoryItem0() (InventoryItem0, error) {
	var body InventoryItem0
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromInventoryItem0 overwrites any union data inside the InventoryItem as the provided InventoryItem0
func (t *InventoryItem) FromInventoryItem0(v InventoryItem0) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeInventoryItem0 performs a merge with any union data inside the InventoryItem, using the provided InventoryItem0
func (t *InventoryItem) MergeInventoryItem0(v InventoryItem0) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsInventoryItem1 returns the union data inside the InventoryItem as a InventoryItem1
func (t InventoryItem) AsInventoryItem1() (InventoryItem1, error) {
	var body InventoryItem1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromInventoryItem1 overwrites any union data inside the InventoryItem as the provided InventoryItem1
func (t *InventoryItem) FromInventoryItem1(v InventoryItem1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeInventoryItem1 performs a merge with any union data inside the InventoryItem, using the provided InventoryItem1
func (t *InventoryItem) MergeInventoryItem1(v InventoryItem1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t InventoryItem) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	if err != nil {
		return nil, err
	}
	object := make(map[string]json.RawMessage)
	if t.union != nil {
		err = json.Unmarshal(b, &object)
		if err != nil {
			return nil, err
		}
	}

	if t.EntityCategory != nil {
		object["entity-category"], err = json.Marshal(t.EntityCategory)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'entity-category': %w", err)
		}
	}

	if t.EntityData != nil {
		object["entity-data"], err = json.Marshal(t.EntityData)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'entity-data': %w", err)
		}
	}

	if t.EntityName != nil {
		object["entity-name"], err = json.Marshal(t.EntityName)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'entity-name': %w", err)
		}
	}

	if t.EntityType != nil {
		object["entity-type"], err = json.Marshal(t.EntityType)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'entity-type': %w", err)
		}
	}

	if t.ExternalIds != nil {
		object["external-ids"], err = json.Marshal(t.ExternalIds)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'external-ids': %w", err)
		}
	}

	if t.ItemId != nil {
		object["item-id"], err = json.Marshal(t.ItemId)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'item-id': %w", err)
		}
	}

	if t.ItemType != nil {
		object["item-type"], err = json.Marshal(t.ItemType)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'item-type': %w", err)
		}
	}

	if t.Labels != nil {
		object["labels"], err = json.Marshal(t.Labels)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'labels': %w", err)
		}
	}
	b, err = json.Marshal(object)
	return b, err
}

func (t *InventoryItem) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	object := make(map[string]json.RawMessage)
	err = json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["entity-category"]; found {
		err = json.Unmarshal(raw, &t.EntityCategory)
		if err != nil {
			return fmt.Errorf("error reading 'entity-category': %w", err)
		}
	}

	if raw, found := object["entity-data"]; found {
		err = json.Unmarshal(raw, &t.EntityData)
		if err != nil {
			return fmt.Errorf("error reading 'entity-data': %w", err)
		}
	}

	if raw, found := object["entity-name"]; found {
		err = json.Unmarshal(raw, &t.EntityName)
		if err != nil {
			return fmt.Errorf("error reading 'entity-name': %w", err)
		}
	}

	if raw, found := object["entity-type"]; found {
		err = json.Unmarshal(raw, &t.EntityType)
		if err != nil {
			return fmt.Errorf("error reading 'entity-type': %w", err)
		}
	}

	if raw, found := object["external-ids"]; found {
		err = json.Unmarshal(raw, &t.ExternalIds)
		if err != nil {
			return fmt.Errorf("error reading 'external-ids': %w", err)
		}
	}

	if raw, found := object["item-id"]; found {
		err = json.Unmarshal(raw, &t.ItemId)
		if err != nil {
			return fmt.Errorf("error reading 'item-id': %w", err)
		}
	}

	if raw, found := object["item-type"]; found {
		err = json.Unmarshal(raw, &t.ItemType)
		if err != nil {
			return fmt.Errorf("error reading 'item-type': %w", err)
		}
	}

	if raw, found := object["labels"]; found {
		err = json.Unmarshal(raw, &t.Labels)
		if err != nil {
			return fmt.Errorf("error reading 'labels': %w", err)
		}
	}

	return err
}

// AsVMData returns the union data inside the InventoryItem_EntityData as a VMData
func (t InventoryItem_EntityData) AsVMData() (VMData, error) {
	var body VMData
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromVMData overwrites any union data inside the InventoryItem_EntityData as the provided VMData
func (t *InventoryItem_EntityData) FromVMData(v VMData) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeVMData performs a merge with any union data inside the InventoryItem_EntityData, using the provided VMData
func (t *InventoryItem_EntityData) MergeVMData(v VMData) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsIOT returns the union data inside the InventoryItem_EntityData as a IOT
func (t InventoryItem_EntityData) AsIOT() (IOT, error) {
	var body IOT
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromIOT overwrites any union data inside the InventoryItem_EntityData as the provided IOT
func (t *InventoryItem_EntityData) FromIOT(v IOT) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeIOT performs a merge with any union data inside the InventoryItem_EntityData, using the provided IOT
func (t *InventoryItem_EntityData) MergeIOT(v IOT) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsManagedServiceData returns the union data inside the InventoryItem_EntityData as a ManagedServiceData
func (t InventoryItem_EntityData) AsManagedServiceData() (ManagedServiceData, error) {
	var body ManagedServiceData
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromManagedServiceData overwrites any union data inside the InventoryItem_EntityData as the provided ManagedServiceData
func (t *InventoryItem_EntityData) FromManagedServiceData(v ManagedServiceData) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeManagedServiceData performs a merge with any union data inside the InventoryItem_EntityData, using the provided ManagedServiceData
func (t *InventoryItem_EntityData) MergeManagedServiceData(v ManagedServiceData) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsInventoryItemEntityData3 returns the union data inside the InventoryItem_EntityData as a InventoryItemEntityData3
func (t InventoryItem_EntityData) AsInventoryItemEntityData3() (InventoryItemEntityData3, error) {
	var body InventoryItemEntityData3
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromInventoryItemEntityData3 overwrites any union data inside the InventoryItem_EntityData as the provided InventoryItemEntityData3
func (t *InventoryItem_EntityData) FromInventoryItemEntityData3(v InventoryItemEntityData3) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeInventoryItemEntityData3 performs a merge with any union data inside the InventoryItem_EntityData, using the provided InventoryItemEntityData3
func (t *InventoryItem_EntityData) MergeInventoryItemEntityData3(v InventoryItemEntityData3) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t InventoryItem_EntityData) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *InventoryItem_EntityData) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsAzureNetworkAccess returns the union data inside the ManagedServiceData_NetworkAccess as a AzureNetworkAccess
func (t ManagedServiceData_NetworkAccess) AsAzureNetworkAccess() (AzureNetworkAccess, error) {
	var body AzureNetworkAccess
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAzureNetworkAccess overwrites any union data inside the ManagedServiceData_NetworkAccess as the provided AzureNetworkAccess
func (t *ManagedServiceData_NetworkAccess) FromAzureNetworkAccess(v AzureNetworkAccess) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAzureNetworkAccess performs a merge with any union data inside the ManagedServiceData_NetworkAccess, using the provided AzureNetworkAccess
func (t *ManagedServiceData_NetworkAccess) MergeAzureNetworkAccess(v AzureNetworkAccess) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsManagedServiceDataNetworkAccess1 returns the union data inside the ManagedServiceData_NetworkAccess as a ManagedServiceDataNetworkAccess1
func (t ManagedServiceData_NetworkAccess) AsManagedServiceDataNetworkAccess1() (ManagedServiceDataNetworkAccess1, error) {
	var body ManagedServiceDataNetworkAccess1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromManagedServiceDataNetworkAccess1 overwrites any union data inside the ManagedServiceData_NetworkAccess as the provided ManagedServiceDataNetworkAccess1
func (t *ManagedServiceData_NetworkAccess) FromManagedServiceDataNetworkAccess1(v ManagedServiceDataNetworkAccess1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeManagedServiceDataNetworkAccess1 performs a merge with any union data inside the ManagedServiceData_NetworkAccess, using the provided ManagedServiceDataNetworkAccess1
func (t *ManagedServiceData_NetworkAccess) MergeManagedServiceDataNetworkAccess1(v ManagedServiceDataNetworkAccess1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t ManagedServiceData_NetworkAccess) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *ManagedServiceData_NetworkAccess) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsAzureNetworkTopology returns the union data inside the ManagedServiceData_NetworkTopology as a AzureNetworkTopology
func (t ManagedServiceData_NetworkTopology) AsAzureNetworkTopology() (AzureNetworkTopology, error) {
	var body AzureNetworkTopology
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAzureNetworkTopology overwrites any union data inside the ManagedServiceData_NetworkTopology as the provided AzureNetworkTopology
func (t *ManagedServiceData_NetworkTopology) FromAzureNetworkTopology(v AzureNetworkTopology) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAzureNetworkTopology performs a merge with any union data inside the ManagedServiceData_NetworkTopology, using the provided AzureNetworkTopology
func (t *ManagedServiceData_NetworkTopology) MergeAzureNetworkTopology(v AzureNetworkTopology) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsManagedServiceDataNetworkTopology1 returns the union data inside the ManagedServiceData_NetworkTopology as a ManagedServiceDataNetworkTopology1
func (t ManagedServiceData_NetworkTopology) AsManagedServiceDataNetworkTopology1() (ManagedServiceDataNetworkTopology1, error) {
	var body ManagedServiceDataNetworkTopology1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromManagedServiceDataNetworkTopology1 overwrites any union data inside the ManagedServiceData_NetworkTopology as the provided ManagedServiceDataNetworkTopology1
func (t *ManagedServiceData_NetworkTopology) FromManagedServiceDataNetworkTopology1(v ManagedServiceDataNetworkTopology1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeManagedServiceDataNetworkTopology1 performs a merge with any union data inside the ManagedServiceData_NetworkTopology, using the provided ManagedServiceDataNetworkTopology1
func (t *ManagedServiceData_NetworkTopology) MergeManagedServiceDataNetworkTopology1(v ManagedServiceDataNetworkTopology1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t ManagedServiceData_NetworkTopology) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *ManagedServiceData_NetworkTopology) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsAzureNetworkTopology returns the union data inside the VMData_NetworkTopology as a AzureNetworkTopology
func (t VMData_NetworkTopology) AsAzureNetworkTopology() (AzureNetworkTopology, error) {
	var body AzureNetworkTopology
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAzureNetworkTopology overwrites any union data inside the VMData_NetworkTopology as the provided AzureNetworkTopology
func (t *VMData_NetworkTopology) FromAzureNetworkTopology(v AzureNetworkTopology) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAzureNetworkTopology performs a merge with any union data inside the VMData_NetworkTopology, using the provided AzureNetworkTopology
func (t *VMData_NetworkTopology) MergeAzureNetworkTopology(v AzureNetworkTopology) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsVMDataNetworkTopology1 returns the union data inside the VMData_NetworkTopology as a VMDataNetworkTopology1
func (t VMData_NetworkTopology) AsVMDataNetworkTopology1() (VMDataNetworkTopology1, error) {
	var body VMDataNetworkTopology1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromVMDataNetworkTopology1 overwrites any union data inside the VMData_NetworkTopology as the provided VMDataNetworkTopology1
func (t *VMData_NetworkTopology) FromVMDataNetworkTopology1(v VMDataNetworkTopology1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeVMDataNetworkTopology1 performs a merge with any union data inside the VMData_NetworkTopology, using the provided VMDataNetworkTopology1
func (t *VMData_NetworkTopology) MergeVMDataNetworkTopology1(v VMDataNetworkTopology1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t VMData_NetworkTopology) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *VMData_NetworkTopology) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}