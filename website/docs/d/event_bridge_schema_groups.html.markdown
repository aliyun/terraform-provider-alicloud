---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_schema_groups"
sidebar_current: "docs-alicloud-datasource-event-bridge-schema-groups"
description: |-
  Provides a list of Event Bridge Schema Groups to the user.
---

# alicloud\_event\_bridge\_schema\_groups

This data source provides the Event Bridge Schema Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.126.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_event_bridge_schema_groups" "example" {
  ids   = ["the_group_id"]
}

output "first_event_bridge_schema_group_id" {
  value = data.alicloud_event_bridge_schema_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Schema Group IDs. The element value is same as `group_id`.
* `description_regex` - (Optional) A regex string to filter results by resource description.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `prefix` - (Optional, ForceNew) The prefix.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `descriptions` - A list of Event Bridge Schema Group description.
* `groups` - A list of Event Bridge Schema Groups. Each element contains the following attributes:
	* `description` - the description of group.
	* `format` - The format of schemas.
	* `group_id` - The first ID of the resource.
	* `id` - The ID of the Schema Group.