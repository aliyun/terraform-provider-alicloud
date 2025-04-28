---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_security_group"
description: |-
  Provides a Alicloud ENS Security Group resource.
---

# alicloud_ens_security_group

Provides a ENS Security Group resource. 

For information about ENS Security Group and how to use it, see [What is Security Group](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createsnapshot).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_security_group&exampleId=b5a40fc7-1e42-e096-f7a2-0e286603674b1fd46a98&activeTab=example&spm=docs.r.ens_security_group.0.b5a40fc71e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_security_group" "default" {
  description         = var.name
  security_group_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Security group description informationIt must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`.
* `security_group_name` - (Optional) Security group nameThe security group name. The length is 2~128 English or Chinese characters. It must start with an uppercase or lowcase letter or a Chinese character and cannot start with `http://` or `https`. Can contain digits, colons (:), underscores (_), or hyphens (-).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Security Group.
* `delete` - (Defaults to 5 mins) Used when delete the Security Group.
* `update` - (Defaults to 5 mins) Used when update the Security Group.

## Import

ENS Security Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_security_group.example <id>
```