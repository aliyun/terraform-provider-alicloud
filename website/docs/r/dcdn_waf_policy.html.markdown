---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_policy"
sidebar_current: "docs-alicloud-resource-dcdn-waf-policy"
description: |-
  Provides a Alicloud DCDN Waf Policy resource.
---

# alicloud_dcdn_waf_policy

Provides a DCDN Waf Policy resource.

For information about DCDN Waf Policy and how to use it, see [What is Waf Policy](https://www.alibabacloud.com/help/en/dcdn/developer-reference/api-dcdn-2018-01-15-createdcdnwafpolicy).

-> **NOTE:** Available since v1.184.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dcdn_waf_policy&exampleId=3bec2598-8763-5264-d657-6b3a56f9d056bd852371&activeTab=example&spm=docs.r.dcdn_waf_policy.0.3bec259887&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dcdn_waf_policy" "example" {
  defense_scene = "waf_group"
  policy_name   = "${var.name}_${random_integer.default.result}"
  policy_type   = "custom"
  status        = "on"
}
```

## Argument Reference

The following arguments are supported:

* `defense_scene` - (Required, ForceNew) The type of protection policy. Valid values: `waf_group`, `custom_acl`, `whitelist`, `ip_blacklist`, `region_block`.
* `policy_name` - (Required) The name of the protection policy. The name must be 1 to 64 characters in length, and can contain letters, digits,and underscores (_).
* `policy_type` - (Required, ForceNew) The type of the protection policy. Valid values: `default`, `custom`.
* `status` - (Required) The status of the resource. Valid values: `on`, `off`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Waf Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Waf Policy.
* `delete` - (Defaults to 1 mins) Used when deleting the Waf Policy.
* `update` - (Defaults to 1 mins) Used when updating the Waf Policy.

## Import

DCDN Waf Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_waf_policy.example <id>
```