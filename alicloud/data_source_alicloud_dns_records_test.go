package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsRecordsDataSource_host_record_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "alimail"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceTypeConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_value_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichina.com"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_line(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceLineConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_status(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceStatusConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_is_locked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.locked", "false"),
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
  value = "mail.mxhichina.com"
  count = 1
}

data "alicloud_dns_records" "record" {
  domain_name = "${alicloud_dns_record.record.name}"
  value_regex = "^mail"
}`, randInt)
}

func testAccCheckAlicloudDnsRecordsDataSourceStatusConfig(randInt int) string {
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

func testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordislocked%v.abc"
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

func testAccCheckAlicloudDnsRecordsDataSourceLineConfig(randInt int) string {
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
