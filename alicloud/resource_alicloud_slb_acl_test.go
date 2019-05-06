package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_slb_acl", &resource.Sweeper{
		Name: "alicloud_slb_acl",
		F:    testSweepSlbAcl,
	})
}

func testSweepSlbAcl(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	req := slb.CreateDescribeAccessControlListsRequest()
	req.RegionId = client.RegionId
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAccessControlLists(req)
	})
	if err != nil {
		return err
	}
	resp, _ := raw.(*slb.DescribeAccessControlListsResponse)

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
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteAccessControlList(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Slb Acl (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSlbAcl_basic(t *testing.T) {
	var acl *slb.DescribeAccessControlListAttributeResponse

	resourceId := "alicloud_slb_acl.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &acl, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbAclBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl",
						"ip_version":   "ipv4",
						"entry_list.#": "2",
					}),
				),
			},
			{
				Config: testAccSlbAclBasicConfig_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbAcl_name",
					}),
				),
			},
			{
				Config: testAccSlbAclBasicConfig_entry_list,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entry_list.#": "3",
					}),
				),
			},
			{
				Config: testAccSlbAclBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl",
						"ip_version":   "ipv4",
						"entry_list.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbAcl_muilt(t *testing.T) {
	var acl *slb.DescribeAccessControlListAttributeResponse

	resourceId := "alicloud_slb_acl.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &acl, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbAclBasicConfig_mulit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl-9",
						"ip_version":   "ipv4",
						"entry_list.#": "2",
					}),
				),
			},
		},
	})
}

func testAccCheckSlbAclDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_acl" {
			continue
		}

		req := slb.CreateDescribeAccessControlListAttributeRequest()
		req.AclId = rs.Primary.ID
		// Try to find the Slb server group
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeAccessControlListAttribute(req)
		})
		if err != nil {

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

resource "alicloud_slb_acl" "default" {
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

const testAccSlbAclBasicConfig_name = `
variable "name" {
	default = "tf-testAccSlbAcl_name"
}
variable "ip_version" {
	default = "ipv4"
}

resource "alicloud_slb_acl" "default" {
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

const testAccSlbAclBasicConfig_entry_list = `
variable "name" {
	default = "tf-testAccSlbAcl_name"
}
variable "ip_version" {
	default = "ipv4"
}

resource "alicloud_slb_acl" "default" {
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
      comment="third"
    }
  ]
}
`
const testAccSlbAclBasicConfig_mulit = `
variable "name" {
	default = "tf-testAccSlbAcl"
}
variable "ip_version" {
	default = "ipv4"
}
variable "count" {
	default = "10"
}
resource "alicloud_slb_acl" "default" {
	count = "${var.count}"
  name = "${var.name}-${count.index}"
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
