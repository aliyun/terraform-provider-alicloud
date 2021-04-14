package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPrivatelinkVpcEndpointServiceUser_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service_user.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceUserMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpointServiceUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccVpcEndpointServiceUserTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivatelinkVpcEndpointServiceUserBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_id":    "${alicloud_ram_user.user.id}",
					"service_id": "${alicloud_privatelink_vpc_endpoint_service.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":    CHECKSET,
						"service_id": CHECKSET,
					}),
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

var AlicloudPrivatelinkVpcEndpointServiceUserMap = map[string]string{}

func AlicloudPrivatelinkVpcEndpointServiceUserBasicDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name         = "%[1]s"
	  display_name = "user_display_name"
	  mobile       = "86-18688888888"
	  email        = "hello.uuu@aaa.com"
	  comments     = "yoyoyo"
	}
	data "alicloud_vpcs" "default" {
	  is_default = true
	}
	resource "alicloud_security_group" "default" {
	  name        = "tftest"
	  description = "privatelink test security group"
	  vpc_id      = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	 service_description = "%[1]s"
	 connect_bandwidth = 103
     auto_accept_connection = false
	}

`, name)
}
