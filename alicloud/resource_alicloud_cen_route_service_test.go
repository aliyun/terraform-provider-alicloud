package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCENCenRouteService_basic(t *testing.T) {
	var v cbn.RouteServiceEntry
	resourceId := "alicloud_cen_route_service.default"
	ra := resourceAttrInit(resourceId, CenRouteServiceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenRouteService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", CenRouteServiceBasicdependence)
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
					"access_region_id": "${alicloud_cen_instance_attachment.vpc.child_instance_region_id}",
					"cen_id":           "${alicloud_cen_instance_attachment.vpc.instance_id}",
					"host":             "100.118.28.52/32",
					"host_region_id":   "${alicloud_cen_instance_attachment.vpc.child_instance_region_id}",
					"host_vpc_id":      "${alicloud_cen_instance_attachment.vpc.child_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_region_id": defaultRegionToTest,
						"cen_id":           CHECKSET,
						"host":             "100.118.28.52/32",
						"host_region_id":   defaultRegionToTest,
						"host_vpc_id":      CHECKSET,
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

var CenRouteServiceMap = map[string]string{
	"status": CHECKSET,
}

func CenRouteServiceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_vpcs" "default"{
	is_default = true
}
resource "alicloud_cen_instance" "default" {
    cen_instance_name = var.name
}
resource "alicloud_cen_instance_attachment" "vpc" {
    instance_id = alicloud_cen_instance.default.id
    child_instance_id = data.alicloud_vpcs.default.vpcs.0.id
	child_instance_type = "VPC"
    child_instance_region_id = data.alicloud_vpcs.default.vpcs.0.region_id
}
`, name)
}
