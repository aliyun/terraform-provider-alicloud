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
	resource.AddTestSweepers("alicloud_vpc", &resource.Sweeper{
		Name: "alicloud_vpc",
		F:    testSweepVpcs,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_vswitch",
			"alicloud_nat_gateway",
			"alicloud_security_group",
			"alicloud_vpn_gateway",
			"alicloud_ots_instance",
			"alicloud_router_interface",
			"alicloud_route_table",
		},
	})
}

func testSweepVpcs(region string) error {
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

	var vpcs []vpc.Vpc
	req := vpc.CreateDescribeVpcsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpcs(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving VPCs: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeVpcsResponse)
		if resp == nil || len(resp.Vpcs.Vpc) < 1 {
			break
		}
		vpcs = append(vpcs, resp.Vpcs.Vpc...)

		if len(resp.Vpcs.Vpc) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range vpcs {
		name := v.VpcName
		id := v.VpcId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VPC: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting VPC: %s (%s)", name, id)
		req := vpc.CreateDeleteVpcRequest()
		req.VpcId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpc(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPC (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudVpc_basic(t *testing.T) {
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_vpc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "router_id"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "route_table_id"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "name", "tf-testAccVpcConfig"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "description", ""),
				),
			},
		},
	})

}

func TestAccAlicloudVpc_update(t *testing.T) {
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "router_id"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "route_table_id"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "name", "tf-testAccVpcConfig"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "description", ""),
				),
			},
			{
				Config: testAccVpcConfigUpdateName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "router_id"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "route_table_id"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "name", "tf_testAccVpcConfigUpdateName"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "description", ""),
				),
			},
			{
				Config: testAccVpcConfigUpdateDesc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "router_id"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "route_table_id"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "name", "tf_testAccVpcConfigUpdateName"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "description", "hello,world"),
				),
			},
			{
				Config: testAccVpcConfigUpdateNameAndDesc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "router_id"),
					resource.TestCheckResourceAttrSet("alicloud_vpc.foo", "route_table_id"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "name", "tf_testAccVpcConfigUpdateNameAndDesc"),
					resource.TestCheckResourceAttr("alicloud_vpc.foo", "description", "who am i"),
				),
			},
		},
	})
}

func TestAccAlicloudVpc_multi(t *testing.T) {
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.bar_1", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.bar_1", "cidr_block", "172.16.0.0/12"),
					testAccCheckVpcExists("alicloud_vpc.bar_2", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.bar_2", "cidr_block", "192.168.0.0/16"),
					testAccCheckVpcExists("alicloud_vpc.bar_3", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_vpc.bar_3", "cidr_block", "10.1.0.0/21"),
				),
			},
		},
	})
}

func testAccCheckVpcExists(n string, vpc *vpc.DescribeVpcAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		instance, err := vpcService.DescribeVpc(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*vpc = instance
		return nil
	}
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpc" {
			continue
		}

		// Try to find the VPC
		instance, err := vpcService.DescribeVpc(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("VPC %s still exist", instance.VpcId))
	}

	return nil
}

const testAccVpcConfig = `
resource "alicloud_vpc" "foo" {
        name = "tf-testAccVpcConfig"
        cidr_block = "172.16.0.0/12"
}
`

const testAccVpcConfigUpdateName = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf_testAccVpcConfigUpdateName"
}
`
const testAccVpcConfigUpdateDesc = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf_testAccVpcConfigUpdateName"
	description="hello,world"
}
`
const testAccVpcConfigUpdateNameAndDesc = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf_testAccVpcConfigUpdateNameAndDesc"
	description="who am i"
}
`
const testAccVpcConfigMulti = `
variable "name" {
  	default = "tf-testAccVpcConfigMulti"
}
resource "alicloud_vpc" "bar_1" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}-1"
}
resource "alicloud_vpc" "bar_2" {
	cidr_block = "192.168.0.0/16"
	name = "${var.name}-2"
}
resource "alicloud_vpc" "bar_3" {
	cidr_block = "10.1.0.0/21"
	name = "${var.name}-3"
}
`
