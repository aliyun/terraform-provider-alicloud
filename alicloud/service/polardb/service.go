package polardb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type PayType string

const (
	PrePaid  = PayType("PrePaid")
	PostPaid = PayType("PostPaid")
	Prepaid  = PayType("Prepaid")
	Postpaid = PayType("Postpaid")
)

type NetworkType string

const (
	Vpc = NetworkType("Vpc")
)

type TimeType string

const (
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type RenewalStatus string

const (
	RenewAutoRenewal = RenewalStatus("AutoRenewal")
	RenewNormal      = RenewalStatus("Normal")
	RenewNotRenewal  = RenewalStatus("NotRenewal")
)

type Status string

const (
	Running = Status("Running")
	Deleted = Status("Deleted")
)

const (
	DefaultTimeout       = 120
	DefaultTimeoutMedium = 500
	DefaultLongTimeout   = 1000
	LOCAL_HOST_IP        = "127.0.0.1"
)

// Service wraps the PolarDB API calls used by the resources and data
// sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the acceptance
// tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}
