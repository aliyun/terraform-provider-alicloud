package ess

const UserId = "userId"
const ScalingGroup = "scaling_group"
const GroupId = "groupId"

// Status mirrors the root alicloud.Status string type used by WaitFor signatures.
// No leaf package exports it, so the ess package keeps a local compat copy.
type Status string

type ScalingRuleType string

const (
	SimpleScalingRule         = ScalingRuleType("SimpleScalingRule")
	TargetTrackingScalingRule = ScalingRuleType("TargetTrackingScalingRule")
	StepScalingRule           = ScalingRuleType("StepScalingRule")
	PredictiveScalingRule     = ScalingRuleType("PredictiveScalingRule")
)

type BatchSize int

const (
	AttachDetachLoadbalancersBatchsize  = BatchSize(5)
	AttachDetachDbinstancesBatchsize    = BatchSize(5)
	AttachDetachAlbServerGroupBatchsize = BatchSize(5)
)

type MaxItems int
type MinItems int

const (
	MaxScalingConfigurationInstanceTypes = MaxItems(10)
	MaxResourcePoolOptionsPrivatePoolIds = MaxItems(20)
	MinResourcePoolOptionsPrivatePoolIds = MinItems(1)
)

// Root-only constants copied from alicloud/common.go and alicloud/extension_*.go.
// These are defined directly in the flat root package rather than re-exported
// from a leaf package, so the ess subpackage keeps local copies instead of
// importing the root alicloud package (which would form an import cycle).

const (
	Active    = Status("Active")
	Inactive  = Status("Inactive")
	Deleted   = Status("Deleted")
	InService = Status("InService")
	Null      = Status("Null")
)

const DefaultTimeout = 120
const DefaultLongTimeout = 1000
const COLON_SEPARATED = ":"

const (
	Average = "Average"
	Minimum = "Minimum"
	Maximum = "Maximum"
)

type CreditSpecification string

const CreditSpecificationStandard = CreditSpecification("Standard")
const CreditSpecificationUnlimited = CreditSpecification("Unlimited")

type DiskCategory string

const DiskCloudEfficiency = DiskCategory("cloud_efficiency")
