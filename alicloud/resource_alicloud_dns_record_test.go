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
	rand := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_record.record",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecord_create(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimail"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "CNAME"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "600"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "default"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
				),
			},
			{
				Config: testAccDnsRecord_host_record(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimailchange"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "CNAME"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "600"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "default"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
				),
			},
			{
				Config: testAccDnsRecord_type(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimailchange"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "600"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "2"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "default"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
				),
			},

			{
				Config: testAccDnsRecord_priority(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimailchange"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.change.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "600"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "3"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "default"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
				),
			},
			{
				Config: testAccDnsRecord_value(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimailchange"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.change.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "600"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "3"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "default"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
				),
			},
			{
				Config: testAccDnsRecord_ttl(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimailchange"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.change.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "800"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "3"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "default"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
				),
			},
			{
				Config: testAccDnsRecord_routing(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimailchange"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.change.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "800"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "3"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "telecom"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
				),
			},
			{
				Config: testAccDnsRecord_ttl2(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimailchange"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.change.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "600"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "3"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "telecom"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
				),
			},
			{
				Config: testAccDnsRecord_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists("alicloud_dns_record.record", v),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "name",
						fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "host_record", "alimail"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "type", "CNAME"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "value", "mail.mxhichin.com"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "ttl", "600"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "routing", "default"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "status", "ENABLE"),
					resource.TestCheckResourceAttr("alicloud_dns_record.record", "locked", "false"),
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
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccDnsRecord_create(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
}
`, randInt)
}

func testAccDnsRecord_host_record(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimailchange"
  type = "CNAME"
  value = "mail.mxhichin.com"
}
`, randInt)
}

func testAccDnsRecord_type(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimailchange"
  type = "MX"
  priority = "2"
  value = "mail.mxhichin.com"
}
`, randInt)
}

func testAccDnsRecord_priority(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimailchange"
  type = "MX"
  value = "mail.change.com"
  priority = "3"
}
`, randInt)
}

func testAccDnsRecord_value(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimailchange"
  type = "MX"
  priority = "3"
  value = "mail.change.com"
}
`, randInt)
}

func testAccDnsRecord_ttl(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimailchange"
  type = "MX"
  priority = "3"
  value = "mail.change.com"
  ttl = "800"
}
`, randInt)
}

func testAccDnsRecord_routing(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimailchange"
  type = "MX"
  priority = "3"
  value = "mail.change.com"
  ttl = "800"
  routing = "telecom"
}
`, randInt)
}

func testAccDnsRecord_ttl2(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimailchange"
  type = "MX"
  priority = "3"
  value = "mail.change.com"
  ttl = "600"
  routing = "telecom"
}
`, randInt)
}

func testAccDnsRecord_all(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordbasic%v.abc"
}

resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimail"
  type = "CNAME"
  value = "mail.mxhichin.com"
  ttl = "600"
  priority = "1"
  routing = "default"
}
`, randInt)
}

func testAccDnsRecordMulti(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsrecordpriority%v.abc"
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
