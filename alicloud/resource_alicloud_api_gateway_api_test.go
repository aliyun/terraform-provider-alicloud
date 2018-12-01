package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_api", &resource.Sweeper{
		Name: "alicloud_api_gateway_api",
		F:    testSweepApiGatewayApi,
	})
}

func testSweepApiGatewayApi(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	req := cloudapi.CreateDescribeApisRequest()
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApis(req)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Apis: %s", err)
	}
	apis := raw.(*cloudapi.DescribeApisResponse)

	swept := false

	for _, v := range apis.ApiSummarys.ApiSummary {
		name := v.ApiName
		id := v.ApiId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping api: %s", name)
			continue
		}
		swept = true

		log.Printf("[INFO] Deleting Api: %s", name)

		req := cloudapi.CreateDeleteApiRequest()
		req.GroupId = v.GroupId
		req.ApiId = id

		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteApi(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Api (%s): %s", name, err)
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudApigatewayApi_basic(t *testing.T) {
	var api cloudapi.DescribeApiResponse

	resource.Test(t, resource.TestCase{
		IDRefreshName: "alicloud_api_gateway_api.apiTest",
		PreCheck:      func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudApigatewayApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwayApiBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayApiExists("alicloud_api_gateway_api.apiTest", &api),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "name", "tf_testAcc_api"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "description", "tf_testAcc_api description"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "auth_type", "APP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.method", "GET"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.path", `/test/path`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.mode", "MAPPING"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "service_type", "HTTP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.address", `http://apigateway-backend.alicloudapi.com:8080`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.method", "GET"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.path", `/web/cloudapi`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.timeout", "20"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.aone_name", "cloudapi-openapi"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_parameters.1255392691.name", "testparam"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_parameters.1255392691.type", "STRING"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_parameters.1255392691.required", "OPTIONAL"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_parameters.1255392691.in", "QUERY"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_parameters.1255392691.in_service", "QUERY"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_parameters.1255392691.name_service", "testparams"),
				),
			},
		},
	})
}

func TestAccAlicloudApigatewayApi_update(t *testing.T) {
	var api cloudapi.DescribeApiResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayApiDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAlicloudApigatwayApiBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayApiExists("alicloud_api_gateway_api.apiTest", &api),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "name", "tf_testAcc_api"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "description", "tf_testAcc_api description"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "auth_type", "APP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.method", "GET"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.path", `/test/path`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.mode", "MAPPING"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "service_type", "HTTP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.address", `http://apigateway-backend.alicloudapi.com:8080`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.method", "GET"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.path", `/web/cloudapi`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.timeout", "20"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.aone_name", "cloudapi-openapi"),
				),
			},
			resource.TestStep{
				Config: testAccAlicloudApigatwayApiUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayApiExists("alicloud_api_gateway_api.apiTest", &api),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "name", "tf_testAcc_api_update"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "description", "tf_testAcc_api description update"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "auth_type", "APP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.method", "GET"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.path", `/test/path/test`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "request_config.0.mode", "MAPPING"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "service_type", "HTTP"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.address", `http://apigateway-backend.alicloudapi.com:8080`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.method", "GET"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.path", `/web/cloudapi/update`),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.timeout", "20"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_api.apiTest", "http_service_config.0.aone_name", "cloudapi-openapi"),
				),
			},
		},
	})
}

func testAccCheckAlicloudApigatewayApiExists(n string, d *cloudapi.DescribeApiResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Api ID is set")
		}

		fmt.Println(rs.Primary.ID)

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cloudApiService := CloudApiService{client}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		resp, err := cloudApiService.DescribeApi(split[1], split[0])
		if err != nil {

			return fmt.Errorf("Error Describe Api: %#v", err)
		}

		*d = *resp
		return nil
	}
}

func testAccCheckAlicloudApigatewayApiDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_api_gateway_api" {
			continue
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		_, err := cloudApiService.DescribeApi(split[1], split[0])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe Api: %#v", err)
		}
	}

	return nil
}

const testAccAlicloudApigatwayApiBasic = `

variable "apigateway_group_name_test" {
  default = "tf_testAccApiGroupDataSource"
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
      protocol        = "HTTP"
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

`

const testAccAlicloudApigatwayApiUpdate = `

variable "apigateway_group_name_test" {
  default = "tf_testAccApiGroupDataSource"
}

variable "apigateway_group_description_test" {
  default = "tf_testAcc_api group description"
}

resource "alicloud_api_gateway_group" "apiGroupTest" {
  name = "${var.apigateway_group_name_test}"
  description = "${var.apigateway_group_description_test}"
}

resource "alicloud_api_gateway_api" "apiTest" {
  name = "tf_testAcc_api_update"
  group_id = "${alicloud_api_gateway_group.apiGroupTest.id}"
  description = "tf_testAcc_api description update"
  auth_type = "APP"
  request_config = [
    {
      protocol = "HTTP"
      method = "GET"
      path = "/test/path/test"
      mode = "MAPPING"
    },
  ]
  service_type = "HTTP"
  http_service_config = [
    {
      address = "http://apigateway-backend.alicloudapi.com:8080"
      method = "GET"
      path = "/web/cloudapi/update"
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

`
