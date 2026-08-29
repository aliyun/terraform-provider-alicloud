package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig Portal. >>> Resource test cases, hand-written.
func TestAccAliCloudApigPortal_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_portal.default"
	ra := resourceAttrInit(resourceId, AlicloudApigPortalMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigPortal")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapigportal%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigPortalBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": "tfacc portal description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "tfacc portal description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name + "_update",
					"description": "tfacc portal description update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name + "_update",
						"description": "tfacc portal description update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"portal_setting_config": []map[string]interface{}{
						{
							"builtin_auth_enabled":       "true",
							"auto_approve_developers":    "true",
							"auto_approve_subscriptions": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"portal_setting_config.#":                            "1",
						"portal_setting_config.0.builtin_auth_enabled":       "true",
						"portal_setting_config.0.auto_approve_developers":    "true",
						"portal_setting_config.0.auto_approve_subscriptions": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"portal_setting_config": []map[string]interface{}{
						{
							"builtin_auth_enabled":       "true",
							"auto_approve_developers":    "false",
							"auto_approve_subscriptions": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"portal_setting_config.0.builtin_auth_enabled":       "true",
						"portal_setting_config.0.auto_approve_developers":    "false",
						"portal_setting_config.0.auto_approve_subscriptions": "true",
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

var AlicloudApigPortalMap = map[string]string{}

func AlicloudApigPortalBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test Apig Portal. <<< Resource test cases, hand-written.
