package dlfnext_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDlfNextTable(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDlfNextTableConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_dlf_table.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_dlf_table.default", "table_id"),
					resource.TestCheckResourceAttrSet("alicloud_dlf_table.default", "owner"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "name", fmt.Sprintf("tf-dlf-tbl-%d", r)),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "comment", "test table"),
					// columns. A partition key must also be listed as a field:
					// the API rejects "Table column [...] should include all
					// partition fields [...]". A primary-key column's type is
					// returned as "BIGINT NOT NULL"; formatTypeString strips the
					// " NOT NULL" suffix so the type round-trips as "BIGINT".
					// status is not asserted: GetTable returns status: null for
					// newly created tables.
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.#", "4"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.0.name", "id"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.0.type", "BIGINT"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.1.name", "name"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.1.type", "STRING"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.2.name", "tags"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.2.type", "MAP<STRING,STRING>"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.2.description", "tag map"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.3.name", "dt"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "fields.3.type", "STRING"),
					// keys
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "partition_keys.#", "1"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "partition_keys.0", "dt"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "primary_keys.#", "1"),
					resource.TestCheckResourceAttr("alicloud_dlf_table.default", "primary_keys.0", "id"),
				),
			},
			{
				Config:   testAccDlfNextTableConfig(r),
				PlanOnly: true,
			},
			{
				ResourceName:      "alicloud_dlf_table.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDlfNextTableBase(r int) string {
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

resource "alicloud_dlf_database" "default" {
  catalog_id = alicloud_dlf_catalog.default.catalog_id
  name       = "tf-dlf-db-%[1]d"
}
`, r)
}

func testAccDlfNextTableConfig(r int) string {
	return fmt.Sprintf(`
%[1]s
resource "alicloud_dlf_table" "default" {
  catalog_id    = alicloud_dlf_catalog.default.catalog_id
  database_name = alicloud_dlf_database.default.name
  name          = "tf-dlf-tbl-%[2]d"
  comment       = "test table"

  # fields is a repeated block; each block has name, type and an optional
  # description. State addressing (fields.0.name etc.) is the same as for
  # list-of-objects syntax. A partition key column must also be declared as a
  # field (see partition_keys = ["dt"] below).
  fields {
    name = "id"
    type = "BIGINT"
  }

  fields {
    name = "name"
    type = "STRING"
  }

  fields {
    name        = "tags"
    type        = "MAP<STRING,STRING>"
    description = "tag map"
  }

  fields {
    name = "dt"
    type = "STRING"
  }

  partition_keys = ["dt"]
  primary_keys   = ["id"]
}
`, testAccDlfNextTableBase(r), r)
}
