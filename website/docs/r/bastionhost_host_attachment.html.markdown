---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_attachment"
sidebar_current: "docs-alicloud-resource-bastionhost-host-attachment"
description: |-
  Provides a Alicloud Bastion Host Host Attachment resource.
---

# alicloud\_bastionhost\_host\_attachment

Provides a Bastion Host Host Attachment resource to add host into one host group.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_host_attachment" "example" {
  host_group_id = "6"
  host_id       = "15"
  instance_id   = "bastionhost-cn-tl32bh0no30"
}

```

## Argument Reference

The following arguments are supported:

* `host_group_id` - (Required, ForceNew) Specifies the added to the host group ID.
* `host_id` - (Required, ForceNew) Specified to be part of a host group of host ID.
* `instance_id` - (Required, ForceNew) The bastion host instance id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Attachment. The value formats as `<instance_id>:<host_group_id>:<host_id>`.

## Import

Bastion Host Host Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host_attachment.example <instance_id>:<host_group_id>:<host_id>
```
