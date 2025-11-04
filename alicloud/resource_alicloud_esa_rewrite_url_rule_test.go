package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA RewriteUrlRule. >>> Resource test cases, automatically generated.
// Case resource_RewriteUrlRule_test
func TestAccAliCloudESARewriteUrlRuleresource_RewriteUrlRule_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_rewrite_url_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudESARewriteUrlRuleresource_RewriteUrlRule_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRewriteUrlRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARewriteUrlRule%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARewriteUrlRuleresource_RewriteUrlRule_testBasicDependence)
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
					"site_id":                   "${data.alicloud_esa_sites.default.sites.0.id}",
					"rewrite_uri_type":          "static",
					"rule_enable":               "on",
					"rewrite_query_string_type": "static",
					"query_string":              "example=123",
					"rule":                      "http.host eq \\\"video.example.com\\\"",
					"site_version":              "0",
					"uri":                       "/image/example.jpg",
					"rule_name":                 "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "http.host eq \\\"video.example.com\\\"",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "重写URL规则名称示例",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rewrite_uri_type": "static",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"uri": "/rewritten/target-uri",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rewrite_query_string_type": "static",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_string": "重写后的查询字符串示例",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rewrite_uri_type":          "static",
					"rule_enable":               "off",
					"rewrite_query_string_type": "static",
					"query_string":              "新的查询字符串示例",
					"rule":                      "http.host eq \\\"video.change.com\\\"",
					"uri":                       "/rewritten/new-target-uri",
					"rule_name":                 "新规则名称示例",
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

var AliCloudESARewriteUrlRuleresource_RewriteUrlRule_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARewriteUrlRuleresource_RewriteUrlRule_testBasicDependence(name string) string {
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

// Test ESA RewriteUrlRule. <<< Resource test cases, automatically generated.
