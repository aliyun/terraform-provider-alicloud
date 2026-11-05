package helper

import (
	"context"
	"errors"
	"strings"
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

// RetryRpcPost is the shared write path, so a cancelled lifecycle must stop it instead of letting
// the next multi-second backoff run to completion. The first backoff is 3s, which makes the old
// sleep-based behaviour unmistakably slower than the cancellation below.
func TestUnitRetryRpcPostStopsOnCancelledContext(t *testing.T) {
	originalJitter := incrementalWaitJitter
	incrementalWaitJitter = func() float64 { return 1.0 }
	t.Cleanup(func() { incrementalWaitJitter = originalJitter })

	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	done := make(chan error, 1)
	go func() {
		_, err := RetryRpcPost(ctx, "Test", func() (map[string]interface{}, error) {
			calls++
			return nil, errors.New("transient")
		}, time.Minute, func(error) bool { return true })
		done <- err
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("RetryRpcPost returned no error after cancellation")
		}
		if !strings.Contains(err.Error(), context.Canceled.Error()) {
			t.Fatalf("the cancellation was not reported: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RetryRpcPost kept waiting out its backoff after the context was cancelled")
	}
	if calls > 1 {
		t.Fatalf("expected no retry after cancellation, got %d calls", calls)
	}
}

func TestUnitIncrementalWaitContextStopsWhenCancelled(t *testing.T) {
	originalJitter := incrementalWaitJitter
	incrementalWaitJitter = func() float64 { return 1.0 }
	t.Cleanup(func() { incrementalWaitJitter = originalJitter })

	ctx, cancel := context.WithCancel(context.Background())
	wait := incrementalWaitContext(ctx, 3*time.Second, 5*time.Second)
	cancel()

	start := time.Now()
	if wait() {
		t.Fatal("incrementalWaitContext reported a completed wait after cancellation")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("incrementalWaitContext slept %s after cancellation", elapsed)
	}
}
