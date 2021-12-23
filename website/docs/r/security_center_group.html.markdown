---
subcategory: "Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_center_group"
sidebar_current: "docs-alicloud-resource-security-center-group"
description: |-
  Provides a Alicloud Security Center Group resource.
---

# alicloud\_security\_center\_group

Provides a Security Center Group resource.

For information about Security Center Group and how to use it, see [What is Group](https://www.alibabacloud.com/help/doc-detail/129195.htm).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_security_center_group" "example" {
  group_name = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Computed, ForceNew) GroupId.
* `group_name` - (Optional) GroupName.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Group. Its value is same as `group_id`.

## Import

Security Center Group can be imported using the id, e.g.

```
$ terraform import alicloud_security_center_group.example <group_id>
```
