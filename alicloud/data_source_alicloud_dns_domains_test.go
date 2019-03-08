package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsDomainsDataSource_ali_domain(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceAliDomainConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.domain_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_name", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.ali_domain", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.group_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_name", fmt.Sprintf("testaccdnsdomain%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.instance_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.version_code", "mianfei"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.puny_code"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.dns_servers.#"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "ids.0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.0", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudDomainsDataSourceAliDomainEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_group_name_regex(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceInstanceIdEmptyConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.domain_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_name", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.ali_domain", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.group_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_name", fmt.Sprintf("testaccdnsdomain%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.instance_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.version_code", "mianfei"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.puny_code"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.dns_servers.#"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "ids.0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "names.0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.0", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudDomainsDataSourceInstanceIdNonEmptyConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_instance_id(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceInstanceIdEmptyConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.domain_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_name", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.ali_domain", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.group_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_name", fmt.Sprintf("testaccdnsdomain%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.instance_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.version_code", "mianfei"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.puny_code"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.dns_servers.#"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "ids.0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "names.0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.0", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudDomainsDataSourceInstanceIdNonEmptyConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_version_code(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceVersionCodeConfigByFixed_mianfei(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.domain_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_name", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.ali_domain", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.group_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_name", fmt.Sprintf("testaccdnsdomain%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.instance_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.version_code", "mianfei"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.puny_code"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.dns_servers.#"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.0", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudDomainsDataSourceVersionCodeConfigByAnyother(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "0"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "names.#"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_all(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceAliDomainAllConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.domain_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_name", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.ali_domain", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.group_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_name", fmt.Sprintf("testaccdnsdomain%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.instance_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.version_code", "mianfei"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.puny_code"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_domains.domain", "domains.0.dns_servers.#"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "names.0", fmt.Sprintf("testaccdnsalidomain%d.abc", rand)),
				),
			},
		},
	})
}

func testAccCheckAlicloudDomainsDataSourceAliDomainConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  ali_domain = "false"
  domain_name_regex = "${alicloud_dns.dns.name}"
}`, randInt, randInt)
}

func testAccCheckAlicloudDomainsDataSourceAliDomainEmpty(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  ali_domain = "true"
  domain_name_regex = "${alicloud_dns.dns.name}"
}`, randInt, randInt)
}

func testAccCheckAlicloudDomainsDataSourceGroupNameRegexConfig_match(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  ali_domain = "false"
  group_name_regex = "${alicloud_dns_group.name}"
  depends_on = [ "${alicloud_dns.dns}" ]
}`, randInt, randInt)
}

func testAccCheckAlicloudDomainsDataSourceGroupNameRegexConfig_nonMatch(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  ali_domain = "false"
  group_name_regex = "${alicloud_dns_group.name}"
  depends_on = [ "alic" ]
}`, randInt, randInt)
}

func testAccCheckAlicloudDomainsDataSourceInstanceIdEmptyConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  domain_name_regex = "${alicloud_dns.dns.name}"
  instance_id = ""
}`, randInt, randInt)
}

func testAccCheckAlicloudDomainsDataSourceInstanceIdNonEmptyConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  domain_name_regex = "${alicloud_dns.dns.name}"
  instance_id = "122"
}`, randInt, randInt)
}

func testAccCheckAlicloudDomainsDataSourceVersionCodeConfigByFixed_mianfei(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  domain_name_regex = "${alicloud_dns.dns.name}"
  version_code = "mianfei"
}`, randInt, randInt)
}

func testAccCheckAlicloudDomainsDataSourceVersionCodeConfigByAnyother(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  domain_name_regex = "${alicloud_dns.dns.name}"
  version_code = "bumianfei"
}`, randInt, randInt)
}

func testAccCheckAlicloudDomainsDataSourceAliDomainAllConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "testaccdnsdomain%d"
}

resource "alicloud_dns" "dns" {
  name = "testaccdnsalidomain%d.abc"
  group_id = "${alicloud_dns_group.group.id}"
}

data "alicloud_dns_domains" "domain" {
  domain_name_regex = "${alicloud_dns.dns.name}"
  instance_id = ""
  version_code = "mianfei"
  ali_domain = "false"
  group_name_regex = "${alicloud_dns_group.group.name}"
}`, randInt, randInt)
}
