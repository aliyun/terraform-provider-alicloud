---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_trigger"
description: |-
  Provides a Alicloud FCV3 Trigger resource.
---

# alicloud_fcv3_trigger

Provides a FCV3 Trigger resource.

A trigger is a way of triggering the execution of a function. In the event-driven computing model, the event source is the producer of the event, the function is the handler of the event, and the trigger provides a centralized and unified way to manage different event sources. In the event source, when the event occurs, if the rules defined by the trigger are met,.

For information about FCV3 Trigger and how to use it, see [What is Trigger](https://www.alibabacloud.com/help/en/functioncompute/api-fc-2023-03-30-createtrigger).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "function_name" {
  default = "TerraformTriggerResourceAPI"
}

variable "trigger_name" {
  default = "TerraformTrigger_CDN"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

data "alicloud_account" "current" {
}

resource "alicloud_fcv3_trigger" "default" {
  trigger_type    = "cdn_events"
  trigger_name    = var.name
  description     = "create"
  qualifier       = "LATEST"
  trigger_config  = jsonencode({ "eventName" : "CachedObjectsPushed", "eventVersion" : "1.0.0", "notes" : "example", "filter" : { "domain" : ["example.com"] } })
  source_arn      = "acs:cdn:*:${data.alicloud_account.current.id}"
  invocation_role = "acs:ram::${data.alicloud_account.current.id}:role/aliyuncdneventnotificationrole"
  function_name   = alicloud_fcv3_function.function.function_name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description of the trigger
* `function_name` - (Required, ForceNew) Function Name
* `invocation_role` - (Optional) The role required by the event source (such as OSS) to call the function.
* `qualifier` - (Required, ForceNew) The version or alias of the function
* `source_arn` - (Optional, ForceNew) Trigger Event source ARN
* `trigger_config` - (Optional) Trigger configuration. The configuration varies for different types of triggers.
* `trigger_name` - (Optional, ForceNew, Computed) Trigger Name
* `trigger_type` - (Required, ForceNew) The type of the trigger. Currently, the supported types include oss, log, http, timer, tablestore, cdn_events, mns_topic and eventbridge.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<function_name>:<trigger_name>`.
* `create_time` - Creation time
* `status` - The state of the trigger

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Trigger.
* `delete` - (Defaults to 5 mins) Used when delete the Trigger.
* `update` - (Defaults to 5 mins) Used when update the Trigger.

## Import

FCV3 Trigger can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_trigger.example <function_name>:<trigger_name>
```