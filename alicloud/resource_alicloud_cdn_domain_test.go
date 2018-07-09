package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/cdn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCdnDomain_basic(t *testing.T) {
	var v cdn.DomainDetail

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain.domain",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCdnDomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists(
						"alicloud_cdn_domain.domain", &v),
					resource.TestCheckResourceAttr(
						"alicloud_cdn_domain.domain",
						"domain_name",
						"www.xiaozhu.com"),
				),
			},
		},
	})
}

func testAccCheckCdnDomainExists(n string, domain *cdn.DomainDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.cdnconn

		request := cdn.DescribeDomainRequest{
			DomainName: rs.Primary.Attributes["domain_name"],
		}

		response, err := conn.DescribeCdnDomainDetail(request)
		log.Printf("[WARN] Domain id %#v", rs.Primary.ID)

		if err == nil {
			*domain = response.GetDomainDetailModel
			return nil
		}
		return fmt.Errorf("Error finding domain %#v", rs.Primary.ID)
	}
}

func testAccCheckCdnDomainDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cdn_domain" {
			continue
		}

		// Try to find the domain
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.cdnconn

		request := cdn.DescribeDomainRequest{
			DomainName: rs.Primary.Attributes["domain_name"],
		}

		_, err := conn.DescribeCdnDomainDetail(request)

		if err != nil && !IsExceptedErrors(err, []string{InvalidDomainNotFound}) {
			return fmt.Errorf("Error Domain still exist.")
		}
	}

	return nil
}

const testAccCdnDomainConfig = `
resource "alicloud_cdn_domain" "domain" {
  domain_name = "www.xiaozhu.com"
  cdn_type = "web"
  source_type = "oss"
  sources = ["terraformtest.aliyuncs.com"]
  optimize_enable = "off"
  page_compress_enable = "off"
  range_enable = "off"
  video_seek_enable = "off"
}`
