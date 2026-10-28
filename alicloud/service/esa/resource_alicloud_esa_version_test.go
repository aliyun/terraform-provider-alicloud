package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA SiteVersion. >>> Resource test cases, automatically generated.
// Case resource_SiteVersion_test
func TestAccAliCloudEsaVersionresource_SiteVersion_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_version.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaVersionresource_SiteVersion_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaVersion")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaVersionresource_SiteVersion_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"id": acctest.CHECKSET,
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
