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
// Status mirrors the root package status type, mirrored locally so the OSS
// service package is self-contained without importing the root package.
type Status string

const Deleted = Status("Deleted")

// DefaultTimeoutMedium mirrors the root package constant (alicloud/common.go)
// so the OSS service package stays self-contained without importing the root
// package. The value and semantics are identical to the root definition.
const DefaultTimeoutMedium = 500

// COLON_SEPARATED mirrors the root package constant (alicloud/common.go) so
// the OSS service package stays self-contained without importing the root
// package. The value and semantics are identical to the root definition.
const COLON_SEPARATED = ":"

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}
