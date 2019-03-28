package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsRecordsDataSource_domain_name(t *testing.T) {
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceDomainName(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.0", fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("testdnsrecordregex%v.abc", rand))),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.domain_name", fmt.Sprintf("testdnsrecordregex%v.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.locked", "false"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
			// The API DescribeDomainRecords required a correct domain name, otherwise return InvalidDomainName.NoExist
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_host_record_regex(t *testing.T) {
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig_match(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.0", fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("testdnsrecordregex%v.abc", rand))),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.domain_name", fmt.Sprintf("testdnsrecordregex%v.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.locked", "false"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig_mismatch(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_type(t *testing.T) {
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceTypeConfig_nonEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.0", fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("testdnsrecordtype%v.abc", rand))),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordtype%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceTypeConfig_empty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_value_regex(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig_match(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.0", fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("testdnsrecordvalueregex%v.abc", rand))),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordvalueregex%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig_mismatch(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_line(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceLineConfig_default(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.0", fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("testdnsrecordline%v.abc", rand))),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordline%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceLineConfig_nonDefault(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_status(t *testing.T) {
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceStatusConfig_enable(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.0", fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("testdnsrecordstatus%v.abc", rand))),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordstatus%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceStatusConfig_disable(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_is_locked(t *testing.T) {
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig_false(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.0", fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("testdnsrecordislocked%d.abc", rand))),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordislocked%d.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.locked", "false"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig_true(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_all(t *testing.T) {
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceAllConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "urls.0", fmt.Sprintf("%v.%v", "alimail", fmt.Sprintf("testdnsrecordall%d.abc", rand))),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordall%d.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.locked", "false"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDnsRecordsDataSourceDomainName(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordregex%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig_match(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordregex%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  host_record_regex = "^ali"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig_mismatch(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordregex%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  host_record_regex = "anyother"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceTypeConfig_nonEmpty(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordtype%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  type = "CNAME"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceTypeConfig_empty(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordtype%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  type = "TXT"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig_match(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordvalueregex%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  value_regex = "^mail"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig_mismatch(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordvalueregex%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  value_regex = "anyother"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceStatusConfig_disable(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordstatus%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  status = "disable"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceStatusConfig_enable(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordstatus%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  status = "enable"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig_false(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordislocked%d.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  is_locked = false
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig_true(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordislocked%d.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  is_locked = true
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceLineConfig_default(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordline%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  line = "default"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceLineConfig_nonDefault(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordline%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  line = "telecom"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceAllConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordall%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  host_record_regex = "^ali"
  value_regex = "^mail"
  type = "CNAME"
  line = "default"
  status = "enable"
  is_locked = "false"
}`, randInt)
}
