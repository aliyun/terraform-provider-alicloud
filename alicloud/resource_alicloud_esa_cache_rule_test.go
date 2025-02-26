package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA CacheRule. >>> Resource test cases, automatically generated.
// Case resource_cacherule_test
func TestAccAliCloudESACacheRuleresource_cacherule_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_cache_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudESACacheRuleresource_cacherule_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCacheRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESACacheRule%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESACacheRuleresource_cacherule_testBasicDependence)
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
					"edge_cache_ttl":              "300",
					"site_id":                     "${data.alicloud_esa_sites.default.sites.0.id}",
					"edge_cache_mode":             "follow_origin",
					"edge_status_code_cache_ttl":  "300",
					"browser_cache_mode":          "no_cache",
					"user_geo":                    "off",
					"browser_cache_ttl":           "300",
					"include_cookie":              "cookie_exapmle",
					"rule_name":                   "rule_example",
					"rule_enable":                 "off",
					"query_string_mode":           "ignore_all",
					"query_string":                "example",
					"bypass_cache":                "cache_all",
					"check_presence_header":       "headername",
					"sort_query_string_for_cache": "off",
					"check_presence_cookie":       "cookiename",
					"user_device_type":            "off",
					"cache_reserve_eligibility":   "bypass_cache_reserve",
					"additional_cacheable_ports":  "2053",
					"rule":                        "http.host eq \\\"video.example.com\\\"",
					"user_language":               "off",
					"serve_stale":                 "off",
					"cache_deception_armor":       "off",
					"site_version":                "0",
					"include_header":              "example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"edge_cache_ttl":              "400",
					"rule_enable":                 "on",
					"query_string_mode":           "exclude_query_string",
					"query_string":                "shili",
					"edge_cache_mode":             "no_cache",
					"bypass_cache":                "bypass_all",
					"edge_status_code_cache_ttl":  "400",
					"check_presence_header":       "headernamee",
					"sort_query_string_for_cache": "on",
					"check_presence_cookie":       "cookienamee",
					"browser_cache_mode":          "follow_origin",
					"user_device_type":            "on",
					"user_geo":                    "on",
					"browser_cache_ttl":           "400",
					"include_cookie":              "cookie_shili",
					"cache_reserve_eligibility":   "eligible_for_cache_reserve",
					"additional_cacheable_ports":  "2052",
					"rule":                        "http.host eq \\\"videoo.example.com\\\"",
					"user_language":               "on",
					"serve_stale":                 "on",
					"cache_deception_armor":       "on",
					"rule_name":                   "rule_shili",
					"include_header":              "shili",
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

var AliCloudESACacheRuleresource_cacherule_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESACacheRuleresource_cacherule_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

`, name)
}

// Test ESA CacheRule. <<< Resource test cases, automatically generated.
