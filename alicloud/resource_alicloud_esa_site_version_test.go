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
func TestAccAliCloudESASiteVersionresource_SiteVersion_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site_version.default"
	ra := resourceAttrInit(resourceId, AliCloudESASiteVersionresource_SiteVersion_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSiteVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESASiteVersion%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESASiteVersionresource_SiteVersion_testBasicDependence)
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
					"site_id":        "618651327383200",
					"description":    "测试版本",
					"origin_version": "1",
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

var AliCloudESASiteVersionresource_SiteVersion_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESASiteVersionresource_SiteVersion_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test ESA SiteVersion. <<< Resource test cases, automatically generated.
