---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_gateway_qos"
sidebar_current: "docs-alicloud-datasource-ens-gateway-qos"
description: |-
  Provides a list of ENS Gateway Qos owned by an Alibaba Cloud account.
---

# alicloud_ens_gateway_qos

This data source provides ENS Gateway Qos available to the user.[What is Gateway Qos](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateEnsGatewayQos)

-> **NOTE:** Available since v1.287.0.

## Example Usage

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
  network_name  = "镇元-网关限速测试使用"
  cidr_block    = "10.0.0.0/10"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5giQWR" {
  cidr_block    = "10.0.8.0/24"
  vswitch_name  = "镇元-网关限速测试"
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
  instance_name              = "镇元-网关限速测试"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
}


resource "alicloud_ens_gateway_qos" "default" {
  gateway_qos_name = var.name
  bandwidth_in     = "10"
  gateway_qos_type = "Nat"
  network_id       = alicloud_ens_network.defaultC7YqlT.id
  bandwidth_out    = "20"
}

data "alicloud_ens_gateway_qos" "default" {
  ids              = ["${alicloud_ens_gateway_qos.default.id}"]
  name_regex       = alicloud_ens_gateway_qos.default.gateway_qos_name
  gateway_qos_name = var.name
  gateway_qos_type = "Nat"
  network_id       = alicloud_ens_network.defaultC7YqlT.id
}

output "alicloud_ens_gateway_qos_example_id" {
  value = data.alicloud_ens_gateway_qos.default.qos.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ens_region_id` - (Optional) The ID of the ENS node.
* `gateway_qos_id` - (Optional) The first ID of the resource
* `gateway_qos_name` - (Optional) The name of the resource
* `gateway_qos_type` - (Optional) Speed limit type.
* `instances` - (Optional) The instance under the Gateway speed limit. See [`instances`](#instances) below.
* `network_id` - (Optional) The ID of the network where the Gateway speed limit is located.
* `ids` - (Optional, Computed) A list of Gateway Qos IDs.
* `name_regex` - (Optional) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### `instances`

The instances supports the following:
* `instance_id` - (Optional) The business ID of the instance.
* `instance_type` - (Optional) The instance type.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Gateway Qos IDs.
* `names` - A list of name of Gateway Qoss.
* `qos` - A list of Gateway Qos Entries. Each element contains the following attributes:
    * `bandwidth_in` - Inbound bandwidth speed limit, unit: Mbps, range: 0-10000.
    * `bandwidth_out` - Outbound bandwidth speed limit, unit: Mbps, range: 0-10000.
    * `creation_time` - The creation time of the resource.
    * `ens_region_id` - The ID of the ENS node.
    * `gateway_qos_id` - The first ID of the resource.
    * `gateway_qos_name` - The name of the resource.
    * `gateway_qos_type` - Speed limit type.
    * `instances` - The instance under the Gateway speed limit.
      * `status` - The speed limit status of the instance-bound gateway.
    * `network_id` - The ID of the network where the Gateway speed limit is located.
    * `status` - The speed limit status of the Gateway.
    * `id` - The ID of the resource supplied above.
