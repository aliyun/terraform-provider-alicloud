package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test PaiWorkspace Dataset. >>> Resource test cases, automatically generated.
// Case Dataset资源测试用例_副本1730104187187_副本1732006619708 9002
func TestAccAliCloudPaiWorkspaceDataset_basic9002(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_dataset.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetMap9002)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetBasicDependence9002)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/\\\"}",
					"description":      "镇元资源测试用例Dataset",
					"accessibility":    "PRIVATE",
					"dataset_name":     name,
					"data_source_type": "NAS",
					"source_type":      "ITAG",
					"workspace_id":     "${alicloud_pai_workspace_workspace.defaultWorkspace.id}",
					"data_type":        "PIC",
					"property":         "DIRECTORY",
					"uri":              "nas://086b649545.cn-hangzhou/",
					"source_id":        "d-xxxxx_v1",
					"user_id":          "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/\"}",
						"description":      "镇元资源测试用例Dataset",
						"accessibility":    "PRIVATE",
						"dataset_name":     name,
						"data_source_type": "NAS",
						"source_type":      "ITAG",
						"workspace_id":     CHECKSET,
						"data_type":        "PIC",
						"property":         "DIRECTORY",
						"uri":              "nas://086b649545.cn-hangzhou/",
						"source_id":        "d-xxxxx_v1",
						"user_id":          "1511928242963727",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":       "{\\\"mountPath\\\":\\\"/mnt/data/update/\\\"}",
					"description":   "镇元资源测试用例Dataset_update",
					"accessibility": "PUBLIC",
					"dataset_name":  name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "key_new",
							"value": "value_new",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":       "{\"mountPath\":\"/mnt/data/update/\"}",
						"description":   "镇元资源测试用例Dataset_update",
						"accessibility": "PUBLIC",
						"dataset_name":  name + "_update",
						"labels.#":      "1",
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

var AlicloudPaiWorkspaceDatasetMap9002 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetBasicDependence9002(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultWorkspace" {
  description    = "dataset_pop_test_438"
  display_name   = "DatasetResouceTest_956"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Case Dataset资源测试用例_副本1730104187187 8536
func TestAccAliCloudPaiWorkspaceDataset_basic8536(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_dataset.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetMap8536)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetBasicDependence8536)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/\\\"}",
					"description":      "镇元资源测试用例Dataset",
					"accessibility":    "PRIVATE",
					"dataset_name":     name,
					"data_source_type": "NAS",
					"source_type":      "ITAG",
					"workspace_id":     "${alicloud_pai_workspace_workspace.defaultWorkspace.id}",
					"data_type":        "PIC",
					"property":         "DIRECTORY",
					"uri":              "nas://086b649545.cn-hangzhou/",
					"source_id":        "d-xxxxx_v1",
					"user_id":          "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/\"}",
						"description":      "镇元资源测试用例Dataset",
						"accessibility":    "PRIVATE",
						"dataset_name":     name,
						"data_source_type": "NAS",
						"source_type":      "ITAG",
						"workspace_id":     CHECKSET,
						"data_type":        "PIC",
						"property":         "DIRECTORY",
						"uri":              "nas://086b649545.cn-hangzhou/",
						"source_id":        "d-xxxxx_v1",
						"user_id":          "1511928242963727",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":       "{\\\"mountPath\\\":\\\"/mnt/data/update/\\\"}",
					"description":   "镇元资源测试用例Dataset_update",
					"accessibility": "PUBLIC",
					"dataset_name":  name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "key_new",
							"value": "value_new",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":       "{\"mountPath\":\"/mnt/data/update/\"}",
						"description":   "镇元资源测试用例Dataset_update",
						"accessibility": "PUBLIC",
						"dataset_name":  name + "_update",
						"labels.#":      "1",
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

var AlicloudPaiWorkspaceDatasetMap8536 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetBasicDependence8536(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultWorkspace" {
  description    = "dataset_pop_test_732"
  display_name   = "DatasetResouceTest_26"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Case Dataset资源测试用例_demo1 8538
func TestAccAliCloudPaiWorkspaceDataset_basic8538(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_dataset.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetMap8538)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetBasicDependence8538)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/\\\"}",
					"description":      "镇元资源测试用例Dataset",
					"accessibility":    "PRIVATE",
					"dataset_name":     name,
					"data_source_type": "OSS",
					"source_type":      "USER",
					"workspace_id":     "${alicloud_pai_workspace_workspace.defaultWorkspace.id}",
					"data_type":        "COMMON",
					"labels": []map[string]interface{}{
						{
							"key":   "key1",
							"value": "zy_test1",
						},
						{
							"key":   "key2",
							"value": "zy_test2",
						},
						{
							"key":   "key3",
							"value": "zy_test3",
						},
					},
					"property":  "DIRECTORY",
					"uri":       "oss://bucket-ruoli.oss-cn-hangzhou.aliyuncs.com/test/",
					"source_id": "d-xxxxx_v1",
					"user_id":   "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/\"}",
						"description":      "镇元资源测试用例Dataset",
						"accessibility":    "PRIVATE",
						"dataset_name":     name,
						"data_source_type": "OSS",
						"source_type":      "USER",
						"workspace_id":     CHECKSET,
						"data_type":        "COMMON",
						"labels.#":         "3",
						"property":         "DIRECTORY",
						"uri":              "oss://bucket-ruoli.oss-cn-hangzhou.aliyuncs.com/test/",
						"source_id":        "d-xxxxx_v1",
						"user_id":          "1511928242963727",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":       "{\\\"mountPath\\\":\\\"/mnt/data/update/\\\"}",
					"description":   "镇元资源测试用例Dataset_update",
					"accessibility": "PUBLIC",
					"dataset_name":  name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "key_new",
							"value": "value_new",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":       "{\"mountPath\":\"/mnt/data/update/\"}",
						"description":   "镇元资源测试用例Dataset_update",
						"accessibility": "PUBLIC",
						"dataset_name":  name + "_update",
						"labels.#":      "1",
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

var AlicloudPaiWorkspaceDatasetMap8538 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetBasicDependence8538(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultWorkspace" {
  description    = "dataset_pop_test_535"
  display_name   = "DatasetResouceTest_16"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Case Dataset资源测试用例_demo2 8539
func TestAccAliCloudPaiWorkspaceDataset_basic8539(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_dataset.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetMap8539)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetBasicDependence8539)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/\\\"}",
					"description":      "镇元资源测试用例Dataset",
					"accessibility":    "PRIVATE",
					"dataset_name":     name,
					"data_source_type": "OSS",
					"source_type":      "USER",
					"workspace_id":     "${alicloud_pai_workspace_workspace.defaultWorkspace.id}",
					"data_type":        "TEXT",
					"labels": []map[string]interface{}{
						{
							"key":   "key1",
							"value": "zy_test1",
						},
						{
							"key":   "key2",
							"value": "zy_test2",
						},
						{
							"key":   "key3",
							"value": "zy_test3",
						},
					},
					"property":  "FILE",
					"source_id": "d-xxxxx_v1",
					"user_id":   "1511928242963727",
					"uri":       "oss://bucket-ruoli.oss-cn-hangzhou.aliyuncs.com/test/a.csv",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/\"}",
						"description":      "镇元资源测试用例Dataset",
						"accessibility":    "PRIVATE",
						"dataset_name":     name,
						"data_source_type": "OSS",
						"source_type":      "USER",
						"workspace_id":     CHECKSET,
						"data_type":        "TEXT",
						"labels.#":         "3",
						"property":         "FILE",
						"source_id":        "d-xxxxx_v1",
						"user_id":          "1511928242963727",
						"uri":              "oss://bucket-ruoli.oss-cn-hangzhou.aliyuncs.com/test/a.csv",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":       "{\\\"mountPath\\\":\\\"/mnt/data/update/\\\"}",
					"description":   "镇元资源测试用例Dataset_update",
					"accessibility": "PUBLIC",
					"dataset_name":  name + "_update",
					"labels":        REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":       "{\"mountPath\":\"/mnt/data/update/\"}",
						"description":   "镇元资源测试用例Dataset_update",
						"accessibility": "PUBLIC",
						"dataset_name":  name + "_update",
						"labels.#":      "0",
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

var AlicloudPaiWorkspaceDatasetMap8539 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetBasicDependence8539(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultWorkspace" {
  description    = "dataset_pop_test_185"
  display_name   = "DatasetResouceTest_146"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Case Dataset资源测试用例 8519
func TestAccAliCloudPaiWorkspaceDataset_basic8519(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_dataset.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetMap8519)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetBasicDependence8519)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/\\\"}",
					"description":      "镇元资源测试用例Dataset",
					"accessibility":    "PRIVATE",
					"dataset_name":     name,
					"data_source_type": "OSS",
					"source_type":      "USER",
					"workspace_id":     "${alicloud_pai_workspace_workspace.defaultWorkspace.id}",
					"data_type":        "COMMON",
					"labels": []map[string]interface{}{
						{
							"key":   "key1",
							"value": "zy_test1",
						},
						{
							"key":   "key2",
							"value": "zy_test2",
						},
						{
							"key":   "key3",
							"value": "zy_test3",
						},
					},
					"property":  "DIRECTORY",
					"uri":       "oss://bucket-ruoli.oss-cn-hangzhou.aliyuncs.com/test/",
					"source_id": "d-xxxxx_v1",
					"user_id":   "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/\"}",
						"description":      "镇元资源测试用例Dataset",
						"accessibility":    "PRIVATE",
						"dataset_name":     name,
						"data_source_type": "OSS",
						"source_type":      "USER",
						"workspace_id":     CHECKSET,
						"data_type":        "COMMON",
						"labels.#":         "3",
						"property":         "DIRECTORY",
						"uri":              "oss://bucket-ruoli.oss-cn-hangzhou.aliyuncs.com/test/",
						"source_id":        "d-xxxxx_v1",
						"user_id":          "1511928242963727",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":       "{\\\"mountPath\\\":\\\"/mnt/data/update/\\\"}",
					"description":   "镇元资源测试用例Dataset_update",
					"accessibility": "PUBLIC",
					"dataset_name":  name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "key_new",
							"value": "value_new",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":       "{\"mountPath\":\"/mnt/data/update/\"}",
						"description":   "镇元资源测试用例Dataset_update",
						"accessibility": "PUBLIC",
						"dataset_name":  name + "_update",
						"labels.#":      "1",
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

var AlicloudPaiWorkspaceDatasetMap8519 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetBasicDependence8519(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultWorkspace" {
  description    = "dataset_pop_test_832"
  display_name   = "DatasetResouceTest_608"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Case Dataset资源测试用例 7381
func TestAccAliCloudPaiWorkspaceDataset_basic7381(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_dataset.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetMap7381)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetBasicDependence7381)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/\\\"}",
					"description":      "镇元资源测试用例Dataset",
					"accessibility":    "PRIVATE",
					"dataset_name":     name,
					"data_source_type": "OSS",
					"source_type":      "USER",
					"workspace_id":     "${alicloud_pai_workspace_workspace.default6RZ5Rl.id}",
					"data_type":        "COMMON",
					"labels": []map[string]interface{}{
						{
							"key":   "key1",
							"value": "zy_test1",
						},
					},
					"property":  "DIRECTORY",
					"uri":       "oss://bucket-ruoli.oss-cn-hangzhou.aliyuncs.com/test/",
					"source_id": "d-xxxxx_v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/\"}",
						"description":      "镇元资源测试用例Dataset",
						"accessibility":    "PRIVATE",
						"dataset_name":     name,
						"data_source_type": "OSS",
						"source_type":      "USER",
						"workspace_id":     CHECKSET,
						"data_type":        "COMMON",
						"labels.#":         "1",
						"property":         "DIRECTORY",
						"uri":              "oss://bucket-ruoli.oss-cn-hangzhou.aliyuncs.com/test/",
						"source_id":        "d-xxxxx_v1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":       "{\\\"mountPath\\\":\\\"/mnt/data/update/\\\"}",
					"description":   "镇元资源测试用例Dataset_update",
					"accessibility": "PUBLIC",
					"dataset_name":  name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "key_new",
							"value": "value_new",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":       "{\"mountPath\":\"/mnt/data/update/\"}",
						"description":   "镇元资源测试用例Dataset_update",
						"accessibility": "PUBLIC",
						"dataset_name":  name + "_update",
						"labels.#":      "1",
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

var AlicloudPaiWorkspaceDatasetMap7381 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetBasicDependence7381(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "default6RZ5Rl" {
  description    = "dataset_pop_test_56"
  display_name   = "DatasetTestCaseWorkSpace_753"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Case Dataset governance coverage (M-RT-0040 property_completeness + E-TF-0006 property_coverage) gov
func TestAccAliCloudPaiWorkspaceDataset_governanceCoverage(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_dataset.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetMapGov)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetBasicDependenceGov)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// dataset_task_ram_role is edition-dependent: the Create API rejects
				// it for BASIC datasets (400 "DatasetTaskRamRole is not supported
				// for BASIC dataset"). Edition is API-managed (Computed) and cannot
				// be forced to non-BASIC from Terraform, so the success path cannot
				// exercise dataset_task_ram_role. This step still sends
				// DatasetTaskRamRole in the Create request (exercises the schema
				// mapping) and pins the edition-dependent constraint via ExpectError.
				Config: testAccConfig(map[string]interface{}{
					"description":             "governance_coverage_create",
					"accessibility":           "ROLE_PUBLIC",
					"dataset_name":            name,
					"data_source_type":        "NAS",
					"source_type":             "USER",
					"workspace_id":            "${alicloud_pai_workspace_workspace.defaultWorkspace.id}",
					"data_type":               "COMMON",
					"property":                "DIRECTORY",
					"uri":                     "nas://086b649545.cn-hangzhou/",
					"source_id":               "d-xxxxx_v1",
					"user_id":                 "1511928242963727",
					"accessible_role_id_list": []interface{}{"PAI_Workspace_Admin"},
					"dataset_task_ram_role":   "acs:ram::1234567890123456:role/test-role",
					"user_metrics_endpoints": []map[string]interface{}{
						{
							"port": 9090,
							"path": "/metrics",
						},
					},
				}),
				ExpectError: regexp.MustCompile("DatasetTaskRamRole is not supported"),
			},
			{
				// Exercise the user_metrics_endpoints Create mapping (M-RT-0040) via
				// ExpectError: with dataset_task_ram_role and accessible_role_id_list
				// removed, the Create API reaches UserMetricsEndpoints validation and
				// rejects it for NAS datasets ("UserMetricsEndpoints is only supported
				// for JUICEFS datasets"), proving UserMetricsEndpoints is sent in the
				// Create request. accessibility=PRIVATE avoids the workspace role lookup.
				Config: testAccConfig(map[string]interface{}{
					"dataset_task_ram_role":   REMOVEKEY,
					"accessible_role_id_list": REMOVEKEY,
					"accessibility":           "PRIVATE",
				}),
				ExpectError: regexp.MustCompile("UserMetricsEndpoints is only supported for JUICEFS"),
			},
			{
				// Create on BASIC for the success path: remove user_metrics_endpoints
				// (JUICEFS-only, unavailable in the ACC workspace). The dataset_task_ram_role
				// and accessible_role_id_list removed in the prior step stay removed
				// (edition/role constraints); accessibility stays PRIVATE (no role lookup).
				// The subsequent Update/Import/Update-mapping steps build on this resource.
				Config: testAccConfig(map[string]interface{}{
					"user_metrics_endpoints": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "governance_coverage_create",
						"accessibility":    "PRIVATE",
						"dataset_name":     name,
						"data_source_type": "NAS",
						"source_type":      "USER",
						"workspace_id":     CHECKSET,
						"data_type":        "COMMON",
						"property":         "DIRECTORY",
						"uri":              "nas://086b649545.cn-hangzhou/",
						"source_id":        "d-xxxxx_v1",
						"user_id":          "1511928242963727",
					}),
				),
			},
			{
				// Update description + dataset_name (success path). accessibility stays
				// PRIVATE and accessible_role_id_list stays absent (no valid role codes
				// in the ACC workspace); UpdateDataset receives Description + Name only.
				Config: testAccConfig(map[string]interface{}{
					"description":  "governance_coverage_update",
					"dataset_name": name + "_upd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "governance_coverage_update",
						"accessibility": "PRIVATE",
						"dataset_name":  name + "_upd",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset_task_ram_role", "user_metrics_endpoints"},
			},
			{
				// Exercise the accessible_role_id_list Update mapping (M-RT-0040) via
				// ExpectError: the ACC workspace has no valid role codes, so UpdateDataset
				// rejects the role list with "Role codes are not in the workspace".
				// ROLE_PUBLIC != PUBLIC so no /publish; UpdateDataset receives
				// Accessibility + AccessibleRoleIdList and validates the roles.
				Config: testAccConfig(map[string]interface{}{
					"accessibility":           "ROLE_PUBLIC",
					"accessible_role_id_list": []interface{}{"PAI_Workspace_Admin", "role-testcustom"},
				}),
				ExpectError: regexp.MustCompile("Role codes are not in the workspace"),
			},
		},
	})
}

var AlicloudPaiWorkspaceDatasetMapGov = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetBasicDependenceGov(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultWorkspace" {
  description    = "dataset_pop_test_gov"
  display_name   = "DatasetResouceTest_Gov"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Test PaiWorkspace Dataset. <<< Resource test cases, automatically generated.
