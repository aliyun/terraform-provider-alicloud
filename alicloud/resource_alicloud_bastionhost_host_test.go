package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHost_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthost%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostBasicDependence0)
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
					"source":               "Local",
					"instance_id":          "${alicloud_bastionhost_instance.default.id}",
					"host_name":            "tf-testAcc-j1ZFb6NcZb7jkd7BUT77QDoA",
					"active_address_type":  "Private",
					"host_private_address": "172.16.0.10",
					"os_type":              "Linux",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":               "Local",
						"instance_id":          CHECKSET,
						"host_name":            "tf-testAcc-j1ZFb6NcZb7jkd7BUT77QDoA",
						"active_address_type":  "Private",
						"host_private_address": "172.16.0.10",
						"os_type":              "Linux",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "tf-testAcc-j1ZFb6NcZb7jkd7BUT77QDoA",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "tf-testAcc-j1ZFb6NcZb7jkd7BUT77QDoA",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_public_address": "10.0.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_public_address": "10.0.0.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active_address_type": "Public",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active_address_type": "Public",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_public_address": "10.0.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_public_address": "10.0.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"os_type": "Windows",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"os_type": "Windows",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "tf-testAcc-zxttoHFctU2IGUrPU5PWrItq",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": "tf-testAcc-zxttoHFctU2IGUrPU5PWrItq",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":              "tf-testAcc-zxttoHFctU2IGUrPU5PWrItq",
					"active_address_type":  "Private",
					"host_private_address": "172.16.0.11",
					"os_type":              "Linux",
					"host_name":            "tf-testAcc-j1ZFb6NcZb7jkd7BUT77QDoA",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":              "tf-testAcc-zxttoHFctU2IGUrPU5PWrItq",
						"active_address_type":  "Private",
						"host_private_address": "172.16.0.11",
						"os_type":              "Linux",
						"host_name":            "tf-testAcc-j1ZFb6NcZb7jkd7BUT77QDoA",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"instance_region_id"},
			},
		},
	})
}

var AlicloudBastionhostHostMap0 = map[string]string{
	"instance_region_id": NOSET,
	"host_id":            CHECKSET,
	"instance_id":        CHECKSET,
}

func AlicloudBastionhostHostBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
 available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
 zone_id = local.zone_id
 vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "this" {
 count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id       = data.alicloud_vpcs.default.ids.0
 zone_id      = data.alicloud_zones.default.ids.0
 cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
resource "alicloud_security_group" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 name   = var.name
}
locals {
 vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
 zone_id     = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
}
resource "alicloud_bastionhost_instance" "default" {
 description        = var.name
 license_code       = "bhah_ent_50_asset"
 period             = "1"
 vswitch_id         = local.vswitch_id
 security_group_ids = [alicloud_security_group.default.id]
}
`, name)
}
