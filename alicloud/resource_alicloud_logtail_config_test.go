package alicloud

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"
)

func TestAccAlicloudLogTail_basic(t *testing.T) {
	var project sls.LogProject
	var store sls.LogStore
	var config sls.LogConfig
	tailbasic_input_detail := "{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}"
	tailbasic_input_detail_plugin := "{\"plugin\":{\"inputs\":[{\"detail\":{\"ExcludeEnv\":null,\"ExcludeLabel\":null,\"IncludeEnv\":null,\"IncludeLabel\":null,\"Stderr\":true,\"Stdout\":true},\"type\":\"service_docker_stdout\"}]}}"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogTailConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogTailbasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.example", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.example", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.example", &config),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "input_type", "file"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "log_sample", "test"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "name", "evan-terraform-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "output_type", "LogService"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "project", "test-tf"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "logstore", "tf-test-logstore"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "input_detail",tailbasic_input_detail),
				),
			},
			{
				Config: testAlicloudLogTailUpdateOneParamater,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.update", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.update", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.update", &config),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "input_type", "file"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "log_sample", "test-logtail-update"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "name", "evan-terraform-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "output_type", "LogService"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "project", "test-tf2"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "logstore", "tf-test-logstore"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "input_detail",tailbasic_input_detail),
				),
			},
			{
				Config: testAlicloudLogTailUpdateAllParamater,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.update_all", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.update_all", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.update_all", &config),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "input_type", "plugin"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "log_sample", "test-logtail-update-all-paramter"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "name", "tf-terraform-plugin-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "output_type", "LogService"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "project", "test-tf3"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "logstore", "tf-test-logstore3"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "input_detail",tailbasic_input_detail_plugin),
				),
			},
		},
	})
}

func testAccCheckAlicloudLogTailConfigExists(name string, config *sls.LogConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Logtail config ID is set")
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		logService := LogService{client}
		logconfig, err := logService.DescribeLogLogtailConfig(split[0], split[2])
		if err != nil {
			return err
		}
		if logconfig == nil || logconfig.Name == "" {
			return fmt.Errorf("LogConfig %s is not exist.", split[2])
		}
		config = logconfig
		return nil
	}
}

func testAccCheckAlicloudLogTailConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	logService := LogService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_logtail_config" {
			continue
		}
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		_, err := logService.DescribeLogLogtailConfig(split[0], split[2])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Check logtail config got an error: %#v.", err)
		}
		return fmt.Errorf("Logtail config %s still exists.", split[2])
	}
	return nil
}

const testAlicloudLogTailbasic = `
resource "alicloud_log_project" "example"{
	name = "test-tf"
	description = "create by terraform"
}
resource "alicloud_log_store" "example"{
  	project = "${alicloud_log_project.example.name}"
  	name = "tf-test-logstore"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_logtail_config" "example"{
	project = "${alicloud_log_project.example.name}"
  	logstore = "${alicloud_log_store.example.name}"
  	input_type = "file"
  	log_sample = "test"
  	name = "evan-terraform-config"
	output_type = "LogService"
  	input_detail = <<DEFINITION
  	{
		"logPath": "/logPath",
		"filePattern": "access.log",
		"logType": "json_log",
		"topicFormat": "default",
		"discardUnmatch": false,
		"enableRawLog": true,
		"fileEncoding": "gbk",
		"maxDepth": 10
	}
	DEFINITION
}
`
const testAlicloudLogTailUpdateOneParamater = `
resource "alicloud_log_project" "update"{
	name = "test-tf2"
	description = "create by terraform"
}
resource "alicloud_log_store" "update"{
  	project = "${alicloud_log_project.update.name}"
  	name = "tf-test-logstore"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_logtail_config" "update"{
	project = "${alicloud_log_project.update.name}"
  	logstore = "${alicloud_log_store.update.name}"
  	input_type = "file"
  	log_sample = "test-logtail-update"
  	name = "evan-terraform-config"
	output_type = "LogService"
  	input_detail = <<DEFINITION
  	{
		"logPath": "/logPath",
		"filePattern": "access.log",
		"logType": "json_log",
		"topicFormat": "default",
		"discardUnmatch": false,
		"enableRawLog": true,
		"fileEncoding": "gbk",
		"maxDepth": 10
	}
	DEFINITION
}
`
const testAlicloudLogTailUpdateAllParamater = `
resource "alicloud_log_project" "update_all"{
	name = "test-tf3"
	description = "create by terraform"
}
resource "alicloud_log_store" "update_all"{
  	project = "${alicloud_log_project.update_all.name}"
  	name = "tf-test-logstore3"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_logtail_config" "update_all"{
	project = "${alicloud_log_project.update_all.name}"
  	logstore = "${alicloud_log_store.update_all.name}"
  	input_type = "plugin"
  	log_sample = "test-logtail-update-all-paramter"
  	name = "tf-terraform-plugin-config"
	output_type = "LogService"
  	input_detail = <<DEFINITION
  	{
		"plugin": {
            "inputs": [
                {
                    "detail": {
                        "ExcludeEnv": null, 
                        "ExcludeLabel": null, 
                        "IncludeEnv": null, 
                        "IncludeLabel": null, 
                        "Stderr": true, 
                        "Stdout": true
                    }, 
                    "type": "service_docker_stdout"
                }
            ]
        }
	}
	DEFINITION
}
`

