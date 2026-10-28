package vpcpeer_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpcpeer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Case 1
func TestAccAliCloudVpcPeerConnectionAccepter_basic2(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_vpc_peer_connection_accepter.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcPeerConnectionAccepterMap0)
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcpeerconnection%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcPeerConnectionAccepterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckVPCPeerConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_vpc_peer_connection.default.id}",
					"link_type":   "Gold",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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

provider "alicloudalt" {
  region = var.accepting_region
}

resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
}

resource "alicloud_vpc" "default1" {
  provider   = alicloudalt
  vpc_name    = var.name
  enable_ipv6 = "true"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_peer_connection" "default" {
  vpc_id               = alicloud_vpc.default.id
  accepting_ali_uid    = data.alicloud_account.default.id
  accepting_region_id  = var.accepting_region
  accepting_vpc_id     = alicloud_vpc.default1.id
}

`, name)
}

var AlicloudVpcPeerConnectionAccepterMap0 = map[string]string{}

func testAccCheckVPCPeerConnectionAccepterDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := vpcpeer.NewService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpc_peer_connection_accepter" {
			continue
		}

		_, err := service.DescribeVpcPeerConnection(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

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
	service := vpcpeer.NewService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpc_peer_connection_accepter" {
			continue
		}

		_, err := service.DescribeVpcPeerConnection(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}
