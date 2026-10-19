package eflo_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// VccRouteEntry is a child of an Eflo VCC (灵骏连接) instance. There is no
// alicloud_eflo_vcc resource or alicloud_eflo_vccs data source in this provider
// yet, so a test cannot create its parent in-test. EFLO_VCC_ID gates the live CRUD
// exercise: set it to a writable VCC instance id in the test account to run the
// full create/read/import/delete cycle; unset, the test skips rather than fail.
func TestAccAliCloudEfloVccRouteEntry_basic(t *testing.T) {
	vccId := os.Getenv("EFLO_VCC_ID")
	if vccId == "" {
		t.Skip("EFLO_VCC_ID must be set to a writable Eflo VCC instance id; no alicloud_eflo_vcc resource exists to create one in-test")
	}

	// A random third octet keeps the destination CIDR clear of routes the VCC already carries.
	r := tfacctest.RandIntRange(2, 250)
	cidr := fmt.Sprintf("192.168.%d.0/24", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEfloVccRouteEntryConfig(vccId, cidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_eflo_vcc_route_entry.default", "id"),
					resource.TestCheckResourceAttr("alicloud_eflo_vcc_route_entry.default", "vcc_id", vccId),
					resource.TestCheckResourceAttr("alicloud_eflo_vcc_route_entry.default", "destination_cidr_block", cidr),
					resource.TestCheckResourceAttrSet("alicloud_eflo_vcc_route_entry.default", "vcc_route_entry_id"),
					resource.TestCheckResourceAttrSet("alicloud_eflo_vcc_route_entry.default", "region_id"),
					resource.TestCheckResourceAttrSet("alicloud_eflo_vcc_route_entry.default", "status"),
				),
			},
			{
				ResourceName:      "alicloud_eflo_vcc_route_entry.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEfloVccRouteEntryConfig(vccId, cidr string) string {
	return fmt.Sprintf(`
resource "alicloud_eflo_vcc_route_entry" "default" {
  vcc_id                 = "%[1]s"
  destination_cidr_block = "%[2]s"
}
`, vccId, cidr)
}
