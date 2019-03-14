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

func TestAccAlicloudDnsGroup_basic(t *testing.T) {
	var v alidns.DomainGroup
	rand := acctest.RandIntRange(100, 999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_group.group",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsGroupConfig_create(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsGroupExists("alicloud_dns_group.group", &v),
					resource.TestCheckResourceAttr("alicloud_dns_group.group", "name", fmt.Sprintf("tf-testacc-c-%d", rand)),
				),
			},
			{
				Config: testAccDnsGroupConfig_update_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsGroupExists("alicloud_dns_group.group", &v),
					resource.TestCheckResourceAttr("alicloud_dns_group.group", "name", fmt.Sprintf("tf-testacc-name-%d", rand)),
				),
			},
		},
	})

}

func testAccDnsGroupConfig_create(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_dns_group" "group" {
	  name = "tf-testacc-c-%d"
	}
	`, rand)
}

func testAccDnsGroupConfig_update_name(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_dns_group" "group" {
	  name = "tf-testacc-name-%d"
	}
	`, rand)
}

func testAccCheckDnsGroupExists(n string, group *alidns.DomainGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Domain group ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		domaingroup, err := dnsService.DescribeDnsGroup(rs.Primary.Attributes["name"])
		if err != nil {
			return WrapError(err)
		}
		log.Printf("[WARN] Group id %#v", rs.Primary.ID)

		*group = domaingroup
		return nil
	}
}

func testAccCheckDnsGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns_group" {
			continue
		}

		// Try to find the domain group
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		groupId := rs.Primary.Attributes["name"]
		_, err := dnsService.DescribeDnsGroup(groupId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}

		return WrapError(Error("Error groups still exist"))
	}
	return nil
}
