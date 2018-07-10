package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDnsRecord_basic(t *testing.T) {
	var v dns.RecordTypeNew

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_record.record",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsRecordConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(
						"alicloud_dns_record.record", &v),
					resource.TestCheckResourceAttr(
						"alicloud_dns_record.record",
						"type",
						"CNAME"),
				),
			},
		},
	})

}

func testAccCheckDnsRecordExists(n string, record *dns.RecordTypeNew) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain Record ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.dnsconn

		request := &dns.DescribeDomainRecordInfoNewArgs{
			RecordId: rs.Primary.ID,
		}

		response, err := conn.DescribeDomainRecordInfoNew(request)
		log.Printf("[WARN] Domain record id %#v", rs.Primary.ID)

		if err == nil {
			*record = response.RecordTypeNew
			return nil
		}
		return fmt.Errorf("Error finding domain record %#v", rs.Primary.ID)
	}
}

func TestAccAlicloudDnsRecord_priority(t *testing.T) {
	var v dns.RecordTypeNew

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_record.record",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsRecordPriority,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(
						"alicloud_dns_record.record", &v),
					resource.TestCheckResourceAttr(
						"alicloud_dns_record.record", "type", "MX"),
					resource.TestCheckResourceAttr(
						"alicloud_dns_record.record", "priority", "10"),
				),
			},
		},
	})

}

func testAccCheckDnsRecordDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns_record" {
			continue
		}

		// Try to find the domain record
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.dnsconn

		request := &dns.DescribeDomainRecordInfoNewArgs{
			RecordId: rs.Primary.ID,
		}

		response, err := conn.DescribeDomainRecordInfoNew(request)
		if err != nil {
			return err
		}
		if response.RecordId != "" {
			return fmt.Errorf("Error Domain record still exist.")
		}
	}

	return nil
}

const testAccDnsRecordConfig = `
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
`
const testAccDnsRecordPriority = `
resource "alicloud_dns" "dns" {
  name = "yufish.com"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alipriority"
  type = "MX"
  value = "www.aliyun.com"
  count = 1
  priority = 10
}
`
