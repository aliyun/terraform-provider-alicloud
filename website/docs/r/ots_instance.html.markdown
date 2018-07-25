---
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_instance"
sidebar_current: "docs-alicloud-resource-ots-instance"
description: |-
  Provides an OTS (Open Table Service) instance resource.
---

# alicloud\_ots\_instance

This resource will help you to manager a [Table Store](https://www.alibabacloud.com/help/doc-detail/27280.htm) Instance.
It is foundation of creating data table.

## Example Usage

```
# Create an OTS instance
resource "alicloud_ots_instance" "foo" {
  name = "my-ots-instance"
  description = "for table"
  accessed_by = "Vpc"
  tags {
    Created = "TF"
    For = "Building table"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the instance.
* `accessed_by` - The network limitation of accessing instance. Valid values:
    * `Any` - Allow all network to access the instance.
    * `Vpc` - Only can the attached VPC allow to access the instance.
    * `ConsoleOrVpc` - Allow web console or the attached VPC to access the instance.

    Default to "Any".
* `instance_type` - (ForceNew) The type of instance. Valid values are "Capacity" and "HighPerformance". Default to "HighPerformance".
* `description` - (Required, ForceNew) The description of the instance.
* `tags` - A mapping of tags to assign to the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is same as the "name".
* `name` - The instance name.
* `description` - The instance description.
* `accessed_by` - TThe network limitation of accessing instance.
* `instance_type` - The instance type.
* `tags` - The instance tags.

## Import

OTS instance can be imported using instance id or name, e.g.

```
$ terraform import alicloud_ots_instance.foo "my-ots-instance"
```

