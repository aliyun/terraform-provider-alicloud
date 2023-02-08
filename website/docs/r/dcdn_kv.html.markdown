---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_kv"
sidebar_current: "docs-alicloud-resource-dcdn-kv"
description: |-
  Provides a Alicloud Dcdn Kv resource.
---

# alicloud_dcdn_kv

Provides a Dcdn Kv resource.

For information about Dcdn Kv and how to use it, see [What is Kv](https://www.alibabacloud.com/help/en/dynamic-route-for-cdn/latest/putdcdnkv).

-> **NOTE:** Available in v1.198.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dcdn_kv_account" "default" {
  status = "online"
}
resource "alicloud_dcdn_kv_namespace" "default" {
  description = "wkmtest"
  namespace   = var.name
}
resource "alicloud_dcdn_kv" "default" {
  value     = "testvalue"
  key       = var.name
  namespace = alicloud_dcdn_kv_namespace.default.namespace
}
```

## Argument Reference

The following arguments are supported:
* `key` - (Required,ForceNew) The name of the key to Put, the longest 512, cannot contain spaces.
* `namespace` - (Required,ForceNew) The name specified when the customer calls PutDcdnKvNamespace
* `value` - (Required) The content of key, up to 2M(2*1000*1000)



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<namespace>:<key>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Kv.
* `delete` - (Defaults to 5 mins) Used when delete the Kv.
* `update` - (Defaults to 5 mins) Used when update the Kv.

## Import

Dcdn Kv can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_kv.example <namespace>:<key>
```