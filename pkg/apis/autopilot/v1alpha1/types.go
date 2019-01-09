package v1alpha1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StoragePolicyResourceName   = "storagepolicy"
	StoragePolicyResourcePlural = "storagepolicies"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StoragePolicy represents pairing with other clusters
type StoragePolicy struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            StoragePolicySpec `json:"spec"`
	// TODO: add status
	//Status          StoragePolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StoragePolicyList is a list of StoragePolicy objects in Kubernetes
type StoragePolicyList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`

	Items []StoragePolicy `json:"items"`
}

// StoragePolicySpec is the spec to create the cluster pair
type StoragePolicySpec struct {
	// Weight defines the weight of the policy which allows to break the tie with other conflicting policies. A policy with
	// higher weight wins over one with lower weight.
	// (optional)
	Weight int64 `json:"weight,omitempty"`
	// Enforcement specifies the enforcement type for policy. Can take values: required or preferred.
	// (optional)
	Enforcement EnforcementType `json:"enforcement,omitempty"`
	// Object is the entity on which to check the conditions
	Object PolicyObject `json:"object"`
	// Conditions are the conditions to check on the policy objects
	Conditions []*meta.LabelSelectorRequirement `json:"conditions"`
	// Action is the action to run for the policy when the conditions are met
	Action PolicyAction `json:"action"`
}

// PolicyObject defines an object for the policy
type PolicyObject struct {
	// Type is the type of the policy object
	Type PolicyObjectType `json:"type"`
	// LabelSelector selects the policy objects
	meta.LabelSelector
}

// PolicyAction defines an action for the policy
type PolicyAction struct {
	// Name is the name of the policy
	Name PolicyActionName `json:"name"`
	// ActionObject is the target object for the policy (optional)
	ActionObject PolicyObject `json:"actionObject,omitempty"`
}
