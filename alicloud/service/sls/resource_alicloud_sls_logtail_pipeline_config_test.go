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

// Test Sls LogtailPipelineConfig. >>> Resource test cases, automatically generated.
// Case LogtailPipelineConfigTestPL 12633
func TestAccAliCloudSlsLogtailPipelineConfig_basic12633(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_logtail_pipeline_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsLogtailPipelineConfigMap12633)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsLogtailPipelineConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	projectName := fmt.Sprintf("tf-logtail-test-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 0: Create resource with all required and optional fields (inputs, processors, flushers, aggregators)
			{
				Config: fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%s"
  description  = "terraform logtail pipeline config test"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.project_name
  name                  = "test"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_log_store" "default_updated" {
  project               = alicloud_log_project.default.project_name
  name                  = "test-updated"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "default" {
  project      = alicloud_log_project.default.project_name
  config_name  = "%s"
  
  inputs = [{
    Type                     = "input_file"
    FilePaths                = "[\"/home/*.log\"]"
    EnableContainerDiscovery = false
    MaxDirSearchDepth        = 0
    FileEncoding             = "utf8"
  }]
  
  processors = [{
    Type      = "processor_parse_regex_native"
    SourceKey = "content"
    Regex     = ".*"
    Keys      = "[\"key1\",\"key2\"]"
  }]
  
  flushers = [{
    Type          = "flusher_sls"
    Logstore      = alicloud_log_store.default.name
    TelemetryType = "logs"

  }]
  
  aggregators = [{
    Type           = "aggregator_default"
    MaxSizeBytes   = 1048576
    MaxTimeSeconds = 3
  }]
}
`, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":       projectName,
						"config_name":   name,
						"inputs.#":      "1",
						"processors.#":  "1",
						"flushers.#":    "1",
						"aggregators.#": "1",
					}),
				),
			},
			// Step 1: Update processors (modify Regex and Keys)
			{
				Config: fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%s"
  description  = "terraform logtail pipeline config test"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.project_name
  name                  = "test"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_log_store" "default_updated" {
  project               = alicloud_log_project.default.project_name
  name                  = "test-updated"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "default" {
  project      = alicloud_log_project.default.project_name
  config_name  = "%s"
  
  inputs = [{
    Type                     = "input_file"
    FilePaths                = "[\"/home/*.log\"]"
    EnableContainerDiscovery = false
    MaxDirSearchDepth        = 0
    FileEncoding             = "utf8"
  }]
  
  processors = [{
    Type      = "processor_parse_regex_native"
    SourceKey = "content"
    Regex     = "([\\d\\.]+) \\S+ \\S+ \\[(\\S+) \\S+\\].*"
    Keys      = "[\"ip\",\"time\"]"
  }]
  
  flushers = [{
    Type          = "flusher_sls"
    Logstore      = alicloud_log_store.default.name
    TelemetryType = "logs"

  }]
  
  aggregators = [{
    Type           = "aggregator_default"
    MaxSizeBytes   = 1048576
    MaxTimeSeconds = 3
  }]
}
`, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"processors.#": "1",
					}),
				),
			},
			// Step 2: Update flushers (modify Logstore)
			{
				Config: fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%s"
  description  = "terraform logtail pipeline config test"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.project_name
  name                  = "test"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_log_store" "default_updated" {
  project               = alicloud_log_project.default.project_name
  name                  = "test-updated"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "default" {
  project      = alicloud_log_project.default.project_name
  config_name  = "%s"
  
  inputs = [{
    Type                     = "input_file"
    FilePaths                = "[\"/home/*.log\"]"
    EnableContainerDiscovery = false
    MaxDirSearchDepth        = 0
    FileEncoding             = "utf8"
  }]
  
  processors = [{
    Type      = "processor_parse_regex_native"
    SourceKey = "content"
    Regex     = "([\\d\\.]+) \\S+ \\S+ \\[(\\S+) \\S+\\].*"
    Keys      = "[\"ip\",\"time\"]"
  }]
  
  flushers = [{
    Type          = "flusher_sls"
    Logstore      = alicloud_log_store.default_updated.name
    TelemetryType = "logs"

  }]
  
  aggregators = [{
    Type           = "aggregator_default"
    MaxSizeBytes   = 1048576
    MaxTimeSeconds = 3
  }]
}
`, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flushers.#": "1",
					}),
				),
			},
			// Step 3: Update inputs (modify FilePaths and MaxDirSearchDepth)
			{
				Config: fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%s"
  description  = "terraform logtail pipeline config test"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.project_name
  name                  = "test"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_log_store" "default_updated" {
  project               = alicloud_log_project.default.project_name
  name                  = "test-updated"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "default" {
  project      = alicloud_log_project.default.project_name
  config_name  = "%s"
  
  inputs = [{
    Type                     = "input_file"
    FilePaths                = "[\"/var/log/*.log\"]"
    EnableContainerDiscovery = false
    MaxDirSearchDepth        = 1
    FileEncoding             = "utf8"
  }]
  
  processors = [{
    Type      = "processor_parse_regex_native"
    SourceKey = "content"
    Regex     = "([\\d\\.]+) \\S+ \\S+ \\[(\\S+) \\S+\\].*"
    Keys      = "[\"ip\",\"time\"]"
  }]
  
  flushers = [{
    Type          = "flusher_sls"
    Logstore      = alicloud_log_store.default_updated.name
    TelemetryType = "logs"

  }]
  
  aggregators = [{
    Type           = "aggregator_default"
    MaxSizeBytes   = 1048576
    MaxTimeSeconds = 3
  }]
}
`, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inputs.#": "1",
					}),
				),
			},
			// Step 4: Add log_sample field
			{
				Config: fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%s"
  description  = "terraform logtail pipeline config test"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.project_name
  name                  = "test"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_log_store" "default_updated" {
  project               = alicloud_log_project.default.project_name
  name                  = "test-updated"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "default" {
  project      = alicloud_log_project.default.project_name
  config_name  = "%s"
  log_sample   = "sample log for testing"
  
  inputs = [{
    Type                     = "input_file"
    FilePaths                = "[\"/var/log/*.log\"]"
    EnableContainerDiscovery = false
    MaxDirSearchDepth        = 1
    FileEncoding             = "utf8"
  }]
  
  processors = [{
    Type      = "processor_parse_regex_native"
    SourceKey = "content"
    Regex     = "([\\d\\.]+) \\S+ \\S+ \\[(\\S+) \\S+\\].*"
    Keys      = "[\"ip\",\"time\"]"
  }]
  
  flushers = [{
    Type          = "flusher_sls"
    Logstore      = alicloud_log_store.default_updated.name
    TelemetryType = "logs"

  }]
  
  aggregators = [{
    Type           = "aggregator_default"
    MaxSizeBytes   = 1048576
    MaxTimeSeconds = 3
  }]
}
`, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_sample": "sample log for testing",
					}),
				),
			},
			// Step 5: Add globals field
			{
				Config: fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%s"
  description  = "terraform logtail pipeline config test"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.project_name
  name                  = "test"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_log_store" "default_updated" {
  project               = alicloud_log_project.default.project_name
  name                  = "test-updated"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "default" {
  project      = alicloud_log_project.default.project_name
  config_name  = "%s"
  log_sample   = "sample log for testing"
  
  globals = {
    EnableTimestampNanosecond = "true"
    UsingOldContentTag        = "false"
  }
  
  inputs = [{
    Type                     = "input_file"
    FilePaths                = "[\"/var/log/*.log\"]"
    EnableContainerDiscovery = false
    MaxDirSearchDepth        = 1
    FileEncoding             = "utf8"
  }]
  
  processors = [{
    Type      = "processor_parse_regex_native"
    SourceKey = "content"
    Regex     = "([\\d\\.]+) \\S+ \\S+ \\[(\\S+) \\S+\\].*"
    Keys      = "[\"ip\",\"time\"]"
  }]
  
  flushers = [{
    Type          = "flusher_sls"
    Logstore      = alicloud_log_store.default_updated.name
    TelemetryType = "logs"

  }]
  
  aggregators = [{
    Type           = "aggregator_default"
    MaxSizeBytes   = 1048576
    MaxTimeSeconds = 3
  }]
}
`, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"globals.%": "2",
					}),
				),
			},
			// Step 6: Update multiple fields (log_sample, globals, aggregators)
			{
				Config: fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%s"
  description  = "terraform logtail pipeline config test"
}

resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.project_name
  name                  = "test"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_log_store" "default_updated" {
  project               = alicloud_log_project.default.project_name
  name                  = "test-updated"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "default" {
  project      = alicloud_log_project.default.project_name
  config_name  = "%s"
  log_sample   = "updated sample log"
  
  globals = {
    EnableTimestampNanosecond = "false"
    UsingOldContentTag        = "true"
  }
  
  inputs = [{
    Type                     = "input_file"
    FilePaths                = "[\"/var/log/*.log\"]"
    EnableContainerDiscovery = false
    MaxDirSearchDepth        = 1
    FileEncoding             = "utf8"
  }]
  
  processors = [{
    Type      = "processor_parse_regex_native"
    SourceKey = "content"
    Regex     = "([\\d\\.]+) \\S+ \\S+ \\[(\\S+) \\S+\\].*"
    Keys      = "[\"ip\",\"time\"]"
  }]
  
  flushers = [{
    Type          = "flusher_sls"
    Logstore      = alicloud_log_store.default_updated.name
    TelemetryType = "logs"

  }]
  
  aggregators = [{
    Type           = "aggregator_default"
    MaxSizeBytes   = 2097152
    MaxTimeSeconds = 5
  }]
}
`, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_sample":    "updated sample log",
						"aggregators.#": "1",
					}),
				),
			},
			// Step 7: ImportState test to verify import functionality
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"flushers"},
			},
		},
	})
}

var AlicloudSlsLogtailPipelineConfigMap12633 = map[string]string{}

// TestAccAliCloudSlsLogtailPipelineConfig_memberOrder covers the TypeList Order
// Coverage requirement for the aggregators, flushers, inputs and processors
// attributes. Each listed attribute is a TypeList whose elements are TypeMap
// blocks (pipeline config blocks).
//
// Step0 configures every listed list with >=2 distinct map members in order A,
// together with all required creation attributes (project, config_name) and the
// non-list Optional attributes (globals, log_sample, task). Steps 1-4 each
// reorder exactly ONE list (a delta specifying only that attr's permuted
// members; the builder inherits everything else), so each consecutive pair
// covers one list. Step5 modifies globals/log_sample/task only (the reordered
// lists are inherited), putting those non-list Optional attributes into the
// modification coverage set.
func TestAccAliCloudSlsLogtailPipelineConfig_memberOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_logtail_pipeline_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsLogtailPipelineConfigMap12633)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsLogtailPipelineConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsLogtailPipelineConfigBasicDependence12633)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 0: base. Every candidate TypeList (aggregators, flushers,
			// inputs, processors) carries >=2 distinct map members in order A,
			// plus all Required attrs (project, config_name) and the non-list
			// Optional attrs (globals, log_sample, task) set so they enter the
			// modification coverage set.
			{
				Config: testAccConfig(map[string]interface{}{
					"project":     "${alicloud_log_project.default.project_name}",
					"config_name": name,
					"globals": map[string]interface{}{
						"EnableTimestampNanosecond": "true",
					},
					"log_sample": "sample log v1",
					"task": map[string]interface{}{
						"Enable": "true",
					},
					"inputs": []map[string]interface{}{
						{
							"Type":                     "input_file",
							"FilePaths":                "[\"/home/*.log\"]",
							"EnableContainerDiscovery": "false",
							"MaxDirSearchDepth":        "0",
							"FileEncoding":             "utf8",
						},
						{
							"Type":                     "input_file",
							"FilePaths":                "[\"/var/log/*.log\"]",
							"EnableContainerDiscovery": "true",
							"MaxDirSearchDepth":        "1",
							"FileEncoding":             "utf8",
						},
					},
					"processors": []map[string]interface{}{
						{
							"Type":      "processor_parse_regex_native",
							"SourceKey": "content",
							"Regex":     ".*",
							"Keys":      "[\"key1\",\"key2\"]",
						},
						{
							"Type":      "processor_parse_regex_native",
							"SourceKey": "message",
							"Regex":     ".*",
							"Keys":      "[\"ip\",\"time\"]",
						},
					},
					"flushers": []map[string]interface{}{
						{
							"Type":          "flusher_sls",
							"Logstore":      "test",
							"TelemetryType": "logs",
						},
						{
							"Type":          "flusher_sls",
							"Logstore":      "test-updated",
							"TelemetryType": "logs",
						},
					},
					"aggregators": []map[string]interface{}{
						{
							"Type":           "aggregator_default",
							"MaxSizeBytes":   "1048576",
							"MaxTimeSeconds": "3",
						},
						{
							"Type":           "aggregator_default",
							"MaxSizeBytes":   "2097152",
							"MaxTimeSeconds": "5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":       acctest.CHECKSET,
						"config_name":   name,
						"inputs.#":      "2",
						"processors.#":  "2",
						"flushers.#":    "2",
						"aggregators.#": "2",
					}),
				),
			},
			// Step 1: reorder aggregators only (delta specifies aggregators; the
			// builder inherits the rest). Pair step0->step1 covers aggregators.
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregators": []map[string]interface{}{
						{
							"Type":           "aggregator_default",
							"MaxSizeBytes":   "2097152",
							"MaxTimeSeconds": "5",
						},
						{
							"Type":           "aggregator_default",
							"MaxSizeBytes":   "1048576",
							"MaxTimeSeconds": "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregators.#": "2",
					}),
				),
			},
			// Step 2: reorder flushers only. Pair step1->step2 covers flushers.
			{
				Config: testAccConfig(map[string]interface{}{
					"flushers": []map[string]interface{}{
						{
							"Type":          "flusher_sls",
							"Logstore":      "test-updated",
							"TelemetryType": "logs",
						},
						{
							"Type":          "flusher_sls",
							"Logstore":      "test",
							"TelemetryType": "logs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flushers.#": "2",
					}),
				),
			},
			// Step 3: reorder inputs only. Pair step2->step3 covers inputs.
			{
				Config: testAccConfig(map[string]interface{}{
					"inputs": []map[string]interface{}{
						{
							"Type":                     "input_file",
							"FilePaths":                "[\"/var/log/*.log\"]",
							"EnableContainerDiscovery": "true",
							"MaxDirSearchDepth":        "1",
							"FileEncoding":             "utf8",
						},
						{
							"Type":                     "input_file",
							"FilePaths":                "[\"/home/*.log\"]",
							"EnableContainerDiscovery": "false",
							"MaxDirSearchDepth":        "0",
							"FileEncoding":             "utf8",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inputs.#": "2",
					}),
				),
			},
			// Step 4: reorder processors only. Pair step3->step4 covers processors.
			{
				Config: testAccConfig(map[string]interface{}{
					"processors": []map[string]interface{}{
						{
							"Type":      "processor_parse_regex_native",
							"SourceKey": "message",
							"Regex":     ".*",
							"Keys":      "[\"ip\",\"time\"]",
						},
						{
							"Type":      "processor_parse_regex_native",
							"SourceKey": "content",
							"Regex":     ".*",
							"Keys":      "[\"key1\",\"key2\"]",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"processors.#": "2",
					}),
				),
			},
			// Step 5: modify globals/log_sample/task only (delta specifies these
			// non-list attrs; the builder inherits the reordered lists). This pair
			// changes no list, but puts globals/log_sample/task into testModifySet.
			{
				Config: testAccConfig(map[string]interface{}{
					"globals": map[string]interface{}{
						"EnableTimestampNanosecond": "false",
					},
					"log_sample": "sample log v2",
					"task": map[string]interface{}{
						"Enable": "false",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_sample": "sample log v2",
					}),
				),
			},
		},
	})
}

func AlicloudSlsLogtailPipelineConfigBasicDependence12633(name string) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%s"
  description  = "terraform logtail pipeline config member order test"
}

resource "alicloud_log_store" "default" {
  project                = alicloud_log_project.default.project_name
  name                   = "test"
  shard_count            = 2
  auto_split             = true
  max_split_shard_count  = 64
}

resource "alicloud_log_store" "default_updated" {
  project                = alicloud_log_project.default.project_name
  name                   = "test-updated"
  shard_count            = 2
  auto_split             = true
  max_split_shard_count  = 64
}
`, name)
}

// Test Sls LogtailPipelineConfig. <<< Resource test cases, automatically generated.
