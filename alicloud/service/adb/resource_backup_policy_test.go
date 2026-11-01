// Package adb_test.
package adb_test

import (
	"fmt"
	"testing"

	adbsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/adb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudADBBackupPolicy(t *testing.T) {
	var v *adbsdk.DescribeBackupPolicyResponse
	resourceId := "alicloud_adb_backup_policy.default"
	serverFunc := func() interface{} {
		return adb.NewAdbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeAdbBackupPolicy")
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := "tf-testAccAdbBackupPolicy"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAdbBackupPolicyConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":           "${alicloud_adb_db_cluster.default.id}",
					"preferred_backup_period": []string{"Tuesday", "Wednesday"},
					"preferred_backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "10:00Z-11:00Z",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
