package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ESA CompressionRule. >>> Resource test cases, automatically generated.
// Case resource_CompressionRule_test
func TestAccAliCloudESACompressionRuleresource_CompressionRule_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_compression_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudESACompressionRuleresource_CompressionRule_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCompressionRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESACompressionRule%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESACompressionRuleresource_CompressionRule_testBasicDependence)
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
					"site_id":      "${alicloud_esa_site.resource_Site_test_Compression.id}",
					"zstd":         "off",
					"rule_enable":  "off",
					"gzip":         "off",
					"brotli":       "off",
					"rule":         "http.host eq \\\"video.example.com\\\"",
					"site_version": "0",
					"rule_name":    "rule_example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zstd":        "off",
					"rule_enable": "on",
					"gzip":        "on",
					"brotli":      "off",
					"rule":        "http.host eq \\\"videoo.example.com\\\"",
					"rule_name":   "rule_viedoo",
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
					"gzip": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"brotli": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zstd": "on",
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
					testAccCheck(map[string]string{
						"sequence": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zstd": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zstd": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zstd":        "off",
					"rule_enable": "off",
					"gzip":        "off",
					"brotli":      "off",
					"rule":        "http.host eq \\\"video.example.com\\\"",
					"rule_name":   "rule_example",
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
				// Removing optional attributes from the config must not send empty
				// strings: the API rejects a present-but-empty parameter with
				// InvalidRuleEnable/MissingParameter.
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": REMOVEKEY,
					"rule_name":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
					"rule_name":   "rule_readded",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "on",
						"rule_name":   "rule_readded",
					}),
				),
			},
			{
				// Removing the whole triple at once: the config type cannot migrate
				// back in place, so the values follow the state, nothing is sent and
				// the plan stays empty - where the old code sent three empty strings
				// and failed the update.
				Config: testAccConfig(map[string]interface{}{
					"rule":        REMOVEKEY,
					"rule_name":   REMOVEKEY,
					"rule_enable": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": CHECKSET,
					}),
				),
			},
			{
				// Re-adding `rule` with new content while the state still holds the
				// old rule is an in-place update, not a replacement.
				Config: testAccConfig(map[string]interface{}{
					"rule":      "http.host eq \\\"image.example.com\\\"",
					"rule_name": "rule_image2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule":      "http.host eq \"image.example.com\"",
						"rule_name": "rule_image2",
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

var AliCloudESACompressionRuleresource_CompressionRule_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESACompressionRuleresource_CompressionRule_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_test_Compression" {
  site_name   = "compression.alicdn-test.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "domestic"
  access_type = "NS"
  flatten_mode = "flatten_at_root"
}

`, name)
}

// TestAccAliCloudESACompressionRule_globalConfig covers a global compression
// config (the rule/rule_name/rule_enable triple not set): in-place updates, a
// partial triple member accepted and persisted by the API, removing it again
// without sending an empty string, and converting the global config to a
// rule-type config by adding `rule`.
func TestAccAliCloudESACompressionRule_globalConfig(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_compression_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudESACompressionRuleresource_CompressionRule_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCompressionRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESACompressionRuleG%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESACompressionRule_globalTestBasicDependence)
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
					"site_id": "${alicloud_esa_site.resource_Site_test_Compression.id}",
					"gzip":    "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gzip": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip": "off",
					}),
				),
			},
			{
				// A partial triple member on a global config is accepted and
				// persisted by the API; the state mirrors it.
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
				// Removing it again must follow the state value instead of
				// sending an empty RuleEnable, which the API rejects with
				// InvalidRuleEnable.
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "on",
					}),
				),
			},
			{
				// Adding `rule` to the global config replaces the resource with a
				// rule-type config: the config type cannot migrate in place.
				Config: testAccConfig(map[string]interface{}{
					"rule":        "http.host eq \\\"video.example.com\\\"",
					"rule_name":   "rule_example",
					"rule_enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule":        CHECKSET,
						"rule_name":   "rule_example",
						"rule_enable": "on",
					}),
				),
			},
			{
				// Removing `gzip` from the config follows the state value kept on
				// the server instead of sending an empty string.
				Config: testAccConfig(map[string]interface{}{
					"gzip": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip": "off",
					}),
				),
			},
			{
				// The match-all form of `rule` (literal true) is an in-place rule
				// change, not a replacement.
				Config: testAccConfig(map[string]interface{}{
					"rule": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": "true",
					}),
				),
			},
		},
	})
}

func AliCloudESACompressionRule_globalTestBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_test_Compression" {
  site_name   = "compression-global.alicdn-test.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "domestic"
  access_type = "NS"
  flatten_mode = "flatten_at_root"
}

`, name)
}

// Test ESA CompressionRule. <<< Resource test cases, automatically generated.
