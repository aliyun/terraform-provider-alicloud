package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDataWorksQualityCheckResults_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksQualityCheckResultsConfigByEntity(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_data_works_quality_check_results.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_quality_check_results.default", "id"),
				),
			},
		},
	})
}

func TestAccAlicloudDataWorksQualityCheckResults_ruleId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksQualityCheckResultsConfigByRule(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_data_works_quality_check_results.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_quality_check_results.default", "id"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDataWorksQualityCheckResultsConfigByEntity() string {
	return fmt.Sprintf(`
data "alicloud_data_works_quality_check_results" "default" {
	project_name = "tfdsaccchain10293"
	start_date   = "2024-01-01 00:00:00"
	end_date     = "2024-01-02 00:00:00"
	entity_id    = "1"
	page_size    = 10
	page_number  = 1
	output_file  = "/tmp/quality_check_results.txt"
}
`)
}

func testAccCheckAlicloudDataWorksQualityCheckResultsConfigByRule() string {
	return fmt.Sprintf(`
data "alicloud_data_works_quality_check_results" "default" {
	project_name = "tfdsaccchain10293"
	start_date   = "2024-01-01 00:00:00"
	end_date     = "2024-01-02 00:00:00"
	rule_id      = "1"
}
`)
}
