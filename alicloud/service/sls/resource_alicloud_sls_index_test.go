// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package sls_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Sls Index. >>> Resource test cases, automatically generated.
// Case index_terraform 10982
func TestAccAliCloudSlsIndex_basic10982(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsIndex")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            "true",
							"case_sensitive": "true",
							"token": []string{
								"a"},
							"exclude_keys": []string{
								"t"},
						},
					},
					"logstore_name": "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":  "${var.name}",
					"keys":          "{\\\"test\\\":{\\\"caseSensitive\\\":false,\\\"token\\\":[\\\"\\\\n\\\",\\\"\\\\t\\\",\\\",\\\",\\\" \\\",\\\";\\\",\\\"\\\\\\\"\\\",\\\"'\\\",\\\"(\\\",\\\")\\\",\\\"{\\\",\\\"}\\\",\\\"[\\\",\\\"]\\\",\\\"<\\\",\\\">\\\",\\\"?\\\",\\\"/\\\",\\\"#\\\",\\\":\\\"],\\\"type\\\":\\\"text\\\",\\\"doc_value\\\":false,\\\"alias\\\":\\\"\\\",\\\"chn\\\":false}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": acctest.CHECKSET,
						"project_name":  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            "false",
							"case_sensitive": "false",
							"token": []string{
								"tt"},
							"include_keys": []string{
								"tt"},
						},
					},
					"max_text_len": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_text_len": "500",
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

var AlicloudSlsIndexMap10982 = map[string]string{}

func AlicloudSlsIndexBasicDependence10982(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "logstore_name" {
  default = "logstore-test-%s"
}

variable "project_name" {
  default = "tf-testacc-sls-project-%s"
}

resource "alicloud_log_project" "defaultdCM1bA" {
  description = "terrafrom test"
  name        = var.name
}

resource "alicloud_log_store" "default7MW26R" {
  hot_ttl          = "7"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaultdCM1bA.project_name
  name             = var.logstore_name
}


`, name, name, name)
}

// Test Sls Index. <<< Resource test cases, automatically generated.

// TestAccAliCloudSlsIndex_keys_default_value verifies that SLS API defaults in
// an index key do not cause a perpetual diff.
func TestAccAliCloudSlsIndex_keys_default_value(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsIndex")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
    name        = "tf-testacc-log-project-%s"
    description = "Test project for SLS index keys default value problem"
}

resource "alicloud_log_store" "default" {
    project          = alicloud_log_project.default.project_name
    name             = "tf-testacc-log-store-%s"
    retention_period = "30"
    shard_count      = "2"
}

resource "alicloud_sls_index" "default" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = alicloud_log_store.default.logstore_name

    line {
        chn            = true
        case_sensitive = false
        token          = ["\n", "\t", ",", " ", ";"]
    }

    keys = jsonencode({
        "cost_ms" : {
            "doc_value" : true,
            "type" : "long"
        }
    })
}
`, name, name, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": acctest.CHECKSET,
						"project_name":  acctest.CHECKSET,
						"keys":          acctest.CHECKSET,
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
    name        = "tf-testacc-log-project-%s"
    description = "Test project for SLS index keys default value problem"
}

resource "alicloud_log_store" "default" {
    project          = alicloud_log_project.default.project_name
    name             = "tf-testacc-log-store-%s"
    retention_period = "30"
    shard_count      = "2"
}

resource "alicloud_sls_index" "default" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = alicloud_log_store.default.logstore_name

    line {
        chn            = true
        case_sensitive = false
        token          = ["\n", "\t", ",", " ", ";"]
    }

    keys = jsonencode({
        "cost_ms" : {
            "doc_value" : true,
            "type" : "long"
        }
    })
}
`, name, name, name),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccAliCloudSlsIndex_memberOrder_lineToken verifies convergent
// reordering of the nested TypeList `line.token`. Pattern B (per-list 2-step):
// step1 = full base with line[0] fully specified (token=[a,b], exclude_keys=[t,u],
// include_keys=[i,j], chn=true, case_sensitive=false) plus log_reduce_black_list=[x,y],
// log_reduce_white_list=[p,q], log_reduce=true, keys, logstore_name, project_name;
// step2 = the SAME full base with ONLY line.token reordered to [b,a] (everything
// else byte-identical, including the other four lists staying in order A).
func TestAccAliCloudSlsIndex_memberOrder_lineToken(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsIndex")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": acctest.CHECKSET,
						"project_name":  acctest.CHECKSET,
						"keys":          acctest.CHECKSET,
						"log_reduce":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"b", "a"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keys":       acctest.CHECKSET,
						"log_reduce": "true",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsIndex_memberOrder_lineExcludeKeys verifies convergent
// reordering of the nested TypeList `line.exclude_keys`. Pattern B: step1 = full
// base (exclude_keys=[t,u]); step2 = SAME full base with ONLY line.exclude_keys
// reordered to [u,t] (everything else byte-identical).
func TestAccAliCloudSlsIndex_memberOrder_lineExcludeKeys(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsIndex")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": acctest.CHECKSET,
						"project_name":  acctest.CHECKSET,
						"keys":          acctest.CHECKSET,
						"log_reduce":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"u", "t"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keys":       acctest.CHECKSET,
						"log_reduce": "true",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsIndex_memberOrder_lineIncludeKeys verifies convergent
// reordering of the nested TypeList `line.include_keys`. Pattern B: step1 = full
// base (include_keys=[i,j]); step2 = SAME full base with ONLY line.include_keys
// reordered to [j,i] (everything else byte-identical).
func TestAccAliCloudSlsIndex_memberOrder_lineIncludeKeys(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsIndex")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": acctest.CHECKSET,
						"project_name":  acctest.CHECKSET,
						"keys":          acctest.CHECKSET,
						"log_reduce":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"j", "i"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keys":       acctest.CHECKSET,
						"log_reduce": "true",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsIndex_memberOrder_logReduceBlackList verifies convergent
// reordering of the top-level TypeList `log_reduce_black_list`. Pattern B:
// step1 = full base (log_reduce_black_list=[x,y]); step2 = SAME full base with
// ONLY log_reduce_black_list reordered to [y,x] (everything else byte-identical).
func TestAccAliCloudSlsIndex_memberOrder_logReduceBlackList(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsIndex")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": acctest.CHECKSET,
						"project_name":  acctest.CHECKSET,
						"keys":          acctest.CHECKSET,
						"log_reduce":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"y", "x"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keys":       acctest.CHECKSET,
						"log_reduce": "true",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsIndex_memberOrder_logReduceWhiteList verifies convergent
// reordering of the top-level TypeList `log_reduce_white_list`. Pattern B:
// step1 = full base (log_reduce_white_list=[p,q]); step2 = SAME full base with
// ONLY log_reduce_white_list reordered to [q,p] (everything else byte-identical).
func TestAccAliCloudSlsIndex_memberOrder_logReduceWhiteList(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsIndex")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"p", "q"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": acctest.CHECKSET,
						"project_name":  acctest.CHECKSET,
						"keys":          acctest.CHECKSET,
						"log_reduce":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a", "b"},
							"exclude_keys":   []string{"t", "u"},
							"include_keys":   []string{"i", "j"},
						},
					},
					"log_reduce":            true,
					"log_reduce_black_list": []string{"x", "y"},
					"log_reduce_white_list": []string{"q", "p"},
					"logstore_name":         "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":          "${var.name}",
					"keys":                  "{\"k\":{\"type\":\"text\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keys":       acctest.CHECKSET,
						"log_reduce": "true",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsIndex_keysAndLogReduce_update verifies modification of
// the keys (JSON-encoded string) and log_reduce (bool) attributes across two
// config steps in the same function.
func TestAccAliCloudSlsIndex_keysAndLogReduce_update(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_index.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsIndexMap10982)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsIndex")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsIndexBasicDependence10982)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"line": []map[string]interface{}{
						{
							"chn":            true,
							"case_sensitive": false,
							"token":          []string{"a"},
						},
					},
					"logstore_name": "${alicloud_log_store.default7MW26R.logstore_name}",
					"project_name":  "${var.name}",
					"keys":          "{\"k1\":{\"type\":\"text\"}}",
					"log_reduce":    true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": acctest.CHECKSET,
						"project_name":  acctest.CHECKSET,
						"keys":          acctest.CHECKSET,
						"log_reduce":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"keys":       "{\"k2\":{\"type\":\"long\"}}",
					"log_reduce": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keys":       acctest.CHECKSET,
						"log_reduce": "false",
					}),
				),
			},
		},
	})
}
