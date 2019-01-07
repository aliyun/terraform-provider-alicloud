package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_app", &resource.Sweeper{
		Name: "alicloud_api_gateway_app",
		F:    testSweepApiGatewayApp,
	})
}

func testSweepApiGatewayApp(region string) error {
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

	req := cloudapi.CreateDescribeAppAttributesRequest()
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeAppAttributes(req)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Apps: %s", err)
	}
	apps, _ := raw.(*cloudapi.DescribeAppAttributesResponse)

	swept := false

	for _, v := range apps.Apps.AppAttribute {
		name := v.AppName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping app: %s", name)
			continue
		}
		swept = true

		log.Printf("[INFO] Deleting App: %s", name)

		req := cloudapi.CreateDeleteAppRequest()
		req.AppId = requests.Integer(v.AppId)
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteApp(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete App (%s): %s", name, err)
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func SkipTestAccAlicloudApigatewayApp_basic(t *testing.T) {
	var app cloudapi.DescribeAppResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwayAppBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayAppExists("alicloud_api_gateway_app.appTest", &app),
					resource.TestCheckResourceAttr("alicloud_api_gateway_app.appTest", "name", "tf_testAccAppResource"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_app.appTest", "description", "tf_testAcc api gateway description"),
				),
			},
		},
	})
}

// At present, One account only support create 50 apps totally.
func SkipTestAccAlicloudApigatewayApp_update(t *testing.T) {
	var app cloudapi.DescribeAppResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwayAppBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayAppExists("alicloud_api_gateway_app.appTest", &app),
					resource.TestCheckResourceAttr("alicloud_api_gateway_app.appTest", "name", "tf_testAccAppResource"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_app.appTest", "description", "tf_testAcc api gateway description"),
				),
			},
			{
				Config: testAccAlicloudApigatwayAppUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayAppExists("alicloud_api_gateway_app.appTest", &app),
					resource.TestCheckResourceAttr("alicloud_api_gateway_app.appTest", "name", "tf_testAccAppResource_u"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_app.appTest", "description", "tf_testAcc api gateway description update"),
				),
			},
		},
	})
}

func testAccCheckAlicloudApigatewayAppExists(n string, d *cloudapi.DescribeAppResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cloudApiService := CloudApiService{client}

		resp, err := cloudApiService.DescribeApp(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error Describe App: %#v", err)
		}

		*d = *resp
		return nil
	}
}

func testAccCheckAlicloudApigatewayAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_api_gateway_app" {
			continue
		}

		_, err := cloudApiService.DescribeApp(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe App: %#v", err)
		}
	}

	return nil
}

const testAccAlicloudApigatwayAppBasic = `

variable "apigateway_app_name_test" {
  default = "tf_testAccAppResource"
}

variable "apigateway_app_description_test" {
  default = "tf_testAcc api gateway description"
}

resource "alicloud_api_gateway_app" "appTest" {
  name = "${var.apigateway_app_name_test}"
  description = "${var.apigateway_app_description_test}"
}
`

const testAccAlicloudApigatwayAppUpdate = `

variable "apigateway_app_name_test" {
  default = "tf_testAccAppResource_u"
}

variable "apigateway_app_description_test" {
  default = "tf_testAcc api gateway description update"
}

resource "alicloud_api_gateway_app" "appTest" {
  name = "${var.apigateway_app_name_test}"
  description = "${var.apigateway_app_description_test}"
}
`
