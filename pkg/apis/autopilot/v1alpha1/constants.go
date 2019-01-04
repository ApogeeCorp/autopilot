package v1alpha1

type (
	// PolicyObjectType is the type for policy objects
	PolicyObjectType string
	// PolicyConditionName is the type for policy condition names
	PolicyConditionName string
	// PolicyActionName is the type for policy actions
	PolicyActionName string
)

const (
	openstorageDomain = "openstorage.io"
	// PolicyObjectPrefix is the key for any openstorage object used for policies
	PolicyObjectPrefix = openstorageDomain + ".object"
	// PolicyConditionPrefix is the key for any openstorage condition used for policies
	PolicyConditionPrefix = openstorageDomain + ".condition"
	// PolicyActionPrefix is the key for any openstorage action used for policies
	PolicyActionPrefix = openstorageDomain + ".action"

	// PolicyObjectTypeVolume is the key for volume objects
	PolicyObjectTypeVolume PolicyObjectType = PolicyObjectPrefix + "/volume"
	// PolicyObjectTypeStoragePool is the key for storagepool objects
	PolicyObjectTypeStoragePool PolicyObjectType = PolicyObjectPrefix + "/storagepool"
	// PolicyObjectTypeNode  is the key for node objects
	PolicyObjectTypeNode PolicyObjectType = PolicyObjectPrefix + "/node"
	// PolicyObjectTypeDisk is the key for disk objects
	PolicyObjectTypeDisk PolicyObjectType = PolicyObjectPrefix + "/disk"

	// PolicyConditionVolume is the key for volume conditions for policies
	PolicyConditionVolume = PolicyConditionPrefix + ".volume"
	// PolicyConditionStoragePool is the key for storagepool conditions for policies
	PolicyConditionStoragePool = PolicyConditionPrefix + ".storagepool"
	// PolicyConditionNode is the key for node conditions for policies
	PolicyConditionNode = PolicyConditionPrefix + ".node"
	// PolicyConditionDisk is the key for disk conditions for policies
	PolicyConditionDisk = PolicyConditionPrefix + ".disk"

	// PolicyActionVolume is the key for volume actions for policies
	PolicyActionVolume = PolicyActionPrefix + ".volume"
	// PolicyActionStoragePool is the key for storagepool actions for policies
	PolicyActionStoragePool = PolicyActionPrefix + ".storagepool"
	// PolicyActionNode is the key for node actions for policies
	PolicyActionNode = PolicyActionPrefix + ".node"
	// PolicyActionDisk is the key for disk actions for policies
	PolicyActionDisk = PolicyActionPrefix + ".disk"
)

const (
	latencyMS       = "/latency_ms"
	actionMove      = "/move"
	actionRebalance = "/rebalance"
)

// EnforcementType Defines the types of enforcement on the given policy
type EnforcementType string

const (
	// EnforcementRequired specifies that the policy is required and must be strictly enforced
	EnforcementRequired EnforcementType = "required"
	// EnforcementPreferred specifies that the policy is preferred and can be best effort
	EnforcementPreferred EnforcementType = "preferred"
)

const (
	/***** Volume conditions *****/

	// PolicyConditionVolumeLatencyMS is the latency (reads + writes) for a volume in milliseconds
	PolicyConditionVolumeLatencyMS PolicyConditionName = PolicyConditionVolume + latencyMS

	/***** Storage pool conditions *****/

	// PolicyConditionStoragePoolLatencyMS is the latency (reads + writes) for a storage pool in milliseconds
	PolicyConditionStoragePoolLatencyMS PolicyConditionName = PolicyConditionStoragePool + latencyMS

	/***** Disk conditions *****/

	// PolicyConditionDiskLatencyMS is the latency (reads + writes) for a disk in milliseconds
	PolicyConditionDiskLatencyMS PolicyConditionName = PolicyConditionDisk + latencyMS
)

const (
	/***** Volume actions *****/

	// PolicyActionVolumeMove is an action to move volumes
	PolicyActionVolumeMove PolicyActionName = PolicyActionVolume + actionMove
	/***** Node actions *****/

	// PolicyActionNodeRebalance is an action to rebalance a node
	PolicyActionNodeRebalance PolicyActionName = PolicyActionNode + actionRebalance
)
