// Package polardb holds the SDKv2 implementation, service and tests for the
// PolarDB resources whose flat-layout code has been migrated into a product
// package. Each resource keeps its own Describe/wait/tag service methods in a
// service_<type>.go file; this file only holds the Service definition shared by
// them.
package polardb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// Service wraps the PolarDB API calls used by the resources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported, so
// callers outside the package (including the acceptance tests) build values
// with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}
