package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudApigatewayGroup_basic(t *testing.T) {
	var group cloudapi.DescribeApiGroupResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudApigatwayGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayGroupExists("alicloud_apigateway_api_group.apiGroup", &group),
					resource.TestCheckResourceAttr("alicloud_apigateway_api_group.apiGroup", "name", "apigateway—testgroup"),
				),
			},
		},
	})
}

func testAccCheckAlicloudApigatewayGroupExists(n string, d *cloudapi.DescribeApiGroupResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Apigroup ID is set")
		}

		fmt.Println(rs.Primary.ID)

		resp, err := testAccProvider.Meta().(*AliyunClient).DescribeApiGroup(rs.Primary.ID)
		if err != nil {

			return fmt.Errorf("Error Describe Apigroup: %#v", err)
		}

		if resp == nil {
			return fmt.Errorf("Api group is not found in the instance %s.", rs.Primary.ID)
		}

		*d = *resp
		return nil
	}
}

func testAccCheckAlicloudApigatewayGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_apigateway_api_group" {
			continue
		}

		_, err := client.DescribeApiGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe Apigroup: %#v", err)
		}
	}

	return nil
}

const testAlicloudApigatwayGroupBasic = `

variable "apigateway_group_name_test" {
  default = "apigateway—testgroup"
}

variable "apigateway_group_description_test" {
  default = "api gateway description"
}

resource "alicloud_apigateway_api_group" "apiGroup1" {
  name = "${var.apigateway_group_name_test}"
  description = "${var.apigateway_group_description_test}"
}
`
