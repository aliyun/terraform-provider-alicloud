package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_cdn_domain", &resource.Sweeper{
		Name: "alicloud_cdn_domain",
		F:    testSweepCdnDomains,
	})
}

func testSweepCdnDomains(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	var domains []cdn.Domains
	args := cdn.DescribeDomainsRequest{}
	args.PageNumber = 1
	args.PageSize = PageSizeLarge
	for {

		raw, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.DescribeUserDomains(args)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving cdn domains: %s", err)
		}
		resp, _ := raw.(*cdn.DomainsResponse)
		if resp == nil || len(resp.Domains.PageData) < 1 {
			break
		}
		domains = append(domains, resp.Domains.PageData...)

		if resp.NextPage() == nil {
			break
		}
	}

	for _, v := range domains {
		name := v.DomainName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping CDN domain: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting CDN domain: %s", name)
		args := cdn.DescribeDomainRequest{
			DomainName: name,
		}
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.DeleteCdnDomain(args)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN domain (%s): %s", name, err)
		}
	}
	return nil
}

func SkipTestAccAlicloudCDNDomain_basic(t *testing.T) {
	var v cdn.DomainDetail
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain.domain",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists(
						"alicloud_cdn_domain.domain", &v),
					resource.TestCheckResourceAttr(
						"alicloud_cdn_domain.domain",
						"domain_name",
						fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
				),
			},
		},
	})
}

// Skip this testcase because of the old API can not support overseas.
// It should be open after using new API
func SkipTestAccAlicloudCDNDomain_https(t *testing.T) {
	var v cdn.DomainDetail
	rand := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain.domain",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainHttpsConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("alicloud_cdn_domain.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "certificate_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "certificate_config.0.server_certificate_status", "on"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "certificate_config.0.server_certificate", strings.Replace(testServerCertificate, `\n`, "\n", -1)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "certificate_config.0.private_key", strings.Replace(testPrivateKey, `\n`, "\n", -1)),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCDNDomain_domestic(t *testing.T) {
	var v cdn.DomainDetail
	rand := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain.domain",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainDomestic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("alicloud_cdn_domain.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "certificate_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "certificate_config.0.server_certificate_status", "on"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "certificate_config.0.server_certificate", strings.Replace(testServerCertificate, `\n`, "\n", -1)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "certificate_config.0.private_key", strings.Replace(testPrivateKey, `\n`, "\n", -1)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "scope", "domestic"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := cdn.DescribeDomainRequest{
			DomainName: rs.Primary.Attributes["domain_name"],
		}

		raw, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.DescribeCdnDomainDetail(request)
		})
		log.Printf("[WARN] Domain id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.(cdn.DomainResponse)
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := cdn.DescribeDomainRequest{
			DomainName: rs.Primary.Attributes["domain_name"],
		}

		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.DescribeCdnDomainDetail(request)
		})

		if err != nil && !IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return fmt.Errorf("Error Domain still exist.")
		}
	}

	return nil
}

func testAccCdnDomainConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
	  source_type = "oss"
	  sources = ["terraformtest.aliyuncs.com"]
	  optimize_enable = "off"
	  page_compress_enable = "off"
	  range_enable = "off"
	  video_seek_enable = "off"
	}`, rand)
}

func testAccCdnDomainHttpsConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
	  source_type = "oss"
	  scope="overseas"
	  sources = ["terraformtest.aliyuncs.com"]
	  certificate_config {
		server_certificate = "%s"
		private_key = "%s"
 	  }
	}`, rand, testServerCertificate, testPrivateKey)
}

func testAccCdnDomainDomestic(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
	  source_type = "oss"
	  scope="domestic"
	  sources = ["terraformtest.aliyuncs.com"]
	  certificate_config {
		server_certificate = "%s"
		private_key = "%s"
 	  }
	}`, rand, testServerCertificate, testPrivateKey)
}
