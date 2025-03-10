package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Vpc Ipv6Address. >>> Resource test cases, automatically generated.
// Case IPv6Address生命周期 4694
func TestAccAliCloudVpcIpv6Address_basic4694(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipv6_address.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpv6AddressMap4694)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpv6Address")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpv6AddressBasicDependence4694)
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
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vswitch_id":               "${alicloud_vswitch.vswich.id}",
					"ipv6_address_description": "create_description",
					"ipv6_address_name":        name,
					"address_type":             "Ipv6Address",
					"ipv6_address":             "${cidrhost(alicloud_vswitch.vswich.ipv6_cidr_block, 128)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":        CHECKSET,
						"vswitch_id":               CHECKSET,
						"ipv6_address":             CHECKSET,
						"ipv6_address_description": "create_description",
						"ipv6_address_name":        name,
						"address_type":             "Ipv6Address",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"ipv6_address_description": "modify_description",
					"ipv6_address_name":        name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":        CHECKSET,
						"ipv6_address_description": "modify_description",
						"ipv6_address_name":        name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
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
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudVpcIpv6AddressMap4694 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcIpv6AddressBasicDependence4694(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-g"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  cidr_block  = "172.16.0.0/12"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "vswich" {
  vpc_id               = alicloud_vpc.vpc.id
  cidr_block           = "172.16.0.0/24"
  zone_id              = data.alicloud_zones.default.zones.0.id
  vswitch_name         = "tf-testacc"
  ipv6_cidr_block_mask = "1"
}


`, name)
}

// Test Vpc Ipv6Address. <<< Resource test cases, automatically generated.
