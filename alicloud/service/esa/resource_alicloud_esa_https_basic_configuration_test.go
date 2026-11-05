package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA HttpsBasicConfiguration. >>> Resource test cases, automatically generated.
// Case resource_HttpsBasicConfiguration_set_globle_test
func TestAccAliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_https_basic_configuration.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaHttpsBasicConfiguration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAHttpsBasicConfiguration%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":           "${data.alicloud_esa_sites.default.sites.0.id}",
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
					"site_id":           "${data.alicloud_esa_sites.default.sites.0.id}",
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
				Config: testAccConfig(map[string]interface{}{
					"sequence": "1",
				}),
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
	"id": acctest.CHECKSET,
}

func AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_globle_testBasicDependence(name string) string {
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

// Case resource_HttpsBasicConfiguration_set_test
func TestAccAliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_https_basic_configuration.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaHttpsBasicConfiguration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAHttpsBasicConfiguration%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.id}",
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
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.id}",
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
	"id": acctest.CHECKSET,
}

func AliCloudESAHttpsBasicConfigurationresource_HttpsBasicConfiguration_set_testBasicDependence(name string) string {
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

// Test ESA HttpsBasicConfiguration. <<< Resource test cases, automatically generated.
