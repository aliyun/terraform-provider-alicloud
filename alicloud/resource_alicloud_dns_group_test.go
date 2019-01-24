package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDnsGroup_basic(t *testing.T) {
	var v dns.DomainGroupType
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_group.group",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsGroupConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsGroupExists(
						"alicloud_dns_group.group", &v),
					resource.TestCheckResourceAttr(
						"alicloud_dns_group.group",
						"name",
						fmt.Sprintf("tf-testacc-%d", rand)),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := &dns.DescribeDomainGroupsArgs{
			KeyWord: rs.Primary.Attributes["name"],
		}

		raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		log.Printf("[WARN] Group id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.([]dns.DomainGroupType)
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := &dns.DescribeDomainGroupsArgs{
			KeyWord: rs.Primary.Attributes["name"],
		}

		raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		response, _ := raw.([]dns.DomainGroupType)
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

func testAccDnsGroupConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_dns_group" "group" {
	  name = "tf-testacc-%d"
	}
	`, rand)
}
