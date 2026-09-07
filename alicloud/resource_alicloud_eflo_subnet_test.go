package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudEfloSubnet_basic2581(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EfloSupportRegions)
	resourceId := "alicloud_eflo_subnet.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloSubnetMap2581)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloSubnet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccEfloSubnet%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloSubnetBasicDependence2581)
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
					"subnet_name": "${var.name}",
					"zone_id":     "${data.alicloud_zones.default.zones.0.id}",
					"cidr":        "10.0.0.0/16",
					"vpd_id":      "${alicloud_eflo_vpd.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subnet_name": name,
						"zone_id":     CHECKSET,
						"cidr":        "10.0.0.0/16",
						"vpd_id":      CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"subnet_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subnet_name": name + "_update",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAlicloudEfloSubnet_basic2582(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EfloSupportRegions)
	resourceId := "alicloud_eflo_subnet.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloSubnetMap2581)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloSubnet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccEfloSubnet%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloSubnetBasicDependence2581)
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
					"subnet_name": "${var.name}",
					"zone_id":     "${data.alicloud_zones.default.zones.0.id}",
					"cidr":        "10.0.0.0/16",
					"type":        "OOB",
					"vpd_id":      "${alicloud_eflo_vpd.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subnet_name": name,
						"zone_id":     CHECKSET,
						"cidr":        "10.0.0.0/16",
						"type":        "OOB",
						"vpd_id":      CHECKSET,
					}),
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

var AlicloudEfloSubnetMap2581 = map[string]string{}

func AlicloudEfloSubnetBasicDependence2581(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eflo_vpd" "default" {
  cidr      = "10.0.0.0/8"
  vpd_name  = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

`, name)
}
