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
)

func init() {
	resource.AddTestSweepers("alicloud_slb", &resource.Sweeper{
		Name: "alicloud_slb",
		F:    testSweepSLBs,
	})
}

func testSweepSLBs(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var slbs []slb.LoadBalancer
	req := slb.CreateDescribeLoadBalancersRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.slbconn.DescribeLoadBalancers(req)
		if err != nil {
			return fmt.Errorf("Error retrieving SLBs: %s", err)
		}
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
		if _, err := conn.slbconn.DeleteLoadBalancer(req); err != nil {
			log.Printf("[ERROR] Failed to delete SLB (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

//test internet_charge_type is PayByBandwidth and it only support China mainland region
func TestAccAlicloudSlb_paybybandwidth(t *testing.T) {
	if !isRegionSupports(SlbPayByBandwidth) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), SlbPayByBandwidth)
		return
	}

	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.bandwidth",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbPayByBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.bandwidth", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.bandwidth", "name", "tf-testAccSlbPayByBandwidth"),
					resource.TestCheckResourceAttr(
						"alicloud_slb.bandwidth", "internet_charge_type", "PayByBandwidth"),
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
		},

		// module name
		IDRefreshName: "alicloud_slb.vpc",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlb4Vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.vpc", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.vpc", "name", "tf-testAccSlb4Vpc"),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_spec(t *testing.T) {
	if !isRegionSupports(SlbSpecification) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), SlbSpecification)
		return
	}

	var slb slb.DescribeLoadBalancerAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.spec",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbBandSpec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.spec", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.spec", "specification", "slb.s2.small"),
				),
			},
			resource.TestStep{
				Config: testAccSlbBandSpecUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.spec", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.spec", "specification", "slb.s2.medium"),
				),
			},
		},
	})
}

func testAccCheckSlbExists(n string, slb *slb.DescribeLoadBalancerAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SLB ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeLoadBalancerAttribute(rs.Primary.ID)

		if err != nil {
			return err
		}

		*slb = *instance
		return nil
	}
}

func testAccCheckSlbDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb" {
			continue
		}

		// Try to find the Slb
		if _, err := client.DescribeLoadBalancerAttribute(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
		}
		return fmt.Errorf("SLB still exist")
	}

	return nil
}

const testAccSlbPayByBandwidth = `
resource "alicloud_slb" "bandwidth" {
  name = "tf-testAccSlbPayByBandwidth"
  specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
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
}
`
