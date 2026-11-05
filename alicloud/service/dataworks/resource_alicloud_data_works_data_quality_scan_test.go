package dataworks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const dataQualityScanResourceId = "alicloud_data_works_data_quality_scan.default"

// TestAccAliCloudDataWorksDataQualityScan_basic exercises the full CRUD path
// against a real DataWorks project created for the test. The Spec references an
// ODPS data source and table; if the test account lacks the referenced ODPS
// data source/table the create step surfaces the API error rather than
// fabricating a pass.
func TestAccAliCloudDataWorksDataQualityScan_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dqs_%d", rand)
	updatedDescription := "Updated description for the data quality scan."
	updatedOwner := fmt.Sprintf("owner_%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDataWorksDataQualityScanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataWorksDataQualityScanBasicConfig(name, "Daily quality scan of ods tables.", "", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataQualityScanResourceId, "id"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "name", name),
					resource.TestCheckResourceAttrSet(dataQualityScanResourceId, "project_id"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "description", "Daily quality scan of ods tables."),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "parameters.#", "1"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "parameters.0.name", "e2e_autolabel"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "parameters.0.value", "10"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "trigger.type", "BySchedule"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "hooks.#", "1"),
					resource.TestCheckResourceAttrSet(dataQualityScanResourceId, "create_time"),
				),
			},
			{
				// Clear every optional block (parameters, trigger, hooks)
				// and assert the clear sticks, not just parameters.
				Config: testAccDataWorksDataQualityScanBasicConfig(name, updatedDescription, "", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "description", updatedDescription),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "parameters.#", "0"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "trigger.#", "0"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "hooks.#", "0"),
				),
			},
			{
				// Re-add the owner field in-place (no ForceNew on owner) and
				// keep the cleared blocks to exercise the scalar update path.
				Config: testAccDataWorksDataQualityScanBasicConfig(name, updatedDescription, updatedOwner, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "description", updatedDescription),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "owner", updatedOwner),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "parameters.#", "0"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "trigger.#", "0"),
					resource.TestCheckResourceAttr(dataQualityScanResourceId, "hooks.#", "0"),
				),
			},
			{
				ResourceName:      dataQualityScanResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDataWorksDataQualityScanDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_data_works_data_quality_scan" {
			continue
		}
		_, err := dataworks.NewService(client).DescribeDataQualityScanContext(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("data quality scan %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing data quality scan %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

func testAccDataWorksDataQualityScanBasicConfig(name, description, owner string, withOptional bool) string {
	optional := ""
	if withOptional {
		optional = `
  description = "Daily quality scan of ods tables."
  parameters {
    name  = "e2e_autolabel"
    value = "10"
  }
  trigger {
    type     = "BySchedule"
    task_ids = [1023777390]
  }
  hooks {
    condition = "results.any { r -> r.status == 'Fail' && r.rule.severity == 'High' }"
    type      = "BlockTaskInstance"
  }
`
	} else {
		optional = fmt.Sprintf(`
  description = "%s"
`, description)
		if owner != "" {
			optional += fmt.Sprintf("  owner       = \"%s\"\n", owner)
		}
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_data_works_project" "default" {
  description       = "default_pj_dqs"
  project_name      = var.name
  display_name      = "default_pj_dqs"
  pai_task_enabled  = false
}

resource "alicloud_data_works_data_quality_scan" "default" {
  name       = var.name
  project_id = alicloud_data_works_project.default.id
  spec       = <<JSON
{
  "datasets": [
    {
      "type": "Table",
      "dataSource": {"name": "odps_first", "envType": "Prod"},
      "tables": ["ods_d_user_info"],
      "filter": "pt = $[yyyymmdd-1]"
    }
  ],
  "rules": [
    {"assertion": "row_count > 0"}
  ]
}
JSON
%s}
`, name, optional)
}
