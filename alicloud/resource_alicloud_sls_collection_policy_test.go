package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls CollectionPolicy. >>> Resource test cases, automatically generated.
// Case TF验收_资源目录 8160
func TestAccAliCloudSlsCollectionPolicy_basic8160(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  "xc-test-oss-01",
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "metering_log",
						"product_code": "oss",
						"policy_name":  "xc-test-oss-01",
						"enabled":      "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "true",
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members": []string{
								"1936728897040477"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou", "cn-shanghai", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "sdd", "des"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.update_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore002.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members": []string{
								"1863928896950968", "1403628895365909", "1751928895291600"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members": []string{
								"1863928896950968", "1403628895365909"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.update_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore002.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
							"members":            []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "true",
					"product_code":       "oss",
					"policy_name":        "xc-test-oss-01",
					"enabled":            "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members": []string{
								"1936728897040477"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "true",
						"product_code":       "oss",
						"policy_name":        "xc-test-oss-01",
						"enabled":            "true",
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

var AlicloudSlsCollectionPolicyMap8160 = map[string]string{}

func AlicloudSlsCollectionPolicyBasicDependence8160(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "project_create_01" {
  description = "xxccc"
  name        = var.name
}

resource "alicloud_log_store" "logstore_create_01" {
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.project_create_01.name
  name             = format("%%s1", var.name)
}

resource "alicloud_log_project" "update_01" {
  description = "dfdfd"
  name        = format("%%s2", var.name)
}

resource "alicloud_log_store" "logstore002" {
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.update_01.name
  name             = format("%%s3", var.name)
}


`, name)
}

// Case TF验收 8140
func TestAccAliCloudSlsCollectionPolicy_basic8140(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8140)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8140)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"instance_ids": []string{
								"xcd", "sdd"},
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  "xc-test-oss-01",
					"enabled":      "false",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "metering_log",
						"product_code": "oss",
						"policy_name":  "xc-test-oss-01",
						"enabled":      "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
						},
					},
					"centralize_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
								"\"xc-sx\"":   "xxc",
							},
							"regions": []string{
								"cn-hangzhou", "cn-shanghai", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "sdd", "des"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project02.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore002.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
								"\"xc-sx\"":   "xxc",
							},
							"regions": []string{
								"cn-hangzhou", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "des"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"regions":       []string{},
							"instance_ids":  []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"instance_ids": []string{
								"xcd", "sdd"},
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "false",
					"product_code":       "oss",
					"policy_name":        "xc-test-oss-01",
					"enabled":            "false",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "false",
						"product_code":       "oss",
						"policy_name":        "xc-test-oss-01",
						"enabled":            "false",
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

var AlicloudSlsCollectionPolicyMap8140 = map[string]string{}

func AlicloudSlsCollectionPolicyBasicDependence8140(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "all_project01" {
  description = "xccx"
  name        = var.name
}

resource "alicloud_log_store" "all_logstore_01" {
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.all_project01.name
  name             = format("%%s1", var.name)
}

resource "alicloud_log_project" "all_project02" {
  description = "cvcvc"
  name        = format("%%s2", var.name)
}

resource "alicloud_log_store" "all_logstore002" {
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.all_project02.name
  name             = format("%%s3", var.name)
}


`, name)
}

// Case TF验收_ResourceMode 7850
func TestAccAliCloudSlsCollectionPolicy_basic7850(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap7850)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence7850)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
						},
					},
					"data_code":    "operation_log",
					"product_code": "project",
					"policy_name":  "xc-test-project-01",
					"enabled":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "operation_log",
						"product_code": "project",
						"policy_name":  "xc-test-project-01",
						"enabled":      "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "attributeMode",
							"resource_tags": map[string]interface{}{
								"\"xc-pro\"": "xc-sd",
							},
							"regions": []string{
								"cn-hangzhou", "cn-beijing"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "instanceMode",
							"regions":       []string{},
							"instance_ids": []string{
								"xxx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
						},
					},
					"data_code":          "operation_log",
					"centralize_enabled": "false",
					"product_code":       "project",
					"policy_name":        "xc-test-project-01",
					"enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "operation_log",
						"centralize_enabled": "false",
						"product_code":       "project",
						"policy_name":        "xc-test-project-01",
						"enabled":            "false",
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

var AlicloudSlsCollectionPolicyMap7850 = map[string]string{}

func AlicloudSlsCollectionPolicyBasicDependence7850(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case TF验收_资源目录 8160  twin
func TestAccAliCloudSlsCollectionPolicy_basic8160_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "true",
					"product_code":       "oss",
					"policy_name":        "xc-test-oss-01",
					"enabled":            "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members": []string{
								"1936728897040477"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "true",
						"product_code":       "oss",
						"policy_name":        "xc-test-oss-01",
						"enabled":            "true",
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

// Case TF验收 8140  twin
func TestAccAliCloudSlsCollectionPolicy_basic8140_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8140)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8140)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"instance_ids": []string{
								"xcd", "sdd"},
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "false",
					"product_code":       "oss",
					"policy_name":        "xc-test-oss-01",
					"enabled":            "false",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "false",
						"product_code":       "oss",
						"policy_name":        "xc-test-oss-01",
						"enabled":            "false",
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

// Case TF验收_ResourceMode 7850  twin
func TestAccAliCloudSlsCollectionPolicy_basic7850_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap7850)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence7850)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
						},
					},
					"data_code":          "operation_log",
					"centralize_enabled": "false",
					"product_code":       "project",
					"policy_name":        "xc-test-project-01",
					"enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "operation_log",
						"centralize_enabled": "false",
						"product_code":       "project",
						"policy_name":        "xc-test-project-01",
						"enabled":            "false",
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

// Case TF验收_资源目录 8160  raw
func TestAccAliCloudSlsCollectionPolicy_basic8160_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "true",
					"product_code":       "oss",
					"policy_name":        "xc-test-oss-01",
					"enabled":            "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members": []string{
								"1936728897040477"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "true",
						"product_code":       "oss",
						"policy_name":        "xc-test-oss-01",
						"enabled":            "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou", "cn-shanghai", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "sdd", "des"},
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.update_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore002.name}",
							"dest_ttl":      "7",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members": []string{
								"1863928896950968", "1403628895365909", "1751928895291600"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
							"dest_ttl":      "7",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members": []string{
								"1863928896950968", "1403628895365909"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.update_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore002.name}",
							"dest_ttl":      "7",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
							"members":            []string{},
						},
					},
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

// Case TF验收 8140  raw
func TestAccAliCloudSlsCollectionPolicy_basic8140_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8140)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8140)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"instance_ids": []string{
								"xcd", "sdd"},
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "false",
					"product_code":       "oss",
					"policy_name":        "xc-test-oss-01",
					"enabled":            "false",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "false",
						"product_code":       "oss",
						"policy_name":        "xc-test-oss-01",
						"enabled":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
								"\"xc-sx\"":   "xxc",
							},
							"regions": []string{
								"cn-hangzhou", "cn-shanghai", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "sdd", "des"},
						},
					},
					"centralize_enabled": "true",
					"enabled":            "true",
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project02.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore002.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "true",
						"enabled":            "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
								"\"xc-sx\"":   "xxc",
							},
							"regions": []string{
								"cn-hangzhou", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "des"},
						},
					},
					"centralize_enabled": "false",
					"enabled":            "false",
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "false",
						"enabled":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"regions":       []string{},
							"instance_ids":  []string{},
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   defaultRegionToTest,
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
							"dest_ttl":      "7",
						},
					},
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

// Case TF验收_ResourceMode 7850  raw
func TestAccAliCloudSlsCollectionPolicy_basic7850_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap7850)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsCollectionPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence7850)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
						},
					},
					"data_code":          "operation_log",
					"centralize_enabled": "false",
					"product_code":       "project",
					"policy_name":        "xc-test-project-01",
					"enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "operation_log",
						"centralize_enabled": "false",
						"product_code":       "project",
						"policy_name":        "xc-test-project-01",
						"enabled":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "attributeMode",
							"resource_tags": map[string]interface{}{
								"\"xc-pro\"": "xc-sd",
							},
							"regions": []string{
								"cn-hangzhou", "cn-beijing"},
						},
					},
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "instanceMode",
							"regions":       []string{},
							"instance_ids": []string{
								"xxx"},
						},
					},
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
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

// Test Sls CollectionPolicy. <<< Resource test cases, automatically generated.
