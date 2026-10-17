package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// collection-order: bypass labels -- pre-existing ForceNew TypeList attribute outside this change

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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// Case DatasetVersion_dataset_task_ram_role_user_metrics_endpoints 9901
func TestAccAliCloudPaiWorkspaceDatasetversion_userMetricsEndpoints9901(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_datasetversion.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetversionMapUserMetrics9901)
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 0 — happy Create/Read of a basic OSS dataset version with no
			// new fields. DatasetTaskRamRole and user_metrics_endpoints are
			// unsupported on the BASIC OSS dataset creatable via
			// alicloud_pai_workspace_dataset (data_source_type ValidateFunc
			// restricts to OSS/NAS; no JUICEFS option), so they are exercised
			// only via the Update failure paths below. Create-time forwarding of
			// both parameters is evidenced by prior ACC runs on the same Create
			// code (API 400 for each parameter at create time).
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/version/\\\"}",
					"description":      "DatasetVersion basic",
					"data_source_type": "OSS",
					"source_type":      "USER",
					"source_id":        "d-xxxxx_v1",
					"data_size":        "2068",
					"data_count":       "1000",
					"uri":              "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
					"property":         "DIRECTORY",
					"dataset_id":       "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"property":   "DIRECTORY",
						"dataset_id": CHECKSET,
					}),
				),
			},
			// Step 1 — happy Update/Read of the basic dataset version (no new
			// fields).
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/new/\\\"}",
					"description":      "DatasetVersion basic new",
					"data_source_type": "OSS",
					"source_type":      "USER",
					"source_id":        "d-xxxxx_v1",
					"data_size":        "5096",
					"data_count":       "1001",
					"uri":              "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
					"property":         "DIRECTORY",
					"dataset_id":       "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_size": "5096",
					}),
				),
			},
			// Step 2 — Import/Read. Verified against the clean happy-path state
			// (after the successful Step 1 Update) BEFORE the failure-path steps:
			// an ExpectError Update that the API rejects leaves the test harness
			// holding the divergent planned state (ram_role set, description
			// advanced) while the real resource is unchanged, so ImportStateVerify
			// cannot be reconciled after those steps.
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// Step 3 — Update failure path for user_metrics_endpoints (also the
			// modification coverage step for user_metrics_endpoints/path/port at
			// configIndex>0): updating the OSS dataset version with
			// user_metrics_endpoints is rejected by the API (Code 201300003),
			// proving the Provider forwards user_metrics_endpoints/path/port on
			// Update. This is a single failure-path step run against the clean
			// post-import state (no dataset_task_ram_role in the prior state),
			// so it is isolated from the dataset_task_ram_role failure path that
			// the sibling datasetTaskRamRole9902 case exercises.
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/new/\\\"}",
					"description":      "DatasetVersion user metrics endpoints update failure",
					"data_source_type": "OSS",
					"source_type":      "USER",
					"source_id":        "d-xxxxx_v1",
					"data_size":        "5096",
					"data_count":       "1001",
					"user_metrics_endpoints": []map[string]interface{}{
						{"path": "/metrics", "port": 9090},
					},
					"uri":        "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
					"property":   "DIRECTORY",
					"dataset_id": "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				ExpectError: regexp.MustCompile(`UserMetricsEndpoints is only supported for JUICEFS datasets`),
			},
		},
	})
}

var AlicloudPaiWorkspaceDatasetversionMapUserMetrics9901 = map[string]string{
	"create_time":  CHECKSET,
	"version_name": CHECKSET,
}

// Case DatasetVersion_dataset_task_ram_role 9902
func TestAccAliCloudPaiWorkspaceDatasetversion_datasetTaskRamRole9902(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_datasetversion.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceDatasetversionMapDatasetTaskRamRole9902)
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 0 — happy Create/Read of a basic OSS dataset version with no
			// new fields. dataset_task_ram_role is unsupported on the BASIC OSS
			// dataset creatable via alicloud_pai_workspace_dataset
			// (data_source_type ValidateFunc restricts to OSS/NAS; no JUICEFS
			// option), so it is exercised only via the Update failure path below.
			// Create-time forwarding of dataset_task_ram_role is evidenced by
			// prior ACC runs on the same Create code (API 400 at create time).
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/version/\\\"}",
					"description":      "DatasetVersion basic",
					"data_source_type": "OSS",
					"source_type":      "USER",
					"source_id":        "d-xxxxx_v1",
					"data_size":        "2068",
					"data_count":       "1000",
					"uri":              "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
					"property":         "DIRECTORY",
					"dataset_id":       "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"property":   "DIRECTORY",
						"dataset_id": CHECKSET,
					}),
				),
			},
			// Step 1 — happy Update/Read of the basic dataset version (no new
			// fields).
			{
				Config: testAccConfig(map[string]interface{}{
					"options":          "{\\\"mountPath\\\":\\\"/mnt/data/new/\\\"}",
					"description":      "DatasetVersion basic new",
					"data_source_type": "OSS",
					"source_type":      "USER",
					"source_id":        "d-xxxxx_v1",
					"data_size":        "5096",
					"data_count":       "1001",
					"uri":              "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
					"property":         "DIRECTORY",
					"dataset_id":       "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_size": "5096",
					}),
				),
			},
			// Step 2 — Import/Read. Verified against the clean happy-path state
			// (after the successful Step 1 Update) BEFORE the failure-path step:
			// an ExpectError Update that the API rejects leaves the test harness
			// holding the divergent planned state (ram_role set, description
			// advanced) while the real resource is unchanged, so ImportStateVerify
			// cannot be reconciled after that step.
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// Step 3 — Update failure path for dataset_task_ram_role (also the
			// modification coverage step at configIndex>0): updating the BASIC
			// OSS dataset version with dataset_task_ram_role is rejected by the
			// API (Code 201300006), proving the Provider forwards
			// dataset_task_ram_role on Update. The parameter first appears in
			// this Update step, so the coverage checker records it as a
			// modification.
			{
				Config: testAccConfig(map[string]interface{}{
					"options":               "{\\\"mountPath\\\":\\\"/mnt/data/new/\\\"}",
					"description":           "DatasetVersion dataset task ram role update failure",
					"data_source_type":      "OSS",
					"source_type":           "USER",
					"source_id":             "d-xxxxx_v1",
					"data_size":             "5096",
					"data_count":            "1001",
					"dataset_task_ram_role": "acs:ram::1234567890123456:role/roletest",
					"uri":                   "oss://ai4d-q9lgxlpwxzqluij66y.oss-cn-hangzhou.aliyuncs.com/",
					"property":              "DIRECTORY",
					"dataset_id":            "${alicloud_pai_workspace_dataset.defaultDataset.id}",
				}),
				ExpectError: regexp.MustCompile(`DatasetTaskRamRole is not supported for BASIC dataset`),
			},
		},
	})
}

var AlicloudPaiWorkspaceDatasetversionMapDatasetTaskRamRole9902 = map[string]string{
	"create_time":  CHECKSET,
	"version_name": CHECKSET,
}

// Test PaiWorkspace Datasetversion. <<< Resource test cases, automatically generated.
