package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_logtail_config", &resource.Sweeper{
		Name: "alicloud_logtail_config",
		F:    testSweepLogConfigs,
	})
}

func testSweepLogConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.ListProject()
	})
	if err != nil {
		return fmt.Errorf("Error retrieving Log Projects: %s", err)
	}
	names, _ := raw.([]string)

	for _, v := range names {
		name := v
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				cf_name_list, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					cf_names, _, cf_err := slsClient.ListConfig(name, 0, 100)
					return cf_names, cf_err
				})
				if err != nil {
					return fmt.Errorf("Error retrieving Log config: %s", err)
				}
				for _, cf_name := range cf_name_list.([]string) {
					log.Printf("[INFO] Deleting Log config: %s", cf_name)
					_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
						return nil, slsClient.DeleteConfig(name, cf_name)
					})
					if err != nil {
						log.Printf("[ERROR] Failed to delete Log Config (%s): %s", cf_name, err)
					}
				}
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Log Project: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting Log Project: %s", name)
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Log Project (%s): %s", name, err)
		}
	}
	return nil
}

func TestAccAlicloudLogTail_basic(t *testing.T) {
	var project sls.LogProject
	var store sls.LogStore
	var config sls.LogConfig
	tailbasic_input_detail := "{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}"
	tailbasic_input_detail_plugin := "{\"plugin\":{\"inputs\":[{\"detail\":{\"ExcludeEnv\":null,\"ExcludeLabel\":null,\"IncludeEnv\":null,\"IncludeLabel\":null,\"Stderr\":true,\"Stdout\":true},\"type\":\"service_docker_stdout\"}]}}"
	tailbasic_input_delimiter := "{\"autoExtend\":true,\"discardUnmatch\":true,\"enableRawLog\":true,\"fileEncoding\":\"utf8\",\"filePattern\":\"*\",\"key\":[\"test\",\"test2\"],\"logPath\":\"/logPath\",\"logType\":\"delimiter_log\",\"maxDepth\":999,\"quote\":\"\\\"\",\"separator\":\",\",\"timekey\":\"test\",\"topicFormat\":\"default\"}"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogTailConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudLogTailbasic(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.example", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.example", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.example", &config),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "input_type", "file"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "log_sample", "test"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "name", "tf-testacclogtailbasic-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "output_type", "LogService"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "logstore", "tf-testacclogtailbasic-logstore"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.example", "input_detail", tailbasic_input_detail),
				),
			},
			{
				Config: testAlicloudLogTailUpdateOneParamater(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.update", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.update", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.update", &config),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "input_type", "file"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "log_sample", "test-logtail-update"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "name", "tf-testacclogtailupdate-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "output_type", "LogService"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "logstore", "tf-testacclogtailupdate-logstore"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update", "input_detail", tailbasic_input_detail),
				),
			},
			{
				Config: testAlicloudLogTailUpdateAllParamater(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.update_all", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.update_all", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.update_all", &config),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "input_type", "plugin"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "log_sample", "test-logtail-update-all-paramter"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "name", "tf-testacclogtailupdateall-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "output_type", "LogService"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "logstore", "tf-testacclogtailupdateall-logstore"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.update_all", "input_detail", tailbasic_input_detail_plugin),
				),
			},
			{
				Config: testAlicloudLogTailUpdateInputDetail(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.delimiter", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.delimiter", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.delimiter", &config),
					resource.TestCheckResourceAttr("alicloud_logtail_config.delimiter", "input_type", "file"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.delimiter", "log_sample", "test-logtail-delimiter"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.delimiter", "name", "tf-testacclogtaildelimiter-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.delimiter", "output_type", "LogService"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.delimiter", "logstore", "tf-testacclogtaildelimiter-logstore"),
					resource.TestCheckResourceAttr("alicloud_logtail_config.delimiter", "input_detail", tailbasic_input_delimiter),
				),
			},
		},
	})
}

func testAccCheckAlicloudLogTailConfigExists(name string, config *sls.LogConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", name))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No Logtail config ID is set"))
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		logService := LogService{client}
		logconfig, err := logService.DescribeLogLogtailConfig(split[0], split[2])
		if err != nil {
			return WrapError(err)
		}
		if logconfig == nil || logconfig.Name == "" {
			return WrapError(fmt.Errorf("LogConfig %s is not exist.", split[2]))
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
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("Logtail config %s still exists.", split[2]))
	}
	return nil
}

func testAlicloudLogTailbasic(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "example"{
	name = "tf-testacclogtailbasic-%d"
	description = "create by terraform"
}
resource "alicloud_log_store" "example"{
  	project = "${alicloud_log_project.example.name}"
  	name = "tf-testacclogtailbasic-logstore"
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
  	name = "tf-testacclogtailbasic-config"
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
`, rand)
}
func testAlicloudLogTailUpdateOneParamater(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "update"{
	name = "tf-testacclogtailupdate-%d"
	description = "create by terraform"
}
resource "alicloud_log_store" "update"{
  	project = "${alicloud_log_project.update.name}"
  	name = "tf-testacclogtailupdate-logstore"
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
  	name = "tf-testacclogtailupdate-config"
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
`, rand)
}
func testAlicloudLogTailUpdateAllParamater(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "update_all"{
	name = "tf-testacclogtailupdateall-%d"
	description = "create by terraform"
}
resource "alicloud_log_store" "update_all"{
  	project = "${alicloud_log_project.update_all.name}"
  	name = "tf-testacclogtailupdateall-logstore"
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
  	name = "tf-testacclogtailupdateall-config"
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
`, rand)
}
func testAlicloudLogTailUpdateInputDetail(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "delimiter"{
	name = "tf-testacclogtailupdate-%d"
	description = "create by terraform"
}
resource "alicloud_log_store" "delimiter"{
  	project = "${alicloud_log_project.delimiter.name}"
  	name = "tf-testacclogtaildelimiter-logstore"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_logtail_config" "delimiter"{
	project = "${alicloud_log_project.delimiter.name}"
  	logstore = "${alicloud_log_store.delimiter.name}"
  	input_type = "file"
  	log_sample = "test-logtail-delimiter"
  	name = "tf-testacclogtaildelimiter-config"
	output_type = "LogService"
  	input_detail = <<DEFINITION
{
	"logPath": "/logPath",
	"filePattern": "*",
	"logType": "delimiter_log",
	"topicFormat": "default",
	"discardUnmatch": true,
	"enableRawLog": true,
	"fileEncoding": "utf8",
	"maxDepth": 999,
	"separator": ",",
	"quote": "\"",
	"timekey":"test",
	"key": [
		"test",
		"test2"
	],
	"autoExtend": true
}
	DEFINITION
}
`, rand)
}
