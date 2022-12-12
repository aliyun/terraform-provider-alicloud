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

func TestAccAlicloudGaBasicEndpoint_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_endpoint.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudGaBasicEndpointMap)
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
	name := fmt.Sprintf("tf-testAccGaBasicEndpoint-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudGaBasicEndpointBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGaBasicEndpointDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":                  "alicloud.hz",
					"accelerator_id":            "${alicloud_ga_basic_accelerator.default.id}",
					"endpoint_group_id":         "${alicloud_ga_basic_endpoint_group.default.id}",
					"endpoint_type":             "ENI",
					"endpoint_address":          "${alicloud_ecs_network_interface.default.id}",
					"endpoint_sub_address_type": "secondary",
					"endpoint_sub_address":      "192.168.0.1",
					"basic_endpoint_name":       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaBasicEndpointExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"accelerator_id":            CHECKSET,
						"endpoint_group_id":         CHECKSET,
						"endpoint_type":             "ENI",
						"endpoint_address":          CHECKSET,
						"endpoint_sub_address_type": "secondary",
						"endpoint_sub_address":      "192.168.0.1",
						"basic_endpoint_name":       name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"basic_endpoint_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaBasicEndpointExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"basic_endpoint_name": name + "-update",
					}),
				),
			},
		},
	})
}

func testAccCheckGaBasicEndpointDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckGaBasicEndpointDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckGaBasicEndpointDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
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
		} else {
			return fmt.Errorf("Ga Basic Endpoint still exist,  ID %s ", fmt.Sprint(resp["EndPointId"]))
		}
	}

	return nil
}

func testAccCheckGaBasicEndpointExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_ga_basic_endpoint ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			gaService := GaService{client}

			resp, err := gaService.DescribeGaBasicEndpoint(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}
			res = resp
			return nil
		}
		return fmt.Errorf("alicloud_ga_basic_endpoint not found")
	}
}

var resourceAlicloudGaBasicEndpointMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudGaBasicEndpointBasicDependence(name string) string {
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

	resource "alicloud_ga_basic_endpoint_group" "default" {
  		accelerator_id        = alicloud_ga_basic_accelerator.default.id
  		endpoint_group_region = "cn-shenzhen"
	}
`, name)
}
