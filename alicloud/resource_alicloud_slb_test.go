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
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_cs_cluster",
		},
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

	service := SlbService{client}
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
		if err := service.sweepSlb(id); err != nil {
			log.Printf("[ERROR] Failed to delete SLB (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSlb_classictest(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlb_no_specification,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 "tf-testAccslbbasic_name",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccSlb_clissic_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 "tf-testAccslbbasic_name",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccSlb_clissic_specification,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "slb.s2.medium",
					}),
				),
			},
			{
				Config: testAccSlb_clissic_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccslbbasic_namenew",
					}),
				),
			},
			{
				Config: testAccSlb_clissic_tags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "10",
					}),
				),
			},
			{
				Config: testAccSlb_clissic_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 "tf-testAccslbbasic_name",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
						"tags.%":               REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_vpctest(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbVpc_no_specification,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 "tf-testAccSlb4Vpc",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccSlbVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 "tf-testAccSlb4Vpc",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccSlbBandSpecUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "slb.s2.medium",
					}),
				),
			},
			{
				Config: testAccSlbVpcUpInternet_charge_type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type": "PayByBandwidth",
					}),
				),
			},
			{
				Config: testAccSlbVpcUpBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccSlbtagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "5",
					}),
				),
			},
			{
				Config: testAccSlbnamesUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlb4Vpcnew",
					}),
				),
			},
			{
				Config: testAccSlbVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 "tf-testAccSlb4Vpc",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"tags.%":               REMOVEKEY,
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_vpcmulti(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbVpc_multi,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 "tf-testAccSlb4Vpc",
						"internet_charge_type": "PayByTraffic",
						"bandwidth":            CHECKSET,
						"specification":        "slb.s2.small",
						"tags.%":               "10",
						"address":              CHECKSET,
						"master_zone_id":       CHECKSET,
						"slave_zone_id":        CHECKSET,
					}),
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
		instance, err := slbService.DescribeSlb(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		slb = instance
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
		if _, err := slbService.DescribeSlb(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("SLB still exist"))
	}

	return nil
}

const testAccSlb_clissic_basic = `
variable "name" {
  default = "tf-testAccslbbasic_name"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  specification = "slb.s2.small"
	internet = true
}
`

const testAccSlb_clissic_specification = `
variable "name" {
  default = "tf-testAccslbbasic_name"
}
resource "alicloud_slb" "default" {
  	name = "${var.name}"
  	specification = "slb.s2.medium"
	internet = true
}
`

const testAccSlb_clissic_name = `
variable "name" {
  default = "tf-testAccslbbasic_namenew"
}
resource "alicloud_slb" "default" {
  name = "${var.name}"
  specification = "slb.s2.medium"
internet = true
}
`

const testAccSlb_clissic_internet_charge_type = `
variable "name" {
  default = "tf-testAccslbbasic_namenew"
}
resource "alicloud_slb" "default" {
  name = "${var.name}"
  specification = "slb.s2.medium"
internet_charge_type = "PayByBandwidth"
internet = true
}
`

const testAccSlb_clissic_bandwidth = `
variable "name" {
  default = "tf-testAccslbbasic_namenew"
}
resource "alicloud_slb" "default" {
  	name = "${var.name}"
  	specification = "slb.s2.medium"
	internet_charge_type = "PayByBandwidth"
	bandwidth="10"
internet = true
}
`
const testAccSlb_clissic_tags = `
variable "name" {
  default = "tf-testAccslbbasic_namenew"
}
resource "alicloud_slb" "default" {
  	name = "${var.name}"
  	specification = "slb.s2.medium"
	internet_charge_type = "PayByBandwidth"
	bandwidth="10"
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

const testAccSlbVpc = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
`

const testAccSlbBandSpecUpdate = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_slb" "default" {
  	name = "${var.name}"
  	specification = "slb.s2.medium"
  	vswitch_id = "${alicloud_vswitch.default.id}"
}
`

const testAccSlbVpcUpInternet_charge_type = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
  vswitch_id = "${alicloud_vswitch.default.id}"
}
`
const testAccSlbVpcUpBandwidth = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
  bandwidth = "10"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
`

const testAccSlbtagsUpdate = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
  bandwidth = "10"
  vswitch_id = "${alicloud_vswitch.default.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
  }
}
`

const testAccSlbnamesUpdate = `
variable "name" {
  default = "tf-testAccSlb4Vpcnew"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
specification = "slb.s2.medium"
  internet_charge_type = "PayByBandwidth"
  internet = true
  bandwidth = "10"
  vswitch_id = "${alicloud_vswitch.default.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
  }
}
`

const testAccSlbVpc_multi = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "default" {
	count = 10
  name = "${var.name}"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.default.id}"
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
const testAccSlb_no_specification = `
variable "name" {
  default = "tf-testAccslbbasic_name"
}
resource "alicloud_slb" "default" {
  	name = "${var.name}"
	internet = true
}
`
const testAccSlbVpc_no_specification = `
variable "name" {
  default = "tf-testAccSlb4Vpc"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
`
