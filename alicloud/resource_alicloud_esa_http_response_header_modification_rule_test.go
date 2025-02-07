package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA HttpResponseHeaderModificationRule. >>> Resource test cases, automatically generated.
// Case httpResponseHeaderModificationRule_test
func TestAccAliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_http_response_header_modification_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaHttpResponseHeaderModificationRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAHttpResponseHeaderModificationRule%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${alicloud_esa_site.resource_Site_HttpResponseHeaderModificationRule_test.id}",
					"rule_enable": "on",
					"response_header_modification": []map[string]interface{}{

						{
							"value":     "add",
							"operation": "add",
							"name":      "testadd",
						},

						{
							"operation": "del",
							"name":      "testdel",
						},

						{
							"value":     "modify",
							"operation": "modify",
							"name":      "testmodify",
						},
					},
					"rule":         "(http.host eq \\\"video.example.com\\\")",
					"site_version": "0",
					"rule_name":    "testResponseHeader",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "testResponseHeader_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "(http.request.uri eq \\\"/content?page=1234\\\")",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"response_header_modification": []map[string]interface{}{

						{
							"value":     "add1",
							"operation": "add",
							"name":      "testadd1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
					"rule":        "(http.host eq \\\"api.example.com\\\")",
					"rule_name":   "test_httpResponseHeader_last",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_HttpResponseHeaderModificationRule_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_HttpResponseHeaderModificationRule_test" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_HttpResponseHeaderModificationRule_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA HttpResponseHeaderModificationRule. <<< Resource test cases, automatically generated.
