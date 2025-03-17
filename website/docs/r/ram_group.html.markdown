---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_group"
description: |-
  Provides a Alicloud RAM Group resource.
---

# alicloud_ram_group

Provides a RAM Group resource.

The group that users can join.

For information about RAM Group and how to use it, see [What is Group](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-creategroup).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_group&exampleId=ad870fa2-5bad-840b-bcba-2e5e5ca51ff785abc0bc&activeTab=example&spm=docs.r.ram_group.0.ad870fa25b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ram_group" "group" {
  group_name = var.name
  comments   = var.name
  force      = true
}
```

## Argument Reference

The following arguments are supported:

* `comments` - (Optional) The description of the user group. The value can be up to 128 characters in length.
* `group_name` - (Optional, ForceNew, Available since v1.245.0) The name of the user group. The name must be `1` to `64` characters in length and can contain letters, digits, periods (.), hyphens (-), and underscores (_).
* `force` - (Optional, Bool) Specifies whether to force delete the Group. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `name` - (Optional, ForceNew, Deprecated since v1.120.0) Field `name` has been deprecated from provider version 1.245.0. New field `group_name` instead.
-> **NOTE:** From version 1.245.0, You must specify at least one of the `group_name` and `name`.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Group.
* `create_time` - (Available since v1.245.0) The time when the RAM User Group was created.

## Timeouts

-> **NOTE:** Available since v1.245.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Group.
* `delete` - (Defaults to 5 mins) Used when delete the Group.
* `update` - (Defaults to 5 mins) Used when update the Group.

## Import

RAM Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_group.example <id>
```
