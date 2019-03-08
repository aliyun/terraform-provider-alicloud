package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsRecordsDataSource_host_record_regex(t *testing.T) {
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordregex%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.ttl"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.priority"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.line"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.status"),
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
				Config: testAccCheckAlicloudDnsRecordsDataSourceTypeConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordtype%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.ttl"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.priority"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.line"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.status"),
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
				Config: testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordvalueregex%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.ttl"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.priority"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.line"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.status"),
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordline%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.ttl"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.priority"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.status"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceLineConfig_nonDefault(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordStatus%v.abc", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.ttl"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.priority"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.line"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceStatusConfig_disable(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "domain_name", fmt.Sprintf("testdnsrecordislocked%d.abc", rand)),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.locked", "false"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.ttl"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.priority"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.line"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_records.record", "records.0.status"),
				),
			},
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig_true(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceEmpty(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.record_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.domain_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.line"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.host_record"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.type"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.value"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.locked"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_records.record", "records.0.priority"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig(randInt int) string {
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

func testAccCheckAlicloudDnsRecordsDataSourceTypeConfig(randInt int) string {
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

func testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig(randInt int) string {
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
  status = "enable"
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
  status = "enable"
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

func testAccCheckAlicloudDnsRecordsDataSourceEmpty(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testaccdnsrecordline%v.abc"
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns.dns.name}"
  line = "default"
}`, randInt)
}
