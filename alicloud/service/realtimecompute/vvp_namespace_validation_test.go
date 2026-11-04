package realtimecompute

import (
	"context"
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// TestUnitRealtimeComputeVvpNamespaceResponseRows pins the read decoding contract: only a decoded
// array counts as evidence of what DescribeNamespaces returned, so an omitted, null or mistyped
// field is an unreadable response rather than an empty namespace set.
func TestUnitRealtimeComputeVvpNamespaceResponseRows(t *testing.T) {
	cases := []struct {
		name     string
		response map[string]interface{}
		wantLen  int
		wantErr  string
	}{
		{name: "decoded empty collection", response: map[string]interface{}{"Namespaces": []interface{}{}}, wantLen: 0},
		{
			name:     "other namespace present",
			response: map[string]interface{}{"Namespaces": []interface{}{map[string]interface{}{"Namespace": "other"}}},
			wantLen:  1,
		},
		{name: "nil response", response: nil, wantErr: "empty response"},
		{name: "missing field", response: map[string]interface{}{}, wantErr: "does not contain Namespaces"},
		{name: "null field", response: map[string]interface{}{"Namespaces": nil}, wantErr: "is null"},
		{name: "mistyped field", response: map[string]interface{}{"Namespaces": "unexpected"}, wantErr: "unexpected type"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rows, err := vvpNamespaceResponseRows(tc.response)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got rows=%v err=%v", tc.wantErr, rows, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if len(rows) != tc.wantLen {
				t.Fatalf("expected %d rows, got %d", tc.wantLen, len(rows))
			}
		})
	}
}

// TestUnitRealtimeComputeVvpNamespaceSelectReadObject pins the two positive not-found shapes and
// the malformed-row error the read path relies on: an empty decoded collection and a decoded
// collection without the requested identity both return (nil, nil), while a row that is not an
// object is an error instead of a silent absence.
func TestUnitRealtimeComputeVvpNamespaceSelectReadObject(t *testing.T) {
	identities := []fwadapt.ReadIdentity{{Path: "$.Namespace", Value: "target", AllowMissing: false}}

	cases := []struct {
		name    string
		rows    []interface{}
		wantNil bool
		wantErr string
	}{
		{name: "empty collection", rows: []interface{}{}, wantNil: true},
		{
			name:    "target absent but other objects remain",
			rows:    []interface{}{map[string]interface{}{"Namespace": "other"}},
			wantNil: true,
		},
		{
			name: "target present",
			rows: []interface{}{map[string]interface{}{"Namespace": "other"}, map[string]interface{}{"Namespace": "target"}},
		},
		{name: "row is not an object", rows: []interface{}{"unexpected"}, wantErr: "not an object"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object, err := fwadapt.SelectReadObject(tc.rows, identities)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got object=%v err=%v", tc.wantErr, object, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if tc.wantNil && object != nil {
				t.Fatalf("expected no match, got %v", object)
			}
			if !tc.wantNil && object == nil {
				t.Fatalf("expected a match, got none")
			}
		})
	}
}

// TestUnitRealtimeComputeVvpNamespaceResourceSpecValidator pins the plan-time element rules: a
// known null element is rejected even though it passes SizeAtLeast(1), while unknown elements and
// unknown attribute values are deferred to apply time.
func TestUnitRealtimeComputeVvpNamespaceResourceSpecValidator(t *testing.T) {
	objType := vvpNamespaceResourceSpecObjType()
	attrTypes := vvpNamespaceResourceSpecAttrTypes()
	object := func(cpu, memory attr.Value) attr.Value {
		return types.ObjectValueMust(attrTypes, map[string]attr.Value{"cpu": cpu, "memory_gb": memory})
	}

	cases := []struct {
		name    string
		value   types.List
		wantErr string
	}{
		{
			name:  "valid element",
			value: types.ListValueMust(objType, []attr.Value{object(types.Int64Value(2), types.Int64Value(8))}),
		},
		{
			name:    "null element rejected",
			value:   types.ListValueMust(objType, []attr.Value{types.ObjectNull(attrTypes)}),
			wantErr: "is null",
		},
		{
			name:  "unknown element deferred",
			value: types.ListValueMust(objType, []attr.Value{types.ObjectUnknown(attrTypes)}),
		},
		{
			name:  "unknown attribute values deferred",
			value: types.ListValueMust(objType, []attr.Value{object(types.Int64Unknown(), types.Int64Unknown())}),
		},
		{
			name:  "unknown list deferred",
			value: types.ListUnknown(objType),
		},
		{
			name:  "null list skipped",
			value: types.ListNull(objType),
		},
		{
			name:    "invalid ratio rejected",
			value:   types.ListValueMust(objType, []attr.Value{object(types.Int64Value(2), types.Int64Value(9))}),
			wantErr: "must be 4 times",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := &validator.ListResponse{}
			vvpNamespaceResourceSpecValidator{}.ValidateList(context.Background(), validator.ListRequest{
				Path:        path.Root("resource_spec"),
				ConfigValue: tc.value,
			}, response)
			hasError := response.Diagnostics.HasError()
			if tc.wantErr == "" {
				if hasError {
					t.Fatalf("unexpected diagnostics: %v", response.Diagnostics)
				}
				return
			}
			if !hasError {
				t.Fatalf("expected an error containing %q, got none", tc.wantErr)
			}
			if !strings.Contains(response.Diagnostics.Errors()[0].Detail(), tc.wantErr) &&
				!strings.Contains(response.Diagnostics.Errors()[0].Summary(), tc.wantErr) {
				t.Fatalf("expected an error containing %q, got %v", tc.wantErr, response.Diagnostics)
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
