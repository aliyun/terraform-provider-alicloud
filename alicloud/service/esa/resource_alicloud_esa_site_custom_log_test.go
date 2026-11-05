package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const siteCustomLogResourceId = "alicloud_esa_site_custom_log.default"

// TestAccAliCloudEsaSiteCustomLog_basic covers the Create→Read→Update lifecycle,
// Import and the no-op Delete behavior. SiteCustomLog is a site-level singleton
// keyed by SiteId; it reuses an existing ESA site from the data source.
func TestAccAliCloudEsaSiteCustomLog_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// ESA has no DeleteSiteCustomLog API; the resource Delete is a no-op and
		// the server-side configuration remains, so there is no destroy to assert.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEsaSiteCustomLogBasicConfig(
					[]string{"x-acs-request-id", "x-real-ip"},
					[]string{"X_Swift_Error", "Server"},
					[]string{"sessionId", "userId"},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(siteCustomLogResourceId, "id"),
					resource.TestCheckResourceAttrSet(siteCustomLogResourceId, "site_id"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "request_headers.#", "2"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "request_headers.0", "x-acs-request-id"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "request_headers.1", "x-real-ip"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "response_headers.#", "2"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "response_headers.0", "X_Swift_Error"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "response_headers.1", "Server"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "cookies.#", "2"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "cookies.0", "sessionId"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "cookies.1", "userId"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "is_exist", "true"),
					resource.TestCheckResourceAttrSet(siteCustomLogResourceId, "config_id"),
				),
			},
			{
				// Update: change the logged custom fields and refresh with no drift.
				Config: testAccEsaSiteCustomLogBasicConfig(
					[]string{"x-acs-trace-id"},
					[]string{"X_Swift_Error"},
					[]string{"sessionId"},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "request_headers.#", "1"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "request_headers.0", "x-acs-trace-id"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "response_headers.#", "1"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "response_headers.0", "X_Swift_Error"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "cookies.#", "1"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "cookies.0", "sessionId"),
					resource.TestCheckResourceAttr(siteCustomLogResourceId, "is_exist", "true"),
				),
			},
			{
				ResourceName:      siteCustomLogResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEsaSiteCustomLogBasicConfig(requestHeaders, responseHeaders, cookies []string) string {
	return fmt.Sprintf(`
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

resource "alicloud_esa_site_custom_log" "default" {
  site_id = data.alicloud_esa_sites.default.sites.0.id

  request_headers  = %s
  response_headers = %s
  cookies          = %s
}
`, listLiteral(requestHeaders), listLiteral(responseHeaders), listLiteral(cookies))
}

// listLiteral renders a Go []string as an HCL list literal.
func listLiteral(values []string) string {
	out := "["
	for i, v := range values {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%q", v)
	}
	return out + "]"
}
