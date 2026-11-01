package maxcompute_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/maxcompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test MaxCompute Role. >>> Resource test cases, automatically generated.
// Case policyRole_terraform测试_typeAdmin 9618
func TestAccAliCloudMaxComputeRole_basic9618(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMaxComputeRoleMap9618)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return maxcompute.NewMaxComputeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMaxComputeRole")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMaxComputeRoleBasicDependence9618)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"project_name": acctest.CHECKSET,
						"policy":       acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Resource\\\":[\\\"acs:odps:*:projects/project_name/authorization/packages\\\",\\\"acs:odps:*:projects/project_name/authorization/packages/*\\\",\\\"acs:odps:*:projects/project_name/authorization/packages/*/*/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMaxComputeRoleMap9584)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return maxcompute.NewMaxComputeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMaxComputeRole")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMaxComputeRoleBasicDependence9584)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"project_name": acctest.CHECKSET,
						"policy":       acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:Update\\\"],\\\"Resource\\\":[\\\"acs:odps:*:projects/test_project_a/tables/tb_*\\\"],\\\"Effect\\\":\\\"Allow\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": acctest.CHECKSET,
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
