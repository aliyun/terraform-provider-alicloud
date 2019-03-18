package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_slb", &resource.Sweeper{
		Name: "alicloud_slb",
		F:    testSweepSLBs,
	})
}

func testSweepSLBs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var slbs []slb.LoadBalancer
	req := slb.CreateDescribeLoadBalancersRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeLoadBalancers(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving SLBs: %s", err)
		}
		resp, _ := raw.(*slb.DescribeLoadBalancersResponse)
		if resp == nil || len(resp.LoadBalancers.LoadBalancer) < 1 {
			break
		}
		slbs = append(slbs, resp.LoadBalancers.LoadBalancer...)

		if len(resp.LoadBalancers.LoadBalancer) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, loadBalancer := range slbs {
		name := loadBalancer.LoadBalancerName
		id := loadBalancer.LoadBalancerId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping SLB: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting SLB: %s (%s)", name, id)
		req := slb.CreateDeleteLoadBalancerRequest()
		req.LoadBalancerId = id
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteLoadBalancer(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete SLB (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSlb_paybybandwidth(t *testing.T) {
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: "alicloud_slb.bandwidth",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbPayByBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.bandwidth", &slb),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "name", "tf-testAccSlbPayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "bandwidth", "1"),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "specification", "slb.s2.medium"),
					resource.TestCheckResourceAttrSet("alicloud_slb.bandwidth", "address"),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "tags.%", "10"),
				),
			},
			{
				Config: testAccSlbPayByBandwidthUpName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.bandwidth", &slb),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "name", "tf-testAccSlbPayByBandwidthUpName"),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "bandwidth", "1"),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "specification", "slb.s2.medium"),
					resource.TestCheckResourceAttrSet("alicloud_slb.bandwidth", "address"),
					resource.TestCheckResourceAttr("alicloud_slb.bandwidth", "tags.%", "10"),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_vpc(t *testing.T) {
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: "alicloud_slb.vpc",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlb4Vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.vpc", &slb),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "name", "tf-testAccSlb4Vpc"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttrSet("alicloud_slb.vpc", "bandwidth"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "specification", "slb.s2.small"),
					resource.TestCheckResourceAttrSet("alicloud_slb.vpc", "address"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "tags.%", "10"),
				),
			},
			{
				Config: testAccSlb4VpcUpInternet_charge_type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.vpc", &slb),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "name", "tf-testAccSlb4Vpc"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "bandwidth", "5"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "specification", "slb.s2.small"),
					resource.TestCheckResourceAttrSet("alicloud_slb.vpc", "address"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "tags.%", "10"),
				),
			},
			{
				Config: testAccSlb4VpcUpBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.vpc", &slb),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "name", "tf-testAccSlb4Vpc"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "specification", "slb.s2.small"),
					resource.TestCheckResourceAttrSet("alicloud_slb.vpc", "address"),
					resource.TestCheckResourceAttr("alicloud_slb.vpc", "tags.%", "10"),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_spec(t *testing.T) {
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.SlbGuaranteedSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_slb.spec",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbBandSpec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.spec", &slb),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "name", "tf_testAccSlbBandSpec"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "specification", "slb.s2.small"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttrSet("alicloud_slb.spec", "bandwidth"),
					resource.TestCheckResourceAttrSet("alicloud_slb.spec", "address"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.%", "10"),

					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_a", "1"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_b", "2"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_c", "3"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_d", "4"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_e", "5"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_f", "6"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_g", "7"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_h", "8"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_i", "9"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_j", "10"),
				),
			},
			{
				Config: testAccSlbBandSpecUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.spec", &slb),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "specification", "slb.s2.medium"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.%", "8"),
					// tags no changed
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_a", "1"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_b", "2"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_c", "3"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_e", "5"),
					// tags remove
					resource.TestCheckNoResourceAttr("alicloud_slb.spec", "tags.tag_d"),
					resource.TestCheckNoResourceAttr("alicloud_slb.spec", "tags.tag_i"),
					resource.TestCheckNoResourceAttr("alicloud_slb.spec", "tags.tag_j"),
					// tags update
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_f", "66"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_g", "77"),
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_h", "88"),
					// tags add new
					resource.TestCheckResourceAttr("alicloud_slb.spec", "tags.tag_k", "11"),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_pay_type(t *testing.T) {
	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: "alicloud_slb.pay_type",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbPayType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.pay_type", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.pay_type", "name", "tf-testAccSlbPayType"),
					resource.TestCheckResourceAttr(
						"alicloud_slb.pay_type", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_slb.pay_type", "instance_charge_type", "PostPaid"),
				),
			},
		},
	})
}

func testAccCheckSlbExists(n string, slb *slb.DescribeLoadBalancerAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No SLB ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		slbService := SlbService{client}
		instance, err := slbService.DescribeLoadBalancerAttribute(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*slb = *instance
		return nil
	}
}

func testAccCheckSlbDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	slbService := SlbService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb" {
			continue
		}

		// Try to find the Slb
		if _, err := slbService.DescribeLoadBalancerAttribute(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("SLB still exist"))
	}

	return nil
}

const testAccSlbPayByBandwidth = `
resource "alicloud_slb" "bandwidth" {
  name = "tf-testAccSlbPayByBandwidth"
  specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}
`
const testAccSlbPayByBandwidthUpName = `
resource "alicloud_slb" "bandwidth" {
  name = "tf-testAccSlbPayByBandwidthUpName"
  specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}
`

const testAccSlb4Vpc = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_slb" "vpc" {
  name = "${var.name}"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}
`
const testAccSlb4VpcUpInternet_charge_type = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_slb" "vpc" {
  name = "${var.name}"
  specification = "slb.s2.small"
  internet_charge_type = "PayByBandwidth"
  internet = true
  bandwidth = "5"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}
`
const testAccSlb4VpcUpBandwidth = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_slb" "vpc" {
  name = "${var.name}"
  specification = "slb.s2.small"
  internet_charge_type = "PayByBandwidth"
  internet = true
  bandwidth = "10"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}
`
const testAccSlbBandSpec = `
variable "name" {
  default = "tf_testAccSlbBandSpec"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_slb" "spec" {
  name = "${var.name}"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}
`

const testAccSlbBandSpecUpdate = `
variable "name" {
  default = "tf_testAccSlbBandSpecUpdate"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_slb" "spec" {
  name = "${var.name}"
  specification = "slb.s2.medium"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_e = 5
    tag_f = 66
    tag_g = 77
    tag_h = 88
    tag_k = 11
  }
}
`
const testAccSlbPayType = `
resource "alicloud_slb" "pay_type" {
  name = "tf-testAccSlbPayType"
  specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
  instance_charge_type = "PostPaid"
  period = 2

}
`
