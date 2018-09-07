package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsDomainsDataSource_ali_domain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceAliDomainConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.ali_domain", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceNameRegexConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDomainsDataSourceAliDomainConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testdnsalidomain%v"
}

resource "alicloud_dns" "dns" {
  name = "testdnsalidomain%v.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  ali_domain = "${alicloud_dns.dns.name == "" ? false : false}"
  group_name_regex = "${alicloud_dns_group.group.name}"
}`, randInt%1000, randInt%1000)
}

func testAccCheckAlicloudDomainsDataSourceNameRegexConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsnameregex%v.abc"
}
data "alicloud_dns_domains" "domain" {
  domain_name_regex = "${alicloud_dns.dns.name}"
}`, randInt%1000)
}
