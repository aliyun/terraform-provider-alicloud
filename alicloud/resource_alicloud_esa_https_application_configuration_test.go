package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ESA HttpsApplicationConfiguration. >>> Resource test cases, automatically generated.
// Case resource_HttpsApplicationConfiguration_test
func TestAccAliCloudESAHttpsApplicationConfigurationresource_HttpsApplicationConfiguration_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_https_application_configuration.default"
	ra := resourceAttrInit(resourceId, AliCloudESAHttpsApplicationConfigurationresource_HttpsApplicationConfiguration_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaHttpsApplicationConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAHttpsApplicationConfiguration%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAHttpsApplicationConfigurationresource_HttpsApplicationConfiguration_testBasicDependence)
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
					"site_id":                 "${data.alicloud_esa_sites.default.sites.0.id}",
					"hsts_include_subdomains": "off",
					"alt_svc_ma":              "86400",
					"rule_enable":             "off",
					"https_force_code":        "301",
					"alt_svc":                 "off",
					"hsts":                    "off",
					"hsts_preload":            "off",
					"hsts_max_age":            "31536000",
					"alt_svc_persist":         "off",
					"alt_svc_clear":           "off",
					"https_force":             "off",
					"rule":                    "http.host eq \\\"video.example.com\\\"",
					"site_version":            "0",
					"rule_name":               "rule_example",
					"https_no_sni_deny":       "off",
					"https_sni_verify":        "off",
					"https_sni_whitelist":     "example.com test.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hsts_include_subdomains": "on",
					"alt_svc_ma":              "172800",
					"rule_enable":             "on",
					"https_force_code":        "301",
					"alt_svc":                 "on",
					"hsts":                    "on",
					"hsts_preload":            "on",
					"hsts_max_age":            "15768000",
					"alt_svc_persist":         "on",
					"alt_svc_clear":           "on",
					"https_force":             "on",
					"rule":                    "http.host eq \\\"videoo.example.com\\\"",
					"rule_name":               "rule_viedoo",
					"https_no_sni_deny":       "on",
					"https_sni_verify":        "on",
					"https_sni_whitelist":     "videoo.example.com secure.videoo.example.com",
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
					"rule": "http.host eq \\\"image.example.com\\\"",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "rule_image",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_force": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_force_code": "302",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alt_svc": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alt_svc_clear": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alt_svc_persist": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alt_svc_ma": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hsts": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hsts_max_age": "31536000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hsts_include_subdomains": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hsts_preload": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_no_sni_deny": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_sni_verify": "on",
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
					"https_sni_whitelist": "secure.example.com api.example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hsts_include_subdomains": "on",
					"alt_svc_ma":              "172800",
					"rule_enable":             "off",
					"https_force_code":        "301",
					"alt_svc":                 "on",
					"hsts":                    "on",
					"hsts_preload":            "on",
					"hsts_max_age":            "15768000",
					"alt_svc_persist":         "on",
					"alt_svc_clear":           "on",
					"https_force":             "on",
					"rule":                    "http.host eq \\\"video.example.com\\\"",
					"rule_name":               "rule_example",
					"https_no_sni_deny":       "off",
					"https_sni_verify":        "off",
					"https_sni_whitelist":     "video.example.com secure.video.example.com",
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

var AliCloudESAHttpsApplicationConfigurationresource_HttpsApplicationConfiguration_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAHttpsApplicationConfigurationresource_HttpsApplicationConfiguration_testBasicDependence(name string) string {
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

// Test ESA HttpsApplicationConfiguration. <<< Resource test cases, automatically generated.
