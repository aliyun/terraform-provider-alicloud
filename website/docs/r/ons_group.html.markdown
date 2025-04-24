---
subcategory: "RocketMQ (Ons)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_group"
sidebar_current: "docs-alicloud-resource-ons-group"
description: |-
  Provides a Alicloud ONS Group resource.
---

# alicloud\_ons\_group

Provides an ONS group resource.

For more information about how to use it, see [RocketMQ Group Management API](https://www.alibabacloud.com/help/doc-detail/29616.html). 

-> **NOTE:** Available in 1.53.0+

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ons_group&exampleId=25255a87-5baf-76e0-83b5-ffe9bc1441fdd8202d77&activeTab=example&spm=docs.r.ons_group.0.25255a875b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "GID-tf-example"
}

variable "group_name" {
  default = "GID-tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ons_instance" "default" {
  instance_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_ons_group" "default" {
  group_name  = var.group_name
  instance_id = alicloud_ons_instance.default.id
  remark      = "dafault_ons_group_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the groups.
* `group_id` - (Optional) Replaced by `group_name` after version 1.98.0.
* `group_name` - (Optional) Name of the group. Two groups on a single instance cannot have the same name. A `group_name` starts with "GID_" or "GID-", and contains letters, numbers, hyphens (-), and underscores (_).
* `group_type` - (Optional) Specify the protocol applicable to the created Group ID. Valid values: `tcp`, `http`. Default to `tcp`.
* `remark` - (Optional) This attribute is a concise description of group. The length cannot exceed 256.
* `read_enable` - (Optional) This attribute is used to set the message reading enabled or disabled. It can only be set after the group is used by the client.
* `tags` - (Optional, Available in v1.98.0+) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<group_name>`.

## Timeouts

-> **NOTE:** Available in 1.98.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 4 mins) Used when Creating ONS instance. 
* `delete` - (Defaults to 4 mins) Used when terminating the ONS instance. 

## Import

ONS GROUP can be imported using the id, e.g.

```shell
$ terraform import alicloud_ons_group.group MQ_INST_1234567890_Baso1234567:GID-onsGroupDemo
```
