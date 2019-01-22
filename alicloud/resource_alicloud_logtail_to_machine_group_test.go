package alicloud

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"
)

func TestAccAlicloudLogtailToMachineGroup_basic(t *testing.T) {
	var project sls.LogProject
	var group sls.MachineGroup
	var store sls.LogStore
	var config sls.LogConfig
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogtailToMachineGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogtailToMachineGroup_basic(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.test", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.test", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.test", &config),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.test", &group),
					testAccCheckAlicloudLogtailToMachineGroupExists("alicloud_logtail_attachment.test", config.Name),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.test", "logtail_config_name", "tf-testaccapplyconfigbasic-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.test", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group"),
				),
			},
			{
				Config: testAccCheckAlicloudLogtailToMachineGroup_MultipleGroup(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.multiple_group", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.multiple_group", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.multiple_group", &config),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.multiple_group", &group),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.multiple_group2", &group),
					testAccCheckAlicloudLogtailToMachineGroupExists("alicloud_logtail_attachment.multiple_group", config.Name),
					testAccCheckAlicloudLogtailToMachineGroupExists("alicloud_logtail_attachment.multiple_group2", config.Name),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group", "logtail_config_name", "tf-testaccapplyconfigbasic-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group2", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group2"),
				),
			},
			{
				Config: testAccCheckAlicloudLogtailToMachineGroup_MultipleConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.multiple_config", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.multiple_config", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.multiple_config", &config),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.multiple_config2", &config),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.multiple_config", &group),
					testAccCheckAlicloudLogtailToMachineGroupExists("alicloud_logtail_attachment.multiple_config", config.Name),
					testAccCheckAlicloudLogtailToMachineGroupExists("alicloud_logtail_attachment.multiple_config2", config.Name),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config", "logtail_config_name", "tf-testaccapplyconfigbasic-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config2", "logtail_config_name", "tf-testaccapplyconfigbasic-config2"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group"),
				),
			},
		},
	})
}

func testAccCheckAlicloudLogtailToMachineGroupExists(name, group_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", name))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No Log machine group ID is set"))
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		logService := LogService{client}
		groupNames, err := logService.DescribeLogtailToMachineGroup(split[0], split[1])
		if err != nil {
			return WrapError(err)
		}
		for _, name := range groupNames {
			if name == split[2] {
				group_name = name
			}
		}
		return nil
	}
}

func testAccCheckAlicloudLogtailToMachineGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	logService := LogService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "" {
			continue
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if _, err := logService.DescribeLogMachineGroup(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("Logtail to machine group %s still exists.", rs.Primary.ID))
	}
	return nil
}

func testAccCheckAlicloudLogtailToMachineGroup_basic(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "test"{
	name = "tf-testaccapplyconfigbasic-%d"
	description = "create by terraform"
}
resource "alicloud_log_store" "test"{
  	project = "${alicloud_log_project.test.name}"
  	name = "tf-testaccapplyconfigbasic-logstore"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_log_machine_group" "test" {
	    project = "${alicloud_log_project.test.name}"
	    name = "tf-testaccapplyconfigbasic-machine-group"
	    topic = "terraform"
	    identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}
resource "alicloud_logtail_config" "test"{
	project = "${alicloud_log_project.test.name}"
  	logstore = "${alicloud_log_store.test.name}"
  	input_type = "file"
  	log_sample = "test-update"
  	name = "tf-testaccapplyconfigbasic-config"
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
resource "alicloud_logtail_attachment" "test" {
	project = "${alicloud_log_project.test.name}"
	logtail_config_name = "${alicloud_logtail_config.test.name}"
	machine_group_name = "${alicloud_log_machine_group.test.name}"
}`, rand)
}

func testAccCheckAlicloudLogtailToMachineGroup_MultipleGroup(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "multiple_group"{
	name = "tf-testaccapplyconfigbasic-%d"
	description = "create by terraform"
}
resource "alicloud_log_store" "multiple_group"{
  	project = "${alicloud_log_project.multiple_group.name}"
  	name = "tf-testaccapplyconfigbasic-logstore"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_log_machine_group" "multiple_group" {
	    project = "${alicloud_log_project.multiple_group.name}"
	    name = "tf-testaccapplyconfigbasic-machine-group"
	    topic = "terraform"
	    identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}
resource "alicloud_log_machine_group" "multiple_group2" {
	    project = "${alicloud_log_project.multiple_group.name}"
	    name = "tf-testaccapplyconfigbasic-machine-group2"
	    topic = "terraform"
	    identify_list = ["10.1.1.1", "10.1.1.3", "10.1.1.2"]
}
resource "alicloud_logtail_config" "multiple_group"{
	project = "${alicloud_log_project.multiple_group.name}"
  	logstore = "${alicloud_log_store.multiple_group.name}"
  	input_type = "file"
  	log_sample = "test-update"
  	name = "tf-testaccapplyconfigbasic-config"
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
resource "alicloud_logtail_attachment" "multiple_group" {
	project = "${alicloud_log_project.multiple_group.name}"
	logtail_config_name = "${alicloud_logtail_config.multiple_group.name}"
	machine_group_name = "${alicloud_log_machine_group.multiple_group.name}"
}
resource "alicloud_logtail_attachment" "multiple_group2" {
	project = "${alicloud_log_project.multiple_group.name}"
	logtail_config_name = "${alicloud_logtail_config.multiple_group.name}"
	machine_group_name = "${alicloud_log_machine_group.multiple_group2.name}"
}
`, rand)
}

func testAccCheckAlicloudLogtailToMachineGroup_MultipleConfig(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "multiple_config"{
	name = "tf-testaccapplyconfigbasic-%d"
	description = "create by terraform"
}
resource "alicloud_log_store" "multiple_config"{
  	project = "${alicloud_log_project.multiple_config.name}"
  	name = "tf-testaccapplyconfigbasic-logstore"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_log_machine_group" "multiple_config" {
	    project = "${alicloud_log_project.multiple_config.name}"
	    name = "tf-testaccapplyconfigbasic-machine-group"
	    topic = "terraform"
	    identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}

resource "alicloud_logtail_config" "multiple_config"{
	project = "${alicloud_log_project.multiple_config.name}"
  	logstore = "${alicloud_log_store.multiple_config.name}"
  	input_type = "file"
  	log_sample = "test-json-sample"
  	name = "tf-testaccapplyconfigbasic-config"
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
resource "alicloud_logtail_config" "multiple_config2"{
	project = "${alicloud_log_project.multiple_config.name}"
  	logstore = "${alicloud_log_store.multiple_config.name}"
  	input_type = "plugin"
  	log_sample = "test-logtail-plugin-sample"
  	name = "tf-testaccapplyconfigbasic-config2"
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
resource "alicloud_logtail_attachment" "multiple_config" {
	project = "${alicloud_log_project.multiple_config.name}"
	logtail_config_name = "${alicloud_logtail_config.multiple_config.name}"
	machine_group_name = "${alicloud_log_machine_group.multiple_config.name}"
}
resource "alicloud_logtail_attachment" "multiple_config2" {
	project = "${alicloud_log_project.multiple_config.name}"
	logtail_config_name = "${alicloud_logtail_config.multiple_config2.name}"
	machine_group_name = "${alicloud_log_machine_group.multiple_config.name}"
}
`, rand)
}
