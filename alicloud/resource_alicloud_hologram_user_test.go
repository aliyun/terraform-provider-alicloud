package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Hologram User. >>> Resource test cases.
func TestAccAliCloudHologramUser_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hologram_user.default"
	ra := resourceAttrInit(resourceId, AlicloudHologramUserMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHologramUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacchologramuser%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHologramUserBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_hologram_instance.default.id}",
					"user_name":   name,
					"super_user":  false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"user_name":   name,
						"super_user":  "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"super_user": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"super_user": "true",
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

var AlicloudHologramUserMap = map[string]string{
	"instance_id": CHECKSET,
	"user_name":   CHECKSET,
	"super_user":  CHECKSET,
}

func AlicloudHologramUserBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name
}

resource "alicloud_hologram_instance" "default" {
  zone_id       = alicloud_vswitch.defaultVSwitch.zone_id
  instance_name = var.name
  payment_type  = "PayAsYouGo"
  instance_type = "Warehouse"
  pricing_cycle = "Hour"
  cpu           = "32"
  gateway_count = "2"
  endpoints {
    type = "Intranet"
  }
  endpoints {
    type       = "VPCSingleTunnel"
    vswitch_id = alicloud_vswitch.defaultVSwitch.id
    vpc_id     = alicloud_vswitch.defaultVSwitch.vpc_id
  }
}
`, name)
}
