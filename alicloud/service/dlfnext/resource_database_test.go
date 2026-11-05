package dlfnext_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDlfNextDatabase(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDlfNextDatabaseConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_dlf_database.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_dlf_database.default", "database_id"),
					resource.TestCheckResourceAttrSet("alicloud_dlf_database.default", "owner"),
					resource.TestCheckResourceAttr("alicloud_dlf_database.default", "name", fmt.Sprintf("tf-dlf-db-%d", r)),
					resource.TestCheckResourceAttr("alicloud_dlf_database.default", "options.k", "v"),
					// status is not asserted: GetDatabase returns status: null
					// for newly created databases (the OpenAPI example "active"
					// does not apply to fresh databases on this account).
				),
			},
			{
				Config: testAccDlfNextDatabaseConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_dlf_database.default", "options.k", "v2"),
					resource.TestCheckResourceAttr("alicloud_dlf_database.default", "options.k2", "v2"),
				),
			},
			{
				Config: testAccDlfNextDatabaseConfigClear(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("alicloud_dlf_database.default", "options.k"),
					resource.TestCheckNoResourceAttr("alicloud_dlf_database.default", "options.k2"),
				),
			},
			{
				ResourceName:      "alicloud_dlf_database.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDlfNextDatabaseBase(r int) string {
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

func testAccDlfNextDatabaseConfig(r int) string {
	return fmt.Sprintf(`
%[1]s
resource "alicloud_dlf_database" "default" {
  catalog_id = alicloud_dlf_catalog.default.catalog_id
  name       = "tf-dlf-db-%[2]d"
  options = {
    k = "v"
  }
}
`, testAccDlfNextDatabaseBase(r), r)
}

func testAccDlfNextDatabaseConfigUpdate(r int) string {
	return fmt.Sprintf(`
%[1]s
resource "alicloud_dlf_database" "default" {
  catalog_id = alicloud_dlf_catalog.default.catalog_id
  name       = "tf-dlf-db-%[2]d"
  options = {
    k  = "v2"
    k2 = "v2"
  }
}
`, testAccDlfNextDatabaseBase(r), r)
}

func testAccDlfNextDatabaseConfigClear(r int) string {
	return fmt.Sprintf(`
%[1]s
resource "alicloud_dlf_database" "default" {
  catalog_id = alicloud_dlf_catalog.default.catalog_id
  name       = "tf-dlf-db-%[2]d"
}
`, testAccDlfNextDatabaseBase(r), r)
}
