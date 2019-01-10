package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/denverdino/aliyungo/cdn"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

func TestAccAlicloudCdnDomain_basic(t *testing.T) {
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

func TestAccAlicloudCdnDomain_https(t *testing.T) {
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
			resource.TestStep{
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

		if err != nil && !IsExceptedErrors(err, []string{InvalidDomainNotFound}) {
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

const testServerCertificate = `-----BEGIN CERTIFICATE-----\nMIICrDCCAZQCCQDApyXUTYDE+DANBgkqhkiG9w0BAQsFADAYMRYwFAYDVQQDDA0q\nLnhpYW96aHUuY29tMB4XDTE4MTIzMTEwMDY1OVoXDTE5MTIzMTEwMDY1OVowGDEW\nMBQGA1UEAwwNKi54aWFvemh1LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC\nAQoCggEBAM51j0hM0dqY/gHed5LhsSbTShXmw2Sqg18avBpHup/C8YzvpVddvRUe\ndE54idcCfdCVhoWEtcNFXIkKpWDpImFtovoupdNdC2XoBHEC42or9ZA+wVsaCih+\nSRuB3r+yNzXHh+rmUa0FLeij/x0gUovmznUsk7UMMzwJLZwmuyi8LCuiTIlQzQ9R\nTdaYo2t4OuGVkdQiJzsYiRRNPOqCKtYvYcEBLalLFOcVn0aG/I9Fn0P3rc8fK9BE\nHaoYRunmusUCdCcKpisHHKYdtmd3Zgz+Z+PBkjARtufO6kOXSov6u7azLQxZgZJm\neagNkqQ4/R+9b4GQ6crj0Pi655QWZVECAwEAATANBgkqhkiG9w0BAQsFAAOCAQEA\nJAWcQb7942jWDWjFqW5C4eyJBuJxK8dsTsYpOa0ccpp/A2cVCTNMzV/sMCGROy2B\nuhrTMG4q3QawFwV+K4nQX8yJUJp57zyUfdit+yOD4hDzPyF5MOlOwzdqeg9Y9ODL\nSc2O6J3tFsx2332Hb7RXFlYodEk1uFezrj4YLVrAgk1QojaBHpuFiA/O3eCmErjW\nc88bcil60qMvlhBhCWaZzivji6oc3y/6qWRZ3apxp9/6sjOvlm5Q+wwoLHXaU+L6\nw7TZptw6upuSrAS+N4Rwqgn5TfPvbGkdRG0X0TDzE1+F+w377kE71nV8lmi8F786\nbqloD1gWJER071WRz1coHQ==\n-----END CERTIFICATE-----\n`
const testPrivateKey = `-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDOdY9ITNHamP4B\n3neS4bEm00oV5sNkqoNfGrwaR7qfwvGM76VXXb0VHnROeInXAn3QlYaFhLXDRVyJ\nCqVg6SJhbaL6LqXTXQtl6ARxAuNqK/WQPsFbGgoofkkbgd6/sjc1x4fq5lGtBS3o\no/8dIFKL5s51LJO1DDM8CS2cJrsovCwrokyJUM0PUU3WmKNreDrhlZHUIic7GIkU\nTTzqgirWL2HBAS2pSxTnFZ9GhvyPRZ9D963PHyvQRB2qGEbp5rrFAnQnCqYrBxym\nHbZnd2YM/mfjwZIwEbbnzupDl0qL+ru2sy0MWYGSZnmoDZKkOP0fvW+BkOnK49D4\nuueUFmVRAgMBAAECggEAGFrP2zSMsN/JXwkSS/Zpwm28WJcPR6nBs49gzyzU/BGw\nEvMWKxc4vewIxlT71axKkTeCVe/QzUc6YkQqPCNkVd/sEN093JAmTxAurfIsR5MF\n9c0hXBDXT+2NzDvmvfBVCPgPtYsT6Xgp8T6fUp1Ef5JrmnD2v62/wX5Hrhr3ixdp\nFoEJ92/FJN+Dt/Duri481+RYqUUbNqiZAuEMsny/hM+aijtpQgZK8FqnVfjIOo0l\np4TwfDoxtTwkj0QeJxelX/sNSD+iQ5xVmnxQ9Q1qmKMk0M1Wl8318ovcrXDTN3RZ\nR/Fp+lNw0aC4OLuOJDc7A13l5unHmex4Ka1kbjOZwQKBgQD8w0hJiVVmiojJ3/YX\n06rG3ufGlKXu4qvRL6eXNqnDp8AnKXpjKiDV5lGz+vK0pcAyF5AhpddDN20Un+nR\n4VU7NslNJdBfjQxfEKJRkGBP22C+pdlxOBDgdgouN9q7TTr9StaQC0KF+HUcF8fO\nsYpEp1R+27iHie+F6uS78ur0mQKBgQDRGnbv8acTDTBW8nll1wTRa0+iVR2yxglS\nndC97HzYHOFb+caVmaM18QzHaaRvAhQ0bZXOjFV0QVv1nytlmHxsMlG/N6LDz5jA\nNLGvbrzDzki+RUk0TWJuYpvdevsnwaCxMwyY6CQ5MKhKAvwgE8AhdAq5MGcE+6U/\noAOZNyOxeQKBgQCiCB2i5mLUpSIjJ2r+wzXK3sH9zvTAOpaiNsZcbTJOto67jB9k\nynDaLhdaJRjJLSgT9H700vc3o6RNgGXHoYedufU5e3AkkKrJlkQ3vTHAf4V5MaA+\nsA5BlenYzv1s7IlQLlV1aYJvl2Kba7MukSlt8UZ9PCUC3i2pz3Zp9cMgoQKBgQCK\nGkR7bMq/1nIausJa9IwGFC3gNP8MV6dInVqEVXCO+2QL7wetPm+A7NdXzPoBJwpZ\nJhdO93ho89Hcg2eSDgf/HazH8eLaGH32U9cW2rhpShDZOcGDfaiI5y+yM8s1Erki\nz2h+hLOH4g8D8ry6ItE+RvneHY2syNb3EqPNyZEVYQKBgGeC4KiyTpW+SqmkDRU/\n/FxIAlZE2vWmtJptAUkmLEVoBi0/GvJoYYJR4Jhj4gS9tJcrI860X0j05NZHvIX8\nUJ/vZZeiFRiatXYUtKbD/ngoOXTF92ew2J3eCwNHH5OKwEO3aejaYwD9m0QM4kTl\nVLWIf79+G3NrxIeqPnqv+l3J\n-----END PRIVATE KEY-----`

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
