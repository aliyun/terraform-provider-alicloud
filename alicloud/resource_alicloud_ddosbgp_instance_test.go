package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DdosBgp Instance. >>> Resource test cases, automatically generated.
// Case 原生1.0实例_v6实例 11131
func TestAccAliCloudDdosBgpInstance_basic11131(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddosbgp_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudDdosBgpInstanceMap11131)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdosBgpInstanceBasicDependence11131)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":        "-1",
					"ip_count":         "100",
					"ip_type":          "IPv6",
					"normal_bandwidth": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":        "-1",
						"ip_count":         "100",
						"ip_type":          "IPv6",
						"normal_bandwidth": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func TestAccAliCloudDdosBgpInstance_basic11131_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddosbgp_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudDdosBgpInstanceMap11131)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdosBgpInstanceBasicDependence11131)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":         "-1",
					"base_bandwidth":    "20",
					"instance_name":     name,
					"ip_count":          "100",
					"ip_type":           "IPv6",
					"normal_bandwidth":  "300",
					"period":            "12",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"type":              "Enterprise",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":         "-1",
						"base_bandwidth":    "20",
						"instance_name":     name,
						"ip_count":          "100",
						"ip_type":           "IPv6",
						"normal_bandwidth":  "300",
						"resource_group_id": CHECKSET,
						"type":              "Enterprise",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

var AliCloudDdosBgpInstanceMap11131 = map[string]string{
	"resource_group_id": CHECKSET,
	"type":              CHECKSET,
	"status":            CHECKSET,
}

func AliCloudDdosBgpInstanceBasicDependence11131(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}
`, name)
}

// Case 原生1.0实例_专业版 11229 适配废弃字段name
func TestAccAliCloudDdosBgpInstance_basic11229(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddosbgp_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudDdosBgpInstanceMap11229)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdosBgpInstanceBasicDependence11229)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":        "-1",
					"ip_count":         "100",
					"ip_type":          "IPv4",
					"normal_bandwidth": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":        "-1",
						"ip_count":         "100",
						"ip_type":          "IPv4",
						"normal_bandwidth": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func TestAccAliCloudDdosBgpInstance_basic11229_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddosbgp_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudDdosBgpInstanceMap11229)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosBgpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosBgpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdosBgpInstanceBasicDependence11229)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":         "20",
					"base_bandwidth":    "20",
					"name":              name,
					"ip_count":          "1",
					"ip_type":           "IPv4",
					"normal_bandwidth":  "100",
					"period":            "12",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"type":              "Professional",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":         "20",
						"base_bandwidth":    "20",
						"name":              name,
						"ip_count":          "1",
						"ip_type":           "IPv4",
						"normal_bandwidth":  "100",
						"resource_group_id": CHECKSET,
						"type":              "Professional",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

var AliCloudDdosBgpInstanceMap11229 = map[string]string{
	"resource_group_id": CHECKSET,
	"type":              CHECKSET,
	"status":            CHECKSET,
}

func AliCloudDdosBgpInstanceBasicDependence11229(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}
`, name)
}

// Test DdosBgp Instance. <<< Resource test cases, automatically generated.
