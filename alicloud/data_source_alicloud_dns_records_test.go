package alicloud

import (
	"testing"

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
				Config: testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig,
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
				Config: testAccCheckAlicloudDnsRecordsDataSourceTypeConfig,
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
				Config: testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig,
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
				Config: testAccCheckAlicloudDnsRecordsDataSourceLineConfig,
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
				Config: testAccCheckAlicloudDnsRecordsDataSourceStatusConfig,
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
				Config: testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.locked", "false"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
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
}`

const testAccCheckAlicloudDnsRecordsDataSourceTypeConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
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
}`

const testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
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
}`

const testAccCheckAlicloudDnsRecordsDataSourceStatusConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
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
}`

const testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
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
}`

const testAccCheckAlicloudDnsRecordsDataSourceLineConfig = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
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
}`
