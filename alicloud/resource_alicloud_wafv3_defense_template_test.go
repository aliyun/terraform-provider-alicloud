package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Wafv3 DefenseTemplate. >>> Resource test cases, automatically generated.
// Case 接入terraform 5993
func TestAccAliCloudWafv3DefenseTemplate_basic5993(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.WAFV3SupportRegions)
	resourceId := "alicloud_wafv3_defense_template.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseTemplateMap5993)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%swafv3defensetemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseTemplateBasicDependence5993)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckForCleanUpInstances(t, string(connectivity.Hangzhou), "waf", "waf", "waf", "waf")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                             "0",
					"instance_id":                        "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_template_name":              name,
					"template_type":                      "user_custom",
					"template_origin":                    "custom",
					"defense_scene":                      "antiscan",
					"resource_manager_resource_group_id": "test",
					"description":                        "update_template",
					"resources":                          []string{"${alicloud_wafv3_domain.default1.domain_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                             "0",
						"instance_id":                        CHECKSET,
						"defense_template_name":              name,
						"template_type":                      "user_custom",
						"template_origin":                    "custom",
						"defense_scene":                      "antiscan",
						"resource_manager_resource_group_id": "test",
						"description":                        "update_template",
						"resources.#":                        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "createTestDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "createTestDescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []string{"${alicloud_wafv3_domain.default2.domain_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"defense_template_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"defense_template_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "update_template",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "update_template",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                             "1",
					"instance_id":                        "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_template_name":              name + "_update",
					"template_type":                      "user_custom",
					"template_origin":                    "custom",
					"defense_scene":                      "antiscan",
					"resource_manager_resource_group_id": "test",
					"description":                        "createTestDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                             "1",
						"instance_id":                        CHECKSET,
						"defense_template_name":              name + "_update",
						"template_type":                      "user_custom",
						"template_origin":                    "custom",
						"defense_scene":                      "antiscan",
						"resource_manager_resource_group_id": "test",
						"description":                        "createTestDescription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_manager_resource_group_id"},
			},
		},
	})
}

var AlicloudWafv3DefenseTemplateMap5993 = map[string]string{
	"defense_template_id": CHECKSET,
}

func AlicloudWafv3DefenseTemplateBasicDependence5993(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_domain" "default1" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  listen {
    protection_resource = "share"
    http_ports = [
      "81",
      "82",
      "83"
    ]
    https_ports = [

    ]
    xff_header_mode = "2"
    xff_headers = [
      "testa",
      "testb",
      "testc"
    ]
    custom_ciphers = [

    ]
    ipv6_enabled = "true"
  }

  redirect {
    keepalive_timeout = "15"
    backends = [
      "1.1.1.1",
      "3.3.3.3",
      "2.2.2.2"
    ]
    write_timeout      = "5"
    keepalive_requests = "1000"
    request_headers {
      key   = "testkey1"
      value = "testValue1"
    }
    request_headers {
      key   = "key1"
      value = "value1"
    }
    request_headers {
      key   = "key22"
      value = "value22"
    }

    loadbalance        = "iphash"
    focus_http_backend = "false"
    sni_enabled        = "false"
    connect_timeout    = "5"
    read_timeout       = "5"
    keepalive          = "true"
    retry              = "true"
  }

  domain                             = "zctest_250744.wafqax.top"
  access_type                        = "share"
}

resource "alicloud_wafv3_domain" "default2" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  listen {
    protection_resource = "share"
    http_ports = [
      "81",
      "82",
      "83"
    ]
    https_ports = [

    ]
    xff_header_mode = "2"
    xff_headers = [
      "testa",
      "testb",
      "testc"
    ]
    custom_ciphers = [

    ]
    ipv6_enabled = "true"
  }

  redirect {
    keepalive_timeout = "15"
    backends = [
      "1.1.1.1",
      "3.3.3.3",
      "2.2.2.2"
    ]
    write_timeout      = "5"
    keepalive_requests = "1000"
    request_headers {
      key   = "testkey1"
      value = "testValue1"
    }
    request_headers {
      key   = "key1"
      value = "value1"
    }
    request_headers {
      key   = "key22"
      value = "value22"
    }

    loadbalance        = "iphash"
    focus_http_backend = "false"
    sni_enabled        = "false"
    connect_timeout    = "5"
    read_timeout       = "5"
    keepalive          = "true"
    retry              = "true"
  }

  domain                             = "zctest_250745.wafqax.top"
  access_type                        = "share"
}

`, name)
}

// Case 接入terraform 5993  twin
func TestAccAliCloudWafv3DefenseTemplate_basic5993_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.WAFV3SupportRegions)
	resourceId := "alicloud_wafv3_defense_template.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseTemplateMap5993)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%swafv3defensetemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseTemplateBasicDependence5993)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                             "0",
					"instance_id":                        "${data.alicloud_wafv3_instances.default.ids.0}",
					"defense_template_name":              name,
					"template_type":                      "user_custom",
					"template_origin":                    "custom",
					"defense_scene":                      "antiscan",
					"resource_manager_resource_group_id": "test",
					"description":                        "update_template",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                             "0",
						"instance_id":                        CHECKSET,
						"defense_template_name":              name,
						"template_type":                      "user_custom",
						"template_origin":                    "custom",
						"defense_scene":                      "antiscan",
						"resource_manager_resource_group_id": "test",
						"description":                        "update_template",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_manager_resource_group_id"},
			},
		},
	})
}

// Test Wafv3 DefenseTemplate. <<< Resource test cases, automatically generated.
