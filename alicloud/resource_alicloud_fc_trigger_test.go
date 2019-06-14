package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudFCTrigger_log(t *testing.T) {
	var v *fc.GetTriggerOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc%s-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"name":   name,
		"config": CHECKSET,
	}
	resourceId := "alicloud_fc_trigger.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlicloudFCTriggerLog(testTriggerLogTemplate, testFCLogRoleTemplate, testFCLogPolicyTemplate, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "filename", "oss_bucket", "oss_key"},
			},
			{
				Config: testAlicloudFCTriggerLogUpdate(testTriggerLogTemplateUpdate, testFCLogRoleTemplate, testFCLogPolicyTemplate, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudFCTrigger_mnsTopic(t *testing.T) {
	var v *fc.GetTriggerOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc%s-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"name":       name,
		"source_arn": CHECKSET,
		"config_mns": testTriggerMnsTopicTemplate,
		"type":       "mns_topic",
	}
	resourceId := "alicloud_fc_trigger.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlicloudFCTriggerMnsTopic(testTriggerMnsTopicTemplate, testFCMnsTopicRoleTemplate, testFCMnsTopicPolicyTemplate, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudFCTrigger_cdn_events(t *testing.T) {
	var v *fc.GetTriggerOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc%s-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"service":       CHECKSET,
		"function":      CHECKSET,
		"source_arn":    CHECKSET,
		"name":          name,
		"config":        CHECKSET,
		"type":          "cdn_events",
		"last_modified": CHECKSET,
	}
	resourceId := "alicloud_fc_trigger.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlicloudFCTriggerCdnEvents(testTriggerCdnEventsTemplate, testFCcdnRoleTemplate, testFCcdnPolicyTemplate, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "filename", "oss_bucket", "oss_key"},
			},
			{
				Config: testAlicloudFCTriggerUpdateCdnEvents(testTriggerUpdateCdnEventsTemplate, testFCcdnRoleTemplate, testFCcdnPolicyTemplate, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAlicloudFCTriggerMnsTopic(trigger, role, policy, name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
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
resource "alicloud_mns_topic" "default" {
  name = "${var.name}"
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
  runtime = "python2.7"
  handler = "hello.handler"
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
}
resource "alicloud_fc_trigger" "default" {
  service = "${alicloud_fc_service.default.name}"
  function = "${alicloud_fc_function.default.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.default.arn}"
  source_arn = "acs:mns:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:/topics/${alicloud_mns_topic.default.name}"
  type = "mns_topic"
  config_mns = <<EOF
  %s
  EOF
  depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
`, name, role, policy, trigger)
}

var testTriggerMnsTopicTemplate = `{"filterTag":"testTag","notifyContentFormat":"STREAM","notifyStrategy":"BACKOFF_RETRY"}`

var testFCMnsTopicPolicyTemplate = `
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

var testFCMnsTopicRoleTemplate = `
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

func testAlicloudFCTriggerLog(trigger, role, policy, name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
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
  runtime = "python2.7"
  handler = "hello.handler"
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
`, name, role, policy, trigger)
}

func testAlicloudFCTriggerLogUpdate(trigger, role, policy, name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
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
  runtime = "python2.7"
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
}
`, name, trigger, role, policy)
}

var testTriggerLogTemplate = `
    {
        "sourceConfig": {
            "logstore": "${alicloud_log_store.default.name}"
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
            "project": "${alicloud_log_project.default.name}",
            "logstore": "${alicloud_log_store.default1.name}"
        },
        "enable": true
    }
`

var testTriggerLogTemplateUpdate = `
    {
        "sourceConfig": {
            "logstore": "${alicloud_log_store.default.name}"
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
            "project": "${alicloud_log_project.default.name}",
            "logstore": "${alicloud_log_store.default1.name}"
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

var testTriggerCdnEventsTemplate = `{"eventName":"LogFileCreated","eventVersion":"1.0.0","notes":"cdn events trigger","filter":{"domain":["${alicloud_cdn_domain.default.domain_name}"]}}`
var testTriggerUpdateCdnEventsTemplate = `{"eventName":"LogFileCreated","eventVersion":"1.0.0","notes":"update cdn events trigger","filter":{"domain":["${alicloud_cdn_domain.default.domain_name}"]}}`

var testFCcdnPolicyTemplate = `
{
  "Version": "1",
  "Statement": [
    {
      "Action": [
        "fc:InvokeFunction"
      ],
      "Resource": [
        "acs:fc:*:*:services/tf_cdnEvents/functions/*",
        "acs:fc:*:*:services/tf_cdnEvents.*/functions/*"
      ],
      "Effect": "Allow"
    }
  ]
}
`

var testFCcdnRoleTemplate = `
{
   "Version": "1",
   "Statement": [
      {
         "Action": "cdn:Describe*",
         "Resource": "*",
         "Effect": "Allow",
		 "Principal": {
           "Service": [
               "log.aliyuncs.com"
           ]
         }
      }
   ]
}
`

func testAlicloudFCTriggerCdnEvents(trigger, role, policy, name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

data "alicloud_account" "default" {
}

resource "alicloud_cdn_domain" "default" {
	  domain_name = "${var.name}.xiaozhu.com"
	  cdn_type = "web"
	  source_type = "oss"
	  sources = ["terraformtest.aliyuncs.com"]
	  optimize_enable = "off"
	  page_compress_enable = "off"
	  range_enable = "off"
	  video_seek_enable = "off"
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
  runtime = "python2.7"
  handler = "hello.handler"
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
}
resource "alicloud_fc_trigger" "default" {
  service = "${alicloud_fc_service.default.name}"
  function = "${alicloud_fc_function.default.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.default.arn}"
  source_arn = "acs:cdn:*:${data.alicloud_account.default.id}"
  type = "cdn_events"
  config = <<EOF
%sEOF
  depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
`, name, role, policy, trigger)
}

func testAlicloudFCTriggerUpdateCdnEvents(trigger, role, policy, name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

data "alicloud_account" "default" {
}

resource "alicloud_cdn_domain" "default" {
	  domain_name = "${var.name}.xiaozhu.com"
	  cdn_type = "web"
	  source_type = "oss"
	  sources = ["terraformtest.aliyuncs.com"]
	  optimize_enable = "off"
	  page_compress_enable = "off"
	  range_enable = "off"
	  video_seek_enable = "off"
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
  runtime = "python2.7"
  handler = "hello.handler"
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
}
resource "alicloud_fc_trigger" "default" {
  service = "${alicloud_fc_service.default.name}"
  function = "${alicloud_fc_function.default.name}"
  name = "${var.name}"
  role = "${alicloud_ram_role.default.arn}"
  source_arn = "acs:cdn:*:${data.alicloud_account.default.id}"
  type = "cdn_events"
  config = <<EOF
%sEOF
  depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
`, name, role, policy, trigger)
}
