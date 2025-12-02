// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Esa HttpIncomingResponseHeaderModificationRule. >>> Resource test cases, automatically generated.
// Case HttpIncomingResponseHeaderModificationRule_test 11932
func TestAccAliCloudEsaHttpIncomingResponseHeaderModificationRule_basic11932(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_http_incoming_response_header_modification_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaHttpIncomingResponseHeaderModificationRuleMap11932)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaHttpIncomingResponseHeaderModificationRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaHttpIncomingResponseHeaderModificationRuleBasicDependence11932)
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
					"site_id":     "${alicloud_esa_site.resource_Site_HttpIncomingResponseHeaderModificationRule_test.id}",
					"rule_enable": "on",
					"response_header_modification": []map[string]interface{}{
						{
							"type":      "static",
							"value":     "add",
							"operation": "add",
							"name":      "testadd",
						},
						{
							"type":      "static",
							"operation": "del",
							"name":      "testdel",
						},
						{
							"type":      "static",
							"value":     "modify",
							"operation": "modify",
							"name":      "testmodify",
						},
					},
					"rule":         "(http.host eq \\\"video.example.com\\\")",
					"sequence":     "1",
					"site_version": "0",
					"rule_name":    "testResponseHeader",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                        CHECKSET,
						"rule_enable":                    "on",
						"response_header_modification.#": "3",
						"rule":                           "(http.host eq \"video.example.com\")",
						"sequence":                       "1",
						"site_version":                   "0",
						"rule_name":                      "testResponseHeader",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "testResponseHeader_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "testResponseHeader_modify",
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
					"response_header_modification": []map[string]interface{}{
						{
							"type":      "dynamic",
							"value":     "ip.geoip.country",
							"operation": "add",
							"name":      "x-ip-country",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"response_header_modification.#": "1",
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

var AlicloudEsaHttpIncomingResponseHeaderModificationRuleMap11932 = map[string]string{
	"config_id": CHECKSET,
}

func AlicloudEsaHttpIncomingResponseHeaderModificationRuleBasicDependence11932(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_esa_rate_plan_instance" "resource_HttpIncomingResponseHeaderModificationRule_test" {
  type         = "NS"
  auto_renew   = false
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "basic"
}

resource "alicloud_esa_site" "resource_Site_HttpIncomingResponseHeaderModificationRule_test" {
  site_name   = "${var.name}.com"
  instance_id = alicloud_esa_rate_plan_instance.resource_HttpIncomingResponseHeaderModificationRule_test.id
  coverage    = "overseas"
  access_type = "NS"
}


`, name)
}

// Test Esa HttpIncomingResponseHeaderModificationRule. <<< Resource test cases, automatically generated.
