// Package apigateway contains the apigateway service resources and data sources
// migrated from the root alicloud package. The root package cannot be imported
// from a service package, so the encapsulated apigateway calls live here.
package apigateway

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// Status mirrors the root package status type used by the WaitFor helpers.
type Status string

const (
	Deleted            = Status("Deleted")
	Normal             = Status("Normal")
	DefaultTimeout     = 120
	DefaultLongTimeout = 1000
	COLON_SEPARATED    = ":"
)

// NewCloudApiService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the acceptance
// tests) build values with it.
func NewCloudApiService(client *connectivity.AliyunClient) *CloudApiService {
	return &CloudApiService{client: client}
}

// NewApiGatewayServiceV2 is the exported constructor for ApiGatewayServiceV2.
// External test packages use it instead of the unexported client field.
func NewApiGatewayServiceV2(client *connectivity.AliyunClient) *ApiGatewayServiceV2 {
	return &ApiGatewayServiceV2{client: client}
}
