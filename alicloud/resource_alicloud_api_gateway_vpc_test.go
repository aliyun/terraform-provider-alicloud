package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_vpc_access", &resource.Sweeper{
		Name: "alicloud_api_gateway_vpc_access",
		F:    testSweepApiGatewayVpc,
	})
}

func testSweepApiGatewayVpc(region string) error {
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

	req := cloudapi.CreateDescribeVpcAccessesRequest()
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeVpcAccesses(req)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Api Gateway Vpc: %s", err)
	}

	allVpcs, _ := raw.(*cloudapi.DescribeVpcAccessesResponse)

	swept := false

	for _, v := range allVpcs.VpcAccessAttributes.VpcAccessAttribute {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Api Gateway Vpc: %s", name)
			continue
		}
		swept = true

		req := cloudapi.CreateRemoveVpcAccessRequest()
		req.VpcId = v.VpcId
		req.InstanceId = v.InstanceId
		req.Port = requests.NewInteger(v.Port)
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.RemoveVpcAccess(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Api Gaiteway Vpc (%s): %s", name, err)
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudApigatewayVpc_basic(t *testing.T) {
	var vpc cloudapi.VpcAccessAttribute
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.ApiGatewayNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudApigatewayVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudApigatwaVpcBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudApigatewayVpcExists("alicloud_api_gateway_vpc_access.foo", &vpc),
					resource.TestCheckResourceAttr("alicloud_api_gateway_vpc_access.foo", "name", fmt.Sprintf("tf-testAccApiGatewayVpc-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_api_gateway_vpc_access.foo", "port", "8080"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_vpc_access.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_vpc_access.foo", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckAlicloudApigatewayVpcExists(n string, d *cloudapi.VpcAccessAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Api Gateway Vpc ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cloudApiService := CloudApiService{client}

		resp, err := cloudApiService.DescribeVpcAccess(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error Describe Apigateway Vpc: %#v", err)
		}

		*d = *resp
		return nil
	}
}

func testAccCheckAlicloudApigatewayVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_api_gateway_vpc_access" {
			continue
		}

		_, err := cloudApiService.DescribeVpcAccess(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe Vpc: %#v", err)
		}
	}

	return nil
}

func testAccAlicloudApigatwaVpcBasic(rand int) string {
	return fmt.Sprintf(`

	data "alicloud_zones" "default" {
	  "available_disk_category"= "cloud_efficiency"
	  "available_resource_creation"= "VSwitch"
	}

	data "alicloud_instance_types" "default" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  cpu_core_count = 1
	  memory_size = 2
	}

	data "alicloud_images" "default" {
	  name_regex = "^ubuntu_14.*_64"
	  most_recent = true
	  owners = "system"
	}

	variable "name" {
	  default = "tf-testAccApiGatewayVpc-%d"
	}

	resource "alicloud_vpc" "foo" {
	  name = "${var.name}"
	  cidr_block = "172.16.0.0/12"
	}

	resource "alicloud_vswitch" "foo" {
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "172.16.0.0/21"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  name = "${var.name}"
	}

	resource "alicloud_security_group" "tf_test_foo" {
	  name = "${var.name}"
	  description = "foo"
	  vpc_id = "${alicloud_vpc.foo.id}"
	}

	resource "alicloud_instance" "foo" {
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	  image_id = "${data.alicloud_images.default.images.0.id}"

	  # series III
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  system_disk_category = "cloud_efficiency"

	  internet_charge_type = "PayByTraffic"
	  internet_max_bandwidth_out = 5
	  security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	  instance_name = "${var.name}"
	}

	resource "alicloud_api_gateway_vpc_access" "foo" {
	  name        = "${var.name}"
	  vpc_id      = "${alicloud_vpc.foo.id}"
	  instance_id = "${alicloud_instance.foo.id}"
	  port        = 8080
	}

	`, rand)
}
