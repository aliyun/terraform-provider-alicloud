package connectivity

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
)

// AliyunClient 保留旧包的类型及方法集兼容，实现在 client 包中。
type AliyunClient = client.AliyunClient

// NewClient 保留旧构造入口；新调用方应直接使用 client.NewClient。
func NewClient(c *config.Config) (*AliyunClient, error) {
	return client.NewClientWithValidRegions(c, ValidRegions)
}

// UnreadableResponseError 保留旧包的 RPC 错误类型兼容。
type UnreadableResponseError = client.UnreadableResponseError

var (
	ErrUnreadableResponse  = client.ErrUnreadableResponse
	ErrRpcRejectedResponse = client.ErrRpcRejectedResponse
)

// RedactLogValue 保留旧包的日志脱敏入口。
func RedactLogValue(value interface{}) interface{} {
	return client.RedactLogValue(value)
}

// Temporarily maintain map for old ecs client methods and store special endpoint information
var EndpointMap = map[string]string{
	"cn-shenzhen-su18-b01":        "ecs.aliyuncs.com",
	"cn-beijing":                  "ecs.aliyuncs.com",
	"cn-shenzhen-st4-d01":         "ecs.aliyuncs.com",
	"cn-haidian-cm12-c01":         "ecs.aliyuncs.com",
	"cn-hangzhou-internal-prod-1": "ecs.aliyuncs.com",
	"cn-qingdao":                  "ecs.aliyuncs.com",
	"cn-shanghai":                 "ecs.aliyuncs.com",
	"cn-shanghai-finance-1":       "ecs.aliyuncs.com",
	"cn-hongkong":                 "ecs.aliyuncs.com",
	"us-west-1":                   "ecs.aliyuncs.com",
	"cn-shenzhen":                 "ecs.aliyuncs.com",
	"cn-shanghai-et15-b01":        "ecs.aliyuncs.com",
	"cn-hangzhou-bj-b01":          "ecs.aliyuncs.com",
	"cn-zhangbei-na61-b01":        "ecs.aliyuncs.com",
	"cn-shenzhen-finance-1":       "ecs.aliyuncs.com",
	"cn-shanghai-et2-b01":         "ecs.aliyuncs.com",
	"ap-southeast-1":              "ecs.aliyuncs.com",
	"cn-beijing-nu16-b01":         "ecs.aliyuncs.com",
	"us-east-1":                   "ecs.aliyuncs.com",
	"cn-fujian":                   "ecs.aliyuncs.com",
	"cn-hangzhou":                 "ecs.aliyuncs.com",
}
