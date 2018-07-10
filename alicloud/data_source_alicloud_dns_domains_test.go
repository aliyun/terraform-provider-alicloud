package alicloud

import (
	"testing"

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
				Config: testAccCheckAlicloudDomainsDataSourceAliDomainConfig,
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
				Config: testAccCheckAlicloudDomainsDataSourceNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDomainsDataSourceAliDomainConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
}

data "alicloud_dns_domains" "domain" {
  ali_domain = "${alicloud_dns.dns.name == "" ? false : false}"
}`

const testAccCheckAlicloudDomainsDataSourceNameRegexConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
}
data "alicloud_dns_domains" "domain" {
  domain_name_regex = "${alicloud_dns.dns.name}"
}`
