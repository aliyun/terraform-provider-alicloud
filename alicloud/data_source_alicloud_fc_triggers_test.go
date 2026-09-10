package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFCTriggersDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_triggers.default"
	name := fmt.Sprintf("tf-testacc%sfctriggerbasic-%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceFCtriggerLogConfigDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name":  "${alicloud_fc_trigger.default.service}",
			"function_name": "${alicloud_fc_trigger.default.function}",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_name":  "${alicloud_fc_trigger.default.service}",
			"function_name": "${alicloud_fc_trigger.default.function}",
			"ids":           []string{"${alicloud_fc_trigger.default.trigger_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_name":  "${alicloud_fc_trigger.default.service}",
			"function_name": "${alicloud_fc_trigger.default.function}",
			"ids":           []string{"${alicloud_fc_trigger.default.trigger_id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":    "${alicloud_fc_trigger.default.name}",
			"service_name":  "${alicloud_fc_trigger.default.service}",
			"function_name": "${alicloud_fc_trigger.default.function}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":    "${alicloud_fc_trigger.default.name}_fake",
			"service_name":  "${alicloud_fc_trigger.default.service}",
			"function_name": "${alicloud_fc_trigger.default.function}",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":    "${alicloud_fc_trigger.default.name}",
			"service_name":  "${alicloud_fc_trigger.default.service}",
			"function_name": "${alicloud_fc_trigger.default.function}",
			"ids":           []string{"${alicloud_fc_trigger.default.trigger_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":    "${alicloud_fc_trigger.default.name}_fake",
			"service_name":  "${alicloud_fc_trigger.default.service}",
			"function_name": "${alicloud_fc_trigger.default.function}",
			"ids":           []string{"${alicloud_fc_trigger.default.trigger_id}"},
		}),
	}

	var existFCtriggerMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"triggers.#":                        "1",
			"ids.#":                             "1",
			"names.#":                           "1",
			"triggers.0.id":                     CHECKSET,
			"triggers.0.name":                   name,
			"triggers.0.type":                   "log",
			"triggers.0.source_arn":             CHECKSET,
			"triggers.0.invocation_role":        CHECKSET,
			"triggers.0.config":                 CHECKSET,
			"triggers.0.creation_time":          CHECKSET,
			"triggers.0.last_modification_time": CHECKSET,
		}
	}

	var fakeFCtriggerMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"triggers.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var fcTriggerRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existFCtriggerMapFunc,
		fakeMapFunc:  fakeFCtriggerMapFunc,
	}

	fcTriggerRecordsCheckInfo.dataSourceTestCheck(t, rand, basicConf, idsConf, nameRegexConf, allConf)

}

func dataSourceFCtriggerLogConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}
data "alicloud_account" "default" {
}

resource "alicloud_log_project" "default" {
  name = "${var.name}"
  description = "tf unit test"
}
resource "alicloud_log_store" "default" {
  project = "${alicloud_log_project.default.name}"
  name = "${var.name}-source"
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_log_store" "default1" {
  project = "${alicloud_log_project.default.name}"
  name = "${var.name}"
  retention_period = "3000"
  shard_count = 1
}

resource "alicloud_fc_service" "default" {
  name = "${var.name}"
  internet_access = false
}

resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = "${alicloud_oss_bucket.default.id}"
  key = "fc/hello.zip"
  content = <<EOF
  	# -*- coding: utf-8 -*-
	def handler(event, context):
	    print "hello world"
	    return 'hello world'
  EOF
}

resource "alicloud_fc_function" "default" {
  service = "${alicloud_fc_service.default.name}"
  name = "${var.name}"
  oss_bucket = "${alicloud_oss_bucket.default.id}"
  oss_key = "${alicloud_oss_bucket_object.default.key}"
  memory_size = 512
  runtime = "python3.9"
  handler = "hello.handler"
}

resource "alicloud_fc_trigger" "default" {
  service = "${alicloud_fc_service.default.name}"
  function = "${alicloud_fc_function.default.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.default.arn}"
  source_arn = "acs:log:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:project/${alicloud_log_project.default.name}"
  type = "log"
  config = <<EOF
  %s
  EOF
  depends_on = ["alicloud_ram_role_policy_attachment.default"]
}

resource "alicloud_ram_role" "default" {
  name = "${var.name}-trigger"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}

resource "alicloud_ram_policy" "default" {
  name = "${var.name}-trigger"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "${alicloud_ram_policy.default.name}"
  policy_type = "Custom"
}`, name, testTriggerLogTemplateDs, testFCLogRoleTemplateDs, testFCLogPolicyTemplateDs)
}

var testTriggerLogTemplateDs = `
	{
		"sourceConfig":{
			"logstore":"${alicloud_log_store.default.name}",
			"startTime":null
		},
		"jobConfig":{
			"maxRetryTime":3,
			"triggerInterval":60
		},
		"functionParameter":{
			"a":"b",
			"c":"d"
		},
		"logConfig":{
			"project":"${alicloud_log_project.default.name}",
			"logstore":"${alicloud_log_store.default1.name}"
		},
		"enable":true,
		"targetConfig":null
	}
`

var testFCLogRoleTemplateDs = `
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "log.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
`

var testFCLogPolicyTemplateDs = `
    {
      "Version": "1",
      "Statement": [
        {
          "Action": [
            "log:Get*",
            "log:List*",
            "log:PostLogStoreLogs",
            "log:CreateConsumerGroup",
            "log:UpdateConsumerGroup",
            "log:DeleteConsumerGroup",
            "log:ListConsumerGroup",
            "log:ConsumerGroupUpdateCheckPoint",
            "log:ConsumerGroupHeartBeat",
            "log:GetConsumerGroupCheckPoint"
          ],
          "Resource": "*",
          "Effect": "Allow"
        }
      ]
    }
`
