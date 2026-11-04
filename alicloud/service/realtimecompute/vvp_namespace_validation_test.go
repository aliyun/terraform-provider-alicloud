package realtimecompute

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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

// TestUnitRealtimeComputeVvpNamespaceCreateResultUncertain pins the boundary of the post-create
// identity recovery: only failures that do not prove the request was rejected (timeouts, transport
// errors, 5xx) may be turned into a lookup by <instance_id>:<namespace>; deterministic product
// answers must be reported unchanged.
func TestUnitRealtimeComputeVvpNamespaceCreateResultUncertain(t *testing.T) {
	status := func(code int) *int { return &code }
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "cancelled", err: context.Canceled, want: false},
		{name: "deadline exceeded", err: context.DeadlineExceeded, want: true},
		{name: "net timeout", err: &net.OpError{Op: "dial", Err: &timeoutError{}}, want: true},
		{name: "server 503", err: &tea.SDKError{StatusCode: status(503), Message: tea.String("service unavailable")}, want: true},
		{name: "server 500 wrapped", err: errors.Join(errors.New("calling CreateNamespace failed"), &tea.SDKError{StatusCode: status(500)}), want: true},
		{name: "server 503 through the shared error wrapper", err: errs.WrapErrorf(&tea.SDKError{StatusCode: status(503), Message: tea.String("service unavailable")}, "calling CreateNamespace failed"), want: true},
		{name: "validation through the shared error wrapper", err: errs.WrapErrorf(&tea.SDKError{StatusCode: status(400), Code: tea.String("903038")}, "calling CreateNamespace failed"), want: false},
		{name: "validation 400", err: &tea.SDKError{StatusCode: status(400), Message: tea.String("InvalidParameter")}, want: false},
		{name: "product 903038", err: &tea.SDKError{StatusCode: status(400), Code: tea.String("903038")}, want: false},
		{name: "timeout text", err: errors.New("Post https://foasconsole.aliyuncs.com: net/http: request canceled (Client.Timeout exceeded while awaiting headers)"), want: true},
		{name: "connection reset text", err: errors.New("read tcp 10.0.0.1:1234->10.0.0.2:443: connection reset by peer"), want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := vvpNamespaceCreateResultUncertain(tc.err); got != tc.want {
				t.Fatalf("vvpNamespaceCreateResultUncertain(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "i/o timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

// TestUnitRealtimeComputeVvpNamespaceOperationBudget pins the single-budget contract: every stage
// of one operation shares one deadline, the remaining budget shrinks instead of restarting, and a
// cancelled parent context still aborts the operation.
func TestUnitRealtimeComputeVvpNamespaceOperationBudget(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	opCtx, remaining, cancel := vvpNamespaceOperationBudget(parent, 2*time.Second)
	if _, ok := opCtx.Deadline(); !ok {
		t.Fatalf("operation context has no deadline")
	}
	first := remaining()
	if first <= 0 || first > 2*time.Second {
		t.Fatalf("initial remaining budget = %v, want (0, 2s]", first)
	}
	time.Sleep(20 * time.Millisecond)
	second := remaining()
	if second >= first {
		t.Fatalf("remaining budget did not shrink: first=%v second=%v", first, second)
	}
	cancelParent()
	if opCtx.Err() != context.Canceled {
		t.Fatalf("cancelling the parent did not cancel the operation context: %v", opCtx.Err())
	}
	cancel()

	short, shortRemaining, cancelShort := vvpNamespaceOperationBudget(context.Background(), 10*time.Millisecond)
	defer cancelShort()
	time.Sleep(30 * time.Millisecond)
	if shortRemaining() > 0 {
		t.Fatalf("expired budget still reports time left: %v", shortRemaining())
	}
	if short.Err() != context.DeadlineExceeded {
		t.Fatalf("expired operation context error = %v, want deadline exceeded", short.Err())
	}
}

// vvpNamespaceUpdateTestServer serves the foasconsole actions the update path can call and records
// every Action it saw, so a case can assert that a stage did or did not reach the service.
func vvpNamespaceUpdateTestServer(t *testing.T) (*httptest.Server, func() []string) {
	t.Helper()
	var mu sync.Mutex
	var actions []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		action := r.FormValue("Action")
		mu.Lock()
		actions = append(actions, action)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if action == "ModifyNamespaceSpec" {
			fmt.Fprint(w, `{"RequestId":"mock-modify"}`)
			return
		}
		fmt.Fprint(w, `{"RequestId":"mock-describe","Namespaces":[{"InstanceId":"inst-unit","Namespace":"demo","Status":"SUCCESS","Ha":false,"GmtCreate":1700000000,"GmtModified":1700000000,"ResourceSpec":{"Cpu":1,"MemoryGB":4},"ResourceUsed":{"Cpu":1,"Cu":1,"MemoryGB":4},"Tags":[]}]}`)
	}))
	t.Cleanup(server.Close)
	return server, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), actions...)
	}
}

// vvpNamespaceUpdateTestResource configures a fresh resource against a local foasconsole endpoint.
func vvpNamespaceUpdateTestResource(t *testing.T, endpoint string) (*vvpNamespaceResource, fwschema.Schema) {
	t.Helper()
	t.Setenv("NO_PROXY", "127.0.0.1")
	t.Setenv("no_proxy", "127.0.0.1")
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").
		SetAccessKeyId("test-access-key-id").
		SetAccessKeySecret("test-access-key-secret"))
	if err != nil {
		t.Fatalf("build credential failed: %v", err)
	}
	endpoints := new(sync.Map)
	config := connectivity.Config{
		AccessKey:            "test-access-key-id",
		SecretKey:            "test-access-key-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build client failed: %v", err)
	}
	// Both the tea RPC layer and the legacy SDK prepend the scheme themselves.
	endpoints.Store("foasconsole", strings.TrimPrefix(endpoint, "http://"))

	res := &vvpNamespaceResource{}
	var configureResp fwresource.ConfigureResponse
	res.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}
	var schemaResp fwresource.SchemaResponse
	res.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	return res, schemaResp.Schema
}

// vvpNamespaceUpdateTestValue builds one state/plan value with an unchanged 1/4 specification and
// the given timeouts.update, so the update skips ModifyNamespaceSpec and goes straight to the
// confirmation read.
func vvpNamespaceUpdateTestValue(t *testing.T, sch fwschema.Schema, updateTimeout string) tftypes.Value {
	t.Helper()
	objectType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("resource schema does not describe an object")
	}
	specListType, ok := objectType.AttributeTypes["resource_spec"].(tftypes.List)
	if !ok {
		t.Fatalf("resource_spec is not a list: %T", objectType.AttributeTypes["resource_spec"])
	}
	spec := tftypes.NewValue(specListType, []tftypes.Value{
		tftypes.NewValue(specListType.ElementType, map[string]tftypes.Value{
			"cpu":       tftypes.NewValue(tftypes.Number, 1),
			"memory_gb": tftypes.NewValue(tftypes.Number, 4),
		}),
	})
	timeouts := tftypes.NewValue(objectType.AttributeTypes["timeouts"], map[string]tftypes.Value{
		"create": tftypes.NewValue(tftypes.String, nil),
		"update": tftypes.NewValue(tftypes.String, updateTimeout),
		"delete": tftypes.NewValue(tftypes.String, nil),
	})
	return tftypes.NewValue(objectType, map[string]tftypes.Value{
		"id":            tftypes.NewValue(tftypes.String, "inst-unit:demo"),
		"gmt_create":    tftypes.NewValue(tftypes.Number, 1700000000),
		"gmt_modified":  tftypes.NewValue(tftypes.Number, 1700000000),
		"ha":            tftypes.NewValue(tftypes.Bool, false),
		"instance_id":   tftypes.NewValue(tftypes.String, "inst-unit"),
		"namespace":     tftypes.NewValue(tftypes.String, "demo"),
		"resource_spec": spec,
		"resource_used": tftypes.NewValue(objectType.AttributeTypes["resource_used"], nil),
		"status":        tftypes.NewValue(tftypes.String, "SUCCESS"),
		"tags":          tftypes.NewValue(objectType.AttributeTypes["tags"], nil),
		"timeouts":      timeouts,
	})
}

// TestUnitRealtimeComputeVvpNamespaceUpdateConfirmationReadBudget drives the public Update request
// path and pins the confirmation read to the operation budget: an update whose timeouts.update is
// already spent reports the exhausted budget and keeps the namespace in state without issuing a
// read, a cancelled lifecycle behaves the same way, and a budget that still has time refreshes the
// state with one read.
func TestUnitRealtimeComputeVvpNamespaceUpdateConfirmationReadBudget(t *testing.T) {
	cases := []struct {
		name             string
		updateTimeout    string
		cancelLifecycle  bool
		wantErrSubstring string
		wantReads        int
	}{
		{
			name:             "exhausted budget keeps state without a read",
			updateTimeout:    "1ns",
			wantErrSubstring: "timeouts.update",
			wantReads:        0,
		},
		{
			name:             "exhausted budget on a cancelled lifecycle keeps state without a read",
			updateTimeout:    "1ns",
			cancelLifecycle:  true,
			wantErrSubstring: "timeouts.update",
			wantReads:        0,
		},
		{
			name:          "remaining budget refreshes the state",
			updateTimeout: "1m",
			wantReads:     1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server, observed := vvpNamespaceUpdateTestServer(t)
			res, sch := vvpNamespaceUpdateTestResource(t, server.URL)
			value := vvpNamespaceUpdateTestValue(t, sch, tc.updateTimeout)

			ctx := context.Background()
			if tc.cancelLifecycle {
				cancelled, cancel := context.WithCancel(ctx)
				cancel()
				ctx = cancelled
			}
			req := fwresource.UpdateRequest{
				Plan:  tfsdk.Plan{Schema: sch, Raw: value},
				State: tfsdk.State{Schema: sch, Raw: value},
			}
			resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: value}}
			res.Update(ctx, req, resp)

			reads := 0
			for _, action := range observed() {
				if action == "DescribeNamespaces" {
					reads++
				}
			}
			if reads != tc.wantReads {
				t.Fatalf("DescribeNamespaces calls = %d, want %d (actions: %v)", reads, tc.wantReads, observed())
			}

			if tc.wantErrSubstring == "" {
				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
				}
			} else {
				if !resp.Diagnostics.HasError() {
					t.Fatalf("expected an error containing %q, got none", tc.wantErrSubstring)
				}
				detail := resp.Diagnostics.Errors()[0].Summary() + " " + resp.Diagnostics.Errors()[0].Detail()
				if !strings.Contains(detail, tc.wantErrSubstring) {
					t.Fatalf("expected an error containing %q, got %q", tc.wantErrSubstring, detail)
				}
			}

			if resp.State.Raw.IsNull() {
				t.Fatalf("the update dropped the namespace from state")
			}
			var id types.String
			resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
			if id.ValueString() != "inst-unit:demo" {
				t.Fatalf("state id = %q, want inst-unit:demo", id.ValueString())
			}
		})
	}
}
