// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudSSO UserProvisioning. >>> Resource test cases, automatically generated.
// Case 全生命周期测试 4630
func TestAccAliCloudCloudSSOUserProvisioning_basic4630(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_user_provisioning.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOUserProvisioningMap4630)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSSOUserProvisioning")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOUserProvisioningBasicDependence4630)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"principal_id":         "${alicloud_cloud_sso_user.default.user_id}",
					"target_type":          "RD-Account",
					"deletion_strategy":    "Keep",
					"duplication_strategy": "KeepBoth",
					"principal_type":       "User",
					"target_id":            "${data.alicloud_account.default.id}",
					"directory_id":         "${alicloud_cloud_sso_user.default.directory_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"principal_id":         CHECKSET,
						"target_type":          "RD-Account",
						"deletion_strategy":    "Keep",
						"duplication_strategy": "KeepBoth",
						"principal_type":       "User",
						"target_id":            CHECKSET,
						"directory_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_strategy": "Delete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_strategy": "Delete",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"duplication_strategy": "Takeover",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duplication_strategy": "Takeover",
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

func TestAccAliCloudCloudSSOUserProvisioning_basic4630_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_user_provisioning.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOUserProvisioningMap4630)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSSOUserProvisioning")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOUserProvisioningBasicDependence4630)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "description",
					"principal_id":         "${alicloud_cloud_sso_user.default.user_id}",
					"target_type":          "RD-Account",
					"deletion_strategy":    "Keep",
					"duplication_strategy": "KeepBoth",
					"principal_type":       "User",
					"target_id":            "${data.alicloud_account.default.id}",
					"directory_id":         "${alicloud_cloud_sso_user.default.directory_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "description",
						"principal_id":         CHECKSET,
						"target_type":          "RD-Account",
						"deletion_strategy":    "Keep",
						"duplication_strategy": "KeepBoth",
						"principal_type":       "User",
						"target_id":            CHECKSET,
						"directory_id":         CHECKSET,
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

var AliCloudCloudSSOUserProvisioningMap4630 = map[string]string{
	"status":                         CHECKSET,
	"user_provisioning_statistics.#": CHECKSET,
	"user_provisioning_id":           CHECKSET,
	"create_time":                    CHECKSET,
}

func AliCloudCloudSSOUserProvisioningBasicDependence4630(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

	data "alicloud_account" "default" {
	}

	data "alicloud_cloud_sso_directories" "default" {
	}

	resource "alicloud_cloud_sso_directory" "default" {
  		count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  		directory_name = var.name
	}

	resource "alicloud_cloud_sso_group" "default" {
  		directory_id = local.directory_id
  		group_name   = var.name
  		description  = var.name
	}

	resource "alicloud_cloud_sso_user" "default" {
  		directory_id = local.directory_id
  		user_name    = var.name
	}

	locals {
  		directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
	}
`, name)
}

// Test CloudSSO UserProvisioning. <<< Resource test cases, automatically generated.
