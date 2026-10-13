package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAliCloudGaBasicEndpoint_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_endpoint.default"
	ra := resourceAttrInit(resourceId, resourceAliCloudGaBasicEndpointMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccGaBasicEndpoint-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAliCloudGaBasicEndpointBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckGaBasicEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":            "${alicloud_ga_basic_accelerator.default.id}",
					"endpoint_group_id":         "${alicloud_ga_basic_endpoint_group.default.id}",
					"endpoint_type":             "ENI",
					"endpoint_address":          "${alicloud_ecs_network_interface.default.id}",
					"endpoint_sub_address_type": "secondary",
					"endpoint_sub_address":      "192.168.0.1",
					"endpoint_zone_id":          "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"basic_endpoint_name":       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaBasicEndpointExists(resourceId, &v),
					testAccCheck(map[string]string{
						"accelerator_id":            CHECKSET,
						"endpoint_group_id":         CHECKSET,
						"endpoint_type":             "ENI",
						"endpoint_address":          CHECKSET,
						"endpoint_sub_address_type": "secondary",
						"endpoint_sub_address":      "192.168.0.1",
						"endpoint_zone_id":          CHECKSET,
						"basic_endpoint_name":       name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"basic_endpoint_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaBasicEndpointExists(resourceId, &v),
					testAccCheck(map[string]string{
						"basic_endpoint_name": name + "-update",
					}),
				),
			},
		},
	})
}

func testAccCheckGaBasicEndpointDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	gaService := GaService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ga_basic_endpoint" {
			continue
		}
		resp, err := gaService.DescribeGaBasicEndpoint(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Ga Basic Endpoint still exist, ID %s ", fmt.Sprint(resp["EndPointId"]))
	}

	return nil
}

func testAccCheckGaBasicEndpointExists(n string, res *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_ga_basic_endpoint ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		gaService := GaService{client}

		resp, err := gaService.DescribeGaBasicEndpoint(rs.Primary.ID)
		if err != nil {
			return err
		}
		*res = resp
		return nil
	}
}

var resourceAliCloudGaBasicEndpointMap = map[string]string{
	"endpoint_id": CHECKSET,
	"status":      CHECKSET,
}

func resourceAliCloudGaBasicEndpointBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	`, name) + configAlternateRegionProvider("cn-shenzhen") + `

	data "alicloud_vpcs" "default" {
  		provider   = alicloudalt
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		provider = alicloudalt
  		vpc_id   = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		provider = alicloudalt
  		vpc_id   = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_ecs_network_interface" "default" {
  		provider           = alicloudalt
  		vswitch_id         = data.alicloud_vswitches.default.ids.0
  		security_group_ids = [alicloud_security_group.default.id]
	}

	resource "alicloud_ga_basic_accelerator" "default" {
  		duration               = 1
  		pricing_cycle          = "Month"
  		basic_accelerator_name = var.name
  		description            = var.name
  		bandwidth_billing_type = "CDT"
  		auto_pay               = true
  		auto_use_coupon        = "true"
  		auto_renew             = false
  		auto_renew_duration    = 1
	}

	resource "alicloud_ga_basic_endpoint_group" "default" {
  		accelerator_id        = alicloud_ga_basic_accelerator.default.id
  		endpoint_group_region = "cn-shenzhen"
	}
`
}
