package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks Project. >>> Resource test cases, automatically generated.
// Case Project_对接Terraform 7051
func TestAccAliCloudDataWorksProject_basic7051(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_project.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksProjectMap7051)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksProjectBasicDependence7051)
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
					"project_name": name,
					"project_mode": "2",
					"description":  "对接terraform测试",
					"display_name": "对接terraform测试",
					"status":       "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name": name,
						"project_mode": "2",
						"description":  "对接terraform测试",
						"display_name": "对接terraform测试",
						"status":       "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "变更terraform测试描述",
					"display_name": "变更terraform测试显示名",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "变更terraform测试描述",
						"display_name": "变更terraform测试显示名",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "被手动禁用",
					"status":      "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "被手动禁用",
						"status":      "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "解除禁用",
					"status":      "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "解除禁用",
						"status":      "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "对接terraform测试",
					"display_name": "对接terraform测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "对接terraform测试",
						"display_name": "对接terraform测试",
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

var AlicloudDataWorksProjectMap7051 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudDataWorksProjectBasicDependence7051(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test DataWorks Project. <<< Resource test cases, automatically generated.
