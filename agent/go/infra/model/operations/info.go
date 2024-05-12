// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package models

import (
	"encoding/json"

	externalRef0 "common/common.yaml"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// EnvUnit describing a logical unit of the integration environment. e.g. azure subscription, gcp project etc
type EnvUnit struct {
	Id   openapi_types.UUID `json:"id"`
	Name string             `json:"name"`

	// Parent the location in the organization Hierarchy specified by the Hierarchy id. can be empty string if integration does not have a Hierarchy structure.
	Parent string `json:"parent"`
}

// EnvUnits defines model for EnvUnits.
type EnvUnits = []EnvUnit

// GeneralInfo defines model for GeneralInfo.
type GeneralInfo struct {
	General *map[string]interface{} `json:"general,omitempty"`
}

// Hierarchy defines model for Hierarchy.
type Hierarchy struct {
	Children *[]Hierarchy `json:"children,omitempty"`
	Id       string       `json:"id"`
	Name     string       `json:"name"`
}

// IntegrationEnvInfo defines model for IntegrationEnvInfo.
type IntegrationEnvInfo struct {
	ConnectorId   openapi_types.UUID      `json:"connector-id"`
	ConnectorType externalRef0.Connector  `json:"connector-type"`
	Info          IntegrationEnvInfo_Info `json:"info"`
}

// IntegrationEnvInfo_Info defines model for IntegrationEnvInfo.Info.
type IntegrationEnvInfo_Info struct {
	union json.RawMessage
}

// OrganizationTreeInfo defines model for OrganizationTreeInfo.
type OrganizationTreeInfo struct {
	OrgTree Hierarchy `json:"org-tree"`
}

// GetV1OperationsEnvUnitsListParams defines parameters for GetV1OperationsEnvUnitsList.
type GetV1OperationsEnvUnitsListParams struct {
	// Cursor Cursor is the starting position in the result set.
	Cursor *int `form:"cursor,omitempty" json:"cursor,omitempty"`

	// PageSize Maximum number of unit items per page
	PageSize *int `form:"page_size,omitempty" json:"page_size,omitempty"`
}

// AsOrganizationTreeInfo returns the union data inside the IntegrationEnvInfo_Info as a OrganizationTreeInfo
func (t IntegrationEnvInfo_Info) AsOrganizationTreeInfo() (OrganizationTreeInfo, error) {
	var body OrganizationTreeInfo
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromOrganizationTreeInfo overwrites any union data inside the IntegrationEnvInfo_Info as the provided OrganizationTreeInfo
func (t *IntegrationEnvInfo_Info) FromOrganizationTreeInfo(v OrganizationTreeInfo) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeOrganizationTreeInfo performs a merge with any union data inside the IntegrationEnvInfo_Info, using the provided OrganizationTreeInfo
func (t *IntegrationEnvInfo_Info) MergeOrganizationTreeInfo(v OrganizationTreeInfo) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsGeneralInfo returns the union data inside the IntegrationEnvInfo_Info as a GeneralInfo
func (t IntegrationEnvInfo_Info) AsGeneralInfo() (GeneralInfo, error) {
	var body GeneralInfo
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromGeneralInfo overwrites any union data inside the IntegrationEnvInfo_Info as the provided GeneralInfo
func (t *IntegrationEnvInfo_Info) FromGeneralInfo(v GeneralInfo) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeGeneralInfo performs a merge with any union data inside the IntegrationEnvInfo_Info, using the provided GeneralInfo
func (t *IntegrationEnvInfo_Info) MergeGeneralInfo(v GeneralInfo) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t IntegrationEnvInfo_Info) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *IntegrationEnvInfo_Info) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}