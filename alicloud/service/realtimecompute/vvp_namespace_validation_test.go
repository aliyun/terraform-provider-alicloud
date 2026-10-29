package realtimecompute

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
)

// TestUnitRealtimeComputeVvpNamespaceResourceSpecRatio covers the plan-time rule that mirrors the
// product contract: a namespace resource spec must pair memory in GB as exactly four times CPU.
func TestUnitRealtimeComputeVvpNamespaceResourceSpecRatio(t *testing.T) {
	cases := []struct {
		name   string
		cpu    int64
		memory int64
		want   string
	}{
		{name: "valid 2/8", cpu: 2, memory: 8},
		{name: "valid 3/12", cpu: 3, memory: 12},
		{name: "valid 4/16", cpu: 4, memory: 16},
		{name: "memory too small", cpu: 2, memory: 4, want: "must be 4 times"},
		{name: "memory too large", cpu: 2, memory: 9, want: "must be 4 times"},
		{name: "zero cpu", cpu: 0, memory: 8, want: "positive"},
		{name: "negative memory", cpu: 2, memory: -8, want: "positive"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := vvpNamespaceResourceSpecRatioError(tc.cpu, tc.memory)
			switch {
			case tc.want == "" && got != "":
				t.Fatalf("cpu=%d memory_gb=%d: unexpected error %q", tc.cpu, tc.memory, got)
			case tc.want != "" && !strings.Contains(got, tc.want):
				t.Fatalf("cpu=%d memory_gb=%d: error %q does not mention %q", tc.cpu, tc.memory, got, tc.want)
			}
		})
	}
}

// TestUnitRealtimeComputeVvpNamespaceNeedRetry pins the transient classification: 903047 stays
// retryable and a plain parameter error stays fatal.
func TestUnitRealtimeComputeVvpNamespaceNeedRetry(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "903047 is transient", err: &tea.SDKError{Code: tea.String("903047"), Message: tea.String("code: 400, a resource spec change is in progress")}, want: true},
		{name: "903043 is fatal", err: &tea.SDKError{Code: tea.String("903043"), Message: tea.String("code: 400, cannot operate")}, want: false},
		{name: "InvalidParameter is fatal", err: &tea.SDKError{Code: tea.String("InvalidParameter"), Message: tea.String("code: 400, bad parameter")}, want: false},
		{name: "throttling is transient", err: &tea.SDKError{Code: tea.String("Throttling.User"), Message: tea.String("code: 400, throttled")}, want: true},
		{name: "nil stays fatal", err: nil, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := needRetry(tc.err); got != tc.want {
				t.Fatalf("needRetry(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// TestUnitRealtimeComputeVvpNamespaceWaitStopsOnCancel proves the retry backoff honours the
// lifecycle context instead of sleeping out the whole interval.
func TestUnitRealtimeComputeVvpNamespaceWaitStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	wait := incrementalWaitContext(ctx, time.Hour, time.Hour)
	start := time.Now()
	if wait() {
		t.Fatal("wait reported that it slept after the context was cancelled")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("wait blocked for %s after cancellation", elapsed)
	}
}
