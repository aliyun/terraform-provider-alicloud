package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudLogtailAttachmentBasic(t *testing.T) {
	var v string
	resourceId := "alicloud_logtail_attachment.default"
	ra := resourceAttrInit(resourceId, logtailAttachmentMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogtailattachment-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogtailAttachmentDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":             "${alicloud_log_project.default.name}",
					"logtail_config_name": "${alicloud_logtail_config.default.name}",
					"machine_group_name":  "${alicloud_log_machine_group.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":             name,
						"logtail_config_name": name,
						"machine_group_name":  name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"machine_group_name"},
			},
		},
	})
}

func TestAccAlicloudLogtailAttachmentMultipleGroup(t *testing.T) {
	var v string
	resourceId := "alicloud_logtail_attachment.default.1"
	ra := resourceAttrInit(resourceId, logtailAttachmentMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogtailattachment-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogtailAttachmentDependenceMultipleGroup)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":             "${alicloud_log_project.default.name}",
					"logtail_config_name": "${alicloud_logtail_config.default.name}",
					"machine_group_name":  "${element(alicloud_log_machine_group.default.*.name,count.index)}",
					"count":               "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudLogtailAttachmentMultipleConfig(t *testing.T) {
	var v string
	resourceId := "alicloud_logtail_attachment.default.1"
	ra := resourceAttrInit(resourceId, logtailAttachmentMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogtailattachment-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogtailAttachmentDependenceMultipleConfig)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":             "${alicloud_log_project.default.name}",
					"logtail_config_name": "${element(alicloud_logtail_config.default.*.name, count.index)}",
					"machine_group_name":  "${alicloud_log_machine_group.default.name}",
					"count":               "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceLogtailAttachmentDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_log_project" "default"{
	name = "${var.name}"
	description = "create by terraform"
}
resource "alicloud_log_store" "default"{
  	project = "${alicloud_log_project.default.name}"
  	name = "${var.name}"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_log_machine_group" "default" {
	    project = "${alicloud_log_project.default.name}"
	    name = "${var.name}"
	    topic = "terraform"
	    identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}
resource "alicloud_logtail_config" "default"{
	project = "${alicloud_log_project.default.name}"
  	logstore = "${alicloud_log_store.default.name}"
  	input_type = "file"
  	log_sample = "test-update"
  	name = "${var.name}"
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
`, name)
}

func resourceLogtailAttachmentDependenceMultipleGroup(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default"{
	name = "${var.name}"
	description = "create by terraform"
}
resource "alicloud_log_store" "default"{
  	project = "${alicloud_log_project.default.name}"
  	name = "${var.name}"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_log_machine_group" "default" {
	count = 2
	project = "${alicloud_log_project.default.name}"
	name = "${var.name}-${count.index}"
	topic = "terraform"
	identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}
resource "alicloud_logtail_config" "default"{
	project = "${alicloud_log_project.default.name}"
  	logstore = "${alicloud_log_store.default.name}"
  	input_type = "file"
  	log_sample = "test-update"
  	name = "${var.name}"
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
`, name)
}

func resourceLogtailAttachmentDependenceMultipleConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default"{
	name = "${var.name}"
	description = "create by terraform"
}
resource "alicloud_log_store" "default"{
  	project = "${alicloud_log_project.default.name}"
  	name = "${var.name}"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_log_machine_group" "default" {
	    project = "${alicloud_log_project.default.name}"
	    name = "${var.name}"
	    topic = "terraform"
	    identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}
resource "alicloud_logtail_config" "default"{
	count = 2
	project = "${alicloud_log_project.default.name}"
  	logstore = "${alicloud_log_store.default.name}"
  	input_type = "file"
  	log_sample = "test-json-sample"
  	name = "${var.name}-${count.index}"
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
`, name)
}

var logtailAttachmentMap = map[string]string{
	"logtail_config_name": CHECKSET,
	"project":             CHECKSET,
	"machine_group_name":  CHECKSET,
}
