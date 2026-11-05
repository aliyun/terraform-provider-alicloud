package dataworks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test DataWorks Project. >>> Resource test cases, automatically generated.
// Case 租户新版空间API对接terraform_新版_成都 8903
func TestAccAliCloudDataWorksProject_basic8903(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_project.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDataWorksProjectMap8903)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dataworks.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDataWorksProject")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDataWorksProjectBasicDependence8903)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":             "tf_desc",
					"project_name":            name,
					"pai_task_enabled":        "false",
					"display_name":            "tf_new_api_display",
					"dev_role_disabled":       "true",
					"dev_environment_enabled": "false",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"status":                  "Available",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":             "tf_desc",
						"project_name":            name,
						"pai_task_enabled":        "false",
						"display_name":            "tf_new_api_display",
						"dev_role_disabled":       "true",
						"dev_environment_enabled": "false",
						"resource_group_id":       acctest.CHECKSET,
						"status":                  "Available",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Forbidden",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Forbidden",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "tf_desc_update",
					"display_name": "tf_new_api_display_update",
					"status":       "Available",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "tf_desc_update",
						"display_name": "tf_new_api_display_update",
						"status":       "Available",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pai_task_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pai_task_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dev_environment_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dev_environment_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dev_role_disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dev_role_disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AlicloudDataWorksProjectMap8903 = map[string]string{}

func AlicloudDataWorksProjectBasicDependence8903(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test DataWorks Project. <<< Resource test cases, automatically generated.
