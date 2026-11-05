package helper

import (
	"context"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// ServiceUnavailable is a retryable error code returned by several
// Alibaba Cloud APIs.
const ServiceUnavailable = "ServiceUnavailable"

// NeedRetry reports whether err is transient and the caller should retry.
// It is reused verbatim from the root package (alicloud/errors.go).
// NeedRetry and IncrementalWait lived in the deleted package
// alicloud/retry, whose name collided with the SDK helper/retry and forced
// a sdkretry alias at every call site; they live in this package so
// service packages can import the SDK helper/retry under its natural name.
func NeedRetry(err error) bool {
	if err == nil {
		return false
	}

	postRegex := regexp.MustCompile("^Post [\"]*https://.*")
	if postRegex.MatchString(err.Error()) {
		return true
	}

	if e, ok := err.(*errors.ClientError); ok && e.ErrorCode() == errors.TimeoutErrorCode {
		return true
	}

	throttlingRegex := regexp.MustCompile("Throttling")
	codeRegex := regexp.MustCompile("^code: 5[\\d]{2}")

	if e, ok := err.(*tea.SDKError); ok {
		if strings.Contains(*e.Message, "code: 500, 您已开通过") {
			return false
		}
		if strings.Contains(*e.Message, "Client.Timeout") {
			return true
		}
		if *e.Code == ServiceUnavailable || *e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(*e.Code) || codeRegex.MatchString(*e.Message) {
			return true
		}
	}

	if e, ok := err.(*errors.ServerError); ok {
		return e.ErrorCode() == ServiceUnavailable || e.ErrorCode() == "Rejected.Throttling" || throttlingRegex.MatchString(e.ErrorCode()) || codeRegex.MatchString(e.Message())
	}

	if e, ok := err.(*common.Error); ok {
		return e.Code == ServiceUnavailable || e.Code == "Rejected.Throttling" || throttlingRegex.MatchString(e.Code) || codeRegex.MatchString(e.Message)
	}

	return false
}

var incrementalWaitJitter = rand.Float64

// incrementalBackoff returns the jittered wait before the (retryCount+1)-th retry.
func incrementalBackoff(retryCount int, firstDuration time.Duration, increaseDuration time.Duration) time.Duration {
	upper := float64(firstDuration) + (math.Pow(2, float64(retryCount))-1)*float64(increaseDuration)
	if upper > float64(20*time.Second) {
		upper = float64(20 * time.Second)
	}
	return time.Duration(upper * incrementalWaitJitter())
}

// IncrementalWait returns a wait function with incremental backoff; it is
// reused verbatim from the root package (alicloud/common.go).
func IncrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	retryCount := 0
	return func() {
		time.Sleep(incrementalBackoff(retryCount, firstDuration, increaseDuration))
		retryCount++
	}
}

// incrementalWaitContext is IncrementalWait for a lifecycle operation: the backoff sleep is
// abandoned as soon as ctx is done and the function reports that it did not wait, so a cancelled
// Terraform operation stops retrying instead of blocking on the next backoff.
func incrementalWaitContext(ctx context.Context, firstDuration time.Duration, increaseDuration time.Duration) func() bool {
	retryCount := 0
	return func() bool {
		timer := time.NewTimer(incrementalBackoff(retryCount, firstDuration, increaseDuration))
		defer timer.Stop()
		retryCount++
		select {
		case <-ctx.Done():
			return false
		case <-timer.C:
			return true
		}
	}
}

// retryContextError is the error reported when the caller's context ended before a retry could run.
func retryContextError(ctx context.Context, fallback error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fallback
}

// RetryRpcPost runs call under retry.RetryContext with the shared incremental backoff, then wraps
// the final error with the action name. retryable decides which errors are transient. The call closure
// keeps the client, endpoint and request in the calling package, so this shared wrapper does not
// depend on connectivity and product packages keep one retrying write path instead of one each.
// ctx is the lifecycle context: cancelling it aborts the backoff and prevents further calls.
func RetryRpcPost(ctx context.Context, action string, call func() (map[string]interface{}, error), timeout time.Duration, retryable func(error) bool) (map[string]interface{}, error) {
	wait := incrementalWaitContext(ctx, 3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		resp, e := call()
		response = resp
		if e != nil {
			if retryable(e) {
				if !wait() {
					return retry.NonRetryableError(retryContextError(ctx, e))
				}
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		return response, errs.WrapErrorf(err, "calling %s failed", action)
	}
	return response, nil
}
