package alicloud

import (
	"fmt"
	"strings"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudLogtailAttachment_basic(t *testing.T) {
	var project sls.LogProject
	var group sls.MachineGroup
	var store sls.LogStore
	var config sls.LogConfig
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogtailAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogtailAttachment_basic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.test", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.test", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.test", &config),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.test", &group),
					testAccCheckAlicloudogtailAttachmentExists("alicloud_logtail_attachment.test", config.Name),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.test", "project", fmt.Sprintf("tf-testaccapplyconfigbasic-%v", randInt)),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.test", "logtail_config_name", "tf-testaccapplyconfigbasic-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.test", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group"),
				),
			},
		},
	})
}

func TestAccAlicloudLogtailAttachment_MultipleGroup(t *testing.T) {
	var project sls.LogProject
	var group sls.MachineGroup
	var store sls.LogStore
	var config sls.LogConfig
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogtailAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogtailAttachment_MultipleGroup(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.multiple_group", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.multiple_group", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.multiple_group", &config),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.multiple_group.0", &group),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.multiple_group.1", &group),
					testAccCheckAlicloudogtailAttachmentExists("alicloud_logtail_attachment.multiple_group.0", config.Name),
					testAccCheckAlicloudogtailAttachmentExists("alicloud_logtail_attachment.multiple_group.1", config.Name),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group.0", "project", fmt.Sprintf("tf-testaccapplyconfigbasic-%v", randInt)),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group.1", "project", fmt.Sprintf("tf-testaccapplyconfigbasic-%v", randInt)),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group.0", "logtail_config_name", "tf-testaccapplyconfigbasic-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group.1", "logtail_config_name", "tf-testaccapplyconfigbasic-config"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group.0", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group-0"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_group.1", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group-1"),
				),
			},
		},
	})
}

func TestAccAlicloudLogtailAttachment_MultipleConfig(t *testing.T) {
	var project sls.LogProject
	var group sls.MachineGroup
	var store sls.LogStore
	var config sls.LogConfig
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogtailAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogtailAttachment_MultipleConfig(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.multiple_config", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.multiple_config", &store),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.multiple_config.0", &config),
					testAccCheckAlicloudLogTailConfigExists("alicloud_logtail_config.multiple_config.1", &config),
					testAccCheckAlicloudLogMachineGroupExists("alicloud_log_machine_group.multiple_config", &group),
					testAccCheckAlicloudogtailAttachmentExists("alicloud_logtail_attachment.multiple_config.0", config.Name),
					testAccCheckAlicloudogtailAttachmentExists("alicloud_logtail_attachment.multiple_config.1", config.Name),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config.0", "project", fmt.Sprintf("tf-testaccapplyconfigbasic-%v", randInt)),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config.1", "project", fmt.Sprintf("tf-testaccapplyconfigbasic-%v", randInt)),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config.0", "logtail_config_name", "tf-testaccapplyconfigbasic-config-0"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config.1", "logtail_config_name", "tf-testaccapplyconfigbasic-config-1"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config.0", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group"),
					resource.TestCheckResourceAttr("alicloud_logtail_attachment.multiple_config.1", "machine_group_name", "tf-testaccapplyconfigbasic-machine-group"),
				),
			},
		},
	})
}

func testAccCheckAlicloudogtailAttachmentExists(name, group_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", name))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No Log machine group ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		logService := LogService{client}
		groupNames, err := logService.DescribeLogtailAttachment(rs.Primary.ID)
		if err != nil {
			return WrapError(err)
		}
		for _, name := range groupNames {
			if name == strings.Split(rs.Primary.ID, COLON_SEPARATED)[2] {
				group_name = name
			}
		}
		return nil
	}
}

func testAccCheckAlicloudLogtailAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	logService := LogService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_logtail_attachment" {
			continue
		}

		if _, err := logService.DescribeLogtailAttachment(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("Logtail Attachment %s still exists.", rs.Primary.ID))
	}
	return nil
}

func testAccCheckAlicloudLogtailAttachment_basic(rand int) string {
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

func testAccCheckAlicloudLogtailAttachment_MultipleGroup(rand int) string {
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
	count = 2
	project = "${alicloud_log_project.multiple_group.name}"
	name = "tf-testaccapplyconfigbasic-machine-group-${count.index}"
	topic = "terraform"
	identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
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
	count = 2
	project = "${alicloud_log_project.multiple_group.name}"
	logtail_config_name = "${alicloud_logtail_config.multiple_group.name}"
	machine_group_name = "${element(alicloud_log_machine_group.multiple_group.*.name,count.index)}"
}
`, rand)
}

func testAccCheckAlicloudLogtailAttachment_MultipleConfig(rand int) string {
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
	count = 2
	project = "${alicloud_log_project.multiple_config.name}"
  	logstore = "${alicloud_log_store.multiple_config.name}"
  	input_type = "file"
  	log_sample = "test-json-sample"
  	name = "tf-testaccapplyconfigbasic-config-${count.index}"
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
resource "alicloud_logtail_attachment" "multiple_config" {
	count = 2
	project = "${alicloud_log_project.multiple_config.name}"
	logtail_config_name = "${element(alicloud_logtail_config.multiple_config.*.name, count.index)}"
	machine_group_name = "${alicloud_log_machine_group.multiple_config.name}"
}
`, rand)
}
