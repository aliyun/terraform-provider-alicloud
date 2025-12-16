---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_trigger"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides a Alicloud Function Compute Trigger resource.
---

# alicloud_fc_trigger

Provides an Alicloud Function Compute Trigger resource. Based on trigger, execute your code in response to events in Alibaba Cloud.
 For information about Service and how to use it, see [What is Function Compute](https://www.alibabacloud.com/help/en/fc/developer-reference/api-fc-open-2021-04-06-createtrigger).

-> **NOTE:** The resource requires a provider field 'account_id'. [See account_id](https://www.terraform.io/docs/providers/alicloud/index.html#account_id).

-> **NOTE:** Available since v1.93.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_trigger&exampleId=f10270d0-14d4-416c-f574-29002e67f3a40b4f4fbd&activeTab=example&spm=docs.r.fc_trigger.0.f10270d014&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "default" {}
data "alicloud_regions" "default" {
  current = true
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "default" {
  project_name = "example-value-${random_integer.default.result}"
}

resource "alicloud_log_store" "default" {
  project_name  = alicloud_log_project.default.project_name
  logstore_name = "example-value"
}

resource "alicloud_log_store" "source_store" {
  project_name  = alicloud_log_project.default.project_name
  logstore_name = "example-source-store"
}

resource "alicloud_ram_role" "default" {
  name        = "fcservicerole-${random_integer.default.result}"
  document    = <<EOF
  {
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "fc.aliyuncs.com"
            ]
          }
        }
      ],
      "Version": "1"
  }
  EOF
  description = "this is a example"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_fc_service" "default" {
  name        = "example-value-${random_integer.default.result}"
  description = "example-value"
  role        = alicloud_ram_role.default.arn
  log_config {
    project                 = alicloud_log_project.default.project_name
    logstore                = alicloud_log_store.default.logstore_name
    enable_instance_metrics = true
    enable_request_metrics  = true
  }
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "index.py"
  content = "import logging \ndef handler(event, context): \nlogger = logging.getLogger() \nlogger.info('hello world') \nreturn 'hello world'"
}

resource "alicloud_fc_function" "default" {
  service     = alicloud_fc_service.default.name
  name        = "terraform-example"
  description = "example"
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = "512"
  runtime     = "python3.10"
  handler     = "hello.handler"
}


resource "alicloud_fc_trigger" "default" {
  service    = alicloud_fc_service.default.name
  function   = alicloud_fc_function.default.name
  name       = "terraform-example"
  role       = alicloud_ram_role.default.arn
  source_arn = "acs:log:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:project/${alicloud_log_project.default.project_name}"
  type       = "log"
  config     = <<EOF
    {
        "sourceConfig": {
            "logstore": "${alicloud_log_store.source_store.logstore_name}",
            "startTime": null
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
             "project": "${alicloud_log_project.default.project_name}",
            "logstore": "${alicloud_log_store.default.logstore_name}"
        },
        "enable": true
    }
  
EOF
}
```

MNS topic trigger:

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_trigger&exampleId=04f4d0c3-a376-e4c2-9aab-11d153798b2d63b483eb&activeTab=example&spm=docs.r.fc_trigger.1.04f4d0c3a3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_account" "default" {}
data "alicloud_regions" "default" {
  current = true
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_mns_topic" "default" {
  name = "example-value-${random_integer.default.result}"
}

resource "alicloud_ram_role" "default" {
  name        = "fcservicerole-${random_integer.default.result}"
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
  description = "this is a example"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunMNSNotificationRolePolicy"
  policy_type = "System"
}

resource "alicloud_fc_service" "default" {
  name            = "example-value-${random_integer.default.result}"
  description     = "example-value"
  internet_access = false
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "index.py"
  content = "import logging \ndef handler(event, context): \nlogger = logging.getLogger() \nlogger.info('hello world') \nreturn 'hello world'"
}

resource "alicloud_fc_function" "default" {
  service     = alicloud_fc_service.default.name
  name        = "terraform-example-${random_integer.default.result}"
  description = "example"
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = "512"
  runtime     = "python3.10"
  handler     = "hello.handler"
}

resource "alicloud_fc_trigger" "default" {
  service    = alicloud_fc_service.default.name
  function   = alicloud_fc_function.default.name
  name       = "terraform-example"
  role       = alicloud_ram_role.default.arn
  source_arn = "acs:mns:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:/topics/${alicloud_mns_topic.default.name}"
  type       = "mns_topic"
  config_mns = <<EOF
  {
    "filterTag":"exampleTag",
    "notifyContentFormat":"STREAM",
    "notifyStrategy":"BACKOFF_RETRY"
  }
  EOF
}
```

CDN events trigger:

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_trigger&exampleId=33bbfc00-315e-a4c1-44ac-4df032a1b46a0167ff13&activeTab=example&spm=docs.r.fc_trigger.2.33bbfc0031&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_account" "default" {}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_cdn_domain_new" "default" {
  domain_name = "example${random_integer.default.result}.tf.com"
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

resource "alicloud_fc_service" "default" {
  name            = "example-value-${random_integer.default.result}"
  description     = "example-value"
  internet_access = false
}

resource "alicloud_ram_role" "default" {
  name        = "fcservicerole-${random_integer.default.result}"
  document    = <<EOF
    {
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "cdn.aliyuncs.com"
            ]
          }
        }
      ],
      "Version": "1"
  }
  EOF
  description = "this is a example"
  force       = true
}

resource "alicloud_ram_policy" "default" {
  policy_name     = "fcservicepolicy-${random_integer.default.result}"
  policy_document = <<EOF
    {
        "Version": "1",
        "Statement": [
        {
            "Action": [
            "fc:InvokeFunction"
            ],
        "Resource": [
            "acs:fc:*:*:services/${alicloud_fc_service.default.name}/functions/*",
            "acs:fc:*:*:services/${alicloud_fc_service.default.name}.*/functions/*"
        ],
        "Effect": "Allow"
        }
        ]
    }
    EOF
  description     = "this is a example"
  force           = true
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = alicloud_ram_policy.default.policy_name
  policy_type = "Custom"
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "index.py"
  content = "import logging \ndef handler(event, context): \nlogger = logging.getLogger() \nlogger.info('hello world') \nreturn 'hello world'"
}

resource "alicloud_fc_function" "default" {
  service     = alicloud_fc_service.default.name
  name        = "terraform-example-${random_integer.default.result}"
  description = "example"
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = "512"
  runtime     = "python3.10"
  handler     = "hello.handler"
}

resource "alicloud_fc_trigger" "default" {
  service    = alicloud_fc_service.default.name
  function   = alicloud_fc_function.default.name
  name       = "terraform-example"
  role       = alicloud_ram_role.default.arn
  source_arn = "acs:cdn:*:${data.alicloud_account.default.id}"
  type       = "cdn_events"
  config     = <<EOF
      {"eventName":"LogFileCreated",
     "eventVersion":"1.0.0",
     "notes":"cdn events trigger",
     "filter":{
        "domain": ["${alicloud_cdn_domain_new.default.domain_name}"]
        }
    }
EOF
}
```

EventBridge trigger:

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_trigger&exampleId=5012b39c-cc91-006c-2926-09ba180ba3ac7b724656&activeTab=example&spm=docs.r.fc_trigger.3.5012b39ccc&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "default" {}
data "alicloud_regions" "default" {
  current = true
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_event_bridge_service_linked_role" "service_linked_role" {
  product_name = "AliyunServiceRoleForEventBridgeSendToFC"
}

resource "alicloud_fc_service" "default" {
  name            = "example-value-${random_integer.default.result}"
  description     = "example-value"
  internet_access = false
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "index.py"
  content = "import logging \ndef handler(event, context): \nlogger = logging.getLogger() \nlogger.info('hello world') \nreturn 'hello world'"
}

resource "alicloud_fc_function" "default" {
  service     = alicloud_fc_service.default.name
  name        = "terraform-example"
  description = "example"
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = "512"
  runtime     = "python3.10"
  handler     = "hello.handler"
}

resource "alicloud_fc_trigger" "oss_trigger" {
  service  = alicloud_fc_service.default.name
  function = alicloud_fc_function.default.name
  name     = "terraform-example-oss"
  type     = "eventbridge"
  config = jsonencode(
    {
      "triggerEnable" : false,
      "asyncInvocationType" : false,
      "eventSourceConfig" : {
        "eventSourceType" : "Default"
      },
      "eventRuleFilterPattern" : "{\"source\":[\"acs.oss\"],\"type\":[\"oss:BucketCreated:PutBucket\"]}",
      "eventSinkConfig" : {
        "deliveryOption" : {
          "mode" : "event-driven",
          "eventSchema" : "CloudEvents"
        }
      },
      "runOptions" : {
        "retryStrategy" : {
          "PushRetryStrategy" : "BACKOFF_RETRY"
        },
        "errorsTolerance" : "ALL",
        "mode" : "event-driven"
      }
    }
  )


}

resource "alicloud_fc_trigger" "mns_trigger" {
  service  = alicloud_fc_service.default.name
  function = alicloud_fc_function.default.name
  name     = "terraform-example-mns"
  type     = "eventbridge"
  config = jsonencode(
    {
      "triggerEnable" : false,
      "asyncInvocationType" : false,
      "eventSourceConfig" : {
        "eventSourceType" : "MNS",
        "eventSourceParameters" : {
          "sourceMNSParameters" : {
            "RegionId" : "${data.alicloud_regions.default.regions.0.id}",
            "QueueName" : "mns-queue",
            "IsBase64Decode" : true
          }
        }
      },
      "eventRuleFilterPattern" : "{}",
      "eventSinkConfig" : {
        "deliveryOption" : {
          "mode" : "event-driven",
          "eventSchema" : "CloudEvents"
        }
      },
      "runOptions" : {
        "retryStrategy" : {
          "PushRetryStrategy" : "BACKOFF_RETRY"
        },
        "errorsTolerance" : "ALL",
        "mode" : "event-driven"
      }
    }
  )
}

resource "alicloud_ons_instance" "default" {
  instance_name = "terraform-example-${random_integer.default.result}"
  remark        = "terraform-example"
}
resource "alicloud_ons_group" "default" {
  group_name  = "GID-example"
  instance_id = alicloud_ons_instance.default.id
  remark      = "terraform-example"
}
resource "alicloud_ons_topic" "default" {
  topic_name   = "mytopic"
  instance_id  = alicloud_ons_instance.default.id
  message_type = 0
  remark       = "terraform-example"
}

resource "alicloud_fc_trigger" "rocketmq_trigger" {
  service  = alicloud_fc_service.default.name
  function = alicloud_fc_function.default.name
  name     = "terraform-example-rocketmq"
  type     = "eventbridge"
  config = jsonencode(
    {
      "triggerEnable" : false,
      "asyncInvocationType" : false,
      "eventRuleFilterPattern" : "{}",
      "eventSinkConfig" : {
        "deliveryOption" : {
          "mode" : "event-driven",
          "eventSchema" : "CloudEvents"
        }
      },
      "eventSourceConfig" : {
        "eventSourceType" : "RocketMQ",
        "eventSourceParameters" : {
          "sourceRocketMQParameters" : {
            "RegionId" : "${data.alicloud_regions.default.regions.0.id}",
            "InstanceId" : "${alicloud_ons_instance.default.id}",
            "GroupID" : "${alicloud_ons_group.default.group_name}",
            "Topic" : "${alicloud_ons_topic.default.topic_name}",
            "Timestamp" : 1686296162,
            "Tag" : "example-tag",
            "Offset" : "CONSUME_FROM_LAST_OFFSET"
          }
        }
      },
      "runOptions" : {
        "retryStrategy" : {
          "PushRetryStrategy" : "BACKOFF_RETRY"
        },
        "errorsTolerance" : "ALL",
        "mode" : "event-driven"
      }
    }
  )
}

resource "alicloud_amqp_instance" "default" {
  instance_name  = "terraform-example-${random_integer.default.result}"
  instance_type  = "professional"
  max_tps        = 1000
  queue_capacity = 50
  support_eip    = true
  max_eip_tps    = 128
  payment_type   = "Subscription"
  period         = 1
}

resource "alicloud_amqp_virtual_host" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = "example-VirtualHost"
}

resource "alicloud_amqp_queue" "default" {
  instance_id       = alicloud_amqp_virtual_host.default.instance_id
  queue_name        = "example-queue"
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
}

resource "alicloud_fc_trigger" "rabbitmq_trigger" {
  service  = alicloud_fc_service.default.name
  function = alicloud_fc_function.default.name
  name     = "terraform-example-rabbitmq"
  type     = "eventbridge"
  config = jsonencode(
    {
      "triggerEnable" : false,
      "asyncInvocationType" : false,
      "eventRuleFilterPattern" : "{}",
      "eventSourceConfig" : {
        "eventSourceType" : "RabbitMQ",
        "eventSourceParameters" : {
          "sourceRabbitMQParameters" : {
            "RegionId" : "${data.alicloud_regions.default.regions.0.id}",
            "InstanceId" : "${alicloud_amqp_instance.default.id}",
            "VirtualHostName" : "${alicloud_amqp_virtual_host.default.virtual_host_name}",
            "QueueName" : "${alicloud_amqp_queue.default.queue_name}"
          }
        }
      },
      "eventSinkConfig" : {
        "deliveryOption" : {
          "mode" : "event-driven",
          "eventSchema" : "CloudEvents"
        }
      },
      "runOptions" : {
        "retryStrategy" : {
          "PushRetryStrategy" : "BACKOFF_RETRY"
        },
        "errorsTolerance" : "ALL",
        "mode" : "event-driven"
      }
    }
  )
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_fc_trigger&spm=docs.r.fc_trigger.example&intl_lang=EN_US)

## Module Support

You can use to the existing [fc module](https://registry.terraform.io/modules/terraform-alicloud-modules/fc/alicloud) 
to create several triggers quickly.

## Argument Reference

The following arguments are supported:

* `service` - (Required, ForceNew) The Function Compute service name.
* `function` - (Required, ForceNew) The Function Compute function name.
* `name` - (ForceNew, Optional) The Function Compute trigger name. It is the only in one service and is conflict with "name_prefix".
* `name_prefix` - (ForceNew, Optional) Setting a prefix to get a only trigger name. It is conflict with "name".
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

```shell
$ terraform import alicloud_fc_trigger.foo my-fc-service:hello-world:hello-trigger
```
