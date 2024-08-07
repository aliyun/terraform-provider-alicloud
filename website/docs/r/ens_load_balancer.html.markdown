---
subcategory: "Ens"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer"
description: |-
  Provides a Alicloud Ens Load Balancer resource.
---

# alicloud_ens_load_balancer

Provides a Ens Load Balancer resource.

Load balancing. When you use it for the first time, please contact the product classmates to add a resource whitelist.

For information about Ens Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createloadbalancer).

-> **NOTE:** Available since v1.228.0.

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
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "network" {
  network_name  = var.name
  description   = "LoadBalancerNetworkDescription_autoexample"
  cidr_block    = "192.168.0.0/16"
  ens_region_id = "cn-hangzhou-44"
}

resource "alicloud_ens_vswitch" "switch" {
  description   = "LoadBalancerVSwitchDescription_autoexample"
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = format("%s1", var.name)
  ens_region_id = "cn-hangzhou-44"
  network_id    = alicloud_ens_network.network.id
}

resource "alicloud_ens_instance" "defaultfGH5i7" {
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  vswitch_id                 = alicloud_ens_vswitch.switch.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = format("%s2", var.name)
  internet_max_bandwidth_out = "0"
  auto_use_coupon            = "true"
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}


resource "alicloud_ens_load_balancer" "default" {
  load_balancer_name = var.name
  payment_type       = "PayAsYouGo"
  ens_region_id      = "cn-hangzhou-44"
  load_balancer_spec = "elb.s1.small"
  vswitch_id         = alicloud_ens_vswitch.switch.id
  network_id         = alicloud_ens_network.network.id
}
```

## Argument Reference

The following arguments are supported:
* `backend_servers` - (Optional) The list of backend servers. See [`backend_servers`](#backend_servers) below.
* `ens_region_id` - (Required, ForceNew) The ID of the ENS node.
* `load_balancer_name` - (Optional) Name of the Server Load Balancer instance

  Rules:

  The length is 1~80 English or Chinese characters. When this parameter is not specified, the system randomly assigns an instance name

  Cannot start with http:// and https.
* `load_balancer_spec` - (Required, ForceNew) Specifications of the Server Load Balancer instance

  Example value: elb.s2.medium

  Optional values: elb.s1.small,elb.s3.medium,elb.s2.small,elb.s2.medium,elb.s3.small
* `network_id` - (Required, ForceNew) The network ID of the created edge load balancing (ELB) instance.

  Example value: n-5sax03dh2eyagujgsn7z9 * * * *
* `payment_type` - (Required, ForceNew) Server Load Balancer Instance Payment Type

  Value:PayAsYouGo
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch to which the VPC instance belongs

  Example value: vsw-5s78haoys9oylle6ln71m * * * *

### `backend_servers`

The backend_servers supports the following:
* `ip` - (Optional) IP address of the backend server  Example value: 192.168.0.5.
* `port` - (Optional, Computed) Port used by the backend server.
* `server_id` - (Required) Backend server instance ID  Example value: i-5vb5h5njxiuhn48a * * * *.
* `type` - (Optional) Backend server type  Example value: ens.
* `weight` - (Optional) Weight of the backend server  Example value: 100.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation Time (UTC) of the load balancing instance.
* `status` - The status of the SLB instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

Ens Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_load_balancer.example <id>
```