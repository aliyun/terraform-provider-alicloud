package apigateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAlicloudApiGatewayDatasetsDataSource covers the datasets data source:
// it lists the dataset created in the same config, filtered by id, and asserts
// that name_regex narrows the result to one (match) then zero (no match).
func TestAccAlicloudApiGatewayDatasetsDataSource(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayDatasetsDataSourceConfig(r, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_datasets.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_datasets.default", "datasets.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_datasets.default", "datasets.0.dataset_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_datasets.default", "datasets.0.dataset_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_datasets.default", "datasets.0.dataset_type"),
				),
			},
			{
				Config: testAccApiGatewayDatasetsDataSourceConfig(r, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_datasets.default", "datasets.#", "0"),
				),
			},
		},
	})
}

func testAccApiGatewayDatasetsDataSourceConfig(r int, match bool) string {
	regex := fmt.Sprintf("tf-testacc-ds-%d", r)
	if !match {
		regex = regex + "-no-such-dataset"
	}
	return fmt.Sprintf(`
resource "alicloud_api_gateway_dataset" "default" {
  dataset_name = "tf-testacc-ds-%[1]d"
  dataset_type = "JWT_BLOCKING"
  description  = "tf test acc dataset for datasource"
}

data "alicloud_api_gateway_datasets" "default" {
  ids        = ["${alicloud_api_gateway_dataset.default.id}"]
  name_regex = "%[2]s"
}
`, r, regex)
}
