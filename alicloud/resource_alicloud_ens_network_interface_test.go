package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens NetworkInterface. >>> Resource test cases (hand-written, 100% attribute coverage).

func TestAccAliCloudEnsNetworkInterface_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_network_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNetworkInterfaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNetworkInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensnetworkinterface%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsNetworkInterfaceBasicDependence0)
	vswitchId := "${alicloud_ens_vswitch.default.id}"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// create with all Optional/Required fields set
				Config: testAccConfig(map[string]interface{}{
					"description":            "desc",
					"network_interface_name": name,
					"security_group_ids":     []interface{}{"${alicloud_ens_security_group.default.id}"},
					"vswitch_id":             vswitchId,
					"vmnc_learn":             true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "desc",
						"network_interface_name": name,
						"security_group_ids":     CHECKSET,
						"vswitch_id":             CHECKSET,
						"vmnc_learn":             "true",
					}),
				),
			},
			{
				// update description and network_interface_name (ModifyNetworkInterfaceAttribute)
				Config: testAccConfig(map[string]interface{}{
					"description":            "desc_update",
					"network_interface_name": name + "_update",
					"security_group_ids":     []interface{}{"${alicloud_ens_security_group.default.id}"},
					"vswitch_id":             vswitchId,
					"vmnc_learn":             true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "desc_update",
						"network_interface_name": name + "_update",
					}),
				),
			},
			{
				// update vmnc_learn to false (ModifyNetworkInterfaceVmncLearn)
				Config: testAccConfig(map[string]interface{}{
					"description":            "desc_update",
					"network_interface_name": name + "_update",
					"security_group_ids":     []interface{}{"${alicloud_ens_security_group.default.id}"},
					"vswitch_id":             vswitchId,
					"vmnc_learn":             false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vmnc_learn": "false",
					}),
				),
			},
			{
				// clear optional fields (description, network_interface_name); keep Required
				Config: testAccConfig(map[string]interface{}{
					"description":            REMOVEKEY,
					"network_interface_name": REMOVEKEY,
					"security_group_ids":     []interface{}{"${alicloud_ens_security_group.default.id}"},
					"vswitch_id":             vswitchId,
					"vmnc_learn":             false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEnsNetworkInterfaceMap0 = map[string]string{
	"status":          CHECKSET,
	"create_time":     CHECKSET,
	"ens_region_id":   CHECKSET,
	"instance_id":     CHECKSET,
	"mac_address":     CHECKSET,
	"network_id":      CHECKSET,
	"primary_ip":      CHECKSET,
	"primary_ip_type": CHECKSET,
}

func AlicloudEnsNetworkInterfaceBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_security_group" "default" {
  security_group_name = "tf-testacc-ens-ni-sg"
}

resource "alicloud_ens_network" "default" {
  network_name  = "tf-testacc-ens-ni-network"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default" {
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = "tf-testacc-ens-ni-vswitch"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default.id
}

`, name)
}

// Test Ens NetworkInterface. <<< Resource test cases.
