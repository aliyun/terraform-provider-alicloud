package oss

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// Service wraps the OSS API calls used by the SDKv2 resources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the acceptance tests) build values
// with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}
