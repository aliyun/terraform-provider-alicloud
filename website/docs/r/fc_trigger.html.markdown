---
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_trigger"
sidebar_current: "docs-alicloud-resource-fc-trigger"
description: |-
  Provides a Alicloud Function Compute Trigger resource.
---

# alicloud\_fc\_function

Provides a Alicloud Function Compute Trigger resource. Based on trigger, execute your code in response to events in Alibaba Cloud.
 For information about Service and how to use it, see [What is Function Compute](https://www.alibabacloud.com/help/doc-detail/52895.htm).

-> **NOTE:** The resource requires a provider field 'account_id'. [See account_id](https://www.terraform.io/docs/providers/alicloud/index.html#account_id).

## Example Usage

Basic Usage

```
variable "region" {
  default = "cn-hangzhou"
}
variable "account" {
  default = "12345"
}

provider "alicloud" {
  account_id = "${var.account}"
  region = "${var.region}"
}

resource "alicloud_fc_trigger" "foo" {
  service = "my-fc-service"
  function = "hello-world"
  name = "hello-trigger"
  role = "${alicloud_ram_role.foo.arn}"
  source_arn = "acs:log:${var.region}:${var.account}:project/${alicloud_log_project.foo.name}"
  type = "log"
  config = <<EOF
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
  depends_on = ["alicloud_ram_role_policy_attachment.foo"]
}

resource "alicloud_ram_role" "foo" {
  name = "${var.name}-trigger"
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
  force = true
}

resource "alicloud_ram_role_policy_attachment" "foo" {
  role_name = "${alicloud_ram_role.foo.name}"
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

```
## Argument Reference

The following arguments are supported:

* `service` - (Required, ForceNew) The Function Compute service name.
* `function` - (Required, ForceNew) The Function Compute function name.
* `name` - (ForceNew) The Function Compute trigger name. It is the only in one service and is conflict with "name_prefix".
* `name_prefix` - (ForceNew) Setting a prefix to get a only trigger name. It is conflict with "name".
* `role` - (Optional) RAM role arn attached to the Function Compute trigger. Role used by the event source to call the function. The value format is "acs:ram::$account-id:role/$role-name". See [Create a trigger](https://www.alibabacloud.com/help/doc-detail/53102.htm) for more details.
* `source_arn` - (Optional) Event source resource address. See [Create a trigger](https://www.alibabacloud.com/help/doc-detail/53102.htm) for more details.
* `config` - (Optional) The config of Function Compute trigger. See [Configure triggers and events](https://www.alibabacloud.com/help/doc-detail/70140.htm) for more details.
* `type` - (Required, ForceNew) The Type of the trigger. Valid values: ["oss", "log", "timer", "http"].

## Attributes Reference

The following arguments are exported:

* `id` - The ID of the function. The value is formate as "<service>:<function>:<name>".
* `last_modified` - The date this resource was last modified.

## Import

Function Compute trigger can be imported using the id, e.g.

```
$ terraform import alicloud_fc_service.foo my-fc-service:hello-world:hello-trigger
```
