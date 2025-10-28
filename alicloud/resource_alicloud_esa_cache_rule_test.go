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
					"edge_cache_ttl": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"edge_cache_ttl": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sequence": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sequence": "1",
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
					"query_string_mode": "exclude_query_string",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_string_mode": "exclude_query_string",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_string": "shili",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_string": "shili",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"edge_cache_mode": "no_cache",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"edge_cache_mode": "no_cache",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bypass_cache": "bypass_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bypass_cache": "bypass_all",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"edge_status_code_cache_ttl": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"edge_status_code_cache_ttl": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_presence_header": "headernamee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_presence_header": "headernamee",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sort_query_string_for_cache": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sort_query_string_for_cache": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_presence_cookie": "cookienamee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_presence_cookie": "cookienamee",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"browser_cache_mode": "follow_origin",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"browser_cache_mode": "follow_origin",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_device_type": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_device_type": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_geo": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_geo": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"browser_cache_ttl": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"browser_cache_ttl": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"include_cookie": "cookie_shili",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"include_cookie": "cookie_shili",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cache_reserve_eligibility": "eligible_for_cache_reserve",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cache_reserve_eligibility": "eligible_for_cache_reserve",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"additional_cacheable_ports": "2052",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"additional_cacheable_ports": "2052",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "http.host eq \\\"videoo.example.com\\\"",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_language": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_language": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serve_stale": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serve_stale": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cache_deception_armor": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cache_deception_armor": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "rule_shili",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "rule_shili",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"include_header": "shili",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"include_header": "shili",
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
