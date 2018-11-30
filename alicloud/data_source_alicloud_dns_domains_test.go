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

func TestAccAlicloudDnsDomainsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_domains.domain", "domains.0.ali_domain"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_domains.domain", "domains.0.version_code"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_domains.domain", "domains.0.puny_code"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_domains.domain", "domains.0.dns_servers.#"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDomainsDataSourceAliDomainConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsalidomain%v"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%v.abc"
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
  name = "testaccdnsnameregex%v.abc"
}
data "alicloud_dns_domains" "domain" {
  domain_name_regex = "${alicloud_dns.dns.name}"
}`, randInt%1000)
}

const testAccCheckAlicloudDomainsDataSourceEmpty = `
data "alicloud_dns_domains" "domain" {
  domain_name_regex = "^tf-testacc-fake-name"
}`
