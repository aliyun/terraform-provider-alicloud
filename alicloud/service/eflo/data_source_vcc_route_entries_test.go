package eflo_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The data source lists route entries for a given Eflo VCC. Like the resource test,
// it is gated on EFLO_VCC_ID because no alicloud_eflo_vcc resource exists to create
// the parent in-test. A VCC with no custom route entries yields an empty ids list,
// which is a valid result; the assertions therefore check the framing attributes only.
func TestAccAliCloudEfloVccRouteEntries_basic(t *testing.T) {
	vccId := os.Getenv("EFLO_VCC_ID")
	if vccId == "" {
		t.Skip("EFLO_VCC_ID must be set to an Eflo VCC instance id; no alicloud_eflo_vcc resource exists to create one in-test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEfloVccRouteEntriesConfig(vccId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_eflo_vcc_route_entries.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_eflo_vcc_route_entries.default", "vcc_id", vccId),
					resource.TestCheckResourceAttrSet("data.alicloud_eflo_vcc_route_entries.default", "region_id"),
				),
			},
		},
	})
}

func testAccEfloVccRouteEntriesConfig(vccId string) string {
	return fmt.Sprintf(`
data "alicloud_eflo_vcc_route_entries" "default" {
  vcc_id = "%[1]s"
}
`, vccId)
}
