package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDns_basic(t *testing.T) {
	var v dns.DomainType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns.dns",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsExists(
						"alicloud_dns.dns", &v),
					resource.TestCheckResourceAttr(
						"alicloud_dns.dns",
						"name",
						"yufish.com"),
				),
			},
		},
	})

}

func testAccCheckDnsExists(n string, domain *dns.DomainType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.dnsconn

		request := &dns.DescribeDomainInfoArgs{
			DomainName: rs.Primary.Attributes["name"],
		}

		response, err := conn.DescribeDomainInfo(request)
		log.Printf("[WARN] Domain id %#v", rs.Primary.ID)

		if err == nil {
			*domain = response
			return nil
		}
		return fmt.Errorf("Error finding domain %#v", rs.Primary.ID)
	}
}

func testAccCheckDnsDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns" {
			continue
		}

		// Try to find the domain
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.dnsconn

		request := &dns.DescribeDomainInfoArgs{
			DomainName: rs.Primary.Attributes["name"],
		}

		_, err := conn.DescribeDomainInfo(request)

		if err != nil && !IsExceptedErrors(err, []string{InvalidDomainNameNoExist}) {
			return fmt.Errorf("Error Domain still exist.")
		}
	}

	return nil
}

const testAccDnsConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
}
`
