package alicloud

import "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"

// WrapError/WrapErrorf must stay var aliases rather than forwarding
// functions: they read runtime.Caller(1) to attribute errors to the caller,
// and a wrapper would insert a frame and misattribute every error.
var (
	WrapError  = errs.WrapError
	WrapErrorf = errs.WrapErrorf
)
