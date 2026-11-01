package realtimecompute

import (
	"strings"
	"testing"

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

// TestUnitRealtimeComputeVvpNamespaceRetryable pins the classification the resource passes to the
// shared helper.RetryRpcPost path: 903047 and the settling conflict/status codes stay retryable,
// 903099 stays fatal, and the general transient branches come from helper.NeedRetry.
func TestUnitRealtimeComputeVvpNamespaceRetryable(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "903047 is transient", err: &tea.SDKError{Code: tea.String("903047"), Message: tea.String("code: 400, a resource spec change is in progress")}, want: true},
		{name: "903099 is fatal", err: &tea.SDKError{Code: tea.String("903099"), Message: tea.String("code: 400, target is below the used capacity")}, want: false},
		{name: "903043 is fatal", err: &tea.SDKError{Code: tea.String("903043"), Message: tea.String("code: 400, cannot operate")}, want: false},
		{name: "InvalidParameter is fatal", err: &tea.SDKError{Code: tea.String("InvalidParameter"), Message: tea.String("code: 400, bad parameter")}, want: false},
		{name: "throttling is transient through helper.NeedRetry", err: &tea.SDKError{Code: tea.String("Throttling.User"), Message: tea.String("code: 400, throttled")}, want: true},
		{name: "dependency violation is transient", err: &tea.SDKError{Code: tea.String("DependencyViolation.HaVip"), Message: tea.String("code: 400, dependent resource")}, want: true},
		{name: "nil stays fatal", err: nil, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := vvpNamespaceRetryable(tc.err); got != tc.want {
				t.Fatalf("vvpNamespaceRetryable(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
