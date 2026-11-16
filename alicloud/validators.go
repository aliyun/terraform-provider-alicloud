package alicloud

import "github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

var (
	IntBetween       = validators.IntBetween
	IntAtLeast       = validators.IntAtLeast
	StringInSlice    = validators.StringInSlice
	StringLenBetween = validators.StringLenBetween
)
