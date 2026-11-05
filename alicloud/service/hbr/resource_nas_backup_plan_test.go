package hbr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hbr"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass path -- pre-existing attribute outside this change

func TestAccAliCloudHbrNasBackupPlan_basic0(t *testing.T) {
	acctest.CheckoutAccount(t, true)
	defer acctest.CheckoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_nas_backup_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBRNasBackupPlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrNasBackupPlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%shbrnasbackupplan%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBRNasBackupPlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_type":             "COMPLETE",
					"vault_id":                "${alicloud_hbr_vault.default.id}",
					"file_system_id":          "${alicloud_nas_file_system.default.id}",
					"schedule":                "I|1602673264|PT2H",
					"nas_backup_plan_name":    "tf-testAccCase",
					"retention":               "1",
					"cross_account_type":      "SELF_ACCOUNT",
					"cross_account_user_id":   "${data.alicloud_account.default.id}",
					"cross_account_role_name": "${alicloud_ram_role.default.id}",
					"path":                    []string{"/"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_type":             "COMPLETE",
						"schedule":                "I|1602673264|PT2H",
						"nas_backup_plan_name":    "tf-testAccCase",
						"retention":               "1",
						"cross_account_type":      "SELF_ACCOUNT",
						"cross_account_user_id":   acctest.CHECKSET,
						"cross_account_role_name": acctest.CHECKSET,
						"path.#":                  "1",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudHBRNasBackupPlanMap0 = map[string]string{
	"path.#":                              acctest.NOSET,
	"retention":                           "",
	"disk_id":                             acctest.NOSET,
	"options":                             "",
	"exclude":                             acctest.NOSET,
	"resource":                            acctest.NOSET,
	"rule":                                acctest.NOSET,
	"udm_region_id":                       acctest.NOSET,
	"speed_limit":                         acctest.NOSET,
	"include":                             acctest.NOSET,
	"detail":                              acctest.NOSET,
	"prefix":                              acctest.NOSET,
	"update_paths":                        acctest.NOSET,
	"bucket":                              acctest.NOSET,
	"instance_id":                         acctest.NOSET,
	"auto_remove_file_system_mount_point": acctest.NOSET,
	"file_system_id":                      acctest.CHECKSET,
	"create_time":                         acctest.CHECKSET,
	"vault_id":                            acctest.CHECKSET,
	"schedule":                            "I|1602673264|PT2H",
	"nas_backup_plan_name":                "tf-testAccCase",
	"backup_type":                         "COMPLETE",
}

func AlicloudHBRNasBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
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
  		protocol_type     = "NFS"
  		description_regex = alicloud_nas_file_system.default.description
	}

	resource "alicloud_ram_role" "default" {
  		name     = var.name
  		document = <<EOF
		{
			"Statement": [
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": [
						"crossbackup.hbr.aliyuncs.com"
					]
				}
			}
			],
  			"Version": "1"
		}
  		EOF
  		force    = true
	}
`, name)
}
