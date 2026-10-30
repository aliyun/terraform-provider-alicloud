package expressconnect

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
)

// Local copies of root-package typed constants used by the expressconnect
// product resources (router_interface, physical_connection). The leaf-package
// rule forbids importing the root alicloud package, so the minimal set used
// here is re-declared with identical values.

type Status string

const (
	Active   Status = "Active"
	Inactive Status = "Inactive"
	Idle     Status = "Idle"
	Deleted  Status = "Deleted"
)

type PayType string

const (
	PrePaid  PayType = "PrePaid"
	PostPaid PayType = "PostPaid"
)

type TimeType string

const (
	Month TimeType = "Month"
	Year  TimeType = "Year"
)

const (
	DefaultTimeout       = 120
	DefaultTimeoutMedium = 500
	DefaultLongTimeout   = 1000
)

// Root-package alias replicas (data_source_alicloud_common.go / common.go /
// tags.go). The leaf-package rule forbids importing root alicloud; these are
// re-bound to the same leaf symbols so bare call sites resolve without package
// qualification and survive local-variable shadowing of the "tags" name.
var (
	dataResourceIdHash      = helper.DataResourceIdHash
	tagsToMap               = tags.TagsToMap
	expandTagsToMapWithTags = tags.ExpandTagsToMapWithTags
	tagsSchema              = tags.TagsSchema
	ignoredTags             = tags.IgnoredTags
	parsingTags             = tags.ParsingTags
)

// validators aliases (validators.go): bare call sites in resource schema definitions.
var (
	IntBetween    = validators.IntBetween
	IntInSlice    = validators.IntInSlice
	StringInSlice = validators.StringInSlice
)

// errs const aliases (errors.go): bare call sites in service.go / service_v2.go /
// service_vpc.go (BuildStateConf target lists, error message constants). The
// leaf-package rule forbids importing root alicloud, so the minimal set used
// bare here is re-bound to errs leaf consts.
const (
	FailedToReachTargetStatus = errs.FailedToReachTargetStatus
	IdMsg                     = errs.IdMsg
	NotFoundMsg               = errs.NotFoundMsg
	WaitTimeoutMsg            = errs.WaitTimeoutMsg
)

// RouterType identifiers (extension_vpc.go): used for router_interface
// router_type validation.
type RouterType string

const (
	VRouter = RouterType("VRouter")
	VBR     = RouterType("VBR")
)
