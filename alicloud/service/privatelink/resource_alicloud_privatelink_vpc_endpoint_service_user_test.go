package privatelink_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	privatelink "github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatelink"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPrivatelinkVpcEndpointServiceUser_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service_user.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceUserMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivatelinkVpcEndpointServiceUser")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccVpcEndpointServiceUserTest%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivatelinkVpcEndpointServiceUserBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_id":    "${alicloud_ram_user.user.id}",
					"service_id": "${alicloud_privatelink_vpc_endpoint_service.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":    acctest.CHECKSET,
						"service_id": acctest.CHECKSET,
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
	  name_regex = "^default-NODELETING$"
	}
	resource "alicloud_security_group" "default" {
	  name        = "tf-testAcc-for-privatelink"
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

// Test PrivateLink VpcEndpointServiceUser. >>> Resource test cases, automatically generated.
// Case 生命周期测试- VpcEndpointServiceUser 4876
func TestAccAliCloudPrivateLinkVpcEndpointServiceUser_basic4876(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service_user.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointServiceUserMap4876)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return privatelink.NewPrivateLinkServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePrivateLinkVpcEndpointServiceUser")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointserviceuser%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudPrivateLinkVpcEndpointServiceUserBasicDependence4876)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_id":    "1419992691371925",
					"dry_run":    "false",
					"service_id": "${alicloud_privatelink_vpc_endpoint_service.defaultSaMoag.id}",
					"user_arn":   "acs:ram:*:1419992691371925:*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":    acctest.CHECKSET,
						"dry_run":    "false",
						"service_id": acctest.CHECKSET,
						"user_arn":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "user_arn"},
			},
		},
	})
}

var AlicloudPrivateLinkVpcEndpointServiceUserMap4876 = map[string]string{}

func AlicloudPrivateLinkVpcEndpointServiceUserBasicDependence4876(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "current" {
}

resource "alicloud_privatelink_vpc_endpoint_service" "defaultSaMoag" {
  service_description   = "test-zejun"
  service_resource_type = "nlb"
}


`, name)
}

// Test PrivateLink VpcEndpointServiceUser. <<< Resource test cases, automatically generated.
