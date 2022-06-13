---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_tunnels"
sidebar_current: "docs-alicloud-datasource-ots-tunnels"
description: |- 
  Provides a list of ots tunnels to the user.
---

# alicloud\_ots\_tunnels

This data source provides the ots tunnels of the current Alibaba Cloud user.

For information about OTS tunnel and how to use it, see [Tunnel overview](https://www.alibabacloud.com/help/en/tablestore/latest/tunnel-service-overview).

-> **NOTE:** Available in v1.172.0+.

## Example Usage

```
data "alicloud_ots_tunnels" "tunnels_ds" {
  instance_name = "sample-instance"
  table_name = "sample-table"
  name_regex = "sample-tunnel"
  output_file = "tunnels.txt"
}

output "first_tunnel_id" {
  value = "${data.alicloud_ots_tunnels.tunnels_ds.tunnels.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of OTS instance.
* `table_name` - (Required) The name of OTS table.
* `ids` - (Optional) A list of tunnel IDs.
* `name_regex` - (Optional) A regex string to filter results by tunnel name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of tunnel IDs.
* `names` - A list of tunnel names.
* `tunnels` - A list of tunnels. Each element contains the following attributes:
  * `id` - The resource ID. The value is `<instance_name>:<table_name>:<tunnel_name>`.
  * `instance_name` - The OTS instance name.
  * `table_name` - The table name of the OTS which could not be changed.
  * `tunnel_name` - The tunnel name of the OTS which could not be changed.
  * `tunnel_id` - The tunnel id of the OTS which could not be changed.
  * `tunnel_rpo` - The latest consumption time of the tunnel, unix time in nanosecond.
  * `tunnel_type` - The type of the OTS tunnel, valid values: `BaseAndStream`, `BaseData`, `Stream`.
  * `tunnel_stage` -  The stage of OTS tunnel, valid values: `InitBaseDataAndStreamShard`, `ProcessBaseData`, `ProcessStream`.
  * `expired` - Whether the tunnel has expired.
  * `create_time` - The creation time of the Tunnel.
  * `channels` - The channels of OTS tunnel. Each element contains the following attributes:
    * `channel_id` - The id of the channel.
    * `channel_type` - The type of the channel, valid values: `BaseData`, `Stream`.
    * `channel_status` - The status of the channel, valid values: `WAIT`, `OPEN`, `CLOSING`, `CLOSE`, `TERMINATED`.
    * `client_id` - The client id of the channel.
    * `channel_rpo` - The latest consumption time of the channel, unix time in nanosecond