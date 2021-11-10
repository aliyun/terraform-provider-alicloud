---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_group"
sidebar_current: "docs-alicloud-resource-bastionhost-host-group"
description: |-
  Provides a Alicloud Bastion Host Host Group resource.
---

# alicloud\_bastionhost\_host\_group

Provides a Bastion Host Host Group resource.

For information about Bastion Host Host Group and how to use it, see [What is Host Group](https://www.alibabacloud.com/help/en/doc-detail/204307.htm).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_host_group" "example" {
  host_group_name = "example_value"
  instance_id     = "bastionhost-cn-tl3xxxxxxx"
}

```

## Argument Reference

The following arguments are supported:

* `comment` - (Optional) Specify the New Host Group of Notes, Supports up to 500 Characters.
* `host_group_name` - (Required) Specify the New Host Group Name, Supports up to 128 Characters.
* `instance_id` - (Required, ForceNew) Specify the New Host Group Where the Bastion Host ID of.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Group. The value formats as `<instance_id>:<host_group_id>`.
* `host_group_id` - Host Group ID.

## Import

Bastion Host Host Group can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host_group.example <instance_id>:<host_group_id>
```
