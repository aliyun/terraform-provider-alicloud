---
subcategory: "PAI Plugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_group"
sidebar_current: "docs-alicloud-resource-pai-plugin-group"
description: |-
  Provides a PAI Plugin Group resource.
---

# alicloud_pai_plugin_group

Provides a PAI Plugin Group resource.

For information about PAI Plugin Group and how to use it, see [PAI Plugin - CreateGroup](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiplugin-2022-01-12-creategroup).

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_plugin_group" "default" {
  name   = var.name
  source = 0
  text   = "13900000000"
  remark = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the group.
* `source` - (Optional, ForceNew) The data source of the group. Valid values: `0` to `4`. `0`: text (one phone number per line, up to 100); `4`: algorithm output.
* `project` - (Optional, ForceNew) The project of the group.
* `algorithm` - (Optional, ForceNew) The algorithm of the group, required when `source` is `4`.
* `filter` - (Optional, ForceNew) The filter of the group.
* `column` - (Optional, ForceNew) The column of the group.
* `table` - (Optional, ForceNew) The table of the group.
* `text` - (Optional, ForceNew) The text content of the group, one phone number per line and up to 100 when `source` is `0`.
* `uri` - (Optional, ForceNew) The URI of the group.
* `remark` - (Optional, ForceNew) The remark of the group.

-> **NOTE:** A group has no in-place update. Changing any argument forces a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `plugin_group_id`.
* `plugin_group_id` - The ID of the PAI Plugin Group.
* `phone_number` - Whether the group contains phone numbers.
* `amount` - The amount of the group.
* `campaign_id` - The campaign ID of the group.
* `created_time` - The time when the group was created.
* `updated_time` - The time when the group was last updated.
* `status` - The status of the group.
* `inference_job_id` - The inference job ID of the group.
* `history` - The history of the group.

## Import

PAI Plugin Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_plugin_group.example <plugin_group_id>
```
