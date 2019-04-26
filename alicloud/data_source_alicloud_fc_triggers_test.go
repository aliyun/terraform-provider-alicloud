package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudFCTriggersDataSource_basic(t *testing.T) {
	randInt := acctest.RandInt()
	serviceName := fmt.Sprintf("tf-testacc-fc-trigger-ds-basic-%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFcTriggersDataSourceBasic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_fc_triggers.triggers"),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.name", serviceName),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.source_arn"),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.type", "log"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.invocation_role"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.config"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.last_modification_time"),
				),
			},
		},
	})
}

func TestAccAlicloudFCTriggersDataSource_empty(t *testing.T) {
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFcTriggersDataSourceEmpty(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_fc_triggers.triggers"),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.source_arn"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.type"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.invocation_role"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.config"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.creation_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.last_modification_time"),
				),
			},
		},
	})
}

func TestAccAlicloudFCTriggersDataSource_mnsTopic(t *testing.T) {
	randInt := acctest.RandInt()
	serviceName := fmt.Sprintf("tf-testacc-fc-trigger-ds-mns-topic-%d", randInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFcTriggersDataSourceMnsTopic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_fc_triggers.triggers"),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.name", serviceName),
					resource.TestMatchResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.source_arn", regexp.MustCompile(fmt.Sprintf("acs:mns:[a-z0-9-]+:[a-z0-9]+:/topics/tf-testacc-fc-trigger-ds-mns-topic-%v", randInt))),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.type", "mns_topic"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.invocation_role"),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.0.config", testTriggerMnsTopicTemplateDs),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_fc_triggers.triggers", "triggers.0.last_modification_time"),
				),
			},
			{
				Config: testAccCheckAlicloudFcTriggersDataSourceMnsTopicFake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_fc_triggers.triggers"),
					resource.TestCheckResourceAttr("data.alicloud_fc_triggers.triggers", "triggers.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudFcTriggersDataSourceBasic(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-fc-trigger-ds-basic-%d"
}

data "alicloud_regions" "current_region" {
  current = true
}
data "alicloud_account" "current" {
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
  handler = "hello.handler"
}

resource "alicloud_fc_trigger" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:log:${data.alicloud_regions.current_region.regions.0.id}:${data.alicloud_account.current.id}:project/${alicloud_log_project.foo.name}"
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

data "alicloud_fc_triggers" "triggers" {
	service_name = "${alicloud_fc_service.foo.name}"
	function_name = "${alicloud_fc_function.foo.name}"
    name_regex = "${alicloud_fc_trigger.foo.name}"
}
`, randInt, testTriggerLogTemplateDs, testFCLogRoleTemplateDs, testFCLogPolicyTemplateDs)
}

var testTriggerLogTemplateDs = `
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

func testAccCheckAlicloudFcTriggersDataSourceEmpty(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-fc-trigger-ds-basic-%d"
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
  handler = "hello.handler"
}

data "alicloud_fc_triggers" "triggers" {
	service_name = "${alicloud_fc_service.foo.name}"
	function_name = "${alicloud_fc_function.foo.name}"
    name_regex = "^tf-testacc-fake-name"
}
`, randInt)
}

func testAccCheckAlicloudFcTriggersDataSourceMnsTopic(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-fc-trigger-ds-mns-topic-%d"
}
data "alicloud_regions" "current_region" {
  current = true
}
data "alicloud_account" "current" {
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
resource "alicloud_mns_topic" "foo" {
  name = "${var.name}"
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
  handler = "hello.handler"
}
resource "alicloud_fc_trigger" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:mns:${data.alicloud_regions.current_region.regions.0.id}:${data.alicloud_account.current.id}:/topics/${alicloud_mns_topic.foo.name}"
  type = "mns_topic"
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
data "alicloud_fc_triggers" "triggers" {
	service_name = "${alicloud_fc_service.foo.name}"
	function_name = "${alicloud_fc_function.foo.name}"
    name_regex = "${alicloud_fc_trigger.foo.name}"
}
`, randInt, testTriggerMnsTopicTemplateDs, testFCMnsTopicRoleTemplateDs, testFCMnsTopicPolicyTemplateDs)
}

func testAccCheckAlicloudFcTriggersDataSourceMnsTopicFake(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-fc-trigger-ds-mns-topic-%d"
}
data "alicloud_regions" "current_region" {
  current = true
}
data "alicloud_account" "current" {
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
resource "alicloud_mns_topic" "foo" {
  name = "${var.name}"
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
  handler = "hello.handler"
}
resource "alicloud_fc_trigger" "foo" {
  service = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:mns:${data.alicloud_regions.current_region.regions.0.id}:${data.alicloud_account.current.id}:/topics/${alicloud_mns_topic.foo.name}"
  type = "mns_topic"
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
data "alicloud_fc_triggers" "triggers" {
	service_name = "${alicloud_fc_service.foo.name}"
	function_name = "${alicloud_fc_function.foo.name}"
    name_regex = "${alicloud_fc_trigger.foo.name}_fake"
}
`, randInt, testTriggerMnsTopicTemplateDs, testFCMnsTopicRoleTemplateDs, testFCMnsTopicPolicyTemplateDs)
}

var testTriggerMnsTopicTemplateDs = `{"filterTag":"testTag","notifyContentFormat":"STREAM","notifyStrategy":"BACKOFF_RETRY"}`

var testFCMnsTopicRoleTemplateDs = `
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "mns.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
`

var testFCMnsTopicPolicyTemplateDs = `
    {
      "Version": "1",
      "Statement": [
        {
          "Action": [
            "log:PostLogStoreLogs"
          ],
          "Resource": "*",
          "Effect": "Allow"
        }
      ]
    }
`
