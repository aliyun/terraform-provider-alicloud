package helper

import (
	"errors"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	sdkErrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
)

func TestNeedRetry(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"transport failure", errors.New(`Post "https://cs.aliyuncs.com": dial tcp: connection reset`), true},
		{"client timeout", sdkErrors.NewClientError(sdkErrors.TimeoutErrorCode, "timeout", nil), true},
		{"tea throttling code", &tea.SDKError{Code: tea.String("Throttling"), Message: tea.String("Requests throttled")}, true},
		{"tea service unavailable", &tea.SDKError{Code: tea.String("ServiceUnavailable"), Message: tea.String("service unavailable")}, true},
		{"tea rejected throttling", &tea.SDKError{Code: tea.String("Rejected.Throttling"), Message: tea.String("request rejected")}, true},
		{"tea server 5xx message", &tea.SDKError{Code: tea.String("InternalError"), Message: tea.String("code: 500, internal error")}, true},
		{"tea client timeout message", &tea.SDKError{Code: tea.String("Timeout"), Message: tea.String("Client.Timeout exceeded while awaiting headers")}, true},
		{"tea excluded already-enabled message", &tea.SDKError{Code: tea.String("OrderFailed"), Message: tea.String("code: 500, 您已开通过")}, false},
		{"tea fatal error", &tea.SDKError{Code: tea.String("ErrorClusterNotFound"), Message: tea.String("cluster not found")}, false},
		{"plain error", errors.New("ErrorClusterNotFound"), false},
	}

	for _, testCase := range testCases {
		if got := NeedRetry(testCase.err); got != testCase.want {
			t.Errorf("case %q: NeedRetry() = %v, want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestIncrementalWait(t *testing.T) {
	firstDuration := 100 * time.Millisecond
	increaseDuration := 50 * time.Millisecond

	originalJitter := incrementalWaitJitter
	incrementalWaitJitter = func() float64 { return 1.0 }
	t.Cleanup(func() { incrementalWaitJitter = originalJitter })

	waitFunc := IncrementalWait(firstDuration, increaseDuration)

	expected := []time.Duration{
		firstDuration,
		firstDuration + increaseDuration,
		firstDuration + 3*increaseDuration,
		firstDuration + 7*increaseDuration,
		firstDuration + 15*increaseDuration,
	}
	const tolerance = 30 * time.Millisecond

	elapsed := make([]time.Duration, len(expected))
	for i := range expected {
		start := time.Now()
		waitFunc()
		elapsed[i] = time.Since(start)
		t.Logf("wait call %d: measured=%v (expected>=%v)", i+1, elapsed[i], expected[i]-tolerance)
	}

	for i := 1; i < len(expected); i++ {
		if elapsed[i] < expected[i]-tolerance {
			t.Errorf("wait call %d: backoff did not accumulate; expected >= %v, measured %v (observed sequence %v)",
				i+1, expected[i]-tolerance, elapsed[i], elapsed)
		}
	}
}
