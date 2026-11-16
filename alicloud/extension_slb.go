package alicloud

type SchedulerType string

const (
	WRRScheduler = SchedulerType("wrr")
	WLCScheduler = SchedulerType("wlc")
	RRScheduler  = SchedulerType("rr")
)

type FlagType string

const (
	OnFlag  = FlagType("on")
	OffFlag = FlagType("off")
)

type RsType string

const (
	ENI = FlagType("eni")
	ECS = FlagType("ecs")
)

type AclType string

const (
	AclTypeBlack = AclType("black")
	AclTypeWhite = AclType("white")
)

type StickySessionType string

const (
	InsertStickySessionType = StickySessionType("insert")
	ServerStickySessionType = StickySessionType("server")
)

const BackendServerPort = -520

type HealthCheckHttpCodeType string

const (
	HTTP_2XX = HealthCheckHttpCodeType("http_2xx")
	HTTP_3XX = HealthCheckHttpCodeType("http_3xx")
	HTTP_4XX = HealthCheckHttpCodeType("http_4xx")
	HTTP_5XX = HealthCheckHttpCodeType("http_5xx")
)

type HealthCheckType string

const (
	TCPHealthCheckType  = HealthCheckType("tcp")
	HTTPHealthCheckType = HealthCheckType("http")
)

type LoadBalancerSpecType string

const (
	S1Small  = "slb.s1.small"
	S2Small  = "slb.s2.small"
	S2Medium = "slb.s2.medium"
	S3Small  = "slb.s3.small"
	S3Medium = "slb.s3.medium"
	S3Large  = "slb.s3.large"
	S4Large  = "slb.s4.large"
)

type ListenerErr struct {
	ErrType string
	Err     error
}

func (e *ListenerErr) Error() string {
	return e.ErrType + " " + e.Err.Error()

}
