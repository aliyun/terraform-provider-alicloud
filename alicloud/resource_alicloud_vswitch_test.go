package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_vswitch", &resource.Sweeper{
		Name: "alicloud_vswitch",
		F:    testSweepVSwitches,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_instance",
			"alicloud_db_instance",
			"alicloud_slb",
			"alicloud_ess_scalinggroup",
			"alicloud_fc_service",
			"alicloud_cs_swarm",
			"alicloud_cs_kubernetes",
			"alicloud_kvstore_instance",
			"alicloud_route_table_attachment",
			"alicloud_havip",
			"alicloud_network_interface",
			"alicloud_drds_instance",
		},
	})
}

func testSweepVSwitches(region string) error {
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

	var vswitches []vpc.VSwitch
	req := vpc.CreateDescribeVSwitchesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitches(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving VSwitches: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeVSwitchesResponse)
		if resp == nil || len(resp.VSwitches.VSwitch) < 1 {
			break
		}
		vswitches = append(vswitches, resp.VSwitches.VSwitch...)

		if len(resp.VSwitches.VSwitch) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, vsw := range vswitches {
		name := vsw.VSwitchName
		id := vsw.VSwitchId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VSwitch: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting VSwitch: %s (%s)", name, id)
		req := vpc.CreateDeleteVSwitchRequest()
		req.VSwitchId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVSwitch(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete VSwitch (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudVswitch_basic(t *testing.T) {
	var vsw vpc.DescribeVSwitchAttributesResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_vswitch.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVswitchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVswitchConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVswitchExists("alicloud_vswitch.foo", &vsw),
					resource.TestCheckResourceAttr(
						"alicloud_vswitch.foo", "cidr_block", "172.16.0.0/21"),
				),
			},
		},
	})

}

func TestAccAlicloudVswitch_multi(t *testing.T) {
	var vsw vpc.DescribeVSwitchAttributesResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVswitchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVswitchMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVswitchExists("alicloud_vswitch.foo_0", &vsw),
					resource.TestCheckResourceAttr(
						"alicloud_vswitch.foo_0", "cidr_block", "172.16.0.0/24"),
					testAccCheckVswitchExists("alicloud_vswitch.foo_1", &vsw),
					resource.TestCheckResourceAttr(
						"alicloud_vswitch.foo_1", "cidr_block", "172.16.1.0/24"),
					testAccCheckVswitchExists("alicloud_vswitch.foo_2", &vsw),
					resource.TestCheckResourceAttr(
						"alicloud_vswitch.foo_2", "cidr_block", "172.16.2.0/24"),
				),
			},
		},
	})

}

func testAccCheckVswitchExists(n string, vsw *vpc.DescribeVSwitchAttributesResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vswitch ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		instance, err := vpcService.DescribeVswitch(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vsw = instance
		return nil
	}
}

func testAccCheckVswitchDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vswitch" {
			continue
		}

		// Try to find the Vswitch
		if _, err := vpcService.DescribeVswitch(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Vswitch still exist")
	}

	return nil
}

const testAccVswitchConfig = `
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}
variable "name" {
  default = "tf-testAccVswitchConfig"
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
`

const testAccVswitchMulti = `
variable "name" {
  default = "tf-testAccVswitchMulti"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo_0" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}-1"
}
resource "alicloud_vswitch" "foo_1" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}-2"
}
resource "alicloud_vswitch" "foo_2" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.2.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}-3"
}

`
