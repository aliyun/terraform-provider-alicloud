package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA HttpsBasicConfiguration. >>> Resource test cases, automatically generated.
// Case resource_HttpsBasicConfiguration_set_globle_test
func TestAccAliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_https_basic_configuration.default"
	ra := resourceAttrInit(resourceId, AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaHttpsBasicConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAHttpsBasicConfiguration%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_testBasicDependence)
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
					"site_id":           "${alicloud_esa_site.resource_HttpBasicConfiguration_set_global_test.id}",
					"ciphersuite":       "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
					"rule_enable":       "on",
					"https":             "on",
					"http3":             "on",
					"http2":             "on",
					"tls10":             "on",
					"tls11":             "on",
					"tls12":             "on",
					"tls13":             "on",
					"ciphersuite_group": "all",
					"rule":              "true",
					"rule_name":         "test_global1",
					"ocsp_stapling":     "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":           "${alicloud_esa_site.resource_HttpBasicConfiguration_set_global_test.id}",
					"ciphersuite":       "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
					"rule_enable":       "off",
					"https":             "off",
					"http3":             "off",
					"http2":             "off",
					"tls10":             "off",
					"tls11":             "off",
					"tls12":             "off",
					"tls13":             "off",
					"ciphersuite_group": "custom",
					"rule":              "true",
					"rule_name":         "test_global1",
					"ocsp_stapling":     "off",
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

var AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_set_globle_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_HttpBasicConfiguration_set_global_test" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_set_globle_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_HttpsBasicConfiguration_set_test
func TestAccAliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_https_basic_configuration.default"
	ra := resourceAttrInit(resourceId, AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaHttpsBasicConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAHttpsBasicConfiguration%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_testBasicDependence)
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
					"site_id":     "${alicloud_esa_site.resource_HttpBasicConfiguration_set_test.id}",
					"rule_enable": "on",
					"https":       "on",
					"rule":        "true",
					"rule_name":   "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${alicloud_esa_site.resource_HttpBasicConfiguration_set_test.id}",
					"rule_enable": "on",
					"https":       "off",
					"rule":        "true",
					"rule_name":   "test2",
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

var AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_rate_plan_instance" "resource_RatePlanInstance_set_test" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_HttpBasicConfiguration_set_test" {
  site_name   = "gositecdn.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_RatePlanInstance_set_test.id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA HttpsBasicConfiguration. <<< Resource test cases, automatically generated.
