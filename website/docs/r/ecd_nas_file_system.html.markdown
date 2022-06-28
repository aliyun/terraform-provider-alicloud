---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_nas_file_system"
sidebar_current: "docs-alicloud-resource-ecd-nas-file-system"
description: |-
  Provides a Alicloud ECD Nas File System resource.
---

# alicloud\_ecd\_nas\_file\_system

Provides a ECD Nas File System resource.

For information about ECD Nas File System and how to use it, see [What is Nas File System](https://help.aliyun.com/document_detail/188382.html).

-> **NOTE:** Available in v1.141.0+.

## Example Usage

Basic Usage

```terraform

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block             = "172.16.0.0/12"
  desktop_access_type    = "Internet"
  office_site_name       = "your_office_site_name"
  enable_internet_access = false
}

resource "alicloud_ecd_nas_file_system" "example" {
  nas_file_system_name = "example_value"
  office_site_id       = alicloud_ecd_simple_office_site.default.id
  description          = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Optional) The filesystem id of nas file system.
* `description` - (Optional, ForceNew) The description of nas file system.
* `mount_target_domain` - (Optional) The domain of mount target.
* `nas_file_system_name` - (Optional, ForceNew) The name of nas file system.
* `office_site_id` - (Required, ForceNew) The ID of office site.
* `reset` - (Optional) The mount point is in an inactive state, reset the mount point of the NAS file system. Default to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Nas File System.
* `status` - The status of nas file system. Valid values: `Pending`, `Running`, `Stopped`,`Deleting`, `Deleted`, `Invalid`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Nas File System.
* `delete` - (Defaults to 10 mins) Used when delete the Nas File System.
* `update` - (Defaults to 10 mins) Used when update the Nas File System.

## Import

ECD Nas File System can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_nas_file_system.example <id>
```
