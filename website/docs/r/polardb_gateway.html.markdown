---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway"
sidebar_current: "docs-alicloud-resource-polardb-gateway"
description: |-
  Provides a PolarDB Gateway resource.
---

# alicloud_polardb_gateway

Provides a PolarDB Gateway resource.

-> **NOTE:** Available since v1.290.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_gateway&exampleId=a987fdef-2f38-783f-eef3-12a032f8e5c74ff6697d&activeTab=example&spm=docs.r.polardb_gateway.0.a987fdef2f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_gateway" "default" {
  zone_id           = "cn-beijing-l"
  db_cluster_class  = "polar.app.g2.small"
  pay_type          = "Postpaid"
  vpc_id            = "vpc-xxx"
  vswitch_id        = "vsw-xxx"
  security_group_id = "sg-xxx"
  db_type           = "PostgreSQL"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_gateway&spm=docs.r.polardb_gateway.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional, ForceNew, Available since v1.290.0) The zone ID of the PolarDB gateway.
* `db_cluster_class` - (Optional, ForceNew, Available since v1.290.0) The specifications of the PolarDB gateway.
* `pay_type` - (Required, ForceNew, Available since v1.290.0) The billing method. Valid values: `Postpaid`, `Prepaid`.
* `auto_renew` - (Optional, ForceNew, Available since v1.290.0) Whether to enable automatic renewal for a subscription gateway. Default value: `false`.
* `period` - (Optional, ForceNew, Available since v1.290.0) The unit of the subscription duration. Valid values: `Month`, `Year`. This argument is required when `pay_type` is `Prepaid`.
* `used_time` - (Optional, ForceNew, Available since v1.290.0) The subscription duration. Valid values are `1` to `9` when `period` is `Month`, and `1` to `3` when `period` is `Year`. This argument is required when `pay_type` is `Prepaid`.
* `vpc_id` - (Required, ForceNew, Available since v1.290.0) The ID of the VPC.
* `vswitch_id` - (Required, ForceNew, Available since v1.290.0) The ID of the vSwitch.
* `security_group_id` - (Optional, ForceNew, Available since v1.290.0) The ID of the security group.
* `db_type` - (Optional, ForceNew, Available since v1.290.0) The database engine. Valid values: `MySQL`, `PostgreSQL`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the PolarDB gateway.
* `region_id` - (Available since v1.290.0) The region ID of the PolarDB gateway.
* `status` - (Available since v1.290.0) The status of the PolarDB gateway.
* `description` - (Available since v1.290.0) The description of the PolarDB gateway.
* `create_time` - (Available since v1.290.0) The time when the PolarDB gateway was created.
* `modify_time` - (Available since v1.290.0) The time when the PolarDB gateway was last modified.
* `expire_time` - (Available since v1.290.0) The expiration time of the PolarDB gateway.
* `expired` - (Available since v1.290.0) Indicates whether the PolarDB gateway has expired.
* `latest_version` - (Available since v1.290.0) The latest available gateway version.
* `current_version` - (Available since v1.290.0) The current gateway version.
* `running_version` - (Available since v1.290.0) The running gateway version.
* `endpoints` - (Available since v1.290.0) The endpoints of the PolarDB gateway. Each item contains the following attributes:
  * `address` - The endpoint address.
  * `endpoint_id` - The endpoint ID.
  * `gateway_id` - The gateway ID.
  * `port` - The endpoint port.
  * `tunnel_id` - The tunnel ID.
  * `vpc_id` - The VPC ID.
  * `network_type` - The network type. Valid values: `Private`, `Public`.
* `security_ip_arrays` - (Available since v1.290.0) The IP whitelist groups of the PolarDB gateway. Each item contains the following attributes:
  * `name` - The whitelist group name.
  * `ip_list` - The IP addresses in the whitelist group.

## Timeouts

The `timeouts` block allows you to specify timeouts for certain actions:

* `create` - (Defaults to 50 mins) Used when creating the PolarDB Gateway.
* `delete` - (Defaults to 10 mins) Used when deleting the PolarDB Gateway.

## Import

PolarDB Gateway can be imported using the ID, e.g.

```shell
$ terraform import alicloud_polardb_gateway.example gw-abc12345678
```
