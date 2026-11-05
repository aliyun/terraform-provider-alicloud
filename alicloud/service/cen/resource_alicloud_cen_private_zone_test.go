package cen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCenPrivateZone_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.CenSupportRegions)
	resourceId := "alicloud_cen_private_zone.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCenPrivateZoneMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenPrivateZone")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenPrivateZone%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCenPrivateZoneBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"cen_id":           acctest.CHECKSET,
						"access_region_id": acctest.CHECKSET,
						"host_vpc_id":      acctest.CHECKSET,
						"host_region_id":   acctest.CHECKSET,
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
	"status": acctest.CHECKSET,
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
