// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package models

import (
	"encoding/json"

	externalRef0 "agentless/infra/model/common"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Defines values for AzureSubscriptionsFilterScope.
const (
	ManagementGroup AzureSubscriptionsFilterScope = "management_group"
	Subscription    AzureSubscriptionsFilterScope = "subscription"
)

// Defines values for OnboardConfigOperationMode.
const (
	Enforcement OnboardConfigOperationMode = "enforcement"
	Visibility  OnboardConfigOperationMode = "visibility"
)

// Defines values for Status.
const (
	Completed    Status = "completed"
	Failed       Status = "failed"
	InProgress   Status = "in_progress"
	NotOnboarded Status = "not_onboarded"
	Partial      Status = "partial"
)

// Defines values for StepName.
const (
	StepNameAzureOnboarding StepName = "AzureOnboarding"
	StepNameMiscellaneous   StepName = "miscellaneous"
)

// AzureOnboarding defines model for AzureOnboarding.
type AzureOnboarding struct {
	Subscriptions AzureSubscriptionsFilter `json:"subscriptions"`
}

// AzureSubscriptionsFilter defines model for AzureSubscriptionsFilter.
type AzureSubscriptionsFilter struct {
	Scope  AzureSubscriptionsFilterScope `json:"scope"`
	Values []string                      `json:"values"`
}

// AzureSubscriptionsFilterScope defines model for AzureSubscriptionsFilter.Scope.
type AzureSubscriptionsFilterScope string

// AzureSubscriptionsStatus defines model for AzureSubscriptionsStatus.
type AzureSubscriptionsStatus = []map[string]Status

// Onboard defines model for Onboard.
type Onboard struct {
	ComponentId   openapi_types.UUID     `json:"component-id"`
	ConnectorType externalRef0.Connector `json:"connector-type"`
	General       OnboardConfig          `json:"general"`
	Steps         []Step                 `json:"steps"`
}

// OnboardConfig defines model for OnboardConfig.
type OnboardConfig struct {
	AutoDiscovery *bool                      `json:"auto-discovery,omitempty"`
	OperationMode OnboardConfigOperationMode `json:"operation-mode"`
	Revision      *int64                     `json:"revision,omitempty"`
}

// OnboardConfigOperationMode defines model for OnboardConfig.OperationMode.
type OnboardConfigOperationMode string

// Status defines model for Status.
type Status string

// StatusRequest defines model for StatusRequest.
type StatusRequest struct {
	DetailedStatus *StatusRequest_DetailedStatus `json:"detailed-status,omitempty"`
}

// StatusRequest_DetailedStatus defines model for StatusRequest.DetailedStatus.
type StatusRequest_DetailedStatus struct {
	union json.RawMessage
}

// StatusResponse defines model for StatusResponse.
type StatusResponse struct {
	DetailedStatus *StatusResponse_DetailedStatus `json:"detailed-status,omitempty"`
	OverallStatus  Status                         `json:"overall-status"`
	Steps          *[]struct {
		Status *Status `json:"status,omitempty"`
		Step   *string `json:"step,omitempty"`
	} `json:"steps,omitempty"`
}

// StatusResponse_DetailedStatus defines model for StatusResponse.DetailedStatus.
type StatusResponse_DetailedStatus struct {
	union json.RawMessage
}

// Step defines model for Step.
type Step struct {
	Action Step_Action `json:"action"`
	Name   StepName    `json:"name"`
}

// StepAction1 defines model for .
type StepAction1 = map[string]interface{}

// Step_Action defines model for Step.Action.
type Step_Action struct {
	union json.RawMessage
}

// StepName defines model for Step.Name.
type StepName string

// PostV1OperationsOnboardingJSONRequestBody defines body for PostV1OperationsOnboarding for application/json ContentType.
type PostV1OperationsOnboardingJSONRequestBody = Onboard

// PostV1OperationsStatusJSONRequestBody defines body for PostV1OperationsStatus for application/json ContentType.
type PostV1OperationsStatusJSONRequestBody = StatusRequest

// AsAzureSubscriptionsFilter returns the union data inside the StatusRequest_DetailedStatus as a AzureSubscriptionsFilter
func (t StatusRequest_DetailedStatus) AsAzureSubscriptionsFilter() (AzureSubscriptionsFilter, error) {
	var body AzureSubscriptionsFilter
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAzureSubscriptionsFilter overwrites any union data inside the StatusRequest_DetailedStatus as the provided AzureSubscriptionsFilter
func (t *StatusRequest_DetailedStatus) FromAzureSubscriptionsFilter(v AzureSubscriptionsFilter) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAzureSubscriptionsFilter performs a merge with any union data inside the StatusRequest_DetailedStatus, using the provided AzureSubscriptionsFilter
func (t *StatusRequest_DetailedStatus) MergeAzureSubscriptionsFilter(v AzureSubscriptionsFilter) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t StatusRequest_DetailedStatus) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *StatusRequest_DetailedStatus) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsAzureSubscriptionsStatus returns the union data inside the StatusResponse_DetailedStatus as a AzureSubscriptionsStatus
func (t StatusResponse_DetailedStatus) AsAzureSubscriptionsStatus() (AzureSubscriptionsStatus, error) {
	var body AzureSubscriptionsStatus
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAzureSubscriptionsStatus overwrites any union data inside the StatusResponse_DetailedStatus as the provided AzureSubscriptionsStatus
func (t *StatusResponse_DetailedStatus) FromAzureSubscriptionsStatus(v AzureSubscriptionsStatus) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAzureSubscriptionsStatus performs a merge with any union data inside the StatusResponse_DetailedStatus, using the provided AzureSubscriptionsStatus
func (t *StatusResponse_DetailedStatus) MergeAzureSubscriptionsStatus(v AzureSubscriptionsStatus) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t StatusResponse_DetailedStatus) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *StatusResponse_DetailedStatus) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsAzureOnboarding returns the union data inside the Step_Action as a AzureOnboarding
func (t Step_Action) AsAzureOnboarding() (AzureOnboarding, error) {
	var body AzureOnboarding
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAzureOnboarding overwrites any union data inside the Step_Action as the provided AzureOnboarding
func (t *Step_Action) FromAzureOnboarding(v AzureOnboarding) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAzureOnboarding performs a merge with any union data inside the Step_Action, using the provided AzureOnboarding
func (t *Step_Action) MergeAzureOnboarding(v AzureOnboarding) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsStepAction1 returns the union data inside the Step_Action as a StepAction1
func (t Step_Action) AsStepAction1() (StepAction1, error) {
	var body StepAction1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStepAction1 overwrites any union data inside the Step_Action as the provided StepAction1
func (t *Step_Action) FromStepAction1(v StepAction1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStepAction1 performs a merge with any union data inside the Step_Action, using the provided StepAction1
func (t *Step_Action) MergeStepAction1(v StepAction1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t Step_Action) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Step_Action) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
