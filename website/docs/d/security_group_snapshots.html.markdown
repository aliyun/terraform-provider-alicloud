---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group_snapshots"
description: |-
    Provides a list of Security Group Snapshots available to the user.
---

# alicloud_security_group_snapshots

This data source provides a list of Security Group Snapshots in an Alibaba Cloud account according to the specified filters.

Security Group Snapshot is a read-only resource; this data source queries the `DescribeSecurityGroupSnapshots` API and returns the matching snapshots.

-> **NOTE:** Available since v1.295.0

## Example Usage

```terraform
# Filter security group snapshots and print the results into a file
data "alicloud_security_group_snapshots" "default" {
  output_file = "security_group_snapshots.json"
}

output "first_snapshot_id" {
  value = "${data.alicloud_security_group_snapshots.default.snapshots.0.id}"
}
```

```terraform
# Query by snapshot ids
data "alicloud_security_group_snapshots" "by_ids" {
  security_group_snapshot_ids = ["ss-xxxxxxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `security_group_snapshot_ids` - (Optional) A list of Security Group Snapshot IDs used to filter the results.
* `region_id` - (Optional) The region ID. Defaults to the region of the provider.
* `filter` - (Optional) The filter parameter for the list operation.
* `output_file` - (Optional) File path where to save the results (in JSON format).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source. It is the hash of the returned snapshot IDs.
* `ids` - A list of Security Group Snapshot IDs.
* `snapshots` - A list of Security Group Snapshots. Each element contains the following attributes:
    * `id` - The ID of the Security Group Snapshot. It is the value of `security_group_snapshot_id`.
    * `security_group_snapshot_id` - The ID of the Security Group Snapshot.
    * `region_id` - The region ID of the Security Group Snapshot.
