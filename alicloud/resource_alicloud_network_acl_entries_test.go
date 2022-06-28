package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func SkipTestAccAlicloudVPCNetworkAclEntries_basic(t *testing.T) {
	resourceId := "alicloud_network_acl_entries.default"
	ra := resourceAttrInit(resourceId, testAccNaclEntriesCheckMap)
	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkAclEntriesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclEntries_create(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclEntriesExists(resourceId),
					testAccCheck(map[string]string{
						"ingress.#":                    "1",
						"egress.#":                     "1",
						"ingress.0.description":        "tf-testAcc_network_acl",
						"ingress.0.entry_type":         "custom",
						"ingress.0.name":               "tf-testAcc_network_acl",
						"ingress.0.policy":             "accept",
						"ingress.0.port":               "-1/-1",
						"ingress.0.protocol":           "all",
						"ingress.0.source_cidr_ip":     "0.0.0.0/32",
						"egress.0.description":         "tf-testAcc_network_acl",
						"egress.0.destination_cidr_ip": "0.0.0.0/32",
						"egress.0.entry_type":          "custom",
						"egress.0.name":                "tf-testAcc_network_acl",
						"egress.0.policy":              "accept",
						"egress.0.port":                "-1/-1",
						"egress.0.protocol":            "all",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkAclEntries_modify(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclEntriesExists(resourceId),
					testAccCheck(map[string]string{
						"ingress.#":                    "2",
						"egress.#":                     "2",
						"ingress.0.description":        "tf-testAcc_network_acl",
						"ingress.0.entry_type":         "custom",
						"ingress.0.name":               "tf-testAcc_network_acl",
						"ingress.0.policy":             "accept",
						"ingress.0.port":               "-1/-1",
						"ingress.0.protocol":           "all",
						"ingress.0.source_cidr_ip":     "0.0.0.0/32",
						"egress.0.description":         "tf-testAcc_network_acl",
						"egress.0.destination_cidr_ip": "0.0.0.0/32",
						"egress.0.entry_type":          "custom",
						"egress.0.name":                "tf-testAcc_network_acl",
						"egress.0.policy":              "accept",
						"egress.0.port":                "-1/-1",
						"egress.0.protocol":            "all",
					}),
				),
			},
			{
				Config: testAccNetworkAclEntries_delete(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclEntriesExists(resourceId),
					testAccCheck(map[string]string{
						"ingress.#":                    "1",
						"egress.#":                     "1",
						"ingress.0.description":        "tf-testAcc_network_acl",
						"ingress.0.entry_type":         "custom",
						"ingress.0.name":               "tf-testAcc_network_acl",
						"ingress.0.policy":             "accept",
						"ingress.0.port":               "-1/-1",
						"ingress.0.protocol":           "all",
						"ingress.0.source_cidr_ip":     "0.0.0.0/32",
						"egress.0.description":         "tf-testAcc_network_acl",
						"egress.0.destination_cidr_ip": "0.0.0.0/32",
						"egress.0.entry_type":          "custom",
						"egress.0.name":                "tf-testAcc_network_acl",
						"egress.0.policy":              "accept",
						"egress.0.port":                "-1/-1",
						"egress.0.protocol":            "all",
					}),
				),
			},
		},
	})
}

func testAccCheckNetworkAclEntriesExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(Error("Not found: %s", n))
		}
		if rs.Primary.ID == "" {
			return WrapError(Error("No Network Acl Entries ID is set"))
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		networkAclId := parts[0]

		if err := vpcService.WaitForNetworkAcl(networkAclId, Available, DefaultTimeout); err != nil {
			return WrapError(err)
		}
		return nil
	}
}

func testAccCheckNetworkAclEntriesDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_network_acl_entries" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		networkAclId := parts[0]

		object, err := vpcService.DescribeNetworkAcl(networkAclId)
		vpcResource := []vpc.Resource{}
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		for _, e := range resources {
			item := e.(map[string]interface{})
			vpcResource = append(vpcResource, vpc.Resource{
				ResourceId:   item["ResourceId"].(string),
				ResourceType: item["ResourceType"].(string),
			})
		}
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func testAccNetworkAclEntries_create(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}

resource "alicloud_network_acl_entries" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  ingress {
      protocol = "all"
      port = "-1/-1"
      source_cidr_ip = "0.0.0.0/32"
      name = "${var.name}"
      entry_type = "custom"
      policy = "accept"
      description = "${var.name}"
	}
  
  egress {
      protocol = "all"
      port = "-1/-1"
      destination_cidr_ip = "0.0.0.0/32"
      name = "${var.name}"
      entry_type = "custom"
      policy = "accept"
      description = "${var.name}"
    }
}
`, randInt)
}

func testAccNetworkAclEntries_modify(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}

resource "alicloud_network_acl_entries" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  ingress  {
      protocol = "all"
      port = "-1/-1"
      source_cidr_ip = "0.0.0.0/32"
      name = "${var.name}"
      entry_type = "custom"
      policy = "accept"
      description = "${var.name}"
    }

  ingress  {
      protocol = "all"
      port = "-1/-1"
      source_cidr_ip = "0.0.0.1/32"
      name = "${var.name}"
      entry_type = "custom"
      policy = "accept"
      description = "${var.name}"
    }
  
  egress {
      protocol = "all"
      port = "-1/-1"
      destination_cidr_ip = "0.0.0.0/32"
      name = "${var.name}"
      entry_type = "custom"
      policy = "accept"
      description = "${var.name}"
    }

  egress {
      protocol = "all"
      port = "-1/-1"
      destination_cidr_ip = "0.0.0.1/32"
      name = "${var.name}"
      entry_type = "custom"
      policy = "accept"
      description = "${var.name}"
    }
}
`, randInt)
}

func testAccNetworkAclEntries_delete(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}

resource "alicloud_network_acl_entries" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  ingress {
      protocol = "all"
      port = "-1/-1"
      source_cidr_ip = "0.0.0.0/32"
      name = "${var.name}"
      entry_type = "custom"
      policy = "accept"
      description = "${var.name}"
    }
  
  egress {
      protocol = "all"
      port = "-1/-1"
      destination_cidr_ip = "0.0.0.0/32"
      name = "${var.name}"
      entry_type = "custom"
      policy = "accept"
      description = "${var.name}"
    }
}
`, randInt)
}

var testAccNaclEntriesCheckMap = map[string]string{
	"network_acl_id": CHECKSET,
}
