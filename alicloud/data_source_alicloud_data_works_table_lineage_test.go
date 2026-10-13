package alicloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAlicloudDataWorksTableLineageDataSource verifies the
// alicloud_data_works_table_lineage data source can query upstream and
// downstream lineage for a real table in the test account.
//
// The test is environment-guarded: it only runs when
// ALICLOUD_DATA_WORKS_LINEAGE_TABLE_NAME (and optionally
// ALICLOUD_DATA_WORKS_LINEAGE_DATABASE_NAME) are set, because the test account
// is not guaranteed to own a table with lineage relationships. This mirrors
// how other pure-query data sources handle pre-existing infrastructure that
// cannot be fixture-created inside the ACC run.
func TestAccAlicloudDataWorksTableLineageDataSource(t *testing.T) {
	tableName := os.Getenv("ALICLOUD_DATA_WORKS_LINEAGE_TABLE_NAME")
	databaseName := os.Getenv("ALICLOUD_DATA_WORKS_LINEAGE_DATABASE_NAME")
	if tableName == "" {
		t.Skip("Skipping test: ALICLOUD_DATA_WORKS_LINEAGE_TABLE_NAME is not set; this data source queries pre-existing table lineage and cannot fixture its own test data.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksTableLineageDataSourceConfig(tableName, databaseName, "up"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_data_works_table_lineage.up"),
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_table_lineage.up", "id"),
					resource.TestCheckResourceAttr("data.alicloud_data_works_table_lineage.up", "direction", "up"),
				),
			},
			{
				Config: testAccCheckAlicloudDataWorksTableLineageDataSourceConfig(tableName, databaseName, "down"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_data_works_table_lineage.down"),
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_table_lineage.down", "id"),
					resource.TestCheckResourceAttr("data.alicloud_data_works_table_lineage.down", "direction", "down"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDataWorksTableLineageDataSourceConfig(tableName, databaseName, direction string) string {
	config := `
data "alicloud_data_works_table_lineage" "` + direction + `" {
  direction  = "` + direction + `"
  table_name = "` + tableName + `"
`
	if databaseName != "" {
		config += `  database_name = "` + databaseName + `"
`
	}
	config += "}\n"
	return config
}
