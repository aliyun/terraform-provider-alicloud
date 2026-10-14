package helper

import (
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/denverdino/aliyungo/common"
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

// IncrementalWait returns a wait function with incremental backoff; it is
// reused verbatim from the root package (alicloud/common.go).
func IncrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	retryCount := 0
	return func() {
		upper := float64(firstDuration) + (math.Pow(2, float64(retryCount))-1)*float64(increaseDuration)
		if upper > float64(20*time.Second) {
			upper = float64(20 * time.Second)
		}
		waitTime := time.Duration(upper * incrementalWaitJitter())
		time.Sleep(waitTime)
		retryCount++
	}
}
