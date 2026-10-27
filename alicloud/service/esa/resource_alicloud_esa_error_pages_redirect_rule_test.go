// Package esa. This file is generated automatically. Please do not modify it manually, thank you!
package esa_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Test Esa ErrorPagesRedirectRule. >>> Resource test cases, automatically generated.
// Case resource_ErrorPagesRedirectRule_test 12504
func TestAccAliCloudEsaErrorPagesRedirectRule_basic12504(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccesa%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEsaErrorPagesRedirectRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule", "(http.host eq \"video.example.com\")"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", name),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "3"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "500"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.1.target_url", "https://example.com/foo"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.1.status_code", "400"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.2.target_url", "https://example.com/test"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.2.status_code", "503"),

					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "config_id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_esa_error_pages_redirect_rule.default", tfjsonpath.New("config_type"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "off"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule", "(http.host eq \"video.example.com\")"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", fmt.Sprintf("%s-mod", name)),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "3"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "500"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.1.target_url", "https://example.com/foo"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.1.status_code", "400"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.2.target_url", "https://example.com/test"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.2.status_code", "503"),
				),
			},
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule", "(http.host eq \"video.example.com\")"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", fmt.Sprintf("%s-mod", name)),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "3"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "500"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.1.target_url", "https://example.com/foo"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.1.status_code", "400"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.2.target_url", "https://example.com/test"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.2.status_code", "503"),
				),
			},
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule", "(http.request.uri eq \"/content?page=1234\")"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", fmt.Sprintf("%s-mod", name)),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "3"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "500"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.1.target_url", "https://example.com/foo"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.1.status_code", "400"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.2.target_url", "https://example.com/test"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.2.status_code", "503"),
				),
			},
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule", "(http.request.uri eq \"/content?page=1234\")"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", fmt.Sprintf("%s-mod", name)),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "1"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/test"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "404"),
				),
			},
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step5(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule", "(http.request.uri eq \"/content?page=1234\")"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", fmt.Sprintf("%s-mod", name)),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "1"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/foo"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "400"),
				),
			},
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step6(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule", "(http.request.uri eq \"/content?page=1234\")"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", fmt.Sprintf("%s-mod", name)),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "1"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "500"),
				),
			},
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step7(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "on"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule", "(http.request.uri eq \"/content?page=1234\")"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", fmt.Sprintf("%s-mod", name)),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "1"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "500"),
				),
			},
			{
				Config: testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step8(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "site_id"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_enable", "off"),

					resource.TestCheckResourceAttrSet("alicloud_esa_error_pages_redirect_rule.default", "rule"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "sequence", "1"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "site_version", "0"),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "rule_name", name),

					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.#", "1"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.target_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr("alicloud_esa_error_pages_redirect_rule.default", "error_pages_redirect.0.status_code", "500"),
				),
			},
			{
				ResourceName:            "alicloud_esa_error_pages_redirect_rule.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// testAccCheckEsaErrorPagesRedirectRuleDestroy asserts that every ESA Error
// Pages Redirect Rule created during the test is gone after destruction. It is
// what makes a no-op Delete (one that returns nil without issuing the RPC)
// fail the suite: the resource would still exist and be reported here.
func testAccCheckEsaErrorPagesRedirectRuleDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_esa_error_pages_redirect_rule" {
			continue
		}
		object, err := esa.DescribeErrorPagesRedirectRule(context.Background(), client, rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("checking alicloud_esa_error_pages_redirect_rule %s destroy: %s", rs.Primary.ID, err)
		}
		if object != nil {
			return fmt.Errorf("alicloud_esa_error_pages_redirect_rule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_test_ErrorPagesRedirectRule" {
  type         = "NS"
  auto_renew   = false
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_test_ErrorPagesRedirectRule" {
  site_name   = "${var.name}.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_test_ErrorPagesRedirectRule.id
  coverage    = "overseas"
  access_type = "NS"
  flatten_mode = "flatten_at_root"
}


`, name)
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step0(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "off"

        rule = "(http.host eq \"video.example.com\")"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}"

        error_pages_redirect = [
        {
            target_url = "https://example.com/foo/bar"
            status_code = "500"
        },
        {
            target_url = "https://example.com/foo"
            status_code = "400"
        },
        {
            target_url = "https://example.com/test"
            status_code = "503"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step1(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "off"

        rule = "(http.host eq \"video.example.com\")"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}-mod"

        error_pages_redirect = [
        {
            target_url = "https://example.com/foo/bar"
            status_code = "500"
        },
        {
            target_url = "https://example.com/foo"
            status_code = "400"
        },
        {
            target_url = "https://example.com/test"
            status_code = "503"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step2(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "on"

        rule = "(http.host eq \"video.example.com\")"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}-mod"

        error_pages_redirect = [
        {
            target_url = "https://example.com/foo/bar"
            status_code = "500"
        },
        {
            target_url = "https://example.com/foo"
            status_code = "400"
        },
        {
            target_url = "https://example.com/test"
            status_code = "503"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step3(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "on"

        rule = "(http.request.uri eq \"/content?page=1234\")"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}-mod"

        error_pages_redirect = [
        {
            target_url = "https://example.com/foo/bar"
            status_code = "500"
        },
        {
            target_url = "https://example.com/foo"
            status_code = "400"
        },
        {
            target_url = "https://example.com/test"
            status_code = "503"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step4(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "on"

        rule = "(http.request.uri eq \"/content?page=1234\")"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}-mod"

        error_pages_redirect = [
        {
            target_url = "https://example.com/test"
            status_code = "404"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step5(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "on"

        rule = "(http.request.uri eq \"/content?page=1234\")"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}-mod"

        error_pages_redirect = [
        {
            target_url = "https://example.com/foo"
            status_code = "400"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step6(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "on"

        rule = "(http.request.uri eq \"/content?page=1234\")"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}-mod"

        error_pages_redirect = [
        {
            target_url = "https://example.com/foo/bar"
            status_code = "500"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step7(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "on"

        rule = "(http.request.uri eq \"/content?page=1234\")"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}-mod"

        error_pages_redirect = [
        {
            target_url = "https://example.com/foo/bar"
            status_code = "500"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

func testAccAliCloudEsaErrorPagesRedirectRule_basic12504_step8(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "alicloud_esa_error_pages_redirect_rule" "default" {
        site_id = "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}"

        rule_enable = "off"

        rule = "true"

        sequence = 1

        site_version = 0

        rule_name = "${var.name}"

        error_pages_redirect = [
        {
            target_url = "https://example.com/foo/bar"
            status_code = "500"
        },
    ]

}
`, name, testAccAliCloudEsaErrorPagesRedirectRule_basic12504_base(name))
}

// Test Esa ErrorPagesRedirectRule. <<< Resource test cases, automatically generated.
