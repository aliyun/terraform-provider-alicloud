package alicloud

import "github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

const UserId = "userId"
const ScalingGroup = "scaling_group"

type ComparisonOperator string

const (
	Gt  = ComparisonOperator(">")
	Gte = ComparisonOperator(">=")
	Lt  = ComparisonOperator("<")
	Lte = ComparisonOperator("<=")
)

type Statistics string

const (
	Avg = Statistics("Average")
	Min = Statistics("Minimum")
	Max = Statistics("Maximum")
)

type Period int

const (
	OneMinite     = Period(60)
	TwoMinite     = Period(120)
	FiveMinite    = Period(300)
	FifteenMinite = Period(900)
)

type MetricType string

const (
	System = MetricType("system")
	Custom = MetricType("custom")
)

type ActionResult string

const (
	Continue = ActionResult("CONTINUE")
	Abandon  = ActionResult("ABANDON")
)

type LifecycleTransition string

const (
	ScaleOut = LifecycleTransition("SCALE_OUT")
	ScaleIn  = LifecycleTransition("SCALE_IN")
)

type AdjustmentType string

const (
	QuantityChangeInCapacity = AdjustmentType("QuantityChangeInCapacity")
	PercentChangeInCapacity  = AdjustmentType("PercentChangeInCapacity")
	TotalCapacity            = AdjustmentType("TotalCapacity")
)

type RecurrenceType string

const (
	Daily   = RecurrenceType("Daily")
	Weekly  = RecurrenceType("Weekly")
	Monthly = RecurrenceType("Monthly")
)

type InstanceCreationType string

const (
	AutoCreated = InstanceCreationType("AutoCreated")
	Attached    = InstanceCreationType("Attached")
)

func EssCommonRequestInit(region string, code ServiceCode, domain CommonRequestDomain) *requests.CommonRequest {
	request := CommonRequestInit(region, code, domain)
	request.Version = ApiVersion20140828
	return request
}
