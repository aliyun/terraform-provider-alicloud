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
				Config: testAccCheckAlicloudLogtailToMachineGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.test", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.test", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.test", &config),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.test", &group),
					testAccCheckAlicloudLogtailToMachineGroupExists("alicloud_logtail_to_machine_group.test", config.Name),
					resource.TestCheckResourceAttr("alicloud_logtail_to_machine_group.test", "project", "test-tf2"),
					resource.TestCheckResourceAttr("alicloud_logtail_to_machine_group.test", "logtail_config_name", "evan-terraform-update"),
					resource.TestCheckResourceAttr("alicloud_logtail_to_machine_group.test", "machine_group_name", "evan-machine-group"),
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

const testAccCheckAlicloudLogtailToMachineGroup_basic = `
resource "alicloud_log_project" "test"{
	name = "test-tf2"
	description = "create by terraform"
}
resource "alicloud_log_store" "test"{
  	project = "${alicloud_log_project.test.name}"
  	name = "tf-test-logstore"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_log_machine_group" "test" {
	    project = "${alicloud_log_project.test.name}"
	    name = "evan-machine-group"
	    topic = "terraform"
	    identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}
resource "alicloud_logtail_config" "test"{
	project = "${alicloud_log_project.test.name}"
  	logstore = "${alicloud_log_store.test.name}"
  	input_type = "file"
  	log_sample = "test-update"
  	name = "evan-terraform-update"
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
resource "alicloud_logtail_to_machine_group" "test" {
	project = "${alicloud_log_project.test.name}"
	logtail_config_name = "${alicloud_logtail_config.test.name}"
	machine_group_name = "${alicloud_log_machine_group.test.name}"
}`
