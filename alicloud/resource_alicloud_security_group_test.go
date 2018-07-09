package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudSecurityGroup_basic(t *testing.T) {
	var sg ecs.DescribeSecurityGroupAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(
						"alicloud_security_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_security_group.foo",
						"name",
						"sg_test"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroup_withVpc(t *testing.T) {
	var sg ecs.DescribeSecurityGroupAttributeResponse
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupConfig_withVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(
						"alicloud_security_group.foo", &sg),
					testAccCheckVpcExists(
						"alicloud_vpc.vpc", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_security_group.foo",
						"inner_access",
						"true"),
				),
			},
		},
	})

}

func testAccCheckSecurityGroupExists(n string, sg *ecs.DescribeSecurityGroupAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SecurityGroup ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		d, err := client.DescribeSecurityGroupAttribute(rs.Primary.ID)

		log.Printf("[WARN] security group id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}
		if d.SecurityGroupId == rs.Primary.ID {
			*sg = d
		}
		return nil
	}
}

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group" {
			continue
		}

		group, err := client.DescribeSecurityGroupAttribute(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidSecurityGroupIdNotFound}) {
				continue
			}
			return err
		}

		if group.SecurityGroupId != "" {
			return fmt.Errorf("Error SecurityGroup still exist")
		}
	}
	return nil
}

const testAccSecurityGroupConfig = `
resource "alicloud_security_group" "foo" {
  name = "sg_test"
}
`

const testAccSecurityGroupConfig_withVpc = `
resource "alicloud_security_group" "foo" {
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}
`
