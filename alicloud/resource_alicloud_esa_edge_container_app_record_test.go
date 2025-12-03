package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA EdgeContainerAppRecord. >>> Resource test cases, automatically generated.
// Case resource_EdgeContainerAppRecord_test
func TestAccAliCloudESAEdgeContainerAppRecordresource_EdgeContainerAppRecord_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_edge_container_app_record.default"
	ra := resourceAttrInit(resourceId, AliCloudESAEdgeContainerAppRecordresource_EdgeContainerAppRecord_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaEdgeContainerAppRecord")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAEdgeContainerAppRecordresource_EdgeContainerAppRecord_testBasicDependence)
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
					"record_name": "tf." + name,
					"site_id":     "${alicloud_esa_site.resource_Site_OriginPool_test.id}",
					"app_id":      "${alicloud_esa_edge_container_app.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
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

var AliCloudESAEdgeContainerAppRecordresource_EdgeContainerAppRecord_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAEdgeContainerAppRecordresource_EdgeContainerAppRecord_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_OriginPool_test" {
  site_name   = var.name
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_edge_container_app" "default" {
  health_check_host = "example.com"
  health_check_type = "l7"
  service_port = "80"
  health_check_interval = "5"
  edge_container_app_name = "terraform-app1"
  health_check_http_code = "http_2xx"
  health_check_uri = "/"
  health_check_timeout = "3"
  health_check_succ_times = "2"
  remarks = "无备注信息1"
  health_check_method = "HEAD"
  health_check_port = "80"
  health_check_fail_times = "5"
  target_port = "3000"
}

`, name)
}

// Test ESA EdgeContainerAppRecord. <<< Resource test cases, automatically generated.
