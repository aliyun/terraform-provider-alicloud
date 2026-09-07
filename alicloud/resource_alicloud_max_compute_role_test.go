package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MaxCompute Role. >>> Resource test cases, automatically generated.
// Case policyRole_terraform测试_typeAdmin 9618
func TestAccAliCloudMaxComputeRole_basic9618(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleMap9618)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleBasicDependence9618)
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
					"role_name":    name,
					"type":         "admin",
					"project_name": "${alicloud_maxcompute_project.default.id}",
					"policy":       "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Resource\\\":[\\\"acs:odps:*:projects/project_name/authorization/roles\\\",\\\"acs:odps:*:projects/project_name/authorization/roles/*/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    name,
						"type":         "admin",
						"project_name": CHECKSET,
						"policy":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Resource\\\":[\\\"acs:odps:*:projects/project_name/authorization/packages\\\",\\\"acs:odps:*:projects/project_name/authorization/packages/*\\\",\\\"acs:odps:*:projects/project_name/authorization/packages/*/*/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"type"},
			},
		},
	})
}

var AlicloudMaxComputeRoleMap9618 = map[string]string{}

func AlicloudMaxComputeRoleBasicDependence9618(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "project_name" {
  default = "tf_test_project_20250109"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

`, name)
}

// Case policyRole_terraform测试_typeResource 9584
func TestAccAliCloudMaxComputeRole_basic9584(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleMap9584)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleBasicDependence9584)
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
					"role_name":    name,
					"type":         "resource",
					"project_name": "${alicloud_maxcompute_project.default.id}",
					"policy":       "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:Select\\\"],\\\"Resource\\\":[\\\"acs:odps:*:projects/test_project_a/tables/tb_*\\\"],\\\"Effect\\\":\\\"Allow\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    name,
						"type":         "resource",
						"project_name": CHECKSET,
						"policy":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:Update\\\"],\\\"Resource\\\":[\\\"acs:odps:*:projects/test_project_a/tables/tb_*\\\"],\\\"Effect\\\":\\\"Allow\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"type"},
			},
		},
	})
}

var AlicloudMaxComputeRoleMap9584 = map[string]string{}

func AlicloudMaxComputeRoleBasicDependence9584(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

`, name)
}

// Test MaxCompute Role. <<< Resource test cases, automatically generated.
