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
	client, err := config.Client()
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
