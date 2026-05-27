package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Vpc Ipv6Address. >>> Resource test cases, automatically generated.
// Case IPv6Address生命周期 4694
func TestAccAliCloudVpcIpv6Address_basic4694(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipv6_address.default"
	ra := resourceAttrInit(resourceId, AliCloudVpcIpv6AddressMap4694)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpv6Address")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVpcIpv6AddressBasicDependence4694)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id": "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_address_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_address_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_address_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_address_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
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

func TestAccAliCloudVpcIpv6Address_basic4694_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipv6_address.default"
	ra := resourceAttrInit(resourceId, AliCloudVpcIpv6AddressMap4694)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpv6Address")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVpcIpv6AddressBasicDependence4694)
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
					"address_type":             "Ipv6Address",
					"ipv6_address":             "${cidrhost(alicloud_vswitch.default.ipv6_cidr_block, 128)}",
					"vswitch_id":               "${alicloud_vswitch.default.id}",
					"ipv6_address_description": name,
					"ipv6_address_name":        name,
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type":             "Ipv6Address",
						"ipv6_address":             CHECKSET,
						"vswitch_id":               CHECKSET,
						"ipv6_address_description": name,
						"ipv6_address_name":        name,
						"resource_group_id":        CHECKSET,
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "Test",
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

var AliCloudVpcIpv6AddressMap4694 = map[string]string{
	"address_type":      CHECKSET,
	"ipv6_address":      CHECKSET,
	"resource_group_id": CHECKSET,
	"create_time":       CHECKSET,
	"status":            CHECKSET,
}

func AliCloudVpcIpv6AddressBasicDependence4694(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  cidr_block  = "192.168.0.0/16"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "default" {
  vswitch_name         = var.name
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "192.168.192.0/24"
  zone_id              = data.alicloud_zones.default.zones.0.id
  ipv6_cidr_block_mask = "1"
}
`, name)
}

// Test Vpc Ipv6Address. <<< Resource test cases, automatically generated.
