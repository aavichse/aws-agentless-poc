// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package models

// Defines values for LookupRequestDirection.
const (
	Inbound  LookupRequestDirection = "inbound"
	Outbound LookupRequestDirection = "outbound"
)

// LookupRequest defines model for LookupRequest.
type LookupRequest struct {
	// AgentIp Source IP address.
	AgentIp string `json:"agent-ip"`

	// AgentItemId reporting agent azure resource id
	AgentItemId *string `json:"agent-item-id,omitempty"`

	// AgentMac reporting agent mac address
	AgentMac string `json:"agent-mac"`

	// ConnectionTime Start time of the network event.** for aggregated events (epoch in seconds)
	ConnectionTime *int64 `json:"connection-time,omitempty"`

	// Direction Indicates the direction of the network event
	Direction *LookupRequestDirection `json:"direction,omitempty"`

	// RemoteIp Destination IP address.
	RemoteIp string `json:"remote-ip"`

	// RemotePort Destination port.
	RemotePort *int64 `json:"remote-port,omitempty"`
}

// LookupRequestDirection Indicates the direction of the network event
type LookupRequestDirection string

// PostV1ProviderLookupJSONRequestBody defines body for PostV1ProviderLookup for application/json ContentType.
type PostV1ProviderLookupJSONRequestBody = LookupRequest
