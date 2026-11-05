package dlfnext_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDlfNextCatalog(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDlfNextCatalogConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_dlf_catalog.default", "id", fmt.Sprintf("tf-dlf-cat-%d", r)),
					resource.TestCheckResourceAttr("alicloud_dlf_catalog.default", "name", fmt.Sprintf("tf-dlf-cat-%d", r)),
					resource.TestCheckResourceAttr("alicloud_dlf_catalog.default", "type", "PAIMON"),
					resource.TestCheckResourceAttr("alicloud_dlf_catalog.default", "is_shared", "false"),
					resource.TestCheckResourceAttrSet("alicloud_dlf_catalog.default", "catalog_id"),
					resource.TestCheckResourceAttr("alicloud_dlf_catalog.default", "status", "RUNNING"),
					resource.TestCheckResourceAttrSet("alicloud_dlf_catalog.default", "owner"),
					resource.TestCheckResourceAttrSet("alicloud_dlf_catalog.default", "created_at"),
				),
			},
			{
				// A catalog with no user-set options should be stable: options is
				// read-only (Computed) and carries server-managed storage defaults,
				// so a refresh must produce an empty plan.
				Config:   testAccDlfNextCatalogConfig(r),
				PlanOnly: true,
			},
			{
				ResourceName:      "alicloud_dlf_catalog.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDlfNextCatalogConfig(r int) string {
	return fmt.Sprintf(`
# DlfNext ACC targets cn-hangzhou. The shared cn-beijing test account
# accumulates deleted PAIMON catalogs (~48h cleanup) and rate-limits new
# catalog creation with INITIALIZE_FAILED once ~18 are pending cleanup;
# cn-hangzhou is a clean, supported region for the same account.
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_dlf_catalog" "default" {
  name      = "tf-dlf-cat-%[1]d"
  type      = "PAIMON"
  is_shared = false
}
`, r)
}
