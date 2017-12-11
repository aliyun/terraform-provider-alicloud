package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.record_id", "3438492787133440"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.domain_name", "heguimin.top"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "smtp"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "ENABLE"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "smtp.mxhichina.com"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "7"),
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "3"),
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "17"),
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "17"),
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
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "17"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig = `
data "alicloud_dns_records" "record" {
  domain_name = "heguimin.top"
  host_record_regex = ".*smtp.*"
}`

const testAccCheckAlicloudDnsRecordsDataSourceTypeConfig = `
data "alicloud_dns_records" "record" {
  domain_name = "heguimin.top"
  type = "CNAME"
}`

const testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig = `
data "alicloud_dns_records" "record" {
  domain_name = "heguimin.top"
  value_regex = "^mail.mxhichina"
}`

const testAccCheckAlicloudDnsRecordsDataSourceStatusConfig = `
data "alicloud_dns_records" "record" {
  domain_name = "heguimin.top"
  status = "enable"
}`

const testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig = `
data "alicloud_dns_records" "record" {
  domain_name = "heguimin.top"
  is_locked = false
}`

const testAccCheckAlicloudDnsRecordsDataSourceLineConfig = `
data "alicloud_dns_records" "record" {
  domain_name = "heguimin.top"
  line = "default"
}`
