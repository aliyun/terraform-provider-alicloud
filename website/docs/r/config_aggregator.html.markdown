---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregator"
description: |-
  Provides a Alicloud Cloud Config (Config) Aggregator resource.
---

# alicloud_config_aggregator

Provides a Cloud Config (Config) Aggregator resource.



For information about Cloud Config (Config) Aggregator and how to use it, see [What is Aggregator](https://www.alibabacloud.com/help/en/cloud-config/latest/api-config-2020-09-07-createaggregator).

-> **NOTE:** Available since v1.124.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_aggregator&exampleId=649b656e-4929-2258-a2fd-fccccad863e8f43eefd3&activeTab=example&spm=docs.r.config_aggregator.0.649b656e49&intl_lang=EN_US" target="_blank">
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_config_aggregator&spm=docs.r.config_aggregator.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `aggregator_accounts` - (Optional, Set) The member accounts of the account group. See [`aggregator_accounts`](#aggregator_accounts) below.
-> **NOTE:** If `aggregator_type` is set to `CUSTOM`, `aggregator_accounts` is required.
* `aggregator_name` - (Required) The name of the account group.
* `aggregator_type` - (Optional, ForceNew) The type of the account group. Default value: `CUSTOM`. Valid values:
  - `RD`: Global account group.
  - `FOLDER`: Folder account group.
  - `CUSTOM`: Custom account group.
* `description` - (Required) The description of the account group.
* `folder_id` - (Optional, Available since v1.262.0) The ID of the attached folder. You can specify multiple folder IDs. Separate the IDs with commas (,). **NOTE:** If `aggregator_type` is set to `FOLDER`, `folder_id` is required.

### `aggregator_accounts`

The aggregator_accounts supports the following:
* `account_id` - (Optional) The member ID.
* `account_name` - (Optional) The member name.
* `account_type` - (Optional) The affiliation of the member. Valid values: `ResourceDirectory`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - (Available since v1.262.0) The timestamp when the account group was created.
* `status` - The status of the account group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Aggregator.
* `delete` - (Defaults to 5 mins) Used when delete the Aggregator.
* `update` - (Defaults to 5 mins) Used when update the Aggregator.

## Import

Cloud Config (Config) Aggregator can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_aggregator.example <id>
```
