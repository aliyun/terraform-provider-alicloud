// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Esa ErrorPagesRedirectRule. >>> Resource test cases, automatically generated.
// Case resource_ErrorPagesRedirectRule_test 12504
func TestAccAliCloudEsaErrorPagesRedirectRule_basic12504(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_error_pages_redirect_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaErrorPagesRedirectRuleMap12504)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaErrorPagesRedirectRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaErrorPagesRedirectRuleBasicDependence12504)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${alicloud_esa_site.resource_Site_test_ErrorPagesRedirectRule.id}",
					"rule_enable":  "off",
					"rule":         "(http.host eq \\\"video.example.com\\\")",
					"sequence":     "1",
					"site_version": "0",
					"rule_name":    "rule_example",
					"error_pages_redirect": []map[string]interface{}{
						{
							"target_url":  "https://example.com/foo/bar",
							"status_code": "500",
						},
						{
							"target_url":  "https://example.com/foo",
							"status_code": "400",
						},
						{
							"target_url":  "https://example.com/test",
							"status_code": "503",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                CHECKSET,
						"rule_enable":            "off",
						"rule":                   "(http.host eq \"video.example.com\")",
						"sequence":               "1",
						"site_version":           "0",
						"rule_name":              "rule_example",
						"error_pages_redirect.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "test_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "test_modify",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "(http.request.uri eq \\\"/content?page=1234\\\")",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": "(http.request.uri eq \"/content?page=1234\")",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"error_pages_redirect": []map[string]interface{}{
						{
							"target_url":  "https://example.com/test",
							"status_code": "404",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"error_pages_redirect.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"error_pages_redirect": []map[string]interface{}{
						{
							"target_url":  "https://example.com/foo",
							"status_code": "400",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"error_pages_redirect.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"error_pages_redirect": []map[string]interface{}{
						{
							"target_url":  "https://example.com/foo/bar",
							"status_code": "500",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"error_pages_redirect.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "off",
					"rule":        "true",
					"rule_name":   "rule_example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "off",
						"rule":        CHECKSET,
						"rule_name":   "rule_example",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEsaErrorPagesRedirectRuleMap12504 = map[string]string{
	"config_id":   CHECKSET,
	"config_type": CHECKSET,
}

func AlicloudEsaErrorPagesRedirectRuleBasicDependence12504(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
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
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_test_ErrorPagesRedirectRule.id
  coverage    = "overseas"
  access_type = "NS"
}


`, name)
}

// Test Esa ErrorPagesRedirectRule. <<< Resource test cases, automatically generated.
