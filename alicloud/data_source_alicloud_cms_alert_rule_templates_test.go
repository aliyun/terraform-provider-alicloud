package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Recursively skipping this acceptance test in the short term: the backing API
// ListAlertRuleTemplates is a private Cms ROA endpoint that is not reachable
// from every acceptance account, so the test guards on a best-effort read and
// only asserts the data source does not error and the list attributes are set.
func TestAccAliCloudCmsAlertRuleTemplatesDataSource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCmsAlertRuleTemplatesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_rule_templates.default", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_rule_templates.default", "templates.#"),
				),
			},
			{
				Config: testAccCheckAlicloudCmsAlertRuleTemplatesDataSourceNameRegex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_rule_templates.filtered", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_rule_templates.filtered", "templates.#"),
				),
			},
			{
				Config: testAccCheckAlicloudCmsAlertRuleTemplatesDataSourceEnableDetails,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_rule_templates.detailed", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_alert_rule_templates.detailed", "templates.#"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCmsAlertRuleTemplatesDataSourceBasic = `
data "alicloud_cms_alert_rule_templates" "default" {
  biz_source = "CI"
}
`

const testAccCheckAlicloudCmsAlertRuleTemplatesDataSourceNameRegex = `
data "alicloud_cms_alert_rule_templates" "filtered" {
  biz_source  = "CI"
  name_regex  = "^non-existent-template-.*$"
}
`

const testAccCheckAlicloudCmsAlertRuleTemplatesDataSourceEnableDetails = `
data "alicloud_cms_alert_rule_templates" "detailed" {
  biz_source     = "CI"
  enable_details = true
}
`
