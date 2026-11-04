package alikafka

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
)

// Bare-symbol aliases mirroring the root alicloud package's unqualified
// re-exports so migrated files keep their original call sites without
// reverse-importing the root package:
//   - TagResourceType    (root common.go: type TagResourceType = tags.TagResourceType)
//   - dataResourceIdHash (root data_source_alicloud_common.go: = helper.DataResourceIdHash)
//   - ignoredTags        (root tags.go: = tags.IgnoredTags)
type TagResourceType = tags.TagResourceType

var (
	dataResourceIdHash = helper.DataResourceIdHash
	ignoredTags        = tags.IgnoredTags
)

// Local typed constants mirroring alicloud.PayType / alicloud.Status so the
// migrated resource and data-source files (schema defaults, request values,
// wait/retry helpers) are self-contained within the product package and do
// not reverse-import the root alicloud package.

// PayType mirrors alicloud.PayType (common.go).
type PayType string

const (
	PrePaid  = PayType("PrePaid")
	PostPaid = PayType("PostPaid")
)

// Status mirrors alicloud.Status (common.go) for the local service wait/retry
// helpers (WaitForAlikafkaInstance / ConsumerGroup / Topic / SaslUser / SaslAcl).
type Status string

const (
	Creating = Status("Creating")
	Running  = Status("Running")
	Deleted  = Status("Deleted")
)
