package connectivity

import "github.com/aliyun/terraform-provider-alicloud/alicloud/client"

// ConvertKebabToSnake 委托客户端层实现，保留旧调用方兼容。
func ConvertKebabToSnake(s string) string {
	return client.ConvertKebabToSnake(s)
}
