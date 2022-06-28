---
subcategory: "Schedulerx"
layout: "alicloud"
page_title: "Alicloud: alicloud_schedulerx_namespace"
sidebar_current: "docs-alicloud-resource-schedulerx-namespace"
description: |- 
    Provides a Alicloud Schedulerx Namespace resource.
---

# alicloud\_schedulerx\_namespace

Provides a Schedulerx Namespace resource.

For information about Schedulerx Namespace and how to use it, see [What is Namespace](https://help.aliyun.com/document_detail/206088.html).

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_schedulerx_namespace" "example" {
  namespace_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the resource.
* `namespace_name` - (Required) The name of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace. Its value is same as `namespace_id`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the resource.
* `update` - (Defaults to 1 mins) Used when update the resource.



## Import

Schedulerx Namespace can be imported using the id, e.g.

```
$ terraform import alicloud_schedulerx_namespace.example <id>
```