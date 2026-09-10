package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudADBBackupPolicy(t *testing.T) {
	var v *adb.DescribeBackupPolicyResponse
	resourceId := "alicloud_adb_backup_policy.default"
	serverFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeAdbBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccAdbBackupPolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbBackupPolicyConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create the backup policy with log backup enabled and a retention period.
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":               "${alicloud_adb_db_cluster.default.id}",
					"preferred_backup_period":     []string{"Tuesday", "Wednesday"},
					"preferred_backup_time":       "10:00Z-11:00Z",
					"enable_backup_log":           "Enable",
					"log_backup_retention_period": 7,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#":   "2",
						"preferred_backup_time":       "10:00Z-11:00Z",
						"enable_backup_log":           "Enable",
						"log_backup_retention_period": "7",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update the log backup retention period.
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": 30,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "30",
						"enable_backup_log":           "Enable",
					}),
				),
			},
			{
				// Disable log backup. The retention period is only meaningful when
				// log backup is enabled, so it is removed from the config to let
				// the server value settle without a perpetual diff.
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log":           "Disable",
					"log_backup_retention_period": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log": "Disable",
					}),
				),
			},
			{
				// Re-enable log backup with a retention period to verify the
				// set/read roundtrip works after a disable.
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log":           "Enable",
					"log_backup_retention_period": 7,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log":           "Enable",
						"log_backup_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Monday", "Saturday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "15:00Z-16:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "15:00Z-16:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Tuesday", "Thursday", "Friday", "Sunday"},
					"preferred_backup_time":   "17:00Z-18:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "4",
						"preferred_backup_time":     "17:00Z-18:00Z",
					}),
				),
			}},
	})
}

func resourceAdbBackupPolicyConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	resource "alicloud_adb_db_cluster" "default" {
		db_cluster_category = "MixedStorage"
		mode = "flexible"
		compute_resource = "8Core32GB"
		vswitch_id              = local.vswitch_id
		description             = "${var.name}"
	}`, AdbCommonTestCase, name)
}
