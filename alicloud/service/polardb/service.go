package polardb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// NewService returns a PolarDB service wrapper around the given client so that
// callers outside this package (for example the package's acceptance tests) can
// use the generated Describe helpers.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}
