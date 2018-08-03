package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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

func buildSlbCommonRequest() *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.Domain = "slb.aliyuncs.com"
	request.Version = ApiVersion20140515
	return request
}
