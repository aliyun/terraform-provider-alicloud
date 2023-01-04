package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudVpcPeerConnectionAccepter_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_peer_connection_accepter.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcPeerConnectionAccepterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcPeerConnectionAccepter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sVpcPeerConnectionAccepter%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcPeerConnectionAccepterBasicDependence)
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
					"instance_id": "${alicloud_vpc_peer_connection.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudVpcPeerConnectionAccepterMap = map[string]string{}

func AlicloudVpcPeerConnectionAccepterBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {}

variable "accepting_region" {
  default = "cn-beijing"
}

provider "alicloud" {
  alias  = "local"
  region = "%s"
}

provider "alicloud" {
  alias  = "accepting"
  region = var.accepting_region
}

data "alicloud_vpcs" "default" {
  provider   = alicloud.local
  name_regex = "default-NODELETING"
}

data "alicloud_vpcs" "defaultone" {
  provider   = alicloud.accepting
  name_regex = "default-NODELETING"
}


resource "alicloud_vpc_peer_connection" "default" {
  peer_connection_name = var.name
  vpc_id               = data.alicloud_vpcs.default.ids.0
  accepting_ali_uid    = data.alicloud_account.default.id
  accepting_region_id  = var.accepting_region
  accepting_vpc_id     = data.alicloud_vpcs.defaultone.ids.0
  description          = var.name
  provider             = alicloud.local
}


`, name, defaultRegionToTest)
}
