package alicloud

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudApigatewayGroupsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudApiGatewayGroupDataSource(acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_api_gateway_groups.data_apigatway_groups"),
					resource.TestMatchResourceAttr("data.alicloud_api_gateway_groups.data_apigatway_groups", "groups.0.name", regexp.MustCompile("^tf_testAccGroupDataSource_*")),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_groups.data_apigatway_groups", "groups.0.description", "tf_testAcc api gateway description"),
				),
			},
		},
	})
}

func TestAccAlicloudApigatewayGroupsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudApiGatewayGroupDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_api_gateway_groups.data_apigatway_groups"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_groups.data_apigatway_groups", "groups.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_api_gateway_groups.data_apigatway_groups", "groups.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_api_gateway_groups.data_apigatway_groups", "groups.0.description"),
				),
			},
		},
	})
}

func testAccCheckAlicloudApiGatewayGroupDataSource(rand int) string {
	return fmt.Sprintf(`

	variable "apigateway_group_name_test" {
	  default = "tf_testAccGroupDataSource_%d"
	}

	variable "apigateway_group_description_test" {
	  default = "tf_testAcc api gateway description"
	}

	resource "alicloud_api_gateway_group" "apiGroupTest" {
	  name = "${var.apigateway_group_name_test}"
	  description = "${var.apigateway_group_description_test}"
	}

	data "alicloud_api_gateway_groups" "data_apigatway_groups"{
	  name_regex = "${alicloud_api_gateway_group.apiGroupTest.name}"
	}
	`, rand)
}

const testAccCheckAlicloudApiGatewayGroupDataSourceEmpty = `
data "alicloud_api_gateway_groups" "data_apigatway_groups"{
  name_regex = "^tf-testacc-fake-name"
}
`
