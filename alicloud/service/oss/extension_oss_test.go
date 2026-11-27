package oss

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

// OSS bucket configuration sub-resources (PutBucketVersioning and siblings)
// fail with ConcurrentUpdateBucketFailed when independent resources update
// the same bucket at once. Their write retries route through
// needOssBucketConfigRetry, which classifies the conflict as transient while
// the shared helper.NeedRetry predicate stays free of OSS retry policy.
func TestUnitNeedOssBucketConfigRetryRetriesConcurrentUpdateConflict(t *testing.T) {
	conflict := &tea.SDKError{
		Code:    tea.String("ConcurrentUpdateBucketFailed"),
		Message: tea.String("The bucket is being updated concurrently. Please try again later."),
		Data:    tea.String(""),
	}
	if !needOssBucketConfigRetry(conflict) {
		t.Fatal("the OSS bucket configuration retry predicate did not classify ConcurrentUpdateBucketFailed as retryable")
	}
	if helper.NeedRetry(conflict) {
		t.Fatal("the shared NeedRetry predicate must not classify the OSS conflict code as retryable")
	}

	// A fatal OSS code keeps failing fast, so the whitelist did not loosen
	// non-conflict error handling.
	fatal := &tea.SDKError{Code: tea.String("NoSuchBucket"), Message: tea.String("The specified bucket does not exist."), Data: tea.String("")}
	if needOssBucketConfigRetry(fatal) {
		t.Fatal("the OSS bucket configuration retry predicate classified NoSuchBucket as retryable")
	}
	if needOssBucketConfigRetry(nil) {
		t.Fatal("the OSS bucket configuration retry predicate classified a nil error as retryable")
	}
}
