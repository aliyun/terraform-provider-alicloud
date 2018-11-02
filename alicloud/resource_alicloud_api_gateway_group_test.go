package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_group", &resource.Sweeper{
		Name: "alicloud_api_gateway_group",
		F:    testSweepApiGatewayGroup,
		Dependencies: []string{
			"alicloud_api_gateway_api",
		},
	})
}

func testSweepApiGatewayGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	req := cloudapi.CreateDescribeApiGroupsRequest()
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApiGroups(req)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Api Groups: %s", err)
	}
	apiGroups, _ := raw.(*cloudapi.DescribeApiGroupsResponse)

	sweeped := false

	for _, v := range apiGroups.ApiGroupAttributes.ApiGroupAttribute {
		name := v.GroupName
		id := v.GroupId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping api group: %s", name)
			continue
		}
		sweeped = true

		log.Printf("[INFO] Deleting Api Group: %s", name)

		req := cloudapi.CreateDeleteApiGroupRequest()
		req.GroupId = id
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteApiGroup(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Api Group (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudApigatewayGroup_basic(t *testing.T) {
	var group cloudapi.DescribeApiGroupResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwayGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayGroupExists("alicloud_api_gateway_group.apiGroupTest", &group),
					resource.TestCheckResourceAttr("alicloud_api_gateway_group.apiGroupTest", "name", "tf_testAccGroupResource"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_group.apiGroupTest", "description", "tf_testAcc api gateway description"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cloudApiService := CloudApiService{client}

		resp, err := cloudApiService.DescribeApiGroup(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error Describe Apigroup: %#v", err)
		}

		*d = *resp
		return nil
	}
}

func testAccCheckAlicloudApigatewayGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_api_gateway_group" {
			continue
		}

		_, err := cloudApiService.DescribeApiGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe Apigroup: %#v", err)
		}
	}

	return nil
}

const testAccAlicloudApigatwayGroupBasic = `

variable "apigateway_group_name_test" {
  default = "tf_testAccGroupResource"
}

variable "apigateway_group_description_test" {
  default = "tf_testAcc api gateway description"
}

resource "alicloud_api_gateway_group" "apiGroupTest" {
  name = "${var.apigateway_group_name_test}"
  description = "${var.apigateway_group_description_test}"
}
`
