package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCdnDomainConfig_basic(t *testing.T) {
	var v cdn.DomainConfig
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain_config.config",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainConfig_basic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnConfigExists("alicloud_cdn_domain_config.config", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_config.config", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_config.config", "function_name", "ip_allow_list_set"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_config.config", "function_args.0.arg_name", "ip_list"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_config.config", "function_args.0.arg_value", "110.110.110.110"),
				),
			},
		},
	})
}

func TestAccAlicloudCdnDomainConfig_other(t *testing.T) {
	var v cdn.DomainConfig
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cdn_domain_config.config",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainConfig_other(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnConfigExists("alicloud_cdn_domain_config.config", &v),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_config.config", "domain_name", fmt.Sprintf("tf-testacc%d.xiaozhu.com", rand)),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_config.config", "function_name", "referer_white_list_set"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_config.config", "function_args.0.arg_name", "refer_domain_allow_list"),
					resource.TestCheckResourceAttr("alicloud_cdn_domain_config.config", "function_args.0.arg_value", "110.110.110.110"),
				),
			},
		},
	})
}

func testAccCheckCdnConfigExists(n string, config *cdn.DomainConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No configId ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		cdnservice := &CdnService{client: client}
		domainconfig, err := cdnservice.DescribeDomainConfig(rs.Primary.ID)
		config = domainconfig

		log.Printf("[WARN] config id %#v", rs.Primary.ID)
		if err != nil {
			return WrapError(err)
		}

		return nil
	}
}

func testAccCheckCdnConfigDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cdn_domain_config" {
			continue
		}

		// Try to find the config
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		cdnservice := &CdnService{client: client}
		_, err := cdnservice.DescribeDomainConfig(rs.Primary.ID)

		if err != nil && !IsExceptedError(err, InvalidDomainNotFound) && !NotFoundError(err) {
			return WrapError(err)
		}
	}
	return nil
}

func testAccCdnDomainConfig_basic(rand int) string {
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
         weight = 10
      }
	}
    resource "alicloud_cdn_domain_config" "config" {
	  domain_name = "${alicloud_cdn_domain_new.domain.domain_name}"
	  function_name = "ip_allow_list_set"
      function_args {
            arg_name = "ip_list"
            arg_value = "110.110.110.110"
      }
	}`, rand)
}

func testAccCdnDomainConfig_other(rand int) string {
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
         weight = 10
      }
	}
    resource "alicloud_cdn_domain_config" "config" {
	  domain_name = "${alicloud_cdn_domain_new.domain.domain_name}"
	  function_name = "referer_white_list_set"
      function_args {
            arg_name = "refer_domain_allow_list"
            arg_value = "110.110.110.110"
      }
	}`, rand)
}
