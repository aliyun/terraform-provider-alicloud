// Package retry is reused verbatim from alicloud/errors.go (NeedRetry) and
// alicloud/common.go (IncrementalWait): service packages cannot import the
// root monolith, so the retry classification and wait timing live here.
// TODO refactor to the alicloud root package, currently is tmp workaround
package retry

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

const ServiceUnavailable = "ServiceUnavailable"

// NeedRetry is reused verbatim from the root package (alicloud/errors.go).
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

// IncrementalWait is reused verbatim from the root package
// (alicloud/common.go).
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
