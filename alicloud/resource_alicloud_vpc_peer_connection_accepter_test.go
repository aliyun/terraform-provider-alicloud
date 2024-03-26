package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Case 1
func TestAccAliCloudVpcPeerConnectionAccepter_basic2(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_vpc_peer_connection_accepter.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcPeerConnectionAccepterMap0)
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p.(*schema.Provider), nil
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
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
  peer_connection_name = var.name
  vpc_id               = alicloud_vpc.default.id
  accepting_ali_uid    = data.alicloud_account.default.id
  accepting_region_id  = var.accepting_region
  accepting_vpc_id     = alicloud_vpc.default1.id
  description          = var.name
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
