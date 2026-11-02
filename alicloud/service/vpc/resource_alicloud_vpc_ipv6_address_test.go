package vpc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Vpc Ipv6Address. >>> Resource test cases, automatically generated.
// Case IPv6Address生命周期 4694
func TestAccAliCloudVpcIpv6Address_basic4694(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipv6_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudVpcIpv6AddressMap4694)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return NewVpcServiceV2(testAccProvider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcIpv6Address")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudVpcIpv6AddressBasicDependence4694)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id": "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudVpcIpv6AddressMap4694)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return NewVpcServiceV2(testAccProvider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcIpv6Address")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudVpcIpv6AddressBasicDependence4694)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"ipv6_address":             acctest.CHECKSET,
						"vswitch_id":               acctest.CHECKSET,
						"ipv6_address_description": name,
						"ipv6_address_name":        name,
						"resource_group_id":        acctest.CHECKSET,
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
	"address_type":      acctest.CHECKSET,
	"ipv6_address":      acctest.CHECKSET,
	"resource_group_id": acctest.CHECKSET,
	"create_time":       acctest.CHECKSET,
	"status":            acctest.CHECKSET,
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
