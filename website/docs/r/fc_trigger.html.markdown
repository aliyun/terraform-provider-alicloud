---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_trigger"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides a Alicloud Function Compute Trigger resource.
---

# alicloud\_fc\_trigger

Provides an Alicloud Function Compute Trigger resource. Based on trigger, execute your code in response to events in Alibaba Cloud.
 For information about Service and how to use it, see [What is Function Compute](https://www.alibabacloud.com/help/doc-detail/52895.htm).

-> **NOTE:** The resource requires a provider field 'account_id'. [See account_id](https://www.terraform.io/docs/providers/alicloud/index.html#account_id).

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

variable "account" {
  default = "12345"
}

provider "alicloud" {
  account_id = var.account
  region     = var.region
}

resource "alicloud_fc_trigger" "foo" {
  service    = "my-fc-service"
  function   = "hello-world"
  name       = "hello-trigger"
  role       = alicloud_ram_role.foo.arn
  source_arn = "acs:log:${var.region}:${var.account}:project/${alicloud_log_project.foo.name}"
  type       = "log"
  config     = <<EOF
    {
        "sourceConfig": {
            "project": "project-for-fc",
            "logstore": "project-for-fc"
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
            "project": "project-for-fc",
            "logstore": "project-for-fc"
        },
        "enable": true
    }
  
EOF


  depends_on = [alicloud_ram_role_policy_attachment.foo]
}

resource "alicloud_ram_role" "foo" {
  name     = "${var.name}-trigger"
  document = <<EOF
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
  
EOF


  description = "this is a test"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name   = alicloud_ram_role.foo.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}
```

MNS topic trigger:

```terraform
variable "name" {
  default = "fctriggermnstopic"
}
data "alicloud_regions" "current_region" {
  current = true
}
data "alicloud_account" "current" {
}
resource "alicloud_log_project" "foo" {
  name        = "${var.name}"
  description = "tf unit test"
}
resource "alicloud_log_store" "bar" {
  project          = "${alicloud_log_project.foo.name}"
  name             = "${var.name}-source"
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_log_store" "foo" {
  project          = "${alicloud_log_project.foo.name}"
  name             = "${var.name}"
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_mns_topic" "foo" {
  name = "${var.name}"
}
resource "alicloud_fc_service" "foo" {
  name            = "${var.name}"
  internet_access = false
}
resource "alicloud_oss_bucket" "foo" {
  bucket = "${var.name}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "foo" {
  bucket = "${alicloud_oss_bucket.foo.id}"
  key    = "fc/hello.zip"
  source = "./hello.zip"
}
resource "alicloud_fc_function" "foo" {
  service     = "${alicloud_fc_service.foo.name}"
  name        = "${var.name}"
  oss_bucket  = "${alicloud_oss_bucket.foo.id}"
  oss_key     = "${alicloud_oss_bucket_object.foo.key}"
  memory_size = 512
  runtime     = "python2.7"
  handler     = "hello.handler"
}
resource "alicloud_ram_role" "foo" {
  name        = "${var.name}-trigger"
  document    = <<EOF
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
  EOF
  description = "this is a test"
  force       = true
}
resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name   = "${alicloud_ram_role.foo.name}"
  policy_name = "AliyunMNSNotificationRolePolicy"
  policy_type = "System"
}
resource "alicloud_fc_trigger" "foo" {
  service    = "${alicloud_fc_service.foo.name}"
  function   = "${alicloud_fc_function.foo.name}"
  name       = "${var.name}"
  role       = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:mns:${data.alicloud_regions.current_region.regions.0.id}:${data.alicloud_account.current.id}:/topics/${alicloud_mns_topic.foo.name}"
  type       = "mns_topic"
  config_mns = <<EOF
  {
    "filterTag":"testTag",
    "notifyContentFormat":"STREAM",
    "notifyStrategy":"BACKOFF_RETRY"
  }
  EOF
  depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}
```

CDN events trigger:

```terraform
variable "name" {
  default = "fctriggercdneventsconfig"
}

data "alicloud_account" "current" {
}

resource "alicloud_cdn_domain_new" "domain" {
  domain_name = "${var.name}.tf.com"
  cdn_type    = "web"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    type     = "ipaddr"
    priority = 20
    port     = 80
    weight   = 10
  }
}

resource "alicloud_fc_service" "foo" {
  name            = "${var.name}"
  internet_access = false
}
resource "alicloud_oss_bucket" "foo" {
  bucket = "${var.name}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "foo" {
  bucket = "${alicloud_oss_bucket.foo.id}"
  key    = "fc/hello.zip"
  source = "./hello.zip"
}
resource "alicloud_fc_function" "foo" {
  service     = "${alicloud_fc_service.foo.name}"
  name        = "${var.name}"
  oss_bucket  = "${alicloud_oss_bucket.foo.id}"
  oss_key     = "${alicloud_oss_bucket_object.foo.key}"
  memory_size = 512
  runtime     = "python2.7"
  handler     = "hello.handler"
}
resource "alicloud_ram_role" "foo" {
  name        = "${var.name}-trigger"
  document    = <<EOF
    {
        "Version": "1",
        "Statement": [
            {
                "Action": "cdn:Describe*",
                "Resource": "*",
                "Effect": "Allow",
		        "Principal": {
                "Service":
                    ["log.aliyuncs.com"]
                }
            }
        ]
    }
    EOF
  description = "this is a test"
  force       = true
}

resource "alicloud_ram_policy" "foo" {
  name        = "${var.name}-trigger"
  document    = <<EOF
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
    EOF
  description = "this is a test"
  force       = true
}
resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name   = "${alicloud_ram_role.foo.name}"
  policy_name = "${alicloud_ram_policy.foo.name}"
  policy_type = "Custom"
}
resource "alicloud_fc_trigger" "default" {
  service    = "${alicloud_fc_service.foo.name}"
  function   = "${alicloud_fc_function.foo.name}"
  name       = "${var.name}"
  role       = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:cdn:*:${data.alicloud_account.current.id}"
  type       = "cdn_events"
  config     = <<EOF
      {"eventName":"LogFileCreated",
     "eventVersion":"1.0.0",
     "notes":"cdn events trigger",
     "filter":{
        "domain": ["${alicloud_cdn_domain_new.domain.domain_name}"]
        }
    }
EOF
  depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}
```

EventBridge trigger:

```terraform
variable "name" {
  default = "fctriggereventbridgeconfig"
}

data "alicloud_account" "current" {
}

# Please make eventbridge available and then assume a specific service-linked role, which refers to https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/event_bridge_service_linked_role
resource "alicloud_event_bridge_service_linked_role" "service_linked_role" {
  product_name = "AliyunServiceRoleForEventBridgeSendToFC"
}

resource "alicloud_fc_service" "foo" {
  name            = "${var.name}"
  internet_access = false
}
resource "alicloud_oss_bucket" "foo" {
  bucket = "${var.name}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "foo" {
  bucket = "${alicloud_oss_bucket.foo.id}"
  key    = "fc/hello.zip"
  source = "./hello.zip"
}
resource "alicloud_fc_function" "foo" {
  service     = "${alicloud_fc_service.foo.name}"
  name        = "${var.name}"
  oss_bucket  = "${alicloud_oss_bucket.foo.id}"
  oss_key     = "${alicloud_oss_bucket_object.foo.key}"
  memory_size = 512
  runtime     = "python2.7"
  handler     = "hello.handler"
}

resource "alicloud_fc_trigger" "default" {
  service  = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name     = "${var.name}"
  type     = "eventbridge"
  config   = <<EOF
    {
        "triggerEnable": false,
        "asyncInvocationType": false,
        "eventRuleFilterPattern": "{\"source\":[\"acs.oss\"],\"type\":[\"oss:BucketCreated:PutBucket\"]}",
        "eventSourceConfig": {
            "eventSourceType": "Default"
        }
    }
EOF
}

resource "alicloud_fc_trigger" "mns" {
  service  = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name     = "${var.name}"
  type     = "eventbridge"
  config   = <<EOF
    {
        "triggerEnable": false,
        "asyncInvocationType": false,
        "eventRuleFilterPattern": "{}",
        "eventSourceConfig": {
            "eventSourceType": "MNS",
            "eventSourceParameters": {
                "sourceMNSParameters": {
                    "RegionId": "cn-hangzhou",
                    "QueueName": "mns-queue",
                    "IsBase64Decode": true
                }
            }
        }
    }
EOF
}

resource "alicloud_fc_trigger" "rocketmq" {
  service  = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name     = "${var.name}"
  type     = "eventbridge"
  config   = <<EOF
    {
        "triggerEnable": false,
        "asyncInvocationType": false,
        "eventRuleFilterPattern": "{}",
        "eventSourceConfig": {
            "eventSourceType": "RocketMQ",
            "eventSourceParameters": {
                "sourceRocketMQParameters": {
                    "RegionId": "cn-hangzhou",
                    "InstanceId": "MQ_INST_164901546557****_BAAN****",
                    "GroupID": "GID_group1",
                    "Topic": "mytopic",
                    "Timestamp": 1636597951984,
                    "Tag": "test-tag",
                    "Offset": "CONSUME_FROM_LAST_OFFSET"
                }
            }
        }
    }
EOF
}

resource "alicloud_fc_trigger" "rabbitmq" {
  service  = "${alicloud_fc_service.foo.name}"
  function = "${alicloud_fc_function.foo.name}"
  name     = "${var.name}"
  type     = "eventbridge"
  config   = <<EOF
    {
        "triggerEnable": false,
        "asyncInvocationType": false,
        "eventRuleFilterPattern": "{}",
        "eventSourceConfig": {
            "eventSourceType": "RabbitMQ",
            "eventSourceParameters": {
                "sourceRabbitMQParameters": {
                    "RegionId": "cn-hangzhou",
                    "InstanceId": "amqp-cn-****** ",
                    "VirtualHostName": "test-virtual",
                    "QueueName": "test-queue"
                }
            }
        }
    }
EOF
}
```

## Module Support

You can use to the existing [fc module](https://registry.terraform.io/modules/terraform-alicloud-modules/fc/alicloud) 
to create several triggers quickly.

## Argument Reference

The following arguments are supported:

* `service` - (Required, ForceNew) The Function Compute service name.
* `function` - (Required, ForceNew) The Function Compute function name.
* `name` - (ForceNew) The Function Compute trigger name. It is the only in one service and is conflict with "name_prefix".
* `name_prefix` - (ForceNew) Setting a prefix to get a only trigger name. It is conflict with "name".
* `role` - (Optional) RAM role arn attached to the Function Compute trigger. Role used by the event source to call the function. The value format is "acs:ram::$account-id:role/$role-name". See [Create a trigger](https://www.alibabacloud.com/help/doc-detail/53102.htm) for more details.
* `source_arn` - (Optional, ForceNew) Event source resource address. See [Create a trigger](https://www.alibabacloud.com/help/doc-detail/53102.htm) for more details.
* `config` - (Optional) The config of Function Compute trigger.It is valid when `type` is not "mns_topic".See [Configure triggers and events](https://www.alibabacloud.com/help/doc-detail/70140.htm) for more details.
* `config_mns` - (Optional, ForceNew, Available in 1.41.0) The config of Function Compute trigger when the type is "mns_topic".It is conflict with `config`.
* `type` - (Required, ForceNew) The Type of the trigger. Valid values: ["oss", "log", "timer", "http", "mns_topic", "cdn_events", "eventbridge"].

-> **NOTE:** Config does not support modification when type is mns_topic.
-> **NOTE:** type = cdn_events, available in 1.47.0+.
-> **NOTE:** type = eventbridge, available in 1.173.0+.

## Attributes Reference

The following arguments are exported:

* `id` - The ID of the function. The value is formate as `<service>:<function>:<name>`.
* `last_modified` - The date this resource was last modified.
* `trigger_id` - The Function Compute trigger ID.

## Import

Function Compute trigger can be imported using the id, e.g.

```
$ terraform import alicloud_fc_trigger.foo my-fc-service:hello-world:hello-trigger
```
