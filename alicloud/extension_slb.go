package alicloud

import (
	"fmt"
	"strings"
)

type SchedulerType string

const (
	WRRScheduler = SchedulerType("wrr")
	WLCScheduler = SchedulerType("wlc")
)

type FlagType string

const (
	OnFlag  = FlagType("on")
	OffFlag = FlagType("off")
)

type TlsCipherPolicy string

const (
	TlsCipherPolicy_1_0        = TlsCipherPolicy("tls_cipher_policy_1_0")
	TlsCipherPolicy_1_1        = TlsCipherPolicy("tls_cipher_policy_1_1")
	TlsCipherPolicy_1_2        = TlsCipherPolicy("tls_cipher_policy_1_2")
	TlsCipherPolicy_1_2_STRICT = TlsCipherPolicy("tls_cipher_policy_1_2_strict")
)

type AclType string

const (
	AclTypeBlack = AclType("black")
	AclTypeWhite = AclType("white")
)

type IPVersion string

const (
	IPVersion4 = IPVersion("ipv4")
	IPVersion6 = IPVersion("ipv6")
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
)

type ListenerErr struct {
	ErrType string
	Err     error
}

func (e *ListenerErr) Error() string {
	return e.ErrType + " " + e.Err.Error()

}

func expandBackendServersToString(list []interface{}, weight int) string {
	if len(list) < 1 {
		return ""
	}
	var items []string
	for _, id := range list {
		items = append(items, fmt.Sprintf("{'ServerId':'%s','Weight':'%d'}", id, weight))
	}
	return fmt.Sprintf("[%s]", strings.Join(items, COMMA_SEPARATED))
}

func expandBackendServersWithPortToString(items []interface{}) string {

	if len(items) < 1 {
		return ""
	}
	var servers []string
	for _, server := range items {
		s := server.(map[string]interface{})

		var server_ids []interface{}
		var port, weight int
		if v, ok := s["server_ids"]; ok {
			server_ids = v.([]interface{})
		}
		if v, ok := s["port"]; ok {
			port = v.(int)
		}
		if v, ok := s["weight"]; ok {
			weight = v.(int)
		}

		for _, id := range server_ids {
			str := fmt.Sprintf("{'ServerId':'%s','Port':'%d','Weight':'%d'}", strings.Trim(id.(string), " "), port, weight)

			servers = append(servers, str)
		}

	}
	return fmt.Sprintf("[%s]", strings.Join(servers, COMMA_SEPARATED))
}

func getLoadBalancerSpecOrder(spec string) int {
	order := 0
	switch spec {
	case S1Small:
		order = 0
	case S2Small:
		order = 1
	case S2Medium:
		order = 2
	case S3Small:
		order = 3
	case S3Medium:
		order = 4
	case S3Large:
		order = 5
	}

	return order
}
