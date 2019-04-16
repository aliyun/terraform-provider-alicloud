package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDnsRecord_basic(t *testing.T) {
	var v *alidns.DescribeDomainRecordInfoResponse
	resourceId := "alicloud_dns_record.record"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
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
					testAccCheck(map[string]string{"name": fmt.Sprintf("tf-testaccdnsrecordbasic%v.abc", rand)}),
				),
			},
			{
				Config: testAccDnsRecord_host_record(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"host_record": "alimailchange"}),
				),
			},
			{
				Config: testAccDnsRecord_type(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "MX",
						"priority": "2",
					}),
				),
			},

			{
				Config: testAccDnsRecord_priority(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"priority": "3"}),
				),
			},
			{
				Config: testAccDnsRecord_value(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"value": "mail.change.com"}),
				),
			},
			{
				Config: testAccDnsRecord_ttl(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"ttl": "800"}),
				),
			},
			{
				Config: testAccDnsRecord_routing(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"routing": "telecom"}),
				),
			},
			{
				Config: testAccDnsRecord_ttl2(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"ttl": "600"}),
				),
			},
			{
				Config: testAccDnsRecord_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})

}

func TestAccAlicloudDnsRecord_multi(t *testing.T) {
	var v *alidns.DescribeDomainRecordInfoResponse

	ra := resourceAttrInit("alicloud_dns_record.record.9", basicMap)
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit("alicloud_dns_record.record.9", &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
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
					testAccCheck(map[string]string{
						"value": "mail.mxhichina9.com",
					}),
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
  value = "mail.mxhichin.com"
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

var basicMap = map[string]string{
	"host_record": "alimail",
	"type":        "CNAME",
	"value":       "mail.mxhichin.com",
	"ttl":         "600",
	"priority":    "0",
	"routing":     "default",
	"status":      "ENABLE",
	"locked":      "false",
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

var multiMap = map[string]string{
	"host_record": "alimail",
	"type":        "CNAME",
	"value":       "mail.mxhichina9.com",
	"ttl":         "600",
	"priority":    "0",
	"name":        CHECKSET,
	"routing":     "default",
	"status":      "ENABLE",
	"locked":      "false",
}
