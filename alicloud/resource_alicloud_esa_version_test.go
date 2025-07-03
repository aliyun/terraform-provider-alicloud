package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA SiteVersion. >>> Resource test cases, automatically generated.
// Case resource_SiteVersion_test
func TestAccAliCloudEsaVersionresource_SiteVersion_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_version.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaVersionresource_SiteVersion_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaVersionresource_SiteVersion_testBasicDependence)
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
					"site_id":        "${alicloud_esa_site.default.id}",
					"description":    "测试版本",
					"origin_version": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "测试版本update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "测试版本update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"origin_version"},
			},
		},
	})
}

var AliCloudEsaVersionresource_SiteVersion_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudEsaVersionresource_SiteVersion_testBasicDependence(name string) string {
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

// Test ESA SiteVersion. <<< Resource test cases, automatically generated.
