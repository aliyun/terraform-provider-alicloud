package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCenPrivateZone_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenSupportRegions)
	resourceId := "alicloud_cen_private_zone.default"
	ra := resourceAttrInit(resourceId, AliCloudCenPrivateZoneMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenPrivateZone")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenPrivateZone%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenPrivateZoneBasicDependence0)
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
					"cen_id":           "${alicloud_cen_instance_attachment.default.instance_id}",
					"access_region_id": "${data.alicloud_regions.default.regions.0.id}",
					"host_vpc_id":      "${alicloud_cen_instance_attachment.default.child_instance_id}",
					"host_region_id":   "${data.alicloud_regions.default.regions.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":           CHECKSET,
						"access_region_id": CHECKSET,
						"host_vpc_id":      CHECKSET,
						"host_region_id":   CHECKSET,
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

var AliCloudCenPrivateZoneMap0 = map[string]string{
	"status": CHECKSET,
}

func AliCloudCenPrivateZoneBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_regions" "default" {
  		current = true
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.17.3.0/24"
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
  		description       = var.name
	}

	resource "alicloud_cen_instance_attachment" "default" {
  		instance_id              = alicloud_cen_instance.default.id
  		child_instance_id        = alicloud_vpc.default.id
  		child_instance_type      = "VPC"
  		child_instance_region_id = data.alicloud_regions.default.regions.0.id
	}
`, name)
}
