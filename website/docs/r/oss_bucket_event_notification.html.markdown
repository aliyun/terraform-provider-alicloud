---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_event_notification"
description: |-
  Provides a Alicloud OSS Bucket Event Notification resource.
---

# alicloud_oss_bucket_event_notification

Provides a OSS Bucket Event Notification resource. The event notification configuration routes OSS bucket events (such as object creation or deletion) to Function Compute functions, so that a function is automatically invoked when matching events occur on the bucket.

For information about OSS Bucket Event Notification and how to use it, see [Configure event notification rules](https://www.alibabacloud.com/help/en/oss/user-guide/configure-event-notification-rules).

-> **NOTE:** Available since v2.1.0.

The `PutBucketEventNotification` operation replaces the entire `NotificationConfiguration` of a bucket. This resource therefore owns the full event-notification configuration of one bucket: the `bucket` is the resource identity, and any change to the `function_compute_configuration` list is applied by re-submitting the complete set of configurations. To remove a single configuration, drop its block from the configuration; to clear all configurations, destroy the resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region = var.region
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}-${random_integer.default.result}"
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_fc_service" "default" {
  name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_oss_bucket_object" "code" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "fc/hello.zip"
  content = <<EOF
# -*- coding: utf-8 -*-
def handler(event, context):
    print("hello world")
    return 'hello world'
EOF
}

resource "alicloud_fc_function" "default" {
  service    = alicloud_fc_service.default.name
  name       = "${var.name}-${random_integer.default.result}"
  oss_bucket = alicloud_oss_bucket.default.id
  oss_key    = alicloud_oss_bucket_object.code.key
  runtime    = "python3.10"
  handler    = "hello.handler"
}

resource "alicloud_ram_role" "default" {
  name        = "${var.name}-${random_integer.default.result}"
  document    = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "oss.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
EOF
  description = "Allow OSS to invoke the FC function"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunFCInvocationAccess"
  policy_type = "System"
}

locals {
  function_arn = "acs:fc:${var.region}:${data.alicloud_caller_identity.current.account_id}:services/${alicloud_fc_service.default.name}.LATEST/functions/${alicloud_fc_function.default.name}"
  assume_role  = alicloud_ram_role.default.arn
}

resource "alicloud_oss_bucket_event_notification" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  function_compute_configuration {
    events       = ["oss:ObjectCreated:PutObject"]
    function_arn = local.function_arn
    assume_role  = local.assume_role

    filter {
      key {
        prefix = "source"
        suffix = ".png"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket whose events trigger the notification.
* `function_compute_configuration` - (Required) The list of Function Compute configurations. Each configuration describes one event-to-function rule. PutBucketEventNotification replaces the whole `NotificationConfiguration`, so this list is the complete set owned by the resource. See [`function_compute_configuration`](#function_compute_configuration) below.

### function_compute_configuration

The `function_compute_configuration` block supports:

* `events` - (Required) The list of event types that trigger the function, e.g. `["oss:ObjectCreated:PutObject"]`.
* `function_arn` - (Required) The ARN of the Function Compute function that OSS invokes, in the form `acs:fc:<region>:<account_id>:services/<service>.LATEST/functions/<function>`.
* `assume_role` - (Optional) The ARN of the RAM role that OSS assumes to invoke the function. The role must trust `oss.aliyuncs.com` and be granted function-invocation permissions.
* `filter` - (Optional) The object name filter for the rule. A notification fires only for objects matching the filter. See [`filter`](#filter) below.

#### filter

The `filter` block supports:

* `key` - (Optional) The object key filter. See [`key`](#key) below.

##### key

The `key` block supports:

* `prefix` - (Optional) The object key prefix to match.
* `suffix` - (Optional) The object key suffix to match.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which equals the bucket name. It is formulated as `<bucket>`.
* `function_compute_configuration.0.id` - The server-assigned ID of the configuration rule. This value is set by OSS and may change when the whole `NotificationConfiguration` is replaced; it should not be relied upon as a stable identifier across updates.

## Import

OSS Bucket Event Notification can be imported using the bucket name, e.g.

```shell
terraform import alicloud_oss_bucket_event_notification.example example-bucket-name
```
