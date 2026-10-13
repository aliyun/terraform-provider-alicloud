package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDataWorksDataServiceApplicationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksDataServiceApplicationDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_data_works_data_service_application.default"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDataWorksDataServiceApplicationDataSource = `
data "alicloud_data_works_data_service_application" "default" {
  project_id = "638"
}
`
