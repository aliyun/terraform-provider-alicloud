package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudHBRNasBackupPlan_basic0(t *testing.T) {
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_nas_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudHBRNasBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrNasBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%shbrnasbackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRNasBackupPlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_type":          "COMPLETE",
					"vault_id":             "${alicloud_hbr_vault.default.id}",
					"file_system_id":       "${alicloud_nas_file_system.default.id}",
					"schedule":             "I|1602673264|PT2H",
					"nas_backup_plan_name": "tf-testAccCase",
					"retention":            "1",
					"path":                 []string{"/"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_type":          "COMPLETE",
						"schedule":             "I|1602673264|PT2H",
						"nas_backup_plan_name": "tf-testAccCase",
						"retention":            "1",
						"path.#":               "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nas_backup_plan_name": "tf-testAccCase2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nas_backup_plan_name": "tf-testAccCase2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": "I|1602673264|P1D",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule": "I|1602673264|P1D",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"path": []string{"/home/test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": "{\\\"UseVSS\\\":false}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options": "{\"UseVSS\":false}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nas_backup_plan_name": "tf-testAccCase3",
					"schedule":             "I|1602673264|PT2H",
					"retention":            "4",
					"path":                 []string{"/home/test2", "/home/test2"},
					"options":              "{\\\"UseVSS\\\":true}",
					"disabled":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nas_backup_plan_name": "tf-testAccCase3",
						"schedule":             "I|1602673264|PT2H",
						"retention":            "4",
						"path.#":               "2",
						"options":              "{\"UseVSS\":true}",
						"disabled":             "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_paths", "auto_remove_file_system_mount_point"},
			},
		},
	})
}

var AlicloudHBRNasBackupPlanMap0 = map[string]string{
	"path.#":                              NOSET,
	"retention":                           "",
	"disk_id":                             NOSET,
	"options":                             "",
	"exclude":                             NOSET,
	"resource":                            NOSET,
	"rule":                                NOSET,
	"udm_region_id":                       NOSET,
	"speed_limit":                         NOSET,
	"include":                             NOSET,
	"detail":                              NOSET,
	"prefix":                              NOSET,
	"update_paths":                        NOSET,
	"bucket":                              NOSET,
	"instance_id":                         NOSET,
	"auto_remove_file_system_mount_point": NOSET,
	"file_system_id":                      CHECKSET,
	"create_time":                         CHECKSET,
	"vault_id":                            CHECKSET,
	"schedule":                            "I|1602673264|PT2H",
	"nas_backup_plan_name":                "tf-testAccCase",
	"backup_type":                         "COMPLETE",
}

func AlicloudHBRNasBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = var.name
  encrypt_type  = "1"
}

data "alicloud_nas_file_systems" "default" {
  protocol_type       = "NFS"
  description_regex   = alicloud_nas_file_system.default.description
}
`, name)
}
