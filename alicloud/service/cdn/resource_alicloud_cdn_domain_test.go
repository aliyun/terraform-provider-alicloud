package cdn_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	alicdn "github.com/denverdino/aliyungo/cdn"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// collection-order: bypass parameter_filter_config.hash_key_args -- pre-existing stock TypeList on the deprecated alicloud_cdn_domain resource; schema and tests moved unchanged in this whole-product migration, reorder coverage was absent on master and remains pending the cdn_domain deprecation decision
// collection-order: bypass refer_config.refer_list -- pre-existing stock TypeList on the deprecated alicloud_cdn_domain resource; schema and tests moved unchanged in this whole-product migration, reorder coverage was absent on master and remains pending the cdn_domain deprecation decision

func init() {
	resource.AddTestSweepers("alicloud_cdn_domain", &resource.Sweeper{
		Name: "alicloud_cdn_domain",
		F:    testSweepCdnDomains,
	})
}

func testSweepCdnDomains(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	var domains []alicdn.Domains
	args := alicdn.DescribeDomainsRequest{}
	args.PageNumber = 1
	args.PageSize = helper.PageSizeLarge
	for {

		raw, err := client.WithCdnClient(func(cdnClient *alicdn.CdnClient) (interface{}, error) {
			return cdnClient.DescribeUserDomains(args)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving cdn domains: %s", err)
		}
		resp, _ := raw.(*alicdn.DomainsResponse)
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
		args := alicdn.DescribeDomainRequest{
			DomainName: name,
		}
		_, err := client.WithCdnClient(func(cdnClient *alicdn.CdnClient) (interface{}, error) {
			return cdnClient.DeleteCdnDomain(args)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN domain (%s): %s", name, err)
		}
	}
	return nil
}

// Converted from SkipTestAccAlicloudCDNDomain_basic to a runnable TestAccAliCloud
// representative for the deprecated alicloud_cdn_domain resource (denverdino/aliyungo/cdn
// SDK, factory resourceAlicloudCdnDomain, independent from alicloud_cdn_domain_new).
// Fixture adapted to the authorized domestic test account only: a filed subdomain of
// pfytlm.xyz and an account-owned OSS bucket as the oss source. The SkipTest https and
// domestic-cert scenarios, all assertions, the sweep and the production Schema/CRUD/SDK
// remain unchanged, pending the cdn_domain deprecation decision.
func TestAccAliCloudCDNDomain_basic(t *testing.T) {
	var v alicdn.DomainDetail
	rand := tfacctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists(
						"alicloud_cdn_domain.domain", &v),
					resource.TestCheckResourceAttr(
						"alicloud_cdn_domain.domain",
						"domain_name",
						fmt.Sprintf("tf-testacc%d.pfytlm.xyz", rand)),
				),
			},
		},
	})
}

// Skip this testcase because of the old API can not support overseas.
// It should be open after using new API
func SkipTestAccAlicloudCDNDomain_https(t *testing.T) {
	var v alicdn.DomainDetail
	rand := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainHttpsConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("alicloud_cdn_domain.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "domain_name", fmt.Sprintf("tf-testacc%d.alicloud-provider.cn", rand)),
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
	var v alicdn.DomainDetail
	rand := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainDomestic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("alicloud_cdn_domain.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain.domain", "domain_name", fmt.Sprintf("tf-testacc%d.alicloud-provider.cn", rand)),
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

func testAccCheckCdnDomainExists(n string, domain *alicdn.DomainDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain ID is set")
		}

		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		request := alicdn.DescribeDomainRequest{
			DomainName: rs.Primary.Attributes["domain_name"],
		}

		raw, err := client.WithCdnClient(func(cdnClient *alicdn.CdnClient) (interface{}, error) {
			return cdnClient.DescribeCdnDomainDetail(request)
		})
		log.Printf("[WARN] Domain id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.(alicdn.DomainResponse)
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
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		request := alicdn.DescribeDomainRequest{
			DomainName: rs.Primary.Attributes["domain_name"],
		}

		_, err := client.WithCdnClient(func(cdnClient *alicdn.CdnClient) (interface{}, error) {
			return cdnClient.DescribeCdnDomainDetail(request)
		})

		if err != nil && !errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return fmt.Errorf("Error Domain still exist.")
		}
	}

	return nil
}

func testAccCdnDomainConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_oss_bucket" "default" {
	  bucket = "tf-testacc-cdn-%d"
	}

	resource "alicloud_cdn_domain" "domain" {
	  domain_name = "tf-testacc%d.pfytlm.xyz"
	  cdn_type = "web"
	  source_type = "oss"
	  sources = ["${alicloud_oss_bucket.default.bucket}.${alicloud_oss_bucket.default.extranet_endpoint}"]
	  optimize_enable = "off"
	  page_compress_enable = "off"
	  range_enable = "off"
	  video_seek_enable = "off"
	}`, rand, rand)
}

func testAccCdnDomainHttpsConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain" "domain" {
	  domain_name = "tf-testacc%d.alicloud-provider.cn"
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
	  domain_name = "tf-testacc%d.alicloud-provider.cn"
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
