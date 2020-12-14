package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPrivatelinkVpcEndpoint_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccPrivatelinkVpcEndpointTest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivatelinkVpcEndpointBasicDependence)
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
					"service_id":        "${alicloud_privatelink_vpc_endpoint_service.default.id}",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"security_group_id": []string{"${alicloud_security_group.default.id}"},
					"vpc_endpoint_name": "TerraformTest",
					"depends_on":        []string{"alicloud_privatelink_vpc_endpoint_service.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_id":          CHECKSET,
						"vpc_id":              CHECKSET,
						"security_group_id.#": "1",
						"vpc_endpoint_name":   "TerraformTest",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "Terraform Test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "Terraform Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name": "TerraformTestUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name": "TerraformTestUpdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "Terraform Test Update",
					"vpc_endpoint_name":    "TerraformTestUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "Terraform Test Update",
						"vpc_endpoint_name":    "TerraformTestUpdate",
					}),
				),
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointMap = map[string]string{
	"bandwidth":                CHECKSET,
	"connection_status":        CHECKSET,
	"endpoint_business_status": CHECKSET,
	"endpoint_domain":          CHECKSET,
	"service_name":             CHECKSET,
	"status":                   CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  is_default = true
	}
	resource "alicloud_security_group" "default" {
	  name        = "tftest"
	  description = "privatelink test security group"
	  vpc_id      = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	 service_description = "%s"
	 connect_bandwidth = 103
     auto_accept_connection = false
	}
`, name)
}
