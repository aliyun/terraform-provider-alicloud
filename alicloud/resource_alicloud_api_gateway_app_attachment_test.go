package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

// At present, One account only support create 50 apps totally.
func SkipTestAccAlicloudApigatewayAppAttachment_basic(t *testing.T) {
	var appAttachment cloudapi.AuthorizedApp
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayAppAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwaAppAttachmentBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayAppAttachmentExists("alicloud_api_gateway_app_attachment.foo", &appAttachment),
					resource.TestCheckResourceAttr("alicloud_api_gateway_app_attachment.foo", "stage_name", "PRE"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_app_attachment.foo", "api_id"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_app_attachment.foo", "group_id"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_app_attachment.foo", "app_id"),
				),
			},
		},
	})
}
func testAccCheckAlicloudApigatewayAppAttachmentExists(n string, d *cloudapi.AuthorizedApp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Api Gateway Authorization ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cloudApiService := CloudApiService{client}
		resp, err := cloudApiService.DescribeAuthorization(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error Describe Apigateway Authorization: %#v", err)
		}
		*d = *resp
		return nil
	}
}
func testAccCheckAlicloudApigatewayAppAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_api_gateway_app_attachment" {
			continue
		}
		_, err := cloudApiService.DescribeAuthorization(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe Authorization: %#v", err)
		}
	}
	return nil
}

const testAccAlicloudApigatwaAppAttachmentBasic = `
resource "alicloud_api_gateway_group" "apiGatewayGroup" {
  name        = "tf_testAccApiGroup"
  description = "tf_testAccApiGroup Description"
}
resource "alicloud_api_gateway_api" "apiGatewayApi" {
  name        = "tf_testAccApi"
  group_id    = "${alicloud_api_gateway_group.apiGatewayGroup.id}"
  description = "description"
  auth_type   = "APP"

  request_config = {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path"
    mode     = "MAPPING"
  }

  service_type = "HTTP"

  http_service_config = {
    address   = "http://apigateway-backend.alicloudapi.com:8080"
    method    = "GET"
    path      = "/web/cloudapi"
    timeout   = 22
    aone_name = "cloudapi-openapi"
  }

  request_parameters = [
    {
      name         = "aa"
      type         = "STRING"
      required     = "OPTIONAL"
      in           = "QUERY"
      in_service   = "QUERY"
      name_service = "testparams"
    },
  ]
}

resource "alicloud_api_gateway_app" "apiGatewayApp" {
  name        = "tf_testAccApiAPP"
  description = "tf_testAccApiAPP Description"
}

resource "alicloud_api_gateway_app_attachment" "foo" {
  api_id     = "${alicloud_api_gateway_api.apiGatewayApi.api_id}"
  group_id   = "${alicloud_api_gateway_group.apiGatewayGroup.id}"
  stage_name = "PRE"
  app_id     = "${alicloud_api_gateway_app.apiGatewayApp.id}"
}
 `
