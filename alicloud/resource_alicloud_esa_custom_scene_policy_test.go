package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA CustomScenePolicy. >>> Resource test cases, automatically generated.
// Case resource_CustomScenePolicy_test_1
func TestAccAliCloudESACustomScenePolicyresource_CustomScenePolicy_test_1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_scene_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudESACustomScenePolicyresource_CustomScenePolicy_test_1Map)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomScenePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESACustomScenePolicyresource_CustomScenePolicy_test_1BasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_scene_policy_name": "example-policy",
					"end_time":                 "2025-08-07T17:00:00Z",
					"create_time":              "2025-07-07T17:00:00Z",
					"site_ids":                 "${alicloud_esa_site.default.id}",
					"template":                 "promotion",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template": "promotion",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"create_time": "2025-07-08T17:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time": "2025-08-08T17:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"site_ids": "618651327383200,569031282787888",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_scene_policy_name": "test-policy",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AliCloudESACustomScenePolicyresource_CustomScenePolicy_test_1Map = map[string]string{
	"id": CHECKSET,
}

func AliCloudESACustomScenePolicyresource_CustomScenePolicy_test_1BasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
 plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
 site_name   = var.name
 instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
 coverage    = "overseas"
 access_type = "NS"
 version_management = true
}

`, name)
}

// Test ESA CustomScenePolicy. <<< Resource test cases, automatically generated.
