// Package cdn holds the SDKv2 resources, data sources and service helpers for the
// CDN namespace. The declarations below reproduce a small set of legacy root-package
// symbols that the migrated code still depends on and that have no leaf-package
// counterpart; their values and behaviour match the root package exactly so the
// migrated resources keep their original semantics.
package cdn

import "github.com/aliyun/terraform-provider-alicloud/alicloud/tags"

// Status mirrors the root-package status type used by the legacy CDN wait helpers.
type Status string

const (
	// Online is the CDN domain online status, matching the root-package value.
	Online = Status("online")
	// Configuring is the CDN domain configuring status, matching the root-package value.
	Configuring = Status("configuring")
	// Deleted is the CDN domain deleted status, matching the root-package value.
	Deleted = Status("Deleted")
)

// DefaultTimeout and DefaultTimeoutMedium are the legacy default timeout constants
// (in seconds) used by the CDN resources; they mirror the root-package values.
const (
	DefaultTimeout       = 120
	DefaultTimeoutMedium = 500
)

// CharityPageUrl is the legacy CDN error-page URL constant used by the deprecated
// alicloud_cdn_domain resource; it mirrors the root-package value.
const CharityPageUrl = "http://promotion.alicdn.com/help/oss/error.html"

// TagResourceType aliases the shared tag resource type so the migrated CDN helpers
// keep the same signature as the root-package originals.
type TagResourceType = tags.TagResourceType
