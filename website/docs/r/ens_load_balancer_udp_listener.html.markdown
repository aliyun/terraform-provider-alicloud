---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer_udp_listener"
description: |-
  Provides a Alicloud ENS Load Balancer Udp Listener resource.
---

# alicloud_ens_load_balancer_udp_listener

Provides a ENS Load Balancer Udp Listener resource.

Load-balanced UDP listener.

For information about ENS Load Balancer Udp Listener and how to use it, see [What is Load Balancer Udp Listener](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-createloadbalancerudplistener).

-> **NOTE:** Available since v1.291.0.

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
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default" {
  cidr_block    = "10.0.6.0/24"
  vswitch_name  = var.name
  ens_region_id = alicloud_ens_network.default.ens_region_id
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_ens_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_ens_vswitch.default.id
  payment_type       = "PayAsYouGo"
  ens_region_id      = alicloud_ens_vswitch.default.ens_region_id
  network_id         = alicloud_ens_vswitch.default.network_id
  load_balancer_spec = "elb.s1.small"
}

resource "alicloud_ens_load_balancer_udp_listener" "default" {
  load_balancer_id             = alicloud_ens_load_balancer.default.id
  listener_port                = 53
  backend_server_port          = 53
  description                  = "example-udp-listener"
  scheduler                    = "wrr"
  healthy_threshold            = 2
  unhealthy_threshold          = 2
  health_check_connect_timeout = 5
  health_check_interval        = 2
  health_check_connect_port    = 53
  health_check_req             = "hello"
  health_check_exp             = "rep"
  eip_transmit                 = "off"
  established_timeout          = 900
  status                       = "Stopped"
}
```

## Argument Reference

The following arguments are supported:
* `backend_server_port` - (Optional, ForceNew, Int) The port used by the backend server of the load balancer instance. Valid values: `1` to `65535`.
* `description` - (Optional) The description of the listener. The description must be `1` to `80` characters in length and cannot start with `http://` or `https://`.
* `eip_transmit` - (Optional) Whether to enable EIP transparent transmission. Valid values:
  - `on`: enabled.
  - `off` (default): disabled.
* `established_timeout` - (Optional, Int) The timeout period of the connection. Unit: seconds. Valid values: `10` to `900`. Default value: `900`.
* `health_check_connect_port` - (Optional, Int) The port used for health checks. Valid values: `1` to `65535`. If this parameter is not set, the backend server port (`backend_server_port`) is used.
* `health_check_connect_timeout` - (Optional, Int) The amount of time to wait for a response from the health check. If the backend server does not respond within the specified time, the health check fails. Unit: seconds. Valid values: `1` to `300`. Default value: `5`. If the value of `health_check_connect_timeout` is smaller than the value of `health_check_interval`, `health_check_connect_timeout` is invalid and the timeout period is the value of `health_check_interval`.
* `health_check_exp` - (Optional) The expected response string for the UDP listener health check. It can contain only letters and digits. Maximum length: 64 characters.
* `health_check_interval` - (Optional, Int) The interval between two consecutive health checks. Unit: seconds. Valid values: `1` to `50`. Default value: `2`.
* `health_check_req` - (Optional) The request string for the UDP listener health check. It can contain only letters and digits. Maximum length: 64 characters.
* `healthy_threshold` - (Optional, Int) The number of consecutive successful health checks that must occur before a backend server is declared healthy (from `fail` to `success`). Valid values: `2` to `10`. Default value: `3`.
* `listener_port` - (Required, ForceNew, Int) The frontend port used by the load balancer instance. Valid values: `1` to `65535`. Ports `250`, `4789`, and `4790` are reserved and cannot be used.
* `load_balancer_id` - (Required, ForceNew) The ID of the load balancer instance.
* `scheduler` - (Optional) The scheduling algorithm. Valid values:
  - `wrr` (default): Backend servers with higher weights receive more requests.
  - `wlc`: Requests are distributed based on the weights and the actual load (number of connections) of each backend server.
  - `rr`: Requests are distributed to backend servers in sequence.
  - `sch`: Consistent hashing based on source IP addresses. Requests from the same source IP address are distributed to the same backend server.
  - `tch`: Consistent hashing based on the four-tuple (source IP address, destination IP address, source port, and destination port). The same flow is distributed to the same backend server.
  - `qch`: Consistent hashing based on QUIC Connection IDs. Requests with the same QUIC Connection ID are distributed to the same backend server.
  - `iqch`: Consistent hashing based on three specific bytes of the iQUIC CID. Requests whose second to fourth bytes are the same are distributed to the same backend server.
* `status` - (Optional, Computed) The status of the listener. Valid values:
  - `Running`: The listener is running.
  - `Stopped`: The listener is stopped.
  - `Starting`: The listener is starting.
  - `Configuring`: The listener is being configured.
  - `Stopping`: The listener is being stopped.
* `unhealthy_threshold` - (Optional, Int) The number of consecutive failed health checks that must occur before a backend server is declared unhealthy (from `success` to `fail`). Valid values: `2` to `10`. Default value: `3`.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<load_balancer_id>:<listener_port>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Udp Listener.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Udp Listener.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer Udp Listener.

## Import

ENS Load Balancer Udp Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_load_balancer_udp_listener.example <load_balancer_id>:<listener_port>
```
