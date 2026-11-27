package oss

import (
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

type LifecycleRuleStatus string

const (
	ExpirationStatusEnabled  = LifecycleRuleStatus("Enabled")
	ExpirationStatusDisabled = LifecycleRuleStatus("Disabled")
)

func ossNotFoundError(err error) bool {
	if e, ok := err.(oss.ServiceError); ok &&
		(e.StatusCode == 404 || strings.HasPrefix(e.Code, "NoSuch") || strings.HasPrefix(e.Message, "No Row found")) {
		return true
	}
	return false
}

// concurrentUpdateBucketFailed is the transient conflict code OSS returns
// when several configuration updates hit the same bucket at once; the bucket
// configuration writes are idempotent, so retrying resolves it.
const concurrentUpdateBucketFailed = "ConcurrentUpdateBucketFailed"

// needOssBucketConfigRetry classifies request errors for OSS bucket
// configuration sub-resource writes (PutBucketVersioning and siblings):
// the shared transient-error predicate plus the OSS-scoped concurrent-update
// conflict code, which only these writes can observe. Keeping the code here
// leaves the common helper.NeedRetry free of OSS-specific retry policy.
func needOssBucketConfigRetry(err error) bool {
	return errs.IsExpectedErrors(err, []string{concurrentUpdateBucketFailed}) || helper.NeedRetry(err)
}
