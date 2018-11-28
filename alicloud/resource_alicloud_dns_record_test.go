package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
				Config: testAccDnsRecordConfig(acctest.RandInt()),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := &dns.DescribeDomainRecordInfoNewArgs{
			RecordId: rs.Primary.ID,
		}

		raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainRecordInfoNew(request)
		})
		log.Printf("[WARN] Domain record id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.(*dns.DescribeDomainRecordInfoNewResponse)
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
				Config: testAccDnsRecordPriority(acctest.RandInt()),
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := &dns.DescribeDomainRecordInfoNewArgs{
			RecordId: rs.Primary.ID,
		}

		raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainRecordInfoNew(request)
		})
		if err != nil {
			return err
		}
		response, _ := raw.(*dns.DescribeDomainRecordInfoNewResponse)
		if response.RecordId != "" {
			return fmt.Errorf("Error Domain record still exist.")
		}
	}

	return nil
}

func testAccDnsRecordConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}

resource "alicloud_dns_record" "record_a1" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.1"
  count = 1
}

resource "alicloud_dns_record" "record_a2" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.2"
  count = 1
}

resource "alicloud_dns_record" "record_a3" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.3"
  count = 1
}

resource "alicloud_dns_record" "record_a4" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.4"
  count = 1
}

resource "alicloud_dns_record" "record_a5" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.5"
  count = 1
}

resource "alicloud_dns_record" "record_a6" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.6"
  count = 1
}

resource "alicloud_dns_record" "record_a7" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.7"
  count = 1
}

resource "alicloud_dns_record" "record_a8" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.8"
  count = 1
}

resource "alicloud_dns_record" "record_a9" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.9"
  count = 1
}

resource "alicloud_dns_record" "record_a10" {
  name = "${alicloud_dns.dns.name}"
  host_record = "rr_a"
  type = "A"
  value = "1.1.1.10"
  count = 1
}
`, randInt)
}

func testAccDnsRecordPriority(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordpriority%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alipriority"
  type = "MX"
  value = "www.aliyun.com"
  count = 1
  priority = 10
}
`, randInt)
}
