package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA HttpRequestHeaderModificationRule. >>> Resource test cases, automatically generated.
// Case httprequestheadermodificationrule_test
func TestAccAliCloudESAHttpRequestHeaderModificationRulehttprequestheadermodificationrule_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_http_request_header_modification_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudESAHttpRequestHeaderModificationRulehttprequestheadermodificationrule_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaHttpRequestHeaderModificationRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAHttpRequestHeaderModificationRule%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAHttpRequestHeaderModificationRulehttprequestheadermodificationrule_testBasicDependence)
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
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"rule_enable":  "on",
					"rule":         "(http.host eq \\\"video.example.com\\\")",
					"site_version": "0",
					"rule_name":    "test",
					"request_header_modification": []map[string]interface{}{

						{
							"value":     "add",
							"operation": "add",
							"name":      "testadd",
							"type":      "static",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "test_modify",
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
					"sequence": "1",
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
					"request_header_modification": []map[string]interface{}{

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
					"request_header_modification": []map[string]interface{}{

						{
							"operation": "del",
							"name":      "testdel1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"request_header_modification": []map[string]interface{}{

						{
							"value":     "modify1",
							"operation": "modify",
							"name":      "testmodify1",
						},
					},
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

var AliCloudESAHttpRequestHeaderModificationRulehttprequestheadermodificationrule_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAHttpRequestHeaderModificationRulehttprequestheadermodificationrule_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name = "gositecdn.cn"
}

`, name)
}

// Test ESA HttpRequestHeaderModificationRule. <<< Resource test cases, automatically generated.
