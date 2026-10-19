package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudEnsNetworkInterfacesDataSource covers the data source with
// enable_details: a matching name_regex returns the created ENI, a non-matching
// name_regex returns an empty list.
func TestAccAliCloudEnsNetworkInterfacesDataSource(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-ens-ni-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccEnsNetworkInterfacesDataSourceConfig(name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ens_network_interfaces.default", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ens_network_interfaces.default", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_ens_network_interfaces.default", "names.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_ens_network_interfaces.default", "interfaces.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ens_network_interfaces.default", "interfaces.0.security_group_ids.#"),
				),
			},
			{
				Config: testAccEnsNetworkInterfacesDataSourceConfig(name, name+"_does_not_exist"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ens_network_interfaces.default", "interfaces.#", "0"),
				),
			},
		},
	})
}

// testAccEnsNetworkInterfacesDataSourceConfig creates an ENI and then queries
// it through the data source filtered by name_regex, with enable_details=true
// so the per-interface security_group_ids are populated.
func testAccEnsNetworkInterfacesDataSourceConfig(name, regex string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_security_group" "default" {
  security_group_name = "tf-testacc-ens-ni-ds-sg-%[1]s"
}

resource "alicloud_ens_network" "default" {
  network_name  = "tf-testacc-ens-ni-ds-network-%[1]s"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default" {
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = "tf-testacc-ens-ni-ds-vswitch-%[1]s"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_ens_network_interface" "default" {
  description            = "desc"
  network_interface_name = "%[1]s"
  security_group_ids     = [alicloud_ens_security_group.default.id]
  vswitch_id             = alicloud_ens_vswitch.default.id
}

data "alicloud_ens_network_interfaces" "default" {
  enable_details = true
  name_regex     = "%[2]s"
  ids            = [alicloud_ens_network_interface.default.id]
}
`, name, regex)
}
