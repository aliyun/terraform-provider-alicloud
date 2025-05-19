---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_user_config"
description: |-
  Provides a Alicloud PAI Workspace User Config resource.
---

# alicloud_pai_workspace_user_config

Provides a PAI Workspace User Config resource.



For information about PAI Workspace User Config and how to use it, see [What is User Config](https://www.alibabacloud.com/help/en/pai/developer-reference/api-aiworkspace-2021-02-04-setuserconfigs).

-> **NOTE:** Available since v1.250.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_workspace_user_config" "default" {
  category_name = "DataPrivacyConfig"
  config_key    = "customizePAIAssumedRole"
  config_value  = var.name
}
```

## Argument Reference

The following arguments are supported:
* `category_name` - (Required, ForceNew) The category. Valid values: `DataPrivacyConfig`.
* `config_key` - (Required, ForceNew) The key of the configuration.
* `config_value` - (Required) The value of the configuration.
* `scope` - (Optional, ForceNew) The scope. Default value: `owner`. Valid values: `owner`, `subUser`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<category_name>:<config_key>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the User Config.
* `delete` - (Defaults to 5 mins) Used when delete the User Config.
* `update` - (Defaults to 5 mins) Used when update the User Config.

## Import

PAI Workspace User Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_user_config.example <category_name>:<config_key>
```
