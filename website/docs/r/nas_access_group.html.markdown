---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_group"
description: |-
  Provides a Alicloud File Storage (NAS) Access Group resource.
---

# alicloud_nas_access_group

Provides a File Storage (NAS) Access Group resource.

File system Access Group.

For information about File Storage (NAS) Access Group and how to use it, see [What is Access Group](https://www.alibabacloud.com/help/en/nas/developer-reference/api-nas-2017-06-26-createaccessgroup).

-> **NOTE:** Available since v1.33.0.

## Example Usage

Basic Usage

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_nas_access_group" "foo" {
  access_group_name = "terraform-example-${random_integer.default.result}"
  access_group_type = "Vpc"
  description       = "terraform-example"
  file_system_type  = "extreme"
}
```

## Argument Reference

The following arguments are supported:
* `access_group_name` - (Optional, ForceNew) The name of the permission group.
* `access_group_type` - (Optional, ForceNew) Permission group types, including Vpc.
* `description` - (Optional) Permission group description information.
* `file_system_type` - (Optional, ForceNew, Computed) File system type.

Value:
  - standard (default): Universal NAS
  - extreme: extreme NAS


The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.218.0) Field 'name' has been deprecated from provider version 1.218.0. New field 'access_group_name' instead.
* `type` - (Deprecated since v1.218.0) Field 'type' has been deprecated from provider version 1.218.0. New field 'access_group_type' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<access_group_name>:<file_system_type>`.
* `create_time` - (Available since v1.218.0) Creation time.
* `region_id` - (Available since v1.256.0) The region ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Group.
* `delete` - (Defaults to 5 mins) Used when delete the Access Group.
* `update` - (Defaults to 5 mins) Used when update the Access Group.

## Import

File Storage (NAS) Access Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_access_group.example <access_group_name>:<file_system_type>
```