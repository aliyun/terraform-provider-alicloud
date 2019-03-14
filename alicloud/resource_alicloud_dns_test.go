package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDns_basic(t *testing.T) {
	var v *alidns.DescribeDomainInfoResponse

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
			{
				Config: testAccDnsConfig_create(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsExists("alicloud_dns.dns", v),
					resource.TestCheckResourceAttr("alicloud_dns.dns", "name", fmt.Sprintf("testdnsbasic%v.abc", randInt)),
					resource.TestCheckResourceAttrSet("alicloud_dns.dns", "dns_server.#"),
				),
			},
			{
				Config: testAccDnsConfig_group_id(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsExists("alicloud_dns.dns", v),
					resource.TestCheckResourceAttr("alicloud_dns.dns", "name", fmt.Sprintf("testdnsbasic%v.abc", randInt)),
					resource.TestCheckResourceAttrSet("alicloud_dns.dns", "group_id"),
					resource.TestCheckResourceAttrSet("alicloud_dns.dns", "dns_server.#"),
				),
			},
		},
	})

}

func testAccCheckDnsExists(n string, domain *alidns.DescribeDomainInfoResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Domain ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		domainInfo, err := dnsService.DescribeDns(rs.Primary.Attributes["name"])
		log.Printf("[WARN] Domain id %#v", rs.Primary.ID)

		if err == nil {
			domain = domainInfo
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckDnsDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns" {
			continue
		}

		// Try to find the domain
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		_, err := dnsService.DescribeDns(rs.Primary.Attributes["name"])

		if err != nil && !IsExceptedErrors(err, []string{InvalidDomainNameNoExist}) {
			return WrapError(err)
		}
	}

	return nil
}

func testAccDnsConfig_create(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsbasic%v.abc"
}
`, randInt)
}

func testAccDnsConfig_group_id(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "test-dns-group"
}

resource "alicloud_dns" "dns" {
  name = "testdnsbasic%v.abc"
  group_id = "${alicloud_dns_group.group.id}"
}
`, randInt)
}
