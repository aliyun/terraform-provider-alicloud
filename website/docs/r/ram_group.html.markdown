---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_group"
sidebar_current: "docs-alicloud-resource-ram-group"
description: |-
  Provides a RAM Group resource.
---

# alicloud_ram_group

Provides a RAM Group resource.

-> **NOTE:** When you want to destroy this resource forcefully(means remove all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully. 

-> **NOTE:** Available since v1.0.0+.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_group&exampleId=6d7f0720-5959-4789-a198-657d7aa5c525dbbb4ea6&activeTab=example&spm=docs.r.ram_group.0.6d7f072059&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a new RAM Group.
resource "alicloud_ram_group" "group" {
  name     = "groupName"
  comments = "this is a group comments."
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the RAM group. This name can have a string of 1 to 128 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `comments` - (Optional) Comment of the RAM group. This parameter can have a string of 1 to 128 characters.
* `force` - (Optional) This parameter is used for resource destroy. Default value is `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The group ID.

## Import

RAM group can be imported using the id or name, e.g.

```shell
$ terraform import alicloud_ram_group.example my-group
```
