package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudHBREcsBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_ecs_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudHBREcsBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrEcsBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%shbrecsbackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBREcsBackupPlanBasicDependence0)
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
					"vault_id":             "${alicloud_hbr_vault.default.id}",
					"instance_id":          "${data.alicloud_instances.default.instances.0.id}",
					"backup_type":          "COMPLETE",
					"schedule":             "I|1602673264|PT2H",
					"ecs_backup_plan_name": "tf-testAcc-hbr-backup-plan",
					"retention":            "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_type":          "COMPLETE",
						"schedule":             "I|1602673264|PT2H",
						"ecs_backup_plan_name": "tf-testAcc-hbr-backup-plan",
						"retention":            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecs_backup_plan_name": "tf-testAcc-hbr-backup-plan2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_backup_plan_name": "tf-testAcc-hbr-backup-plan2",
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
					"include": "[\\\"/home\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"include": "[\"/home\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude": "[\\\"/proc\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude": "[\"/proc\"]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"speed_limit": "0:24:5120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"speed_limit": "0:24:5120",
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
					"ecs_backup_plan_name": "tf-testAcc-hbr-backup-plan3",
					"schedule":             "I|1602673264|PT2H",
					"retention":            "4",
					"path":                 []string{"/home/test2", "/home/test2"},
					"include":              "[\\\"/proc\\\"]",
					"exclude":              "[\\\"/home\\\", \\\"/var/\\\"]",
					"speed_limit":          "0:24:1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_backup_plan_name": "tf-testAcc-hbr-backup-plan3",
						"schedule":             "I|1602673264|PT2H",
						"retention":            "4",
						"path.#":               "2",
						"include":              "[\"/proc\"]",
						"exclude":              "[\"/home\", \"/var/\"]",
						"speed_limit":          "0:24:1024",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_paths"},
			},
		},
	})
}

var AlicloudHBREcsBackupPlanMap0 = map[string]string{
	"path.#":               NOSET,
	"retention":            "",
	"disk_id":              NOSET,
	"options":              "",
	"exclude":              "",
	"resource":             NOSET,
	"rule":                 NOSET,
	"file_system_id":       NOSET,
	"udm_region_id":        NOSET,
	"speed_limit":          "",
	"include":              "",
	"detail":               "",
	"prefix":               NOSET,
	"update_paths":         NOSET,
	"bucket":               NOSET,
	"instance_id":          CHECKSET,
	"schedule":             "I|1602673264|PT2H",
	"ecs_backup_plan_name": "tf-testAcc-hbr-backup-plan",
	"backup_type":          "COMPLETE",
	"vault_id":             CHECKSET,
}

func AlicloudHBREcsBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = "${var.name}"
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-hbr-ecs-backup-plan"
  status = "Running"
}
`, name)
}
