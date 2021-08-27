package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRdsParameterGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudRdsParameterGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence0)
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
					"engine":         "mysql",
					"engine_version": `5.7`,
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `3000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86400`,
						},
					},
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "mysql",
						"engine_version":       "5.7",
						"param_detail.#":       "2",
						"parameter_group_name": name,
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
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `4000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86460`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_desc": "update_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_desc": "update_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `3000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86400`,
						},
					},
					"parameter_group_desc": "test",
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#":       "2",
						"parameter_group_desc": "test",
						"parameter_group_name": name,
					}),
				),
			},
		},
	})
}

var AlicloudRdsParameterGroupMap0 = map[string]string{}

func AlicloudRdsParameterGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
`, name)
}
