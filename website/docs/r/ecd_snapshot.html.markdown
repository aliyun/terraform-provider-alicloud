---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_snapshot"
sidebar_current: "docs-alicloud-resource-ecd-snapshot"
description: |-
  Provides a Alicloud ECD Snapshot resource.
---

# alicloud\_ecd\_snapshot

Provides a ECD Snapshot resource.

For information about ECD Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/en/elastic-desktop-service/latest/createsnapshot).

-> **NOTE:** Available in v1.169.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block             = "172.16.0.0/12"
  desktop_access_type    = "Internet"
  office_site_name       = var.name
  enable_internet_access = false
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = "example_value"
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "example_value"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}


resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.0.id
  desktop_name    = var.name
}

resource "alicloud_ecd_snapshot" "default" {
  description      = var.name
  desktop_id       = alicloud_ecd_desktop.default.id
  snapshot_name    = var.name
  source_disk_type = "SYSTEM"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the Snapshot.
* `desktop_id` - (Required, ForceNew) The ID of the Desktop.
* `snapshot_name` - (Required, ForceNew) The name of the Snapshot.
* `source_disk_type` - (Required, ForceNew) The type of the disk for which to create a snapshot. Valid values: `SYSTEM`, `DATA`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot.
* `status` - The status of the snapshot.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mines) Used when create the Snapshot.
* `delete` - (Defaults to 1 mines) Used when delete the Snapshot.

## Import

ECD Snapshot can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_snapshot.example <id>
```