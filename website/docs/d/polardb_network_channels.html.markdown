---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_network_channels"
description: |-
  This data source provides the network channels of a PolarDB cluster.
---

# alicloud_polardb_network_channels

This data source provides the network channels available for a specified PolarDB cluster.

-> **NOTE:** Available since v2.2.0.

## Example Usage

```terraform
data "alicloud_polardb_network_channels" "default" {
  db_cluster_id = "pc-xxxxxxxxxxxx"
}

# Use the channel_name for PolarDB-PG logical replication
output "channel_name" {
  value = data.alicloud_polardb_network_channels.default.channel_infos.0.channel_name
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required) The ID of the PolarDB cluster.
* `channel_name` - (Optional) The name of the network channel to query.
* `vpc_id` - (Optional) The ID of the VPC to filter network channels.
* `name_regex` - (Optional) A regex string to filter results by channel name.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id` - The ID of the PolarDB cluster supplied above.
* `ids` - A list of network channel IDs. The value is formulated as `<db_cluster_id>:<channel_name>`.
* `names` - A list of network channel names.
* `channel_infos` - A list of network channels. Each element contains the following attributes:
  * `channel_name` - The name of the network channel.
  * `db_cluster_id` - The ID of the source cluster.
  * `target_db_cluster_id` - The ID of the target cluster.
  * `target_ip` - The IP address of the target instance.
  * `target_port` - The port of the target instance.
  * `target_type` - The type of the target instance.
  * `vpc_id` - The ID of the VPC to which the access endpoint resides.
  * `notes` - The note information.
  * `region_id` - The region ID.
