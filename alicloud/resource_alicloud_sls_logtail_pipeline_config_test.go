// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls LogtailPipelineConfig. >>> Resource test cases, automatically generated.
// Case LogtailPipelineConfigTestPL 12633
func TestAccAliCloudSlsLogtailPipelineConfig_basic12633(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_logtail_pipeline_config.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsLogtailPipelineConfigMap12633)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsLogtailPipelineConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	projectName := fmt.Sprintf("tf-logtail-test-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
    Region        = "eu-central-1"
    Endpoint      = "eu-central-1-intranet.log.aliyuncs.com"
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
    Region        = "eu-central-1"
    Endpoint      = "eu-central-1-intranet.log.aliyuncs.com"
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
    Region        = "eu-central-1"
    Endpoint      = "eu-central-1-intranet.log.aliyuncs.com"
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
    Region        = "eu-central-1"
    Endpoint      = "eu-central-1-intranet.log.aliyuncs.com"
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
    Region        = "eu-central-1"
    Endpoint      = "eu-central-1-intranet.log.aliyuncs.com"
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
    Region        = "eu-central-1"
    Endpoint      = "eu-central-1-intranet.log.aliyuncs.com"
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
    Region        = "eu-central-1"
    Endpoint      = "eu-central-1-intranet.log.aliyuncs.com"
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudSlsLogtailPipelineConfigMap12633 = map[string]string{}

// Test Sls LogtailPipelineConfig. <<< Resource test cases, automatically generated.
