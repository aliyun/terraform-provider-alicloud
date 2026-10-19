package dataworks_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAlicloudDataWorksExtensionsDataSource_basic lists DataWorks extensions
// and verifies the name_regex filter produces an empty set when nothing matches.
func TestAccAlicloudDataWorksExtensionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudDataWorksExtensionsDataSourceConfig_exist(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_extensions.default", "extensions.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_data_works_extensions.default", "ids.#"),
				),
			},
			{
				Config: testAccAlicloudDataWorksExtensionsDataSourceConfig_fake(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_data_works_extensions.default", "extensions.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_data_works_extensions.default", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccAlicloudDataWorksExtensionsDataSourceConfig_exist() string {
	return `
data "alicloud_data_works_extensions" "default" {
  output_file = "data_works_extensions.json"
}
`
}

func testAccAlicloudDataWorksExtensionsDataSourceConfig_fake() string {
	return `
data "alicloud_data_works_extensions" "default" {
  name_regex = "___non_existent_extension_regex___12345"
}
`
}
