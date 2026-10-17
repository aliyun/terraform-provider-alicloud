package dataworks_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudDataWorksTable exercises the Framework-native alicloud_dataworks_table
// resource lifecycle: create, update, import and destroy.
//
// A DataWorks table is created in a workspace's bound MaxCompute engine and the API
// requires DataWorks Enterprise/Flagship Edition (Forbidden.Access otherwise). There is no
// Terraform resource to bind a MaxCompute engine to a freshly created DataWorks project,
// so the fixture is a pre-existing Enterprise workspace that already has a MaxCompute
// engine bound. The workspace ID and engine name are supplied through environment
// variables; when they are absent the test is skipped rather than failing.
//
// DATA_WORKS_TABLE_PROJECT_ID — numeric DataWorks workspace ID that owns the table.
// DATA_WORKS_TABLE_ENGINE     — MaxCompute project (engine) name; table_guid is
//
//	built as odps.<engine>.tf-testacc<rand>.
//
// DATA_WORKS_TABLE_CATEGORY_ID        — optional, an existing table category ID.
// DATA_WORKS_TABLE_PHYSICS_LEVEL_ID    — optional, an existing physics level ID.
func TestAccAliCloudDataWorksTable(t *testing.T) {
	if os.Getenv("DATA_WORKS_TABLE_PROJECT_ID") == "" || os.Getenv("DATA_WORKS_TABLE_ENGINE") == "" {
		t.Skip("skip: set DATA_WORKS_TABLE_PROJECT_ID and DATA_WORKS_TABLE_ENGINE to a DataWorks Enterprise workspace with a bound MaxCompute engine")
	}

	r := tfacctest.RandIntRange(100000, 999999)
	tableName := fmt.Sprintf("tf_testacc%d", r)
	engine := os.Getenv("DATA_WORKS_TABLE_ENGINE")
	projectId := os.Getenv("DATA_WORKS_TABLE_PROJECT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudDataWorksTableConfig(engine, projectId, tableName, fmt.Sprintf("tf-testacc-%d", r), fmt.Sprintf("tf-intro-%d", r)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_dataworks_table.default", "id"),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "table_guid", fmt.Sprintf("odps.%s.%s", engine, tableName)),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "table_name", tableName),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "project_id", projectId),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "comment", fmt.Sprintf("tf-testacc-%d", r)),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "life_cycle", "30"),
					// AC2: write-only attributes are not returned by the read API but must
					// survive the apply-time Read without drifting back to null.
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "is_view", "0"),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "intro_wiki", fmt.Sprintf("tf-intro-%d", r)),
					resource.TestCheckResourceAttrSet("alicloud_dataworks_table.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_dataworks_table.default", "last_modify_time"),
					resource.TestCheckResourceAttrSet("alicloud_dataworks_table.default", "region_id"),
				),
			},
			{
				Config: testAccAliCloudDataWorksTableConfigUpdate(engine, projectId, tableName, fmt.Sprintf("tf-testacc-%d-upd", r), fmt.Sprintf("tf-intro-%d-upd", r)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "comment", fmt.Sprintf("tf-testacc-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "life_cycle", "60"),
					// is_view remains 0 after the update Read: write-only drift check (AC2).
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "is_view", "0"),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "intro_wiki", fmt.Sprintf("tf-intro-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_dataworks_table.default", "table_name", tableName),
				),
			},
			{
				ResourceName:      "alicloud_dataworks_table.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Write-only attributes are not returned by GetMetaTableFullInfo, so the
					// post-import Read cannot repopulate them; they are preserved from state
					// during normal applies but not during a cold import.
					"is_view", "category_id", "physics_level_id",
					"intro_wiki", "first_level_theme_id", "second_level_theme_id",
					"level_id", "level_type",
				},
			},
		},
	})
}

func testAccAliCloudDataWorksTableConfig(engine, projectId, tableName, comment, intro string) string {
	return fmt.Sprintf(`
resource "alicloud_dataworks_table" "default" {
  table_guid = "odps.%[1]s.%[3]s"
  project_id = %[2]s
  comment    = "%[4]s"
  life_cycle = 30
  is_view    = 0
  intro_wiki = "%[5]s"
}
`, engine, projectId, tableName, comment, intro)
}

func testAccAliCloudDataWorksTableConfigUpdate(engine, projectId, tableName, comment, intro string) string {
	return fmt.Sprintf(`
resource "alicloud_dataworks_table" "default" {
  table_guid = "odps.%[1]s.%[3]s"
  project_id = %[2]s
  comment    = "%[4]s"
  life_cycle = 60
  is_view    = 0
  intro_wiki = "%[5]s"
}
`, engine, projectId, tableName, comment, intro)
}
