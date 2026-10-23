package alicloud

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestUnitPolarDBAdvancedBackupPolicyActions(t *testing.T) {
	old := []map[string]interface{}{
		{
			"policy_id":       "policy-update",
			"src_type":        "db",
			"dest_type":       "level1",
			"filter_type":     "crontab",
			"filter_key":      "dayOfWeek",
			"retention_type":  "delay",
			"retention_value": "7",
		},
		{
			"policy_id":   "policy-delete",
			"src_type":    "level1",
			"dest_type":   "level2",
			"filter_type": "event",
		},
	}
	desired := []map[string]interface{}{
		{
			"src_type":        "db",
			"dest_type":       "level1",
			"filter_type":     "crontab",
			"filter_key":      "dayOfWeek",
			"retention_type":  "delay",
			"retention_value": "14",
		},
		{
			"src_type":                       "level2",
			"dest_type":                      "level2Cross",
			"filter_type":                    "event",
			"retention_type":                 "never",
			"only_preserve_one_each_day":     false,
			"only_preserve_one_each_day_set": true,
		},
	}
	actions := polarDBAdvancedPolicyActions(old, desired)
	if len(actions) != 3 || actions[0].Action != "UPDATE" || actions[1].Action != "CREATE" || actions[2].Action != "DELETE" {
		t.Fatalf("unexpected actions: %#v", actions)
	}
	request := map[string]interface{}{
		"DBClusterId":       "pc-test",
		"BackupPolicyLevel": "Advanced",
	}
	if err := expandPolarDBAdvancedPolicyActions(request, actions); err != nil {
		t.Fatal(err)
	}
	var encoded []map[string]interface{}
	if err := json.Unmarshal([]byte(request["AdvancedDataPolicies"].(string)), &encoded); err != nil {
		t.Fatal(err)
	}
	if encoded[0]["ActionType"] != "UPDATE" || encoded[0]["PolicyId"] != "policy-update" || encoded[0]["RetentionValue"] != "14" {
		t.Fatalf("update action was not encoded correctly: %#v", encoded[0])
	}
	if encoded[1]["ActionType"] != "CREATE" || encoded[1]["OnlyPreserveOneEachDay"] != false {
		t.Fatalf("create action was not encoded correctly: %#v", encoded[1])
	}
	if _, exists := encoded[1]["PolicyId"]; exists {
		t.Fatal("CREATE must not send PolicyId")
	}
	if encoded[2]["ActionType"] != "DELETE" || encoded[2]["PolicyId"] != "policy-delete" || len(encoded[2]) != 2 {
		t.Fatalf("delete action was not encoded correctly: %#v", encoded[2])
	}
	if _, exists := request["AdvancedDataPolicies.1.ActionType"]; exists {
		t.Fatal("advanced policies must be sent as one JSON array parameter")
	}
}

func TestUnitPolarDBAdvancedBackupPolicyNoopAndAdoption(t *testing.T) {
	previous := []map[string]interface{}{
		{
			"policy_id":                  "policy-existing",
			"src_type":                   "db",
			"dest_type":                  "level1",
			"filter_type":                "event",
			"retention_type":             "never",
			"storage_class":              "",
			"only_preserve_one_each_day": true,
		},
	}
	desired := []map[string]interface{}{
		{
			"src_type":       "db",
			"dest_type":      "level1",
			"filter_type":    "event",
			"retention_type": "never",
		},
	}
	if actions := polarDBAdvancedPolicyActions(previous, desired); len(actions) != 0 {
		t.Fatalf("an unchanged declaration must not produce another API action: %#v", actions)
	}
	adopted := []map[string]interface{}{
		{
			"policy_id":      "policy-existing",
			"src_type":       "db",
			"dest_type":      "level1",
			"filter_type":    "event",
			"retention_type": "never",
		},
	}
	if actions := polarDBAdvancedPolicyActions(nil, adopted); len(actions) != 1 || actions[0].Action != "UPDATE" {
		t.Fatalf("an explicit existing PolicyId must be adopted with UPDATE, not duplicated: %#v", actions)
	}
}

func TestUnitPolarDBAdvancedBackupPolicyRegions(t *testing.T) {
	previous := []map[string]interface{}{
		{
			"policy_id":       "policy-test",
			"src_type":        "db",
			"dest_type":       "level1",
			"filter_type":     "crontab",
			"filter_key":      "dayOfWeek",
			"retention_type":  "delay",
			"retention_value": "7",
		},
	}
	desired := []map[string]interface{}{
		{
			"src_type":        "db",
			"dest_type":       "level1",
			"filter_type":     "crontab",
			"filter_key":      "dayOfWeek",
			"retention_type":  "delay",
			"retention_value": "14",
		},
	}
	if err := preparePolarDBAdvancedPolicyRegions(previous, "cn-beijing", false); err != nil {
		t.Fatal(err)
	}
	if err := preparePolarDBAdvancedPolicyRegions(desired, "cn-beijing", true); err != nil {
		t.Fatal(err)
	}
	if previous[0]["src_region"] != "cn-beijing" || desired[0]["dest_region"] != "cn-beijing" {
		t.Fatalf("same-region rule was not normalized: previous=%#v desired=%#v", previous, desired)
	}
	createActions := polarDBAdvancedPolicyActions(nil, desired)
	createRequest := map[string]interface{}{
		"DBClusterId":       "pc-test",
		"BackupPolicyLevel": "Advanced",
	}
	if err := expandPolarDBAdvancedPolicyActions(createRequest, createActions); err != nil {
		t.Fatal(err)
	}
	var created []map[string]interface{}
	if err := json.Unmarshal([]byte(createRequest["AdvancedDataPolicies"].(string)), &created); err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 || created[0]["ActionType"] != "CREATE" || created[0]["SrcRegion"] != "cn-beijing" || created[0]["DestRegion"] != "cn-beijing" {
		t.Fatalf("CREATE request must contain both API regions: %#v", created)
	}
	actions := polarDBAdvancedPolicyActions(previous, desired)
	if len(actions) != 1 || actions[0].Action != "UPDATE" {
		t.Fatalf("normalized rule must match the existing ID: %#v", actions)
	}
	request := map[string]interface{}{
		"DBClusterId":       "pc-test",
		"BackupPolicyLevel": "Advanced",
	}
	if err := expandPolarDBAdvancedPolicyActions(request, actions); err != nil {
		t.Fatal(err)
	}
	var encoded []map[string]interface{}
	if err := json.Unmarshal([]byte(request["AdvancedDataPolicies"].(string)), &encoded); err != nil {
		t.Fatal(err)
	}
	if encoded[0]["SrcRegion"] != "cn-beijing" || encoded[0]["DestRegion"] != "cn-beijing" {
		t.Fatalf("both required API regions must be sent: %#v", encoded[0])
	}
	observed := map[string]interface{}{
		"src_type":        "db",
		"src_region":      "cn-beijing",
		"dest_type":       "level1",
		"dest_region":     "cn-beijing",
		"filter_type":     "crontab",
		"filter_key":      "dayOfWeek",
		"retention_type":  "delay",
		"retention_value": "14",
	}
	if polarDBAdvancedPolicyIdentity(observed) != polarDBAdvancedPolicyIdentity(desired[0]) || !polarDBAdvancedPolicyApplied(observed, desired[0]) {
		t.Fatal("normalized action must match Describe read-back and waiter checks")
	}
	cross := []map[string]interface{}{
		{
			"src_type":   "level1",
			"dest_type":  "level2Cross",
			"src_region": "cn-beijing",
		},
	}
	if err := preparePolarDBAdvancedPolicyRegions(cross, "cn-beijing", true); err == nil {
		t.Fatal("cross-region rule must not default its destination region")
	}
	cross[0]["dest_region"] = "cn-hangzhou"
	if err := preparePolarDBAdvancedPolicyRegions(cross, "cn-beijing", true); err != nil {
		t.Fatal(err)
	}
}

func TestUnitPolarDBAdvancedBackupPolicyDescribeNormalization(t *testing.T) {
	response := map[string]interface{}{
		"AdvancedDataPolicies": map[string]interface{}{
			"AdvancedDataPolicy": []interface{}{
				map[string]interface{}{
					"PolicyId":       "policy-managed\n",
					"SrcType":        "db",
					"DestType":       "level1",
					"FilterType":     "crontab",
					"FilterKey":      "dayOfWeek",
					"RetentionType":  "delay",
					"RetentionValue": "7",
					"AutoCreated":    false,
				},
				map[string]interface{}{
					"PolicyId":    "policy-system",
					"AutoCreated": true,
				},
			},
		},
	}
	observed := polarDBDescribeAdvancedPolicies(response)
	previous := []map[string]interface{}{
		{
			"src_type":    "db",
			"dest_type":   "level1",
			"filter_type": "crontab",
			"filter_key":  "dayOfWeek",
		},
	}
	managed, auto := reconcilePolarDBAdvancedPolicies(previous, observed)
	if len(managed) != 1 || managed[0]["policy_id"] != "policy-managed" {
		t.Fatalf("managed policy ID was not read back: %#v", managed)
	}
	if len(auto) != 1 || auto[0]["policy_id"] != "policy-system" {
		t.Fatalf("auto-created policy was not separated: %#v", auto)
	}
	if polarDBAdvancedPolicyIdentity(previous[0]) != polarDBAdvancedPolicyIdentity(managed[0]) {
		t.Fatal("omitted optional regions and empty regions must have the same identity")
	}
	if !polarDBAdvancedPolicyApplied(managed[0], map[string]interface{}{
		"src_type":        "db",
		"src_region":      "",
		"dest_type":       "level1",
		"dest_region":     "",
		"filter_type":     "crontab",
		"filter_key":      "dayOfWeek",
		"retention_type":  "delay",
		"retention_value": "7",
	}) {
		t.Fatal("Describe must satisfy a CREATE request despite omitted optional regions")
	}
	if polarDBAdvancedPolicyApplied(managed[0], map[string]interface{}{"retention_value": "14"}) {
		t.Fatal("UPDATE must wait for its changed fields to appear in Describe")
	}
}

func TestUnitPolarDBAdvancedBackupPolicyOptionalBool(t *testing.T) {
	resource := resourceAlicloudPolarDBBackupPolicy()
	base := map[string]interface{}{
		"db_cluster_id":       "pc-test",
		"backup_policy_level": "Advanced",
		"advanced_data_policies": []interface{}{
			map[string]interface{}{
				"src_type":       "db",
				"dest_type":      "level1",
				"filter_type":    "event",
				"retention_type": "never",
			},
		},
	}
	d := schema.TestResourceDataRaw(t, resource.Schema, base)
	if _, exists := d.GetOkExists("advanced_data_policies.0.only_preserve_one_each_day"); exists {
		t.Fatal("unset optional boolean must not be sent")
	}
	base["advanced_data_policies"] = []interface{}{
		map[string]interface{}{
			"src_type":                   "db",
			"dest_type":                  "level1",
			"filter_type":                "event",
			"retention_type":             "never",
			"only_preserve_one_each_day": false,
		},
	}
	d = schema.TestResourceDataRaw(t, resource.Schema, base)
	if _, exists := d.GetOkExists("advanced_data_policies.0.only_preserve_one_each_day"); !exists {
		t.Fatal("explicit false must be distinguishable from an unset optional boolean")
	}
}

func TestUnitPolarDBAdvancedBackupPolicyCustomizeDiff(t *testing.T) {
	resource := resourceAlicloudPolarDBBackupPolicy()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"db_cluster_id":               "pc-test",
		"backup_policy_level":         "Advanced",
		"backup_retention_period":     "7",
		"log_backup_retention_period": 7,
	})
	data.SetId("pc-test")
	state := data.State()
	raw := make(map[string]cty.Value, len(resource.Schema))
	for key := range resource.Schema {
		raw[key] = cty.NullVal(cty.DynamicPseudoType)
	}
	raw["db_cluster_id"] = cty.StringVal("pc-test")
	raw["backup_policy_level"] = cty.StringVal("Advanced")
	state.RawConfig = cty.ObjectVal(raw)

	_, err := resource.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]interface{}{
		"db_cluster_id":       "pc-test",
		"backup_policy_level": "Advanced",
	}), nil)
	if err != nil {
		t.Fatalf("read-back of omitted legacy Optional+Computed fields must not be rejected: %v", err)
	}

	raw["log_backup_retention_period"] = cty.NumberIntVal(14)
	state.RawConfig = cty.ObjectVal(raw)
	_, err = resource.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]interface{}{
		"db_cluster_id":               "pc-test",
		"backup_policy_level":         "Advanced",
		"log_backup_retention_period": 14,
	}), nil)
	if err == nil {
		t.Fatal("explicit log-backup changes must be rejected in Advanced mode")
	}
}

func TestUnitPolarDBAdvancedBackupPolicyOptionalRemovalDiff(t *testing.T) {
	resource := resourceAlicloudPolarDBBackupPolicy()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"db_cluster_id":       "pc-test",
		"backup_policy_level": "Advanced",
		"advanced_data_policies": []interface{}{
			map[string]interface{}{
				"policy_id":                  "policy-test",
				"retention_type":             "delay",
				"retention_value":            "7",
				"filter_type":                "crontab",
				"filter_key":                 "dayOfWeek",
				"filter_value":               "1,2,3,4,5,6,7",
				"src_type":                   "db",
				"src_region":                 "cn-beijing",
				"dest_type":                  "level1",
				"dest_region":                "cn-beijing",
				"dump_action":                "copy",
				"storage_class":              "STANDARD",
				"only_preserve_one_each_day": true,
			},
		},
	})
	data.SetId("pc-test")
	diff, err := resource.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(map[string]interface{}{
		"db_cluster_id":       "pc-test",
		"backup_policy_level": "Advanced",
		"advanced_data_policies": []interface{}{
			map[string]interface{}{
				"retention_type": "delay",
				"filter_type":    "crontab",
				"src_type":       "db",
				"dest_type":      "level1",
			},
		},
	}), nil)
	if err != nil {
		t.Fatalf("removing Optional+Computed fields must not be rejected: %v", err)
	}
	for key, attribute := range diff.Attributes {
		if strings.HasPrefix(key, "advanced_data_policies") && attribute.Old != attribute.New {
			t.Fatalf("omitted Optional+Computed field produced a plan diff: %s: %#v", key, attribute)
		}
	}

	_, err = resource.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(map[string]interface{}{
		"db_cluster_id":       "pc-test",
		"backup_policy_level": "Advanced",
		"advanced_data_policies": []interface{}{
			map[string]interface{}{
				"retention_type": "delay",
				"filter_type":    "crontab",
				"src_type":       "db",
				"dest_type":      "level1",
				"filter_value":   "",
			},
		},
	}), nil)
	if err == nil || !strings.Contains(err.Error(), "filter_value cannot be cleared") {
		t.Fatalf("explicit empty filter_value must be rejected, got %v", err)
	}
}

func TestUnitPolarDBAdvancedBackupPolicyCrossRegionDiff(t *testing.T) {
	resource := resourceAlicloudPolarDBBackupPolicy()
	_, err := resource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(map[string]interface{}{
		"db_cluster_id":       "pc-test",
		"backup_policy_level": "Advanced",
		"advanced_data_policies": []interface{}{
			map[string]interface{}{
				"retention_type":  "delay",
				"retention_value": "30",
				"filter_type":     "event",
				"src_type":        "level1",
				"src_region":      "cn-beijing",
				"dest_type":       "level2Cross",
			},
		},
	}), nil)
	if err == nil || !strings.Contains(err.Error(), "requires explicit src_region and dest_region") {
		t.Fatalf("cross-region rule without destination region must be rejected, got %v", err)
	}
}

func TestUnitPolarDBAdvancedBackupPolicyACCRemovalConfig(t *testing.T) {
	config := resourceTestAccConfigFunc(
		"alicloud_polardb_backup_policy.default",
		"tf-testAccPolarDBAdvancedBackupPolicyConfig",
		resourcePolarDBBackupPolicyConfigDependence,
	)
	created := config(map[string]interface{}{
		"db_cluster_id":       "pc-test",
		"backup_policy_level": "Advanced",
		"advanced_data_policies": []map[string]interface{}{
			{
				"retention_type":  "delay",
				"retention_value": "7",
				"filter_type":     "crontab",
				"filter_key":      "dayOfWeek",
				"src_type":        "db",
				"dest_type":       "level1",
			},
		},
	})
	if !strings.Contains(created, "advanced_data_policies {") || strings.Contains(created, "advanced_data_policies =") {
		t.Fatalf("create configuration must use a nested block: %s", created)
	}
	removed := config(map[string]interface{}{
		"advanced_data_policies": REMOVEKEY,
	})
	if strings.Contains(removed, "advanced_data_policies") {
		t.Fatalf("removal configuration must omit the nested block entirely: %s", removed)
	}
	parser := hclparse.NewParser()
	if _, diagnostics := parser.ParseHCL([]byte(removed), "backup_policy.tf"); diagnostics.HasErrors() {
		t.Fatalf("removal configuration is invalid HCL: %s", diagnostics.Error())
	}
}
