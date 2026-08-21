---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_members"
description: |-
  This data source provides the members of the current Realtime Compute workspace namespace.
---

# alicloud_realtime_compute_members

This data source provides the members of the current Realtime Compute workspace namespace.

-> **NOTE:** Available since v1.265.0.

## Example Usage

```terraform
data "alicloud_realtime_compute_members" "default" {
  resource_id = "xxxxxxx"
  namespace   = "name-default"
}

output "first_member_role" {
  value = data.alicloud_realtime_compute_members.default.members.0.role
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The resource ID of the VVP workspace.
* `namespace` - (Required) The name of the namespace to query members for.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the `arguments` listed above:

* `id` - The data source ID, formatted as `<resource_id>:<namespace>`.
* `members` - A list of members. Each element contains the following attributes:
  * `member` - The RAM user ID of the member.
  * `role` - The role of the member.
* `ids` - A list of member IDs.
* `total_size` - The total number of members in the namespace.
