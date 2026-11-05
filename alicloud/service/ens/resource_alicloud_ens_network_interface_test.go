package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ENS region exercising the edge network; ENS resources use dedicated region
// ids rather than the standard ALICLOUD_REGION.
const ensNetworkInterfaceTestRegion = "cn-wuxi-telecom_unicom_cmcc-2"

func TestAccAliCloudEnsNetworkInterface_basic(t *testing.T) {
	resourceId := "alicloud_ens_network_interface.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snetworkinterface%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEnsNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEnsNetworkInterfaceBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "network_interface_id"),
					resource.TestCheckResourceAttr(resourceId, "network_interface_name", name),
					resource.TestCheckResourceAttr(resourceId, "description", "NetworkInterfaceDescription_autotest"),
					resource.TestCheckResourceAttr(resourceId, "vmnc_learn", "false"),
					resource.TestCheckResourceAttr(resourceId, "private_ip_address_count", "0"),
					resource.TestCheckResourceAttr(resourceId, "status", "Available"),
					resource.TestCheckResourceAttr(resourceId, "security_group_ids.#", "1"),
				),
			},
			{
				Config: testAccAliCloudEnsNetworkInterfaceUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "network_interface_name", name+"_update"),
					resource.TestCheckResourceAttr(resourceId, "description", "NetworkInterfaceDescription_UPDATE_autotest"),
					resource.TestCheckResourceAttr(resourceId, "vmnc_learn", "true"),
					resource.TestCheckResourceAttr(resourceId, "private_ip_address_count", "1"),
					resource.TestCheckResourceAttr(resourceId, "private_ip_addresses.#", "1"),
				),
			},
			{
				Config: testAccAliCloudEnsNetworkInterfaceReduce(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "network_interface_name", name+"_update"),
					resource.TestCheckResourceAttr(resourceId, "vmnc_learn", "false"),
					resource.TestCheckResourceAttr(resourceId, "private_ip_address_count", "0"),
					resource.TestCheckResourceAttr(resourceId, "private_ip_addresses.#", "0"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_ip_address_count"},
			},
		},
	})
}

func testAccCheckEnsNetworkInterfaceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ens_network_interface" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := ens.NewService(client).DescribeEnsNetworkInterface(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing network interface %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("network interface %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudEnsNetworkInterfaceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  description   = "for ens network interface test"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "%s"
}

resource "alicloud_ens_vswitch" "default" {
  description   = "for ens network interface test"
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = var.name
  ens_region_id = "%s"
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_ens_security_group" "default" {
  description          = "for ens network interface test"
  security_group_name  = var.name
}
`, name, ensNetworkInterfaceTestRegion, ensNetworkInterfaceTestRegion)
}

func testAccAliCloudEnsNetworkInterfaceBasic(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_ens_network_interface" "default" {
  vswitch_id           = alicloud_ens_vswitch.default.id
  security_group_ids   = [alicloud_ens_security_group.default.id]
  network_interface_name = var.name
  description          = "NetworkInterfaceDescription_autotest"
}
`, testAccAliCloudEnsNetworkInterfaceDependence(name))
}

func testAccAliCloudEnsNetworkInterfaceUpdate(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_ens_network_interface" "default" {
  vswitch_id             = alicloud_ens_vswitch.default.id
  security_group_ids     = [alicloud_ens_security_group.default.id]
  network_interface_name = "%s_update"
  description            = "NetworkInterfaceDescription_UPDATE_autotest"
  vmnc_learn             = true
  private_ip_address_count = 1
}
`, testAccAliCloudEnsNetworkInterfaceDependence(name), name)
}

func testAccAliCloudEnsNetworkInterfaceReduce(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_ens_network_interface" "default" {
  vswitch_id             = alicloud_ens_vswitch.default.id
  security_group_ids     = [alicloud_ens_security_group.default.id]
  network_interface_name = "%s_update"
  description            = "NetworkInterfaceDescription_UPDATE_autotest"
  vmnc_learn             = false
  private_ip_address_count = 0
}
`, testAccAliCloudEnsNetworkInterfaceDependence(name), name)
}
