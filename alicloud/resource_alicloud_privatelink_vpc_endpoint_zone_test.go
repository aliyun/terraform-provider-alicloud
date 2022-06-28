package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Currently, Private network slb can only be created through the console.
func SkipTestAccAlicloudPrivatelinkVpcEndpointZone_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_zone.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointZoneMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpointZone")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointZoneTest%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivatelinkVpcEndpointZoneBasicDependence)
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
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.default.id}",
					"vswitch_id":  "${data.alicloud_vswitches.default.ids.0}",
					"zone_id":     "eu-central-1a",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointZoneMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointZoneBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	 name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	 	 is_default = true
	}
	resource "alicloud_security_group" "default" {
	 name = "%[1]s"
	 description = "privatelink test security group"
	 vpc_id = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	service_description = "test for privatelink connection"
	connect_bandwidth = 103
	auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint_service_resource" "default" {
	 service_id    =  "${alicloud_privatelink_vpc_endpoint_service.default.id}"
	 resource_id   =  "lb-gw8nuyxxxx"
	 resource_type = "slb"
	}
	resource "alicloud_privatelink_vpc_endpoint" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service_resource.default.service_id
	 vpc_id = data.alicloud_vpcs.default.ids.0
	 security_group_id = [alicloud_security_group.default.id]
	 vpc_endpoint_name = "%[1]s"
	 depends_on = [alicloud_privatelink_vpc_endpoint_service.default]
	}
`, name)
}
