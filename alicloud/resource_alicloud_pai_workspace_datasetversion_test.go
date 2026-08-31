package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace Datasetversion. >>> Resource test cases, automatically generated.
// Case DatasetVersion测试用例1_副本1732008135754 9008
func TestAccAliCloudPaiWorkspaceDatasetversion_basic9008(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_datasetversion.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetversionMap9008)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDatasetversion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetversionBasicDependence9008)
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
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/version/\\\"}",
					"description":      "镇元资源测试用例DatasetVersion",
					"data_source_type": "OSS",
					"source_type":      "USER",
					"source_id":        "d-xxxxx_v1",
					"data_size":        "2068",
					"data_count":       "1000",
					"labels": []map[string]interface{}{
						{
							"key":   "key1",
							"value": "test1",
						},
					},
					"uri":        "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
					"property":   "DIRECTORY",
					"dataset_id": "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/version/\"}",
						"description":      "镇元资源测试用例DatasetVersion",
						"data_source_type": "OSS",
						"source_type":      "USER",
						"source_id":        "d-xxxxx_v1",
						"data_size":        "2068",
						"data_count":       "1000",
						"labels.#":         "1",
						"uri":              "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
						"property":         "DIRECTORY",
						"dataset_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":     "{\\\"mountPath\\\":\\\"/mnt/data/new/\\\"}",
					"description": "镇元资源测试用例DatasetVersion new",
					"data_size":   "5096",
					"data_count":  "1001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":     "{\"mountPath\":\"/mnt/data/new/\"}",
						"description": "镇元资源测试用例DatasetVersion new",
						"data_size":   "5096",
						"data_count":  "1001",
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

var AlicloudPaiWorkspaceDatasetversionMap9008 = map[string]string{
	"create_time":  CHECKSET,
	"version_name": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetversionBasicDependence9008(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultAiWorkspace" {
  description    = "DatasetResouceTest_512"
  display_name   = "DatasetVerResouce_158"
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_dataset" "defaultDataset" {
  accessibility    = "PRIVATE"
  source_type      = "USER"
  data_type        = "PIC"
  workspace_id     = alicloud_pai_workspace_workspace.defaultAiWorkspace.id
  options          = "{\"mountPath\":\"/mnt/data/\"}"
  description      = "镇元资源测试用例Dataset"
  source_id        = "d-xxxxx_v1"
  uri              = "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/"
  dataset_name     = format("%%s1", var.name)
  user_id          = "1511928242963727"
  data_source_type = "OSS"
  property         = "DIRECTORY"
}


`, name)
}

// Case DatasetVersion测试用例1 8541
func TestAccAliCloudPaiWorkspaceDatasetversion_basic8541(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_datasetversion.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetversionMap8541)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDatasetversion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetversionBasicDependence8541)
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
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/version/\\\"}",
					"description":      "镇元资源测试用例DatasetVersion",
					"data_source_type": "OSS",
					"source_type":      "USER",
					"source_id":        "d-xxxxx_v1",
					"data_size":        "2068",
					"data_count":       "1000",
					"labels": []map[string]interface{}{
						{
							"key":   "key1",
							"value": "test1",
						},
					},
					"uri":        "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
					"property":   "DIRECTORY",
					"dataset_id": "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/version/\"}",
						"description":      "镇元资源测试用例DatasetVersion",
						"data_source_type": "OSS",
						"source_type":      "USER",
						"source_id":        "d-xxxxx_v1",
						"data_size":        "2068",
						"data_count":       "1000",
						"labels.#":         "1",
						"uri":              "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
						"property":         "DIRECTORY",
						"dataset_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":     "{\\\"mountPath\\\":\\\"/mnt/data/new/\\\"}",
					"description": "镇元资源测试用例DatasetVersion new",
					"data_size":   "5096",
					"data_count":  "1001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":     "{\"mountPath\":\"/mnt/data/new/\"}",
						"description": "镇元资源测试用例DatasetVersion new",
						"data_size":   "5096",
						"data_count":  "1001",
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

var AlicloudPaiWorkspaceDatasetversionMap8541 = map[string]string{
	"create_time":  CHECKSET,
	"version_name": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetversionBasicDependence8541(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultAiWorkspace" {
  description    = "dataset_pop_test_906t_211"
  display_name   = "DatasetVerResouce_924"
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_dataset" "defaultDataset" {
  options          = "{\"mountPath\":\"/mnt/data/\"}"
  description      = "镇元资源测试用例Dataset"
  accessibility    = "PRIVATE"
  dataset_name     = format("%%s1", var.name)
  data_source_type = "OSS"
  source_type      = "USER"
  workspace_id     = alicloud_pai_workspace_workspace.defaultAiWorkspace.id
  data_type        = "PIC"
  property         = "DIRECTORY"
  uri              = "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/"
  source_id        = "d-xxxxx_v1"
  user_id          = "1511928242963727"
}


`, name)
}

// Case DatasetVersion测试用例2 8551
func TestAccAliCloudPaiWorkspaceDatasetversion_basic8551(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_datasetversion.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetversionMap8551)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceDatasetversion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceDatasetversionBasicDependence8551)
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
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/version/\\\"}",
					"description":      "镇元资源测试用例DatasetVersion",
					"data_source_type": "OSS",
					"source_type":      "USER",
					"source_id":        "d-xxxxx_v1",
					"data_size":        "2068",
					"data_count":       "1000",
					"uri":              "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/catdog.json",
					"property":         "FILE",
					"dataset_id":       "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":          "{\"mountPath\":\"/mnt/data/version/\"}",
						"description":      "镇元资源测试用例DatasetVersion",
						"data_source_type": "OSS",
						"source_type":      "USER",
						"source_id":        "d-xxxxx_v1",
						"data_size":        "2068",
						"data_count":       "1000",
						"uri":              "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/catdog.json",
						"property":         "FILE",
						"dataset_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options":     "{\\\"mountPath\\\":\\\"/mnt/data/new/\\\"}",
					"description": "镇元资源测试用例DatasetVersion new",
					"data_size":   "5096",
					"data_count":  "1001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options":     "{\"mountPath\":\"/mnt/data/new/\"}",
						"description": "镇元资源测试用例DatasetVersion new",
						"data_size":   "5096",
						"data_count":  "1001",
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

var AlicloudPaiWorkspaceDatasetversionMap8551 = map[string]string{
	"create_time":  CHECKSET,
	"version_name": CHECKSET,
}

func AlicloudPaiWorkspaceDatasetversionBasicDependence8551(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultAiWorkspace" {
  description    = "dataset_pop_test_914t_482"
  display_name   = "DatasetVerResouce_100"
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_dataset" "defaultDataset" {
  options          = "{\"mountPath\":\"/mnt/data/\"}"
  description      = "镇元资源测试用例Dataset"
  accessibility    = "PRIVATE"
  dataset_name     = format("%%s1", var.name)
  data_source_type = "OSS"
  source_type      = "USER"
  workspace_id     = alicloud_pai_workspace_workspace.defaultAiWorkspace.id
  data_type        = "PIC"
  property         = "DIRECTORY"
  uri              = "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/"
  source_id        = "d-xxxxx_v1"
  user_id          = "1511928242963727"
}


`, name)
}

// Test PaiWorkspace Datasetversion. <<< Resource test cases, automatically generated.
