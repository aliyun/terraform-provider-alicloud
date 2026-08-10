---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_gateway_qos"
description: |-
  Provides a Alicloud ENS Gateway Qos resource.
---

# alicloud_ens_gateway_qos

Provides a ENS Gateway Qos resource.

Gateway speed limit.

For information about ENS Gateway Qos and how to use it, see [What is Gateway Qos](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateEnsGatewayQos).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "ens_region_id" {
  default = "cn-hangzhou-63"
}

resource "alicloud_ens_network" "defaultC7YqlT" {
  network_name  = "镇元-网关限速example使用"
  cidr_block    = "10.0.0.0/10"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5giQWR" {
  cidr_block    = "10.0.8.0/24"
  vswitch_name  = "镇元-网关限速example"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.defaultC7YqlT.id
}

resource "alicloud_ens_instance" "defaultvuczsY" {
  amount      = "1"
  period_unit = "Month"
  auto_renew  = false
  system_disk {
    size = "20"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.small"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  vswitch_id                 = alicloud_ens_vswitch.default5giQWR.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "镇元-网关限速example"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
}


resource "alicloud_ens_gateway_qos" "default" {
  gateway_qos_name = "example"
  bandwidth_in     = "10"
  gateway_qos_type = "Nat"
  network_id       = alicloud_ens_network.defaultC7YqlT.id
  bandwidth_out    = "20"
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth_in` - (Optional, Int) Inbound bandwidth speed limit, unit: Mbps, range: 0-10000
* `bandwidth_out` - (Optional, Int) Outbound bandwidth speed limit, unit: Mbps, range: 0-10000.
* `gateway_qos_name` - (Optional) The name of the resource
* `gateway_qos_type` - (Required, ForceNew) Speed limit type.
* `instances` - (Optional, List) The instance under the Gateway speed limit. See [`instances`](#instances) below.
* `network_id` - (Required, ForceNew) The ID of the network where the Gateway speed limit is located.

### `instances`

The instances supports the following:
* `instance_id` - (Required) The business ID of the instance.
* `instance_type` - (Required) The instance type.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `creation_time` - The creation time of the resource.
* `ens_region_id` - The ID of the ENS node.
* `instances` - The instance under the Gateway speed limit.
  * `status` - The speed limit status of the instance-bound gateway.
* `status` - The speed limit status of the Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Gateway Qos.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway Qos.
* `update` - (Defaults to 5 mins) Used when update the Gateway Qos.

## Import

ENS Gateway Qos can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_gateway_qos.example <gateway_qos_id>
```