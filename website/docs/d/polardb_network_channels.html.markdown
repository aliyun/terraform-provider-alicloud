---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_network_channels"
sidebar_current: "docs-alicloud-datasource-polardb-network-channels"
description: |-
  Queries PolarDB network channels for a source cluster.
---

# alicloud_polardb_network_channels

Queries existing PolarDB network channels using `DescribeNetworkChannel`. The returned channel names can be used when configuring PostgreSQL logical replication.

-> **NOTE:** Available since v2.4.0.

-> **NOTE:** This data source only queries existing channels. It does not create network channels, enable logical replication, or execute `CREATE SUBSCRIPTION`.

## Example Usage

```terraform
variable "source_db_cluster_id" {
  type        = string
  description = "ID of the existing source PolarDB cluster."
}

variable "channel_name" {
  type        = string
  description = "Name of the existing network channel to retrieve."
}

data "alicloud_polardb_network_channels" "replication" {
  db_cluster_id = var.source_db_cluster_id
  channel_name  = var.channel_name
}

output "channel_names" {
  value = data.alicloud_polardb_network_channels.replication.names
}
```

For a source cluster created in the same Terraform configuration, set `db_cluster_id = alicloud_polardb_cluster.source.id` to establish the dependency. An empty result returns empty `names` and `channels` lists; check that a channel exists before indexing either list.

The caller needs `polardb:DescribeNetworkChannel` permission for the source cluster. Database credentials are not required by this data source and should not be included in its configuration.

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, Available since v2.4.0) ID of the source PolarDB cluster.
* `channel_name` - (Optional, Available since v2.4.0) Exact channel name. Sent to the API and also applied to returned channel names.
* `resource_group_id` - (Optional, Available since v2.4.0) Resource group ID sent to the API.
* `vpc_id` - (Optional, Available since v2.4.0) VPC ID sent to the API. This is a server-side filter; the provider does not discard channels based on an absent VPC field in the response.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Available since v2.4.0) Stable identifier for the query.
* `names` - (Available since v2.4.0) Channel names in the same order as `channels`.
* `channels` - (Available since v2.4.0) Channels ordered by name. Each element contains:
  * `channel_name` - (Available since v2.4.0) Channel name for the replication connection.
  * `db_cluster_id` - (Available since v2.4.0) Source cluster ID. Defaults to the queried cluster if omitted in the response.
  * `eni_id` - (Available since v2.4.0) ID of the elastic network interface used by the channel's reverse-access link, if returned. This is read-only metadata; this data source does not manage the network interface. Empty if the API does not return it.
  * `region_id` - (Available since v2.4.0) Region ID. Defaults to the provider region if omitted in the response.
  * `notes` - (Available since v2.4.0) Channel notes.
  * `target_db_cluster_id` - (Available since v2.4.0) Target cluster ID, if returned.
  * `target_ip` - (Available since v2.4.0) Target IP address, if returned.
  * `target_port` - (Available since v2.4.0) Target port as a string, if returned.
  * `target_type` - (Available since v2.4.0) Target database type, if returned.
  * `vpc_id` - (Available since v2.4.0) VPC ID, if returned.

## Timeouts

The `timeouts` block allows you to specify the timeout for certain actions:

* `read` - (Defaults to 5 minutes) Used when querying network channels.
