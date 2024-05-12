// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package models

import (
	externalRef0 "common/common.yaml"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Defines values for Status.
const (
	Connecting         Status = "connecting"
	Down               Status = "down"
	Error              Status = "error"
	Initializing       Status = "initializing"
	Missing            Status = "missing"
	PartiallyUp        Status = "partially_up"
	Stopped            Status = "stopped"
	Up                 Status = "up"
	VerificationFailed Status = "verification_failed"
	Verifying          Status = "verifying"
)

// ApplicationStatus defines model for ApplicationStatus.
type ApplicationStatus struct {
	ApplicationStatus Status                    `json:"application-status"`
	ServicesStatus    *map[string]ServiceStatus `json:"services-status,omitempty"`
}

// ComponentDetails defines model for ComponentDetails.
type ComponentDetails struct {
	ComponentVersion    *string   `json:"component-version,omitempty"`
	ConfigRevision      *int      `json:"config-revision,omitempty"`
	DcInventoryRevision int       `json:"dc-inventory-revision"`
	Hostname            *string   `json:"hostname,omitempty"`
	PolicyRevision      int       `json:"policy-revision"`
	PrivateIpAddresses  *[]string `json:"private-ip-addresses,omitempty"`
	PublicIpAddresses   *[]string `json:"public-ip-addresses,omitempty"`
}

// ComponentHealth defines model for ComponentHealth.
type ComponentHealth struct {
	ComponentDetails ComponentDetails       `json:"component-details"`
	ComponentId      openapi_types.UUID     `json:"component-id"`
	ComponentType    externalRef0.Connector `json:"component-type"`
	Status           ComponentStatus        `json:"status"`
}

// ComponentMetrics defines model for ComponentMetrics.
type ComponentMetrics struct {
	ComponentId      openapi_types.UUID `json:"component-id"`
	ComponentMetrics []Metric           `json:"component-metrics"`
	ComponentType    string             `json:"component-type"`
}

// ComponentStatus defines model for ComponentStatus.
type ComponentStatus struct {
	ApplicationsStatus *map[string]ApplicationStatus `json:"applications-status,omitempty"`
	OverallStatus      Status                        `json:"overall-status"`
}

// Metric defines model for Metric.
type Metric struct {
	// Name the metric that will be added to telegraf
	Name     string  `json:"name"`
	TagName  *string `json:"tag_name,omitempty"`
	TagValue *string `json:"tag_value,omitempty"`
	Value    float32 `json:"value"`
}

// ServiceStatus defines model for ServiceStatus.
type ServiceStatus struct {
	Msg           *string `json:"msg,omitempty"`
	ServiceStatus Status  `json:"service-status"`
}

// Status defines model for Status.
type Status string
