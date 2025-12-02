// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Esa HttpIncomingRequestHeaderModificationRule. >>> Resource test cases, automatically generated.
// Case HttpIncomingRequestHeadermodificationrule_test 11933
func TestAccAliCloudEsaHttpIncomingRequestHeaderModificationRule_basic11933(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_http_incoming_request_header_modification_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaHttpIncomingRequestHeaderModificationRuleMap11933)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaHttpIncomingRequestHeaderModificationRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaHttpIncomingRequestHeaderModificationRuleBasicDependence11933)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${alicloud_esa_site.resource_Site_HttpIncomingRequestHeaderModificationRule_test.id}",
					"rule_enable":  "on",
					"rule":         "(http.host eq \\\"video.example.com\\\")",
					"sequence":     "1",
					"site_version": "0",
					"rule_name":    "test",
					"request_header_modification": []map[string]interface{}{
						{
							"type":      "static",
							"value":     "add",
							"operation": "add",
							"name":      "testadd",
						},
						{
							"type":      "dynamic",
							"value":     "ip.geoip.country",
							"operation": "modify",
							"name":      "testmodify",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                       CHECKSET,
						"rule_enable":                   "on",
						"rule":                          "(http.host eq \"video.example.com\")",
						"sequence":                      "1",
						"site_version":                  "0",
						"rule_name":                     "test",
						"request_header_modification.#": "2",
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
					"rule_enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "off",
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
					"request_header_modification": []map[string]interface{}{
						{
							"type":      "dynamic",
							"value":     "ip.geoip.country",
							"operation": "add",
							"name":      "testadd1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_header_modification.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"request_header_modification": []map[string]interface{}{
						{
							"type":      "static",
							"value":     "modify1",
							"operation": "modify",
							"name":      "testmodify1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_header_modification.#": "1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEsaHttpIncomingRequestHeaderModificationRuleMap11933 = map[string]string{
	"config_id": CHECKSET,
}

func AlicloudEsaHttpIncomingRequestHeaderModificationRuleBasicDependence11933(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_esa_rate_plan_instance" "resource_HttpIncomingRequestHeaderModificationRule_test" {
  type         = "NS"
  auto_renew   = false
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_HttpIncomingRequestHeaderModificationRule_test" {
  site_name   = "${var.name}.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_HttpIncomingRequestHeaderModificationRule_test.id
  coverage    = "overseas"
  access_type = "NS"
}


`, name)
}

// Test Esa HttpIncomingRequestHeaderModificationRule. <<< Resource test cases, automatically generated.
