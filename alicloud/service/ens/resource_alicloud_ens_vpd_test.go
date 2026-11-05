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

// ENS region used by the existing alicloud_ens_network acceptance tests; VPD
// is created under a network in the same ENS region.
const ensVpdTestRegion = "cn-chenzhou-telecom_unicom_cmcc"

func TestAccAliCloudEnsVpd_basic(t *testing.T) {
	resourceId := "alicloud_ens_vpd.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccensvpd%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEnsVpdDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEnsVpdBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "ens_region_id", ensVpdTestRegion),
					resource.TestCheckResourceAttrSet(resourceId, "network_id"),
					resource.TestCheckResourceAttr(resourceId, "name", name),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
			{
				Config: testAccAliCloudEnsVpdBasicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "name", name+"_update"),
					resource.TestCheckResourceAttr(resourceId, "description", "vpd_description_update"),
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

func testAccCheckEnsVpdDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ens_vpd" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := ens.NewService(client).DescribeEnsVpd(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing ens vpd %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("ens vpd %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudEnsVpdBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ens_network" "default" {
  ens_region_id = "%[2]s"
  cidr_block    = "192.168.2.0/24"
  network_name  = var.name
}

resource "alicloud_ens_vpd" "default" {
  ens_region_id = alicloud_ens_network.default.ens_region_id
  network_id    = alicloud_ens_network.default.id
  name          = var.name
}
`, name, ensVpdTestRegion)
}

func testAccAliCloudEnsVpdBasicUpdate(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ens_network" "default" {
  ens_region_id = "%[2]s"
  cidr_block    = "192.168.2.0/24"
  network_name  = var.name
}

resource "alicloud_ens_vpd" "default" {
  ens_region_id = alicloud_ens_network.default.ens_region_id
  network_id    = alicloud_ens_network.default.id
  name          = "${var.name}_update"
  description   = "vpd_description_update"
}
`, name, ensVpdTestRegion)
}
