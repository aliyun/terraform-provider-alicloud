package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Wafv3 DefenseTemplate. >>> Resource test cases, automatically generated.
// Case 防护模板-防护模板和防护组关联关系默认防护ip_blacklist 11649
func TestAccAliCloudWafv3DefenseTemplate_basic11649(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_template.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseTemplateMap11649)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseTemplateBasicDependence11649)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":      "1",
					"description": "testTP",
					"resource_groups": []string{
						"${alicloud_wafv3_defense_resource_group.defaultUiMej9.group_name}", "${alicloud_wafv3_defense_resource_group.defaultbKntkl.group_name}", "${alicloud_wafv3_defense_resource_group.defaultSBiHAx.group_name}"},
					"resources": []string{
						"${alicloud_wafv3_domain.defaulttiuxAo.domain_id}", "${alicloud_wafv3_domain.default4lgADu.domain_id}", "${alicloud_wafv3_domain.defaultYZzU91.domain_id}"},
					"instance_id":           "${data.alicloud_wafv3_instances.default.ids.0}",
					"template_origin":       "custom",
					"defense_scene":         "ip_blacklist",
					"template_type":         "user_default",
					"defense_template_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                CHECKSET,
						"description":           "testTP",
						"resource_groups.#":     "3",
						"resources.#":           "3",
						"instance_id":           CHECKSET,
						"template_origin":       "custom",
						"defense_scene":         "ip_blacklist",
						"template_type":         "user_default",
						"defense_template_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":      "0",
					"description": "testacccc",
					"resource_groups": []string{
						"${alicloud_wafv3_defense_resource_group.defaultbKntkl.group_name}", "${alicloud_wafv3_defense_resource_group.defaultSBiHAx.group_name}"},
					"defense_template_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                CHECKSET,
						"description":           "testacccc",
						"resource_groups.#":     "2",
						"defense_template_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "testbbbb",
					"resource_groups": []string{
						"${alicloud_wafv3_defense_resource_group.defaultbKntkl.group_name}"},
					"resources": []string{
						"${alicloud_wafv3_domain.default4lgADu.domain_id}"},
					"defense_template_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           "testbbbb",
						"resource_groups.#":     "1",
						"resources.#":           "1",
						"defense_template_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_groups": []string{},
					"resources": []string{
						"${alicloud_wafv3_domain.default4lgADu.domain_id}", "${alicloud_wafv3_domain.defaultYZzU91.domain_id}", "${alicloud_wafv3_domain.defaulttiuxAo.domain_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_groups.#": "0",
						"resources.#":       "3",
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

var AlicloudWafv3DefenseTemplateMap11649 = map[string]string{
	"defense_template_id": CHECKSET,
}

func AlicloudWafv3DefenseTemplateBasicDependence11649(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_resource_group" "defaultUiMej9" {
  group_name  = "testformTF1"
  description = "test"
  instance_id = data.alicloud_wafv3_instances.default.ids.0
}

resource "alicloud_wafv3_defense_resource_group" "defaultbKntkl" {
  description   = "test"
  group_name    = "testformTF2"
  instance_id   = data.alicloud_wafv3_instances.default.ids.0
}

resource "alicloud_wafv3_defense_resource_group" "defaultSBiHAx" {
  description   = "test"
  group_name    = "testformTF3"
  instance_id   = data.alicloud_wafv3_instances.default.ids.0
}

resource "alicloud_wafv3_domain" "defaulttiuxAo" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  access_type = "share"
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_1.wafqax.top"
}

resource "alicloud_wafv3_domain" "default4lgADu" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  access_type = "share"
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_2.wafqax.top"
}

resource "alicloud_wafv3_domain" "defaultYZzU91" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  access_type = "share"
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_3.wafqax.top"
}

`, name)
}

// Case 防护模板_20250714_自定义模板_预期可重入_测试通过_绑定防护对象_未绑定防护对象组 11108
func TestAccAliCloudWafv3DefenseTemplate_basic11108(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_template.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseTemplateMap11108)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseTemplateBasicDependence11108)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${data.alicloud_wafv3_instances.default.ids.0}",
					"template_origin":       "custom",
					"defense_template_name": name,
					"defense_scene":         "ip_blacklist",
					"template_type":         "user_custom",
					"status":                "1",
					"description":           "testCreate",
					"resources":             []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"template_origin":       "custom",
						"defense_template_name": name,
						"defense_scene":         "ip_blacklist",
						"template_type":         "user_custom",
						"status":                CHECKSET,
						"description":           "testCreate",
						"resources.#":           "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"defense_template_name": name + "_update",
					"status":                "0",
					"description":           "testupdate",
					"resources": []string{
						"${alicloud_wafv3_domain.default1.domain_id}", "${alicloud_wafv3_domain.default2.domain_id}", "${alicloud_wafv3_domain.default3.domain_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"defense_template_name": name + "_update",
						"status":                CHECKSET,
						"description":           "testupdate",
						"resources.#":           "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []string{
						"${alicloud_wafv3_domain.default1.domain_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "0",
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

var AlicloudWafv3DefenseTemplateMap11108 = map[string]string{
	"defense_template_id": CHECKSET,
}

func AlicloudWafv3DefenseTemplateBasicDependence11108(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_domain" "default1" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  access_type = "share"
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_1.wafqax.top"
}

resource "alicloud_wafv3_domain" "default2" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  access_type = "share"
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_2.wafqax.top"
}

resource "alicloud_wafv3_domain" "default3" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  access_type = "share"
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_3.wafqax.top"
}


`, name)
}

// Case 防护模板_20250714_自定义模板_预期可重入_测试通过 11016
func TestAccAliCloudWafv3DefenseTemplate_basic11016(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_template.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseTemplateMap11016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseTemplateBasicDependence11016)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                        "${data.alicloud_wafv3_instances.default.ids.0}",
					"template_origin":                    "custom",
					"defense_template_name":              name,
					"defense_scene":                      "ip_blacklist",
					"template_type":                      "user_custom",
					"status":                             "1",
					"description":                        "testCreate",
					"resources":                          []string{},
					"resource_manager_resource_group_id": "${alicloud_resource_manager_resource_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":                        CHECKSET,
						"template_origin":                    "custom",
						"defense_template_name":              name,
						"defense_scene":                      "ip_blacklist",
						"template_type":                      "user_custom",
						"status":                             CHECKSET,
						"description":                        "testCreate",
						"resources.#":                        "0",
						"resource_manager_resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"defense_template_name": name + "_update",
					"status":                "0",
					"description":           "testupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"defense_template_name": name + "_update",
						"status":                CHECKSET,
						"description":           "testupdate",
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

var AlicloudWafv3DefenseTemplateMap11016 = map[string]string{
	"defense_template_id": CHECKSET,
}

func AlicloudWafv3DefenseTemplateBasicDependence11016(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

resource "alicloud_resource_manager_resource_group" "default" {
  name         = "testformTF"
  display_name = "testformTF"
}

data "alicloud_wafv3_instances" "default" {
}


`, name)
}

// Test Wafv3 DefenseTemplate. <<< Resource test cases, automatically generated.
