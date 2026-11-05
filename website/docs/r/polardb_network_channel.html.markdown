---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_network_channel"
description: |-
  Provides a Alicloud Polardb Network Channel resource.
---

# alicloud_polardb_network_channel

Provides a Polardb Network Channel resource.



For information about Polardb Network Channel and how to use it, see [What is Network Channel](https://next.api.alibabacloud.com/document/polardb/2017-08-01/CreateNetworkChannel).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

No example is available because resource test cases were not supplied.

## Argument Reference

The following arguments are supported:
* `channel_infos` - (Optional, ForceNew, List) Channel Network List See [`channel_infos`](#channel_infos) below.
* `channel_name` - (Required, ForceNew) The name of the network channel.
* `db_cluster_id` - (Required, ForceNew) The ID of the cluster.
* `notes` - (Optional, ForceNew) Note information.
* `owner_account` - (Optional, ForceNew) Owner Account
* `owner_id` - (Optional, ForceNew, Int) Owner ID
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `resource_owner_account` - (Optional, ForceNew) Resource Owner Account
* `resource_owner_id` - (Optional, ForceNew, Int) Resource Owner ID
* `target_db_cluster_id` - (Optional, ForceNew) The target instance name.
* `target_ip` - (Required, ForceNew) The IP address of the target instance.
* `target_port` - (Optional, ForceNew) The port of the target instance.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC to which the access endpoint resides.

  -> **NOTE:** PolarDB exposes no update API for network channels; every argument above is ForceNew, so changing any of them triggers resource replacement rather than an in-place update.


### `channel_infos`

The channel_infos supports the following:
* `channel_name` - (Optional, ForceNew) Channel Name
* `db_cluster_id` - (Optional, ForceNew) Database Cluster ID
* `notes` - (Optional, ForceNew) Notes
* `region_id` - (Optional, ForceNew) Region ID
* `target_db_cluster_id` - (Optional, ForceNew) Target Database Cluster ID
* `target_ip` - (Optional, ForceNew) Target IP
* `target_port` - (Optional, ForceNew) Target Port
* `target_type` - (Optional, ForceNew) Target Type
* `vpc_id` - (Optional, ForceNew) VPC ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<db_cluster_id>:<channel_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Channel.
* `delete` - (Defaults to 5 mins) Used when delete the Network Channel.

## Import

Polardb Network Channel can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_network_channel.example <db_cluster_id>:<channel_name>
```