package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDnsGroup_basic(t *testing.T) {
	var v dns.DomainGroupType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_group.group",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsGroupExists(
						"alicloud_dns_group.group", &v),
					resource.TestCheckResourceAttr(
						"alicloud_dns_group.group",
						"name",
						"yutest"),
				),
			},
		},
	})

}

func testAccCheckDnsGroupExists(n string, group *dns.DomainGroupType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain group ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.dnsconn

		request := &dns.DescribeDomainGroupsArgs{
			KeyWord: rs.Primary.Attributes["name"],
		}

		response, err := conn.DescribeDomainGroups(request)
		log.Printf("[WARN] Group id %#v", rs.Primary.ID)

		if err == nil {
			if response != nil && len(response) > 0 {
				*group = response[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding domain group %#v", rs.Primary.ID)
	}
}

func testAccCheckDnsGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns_group" {
			continue
		}

		// Try to find the domain group
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.dnsconn

		request := &dns.DescribeDomainGroupsArgs{
			KeyWord: rs.Primary.Attributes["name"],
		}

		response, err := conn.DescribeDomainGroups(request)

		if response != nil && len(response) > 0 {
			return fmt.Errorf("Error groups still exist")
		}

		if err != nil {
			// Verify the error is what we want
			return err
		}
	}
	return nil
}

const testAccDnsGroupConfig = `
resource "alicloud_dns_group" "group" {
  name = "yutest"
}
`
