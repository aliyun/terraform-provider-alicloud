package expressconnectrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// newVpcAssociationUnitClient points the ExpressConnectRouter endpoint at an
// isolated HTTP fixture so the lifecycle code runs without cloud access.
func newVpcAssociationUnitClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build credential: %v", err)
	}
	endpoint := trimHTTPPrefix(server.URL)
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build client: %v", err)
	}
	endpoints.Store("expressconnectrouter", endpoint)
	return client
}

func trimHTTPPrefix(rawURL string) string {
	return strings.TrimPrefix(rawURL, "http://")
}

func writeVpcAssociationJSON(t *testing.T, w http.ResponseWriter, status int, body interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Errorf("encode response: %v", err)
	}
}

func writeVpcAssociationActive(t *testing.T, w http.ResponseWriter) {
	writeVpcAssociationJSON(t, w, http.StatusOK, map[string]interface{}{
		"RequestId": "unit-test",
		"AssociationList": []map[string]interface{}{{
			"AssociationId":       "assoc-1",
			"AssociationNodeType": "VPC",
			"EcrId":               "ecr-1",
			"VpcId":               "vpc-1",
			"RegionId":            "cn-hangzhou",
			"Status":              "ACTIVE",
			"Description":         "unit-description",
			"AllowedPrefixesMode": "MatchMode",
			"AllowedPrefixes":     []string{"172.16.1.0/24"},
			"VpcOwnerId":          1,
		}},
	})
}

// silenceVpcAssociationWaits removes the production backoff and the first state
// poll delay so retry regressions run instantly.
func silenceVpcAssociationWaits(t *testing.T) {
	t.Helper()
	previousWait, previousDelay := vpcAssociationRetryWait, vpcAssociationStateWaitDelay
	vpcAssociationRetryWait = func() func() { return func() {} }
	vpcAssociationStateWaitDelay = 0
	t.Cleanup(func() {
		vpcAssociationRetryWait = previousWait
		vpcAssociationStateWaitDelay = previousDelay
	})
}

func newVpcAssociationUpdateResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, ResourceVpcAssociation().Schema, raw)
	d.SetId("ecr-1:assoc-1:vpc-1")
	return d
}

func TestUnitVpcAssociationCreateRetryCodes(t *testing.T) {
	silenceVpcAssociationWaits(t)
	// The expected transient set is restated here so the test fails when a
	// documented retryable code silently disappears from the resource.
	for _, code := range []string{
		"Conflict.Lock",
		"IncorrectStatus.Ecr",
		"IncorrectStatus.Vpc",
		"IncorrectStatus.VpcOperating",
		"IncorrectStatus.Vswitch",
		"IncorrectStatus.RouteEntry",
		"OperationDenied.ForceDeleting",
		"OperationFailed.CallingVpcYaochi",
		"InternalErrorTemp",
	} {
		code := code
		t.Run(code, func(t *testing.T) {
			var createCalls int32
			client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
				}
				switch r.Form.Get("Action") {
				case "CreateExpressConnectRouterAssociation":
					if atomic.AddInt32(&createCalls, 1) == 1 {
						writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
							"Code": code, "Message": "transient " + code, "RequestId": "unit-test",
						})
						return
					}
					writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{
						"AssociationId": "assoc-1", "RequestId": "unit-test",
					})
				case "DescribeExpressConnectRouterAssociation":
					writeVpcAssociationActive(t, w)
				default:
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
					writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
						"Code": "InvalidParameter", "Message": "unexpected action",
					})
				}
			})

			d := schema.TestResourceDataRaw(t, ResourceVpcAssociation().Schema, map[string]interface{}{
				"ecr_id":                "ecr-1",
				"vpc_id":                "vpc-1",
				"association_region_id": "cn-hangzhou",
			})
			d.MarkNewResource()
			if err := ResourceVpcAssociation().Create(d, client); err != nil {
				t.Fatalf("create with transient %s should succeed after a retry: %v", code, err)
			}
			if got := atomic.LoadInt32(&createCalls); got != 2 {
				t.Fatalf("CreateExpressConnectRouterAssociation calls = %d, want 2", got)
			}
			if got, want := d.Id(), "ecr-1:assoc-1:vpc-1"; got != want {
				t.Fatalf("resource id = %q, want %q", got, want)
			}
		})
	}
}

func TestUnitVpcAssociationCreateTerminalCodesAreNotRetried(t *testing.T) {
	silenceVpcAssociationWaits(t)
	for _, code := range []string{
		"InvalidParameter",
		"QuotaExceeded.VpcCountPerEcr",
		"IllegalParamFormat.AllowedPrefixes",
		"ResourceAlreadyExist.VpcId",
		"UnsupportedRegion",
	} {
		code := code
		t.Run(code, func(t *testing.T) {
			var createCalls int32
			client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
				}
				if r.Form.Get("Action") != "CreateExpressConnectRouterAssociation" {
					t.Errorf("terminal error must not trigger %q", r.Form.Get("Action"))
					writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
						"Code": "InvalidParameter", "Message": "unexpected action",
					})
					return
				}
				atomic.AddInt32(&createCalls, 1)
				writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
					"Code": code, "Message": "terminal " + code, "RequestId": "unit-test",
				})
			})

			d := schema.TestResourceDataRaw(t, ResourceVpcAssociation().Schema, map[string]interface{}{
				"ecr_id":                "ecr-1",
				"vpc_id":                "vpc-1",
				"association_region_id": "cn-hangzhou",
			})
			d.MarkNewResource()
			err := ResourceVpcAssociation().Create(d, client)
			if err == nil {
				t.Fatalf("terminal error %s must be returned", code)
			}
			if got := atomic.LoadInt32(&createCalls); got != 1 {
				t.Fatalf("CreateExpressConnectRouterAssociation calls = %d, want 1", got)
			}
			if d.Id() != "" {
				t.Fatalf("failed create assigned id %q", d.Id())
			}
		})
	}
}

func TestUnitVpcAssociationCreateSendsModeOnlyWithPrefixes(t *testing.T) {
	silenceVpcAssociationWaits(t)

	t.Run("mode_with_prefixes_is_sent", func(t *testing.T) {
		var mu sync.Mutex
		var lastForm url.Values
		client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Errorf("parse request: %v", err)
			}
			switch r.Form.Get("Action") {
			case "CreateExpressConnectRouterAssociation":
				mu.Lock()
				lastForm = r.Form
				mu.Unlock()
				writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{
					"AssociationId": "assoc-1", "RequestId": "unit-test",
				})
			case "DescribeExpressConnectRouterAssociation":
				writeVpcAssociationActive(t, w)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})

		d := schema.TestResourceDataRaw(t, ResourceVpcAssociation().Schema, map[string]interface{}{
			"ecr_id":                "ecr-1",
			"vpc_id":                "vpc-1",
			"association_region_id": "cn-hangzhou",
			"allowed_prefixes":      []interface{}{"172.16.1.0/24"},
			"allowed_prefixes_mode": "MatchMode",
			"description":           "unit-description",
		})
		d.MarkNewResource()
		if err := ResourceVpcAssociation().Create(d, client); err != nil {
			t.Fatalf("create with mode and prefixes: %v", err)
		}
		mu.Lock()
		form := lastForm
		mu.Unlock()
		if got := form.Get("AllowedPrefixesMode"); got != "MatchMode" {
			t.Errorf("AllowedPrefixesMode = %q, want MatchMode", got)
		}
		if got := form.Get("Description"); got != "unit-description" {
			t.Errorf("Description = %q, want unit-description", got)
		}
		if !formHasArrayValue(form, "AllowedPrefixes", "172.16.1.0/24") {
			t.Errorf("AllowedPrefixes did not carry the configured prefix: %v", form)
		}
	})

	t.Run("mode_without_prefixes_is_rejected_before_request", func(t *testing.T) {
		var requests int32
		client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&requests, 1)
			t.Errorf("unexpected request %q for an invalid mode-only configuration", r.Form.Get("Action"))
			writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
				"Code": "InvalidParameter", "Message": "unexpected request",
			})
		})

		d := schema.TestResourceDataRaw(t, ResourceVpcAssociation().Schema, map[string]interface{}{
			"ecr_id":                "ecr-1",
			"vpc_id":                "vpc-1",
			"association_region_id": "cn-hangzhou",
			"allowed_prefixes_mode": "MatchMode",
		})
		d.MarkNewResource()
		err := ResourceVpcAssociation().Create(d, client)
		if err == nil || !containsAll(err.Error(), "allowed_prefixes_mode", "non-empty", "allowed_prefixes") {
			t.Fatalf("mode-only create error = %v, want the local combination error", err)
		}
		if got := atomic.LoadInt32(&requests); got != 0 {
			t.Fatalf("invalid mode-only create issued %d requests, want 0", got)
		}
	})
}

func TestUnitVpcAssociationUpdateDescriptionRetry(t *testing.T) {
	silenceVpcAssociationWaits(t)
	var modifyCalls int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociation":
			if atomic.AddInt32(&modifyCalls, 1) == 1 {
				writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
					"Code": "IncorrectStatus.Ecr", "Message": "another operation is running", "RequestId": "unit-test",
				})
				return
			}
			if got := r.Form.Get("Description"); got != "unit-description" {
				t.Errorf("Description = %q, want unit-description", got)
			}
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := newVpcAssociationUpdateResourceData(t, map[string]interface{}{
		"ecr_id":                "ecr-1",
		"vpc_id":                "vpc-1",
		"association_region_id": "cn-hangzhou",
		"description":           "unit-description",
	})
	if err := ResourceVpcAssociation().Update(d, client); err != nil {
		t.Fatalf("description update should retry a recoverable status: %v", err)
	}
	if got := atomic.LoadInt32(&modifyCalls); got != 2 {
		t.Fatalf("ModifyExpressConnectRouterAssociation calls = %d, want 2", got)
	}
}

// newVpcAssociationUpdateDataFromState builds a ResourceData whose previous
// state is explicit, so a change to an empty string is a real diff instead of
// an unset value.
func newVpcAssociationUpdateDataFromState(t *testing.T, state map[string]string, config map[string]interface{}) *schema.ResourceData {
	t.Helper()
	m := schema.InternalMap(ResourceVpcAssociation().Schema)
	instanceState := &terraform.InstanceState{ID: "ecr-1:assoc-1:vpc-1", Attributes: state}
	instanceState.Attributes["id"] = instanceState.ID
	diff, err := m.Diff(context.Background(), instanceState, terraform.NewResourceConfigRaw(config), nil, nil, false)
	if err != nil {
		t.Fatalf("build diff: %v", err)
	}
	d, err := m.Data(instanceState, diff)
	if err != nil {
		t.Fatalf("build resource data: %v", err)
	}
	return d
}

func TestUnitVpcAssociationUpdateParameterErrorFailsFast(t *testing.T) {
	silenceVpcAssociationWaits(t)
	var calls int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociation":
			atomic.AddInt32(&calls, 1)
			writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
				"Code": "IllegalParamFormat.Description", "Message": "description format is illegal", "RequestId": "unit-test",
			})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := newVpcAssociationUpdateResourceData(t, map[string]interface{}{
		"ecr_id":                "ecr-1",
		"vpc_id":                "vpc-1",
		"association_region_id": "cn-hangzhou",
		"description":           "unit-description",
	})
	if err := ResourceVpcAssociation().Update(d, client); err == nil {
		t.Fatalf("parameter errors must be returned to the caller")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("ModifyExpressConnectRouterAssociation calls = %d, want 1 (parameter errors must not retry)", got)
	}
}

func TestUnitVpcAssociationUpdateSharesOneDeadline(t *testing.T) {
	silenceVpcAssociationWaits(t)
	start := time.Now()
	var offsetSeconds int64
	previousNow := vpcAssociationNow
	vpcAssociationNow = func() time.Time {
		return start.Add(time.Duration(atomic.LoadInt64(&offsetSeconds)) * time.Second)
	}
	t.Cleanup(func() { vpcAssociationNow = previousNow })

	var modifies, describes int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociation":
			atomic.AddInt32(&modifies, 1)
			// Consume the whole five minute update budget during the first write.
			atomic.StoreInt64(&offsetSeconds, 300)
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "ModifyExpressConnectRouterAssociationAllowedPrefix":
			atomic.AddInt32(&modifies, 1)
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			atomic.AddInt32(&describes, 1)
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := newVpcAssociationUpdateDataFromState(t,
		map[string]string{
			"ecr_id":                "ecr-1",
			"vpc_id":                "vpc-1",
			"association_region_id": "cn-hangzhou",
			"description":           "old-description",
			"allowed_prefixes.#":    "1",
			"allowed_prefixes.0":    "172.16.1.0/24",
		},
		map[string]interface{}{
			"ecr_id":                "ecr-1",
			"vpc_id":                "vpc-1",
			"association_region_id": "cn-hangzhou",
			"description":           "new-description",
			"allowed_prefixes":      []interface{}{"172.16.2.0/24"},
		})
	if !d.HasChange("description") || !d.HasChange("allowed_prefixes") {
		t.Fatalf("both sequential writes must be pending")
	}
	err := ResourceVpcAssociation().Update(d, client)
	if err == nil || !containsAll(err.Error(), "Timeout", "shared update deadline") {
		t.Fatalf("shared budget error = %v, want the shared deadline timeout", err)
	}
	if got := atomic.LoadInt32(&modifies); got != 1 {
		t.Fatalf("modify calls = %d, want 1; a per-step budget would have issued the second write", got)
	}
	if got := atomic.LoadInt32(&describes); got != 0 {
		t.Fatalf("describe calls = %d, want 0 after the budget is exhausted", got)
	}
}

func TestUnitVpcAssociationUpdateDescriptionClearSendsEmptyString(t *testing.T) {
	silenceVpcAssociationWaits(t)
	var forms []url.Values
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociation":
			forms = append(forms, r.Form)
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := newVpcAssociationUpdateDataFromState(t,
		map[string]string{
			"ecr_id":                "ecr-1",
			"vpc_id":                "vpc-1",
			"association_region_id": "cn-hangzhou",
			"description":           "old-description",
		},
		map[string]interface{}{
			"ecr_id":                "ecr-1",
			"vpc_id":                "vpc-1",
			"association_region_id": "cn-hangzhou",
			"description":           "",
		})
	if !d.HasChange("description") {
		t.Fatalf("clearing the description must be a change")
	}
	if err := ResourceVpcAssociation().Update(d, client); err != nil {
		t.Fatalf("clear description: %v", err)
	}
	if len(forms) != 1 {
		t.Fatalf("ModifyExpressConnectRouterAssociation calls = %d, want 1", len(forms))
	}
	values, present := forms[0]["Description"]
	if !present || len(values) != 1 || values[0] != "" {
		t.Fatalf("Description present=%v values=%v, want one explicit empty value", present, values)
	}
}

func TestUnitVpcAssociationUpdatePrefixSendsStoredMode(t *testing.T) {
	silenceVpcAssociationWaits(t)
	var mu sync.Mutex
	var forms []url.Values
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociationAllowedPrefix":
			mu.Lock()
			forms = append(forms, r.Form)
			count := len(forms)
			mu.Unlock()
			if count == 1 {
				writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
					"Code": "IncorrectStatus.EcrAssociation", "Message": "association is busy", "RequestId": "unit-test",
				})
				return
			}
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	d := newVpcAssociationUpdateResourceData(t, map[string]interface{}{
		"ecr_id":                "ecr-1",
		"vpc_id":                "vpc-1",
		"association_region_id": "cn-hangzhou",
		"allowed_prefixes":      []interface{}{"172.16.1.0/24"},
	})
	if err := d.Set("allowed_prefixes_mode", "IncrementalMode"); err != nil {
		t.Fatalf("seed stored mode: %v", err)
	}
	if err := ResourceVpcAssociation().Update(d, client); err != nil {
		t.Fatalf("prefix update with stored mode: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(forms) != 2 {
		t.Fatalf("ModifyExpressConnectRouterAssociationAllowedPrefix calls = %d, want 2", len(forms))
	}
	for i, form := range forms {
		if got := form.Get("AllowedPrefixesMode"); got != "IncrementalMode" {
			t.Errorf("request %d AllowedPrefixesMode = %q, want IncrementalMode", i+1, got)
		}
		if !formHasArrayValue(form, "AllowedPrefixes", "172.16.1.0/24") {
			t.Errorf("request %d did not carry the prefix list: %v", i+1, form)
		}
	}
}

func TestUnitVpcAssociationUpdateRejectsModeOnlyChange(t *testing.T) {
	silenceVpcAssociationWaits(t)
	var requests int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		t.Errorf("unexpected request %q for a mode-only update", r.Form.Get("Action"))
	})

	d := newVpcAssociationUpdateResourceData(t, map[string]interface{}{
		"ecr_id":                "ecr-1",
		"vpc_id":                "vpc-1",
		"association_region_id": "cn-hangzhou",
		"allowed_prefixes_mode": "IncrementalMode",
	})
	err := ResourceVpcAssociation().Update(d, client)
	if err == nil || !containsAll(err.Error(), "allowed_prefixes_mode", "non-empty", "allowed_prefixes") {
		t.Fatalf("mode-only update error = %v, want the local combination error", err)
	}
	if got := atomic.LoadInt32(&requests); got != 0 {
		t.Fatalf("invalid mode-only update issued %d requests, want 0", got)
	}
}

func TestUnitVpcAssociationUpdateBudgetIsShared(t *testing.T) {
	deadline := time.Unix(1_000, 0)
	cases := []struct {
		name string
		now  time.Time
		want time.Duration
	}{
		{name: "first_step_gets_full_window", now: time.Unix(990, 0), want: 10 * time.Second},
		{name: "second_step_only_gets_remainder", now: time.Unix(997, 0), want: 3 * time.Second},
		{name: "exhausted_budget_is_zero", now: time.Unix(1_005, 0), want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := remainingVpcAssociationTimeout(deadline, tc.now); got != tc.want {
				t.Fatalf("remaining timeout = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestUnitVpcAssociationPrefixModeConfigValidation(t *testing.T) {
	cases := []struct {
		name    string
		raw     cty.Value
		wantErr bool
	}{
		{
			name: "mode_with_prefix_is_valid",
			raw: cty.ObjectVal(map[string]cty.Value{
				"allowed_prefixes_mode": cty.StringVal("MatchMode"),
				"allowed_prefixes":      cty.ListVal([]cty.Value{cty.StringVal("172.16.1.0/24")}),
			}),
		},
		{
			name: "mode_with_empty_prefix_is_rejected",
			raw: cty.ObjectVal(map[string]cty.Value{
				"allowed_prefixes_mode": cty.StringVal("MatchMode"),
				"allowed_prefixes":      cty.ListValEmpty(cty.String),
			}),
			wantErr: true,
		},
		{
			name: "null_mode_is_valid",
			raw: cty.ObjectVal(map[string]cty.Value{
				"allowed_prefixes_mode": cty.NullVal(cty.String),
				"allowed_prefixes":      cty.ListValEmpty(cty.String),
			}),
		},
		{
			name: "empty_mode_is_valid",
			raw: cty.ObjectVal(map[string]cty.Value{
				"allowed_prefixes_mode": cty.StringVal(""),
				"allowed_prefixes":      cty.ListValEmpty(cty.String),
			}),
		},
		{
			name: "unknown_prefix_is_skipped",
			raw: cty.ObjectVal(map[string]cty.Value{
				"allowed_prefixes_mode": cty.StringVal("IncrementalMode"),
				"allowed_prefixes":      cty.UnknownVal(cty.List(cty.String)),
			}),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateVpcAssociationPrefixModeRawConfig(tc.raw)
			if tc.wantErr && err == nil {
				t.Fatal("expected the local combination error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestUnitDescribeVpcAssociationTargetIdentity(t *testing.T) {
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		if r.Form.Get("Action") != "DescribeExpressConnectRouterAssociation" {
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
		writeVpcAssociationJSON(t, w, http.StatusOK, map[string]interface{}{
			"RequestId": "unit-test",
			"AssociationList": []map[string]interface{}{
				{
					"AssociationId":       "assoc-a",
					"AssociationNodeType": "VPC",
					"EcrId":               "ecr-1",
					"VpcId":               "vpc-a",
					"Status":              "ACTIVE",
					"Description":         "same-description",
				},
				{
					"AssociationId":       "assoc-b",
					"AssociationNodeType": "VPC",
					"EcrId":               "ecr-1",
					"VpcId":               "vpc-b",
					"Status":              "ACTIVE",
					"Description":         "same-description",
				},
			},
		})
	})
	service := NewService(client)

	object, err := service.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation("ecr-1:assoc-b:vpc-b")
	if err != nil {
		t.Fatalf("describe association B: %v", err)
	}
	if got := object["AssociationId"]; got != "assoc-b" {
		t.Fatalf("describe B returned association %v, want assoc-b", got)
	}

	if _, err := service.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation("ecr-1:assoc-c:vpc-c"); !errs.NotFoundError(err) {
		t.Fatalf("missing association C error = %v, want NotFound even though A exists", err)
	}
}

func formHasArrayValue(form url.Values, key, want string) bool {
	prefix := key + "."
	for formKey, values := range form {
		if formKey != key && !strings.HasPrefix(formKey, prefix) {
			continue
		}
		for _, value := range values {
			if value == want {
				return true
			}
		}
	}
	return false
}

func containsAll(value string, substrings ...string) bool {
	for _, substring := range substrings {
		if !strings.Contains(value, substring) {
			return false
		}
	}
	return true
}

// TestUnitVpcAssociationPrefixOrderSuppress is the bounded net for the
// allowed_prefixes suppressor: an order-only change must be ignored while every
// membership change must still surface as a diff, so an over-broad suppressor
// cannot silently swallow real updates.
func TestUnitVpcAssociationPrefixOrderSuppress(t *testing.T) {
	cases := []struct {
		name     string
		oldValue string
		newValue string
		expected bool
	}{
		{"ReorderOnly", "172.16.1.0/24,172.16.2.0/24", "172.16.2.0/24,172.16.1.0/24", true},
		{"LongerReorder", "172.16.1.0/24,172.16.2.0/24,172.16.3.0/24", "172.16.3.0/24,172.16.1.0/24,172.16.2.0/24", true},
		{"MemberAdded", "172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24,172.16.2.0/24,172.16.3.0/24", false},
		{"MemberRemoved", "172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24", false},
		{"MemberReplaced", "172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24,172.16.9.0/24", false},
		{"MultiplicityChanged", "172.16.1.0/24,172.16.1.0/24,172.16.2.0/24", "172.16.1.0/24,172.16.2.0/24,172.16.2.0/24", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := newVpcAssociationPrefixDiffData(t, tc.oldValue, tc.newValue)
			if got := suppressVpcAssociationPrefixReorder("allowed_prefixes.#", tc.oldValue, tc.newValue, d); got != tc.expected {
				t.Fatalf("suppress(%q, %q) = %t, want %t", tc.oldValue, tc.newValue, got, tc.expected)
			}
		})
	}

	if suppressVpcAssociationPrefixReorder("allowed_prefixes.#", "172.16.1.0/24", "172.16.1.0/24", nil) {
		t.Fatal("a nil ResourceData carries no old/new pair and must not suppress")
	}
	if _, ok := vpcAssociationPrefixStrings([]interface{}{"172.16.1.0/24", 42}); ok {
		t.Fatal("a list with a non-string member must not be treated as a known prefix list")
	}
	if _, ok := vpcAssociationPrefixStrings([]interface{}{}); ok {
		t.Fatal("an empty list must not be treated as a known prefix list")
	}
}

// newVpcAssociationPrefixDiffData rebuilds the planner's old/new pair so the
// suppressor is exercised with the same ResourceData it sees during a real plan
// instead of a hand-assembled stand-in.
func newVpcAssociationPrefixDiffData(t *testing.T, oldValue, newValue string) *schema.ResourceData {
	t.Helper()
	r := &schema.Resource{Schema: map[string]*schema.Schema{
		"allowed_prefixes": {
			Type:             schema.TypeList,
			Optional:         true,
			Elem:             &schema.Schema{Type: schema.TypeString},
			DiffSuppressFunc: suppressVpcAssociationPrefixReorder,
		},
	}}
	oldMembers := strings.Split(oldValue, ",")
	state := &terraform.InstanceState{ID: "ecr-1:vpc-1", Attributes: map[string]string{
		"allowed_prefixes.#": strconv.Itoa(len(oldMembers)),
	}}
	for i, member := range oldMembers {
		state.Attributes[fmt.Sprintf("allowed_prefixes.%d", i)] = member
	}
	newMembers := make([]interface{}, 0, len(strings.Split(newValue, ",")))
	for _, member := range strings.Split(newValue, ",") {
		newMembers = append(newMembers, member)
	}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]interface{}{
		"allowed_prefixes": newMembers,
	}), nil)
	if err != nil {
		t.Fatalf("build prefix diff: %v", err)
	}
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("build prefix resource data: %v", err)
	}
	return d
}
