package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudFCTrigger_log(t *testing.T) {
	var service fc.GetServiceOutput
	var project sls.LogProject
	var store sls.LogStore
	var function fc.GetFunctionOutput
	var trigger fc.GetTriggerOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudFCTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudFCTriggerLog(testTriggerLogTemplate, testFCLogRoleTemplate, testFCLogPolicyTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.foo", &project),
					testAccCheckAlicloudLogStoreExists("alicloud_log_store.foo", &store),
					testAccCheckAlicloudFCServiceExists("alicloud_fc_service.foo", &service),
					testAccCheckAlicloudFCFunctionExists("alicloud_fc_function.foo", &function),
					testAccCheckAlicloudFCTriggerExists("alicloud_fc_trigger.foo", &trigger),
					resource.TestCheckResourceAttr("alicloud_fc_trigger.foo", "name", "test-alicloud-fc-trigger"),
					resource.TestCheckResourceAttrSet("alicloud_fc_trigger.foo", "config"),
				),
			},
			{
				Config: testAlicloudFCTriggerLogUpdate(testTriggerLogTemplateUpdate, testFCLogRoleTemplate, testFCLogPolicyTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudFCTriggerExists("alicloud_fc_trigger.foo", &trigger),
					resource.TestCheckResourceAttr("alicloud_fc_trigger.foo", "name", "test-alicloud-fc-trigger"),
					resource.TestCheckResourceAttrSet("alicloud_fc_trigger.foo", "config"),
				),
			},
		},
	})
}

func testAccCheckAlicloudFCTriggerExists(name string, trigger *fc.GetTriggerOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Log store ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		ser, err := client.DescribeFcTrigger(split[0], split[1], split[2])
		if err != nil {
			return err
		}
		trigger = ser

		return nil
	}
}

func testAccCheckAlicloudFCTriggerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_fc_trigger" {
			continue
		}

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		ser, err := client.DescribeFcTrigger(split[0], split[1], split[2])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Check fc service got an error: %#v.", err)
		}

		if ser == nil {
			return nil
		}

		return fmt.Errorf("FC service %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAlicloudFCTriggerLog(trigger, role, policy string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  account_id = "${var.account}"
  region = "${var.region}"
}
variable "region" {
  default = "cn-hangzhou"
}
variable "account" {
  default = "1204663572767468"
}
variable "name" {
  default = "test-alicloud-fc-trigger"
}

resource "alicloud_log_project" "foo" {
  name = "${var.name}"
  description = "tf unit test"
}
resource "alicloud_log_store" "bar" {
  project = "${alicloud_log_project.foo.name}"
  name = "${var.name}-source"
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_log_store" "foo" {
  project = "${alicloud_log_project.foo.name}"
  name = "${var.name}"
  retention_period = "3000"
  shard_count = 1
}

resource "alicloud_fc_service" "foo" {
  name = "${var.name}"
  internet_access = false
}

resource "alicloud_oss_bucket" "foo" {
  bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "foo" {
  bucket = "${alicloud_oss_bucket.foo.id}"
  key = "fc/hello.zip"
  content = <<EOF
  	# -*- coding: utf-8 -*-
	def handler(event, context):
	    print "hello world"
	    return 'hello world'
  EOF
}

resource "alicloud_fc_function" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  name = "${var.name}"
  oss_bucket = "${alicloud_oss_bucket.foo.id}"
  oss_key = "${alicloud_oss_bucket_object.foo.key}"
  memory_size = 512
  runtime = "python2.7"
}

resource "alicloud_fc_trigger" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:log:${var.region}:${var.account}:project/${alicloud_log_project.foo.name}"
  type = "log"
  config = <<EOF
  %s
  EOF
  depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}

resource "alicloud_ram_role" "foo" {
  name = "${var.name}-trigger"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}

resource "alicloud_ram_policy" "foo" {
  name = "${var.name}-trigger"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name = "${alicloud_ram_role.foo.name}"
  policy_name = "${alicloud_ram_policy.foo.name}"
  policy_type = "Custom"
}
`, trigger, role, policy)
}

func testAlicloudFCTriggerLogUpdate(trigger, role, policy string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  account_id = "${var.account}"
  region = "${var.region}"
}
variable "region" {
  default = "cn-hangzhou"
}
variable "account" {
  default = "1204663572767468"
}
variable "name" {
  default = "test-alicloud-fc-trigger"
}

resource "alicloud_log_project" "foo" {
  name = "${var.name}"
  description = "tf unit test"
}
resource "alicloud_log_store" "bar" {
  project = "${alicloud_log_project.foo.name}"
  name = "${var.name}-source"
  retention_period = "3000"
  shard_count = 1
}
resource "alicloud_log_store" "foo" {
  project = "${alicloud_log_project.foo.name}"
  name = "${var.name}"
  retention_period = "3000"
  shard_count = 1
}

resource "alicloud_fc_service" "foo" {
  name = "${var.name}"
  internet_access = false
}

resource "alicloud_oss_bucket" "foo" {
  bucket = "${var.name}"
}

resource "alicloud_oss_bucket_object" "foo" {
  bucket = "${alicloud_oss_bucket.foo.id}"
  key = "fc/hello.zip"
  content = <<EOF
  	# -*- coding: utf-8 -*-
	def handler(event, context):
	    print "hello world"
	    return 'hello world'
  EOF
}

resource "alicloud_fc_function" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  name = "${var.name}"
  oss_bucket = "${alicloud_oss_bucket.foo.id}"
  oss_key = "${alicloud_oss_bucket_object.foo.key}"
  memory_size = 512
  runtime = "python2.7"
}

resource "alicloud_fc_trigger" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:log:${var.region}:${var.account}:project/${alicloud_log_project.foo.name}"
  type = "log"
  config = <<EOF
  %s
  EOF
  depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}

resource "alicloud_ram_role" "foo" {
  name = "${var.name}-trigger"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}

resource "alicloud_ram_policy" "foo" {
  name = "${var.name}-trigger"
  document = <<EOF
  %s
  EOF
  description = "this is a test"
  force = true
}
resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name = "${alicloud_ram_role.foo.name}"
  policy_name = "${alicloud_ram_policy.foo.name}"
  policy_type = "Custom"
}
`, trigger, role, policy)
}

var testTriggerLogTemplate = `
    {
        "sourceConfig": {
            "project": "${alicloud_log_project.foo.name}",
            "logstore": "${alicloud_log_store.bar.name}"
        },
        "jobConfig": {
            "maxRetryTime": 3,
            "triggerInterval": 60
        },
        "functionParameter": {
            "a": "b",
            "c": "d"
        },
        "logConfig": {
            "project": "${alicloud_log_project.foo.name}",
            "logstore": "${alicloud_log_store.foo.name}"
        },
        "enable": true
    }
`

var testTriggerLogTemplateUpdate = `
    {
        "sourceConfig": {
            "project": "${alicloud_log_project.foo.name}",
            "logstore": "${alicloud_log_store.bar.name}"
        },
        "jobConfig": {
            "maxRetryTime": 4,
            "triggerInterval": 100
        },
        "functionParameter": {
            "a": "bb",
            "c": "dd"
        },
        "logConfig": {
            "project": "${alicloud_log_project.foo.name}",
            "logstore": "${alicloud_log_store.foo.name}"
        },
        "enable": true
    }
`
var testFCLogPolicyTemplate = `
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

var testFCLogRoleTemplate = `
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
