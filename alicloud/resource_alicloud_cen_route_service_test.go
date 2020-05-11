package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCenRouteService_basic(t *testing.T) {
	var v cbn.RouteServiceEntry
	resourceId := "alicloud_cen_route_service.default"
	ra := resourceAttrInit(resourceId, CenRouteServiceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenRouteService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenRouteService%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CenRouteServiceBasicdependence)
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
					"access_region_ids": []string{defaultRegionToTest},
					"cen_id":            "${alicloud_cen_instance.default.id}",
					"host":              "100.64.1.0/24",
					"host_region_id":    defaultRegionToTest,
					"host_vpc_id":       "${alicloud_vpc.default.id}",
					"depends_on":        []string{"alicloud_cen_instance_attachment.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_region_ids.#": "1",
						"cen_id":              CHECKSET,
						"host":                "100.64.1.0/24",
						"host_region_id":      defaultRegionToTest,
						"host_vpc_id":         CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_region_ids"},
			},
		},
	})
}

var CenRouteServiceMap = map[string]string{
	"status": CHECKSET,
}

func CenRouteServiceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_cen_instance" "default" {
	name = "${var.name}"
}
resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}
resource "alicloud_cen_instance_attachment" "default" {
	instance_id = "${alicloud_cen_instance.default.id}"
	child_instance_id = "${alicloud_vpc.default.id}"
	child_instance_region_id = "%s"
  	depends_on = [
		"alicloud_cen_instance.default",
		"alicloud_vpc.default"]
}
`, name, defaultRegionToTest)
}
