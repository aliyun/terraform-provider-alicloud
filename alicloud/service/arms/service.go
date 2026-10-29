package arms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// Service wraps the ARMS API calls used by the resources of this package.
// The root alicloud package cannot be imported from a service package, so the
// encapsulated calls live here with the resources.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}
