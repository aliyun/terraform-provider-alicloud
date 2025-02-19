---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregator"
sidebar_current: "docs-alicloud-resource-config-aggregator"
description: |-
  Provides a Alicloud Cloud Config Aggregator resource.
---

# alicloud_config_aggregator

Provides a Cloud Config Aggregator resource.

For information about Cloud Config Aggregate Config Rule and how to use it, see [What is Aggregator](https://www.alibabacloud.com/help/en/cloud-config/latest/api-config-2020-09-07-createaggregator).

-> **NOTE:** Available since v1.124.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_aggregator&exampleId=b68e66a2-1930-005b-d7ad-19d3a85a2adcf6c7bb6a&activeTab=example&spm=docs.r.config_aggregator.0.b68e66a219&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}

locals {
  last = length(data.alicloud_resource_manager_accounts.default.accounts) - 1
}

resource "alicloud_config_aggregator" "default" {
  aggregator_accounts {
    account_id   = data.alicloud_resource_manager_accounts.default.accounts[local.last].account_id
    account_name = data.alicloud_resource_manager_accounts.default.accounts[local.last].display_name
    account_type = "ResourceDirectory"
  }
  aggregator_name = var.name
  description     = var.name
  aggregator_type = "CUSTOM"
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_accounts` - (Optional) The information of account in aggregator. If the aggregator_type is RD, it is optional and means add all members in the resource directory to the account group. See [`aggregator_accounts`](#aggregator_accounts) below.  **NOTE:** the field `aggregator_accounts` is not required from version 1.148.0.
* `aggregator_name` - (Required) The name of aggregator.
* `aggregator_type` - (Optional, ForceNew) The type of aggregator. Valid values: `CUSTOM`, `RD`. The Default value: `CUSTOM`.
  * `CUSTOM` - The custom account group.
  * `RD` - The global account group.
* `description` - (Required) The description of aggregator.

### `aggregator_accounts`

The aggregator_accounts supports the following: 

* `account_id` - (Required) Aggregator account Uid.
* `account_name` - (Required) Aggregator account name.
* `account_type` - (Required) Aggregator account source type. Valid values: `ResourceDirectory`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Aggregator.
* `status` - The status of the resource. Valid values: `0`: creating `1`: normal `2`: deleting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Aggregator.

## Import

Cloud Config Aggregator can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_aggregator.example <id>
```
