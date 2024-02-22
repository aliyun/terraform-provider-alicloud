---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_hybrid_double_write"
description: |-
  Provides a Alicloud Cloud Monitor Service Hybrid Double Write resource.
---

# alicloud_cloud_monitor_service_hybrid_double_write

Provides a Cloud Monitor Service Hybrid Double Write resource. 

For information about Cloud Monitor Service Hybrid Double Write and how to use it, see [What is Hybrid Double Write](https://next.api.alibabacloud.com/document/Cms/2018-03-08/CreateHybridDoubleWrite).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_account" "default" {
}

resource "alicloud_cms_namespace" "source" {
  namespace = var.name
}

resource "alicloud_cms_namespace" "default" {
  namespace = "${var.name}-source"
}

resource "alicloud_cloud_monitor_service_hybrid_double_write" "default" {
  source_namespace = alicloud_cms_namespace.source.id
  source_user_id   = data.alicloud_account.default.id
  namespace        = alicloud_cms_namespace.default.id
  user_id          = data.alicloud_account.default.id
}
```

## Argument Reference

The following arguments are supported:

* `source_namespace` - (Required, ForceNew) Source Namespace.
* `source_user_id` - (Required, ForceNew) Source UserId.
* `namespace` - (Required, ForceNew) Target Namespace.
* `user_id` - (Required, ForceNew) Target UserId.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hybrid Double Write. It formats as `<source_namespace>:<source_user_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Hybrid Double Write.
* `delete` - (Defaults to 5 mins) Used when delete the Hybrid Double Write.

## Import

Cloud Monitor Service Hybrid Double Write can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_monitor_service_hybrid_double_write.example <source_namespace>:<source_user_id>
```
