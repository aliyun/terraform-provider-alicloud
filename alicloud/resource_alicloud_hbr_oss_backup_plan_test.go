package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudHBROssBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_oss_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudHBROssBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrOssBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrossbackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBROssBackupPlanBasicDependence0)
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
					"bucket":               "${alicloud_oss_bucket.default.bucket}",
					"backup_type":          "COMPLETE",
					"prefix":               "/home",
					"schedule":             "I|1602673264|PT2H",
					"oss_backup_plan_name": "tf-testAccHbrOss",
					"retention":            "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_type":          "COMPLETE",
						"prefix":               "/home",
						"schedule":             "I|1602673264|PT2H",
						"oss_backup_plan_name": "tf-testAccHbrOss",
						"retention":            "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_backup_plan_name": "tf-testAccHbrOss2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_backup_plan_name": "tf-testAccHbrOss2",
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
					"prefix": "var",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix": "var",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_backup_plan_name": "tf-testAccHbrOss3",
					"schedule":             "I|1602673264|PT2H",
					"retention":            "3",
					"prefix":               "root",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_backup_plan_name": "tf-testAccHbrOss3",
						"schedule":             "I|1602673264|PT2H",
						"retention":            "3",
						"prefix":               "root",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudHBROssBackupPlanMap0 = map[string]string{
	"path.#":                              NOSET,
	"disk_id":                             NOSET,
	"options":                             NOSET,
	"exclude":                             NOSET,
	"resource":                            NOSET,
	"rule":                                NOSET,
	"udm_region_id":                       NOSET,
	"speed_limit":                         NOSET,
	"include":                             NOSET,
	"detail":                              NOSET,
	"prefix":                              NOSET,
	"update_paths":                        NOSET,
	"instance_id":                         NOSET,
	"auto_remove_file_system_mount_point": NOSET,
	"file_system_id":                      NOSET,
	"create_time":                         NOSET,
	"bucket":                              CHECKSET,
	"vault_id":                            CHECKSET,
	"retention":                           "",
	"schedule":                            "I|1602673264|PT2H",
	"oss_backup_plan_name":                "tf-testAccHbrOss",
	"backup_type":                         "COMPLETE",
}

func AlicloudHBROssBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}
`, name)
}
