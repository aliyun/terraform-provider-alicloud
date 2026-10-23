package polardb_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudPolarDBBackup_basic exercises the property_completeness [M-RT-0040]
// mappings that apply in the default (Normal) backup-policy level: CreateBackup.Comment
// (write/read-back) and ModifyBackupPolicy.BackupPolicyLevel + PreferredBackupPeriod/Time
// (applied through the Create secondary ModifyBackupPolicy call and through Update).
// RestoreTable [M-RT-0039] is a model-level Backup update operation (verified in
// cloudspec_pre) and is deliberately not part of the declarative update contract.
func TestAccAliCloudPolarDBBackup_basic(t *testing.T) {
	resourceId := "alicloud_polardb_backup.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testacc%spolardbbackup%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// IDRefreshName is intentionally omitted. The backup resource reads backup-set
		// and backup-task progress fields (backup_progress_status, process, create_time,
		// backup_set_size, backup_end_time, status, ...) from the DescribeBackups and
		// DescribeBackupTasks list responses (last-entry mapping); those reflect the
		// latest backup/task and change between reads as the backup progresses, so an
		// ID-only refresh cannot be idempotent. Stable policy fields are still asserted
		// via the Steps 1-2 Check and verified on import (volatile fields ignored below).
		CheckDestroy: testAccCheckPolarDBBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPolarDBBackupBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					// F1/F2: the stable backup set id is materialized after the async
					// CreateBackup task completes; assert it is set (non-empty).
					resource.TestCheckResourceAttrSet(resourceId, "backup_id"),
					resource.TestCheckResourceAttr(resourceId, "comment", name),
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Normal"),
					resource.TestCheckResourceAttr(resourceId, "preferred_backup_period", "Tuesday,Wednesday"),
					resource.TestCheckResourceAttr(resourceId, "preferred_backup_time", "10:00Z-11:00Z"),
				),
			},
			{
				Config: testAccAliCloudPolarDBBackupBasicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "comment", name),
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Normal"),
					// F7: preferred_backup_period uses weekday-set semantic equality, so the API's
					// canonical reordering (Mon..Sun) does not produce a persistent plan diff.
					resource.TestCheckResourceAttr(resourceId, "preferred_backup_period", "Monday,Wednesday,Saturday"),
					resource.TestCheckResourceAttr(resourceId, "preferred_backup_time", "15:00Z-16:00Z"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"comment",
					// Volatile backup-set fields sourced from DescribeBackups (last-entry =
					// latest backup set) and backup-task fields from DescribeBackupTasks
					// (latest task progress). These change between reads as backups are
					// created/progress, so they cannot match between the prior state and a
					// fresh import. Stable policy fields (backup_policy_level,
					// preferred_backup_period/time, backup_frequency,
					// backup_retention_policy_on_cluster_deletion, advanced_data_policies)
					// and identity (id, db_cluster_id) remain verified by ImportStateVerify.
					"backup_end_time", "backup_method", "backup_mode",
					"backup_progress_status", "backup_set_size", "backup_start_time", "backup_type",
					"backups_level", "consistent_time", "create_time", "is_avail", "job_mode",
					"preferred_next_backup_time", "process", "status", "store_status", "task_action"},
			},
		},
	})
}

// TestAccAliCloudPolarDBBackup_advancedDataPolicies exercises the AdvancedDataPolicies
// mapping (ModifyBackupPolicy.AdvancedDataPolicies, 17 sub-fields) under the Advanced
// backup-policy level. AdvancedDataPolicies is only meaningful with BackupPolicyLevel =
// Advanced, which requires the cluster to have advanced backup enabled (申请使用, per
// the BackupPolicyLevel API contract). If the test account/cluster lacks that capability
// the case is reported as a capability gate rather than an implementation defect.
func TestAccAliCloudPolarDBBackup_advancedDataPolicies(t *testing.T) {
	resourceId := "alicloud_polardb_backup.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testacc%spolardbbackupadv%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// IDRefreshName is intentionally omitted. The backup resource reads backup-set
		// and backup-task progress fields (backup_progress_status, process, create_time,
		// backup_set_size, backup_end_time, status, ...) from the DescribeBackups and
		// DescribeBackupTasks list responses (last-entry mapping); those reflect the
		// latest backup/task and change between reads as the backup progresses, so an
		// ID-only refresh cannot be idempotent. Stable policy fields are still asserted
		// via the Steps 1-2 Check and verified on import (volatile fields ignored below).
		CheckDestroy: testAccCheckPolarDBBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPolarDBBackupAdvanced(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Advanced"),
					// QA 45954: Advanced policies require the full src_type/dest_type/
					// src_region/dest_region quartet; assert the configured db->level1
					// chain instead of the previous single src_type entry. The check
					// matches by chain so it does not depend on API ordering or on
					// system-generated entries.
					testAccCheckPolarDBBackupAdvancedPolicy(resourceId, region, "7"),
				),
			},
			{
				Config: testAccAliCloudPolarDBBackupAdvancedUpdate(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Advanced"),
					testAccCheckPolarDBBackupAdvancedPolicy(resourceId, region, "14"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// preferred_backup_period has no read-back under BackupPolicyLevel=Advanced
				// (doc 2319233: 开启高级备份后，该参数不再生效; the API only reports all seven
				// days), so import cannot reproduce the configured Tuesday,Wednesday value.
				ImportStateVerifyIgnore: []string{"comment", "advanced_data_policies", "preferred_backup_period",
					"backup_end_time", "backup_method", "backup_mode",
					"backup_progress_status", "backup_set_size", "backup_start_time", "backup_type",
					"backups_level", "consistent_time", "create_time", "is_avail", "job_mode",
					"preferred_next_backup_time", "process", "status", "store_status", "task_action"},
			},
		},
	})
}

func testAccCheckPolarDBBackupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_backup" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		// F1/F8: the resource id is the composite "DBClusterId:BackupId". Resolve both parts
		// and query DescribeBackups for the exact owned BackupId. The previous check only
		// verified the cluster policy was gone, so a DeleteBackup that silently no-op'd or
		// deleted the wrong (foreign) backup set would still pass once the fixture cluster
		// was destroyed. Querying by BackupId fails the test when the owned backup survives.
		dbClusterId, backupId := rs.Primary.ID, ""
		if idx := strings.Index(rs.Primary.ID, ":"); idx >= 0 {
			dbClusterId, backupId = rs.Primary.ID[:idx], rs.Primary.ID[idx+1:]
		}
		if dbClusterId == "" || backupId == "" {
			return fmt.Errorf("polardb backup destroy check: malformed composite id %q", rs.Primary.ID)
		}
		now := time.Now().UTC()
		request := map[string]interface{}{
			"DBClusterId": dbClusterId,
			"BackupId":    backupId,
			"StartTime":   now.Add(-7 * 24 * time.Hour).Format("2006-01-02T15:04Z"),
			"EndTime":     now.Format("2006-01-02T15:04Z"),
			"RegionId":    client.RegionId,
		}
		resp, err := client.RpcPost("polardb", "2017-08-01", "DescribeBackups", nil, request, true)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing backups for cluster %s backup %s failed: %s", dbClusterId, backupId, err)
		}
		// RpcPost returns map[string]interface{} directly; no interface assertion needed.
		if resp != nil {
			if items, ok := resp["Items"].(map[string]interface{}); ok && items != nil {
				if backups, ok := items["Backup"].([]interface{}); ok && len(backups) > 0 {
					return fmt.Errorf("polardb backup set %s for cluster %s still exists after destroy", backupId, dbClusterId)
				}
			}
		}
	}
	return nil
}

// testAccCheckPolarDBBackupAdvancedPolicy verifies that the configured db->level1
// advanced data policy is present with the expected retention. It matches by chain
// instead of list index so the assertion is independent of API ordering and of
// system-generated policies returned by DescribeBackupPolicy.
func testAccCheckPolarDBBackupAdvancedPolicy(resourceId, region, retentionValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("polardb backup %s not found in state", resourceId)
		}
		attrs := rs.Primary.Attributes
		count, err := strconv.Atoi(attrs["advanced_data_policies.#"])
		if err != nil {
			return fmt.Errorf("polardb backup %s has no readable advanced_data_policies count: %v", resourceId, err)
		}
		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("advanced_data_policies.%d.", i)
			if attrs[prefix+"src_type"] != "db" || attrs[prefix+"dest_type"] != "level1" ||
				attrs[prefix+"src_region"] != region || attrs[prefix+"dest_region"] != region {
				continue
			}
			if got := attrs[prefix+"bak_type"]; got != "F" {
				return fmt.Errorf("advanced policy db->level1 has bak_type %q, expected F", got)
			}
			if got := attrs[prefix+"retention_type"]; got != "delay" {
				return fmt.Errorf("advanced policy db->level1 has retention_type %q, expected delay", got)
			}
			if got := attrs[prefix+"retention_value"]; got != retentionValue {
				return fmt.Errorf("advanced policy db->level1 has retention_value %q, expected %s", got, retentionValue)
			}
			return nil
		}
		return fmt.Errorf("polardb backup %s has no db->level1 advanced policy for region %s among %d entries", resourceId, region, count)
	}
}

func testAccAliCloudPolarDBBackupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  pay_type   = "PostPaid"
  db_type    = "MySQL"
  db_version = "8.0"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = "polar.mysql.x4.large"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
  // encrypt_new_tables is Optional (not Computed) on the cluster resource and the
  // API returns "OFF" for a fresh cluster; set it explicitly so the post-apply plan
  // is empty (otherwise the fixture drifts: encrypt_new_tables = "OFF" -> null).
  encrypt_new_tables = "OFF"
}
`, name)
}

func testAccAliCloudPolarDBBackupBasic(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  comment                 = var.name
  backup_policy_level     = "Normal"
  preferred_backup_period = "Tuesday,Wednesday"
  preferred_backup_time   = "10:00Z-11:00Z"
}
`, testAccAliCloudPolarDBBackupDependence(name))
}

func testAccAliCloudPolarDBBackupBasicUpdate(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  comment                 = var.name
  backup_policy_level     = "Normal"
  preferred_backup_period = "Monday,Wednesday,Saturday"
  preferred_backup_time   = "15:00Z-16:00Z"
}
`, testAccAliCloudPolarDBBackupDependence(name))
}

func testAccAliCloudPolarDBBackupAdvanced(name, region string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  backup_policy_level     = "Advanced"
  preferred_backup_period = "Tuesday,Wednesday"
  preferred_backup_time   = "10:00Z-11:00Z"

  // Advanced policies are validated server-side as a full chain: the request must
  // carry srcType, destType, srcRegion and destRegion together (QA 45954 request id
  // 01A16365-98BA-52EE-9520-26931432B5C4). db -> level1 in the test region is the
  // base chain documented for ModifyBackupPolicy.AdvancedDataPolicies, and the
  // scheduling fields (FilterType/FilterKey/FilterValue, OnlyPreserveOneEachDay)
  // are part of the policy definition: omitting them made DBS reject the save with
  // DBS.Policy.BackupPolicy.SaveFailed (request id 01A166A1-2107-5CFA-B029-82739368B890).
  advanced_data_policies {
    action_type                = "CREATE"
    bak_type                   = "F"
    src_type                   = "db"
    dest_type                  = "level1"
    src_region                 = "%s"
    dest_region                = "%s"
    filter_type                = "crontab"
    filter_key                 = "dayOfWeek"
    filter_value               = "1,2,3,4,5,6,7"
    only_preserve_one_each_day = true
    retention_type             = "delay"
    retention_value            = "7"
  }
}
`, testAccAliCloudPolarDBBackupDependence(name), region, region)
}

func testAccAliCloudPolarDBBackupAdvancedUpdate(name, region string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  backup_policy_level     = "Advanced"
  preferred_backup_period = "Tuesday,Wednesday"
  preferred_backup_time   = "10:00Z-11:00Z"

  advanced_data_policies {
    action_type                = "UPDATE"
    bak_type                   = "F"
    src_type                   = "db"
    dest_type                  = "level1"
    src_region                 = "%s"
    dest_region                = "%s"
    filter_type                = "crontab"
    filter_key                 = "dayOfWeek"
    filter_value               = "1,2,3,4,5,6,7"
    only_preserve_one_each_day = true
    retention_type             = "delay"
    retention_value            = "14"
  }
}
`, testAccAliCloudPolarDBBackupDependence(name), region, region)
}
