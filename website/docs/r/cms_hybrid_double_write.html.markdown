---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_double_write"
sidebar_current: "docs-alicloud-resource-cms-hybrid-double-write"
description: |-
  Provides a Alicloud Cloud Monitor Service Hybrid Double Write resource.
---

# alicloud_cms_hybrid_double_write

Provides a Cloud Monitor Service Hybrid Double Write resource.

For information about Cloud Monitor Service Hybrid Double Write and how to use it, see [What is Hybrid Double Write](https://next.api.aliyun.com/document/Cms/2018-03-08/CreateHybridDoubleWrite).

-> **NOTE:** Available in v1.204.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_account" "default" {}

resource "alicloud_cms_namespace" "default" {
  count         = 2
  description   = var.name
  namespace     = "${var.name}-${count.index}"
  specification = "cms.s1.large"
}

resource "alicloud_cms_hybrid_double_write" "default" {
  user_id          = data.alicloud_account.default.id
  source_namespace = alicloud_cms_namespace.default.0.namespace
  namespace        = alicloud_cms_namespace.default.1.namespace
}
```

## Argument Reference

The following arguments are supported:
* `namespace` - (ForceNew,Required) Double write target Namespace.
* `user_id` - (ForceNew,Required) The double written target user ID.
* `source_namespace` - (ForceNew,Required) Double write source namespace.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is same as `source_namespace`.
* `source_user_id` - The id of the source user.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Hybrid Double Write.
* `delete` - (Defaults to 5 mins) Used when delete the Hybrid Double Write.

## Import

Cloud Monitor Service Hybrid Double Write can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_hybrid_double_write.example <id>
```