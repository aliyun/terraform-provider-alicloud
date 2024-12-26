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

// Case 1
func TestAccAliCloudVpcPeerConnectionAccepter_basic2(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_vpc_peer_connection_accepter.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcPeerConnectionAccepterMap0)
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
	name := fmt.Sprintf("tf-testacc%svpcpeerconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcPeerConnectionAccepterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVPCPeerConnectionAccepterDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_vpc_peer_connection.default.id}",
					"link_type":   "Gold",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"link_type":   "Gold",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
					"link_type":   "Platinum",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"link_type":   "Platinum",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_connection_accepter_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_connection_accepter_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force_delete": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force_delete": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "200",
					}),
				),
			},
		},
	})
}

func AlicloudVpcPeerConnectionAccepterBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
variable "accepting_region" {
  default = "cn-beijing"
}

data "alicloud_account" "default" {}

provider "alicloud" {
  alias = "local"
  region = "%s"
}

provider "alicloud" {
  alias = "accepting"
  region = var.accepting_region
}

resource "alicloud_vpc" "default" {
  provider   = alicloud.local
  vpc_name    = var.name
  enable_ipv6 = "true"
}

resource "alicloud_vpc" "default1" {
  provider   = alicloud.accepting
  vpc_name    = var.name
  enable_ipv6 = "true"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_peer_connection" "default" {
  vpc_id               = alicloud_vpc.default.id
  accepting_ali_uid    = data.alicloud_account.default.id
  accepting_region_id  = var.accepting_region
  accepting_vpc_id     = alicloud_vpc.default1.id
  provider             = alicloud.local
}

`, name, defaultRegionToTest)
}

var AlicloudVpcPeerConnectionAccepterMap0 = map[string]string{}

func testAccCheckVPCPeerConnectionAccepterDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckVPCPeerConnectionAccepterDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckVPCPeerConnectionAccepterDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	vpcPeerService := VpcPeerService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpc_peer_connection_accepter" {
			continue
		}

		_, err := vpcPeerService.DescribeVpcPeerConnection(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}
