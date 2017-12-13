package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
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
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.ali_domain", "true"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_version_code(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceVersionCodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
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

func TestAccAlicloudDnsDomainsDataSource_group_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceGroupNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "2"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDomainsDataSourceAliDomainConfig = `
data "alicloud_dns_domains" "domain" {
  ali_domain = true
}`

const testAccCheckAlicloudDomainsDataSourceVersionCodeConfig = `
data "alicloud_dns_domains" "domain" {
  version_code = "mianfei"
}`

const testAccCheckAlicloudDomainsDataSourceNameRegexConfig = `
data "alicloud_dns_domains" "domain" {
  domain_name_regex = ".*"
}`

const testAccCheckAlicloudDomainsDataSourceGroupNameRegexConfig = `
data "alicloud_dns_domains" "domain" {
  group_name_regex = ".*"
}`
