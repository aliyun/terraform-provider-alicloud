---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_account_logging"
description: |-
  Provides a Alicloud Message Service Account Logging resource.
---

# alicloud_message_service_account_logging

Provides a Message Service Account Logging resource.

Account Logging Configuration.

For information about Message Service Account Logging and how to use it, see [What is Account Logging](https://next.api.alibabacloud.com/document/Mns-open/2022-01-19/SetAccountAttributes).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_account_logging&exampleId=5dc46c29-c718-916b-c43d-fb134b9dae10494c35a6&activeTab=example&spm=docs.r.message_service_account_logging.0.5dc46c29c7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
  description  = "example project for account logging"
}

resource "alicloud_log_store" "default" {
  project_name     = alicloud_log_project.default.project_name
  logstore_name    = var.name
  retention_period = 30
  shard_count      = 2
}

resource "alicloud_message_service_account_logging" "default" {
  log_enabled           = true
  message_trace_enabled = false
  project_name          = alicloud_log_project.default.project_name
  log_store_name        = alicloud_log_store.default.logstore_name
}
```

### Deleting `alicloud_message_service_account_logging` or removing it from your configuration

Terraform cannot destroy resource `alicloud_message_service_account_logging`. Terraform will remove this resource from the state file, however resources may remain.


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_message_service_account_logging&spm=docs.r.message_service_account_logging.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `log_enabled` - (Required) Whether to enable delivering the account operation logs to Simple Log Service (SLS).
* `log_store_name` - (Optional) The name of the SLS Logstore that receives the logs. Required when `log_enabled` is `true`.
* `project_name` - (Optional) The name of the SLS Project that the Logstore belongs to. Required when `log_enabled` is `true`.
* `message_trace_enabled` - (Optional) Whether to enable the message trace feature.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<Alibaba Cloud Account ID>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account Logging.
* `update` - (Defaults to 5 mins) Used when update the Account Logging.

## Import

Message Service Account Logging can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_account_logging.example <Alibaba Cloud Account ID>
```