package v1alpha1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LabelSelectorOperator is the set of operators that can be used in a selector requirement.
type LabelSelectorOperator string

const (
	// StoragePolicyResourceName is the name of the singular StoragePolicy objects
	StoragePolicyResourceName = "storagepolicy"

	// StoragePolicyResourcePlural is the name of the plural StoragePolicy objects
	StoragePolicyResourcePlural = "storagepolicies"

  // LabelSelectorOpIn is operator where the key must have one of the values
	LabelSelectorOpIn LabelSelectorOperator = "In"
	// LabelSelectorOpNotIn is operator where the key must not have any of the values
	LabelSelectorOpNotIn LabelSelectorOperator = "NotIn"
	// LabelSelectorOpExists is operator where the key must exist
	LabelSelectorOpExists LabelSelectorOperator = "Exists"
	// LabelSelectorOpDoesNotExist is operator where the key must not exist
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "DoesNotExist"
	// LabelSelectorOpGt is operator where the key must be greater than the values
	LabelSelectorOpGt LabelSelectorOperator = "Gt"
	// LabelSelectorOpLt is operator where the key must be less than the values
	LabelSelectorOpLt LabelSelectorOperator = "Lt"
)

// LabelSelectorRequirement is a selector that contains values, a key, and an operator that
// relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	// +patchMergeKey=key
	// +patchStrategy=merge
	Key string `json:"key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist, Lt and Gt.
	Operator LabelSelectorOperator `json:"operator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	Values []string `json:"values"`
}

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
	Conditions []*LabelSelectorRequirement `json:"conditions"`
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
