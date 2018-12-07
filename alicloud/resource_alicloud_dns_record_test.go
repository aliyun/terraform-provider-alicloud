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

func TestAccAlicloudDnsRecord_multi(t *testing.T) {
	var v dns.RecordTypeNew
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_dns_record.record.9",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsRecordMulti(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record.9", &v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record.9", "type", "CNAME"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record.9", "ttl", "600"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record.9", "priority", "0"),
					resource.TestCheckResourceAttrSet("alicloud_dns_record.record.9", "name"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record.9", "host_record", "alimail"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record.9", "value", "mail.mxhichina9.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record.9", "routing", "default"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record.9", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record.9", "locked", "false"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecord_routing(t *testing.T) {
	var v dns.RecordTypeNew

	randInt := acctest.RandInt()
	dnsName := fmt.Sprintf("testdnsrecordrouting%v.abc", randInt)

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
				Config: testAccDnsRecordRouting(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", &v),
					resource.TestCheckResourceAttrSet("alicloud_dns_record.record", "id"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name", dnsName),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimail"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "CNAME"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttrSet("alicloud_dns_record.record", "ttl"),
					resource.TestCheckResourceAttrSet("alicloud_dns_record.record", "priority"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "oversea"),
					resource.TestCheckResourceAttrSet("alicloud_dns_record.record", "status"),
					resource.TestCheckResourceAttrSet("alicloud_dns_record.record", "locked"),
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
			if IsExceptedErrors(err, []string{DomainRecordNotBelongToUser}) {
				continue
			}
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

func testAccDnsRecordMulti(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordpriority%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichina${count.index}.com"
  count = 10
}
`, randInt)
}

func testAccDnsRecordRouting(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "testdnsrecordrouting%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  routing = "oversea"
  count = 1
}
`, randInt)
}
