package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudEnsNetworkInterface_basic0 covers the ENI lifecycle:
// create with all Optional/Required arguments, update Description and
// NetworkInterfaceName via ModifyNetworkInterfaceAttribute, clear the
// Optional arguments back to only the Required ones, and finally import.
func TestAccAliCloudEnsNetworkInterface_basic0(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-ens-ni-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccEnsNetworkInterfaceConfig(name, "desc", name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ens_network_interface.default", "description", "desc"),
					resource.TestCheckResourceAttr("alicloud_ens_network_interface.default", "network_interface_name", name),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "security_group_ids.#"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "vswitch_id"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "status"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "ens_region_id"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "network_id"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "primary_ip"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "primary_ip_type"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "mac_address"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "create_time"),
				),
			},
			{
				Config: testAccEnsNetworkInterfaceConfig(name, "desc_update", name+"_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ens_network_interface.default", "description", "desc_update"),
					resource.TestCheckResourceAttr("alicloud_ens_network_interface.default", "network_interface_name", name+"_update"),
				),
			},
			{
				// Clear the Optional arguments; only the Required
				// security_group_ids and the create-only vswitch_id remain.
				Config: testAccEnsNetworkInterfaceConfigCleared(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "security_group_ids.#"),
					resource.TestCheckResourceAttrSet("alicloud_ens_network_interface.default", "vswitch_id"),
				),
			},
			{
				ResourceName:      "alicloud_ens_network_interface.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccEnsNetworkInterfaceConfig builds a config with a security group, an
// ENS network and vswitch as dependencies, and the ENI under test. description
// and network_interface_name are the two updatable arguments.
func testAccEnsNetworkInterfaceConfig(name, description, niName string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_security_group" "default" {
  security_group_name = "tf-testacc-ens-ni-sg-%[1]s"
}

resource "alicloud_ens_network" "default" {
  network_name  = "tf-testacc-ens-ni-network-%[1]s"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default" {
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = "tf-testacc-ens-ni-vswitch-%[1]s"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_ens_network_interface" "default" {
  description            = "%[2]s"
  network_interface_name = "%[3]s"
  security_group_ids     = [alicloud_ens_security_group.default.id]
  vswitch_id             = alicloud_ens_vswitch.default.id
}
`, name, description, niName)
}

// testAccEnsNetworkInterfaceConfigCleared omits description and
// network_interface_name, exercising the empty-value update path.
func testAccEnsNetworkInterfaceConfigCleared(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_security_group" "default" {
  security_group_name = "tf-testacc-ens-ni-sg-%[1]s"
}

resource "alicloud_ens_network" "default" {
  network_name  = "tf-testacc-ens-ni-network-%[1]s"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default" {
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = "tf-testacc-ens-ni-vswitch-%[1]s"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_ens_network_interface" "default" {
  security_group_ids = [alicloud_ens_security_group.default.id]
  vswitch_id         = alicloud_ens_vswitch.default.id
}
`, name)
}
