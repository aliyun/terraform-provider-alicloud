---
subcategory: "MaxCompute"
layout: "alicloud"
page_title: "Alicloud: alicloud_maxcompute_project"
sidebar_current: "docs-alicloud-resource-maxcompute-project"
description: |-
  Provides a Alicloud maxcompute project resource.
---

# alicloud\_maxcompute\_project

The project is the basic unit of operation in maxcompute. It is similar to the concept of Database or Schema in traditional databases, and sets the boundary for maxcompute multi-user isolation and access control. [Refer to details](https://www.alibabacloud.com/help/doc-detail/27818.html).

->**NOTE:** Available in 1.77.0+.

## Example Usage

Basic Usage

```
resource "alicloud_maxcompute_project" "example" {
  name               = "tf_maxcompute_project"
  specification_type = "OdpsStandard"
  order_type         = "PayAsYouGo"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the maxcompute project. 
* `specification_type` - (Required)  The type of resource Specification, only `OdpsStandard` supported currently.
* `order_type` - (Required) The type of payment, only `PayAsYouGo` supported currently.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the maxcompute project. It is the same as its name.

## Import

MaxCompute project can be imported using the *name* or ID, e.g.

```
$ terraform import alicloud_maxcompute_project.example tf_maxcompute_project
```
