---
subcategory: "Internet of Things (Iot)"
layout: "alicloud"
page_title: "Alicloud: alicloud_iot_device_group"
sidebar_current: "docs-alicloud-resource-iot-device-group"
description: |-
  Provides a Alicloud Iot Device Group resource.
---

# alicloud\_iot\_device\_group

Provides a Iot Device Group resource.

For information about Iot Device Group and how to use it, see [What is Device Group](https://www.alibabacloud.com/help/product/30520.htm).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_iot_device_group" "example" {
  group_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `group_desc` - (Optional) The GroupDesc of the device group.
* `group_name` - (Required, ForceNew) The GroupName of the device group.
* `iot_instance_id` - (Optional) The id of the Iot Instance.
* `super_group_id` - (Optional, ForceNew) The id of the SuperGroup.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Device Group.

## Import

Iot Device Group can be imported using the id, e.g.

```
$ terraform import alicloud_iot_device_group.example <id>
```
