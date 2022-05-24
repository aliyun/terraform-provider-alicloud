---
subcategory: "Ddos Basic"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddos_basic_defense_threshold"
sidebar_current: "docs-alicloud-resource-ddos-basic-defense-threshold"
description: |-
  Provides a Alicloud Ddos Basic defense threshold resource.
---

# alicloud\_ddos\_basic\_antiddos

Provides a Ddos Basic defense threshold resource.

For information about Ddos Basic Antiddos and how to use it, see [What is Defense Threshold](https://www.alibabacloud.com/help/en/ddos-protection/latest/modifydefensethreshold).

-> **NOTE:** Available in v1.168.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ddos_basic_defense_threshold" "example" {
  instance_id   = alicloud_eip_address.default.id
  ddos_type     = "defense"
  instance_type = "eip"
  bps           = 390
  pps           = 90000
}

resource "alicloud_eip_address" "default" {
  address_name         = var.name
  isp                  = "BGP"
  internet_charge_type = "PayByBandwidth"
  payment_type         = "PayAsYouGo"
}
```

## Argument Reference

The following arguments are supported:
* `instance_type` - (Required, ForceNew) The instance type of the public IP address asset. Value: `ecs`,`slb`,`eip`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `ddos_type` - (Required, ForceNew) The type of the threshold to query. Valid values: `defense`,`blackhole`.
  -`defense` - scrubbing threshold.
  -`blackhole` - DDoS mitigation threshold.
* `bps` - (Optional) Specifies the traffic scrubbing threshold. Unit: Mbit/s. The traffic scrubbing threshold cannot exceed the peak inbound or outbound Internet traffic, whichever is larger, of the asset.
* `pps` - (Optional) The current message number cleaning threshold. Unit: pps.
* `is_auto` - (Optional, Computed) Whether it is the system default threshold. Value:
  - `true`: indicates yes, that is, the DDoS protection service dynamically adjusts the cleaning threshold according to the traffic load of the cloud server.
  - `false`: indicates no, that is, you manually set the cleaning threshold.
* `internet_ip` - (Optional) The Internet IP address.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Antiddos. The value formats as `<instance_id>:<instance_type>:<ddos_type>`.
* `max_bps` - The maximum traffic scrubbing threshold. Unit: Mbit/s.
* `max_pps` - The maximum packet scrubbing threshold. Unit: pps.

## Import

Ddos Basic Antiddos can be imported using the id, e.g.

```
$ terraform import alicloud_ddos_basic_antiddos.example <instance_id>:<instance_type>:<ddos_type>
```