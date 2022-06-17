---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_tunnel"
sidebar_current: "docs-alicloud-resource-ots-tunnel"
description: |-
  Provides an OTS (Open Table Service) tunnel resource.
---

# alicloud\_ots\_tunnel

Provides an OTS tunnel resource.

For information about OTS tunnel and how to use it, see [Tunnel overview](https://www.alibabacloud.com/help/en/tablestore/latest/tunnel-service-overview).

-> **NOTE:** Available in v1.172.0+.

## Example Usage

```
variable "name" {
  default = "terraformtest"
}

resource "alicloud_ots_instance" "foo" {
  name        = var.name
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "foo" {
  instance_name = alicloud_ots_instance.foo.name
  table_name    = var.name
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  primary_key {
    name = "pk2"
    type = "String"
  }
  primary_key {
    name = "pk3"
    type = "Binary"
  }

  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}

resource "alicloud_ots_tunnel" "foo" {
  instance_name = alicloud_ots_instance.foo.name
  table_name = alicloud_ots_table.foo.table_name
  tunnel_name = var.name
  tunnel_type = "BaseAndStream"
}
```

## Argument Reference

The following arguments are supported:
* `instance_name` - (Required, ForceNew) The name of the OTS instance in which table will located.
* `table_name` - (Required, ForceNew) The name of the OTS table. If changed, a new table would be created.
* `tunnel_name` - (Required, ForceNew) The name of the OTS tunnel. If changed, a new tunnel would be created. 
* `tunnel_type` - (Required, ForceNew) The type of the OTS tunnel. Only `BaseAndStream`, `BaseData` or `Stream` is allowed.

## Attributes Reference

The following attributes are exported:

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
  * `channel_rpo` - The latest consumption time of the channel, unix time in nanosecond.
  
### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the OTS tunnel.
* `delete` - (Defaults to 10 mins) Used when delete the OTS tunnel.

## Import

OTS tunnel can be imported using id, e.g.

```
$ terraform import alicloud_ots_tunnel.foo "<instance_name>:<table_name>:<tunnel_name>"
```
