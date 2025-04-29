---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_tunnel"
sidebar_current: "docs-alicloud-resource-ots-tunnel"
description: |-
  Provides an OTS (Open Table Service) tunnel resource.
---

# alicloud_ots_tunnel

Provides an OTS tunnel resource.

For information about OTS tunnel and how to use it, see [Tunnel overview](https://www.alibabacloud.com/help/en/tablestore/latest/tunnel-service-overview).

-> **NOTE:** Available since v1.172.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ots_tunnel&exampleId=4f2e8e09-a923-817c-473e-470b380fa1c464516321&activeTab=example&spm=docs.r.ots_tunnel.0.4f2e8e09a9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ots_instance" "default" {
  name        = "${var.name}-${random_integer.default.result}"
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_ots_table" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = "tf_example"
  time_to_live  = -1
  max_version   = 1
  enable_sse    = true
  sse_key_type  = "SSE_KMS_SERVICE"
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
}

resource "alicloud_ots_tunnel" "default" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = alicloud_ots_table.default.table_name
  tunnel_name   = "tf_example"
  tunnel_type   = "BaseAndStream"
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
* `tunnel_id` - The tunnel id of the OTS which could not be changed.
* `tunnel_rpo` - The latest consumption time of the tunnel, unix time in nanosecond.
* `tunnel_stage` -  The stage of OTS tunnel, valid values: `InitBaseDataAndStreamShard`, `ProcessBaseData`, `ProcessStream`.
* `expired` - Whether the tunnel has expired.
* `create_time` - The creation time of the Tunnel.
* `channels` - The channels of OTS tunnel. Each element contains the following attributes:
  * `channel_id` - The id of the channel.
  * `channel_type` - The type of the channel, valid values: `BaseData`, `Stream`.
  * `channel_status` - The status of the channel, valid values: `WAIT`, `OPEN`, `CLOSING`, `CLOSE`, `TERMINATED`.
  * `client_id` - The client id of the channel.
  * `channel_rpo` - The latest consumption time of the channel, unix time in nanosecond.
  
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the OTS tunnel.
* `delete` - (Defaults to 10 mins) Used when delete the OTS tunnel.

## Import

OTS tunnel can be imported using id, e.g.

```shell
$ terraform import alicloud_ots_tunnel.foo <instance_name>:<table_name>:<tunnel_name>
```
