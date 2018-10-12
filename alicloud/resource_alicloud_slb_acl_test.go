package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_slb_acl", &resource.Sweeper{
		Name: "alicloud_slb_acl",
		F:    testSweepSlbAcl,
	})
}

func testSweepSlbAcl(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	req := slb.CreateDescribeAccessControlListsRequest()
	req.RegionId = aliyunClient.RegionId
	resp, err := aliyunClient.slbconn.DescribeAccessControlLists(req)
	if err != nil {
		return err
	}

	for _, acl := range resp.Acls.Acl {
		name := acl.AclName
		id := acl.AclId

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Slb Acl: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting Slb Acl : %s (%s)", name, id)
		req := slb.CreateDeleteAccessControlListRequest()
		req.AclId = id
		if _, err := aliyunClient.slbconn.DeleteAccessControlList(req); err != nil {
			log.Printf("[ERROR] Failed to delete Slb Acl (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSlbAcl_basic(t *testing.T) {
	var acl slb.DescribeAccessControlListAttributeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_acl.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbAclDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbAclBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbAclExists("alicloud_slb_acl.foo", &acl),
					resource.TestCheckResourceAttr(
						"alicloud_slb_acl.foo", "name", "tf-testAccSlbAcl"),

					resource.TestCheckResourceAttr(
						"alicloud_slb_acl.foo", "ip_version", "ipv4"),

					resource.TestCheckResourceAttr(
						"alicloud_slb_acl.foo", "entry_list.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccSlbAclBasicConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbAclExists("alicloud_slb_acl.foo", &acl),
					resource.TestCheckResourceAttr(
						"alicloud_slb_acl.foo", "name", "tf-testAccSlbAclUpdate"),

					resource.TestCheckResourceAttr(
						"alicloud_slb_acl.foo", "ip_version", "ipv4"),

					resource.TestCheckResourceAttr(
						"alicloud_slb_acl.foo", "entry_list.#", "3"),
				),
			},
		},
	})
}

func testAccCheckSlbAclExists(n string, acl *slb.DescribeAccessControlListAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SLB ACL ID is set")
		}

		req := slb.CreateDescribeAccessControlListAttributeRequest()
		req.AclId = rs.Primary.ID

		client := testAccProvider.Meta().(*AliyunClient).slbconn
		r, err := client.DescribeAccessControlListAttribute(req)
		if err != nil {
			return fmt.Errorf("No SLB ACL ID %s is set", req.AclId)
		}

		*acl = *r

		return nil
	}
}

func testAccCheckSlbAclDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient).slbconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_acl" {
			continue
		}

		req := slb.CreateDescribeAccessControlListAttributeRequest()
		req.AclId = rs.Primary.ID
		// Try to find the Slb server group
		if _, err := client.DescribeAccessControlListAttribute(req); err != nil {

			if IsExceptedError(err, SlbAclNotExists) {
				continue
			}
			return err
		}
		return fmt.Errorf("SLB Acl %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccSlbAclBasicConfig = `
variable "name" {
	default = "tf-testAccSlbAcl"
}
variable "ip_version" {
	default = "ipv4"
}

resource "alicloud_slb_acl" "foo" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list = [
    {
      entry="10.10.10.0/24"
      comment="first"
    },
    {
      entry="168.10.10.0/24"
      comment="second"
    }
  ]
}
`

const testAccSlbAclBasicConfigUpdate = `
variable "name" {
	default = "tf-testAccSlbAclUpdate"
}
variable "ip_version" {
	default = "ipv4"
}

resource "alicloud_slb_acl" "foo" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list = [
    {
      entry="10.10.10.0/24"
      comment="first"
    },
    {
      entry="168.10.10.0/24"
      comment="second"
    },
    {
      entry="172.10.10.0/24"
      comment="second"
    }
  ]
}
`
