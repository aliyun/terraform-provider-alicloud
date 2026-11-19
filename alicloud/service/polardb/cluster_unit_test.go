package polardb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestUnitPolarDBDistributedNodeSchema(t *testing.T) {
	clusterSchema := ResourceCluster().Schema
	for _, key := range []string{"cn_node_class", "dn_node_class", "cn_node_num", "dn_node_num"} {
		field := clusterSchema[key]
		if !field.Optional || !field.Computed || field.ForceNew {
			t.Fatalf("%s must be Optional and Computed without ForceNew", key)
		}
	}
	for _, key := range []string{"cn_node_num", "dn_node_num"} {
		if clusterSchema[key].Type != schema.TypeInt {
			t.Fatalf("%s must use TypeInt", key)
		}
	}
	if len(clusterSchema["db_node_count"].ConflictsWith) != 2 {
		t.Fatal("db_node_count must conflict with distributed node count fields")
	}
	for _, key := range []string{"modify_dist_nodes", "delete_sub_groups", "restart_cluster"} {
		if _, ok := clusterSchema[key]; ok {
			t.Fatalf("command-style field %s must not be part of desired state", key)
		}
	}
}
func TestUnitClassifyPolarDBDistributedNode(t *testing.T) {
	tests := []struct {
		name     string
		node     map[string]interface{}
		expected string
		wantErr  bool
	}{
		{name: "compute type", node: map[string]interface{}{"DBNodeId": "cn-1", "DBNodeType": "Compute"}, expected: "cn"},
		{name: "data type", node: map[string]interface{}{"DBNodeId": "dn-1", "DBNodeType": "Data"}, expected: "dn"},
		{name: "compute subgroup", node: map[string]interface{}{"DBNodeId": "cn-0", "SubGroupType": "Compute"}, expected: "cn"},
		{name: "data subgroup", node: map[string]interface{}{"DBNodeId": "dn-0", "SubGroupType": "Data"}, expected: "dn"},
		{name: "CN role", node: map[string]interface{}{"DBNodeId": "cn-2", "DBNodeRole": "CN"}, expected: "cn"},
		{name: "DN role", node: map[string]interface{}{"DBNodeId": "dn-2", "DBNodeRole": "DN"}, expected: "dn"},
		{name: "legacy CN", node: map[string]interface{}{"DBNodeId": "cn-3", "DBNodeClass": "polar.pg.x4.medium", "IsPrimaryCN": false}, expected: "cn"},
		{name: "legacy CN reader", node: map[string]interface{}{"DBNodeId": "cn-4", "DBNodeClass": "polar.pg.x4.medium", "DBNodeRole": "Reader", "IsPrimaryCN": false}, expected: "cn"},
		{name: "legacy DN", node: map[string]interface{}{"DBNodeId": "dn-3", "DBNodeClass": "polar.pg.x4.medium"}, expected: "dn"},
		{name: "unknown role", node: map[string]interface{}{"DBNodeId": "node-1", "DBNodeClass": "polar.pg.x4.medium", "DBNodeRole": "Reader"}, wantErr: true},
		{name: "missing role", node: map[string]interface{}{"DBNodeId": "node-2"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := classifyPolarDBDistributedNode(tt.node)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected classification error, got %q", actual)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected classification error: %s", err)
			}
			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
func TestUnitValidatePolarDBDistributedNodeCounts(t *testing.T) {
	clusterInfo := map[string]interface{}{"CnNodeCount": 1, "DnNodeCount": 2}
	if err := validatePolarDBDistributedNodeCounts(clusterInfo, 1, 2); err != nil {
		t.Fatalf("expected subgroup counts to match logical counts: %s", err)
	}
	if err := validatePolarDBDistributedNodeCounts(clusterInfo, 3, 0); err == nil {
		t.Fatal("expected silently misclassified nodes to be rejected")
	}
}
func TestUnitSelectPolarDBDistributedSubGroupsForDeletion(t *testing.T) {
	nodes := []map[string]interface{}{
		{"DBNodeId": "pi-primary-writer", "SubGroupName": "pc-test", "CreationTime": "2026-01-01T00:00:00Z"},
		{"DBNodeId": "pi-primary-reader", "SubGroupName": "pc-test", "CreationTime": "2026-01-01T00:00:00Z"},
		{"DBNodeId": "pi-extra-writer", "SubGroupName": "extra-cn", "CreationTime": "2026-01-02T00:00:00Z"},
		{"DBNodeId": "pi-extra-reader", "SubGroupName": "extra-cn", "CreationTime": "2026-01-02T00:00:00Z"},
		{"DBNodeId": "pi-newest", "SubGroupName": "newest-cn", "CreationTime": "2026-01-03T00:00:00Z"},
		{"DBNodeId": "pi-missing-subgroup"},
	}
	selected := selectPolarDBDistributedSubGroupsForDeletion(nodes, "cn", "pc-test", 1)
	if len(selected) != 1 || selected[0] != "newest-cn" {
		t.Fatalf("expected the newest non-primary logical subgroup, got %#v", selected)
	}
	selected = selectPolarDBDistributedSubGroupsForDeletion(nodes, "dn", "pc-test", 2)
	if len(selected) != 2 || selected[0] != "newest-cn" || selected[1] != "extra-cn" {
		t.Fatalf("expected deterministic DN selection without the primary or missing subgroup, got %#v", selected)
	}
}
func TestUnitPolarDBParameterModificationPendingStates(t *testing.T) {
	states := []string{"Rebooting", "Running"}
	nextState := 0
	stateConf := helper.BuildStateConf(
		polarDBParameterModificationPendingStates(),
		[]string{"Running"},
		time.Second,
		0,
		func() (interface{}, string, error) {
			state := states[nextState]
			nextState++
			return state, state, nil
		},
	)
	stateConf.PollInterval = time.Millisecond

	result, err := stateConf.WaitForState()
	if err != nil {
		t.Fatalf("expected Rebooting to be treated as pending: %s", err)
	}
	if result != "Running" {
		t.Fatalf("expected waiter to finish in Running, got %v", result)
	}
}

func TestUnitPolarDBClusterEncryptNewTablesSchema(t *testing.T) {
	field := ResourceCluster().Schema["encrypt_new_tables"]
	if field == nil {
		t.Fatal("encrypt_new_tables must be defined in alicloud_polardb_cluster schema")
	}
	// Read unconditionally stores the value reported by DescribeDBClusterTDE, including for
	// configurations that never set this attribute; Computed keeps that value from planning a diff.
	if !field.Optional || !field.Computed {
		t.Fatalf("encrypt_new_tables must be Optional and Computed, got Optional=%v Computed=%v", field.Optional, field.Computed)
	}
}

// polarDBClusterTDEHTTPClient returns a client whose polardb endpoint is a local server. The server
// records the serialized ModifyDBClusterTDE request and rejects it, so the update returns before the
// provider's hardcoded three minute TDE state wait; assertions run on the request that was sent.
func polarDBClusterTDEHTTPClient(t *testing.T, capture func(url.Values)) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %s", err)
		}
		if action := r.Form.Get("Action"); action != "ModifyDBClusterTDE" {
			t.Errorf("TDE-only update sent unexpected action %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"UnexpectedAction","Message":"unexpected action"}`))
			return
		}
		capture(r.Form)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unit test stops the update after the TDE request"}`))
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build dummy credential: %s", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatalf("build client: %s", err)
	}
	endpoints.Store("polardb", endpoint)
	return client
}

func TestUnitPolarDBClusterTDEEncryptNewTablesRequest(t *testing.T) {
	// Each case starts from a real configuration, optionally lets Read store the
	// returned value, and diffs against a second configuration. The provider only
	// sees state plus the new configuration, so "never configured" and "removed
	// from configuration" share the absent-attribute shape and differ by the value
	// Read stored before the update.
	baseConfig := func() map[string]interface{} {
		return map[string]interface{}{
			"db_type":                   "MySQL",
			"db_version":                "8.0",
			"db_node_class":             "polar.mysql.x4.large",
			"vswitch_id":                "vsw-tde-unit",
			"tde_status":                "Enabled",
			"enable_automatic_rotation": false,
		}
	}
	cases := []struct {
		name       string
		configured string
		readBack   string
		next       string
		wantSent   bool
		wantValue  string
	}{
		// Attribute never appeared in configuration; DescribeDBClusterTDE supplied OFF.
		{"never configured omits read back default", "", "OFF", "", false, ""},
		// Configuration changes the value, so it must be sent.
		{"explicit value change is sent", "OFF", "OFF", "ON", true, "ON"},
		// Configuration keeps the same value while another TDE parameter changes.
		{"explicit same value is omitted", "ON", "ON", "ON", false, ""},
		// Attribute was configured, then removed from configuration; Read keeps ON.
		{"removed from configuration omits read back value", "ON", "ON", "", false, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clusterResource := ResourceCluster()
			previousConfig := baseConfig()
			if tc.configured != "" {
				previousConfig["encrypt_new_tables"] = tc.configured
			}
			previous := schema.TestResourceDataRaw(t, clusterResource.Schema, previousConfig)
			previous.SetId("pc-polar-tde-unit")
			if tc.readBack != "" {
				if err := previous.Set("encrypt_new_tables", tc.readBack); err != nil {
					t.Fatalf("store read back value: %s", err)
				}
			}
			state := previous.State()

			nextConfig := baseConfig()
			nextConfig["enable_automatic_rotation"] = true
			if tc.next != "" {
				nextConfig["encrypt_new_tables"] = tc.next
			}
			diff, err := clusterResource.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(nextConfig), nil)
			if err != nil {
				t.Fatalf("diff configuration: %s", err)
			}
			if diff == nil || diff.Empty() {
				t.Fatal("expected enable_automatic_rotation to produce a TDE update diff")
			}
			for key, attr := range diff.Attributes {
				// Computed-only entries are unknowns that do not trigger another update path.
				if attr == nil || attr.NewComputed {
					continue
				}
				switch key {
				case "tde_status", "encrypt_new_tables", "encryption_key", "role_arn", "enable_automatic_rotation":
				default:
					t.Fatalf("unexpected diff attribute %q: %#v", key, attr)
				}
			}
			if attr, ok := diff.Attributes["encrypt_new_tables"]; ok && !attr.NewComputed {
				if !tc.wantSent || attr.New != tc.wantValue {
					t.Fatalf("unexpected EncryptNewTables value diff: %#v", attr)
				}
			}

			var captured []url.Values
			client := polarDBClusterTDEHTTPClient(t, func(form url.Values) {
				captured = append(captured, form)
			})
			d, err := schema.InternalMap(clusterResource.Schema).Data(state, diff)
			if err != nil {
				t.Fatalf("build resource data: %s", err)
			}
			d.SetId("pc-polar-tde-unit")

			if err := resourceClusterUpdate(d, client); err == nil {
				t.Fatal("update should return the mocked ModifyDBClusterTDE rejection")
			}
			if len(captured) != 1 {
				t.Fatalf("want exactly one ModifyDBClusterTDE request, got %d", len(captured))
			}
			values, sent := captured[0]["EncryptNewTables"]
			if sent != tc.wantSent {
				t.Fatalf("EncryptNewTables sent = %v, want %v (form: %v)", sent, tc.wantSent, captured[0].Encode())
			}
			if tc.wantSent && (len(values) != 1 || values[0] != tc.wantValue) {
				t.Fatalf("EncryptNewTables = %v, want %q", values, tc.wantValue)
			}
		})
	}
}

// polarDBClusterLocalClient points the polardb endpoint of a test client at a local
// server so update-path unit tests can assert on the requests that were sent.
func polarDBClusterLocalClient(t *testing.T, handler http.Handler) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build dummy credential: %s", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatalf("build client: %s", err)
	}
	endpoints.Store("polardb", endpoint)
	return client
}

// polarDBClusterDeletionLockHTTPClient returns a client whose polardb endpoint is a
// local server. The server answers the live cluster attribute precheck, the deletion
// lock request when expected, and rejects ModifyDBClusterArch so the update returns
// before the provider's StandbyTransing state wait; assertions run on the requests
// that were sent.
func polarDBClusterDeletionLockHTTPClient(t *testing.T, deletionLockExpected bool, capture func(action string, form url.Values)) *connectivity.AliyunClient {
	t.Helper()
	return polarDBClusterLocalClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %s", err)
		}
		action := r.Form.Get("Action")
		capture(action, r.Form)
		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "DescribeDBClusterAttribute":
			_, _ = w.Write([]byte(`{"DBClusterId":"pc-deletion-lock-unit","DBType":"MySQL","DeployUnit":"Advanced","HotStandbyCluster":"StandbyClusterOFF","RequestId":"unit"}`))
		case "ModifyDBClusterDeletion":
			if !deletionLockExpected {
				t.Errorf("deletion lock request must not be sent, form: %v", r.Form)
			}
			_, _ = w.Write([]byte(`{"RequestId":"unit"}`))
		case "ModifyDBClusterArch":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unit test stops the update after the arch request"}`))
		default:
			t.Errorf("combination update sent unexpected action %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"UnexpectedAction","Message":"unexpected action"}`))
		}
	}))
}

// polarDBClusterHotStandbyPrecheckHTTPClient serves DescribeDBClusterAttribute with the
// given body and rejects ModifyDBClusterArch, so precheck unit tests fail loudly if the
// in-place switch is attempted anyway.
func polarDBClusterHotStandbyPrecheckHTTPClient(t *testing.T, attributeBody string, capture func(action string, form url.Values)) *connectivity.AliyunClient {
	t.Helper()
	return polarDBClusterLocalClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %s", err)
		}
		action := r.Form.Get("Action")
		capture(action, r.Form)
		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "DescribeDBClusterAttribute":
			_, _ = w.Write([]byte(attributeBody))
		case "ModifyDBClusterArch":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unit test stops the update after the arch request"}`))
		default:
			t.Errorf("precheck update sent unexpected action %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"UnexpectedAction","Message":"unexpected action"}`))
		}
	}))
}

func TestUnitPolarDBClusterDeletionLockWithHotStandbyUpgrade(t *testing.T) {
	// F-03 regression: an existing cluster that changes deletion_lock together with
	// hot_standby_cluster OFF -> ON must still run the in-place upgrade. The lock
	// operation itself only applies to PostPaid clusters, so a PrePaid cluster skips
	// that single call while the update continues to ModifyDBClusterArch instead of
	// returning early with a silent success.
	baseConfig := func(payType string) map[string]interface{} {
		return map[string]interface{}{
			"db_type":             "MySQL",
			"db_version":          "8.0",
			"db_node_class":       "polar.mysql.x4.large",
			"vswitch_id":          "vsw-deletion-lock-unit",
			"pay_type":            payType,
			"deletion_lock":       0,
			"hot_standby_cluster": "OFF",
		}
	}
	cases := []struct {
		name                string
		payType             string
		wantDeletionLockReq bool
	}{
		{"PrePaid skips the lock call and still upgrades hot standby", "PrePaid", false},
		{"PostPaid applies the lock before upgrading hot standby", "PostPaid", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clusterResource := ResourceCluster()
			previous := schema.TestResourceDataRaw(t, clusterResource.Schema, baseConfig(tc.payType))
			previous.SetId("pc-deletion-lock-unit")
			state := previous.State()

			nextConfig := baseConfig(tc.payType)
			nextConfig["deletion_lock"] = 1
			nextConfig["hot_standby_cluster"] = "ON"
			diff, err := clusterResource.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(nextConfig), nil)
			if err != nil {
				t.Fatalf("diff configuration: %s", err)
			}
			if diff == nil || diff.Empty() {
				t.Fatal("expected deletion_lock and hot_standby_cluster to produce an update diff")
			}
			for key, attr := range diff.Attributes {
				// Computed-only entries are unknowns that do not trigger another update path.
				if attr == nil || attr.NewComputed {
					continue
				}
				switch key {
				case "deletion_lock", "hot_standby_cluster":
				default:
					t.Fatalf("unexpected diff attribute %q: %#v", key, attr)
				}
				if attr.RequiresNew {
					t.Fatalf("attribute %q must update in place", key)
				}
			}

			var archForms, deletionLockForms []url.Values
			client := polarDBClusterDeletionLockHTTPClient(t, tc.wantDeletionLockReq, func(action string, form url.Values) {
				switch action {
				case "ModifyDBClusterDeletion":
					deletionLockForms = append(deletionLockForms, form)
				case "ModifyDBClusterArch":
					archForms = append(archForms, form)
				}
			})
			d, err := schema.InternalMap(clusterResource.Schema).Data(state, diff)
			if err != nil {
				t.Fatalf("build resource data: %s", err)
			}
			d.SetId("pc-deletion-lock-unit")

			if err := resourceClusterUpdate(d, client); err == nil {
				t.Fatal("update should return the mocked ModifyDBClusterArch rejection")
			}
			if len(archForms) != 1 {
				t.Fatalf("want exactly one ModifyDBClusterArch request, got %d", len(archForms))
			}
			arch := archForms[0]
			if got := arch.Get("HotStandbyCluster"); got != "on" {
				t.Fatalf("HotStandbyCluster = %q, want %q", got, "on")
			}
			if got := arch.Get("DBClusterId"); got != "pc-deletion-lock-unit" {
				t.Fatalf("DBClusterId = %q, want %q", got, "pc-deletion-lock-unit")
			}
			if got := len(deletionLockForms) > 0; got != tc.wantDeletionLockReq {
				t.Fatalf("deletion lock request sent = %v, want %v", got, tc.wantDeletionLockReq)
			}
			if tc.wantDeletionLockReq && deletionLockForms[0].Get("Protection") != "true" {
				t.Fatalf("Protection = %q, want %q", deletionLockForms[0].Get("Protection"), "true")
			}
		})
	}
}

func TestUnitPolarDBClusterHotStandbyClusterInPlacePlanBoundaries(t *testing.T) {
	// F-01/F-02 regression: ModifyDBClusterArch only supports switching the storage
	// hot standby from off to on, and only for MySQL clusters ("equal" is retired and
	// clusters with the hot standby already on cannot switch via this API). Every
	// other planned transition keeps the historical replacement behavior, so the plan
	// never promises an in-place call the API does not support.
	clusterResource := ResourceCluster()
	cases := []struct {
		name        string
		dbType      string
		oldMode     string
		newMode     string
		wantInPlace bool
	}{
		{"MySQL off to on upgrades in place", "MySQL", "OFF", "ON", true},
		{"MySQL off to equal recreates", "MySQL", "OFF", "EQUAL", false},
		{"MySQL on to equal recreates", "MySQL", "ON", "EQUAL", false},
		{"MySQL equal to on recreates", "MySQL", "EQUAL", "ON", false},
		{"MySQL on to off recreates", "MySQL", "ON", "OFF", false},
		{"MySQL equal to off recreates", "MySQL", "EQUAL", "OFF", false},
		{"PostgreSQL off to on recreates", "PostgreSQL", "OFF", "ON", false},
		{"Oracle off to on recreates", "Oracle", "OFF", "ON", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := func(mode string) map[string]interface{} {
				return map[string]interface{}{
					"db_type":             tc.dbType,
					"db_version":          "8.0",
					"db_node_class":       "polar.mysql.x4.large",
					"vswitch_id":          "vsw-hot-standby-unit",
					"hot_standby_cluster": mode,
				}
			}
			previous := schema.TestResourceDataRaw(t, clusterResource.Schema, config(tc.oldMode))
			previous.SetId("pc-hot-standby-unit")
			state := previous.State()

			diff, err := clusterResource.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(config(tc.newMode)), nil)
			if err != nil {
				t.Fatalf("diff configuration: %s", err)
			}
			if diff == nil || diff.Empty() {
				t.Fatal("expected hot_standby_cluster to produce a diff")
			}
			attr := diff.Attributes["hot_standby_cluster"]
			if attr == nil || attr.NewComputed {
				t.Fatalf("expected a concrete hot_standby_cluster diff entry, got %#v", attr)
			}
			if got := attr.New; got != tc.newMode {
				t.Fatalf("hot_standby_cluster new value = %q, want %q", got, tc.newMode)
			}
			if wantRequiresNew := !tc.wantInPlace; attr.RequiresNew != wantRequiresNew {
				t.Fatalf("hot_standby_cluster RequiresNew = %v, want %v", attr.RequiresNew, wantRequiresNew)
			}
		})
	}
}

func TestUnitPolarDBClusterHotStandbyClusterUpgradePrecheck(t *testing.T) {
	// F-02 regression: the in-place branch verifies the live cluster attributes
	// before calling ModifyDBClusterArch, so combinations outside the MySQL/Advanced
	// scope fail with the real limitation instead of a server-side rejection.
	cases := []struct {
		name          string
		attributeBody string
		wantErrText   string
	}{
		{"PostgreSQL cluster is rejected before the arch call",
			`{"DBClusterId":"pc-hot-standby-unit","DBType":"PostgreSQL","DeployUnit":"Advanced","HotStandbyCluster":"StandbyClusterOFF","RequestId":"unit"}`, "DBType"},
		{"Non-Advanced cluster is rejected before the arch call",
			`{"DBClusterId":"pc-hot-standby-unit","DBType":"MySQL","DeployUnit":"Basic","HotStandbyCluster":"StandbyClusterOFF","RequestId":"unit"}`, "DeployUnit"},
		{"Cluster with the hot standby already on is rejected",
			`{"DBClusterId":"pc-hot-standby-unit","DBType":"MySQL","DeployUnit":"Advanced","HotStandbyCluster":"StandbyClusterON","RequestId":"unit"}`, "StandbyClusterON"},
		{"Legacy equal cluster is rejected",
			`{"DBClusterId":"pc-hot-standby-unit","DBType":"MySQL","DeployUnit":"Advanced","HotStandbyCluster":"equal","RequestId":"unit"}`, "equal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clusterResource := ResourceCluster()
			config := func(mode string) map[string]interface{} {
				return map[string]interface{}{
					"db_type":             "MySQL",
					"db_version":          "8.0",
					"db_node_class":       "polar.mysql.x4.large",
					"vswitch_id":          "vsw-hot-standby-unit",
					"hot_standby_cluster": mode,
				}
			}
			previous := schema.TestResourceDataRaw(t, clusterResource.Schema, config("OFF"))
			previous.SetId("pc-hot-standby-unit")
			state := previous.State()

			diff, err := clusterResource.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(config("ON")), nil)
			if err != nil {
				t.Fatalf("diff configuration: %s", err)
			}
			var archForms []url.Values
			client := polarDBClusterHotStandbyPrecheckHTTPClient(t, tc.attributeBody, func(action string, form url.Values) {
				if action == "ModifyDBClusterArch" {
					archForms = append(archForms, form)
				}
			})
			d, err := schema.InternalMap(clusterResource.Schema).Data(state, diff)
			if err != nil {
				t.Fatalf("build resource data: %s", err)
			}
			d.SetId("pc-hot-standby-unit")

			if err := resourceClusterUpdate(d, client); err == nil || !strings.Contains(err.Error(), tc.wantErrText) {
				t.Fatalf("update error = %v, want it to mention %q", err, tc.wantErrText)
			}
			if len(archForms) != 0 {
				t.Fatalf("want no ModifyDBClusterArch request, got %d", len(archForms))
			}
		})
	}
}

func TestUnitPolarDBClusterHotStandbyClusterUpdateGuard(t *testing.T) {
	// F-01 regression: even if an unsupported transition (such as OFF -> EQUAL or
	// the reverse ON -> OFF) somehow reaches the update path, the branch refuses
	// it before any request instead of sending a mode ModifyDBClusterArch has
	// retired ("equal") or cannot switch (an already enabled hot standby).
	clusterResource := ResourceCluster()
	config := func(mode string) map[string]interface{} {
		return map[string]interface{}{
			"db_type":             "MySQL",
			"db_version":          "8.0",
			"db_node_class":       "polar.mysql.x4.large",
			"vswitch_id":          "vsw-hot-standby-unit",
			"hot_standby_cluster": mode,
		}
	}
	cases := []struct {
		name    string
		oldMode string
		newMode string
	}{
		{"off to equal is refused", "OFF", "EQUAL"},
		{"on to off is refused", "ON", "OFF"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			previous := schema.TestResourceDataRaw(t, clusterResource.Schema, config(tc.oldMode))
			previous.SetId("pc-hot-standby-unit")
			state := previous.State()

			diff, err := clusterResource.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(config(tc.newMode)), nil)
			if err != nil {
				t.Fatalf("diff configuration: %s", err)
			}
			var requests []string
			client := polarDBClusterHotStandbyPrecheckHTTPClient(t, `{"DBClusterId":"pc-hot-standby-unit","DBType":"MySQL","DeployUnit":"Advanced","HotStandbyCluster":"StandbyClusterOFF","RequestId":"unit"}`, func(action string, form url.Values) {
				requests = append(requests, action)
			})
			d, err := schema.InternalMap(clusterResource.Schema).Data(state, diff)
			if err != nil {
				t.Fatalf("build resource data: %s", err)
			}
			d.SetId("pc-hot-standby-unit")

			if err := resourceClusterUpdate(d, client); err == nil || !strings.Contains(err.Error(), "from OFF to ON") {
				t.Fatalf("update error = %v, want the off-to-on only guard", err)
			}
			if len(requests) != 0 {
				t.Fatalf("want no polardb request, got %v", requests)
			}
		})
	}
}

// polarDBClusterStandbyAzHTTPClient serves the whole update flow of a standby
// zone change: the live cluster attribute precheck and waiter refresh report a
// Running MySQL Advanced cluster with the hot standby off, ModifyDBClusterArch
// succeeds, while ModifyDBClusterPrimaryZone and the trailing Read's
// DescribeDBClusters are rejected, so assertions can run on the requests that
// were sent.
func polarDBClusterStandbyAzHTTPClient(t *testing.T, capture func(action string, form url.Values)) *connectivity.AliyunClient {
	t.Helper()
	return polarDBClusterLocalClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form: %s", err)
		}
		action := r.Form.Get("Action")
		capture(action, r.Form)
		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "DescribeDBClusterAttribute":
			_, _ = w.Write([]byte(`{"DBClusterId":"pc-hot-standby-unit","DBType":"MySQL","DeployUnit":"Advanced","HotStandbyCluster":"StandbyClusterOFF","DBClusterStatus":"Running","RequestId":"unit"}`))
		case "ModifyDBClusterArch":
			_, _ = w.Write([]byte(`{"RequestId":"unit"}`))
		case "ModifyDBClusterPrimaryZone":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unit test: the zone migration must not repeat after the arch upgrade"}`))
		case "DescribeDBClusters":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidParameter","Message":"unit test stops the update after the combined branches"}`))
		default:
			t.Errorf("standby zone update sent unexpected action %q", action)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"UnexpectedAction","Message":"unexpected action"}`))
		}
	}))
}

func TestUnitPolarDBClusterHotStandbyUpgradeConsumesStandbyAzChange(t *testing.T) {
	// F-11 regression: changing the storage hot standby and the standby zone
	// together must send exactly one write. ModifyDBClusterArch establishes the
	// standby zone placement (its StandbyAZ parameter), so it consumes the whole
	// standby_az change; a follow-up ModifyDBClusterPrimaryZone call would only
	// repeat the same zone migration and surface InvalidZoneId.Same.
	clusterResource := ResourceCluster()
	config := func(hotStandby, standbyAz string) map[string]interface{} {
		return map[string]interface{}{
			"db_type":             "MySQL",
			"db_version":          "8.0",
			"db_node_class":       "polar.mysql.x4.large",
			"vswitch_id":          "vsw-hot-standby-unit",
			"hot_standby_cluster": hotStandby,
			"standby_az":          standbyAz,
		}
	}
	previous := schema.TestResourceDataRaw(t, clusterResource.Schema, config("OFF", "cn-hangzhou-a"))
	previous.SetId("pc-hot-standby-unit")
	state := previous.State()

	diff, err := clusterResource.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(config("ON", "cn-hangzhou-b")), nil)
	if err != nil {
		t.Fatalf("diff configuration: %s", err)
	}
	for key, attr := range diff.Attributes {
		// Computed-only entries are unknowns that do not trigger another update path.
		if attr == nil || attr.NewComputed {
			continue
		}
		switch key {
		case "hot_standby_cluster", "standby_az":
		default:
			t.Fatalf("unexpected diff attribute %q: %#v", key, attr)
		}
		if attr.RequiresNew {
			t.Fatalf("attribute %q must update in place", key)
		}
	}

	// Shorten the arch wait lead time so the update reaches the standby zone
	// branch within the test instead of the production four minutes.
	originalWaitDelay := hotStandbyArchWaitDelay
	hotStandbyArchWaitDelay = time.Millisecond
	t.Cleanup(func() { hotStandbyArchWaitDelay = originalWaitDelay })

	var archForms, zoneForms []url.Values
	client := polarDBClusterStandbyAzHTTPClient(t, func(action string, form url.Values) {
		switch action {
		case "ModifyDBClusterArch":
			archForms = append(archForms, form)
		case "ModifyDBClusterPrimaryZone":
			zoneForms = append(zoneForms, form)
		}
	})
	d, err := schema.InternalMap(clusterResource.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("build resource data: %s", err)
	}
	d.SetId("pc-hot-standby-unit")

	// The trailing Read is stopped by the mocked DescribeDBClusters rejection,
	// which proves the update ran both branches and reached its end.
	if err := resourceClusterUpdate(d, client); err == nil || !strings.Contains(err.Error(), "DescribeDBClusters") {
		t.Fatalf("update error = %v, want the Read rejection after the combined branches", err)
	}
	if len(archForms) != 1 {
		t.Fatalf("want exactly one ModifyDBClusterArch request, got %d", len(archForms))
	}
	arch := archForms[0]
	if got := arch.Get("HotStandbyCluster"); got != "on" {
		t.Fatalf("HotStandbyCluster = %q, want %q", got, "on")
	}
	if got := arch.Get("StandbyAZ"); got != "cn-hangzhou-b" {
		t.Fatalf("StandbyAZ = %q, want %q, the arch call must forward the new standby zone", got, "cn-hangzhou-b")
	}
	if got := arch.Get("DBClusterId"); got != "pc-hot-standby-unit" {
		t.Fatalf("DBClusterId = %q, want %q", got, "pc-hot-standby-unit")
	}
	if len(zoneForms) != 0 {
		t.Fatalf("the standby zone change must be consumed by the arch upgrade, got %d ModifyDBClusterPrimaryZone request(s)", len(zoneForms))
	}
}

func TestUnitPolarDBClusterStandbyAzIndependentSwitch(t *testing.T) {
	// The standby zone switch keeps its own path when the hot standby mode does
	// not change: exactly one ModifyDBClusterPrimaryZone request with
	// ZoneType=Standby and no ModifyDBClusterArch request, and the server
	// rejection is returned instead of being swallowed before the state wait.
	clusterResource := ResourceCluster()
	config := func(standbyAz string) map[string]interface{} {
		return map[string]interface{}{
			"db_type":             "MySQL",
			"db_version":          "8.0",
			"db_node_class":       "polar.mysql.x4.large",
			"vswitch_id":          "vsw-hot-standby-unit",
			"hot_standby_cluster": "ON",
			"standby_az":          standbyAz,
		}
	}
	previous := schema.TestResourceDataRaw(t, clusterResource.Schema, config("cn-hangzhou-a"))
	previous.SetId("pc-hot-standby-unit")
	state := previous.State()

	diff, err := clusterResource.Diff(context.Background(), state, sdktf.NewResourceConfigRaw(config("cn-hangzhou-b")), nil)
	if err != nil {
		t.Fatalf("diff configuration: %s", err)
	}
	if attr := diff.Attributes["standby_az"]; attr == nil || attr.NewComputed {
		t.Fatalf("expected a concrete standby_az diff entry, got %#v", attr)
	}

	var archForms, zoneForms []url.Values
	client := polarDBClusterStandbyAzHTTPClient(t, func(action string, form url.Values) {
		switch action {
		case "ModifyDBClusterArch":
			archForms = append(archForms, form)
		case "ModifyDBClusterPrimaryZone":
			zoneForms = append(zoneForms, form)
		}
	})
	d, err := schema.InternalMap(clusterResource.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("build resource data: %s", err)
	}
	d.SetId("pc-hot-standby-unit")

	if err := resourceClusterUpdate(d, client); err == nil || !strings.Contains(err.Error(), "ModifyDBClusterPrimaryZone") {
		t.Fatalf("update error = %v, want the zone request rejection to be returned", err)
	}
	if len(zoneForms) != 1 {
		t.Fatalf("want exactly one ModifyDBClusterPrimaryZone request, got %d", len(zoneForms))
	}
	zone := zoneForms[0]
	if got := zone.Get("ZoneId"); got != "cn-hangzhou-b" {
		t.Fatalf("ZoneId = %q, want %q", got, "cn-hangzhou-b")
	}
	if got := zone.Get("ZoneType"); got != "Standby" {
		t.Fatalf("ZoneType = %q, want %q", got, "Standby")
	}
	if got := zone.Get("DBClusterId"); got != "pc-hot-standby-unit" {
		t.Fatalf("DBClusterId = %q, want %q", got, "pc-hot-standby-unit")
	}
	if len(archForms) != 0 {
		t.Fatalf("want no ModifyDBClusterArch request, got %d", len(archForms))
	}
}
