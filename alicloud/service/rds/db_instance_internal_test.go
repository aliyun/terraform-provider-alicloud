package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestUnitParseDBInstanceParamGroupId is a regression test for the provider
// panic reported when creating alicloud_db_instance.
//
// resourceAliCloudDBInstanceRead used to run, without a comma-ok guard:
//
//	dbParamGroupInfo := response["ParamGroupInfo"].(map[string]interface{})
//	...
//	d.Set("db_param_group_id", dbParamGroupInfo["ParamGroupId"].(string))
//
// A freshly created instance whose DescribeParameters response does not yet
// contain "ParamGroupInfo" therefore crashed the plugin with:
//
//	interface conversion: interface {} is nil, not map[string]interface {}
//
// The parsing is now isolated in parseDBInstanceParamGroupId, which must never
// panic and must report presence via its second return value. The first two
// cases below are the crash inputs (they panic against the old inline code).
func TestUnitParseDBInstanceParamGroupId(t *testing.T) {
	cases := []struct {
		name      string
		response  map[string]interface{}
		wantID    string
		wantFound bool
	}{
		{
			name:     "ParamGroupInfo absent (freshly created instance)",
			response: map[string]interface{}{},
		},
		{
			name:     "ParamGroupInfo present but nil",
			response: map[string]interface{}{"ParamGroupInfo": nil},
		},
		{
			name:     "ParamGroupInfo present but ParamGroupId missing",
			response: map[string]interface{}{"ParamGroupInfo": map[string]interface{}{}},
		},
		{
			name:     "ParamGroupId present but not a string",
			response: map[string]interface{}{"ParamGroupInfo": map[string]interface{}{"ParamGroupId": 123}},
		},
		{
			name:      "ParamGroupInfo and ParamGroupId both present",
			response:  map[string]interface{}{"ParamGroupInfo": map[string]interface{}{"ParamGroupId": "rpg-test123"}},
			wantID:    "rpg-test123",
			wantFound: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, ok := ParseDBInstanceParamGroupId(c.response)
			if ok != c.wantFound {
				t.Fatalf("parseDBInstanceParamGroupId(%v) found = %v, want %v", c.response, ok, c.wantFound)
			}
			if id != c.wantID {
				t.Fatalf("parseDBInstanceParamGroupId(%v) id = %q, want %q", c.response, id, c.wantID)
			}
		})
	}
}

type fakeFieldChangeReporter map[string]bool

func (f fakeFieldChangeReporter) HasChange(field string) bool {
	return f[field]
}

func TestUnitDBInstanceCreateOnlyFieldsUnchanged(t *testing.T) {
	if err := validateDBInstanceCreateOnlyFieldsUnchanged(fakeFieldChangeReporter{}); err != nil {
		t.Fatalf("expected no error when no create-only field changes, got %v", err)
	}
	for _, field := range dbInstanceCreateOnlyFields {
		err := validateDBInstanceCreateOnlyFieldsUnchanged(fakeFieldChangeReporter{field: true})
		if err == nil {
			t.Fatalf("expected an error when %s changes on an existing instance", field)
		}
		if !strings.Contains(err.Error(), "`"+field+"`") {
			t.Fatalf("error %q does not name %s", err, field)
		}
	}
}

// TestUnitDBInstanceBoolFromConfig pins the replacement for the deprecated
// GetOkExists on auto_create_proxy/auto_pay: an explicit false must be sent while
// an omitted or unknown value must not be sent at all.
func TestUnitDBInstanceBoolFromConfig(t *testing.T) {
	cases := []struct {
		name      string
		config    cty.Value
		expected  bool
		present   bool
		attribute string
	}{
		{name: "explicit false", config: cty.ObjectVal(map[string]cty.Value{"auto_pay": cty.False}), attribute: "auto_pay", expected: false, present: true},
		{name: "explicit true", config: cty.ObjectVal(map[string]cty.Value{"auto_pay": cty.True}), attribute: "auto_pay", expected: true, present: true},
		{name: "null value", config: cty.ObjectVal(map[string]cty.Value{"auto_pay": cty.NullVal(cty.Bool)}), attribute: "auto_pay"},
		{name: "unknown value", config: cty.ObjectVal(map[string]cty.Value{"auto_pay": cty.UnknownVal(cty.Bool)}), attribute: "auto_pay"},
		{name: "absent attribute", config: cty.ObjectVal(map[string]cty.Value{"auto_create_proxy": cty.True}), attribute: "auto_pay"},
		{name: "nil config", config: cty.NilVal, attribute: "auto_pay"},
		{name: "null config", config: cty.NullVal(cty.EmptyObject), attribute: "auto_pay"},
	}
	for _, c := range cases {
		got, present := dbInstanceBoolFromConfig(c.config, c.attribute)
		if present != c.present || got != c.expected {
			t.Fatalf("case %s: dbInstanceBoolFromConfig() = (%v, %v), want (%v, %v)", c.name, got, present, c.expected, c.present)
		}
	}
}

// TestUnitNonEmptyStringAttr covers the Read guard used for
// Optional+Computed attributes: an omitted, null or empty API value must not
// overwrite the last known state.
func TestUnitNonEmptyStringAttr(t *testing.T) {
	cases := []struct {
		name     string
		instance map[string]interface{}
		expected string
		present  bool
	}{
		{name: "present", instance: map[string]interface{}{"ConnectionMode": "Safe"}, expected: "Safe", present: true},
		{name: "absent", instance: map[string]interface{}{}},
		{name: "null", instance: map[string]interface{}{"ConnectionMode": nil}},
		{name: "empty", instance: map[string]interface{}{"ConnectionMode": ""}},
	}
	for _, c := range cases {
		got, present := nonEmptyStringAttr(c.instance, "ConnectionMode")
		if got != c.expected || present != c.present {
			t.Fatalf("case %s: nonEmptyStringAttr() = (%q, %v), want (%q, %v)", c.name, got, present, c.expected, c.present)
		}
	}
}

// TestUnitDBInstancePgBouncerEnabledFromInstance covers the Read conversion of
// the PgBouncer flag: only an explicit "true"/"false" from the API is trusted,
// while an omitted or unparsable value is reported as false instead of keeping
// an unconfirmed configured true.
func TestUnitDBInstancePgBouncerEnabledFromInstance(t *testing.T) {
	cases := []struct {
		name     string
		instance map[string]interface{}
		want     bool
	}{
		{name: "enabled", instance: map[string]interface{}{"PGBouncerEnabled": "true"}, want: true},
		{name: "disabled", instance: map[string]interface{}{"PGBouncerEnabled": "false"}},
		{name: "boolean true", instance: map[string]interface{}{"PGBouncerEnabled": true}, want: true},
		{name: "omitted", instance: map[string]interface{}{}},
		{name: "null", instance: map[string]interface{}{"PGBouncerEnabled": nil}},
		{name: "unparsable", instance: map[string]interface{}{"PGBouncerEnabled": "unsupported"}},
	}
	for _, c := range cases {
		if got := pgBouncerEnabledFromInstance(c.instance); got != c.want {
			t.Fatalf("case %s: pgBouncerEnabledFromInstance() = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestUnitDBInstanceUserBackupIdValidation(t *testing.T) {
	legal := map[string]interface{}{
		"engine":               "MySQL",
		"engine_version":       "5.7",
		"instance_charge_type": "Postpaid",
		"category":             "Basic",
	}
	cases := []struct {
		name    string
		raw     map[string]interface{}
		wantErr bool
	}{
		{name: "legal combination", raw: legal},
		{name: "prepaid", raw: copyWith(legal, "instance_charge_type", "Prepaid"), wantErr: true},
		{name: "non MySQL", raw: copyWith(legal, "engine", "PostgreSQL"), wantErr: true},
		{name: "non 5.7", raw: copyWith(legal, "engine_version", "8.0"), wantErr: true},
		{name: "non Basic category", raw: copyWith(legal, "category", "HighAvailability"), wantErr: true},
		{name: "category unset", raw: copyWith(legal, "category", ""), wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAliCloudDBInstance().Schema, c.raw)
			err := validateDBInstanceUserBackupId(d)
			if (err != nil) != c.wantErr {
				t.Fatalf("validateDBInstanceUserBackupId() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func copyWith(src map[string]interface{}, key string, value interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	dst[key] = value
	return dst
}

func TestUnitDBInstanceSQLServerMajorVersionUpgradeRequest(t *testing.T) {
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou", SourceIp: "127.0.0.1"}
	base := map[string]interface{}{
		"engine":                      "SQLServer",
		"engine_version":              "2019_std_ha",
		"instance_type":               "mssql.x4.medium.s2",
		"instance_charge_type":        "Postpaid",
		"category":                    "Basic",
		"zone_id":                     "cn-hangzhou-k",
		"vswitch_id":                  "vsw-primary",
		"allow_major_version_upgrade": true,
	}

	build := func(raw map[string]interface{}) (map[string]interface{}, error) {
		d := schema.TestResourceDataRaw(t, resourceAliCloudDBInstance().Schema, raw)
		d.SetId("rm-test123")
		return buildDBInstanceSQLServerMajorVersionUpgradeRequest(d, client)
	}

	t.Run("basic category carries the required parameters", func(t *testing.T) {
		request, err := build(base)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for key, want := range map[string]interface{}{
			"Category":        "Basic",
			"ZoneId":          "cn-hangzhou-k",
			"VSwitchId":       "vsw-primary",
			"EngineVersion":   "2019_std_ha",
			"DBInstanceId":    "rm-test123",
			"DBInstanceClass": "mssql.x4.medium.s2",
		} {
			if got := request[key]; got != want {
				t.Fatalf("request[%s] = %v, want %v", key, got, want)
			}
		}
		if _, ok := request["ZoneIdSlave1"]; ok {
			t.Fatalf("ZoneIdSlave1 must not be sent for a Basic category instance")
		}
		if token, _ := request["ClientToken"].(string); token == "" {
			t.Fatalf("ClientToken must be set for an idempotent retry")
		}
	})

	t.Run("explicit false is forwarded", func(t *testing.T) {
		request, err := build(copyWith(base, "allow_major_version_upgrade", false))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, ok := request["AllowMajorVersionUpgrade"]
		if !ok || got != false {
			t.Fatalf("AllowMajorVersionUpgrade = %v (present %v), want false", got, ok)
		}
	})

	t.Run("high availability without a slave zone stays single AZ", func(t *testing.T) {
		request, err := build(copyWith(base, "category", "HighAvailability"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := request["ZoneIdSlave1"]; got != "cn-hangzhou-k" {
			t.Fatalf("ZoneIdSlave1 = %v, want the primary zone for a single-AZ deployment", got)
		}
	})

	t.Run("high availability with a different slave zone needs two vswitches", func(t *testing.T) {
		raw := copyWith(copyWith(base, "category", "HighAvailability"), "zone_id_slave_a", "cn-hangzhou-j")
		if _, err := build(raw); err == nil {
			t.Fatalf("expected an error when the slave zone differs but only one vSwitch is configured")
		}
		raw = copyWith(raw, "vswitch_id", "vsw-primary,vsw-standby")
		request, err := build(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := request["ZoneIdSlave1"]; got != "cn-hangzhou-j" {
			t.Fatalf("ZoneIdSlave1 = %v, want cn-hangzhou-j", got)
		}
		if got := request["VSwitchId"]; got != "vsw-primary,vsw-standby" {
			t.Fatalf("VSwitchId = %v, want both vSwitch IDs", got)
		}
	})

	for _, tc := range []struct {
		name  string
		raw   map[string]interface{}
		field string
	}{
		{name: "missing category", raw: copyWith(base, "category", ""), field: "category"},
		{name: "missing zone", raw: copyWith(base, "zone_id", ""), field: "zone_id"},
		{name: "missing vswitch", raw: copyWith(base, "vswitch_id", ""), field: "vswitch_id"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := build(tc.raw); err == nil || !strings.Contains(err.Error(), "`"+tc.field+"` is required") {
				t.Fatalf("expected a %s required error, got %v", tc.field, err)
			}
		})
	}
}

func TestUnitDBInstanceConcurrentErrorClassification(t *testing.T) {
	// The codes are listed explicitly (not taken from the production slices) so
	// the test fails when a documented wait-and-retry code is dropped, and so the
	// create-only / modify-only split is asserted in both directions.
	shared := []string{
		"InvalidConcurrentOperate",
		"System.ConcurrentConflict",
		"OperationDenied.SystemConcurrent",
	}
	for _, code := range shared {
		err := &tea.SDKError{Code: tea.String(code), Message: tea.String("concurrent"), Data: tea.String("")}
		if !isDBInstanceConcurrentError(err) || !isDBInstanceCreateWaitError(err) || !isDBInstanceModifyWaitError(err) {
			t.Fatalf("expected %s to be classified on every write surface", code)
		}
	}

	createOnly := &tea.SDKError{Code: tea.String("Order.UnfinishedRefundOrder"), Message: tea.String("refund"), Data: tea.String("")}
	if !isDBInstanceCreateWaitError(createOnly) {
		t.Fatal("Order.UnfinishedRefundOrder must be retried on create")
	}
	if isDBInstanceModifyWaitError(createOnly) {
		t.Fatal("Order.UnfinishedRefundOrder must not be retried on modify")
	}

	modifyOnly := &tea.SDKError{Code: tea.String("OperationDenied.DurationLimit"), Message: tea.String("duration"), Data: tea.String("")}
	if !isDBInstanceModifyWaitError(modifyOnly) {
		t.Fatal("OperationDenied.DurationLimit must be retried on modify")
	}
	if isDBInstanceCreateWaitError(modifyOnly) {
		t.Fatal("OperationDenied.DurationLimit must not be retried on create")
	}

	permanent := &tea.SDKError{Code: tea.String("InvalidDBInstanceId.NotFound"), Message: tea.String("not found"), Data: tea.String("")}
	if isDBInstanceCreateWaitError(permanent) || isDBInstanceModifyWaitError(permanent) || isDBInstanceConcurrentError(permanent) {
		t.Fatal("a permanent error must not be classified as retryable")
	}
	if isDBInstanceCreateWaitError(nil) || isDBInstanceModifyWaitError(nil) {
		t.Fatal("nil must not be classified as retryable")
	}
}

// TestUnitDBInstanceTemplateIdListOrderSuppress exercises the bounded-suppressor
// unit net for template_id_list. The AttachWhitelistTemplate/Detach operations
// use set-based comparison, so order-only changes are semantically meaningless.
// The cases run the same algorithm the SDK reaches through the DiffSuppressFunc
// (templateIdListEqualIgnoringOrder with the values from d.GetChange), so the
// production path is verified.
func TestUnitDBInstanceTemplateIdListOrderSuppress(t *testing.T) {
	listCases := []struct {
		name     string
		old      []interface{}
		new      []interface{}
		expected bool
	}{
		{"ReorderOnly", []interface{}{1, 2}, []interface{}{2, 1}, true},
		{"ReorderOnlyLonger", []interface{}{1, 2, 3}, []interface{}{3, 1, 2}, true},
		{"MemberAdded", []interface{}{1, 2}, []interface{}{1, 2, 3}, false},
		{"MemberRemoved", []interface{}{1, 2, 3}, []interface{}{1, 2}, false},
		{"MemberReplaced", []interface{}{1, 2}, []interface{}{1, 3}, false},
		{"DuplicateCollapses", []interface{}{1, 1}, []interface{}{1}, false},
		{"DuplicateReorder", []interface{}{1, 1, 2}, []interface{}{2, 1, 1}, true},
		{"SameOrder", []interface{}{1, 2}, []interface{}{1, 2}, true},
		{"Empty", []interface{}{}, []interface{}{}, true},
	}
	for _, tc := range listCases {
		if got := templateIdListEqualIgnoringOrder(tc.old, tc.new); got != tc.expected {
			t.Fatalf("list case %s: got %v, expected %v", tc.name, got, tc.expected)
		}
	}
}

// dbInstanceMockClient builds an AliyunClient whose RDS endpoint is an isolated
// in-process HTTP server, so request-level tests exercise the same code path
// Terraform uses without touching the cloud.
func dbInstanceMockClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("rds", endpoint)
	return client
}

func dbInstanceRawConfig(values map[string]interface{}) cty.Value {
	attrs := map[string]cty.Value{}
	for key, value := range values {
		switch typed := value.(type) {
		case bool:
			attrs[key] = cty.BoolVal(typed)
		case int:
			attrs[key] = cty.NumberIntVal(int64(typed))
		case string:
			attrs[key] = cty.StringVal(typed)
		default:
			attrs[key] = cty.NullVal(cty.DynamicPseudoType)
		}
	}
	return cty.ObjectVal(attrs)
}

func dbInstanceMockCreateConfig(extra map[string]interface{}) map[string]interface{} {
	config := map[string]interface{}{
		"engine":                   "MySQL",
		"engine_version":           "8.0",
		"instance_type":            "mysql.n2.small.2c",
		"instance_storage":         20,
		"category":                 "Basic",
		"instance_charge_type":     "Postpaid",
		"db_instance_storage_type": "cloud_ssd",
	}
	for key, value := range extra {
		config[key] = value
	}
	return config
}

// TestUnitDBInstanceCreateWireWriteOnlyBooleans drives the public Create entry
// point against an isolated HTTP endpoint and asserts the serialized
// CreateDBInstance form. An omitted value must not be sent (the API default is
// true), while an explicit false/true must reach the wire unchanged.
func TestUnitDBInstanceCreateWireWriteOnlyBooleans(t *testing.T) {
	cases := []struct {
		name        string
		extra       map[string]interface{}
		want        map[string]string
		wantMissing []string
	}{
		{
			name:        "omitted",
			extra:       map[string]interface{}{},
			wantMissing: []string{"AutoPay", "AutoCreateProxy"},
		},
		{
			name:  "explicit false",
			extra: map[string]interface{}{"auto_pay": false, "auto_create_proxy": false},
			want:  map[string]string{"AutoPay": "false", "AutoCreateProxy": "false"},
		},
		{
			name:  "explicit true",
			extra: map[string]interface{}{"auto_pay": true, "auto_create_proxy": true},
			want:  map[string]string{"AutoPay": "true", "AutoCreateProxy": "true"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := dbInstanceMockCreateConfig(tc.extra)
			var captured url.Values
			client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
				if action := r.Form.Get("Action"); action != "CreateDBInstance" {
					t.Errorf("unexpected action %q", action)
				}
				captured = r.Form
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Code":    "InvalidParameter",
					"Message": "stop after request capture",
				})
			})

			res := resourceAliCloudDBInstance()
			diff, err := res.Diff(context.Background(), nil, sdkterraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatalf("diff failed: %v", err)
			}
			if diff == nil {
				t.Fatal("expected a create diff")
			}
			diff.RawConfig = dbInstanceRawConfig(config)
			if _, diags := res.Apply(context.Background(), nil, diff, client); !diags.HasError() {
				t.Fatal("expected the mocked CreateDBInstance failure")
			}
			if captured == nil {
				t.Fatal("CreateDBInstance was not called")
			}
			for key, want := range tc.want {
				if got := captured.Get(key); got != want {
					t.Fatalf("%s = %q, want %q", key, got, want)
				}
			}
			for _, key := range tc.wantMissing {
				if _, ok := captured[key]; ok {
					t.Fatalf("%s must not be sent when omitted, got %q", key, captured.Get(key))
				}
			}
		})
	}
}

// TestUnitDBInstanceCreateWireWriteOnlyStrings drives the public Create entry
// point against an isolated HTTP endpoint and asserts that the three string
// write-only parameters are serialized into the CreateDBInstance form when they
// are configured. It proves request construction only: the test account has no
// valid business info, coupon or user backup, so cloud-side acceptance and
// read-back of these values remain unverified (see the implementation report).
func TestUnitDBInstanceCreateWireWriteOnlyStrings(t *testing.T) {
	config := dbInstanceMockCreateConfig(map[string]interface{}{
		"engine_version": "5.7",
		"category":       "Basic",
		"business_info":  "121436975448952",
		"promotion_code": "tf-promotion-code",
		"user_backup_id": "backup-121436975448952",
	})
	var captured url.Values
	client := dbInstanceMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		if action := r.Form.Get("Action"); action != "CreateDBInstance" {
			t.Errorf("unexpected action %q", action)
		}
		captured = r.Form
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code":    "InvalidParameter",
			"Message": "stop after request capture",
		})
	})

	res := resourceAliCloudDBInstance()
	diff, err := res.Diff(context.Background(), nil, sdkterraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatalf("diff failed: %v", err)
	}
	if diff == nil {
		t.Fatal("expected a create diff")
	}
	diff.RawConfig = dbInstanceRawConfig(config)
	if _, diags := res.Apply(context.Background(), nil, diff, client); !diags.HasError() {
		t.Fatal("expected the mocked CreateDBInstance failure")
	}
	if captured == nil {
		t.Fatal("CreateDBInstance was not called")
	}
	for key, want := range map[string]string{
		"BusinessInfo":  "121436975448952",
		"PromotionCode": "tf-promotion-code",
		"UserBackupId":  "backup-121436975448952",
	} {
		if got := captured.Get(key); got != want {
			t.Fatalf("%s = %q, want %q", key, got, want)
		}
	}
}

// TestUnitDBInstanceCreateOnlyFieldsPlanGuard drives the real CustomizeDiff
// path with a prior state, so the create-only rejection is proven for the
// boolean write-only fields and not only through a fake HasChange reporter.
func TestUnitDBInstanceCreateOnlyFieldsPlanGuard(t *testing.T) {
	cases := []struct {
		field string
		from  interface{}
		to    interface{}
	}{
		{field: "auto_create_proxy", from: false, to: true},
		{field: "auto_pay", from: false, to: true},
		{field: "business_info", from: "", to: "info-1"},
		{field: "promotion_code", from: "", to: "CODE-1"},
		{field: "user_backup_id", from: "", to: "backup-1"},
	}

	baseAttributes := func(field string, value interface{}) map[string]string {
		return map[string]string{
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"instance_type":            "mysql.n2.small.2c",
			"instance_storage":         "20",
			"category":                 "Basic",
			"instance_charge_type":     "Postpaid",
			"db_instance_storage_type": "cloud_ssd",
			"connection_mode":          "Standard",
			field:                      fmt.Sprint(value),
		}
	}

	for _, tc := range cases {
		t.Run(tc.field+"/in_place_rejected", func(t *testing.T) {
			config := dbInstanceMockCreateConfig(map[string]interface{}{tc.field: tc.to})
			state := &sdkterraform.InstanceState{ID: "rm-existing", Attributes: baseAttributes(tc.field, tc.from)}

			res := resourceAliCloudDBInstance()
			_, err := res.Diff(context.Background(), state, sdkterraform.NewResourceConfigRaw(config), nil)
			if err == nil || !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("expected a plan-time rejection for %s, got %v", tc.field, err)
			}
		})

		t.Run(tc.field+"/replacement_allowed", func(t *testing.T) {
			// An independent ForceNew change replaces the instance, so the
			// create-only parameter is applied to the new instance and the plan
			// must not be rejected.
			config := dbInstanceMockCreateConfig(map[string]interface{}{
				tc.field:          tc.to,
				"connection_mode": "Safe",
			})
			state := &sdkterraform.InstanceState{ID: "rm-existing", Attributes: baseAttributes(tc.field, tc.from)}

			res := resourceAliCloudDBInstance()
			if _, err := res.Diff(context.Background(), state, sdkterraform.NewResourceConfigRaw(config), nil); err != nil {
				t.Fatalf("a replacement plan must accept the create-only change on %s, got %v", tc.field, err)
			}
		})
	}
}

// TestUnitDBInstanceForceNewFieldListMatchesSchema keeps the replacement list in
// sync with the schema: a ForceNew attribute missing from the list would bring
// back the create-only false rejection, and a listed field that lost ForceNew
// would silently bypass the guard.
func TestUnitDBInstanceForceNewFieldListMatchesSchema(t *testing.T) {
	schemaMap := resourceAliCloudDBInstance().Schema
	listed := make(map[string]bool, len(dbInstanceForceNewFields))
	for _, field := range dbInstanceForceNewFields {
		if _, ok := schemaMap[field]; !ok {
			t.Fatalf("listed ForceNew field %q does not exist in the schema", field)
		}
		listed[field] = true
	}
	for name, item := range schemaMap {
		if item.ForceNew != listed[name] {
			t.Fatalf("ForceNew mismatch for %q: schema=%v list=%v", name, item.ForceNew, listed[name])
		}
	}
}
