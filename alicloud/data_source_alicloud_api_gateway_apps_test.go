package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

// At present, One account only support create 50 apps totally.
func SkipTestAccAlicloudApigatewayAppsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudApiGatewayAppsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_api_gateway_apps.data_apigatway_apps"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_apps.data_apigatway_apps", "apps.0.name", "tf_testAccAppDataSource"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_apps.data_apigatway_apps", "apps.0.description", "tf_testAcc api gateway description"),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_apps.data_apigatway_apps", "apps.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_apps.data_apigatway_apps", "apps.0.created_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_api_gateway_apps.data_apigatway_apps", "apps.0.modified_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudApiGatewayAppsDataSource = `

variable "apigateway_app_name_test" {
  default = "tf_testAccAppDataSource"
}

variable "apigateway_app_description_test" {
  default = "tf_testAcc api gateway description"
}

resource "alicloud_api_gateway_app" "apiAppTest" {
  name = "${var.apigateway_app_name_test}"
  description = "${var.apigateway_app_description_test}"
}

data "alicloud_api_gateway_apps" "data_apigatway_apps"{
  name_regex = "${alicloud_api_gateway_app.apiAppTest.name}"
}

`
