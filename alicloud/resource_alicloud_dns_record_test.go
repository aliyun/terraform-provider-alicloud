package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDnsRecord_basic(t *testing.T) {
	var v *alidns.DescribeDomainRecordInfoResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_record.record",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(
						"alicloud_dns_record.record", v),
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
	var v *alidns.DescribeDomainRecordInfoResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_record.record",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordPriority(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(
						"alicloud_dns_record.record", v),
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
	var v *alidns.DescribeDomainRecordInfoResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_dns_record.record.9",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordMulti(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record.9", v),
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
	var v *alidns.DescribeDomainRecordInfoResponse

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
			{
				Config: testAccDnsRecordRouting(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
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

func testAccCheckDnsRecordExists(n string, record *alidns.DescribeDomainRecordInfoResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Domain Record ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		recordInfo, err := dnsService.DescribeDnsRecord(rs.Primary.ID)
		log.Printf("[WARN] Domain record id %#v", rs.Primary.ID)

		if err == nil {
			record = recordInfo
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckDnsRecordDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns_record" {
			continue
		}

		// Try to find the domain record
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		_, err := dnsService.DescribeDnsRecord(rs.Primary.ID)
		if err != nil {
			if IsExceptedErrors(err, []string{DomainRecordNotBelongToUser}) {
				continue
			}
			return WrapError(err)
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
