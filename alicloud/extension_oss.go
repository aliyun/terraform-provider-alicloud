package alicloud

import (
	"errors"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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

// ossNoSuchBucketError reports whether err is precisely the OSS NoSuchBucket
// service error. Unlike ossNotFoundError it never matches a bare 404.
func ossNoSuchBucketError(err error) bool {
	var serviceErr oss.ServiceError
	if errors.As(err, &serviceErr) {
		return serviceErr.Code == "NoSuchBucket"
	}
	return false
}

// ossRetryableServiceError reports whether err is a transient OSS failure that
// is worth retrying within a bounded wait: transport errors, 5xx responses and
// throttling. Deterministic 4xx failures are not retryable.
func ossRetryableServiceError(err error) bool {
	var serviceErr oss.ServiceError
	if !errors.As(err, &serviceErr) {
		return true
	}
	if serviceErr.StatusCode >= 500 {
		return true
	}
	return serviceErr.Code == "Throttling" || serviceErr.Code == "RequestTimeout"
}
