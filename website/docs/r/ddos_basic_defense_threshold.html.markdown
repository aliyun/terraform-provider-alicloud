---
subcategory: "Ddos Basic"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddos_basic_defense_threshold"
sidebar_current: "docs-alicloud-resource-ddos-basic-defense-threshold"
description: |-
  Provides a Alicloud Ddos Basic defense threshold resource.
---

# alicloud_ddos_basic_defense_threshold

Provides a Ddos Basic defense threshold resource.

For information about Ddos Basic Antiddos and how to use it, see [What is Defense Threshold](https://www.alibabacloud.com/help/en/ddos-protection/latest/modifydefensethreshold).

-> **NOTE:** Available since v1.168.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ddos_basic_defense_threshold&exampleId=5cd3a082-346b-a1e9-fe9b-9bec014af3baeee821a1&activeTab=example&spm=docs.r.ddos_basic_defense_threshold.0.5cd3a08234&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_eip_address" "default" {
  address_name         = var.name
  isp                  = "BGP"
  internet_charge_type = "PayByBandwidth"
  payment_type         = "PayAsYouGo"
}

resource "alicloud_ddos_basic_defense_threshold" "default" {
  instance_id   = alicloud_eip_address.default.id
  ddos_type     = "defense"
  instance_type = "eip"
  bps           = 390
  pps           = 90000
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ddos_basic_defense_threshold&spm=docs.r.ddos_basic_defense_threshold.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `instance_type` - (Required, ForceNew) The instance type of the public IP address asset. Value: `ecs`,`slb`,`eip`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `ddos_type` - (Required, ForceNew) The type of the threshold to query. Valid values: `defense`,`blackhole`.
  -`defense` - scrubbing threshold.
  -`blackhole` - DDoS mitigation threshold.
* `bps` - (Optional) Specifies the traffic scrubbing threshold. Unit: Mbit/s. The traffic scrubbing threshold cannot exceed the peak inbound or outbound Internet traffic, whichever is larger, of the asset.
* `pps` - (Optional) The current message number cleaning threshold. Unit: pps.
* `is_auto` - (Optional) Whether it is the system default threshold. Value:
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

```shell
$ terraform import alicloud_ddos_basic_antiddos.example <instance_id>:<instance_type>:<ddos_type>
```