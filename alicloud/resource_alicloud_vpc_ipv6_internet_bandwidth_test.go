package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVPCIpv6InternetBandwidth_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipv6_internet_bandwidth.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCIpv6InternetBandwidthMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpv6InternetBandwidth")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipv6internetbandwidth%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCIpv6InternetBandwidthBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ECS_WITH_IPV6_ADDRESS")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_address_id":      "${data.alicloud_vpc_ipv6_addresses.default.addresses.0.id}",
					"ipv6_gateway_id":      "${data.alicloud_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id}",
					"internet_charge_type": "PayByBandwidth",
					"bandwidth":            "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_address_id":      CHECKSET,
						"ipv6_gateway_id":      CHECKSET,
						"internet_charge_type": "PayByBandwidth",
						"bandwidth":            "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "50",
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

var AlicloudVPCIpv6InternetBandwidthMap0 = map[string]string{
	"internet_charge_type": CHECKSET,
	"status":               CHECKSET,
}

func AlicloudVPCIpv6InternetBandwidthBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "alicloud_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.alicloud_instances.default.instances.0.id
  status                 = "Available"
}
`, name)
}
