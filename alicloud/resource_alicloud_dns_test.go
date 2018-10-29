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

func TestAccAlicloudDns_basic(t *testing.T) {
	var v dns.DomainType

	randInt := acctest.RandInt()
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
				Config: testAccDnsConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsExists(
						"alicloud_dns.dns", &v),
					resource.TestCheckResourceAttr(
						"alicloud_dns.dns",
						"name",
						fmt.Sprintf("testdnsbasic%v.abc", randInt)),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := &dns.DescribeDomainInfoArgs{
			DomainName: rs.Primary.Attributes["name"],
		}

		raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainInfo(request)
		})
		log.Printf("[WARN] Domain id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.(dns.DomainType)
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := &dns.DescribeDomainInfoArgs{
			DomainName: rs.Primary.Attributes["name"],
		}

		_, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainInfo(request)
		})

		if err != nil && !IsExceptedErrors(err, []string{InvalidDomainNameNoExist}) {
			return fmt.Errorf("Error Domain still exist.")
		}
	}

	return nil
}

func testAccDnsConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsbasic%v.abc"
}
`, randInt)
}
