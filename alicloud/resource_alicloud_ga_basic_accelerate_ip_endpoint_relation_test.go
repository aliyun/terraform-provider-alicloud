package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaBasicAccelerateIpEndpointRelation_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_accelerate_ip_endpoint_relation.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudGaBasicAccelerateIpEndpointRelationMap)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccGaBasicAccelerateIpEndpointRelation-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudGaBasicAccelerateIpEndpointRelationBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGaBasicAccelerateIpEndpointRelationDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":   "${alicloud_ga_basic_accelerate_ip.default.accelerator_id}",
					"accelerate_ip_id": "${alicloud_ga_basic_accelerate_ip.default.id}",
					"endpoint_id":      "${alicloud_ga_basic_endpoint.default.endpoint_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaBasicAccelerateIpEndpointRelationExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"accelerator_id":   CHECKSET,
						"accelerate_ip_id": CHECKSET,
						"endpoint_id":      CHECKSET,
					}),
				),
			},
		},
	})
}

func testAccCheckGaBasicAccelerateIpEndpointRelationDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckGaBasicAccelerateIpEndpointRelationDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckGaBasicAccelerateIpEndpointRelationDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	gaService := GaService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ga_basic_accelerate_ip_endpoint_relation" {
			continue
		}
		resp, err := gaService.DescribeGaBasicAccelerateIpEndpointRelation(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Ga Basic Accelerate Ip Endpoint Relation still exist,  ID %s ", fmt.Sprintf("%v:%v:%v", resp["AcceleratorId"], resp["AccelerateIpId"], resp["EndpointId"]))
		}
	}

	return nil
}

func testAccCheckGaBasicAccelerateIpEndpointRelationExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_ga_basic_accelerate_ip_endpoint_relation ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			gaService := GaService{client}

			resp, err := gaService.DescribeGaBasicAccelerateIpEndpointRelation(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}
			res = resp
			return nil
		}
		return fmt.Errorf("alicloud_ga_basic_accelerate_ip_endpoint_relation not found")
	}
}

var resourceAlicloudGaBasicAccelerateIpEndpointRelationMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudGaBasicAccelerateIpEndpointRelationBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	provider "alicloud" {
  		alias  = "sz"
  		region = "cn-shenzhen"
	}

	provider "alicloud" {
  		alias  = "hz"
  		region = "cn-hangzhou"
	}

	data "alicloud_vpcs" "default" {
  		provider   = "alicloud.sz"
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		provider = "alicloud.sz"
  		vpc_id   = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		provider = "alicloud.sz"
  		vpc_id   = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_ecs_network_interface" "default" {
  		provider           = "alicloud.sz"
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

	resource "alicloud_ga_basic_ip_set" "default" {
  		accelerator_id       = alicloud_ga_basic_accelerator.default.id
  		accelerate_region_id = "cn-hangzhou"
  		isp_type             = "BGP"
  		bandwidth            = "5"
	}

	resource "alicloud_ga_basic_accelerate_ip" "default" {
  		accelerator_id = alicloud_ga_basic_ip_set.default.accelerator_id
  		ip_set_id      = alicloud_ga_basic_ip_set.default.id
	}

	resource "alicloud_ga_basic_endpoint_group" "default" {
  		accelerator_id        = alicloud_ga_basic_accelerator.default.id
  		endpoint_group_region = "cn-shenzhen"
	}

	resource "alicloud_ga_basic_endpoint" "default" {
  		accelerator_id            = alicloud_ga_basic_accelerator.default.id
  		endpoint_group_id         = alicloud_ga_basic_endpoint_group.default.id
  		endpoint_type             = "ENI"
  		endpoint_address          = alicloud_ecs_network_interface.default.id
  		endpoint_sub_address_type = "primary"
  		endpoint_sub_address      = "192.168.0.1"
  		basic_endpoint_name       = var.name
	}
`, name)
}
