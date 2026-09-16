package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDataworksDataServiceAuthorizedApisDataSource_basic(t *testing.T) {
	testAccPreCheckDataWorksDataServiceProjectId(t)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataworksDataServiceAuthorizedApisDataSourceBasicConfig(os.Getenv("ALICLOUD_DATA_WORKS_PROJECT_ID")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dataworks_data_service_authorized_apis.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_dataworks_data_service_authorized_apis.default", "total_count"),
				),
			},
		},
	})
}

func TestAccAliCloudDataworksDataServiceAuthorizedApisDataSource_empty(t *testing.T) {
	testAccPreCheckDataWorksDataServiceProjectId(t)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataworksDataServiceAuthorizedApisDataSourceEmptyConfig(os.Getenv("ALICLOUD_DATA_WORKS_PROJECT_ID")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dataworks_data_service_authorized_apis.default"),
					resource.TestCheckResourceAttr("data.alicloud_dataworks_data_service_authorized_apis.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dataworks_data_service_authorized_apis.default", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dataworks_data_service_authorized_apis.default", "apis.#", "0"),
				),
			},
		},
	})
}

func testAccPreCheckDataWorksDataServiceProjectId(t *testing.T) {
	if v := os.Getenv("ALICLOUD_DATA_WORKS_PROJECT_ID"); v == "" {
		t.Skip("ALICLOUD_DATA_WORKS_PROJECT_ID must be set for the DataWorks Data Service authorized APIs datasource acceptance test")
	}
}

// Covers: name_regex, page_number, page_size, output_file, project_id
func testAccCheckAlicloudDataworksDataServiceAuthorizedApisDataSourceBasicConfig(projectId string) string {
	return fmt.Sprintf(`
data "alicloud_dataworks_data_service_authorized_apis" "default" {
	project_id  = "%s"
	name_regex  = ".*"
	page_number = 1
	page_size   = 10
	output_file = "authorized_apis_basic.txt"
}
`, projectId)
}

// Covers: api_name_keyword, name_regex, ids, page_number, page_size, project_id, tenant_id
func testAccCheckAlicloudDataworksDataServiceAuthorizedApisDataSourceEmptyConfig(projectId string) string {
	return fmt.Sprintf(`
data "alicloud_dataworks_data_service_authorized_apis" "default" {
	project_id       = "%s"
	api_name_keyword = "ZZZ_NONEXISTENT_ZZZ"
	name_regex       = "ZZZ_NONEXISTENT_ZZZ"
	ids              = ["nonexistent-api-id"]
	page_number      = 1
	page_size        = 5
	tenant_id        = ""
}
`, projectId)
}
