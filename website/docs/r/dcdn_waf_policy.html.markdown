---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_policy"
sidebar_current: "docs-alicloud-resource-dcdn-waf-policy"
description: |-
  Provides a Alicloud DCDN Waf Policy resource.
---

# alicloud\_dcdn\_waf\_policy

Provides a DCDN Waf Policy resource.

For information about DCDN Waf Policy and how to use it, see [What is Waf Policy](https://www.alibabacloud.com/help/en/dynamic-route-for-cdn/latest/set-the-protection-policies#doc-api-dcdn-CreateDcdnWafPolicy).

-> **NOTE:** Available in v1.184.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dcdn_waf_policy" "example" {
  defense_scene = "waf_group"
  policy_name   = var.name
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Waf Policy.
* `delete` - (Defaults to 1 mins) Used when deleting the Waf Policy.
* `update` - (Defaults to 1 mins) Used when updating the Waf Policy.

## Import

DCDN Waf Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_waf_policy.example <id>
```