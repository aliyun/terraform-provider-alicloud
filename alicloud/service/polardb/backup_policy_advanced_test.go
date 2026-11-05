package polardb_test

import (
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func TestUnitPolarDBAdvancedBackupPolicyACCRemovalConfig(t *testing.T) {
	config := acctest.ResourceTestAccConfig(
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
		"advanced_data_policies": acctest.REMOVEKEY,
	})
	if strings.Contains(removed, "advanced_data_policies") {
		t.Fatalf("removal configuration must omit the nested block entirely: %s", removed)
	}
	parser := hclparse.NewParser()
	if _, diagnostics := parser.ParseHCL([]byte(removed), "backup_policy.tf"); diagnostics.HasErrors() {
		t.Fatalf("removal configuration is invalid HCL: %s", diagnostics.Error())
	}
}
