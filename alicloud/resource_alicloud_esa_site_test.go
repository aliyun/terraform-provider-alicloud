package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEsaSite_basic8490(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaSiteMap8490)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaSiteBasicDependence8490)
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
					"site_name":         name,
					"coverage":          "overseas",
					"access_type":       "NS",
					"instance_id":       "${alicloud_esa_rate_plan_instance.defaultIEoDfU.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_name":         name,
						"coverage":          "overseas",
						"access_type":       "NS",
						"instance_id":       CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"coverage":          "global",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"coverage":          "global",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudEsaSiteMap8490 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEsaSiteBasicDependence8490(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_esa_rate_plan_instance" "defaultIEoDfU" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "basic"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

`, name)
}

// Case Site 测试用例_副本 8484
func TestAccAliCloudEsaSite_basic8484(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaSiteMap8484)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaSiteBasicDependence8484)
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
					"site_name":         name,
					"coverage":          "overseas",
					"access_type":       "NS",
					"instance_id":       "${alicloud_esa_rate_plan_instance.defaultIEoDfU.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_name":         name,
						"coverage":          "overseas",
						"access_type":       "NS",
						"instance_id":       CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"coverage":          "global",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"coverage":          "global",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudEsaSiteMap8484 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEsaSiteBasicDependence8484(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_esa_rate_plan_instance" "defaultIEoDfU" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  plan_name    = "basic"
  auto_pay     = true
}


`, name)
}

// Case Site 测试用例_副 8288
func TestAccAliCloudEsaSite_basic8288(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaSiteMap8288)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaSiteBasicDependence8288)
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
					"site_name":         name,
					"coverage":          "overseas",
					"access_type":       "NS",
					"instance_id":       "${alicloud_esa_rate_plan_instance.defaultIEoDfU.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_name":         name,
						"coverage":          "overseas",
						"access_type":       "NS",
						"instance_id":       CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"coverage":          "global",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"coverage":          "global",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudEsaSiteMap8288 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEsaSiteBasicDependence8288(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_esa_rate_plan_instance" "defaultIEoDfU" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  plan_name    = "basic"
  auto_pay     = true
}


`, name)
}

// Case Site 测试用例_副本1728637287794 8106
func TestAccAliCloudEsaSite_basic8106(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaSiteMap8106)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaSiteBasicDependence8106)
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
					"site_name":         name,
					"coverage":          "overseas",
					"access_type":       "NS",
					"instance_id":       "${alicloud_esa_rate_plan_instance.defaultIEoDfU.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_name":         name,
						"coverage":          "overseas",
						"access_type":       "NS",
						"instance_id":       CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"coverage":          "global",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"coverage":          "global",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudEsaSiteMap8106 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEsaSiteBasicDependence8106(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_esa_rate_plan_instance" "defaultIEoDfU" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  plan_name    = "basic"
  auto_pay     = true
}


`, name)
}

// Case Site测试用例_站点配置 9798
func TestAccAliCloudEsaSite_basic9798(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaSiteMap9798)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaSiteBasicDependence9798)
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
					"site_name":                     name,
					"coverage":                      "overseas",
					"access_type":                   "NS",
					"instance_id":                   "esa-site-a755h6gbkx6o",
					"cache_architecture_mode":       "edge",
					"add_real_client_ip_header":     "on",
					"add_client_geolocation_header": "off",
					"ipv6_enable":                   "on",
					"site_version":                  "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_name":                     name,
						"coverage":                      "overseas",
						"access_type":                   "NS",
						"instance_id":                   CHECKSET,
						"cache_architecture_mode":       "edge",
						"add_real_client_ip_header":     "on",
						"add_client_geolocation_header": "off",
						"ipv6_enable":                   "on",
						"site_version":                  "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"coverage":                      "global",
					"add_client_geolocation_header": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"coverage":                      "global",
						"add_client_geolocation_header": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudEsaSiteMap9798 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEsaSiteBasicDependence9798(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Esa Site. >>> Resource test cases, automatically generated.
// Case site配置 10677
func TestAccAliCloudEsaSite_basic10677(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_site.default"
	ra := resourceAttrInit(resourceId, AlicloudEsaSiteMap10677)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("bcd%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEsaSiteBasicDependence10677)
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
					"site_name":                     name,
					"instance_id":                   "${data.alicloud_esa_sites.default.sites.0.instance_id}",
					"development_mode":              "on",
					"cache_reserve_instance_id":     "${alicloud_esa_cache_reserve_instance.default.id}",
					"cache_architecture_mode":       "edge",
					"case_insensitive":              "on",
					"coverage":                      "overseas",
					"add_real_client_ip_header":     "on",
					"site_name_exclusive":           "on",
					"ipv6_region":                   "x.x",
					"seo_bypass":                    "on",
					"cache_reserve_enable":          "on",
					"add_client_geolocation_header": "on",
					"cross_border_optimization":     "on",
					"ipv6_enable":                   "on",
					"site_version":                  "0",
					"flatten_mode":                  "flatten_all",
					"access_type":                   "NS",
					"version_management":            "true",
					"tag_name":                      "tagtest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_name":                     name,
						"instance_id":                   CHECKSET,
						"development_mode":              "on",
						"cache_reserve_instance_id":     CHECKSET,
						"cache_architecture_mode":       "edge",
						"case_insensitive":              "on",
						"coverage":                      "overseas",
						"add_real_client_ip_header":     "on",
						"site_name_exclusive":           "on",
						"ipv6_region":                   "x.x",
						"seo_bypass":                    "on",
						"cache_reserve_enable":          "on",
						"add_client_geolocation_header": "on",
						"cross_border_optimization":     "on",
						"ipv6_enable":                   "on",
						"site_version":                  "0",
						"flatten_mode":                  "flatten_all",
						"access_type":                   "NS",
						"version_management":            "true",
						"tag_name":                      "tagtest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"development_mode":              "off",
					"cache_reserve_instance_id":     "${alicloud_esa_cache_reserve_instance.update.id}",
					"cache_architecture_mode":       "edge_smart",
					"case_insensitive":              "off",
					"add_real_client_ip_header":     "off",
					"site_name_exclusive":           "off",
					"ipv6_region":                   "cn.cn",
					"seo_bypass":                    "off",
					"cache_reserve_enable":          "off",
					"add_client_geolocation_header": "off",
					"cross_border_optimization":     "off",
					"ipv6_enable":                   "off",
					"flatten_mode":                  "flatten_at_root",
					"version_management":            "false",
					"tag_name":                      "tst",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"development_mode":              "off",
						"cache_reserve_instance_id":     CHECKSET,
						"cache_architecture_mode":       "edge_smart",
						"case_insensitive":              "off",
						"add_real_client_ip_header":     "off",
						"site_name_exclusive":           "off",
						"ipv6_region":                   "cn.cn",
						"seo_bypass":                    "off",
						"cache_reserve_enable":          "off",
						"add_client_geolocation_header": "off",
						"cross_border_optimization":     "off",
						"ipv6_enable":                   "off",
						"flatten_mode":                  "flatten_at_root",
						"version_management":            "false",
						"tag_name":                      "tst",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudEsaSiteMap10677 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEsaSiteBasicDependence10677(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_cache_reserve_instance" "default" {
  quota_gb     = "10240"
  cr_region    = "CN-beijing"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  auto_pay     = true
}

resource "alicloud_esa_cache_reserve_instance" "update" {
  quota_gb     = "10240"
  cr_region    = "CN-beijing"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  auto_pay     = true
}

`, name)
}

// Test Esa Site. <<< Resource test cases, automatically generated.
