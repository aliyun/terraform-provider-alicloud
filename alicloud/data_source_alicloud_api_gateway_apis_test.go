package alicloud

import (
	"testing"

	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudApigatewayApisDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudApiGatewayApiDataSource(acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_api_gateway_apis.data_apis"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_apis.data_apis", "apis.0.name", "tf_testAcc_api"),
					resource.TestMatchResourceAttr("data.alicloud_api_gateway_apis.data_apis", "apis.0.group_name", regexp.MustCompile("^tf_testAccApiGroupDataSource_*")),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_apis.data_apis", "apis.0.description", "tf_testAcc_api description"),
				),
			},
		},
	})
}

func TestAccAlicloudApigatewayApisDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudApiGatewayApiDataSourceEmpty(acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_api_gateway_apis.data_apis"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_apis.data_apis", "apis.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_api_gateway_apis.data_apis", "apis.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_api_gateway_apis.data_apis", "apis.0.group_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_api_gateway_apis.data_apis", "apis.0.description"),
				),
			},
		},
	})
}

func testAccCheckAlicloudApiGatewayApiDataSource(rand int) string {
	return fmt.Sprintf(`

	variable "apigateway_group_name_test" {
	  default = "tf_testAccApiGroupDataSource_%d"
	}

	variable "apigateway_group_description_test" {
	  default = "tf_testAcc_api group description"
	}

	resource "alicloud_api_gateway_group" "apiGroupTest" {
	  name = "${var.apigateway_group_name_test}"
	  description = "${var.apigateway_group_description_test}"
	}

	resource "alicloud_api_gateway_api" "apiTest" {
	  name = "tf_testAcc_api"
	  group_id = "${alicloud_api_gateway_group.apiGroupTest.id}"
	  description = "tf_testAcc_api description"
	  auth_type = "APP"
	  request_config = [
	    {
	      protocol = "HTTP"
	      method = "GET"
	      path = "/test/path"
	      mode = "MAPPING"
	    },
	  ]
	  service_type = "HTTP"
	  http_service_config = [
	    {
	      address = "http://apigateway-backend.alicloudapi.com:8080"
	      method = "GET"
	      path = "/web/cloudapi"
	      timeout = 20
	      aone_name = "cloudapi-openapi"
	    },
	  ]

	  request_parameters = [
	    {
	      name = "testparam"
	      type = "STRING"
	      required = "OPTIONAL"
	      in = "QUERY"
	      in_service = "QUERY"
	      name_service = "testparams"
	    },
	  ]
	}

	data "alicloud_api_gateway_apis" "data_apis"{
	  group_id = "${alicloud_api_gateway_group.apiGroupTest.id}"
	  api_id = "${alicloud_api_gateway_api.apiTest.id}"
	}

	`, rand)
}

func testAccCheckAlicloudApiGatewayApiDataSourceEmpty(rand int) string {
	return fmt.Sprintf(`
	variable "apigateway_group_name_test" {
	  default = "tf_testAccApiGroupDataSourceEmpty_%d"
	}

	variable "apigateway_group_description_test" {
	  default = "tf_testAcc_api group description"
	}

	resource "alicloud_api_gateway_group" "apiGroupTest" {
	  name = "${var.apigateway_group_name_test}"
	  description = "${var.apigateway_group_description_test}"
	}


	data "alicloud_api_gateway_apis" "data_apis"{
	  group_id = "${alicloud_api_gateway_group.apiGroupTest.id}"
	}
	`, rand)
}
