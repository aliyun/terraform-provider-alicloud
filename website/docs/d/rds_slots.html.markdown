---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_slots"
description: |-
  Provides a list of Rds Replication Slots to the user.
---

# alicloud_rds_slots

This data source provides the Rds Replication Slots of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.204.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_slots" "example" {
  db_instance_id = "example_value"
}

output "first_rds_slots_name" {
  value = data.alicloud_rds_slots.example.slots.0.slot_name
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The db instance id.
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `slots` - A list of Rds Replication Slots. Each element contains the following attributes:
  * `slot_name` - The Replication Slot name.
  * `plugin` - The plugin used by Replication Slot.
  * `slot_type` - The Replication Slot type.
  * `database` - The name of the database where Replication Slot is located.
  * `temporary` - Is the Replication Slot temporary.
  * `slot_status` - The Replication Slot status.
  * `wal_delay` - The amount of logs accumulated by Replication Slot.