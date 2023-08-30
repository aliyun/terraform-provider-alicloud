---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_disk_replica_group"
sidebar_current: "docs-alicloud-resource-ebs-disk-replica-group"
description: |-
  Provides a Alicloud EBS Disk Replica Group resource.
---

# alicloud_ebs_disk_replica_group

Provides a EBS Disk Replica Group resource.

For information about EBS Disk Replica Group and how to use it, see [What is Disk Replica Group](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/creatediskreplicagroup).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_ebs_regions" "default" {
  region_id = data.alicloud_regions.default.regions.0.id
}

resource "alicloud_ebs_disk_replica_group" "default" {
  source_region_id      = data.alicloud_regions.default.regions.0.id
  source_zone_id        = data.alicloud_ebs_regions.default.regions[0].zones[0].zone_id
  destination_region_id = data.alicloud_regions.default.regions.0.id
  destination_zone_id   = data.alicloud_ebs_regions.default.regions[0].zones[1].zone_id
  group_name            = var.name
  description           = var.name
  rpo                   = 900
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the consistent replication group.
* `destination_region_id` - (Required, ForceNew) The ID of the region to which the disaster recovery site belongs.
* `destination_zone_id` - (Required, ForceNew) The ID of the zone to which the disaster recovery site belongs.
* `group_name` - (Optional) Consistent replication group name.
* `rpo` - (Optional, Computed, ForceNew) The recovery point objective (RPO) of the replication pair-consistent group. Unit: seconds.
* `source_region_id` - (Required, ForceNew) The ID of the region to which the production site belongs.
* `source_zone_id` - (Required, ForceNew) The ID of the zone to which the production site belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Disk Replica Group.
* `status` - The status of the consistent replication group. 


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Disk Replica Group.
* `update` - (Defaults to 1 mins) Used when update the Disk Replica Group.
* `delete` - (Defaults to 1 mins) Used when delete the Disk Replica Group.


## Import

EBS Disk Replica Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_disk_replica_group.example <id>
```