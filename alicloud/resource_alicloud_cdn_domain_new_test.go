package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_cdn_domain_new", &resource.Sweeper{
		Name: "alicloud_cdn_domain_new",
		F:    testSweepCdnDomains_new,
	})
}

func testSweepCdnDomains_new(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	var domains []cdn.PageData
	args := cdn.CreateDescribeUserDomainsRequest()
	args.PageNumber = requests.NewInteger(1)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	for {

		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.DescribeUserDomains(args)
		})
		if err != nil {
			log.Printf("Error retrieving CDN Domain new: %s", err)
		}
		addDebug(args.GetActionName(), raw)
		resp, _ := raw.(*cdn.DescribeUserDomainsResponse)
		if resp == nil || len(resp.Domains.PageData) < 1 {
			break
		}
		domains = append(domains, resp.Domains.PageData...)

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			log.Printf("Error get Next page Number: %s", err)
			break
		} else {
			args.PageNumber = page
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
		request := cdn.CreateDeleteCdnDomainRequest()
		request.DomainName = name
		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.DeleteCdnDomain(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CDN domain (%s): %s", name, err)
		}
		addDebug(request.GetActionName(), raw)
	}
	return nil
}

func TestAccAlicloudCdnDomainNew_withTypeOSS(t *testing.T) {
	var v cdn.GetDomainDetailModel
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain_new.domain",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy_new,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainConfig_withTypeOSS(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists_new("alicloud_cdn_domain_new.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "scope", "overseas"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "cdn_type", "web"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.content", "www.aliyuntest.com"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.type", "oss"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.priority", "20"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.weight", "10"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.port", "80"),
				),
			},
			{
				Config: testAccCdnDomainConfig_withTypeOSSChange(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists_new("alicloud_cdn_domain_new.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "scope", "overseas"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "cdn_type", "web"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.content", "www.aliyuntest.com"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.type", "oss"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.priority", "20"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.weight", "30"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.port", "80"),
				),
			},
		},
	})
}

func TestAccAlicloudCdnDomainNew_withTypeDomain(t *testing.T) {
	var v cdn.GetDomainDetailModel
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain_new.domain",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy_new,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainConfig_withTypeDomain(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists_new("alicloud_cdn_domain_new.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "scope", "overseas"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "cdn_type", "web"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.content", "www.aliyuntest.com"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.type", "domain"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.priority", "20"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.weight", "10"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.port", "80"),
				),
			},
			{
				Config: testAccCdnDomainConfig_withTypeDomainChange(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists_new("alicloud_cdn_domain_new.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "scope", "overseas"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "cdn_type", "web"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.content", "www.aliyuntest.com"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.type", "domain"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.priority", "20"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.weight", "30"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.port", "80"),
				),
			},
		},
	})
}

func TestAccAlicloudCdnDomainNew_withTypeIpaddr(t *testing.T) {
	var v cdn.GetDomainDetailModel
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain_new.domain",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy_new,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainConfig_withTypeIpaddr(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists_new("alicloud_cdn_domain_new.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "scope", "domestic"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "cdn_type", "web"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.content", "1.1.1.1"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.type", "ipaddr"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.priority", "20"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.weight", "10"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.port", "80"),
				),
			},
			{
				Config: testAccCdnDomainConfig_withTypeIpaddrChange(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists_new("alicloud_cdn_domain_new.domain", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "scope", "domestic"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "cdn_type", "web"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.content", "1.1.1.1"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.type", "ipaddr"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.priority", "40"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.weight", "10"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_new.domain", "sources.0.port", "80"),
				),
			},
		},
	})
}

func testAccCheckCdnDomainExists_new(n string, detail *cdn.GetDomainDetailModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Domain ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		cdnservice := &CdnService{client: client}
		domain, err := cdnservice.DescribeCdnDomain(rs.Primary.ID)

		log.Printf("[WARN] Domain id %#v", rs.Primary.ID)
		if err == nil {
			detail = domain
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckCdnDomainDestroy_new(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cdn_domain_new" {
			continue
		}

		// Try to find the domain
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		cdnservice := &CdnService{client: client}
		_, err := cdnservice.DescribeCdnDomain(rs.Primary.ID)

		if err != nil && !IsExceptedErrors(err, []string{InvalidDomainNotFound}) {
			return WrapError(err)
		}
	}

	return nil
}

func testAccCdnDomainConfig_withTypeOSS(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
      scope = "overseas"
      sources {
         content = "www.aliyuntest.com"
         type = "oss"
      }
	}`, rand)
}

func testAccCdnDomainConfig_withTypeOSSChange(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
	  scope = "overseas"
      sources {
         content = "www.aliyuntest.com"
         type = "oss"
         priority = 20
         port = 80
         weight = 30
      }
	}`, rand)
}

func testAccCdnDomainConfig_withTypeDomain(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
      scope = "overseas"
      sources {
         content = "www.aliyuntest.com"
         type = "domain"
      }
	}`, rand)
}

func testAccCdnDomainConfig_withTypeDomainChange(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
	  scope = "overseas"
      sources {
         content = "www.aliyuntest.com"
         type = "domain"
         priority = 20
         port = 80
         weight = 30
      }
	}`, rand)
}

func testAccCdnDomainConfig_withTypeIpaddr(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
      scope = "domestic"
      sources {
         content = "1.1.1.1"
         type = "ipaddr"
      }
	}`, rand)
}

func testAccCdnDomainConfig_withTypeIpaddrChange(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_cdn_domain_new" "domain" {
	  domain_name = "tf-testacc%d.xiaozhu.com"
	  cdn_type = "web"
	  scope = "domestic"
	  sources {
         content = "1.1.1.1"
         type = "ipaddr"
         priority = 40
         port = 80
         weight = 20
      }
	}`, rand)
}
