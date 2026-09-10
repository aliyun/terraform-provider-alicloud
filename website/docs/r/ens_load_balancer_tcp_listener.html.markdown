---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer_tcp_listener"
description: |-
  Provides a Alicloud ENS Load Balancer TCP Listener resource.
---

# alicloud_ens_load_balancer_tcp_listener

Provides a ENS Load Balancer TCP Listener resource.

For information about ENS Load Balancer TCP Listener and how to use it, see [What is ENS Load Balancer TCP Listener](https://www.alibabacloud.com/help/en/ens/).

-> **NOTE:** Available since v1.288.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_network" "network" {
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "switch" {
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.network.id
}

resource "alicloud_ens_load_balancer" "default" {
  payment_type       = "PayAsYouGo"
  ens_region_id      = "cn-chenzhou-telecom_unicom_cmcc"
  load_balancer_spec = "elb.s1.small"
  vswitch_id         = alicloud_ens_vswitch.switch.id
  network_id         = alicloud_ens_network.network.id
}

resource "alicloud_ens_load_balancer_tcp_listener" "default" {
  load_balancer_id             = alicloud_ens_load_balancer.default.id
  listener_port                = 80
  backend_server_port          = 8080
  description                  = "tcp-listener"
  scheduler                    = "wrr"
  established_timeout          = 300
  healthy_threshold            = 3
  unhealthy_threshold          = 3
  health_check_connect_timeout = 5
  health_check_interval        = 10
  health_check_connect_port    = 8080
  health_check_type            = "tcp"
  eip_transmit                 = "off"
  status                       = "Running"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The ID of the ENS Load Balancer instance.
* `listener_port` - (Required, ForceNew) The front-end port used by the ENS Load Balancer instance. Valid values: `1` to `65535`.
* `backend_server_port` - (Optional, ForceNew) The back-end port used by the ENS Load Balancer instance. Valid values: `1` to `65535`.
* `description` - (Optional) The description of the listener. The length is limited to `1` to `80` characters. It cannot start with `http://` or `https://`.
* `scheduler` - (Optional) The scheduling algorithm. Valid values: `wrr` (default), `wlc`, `rr`, `sch`, `qch`, `iqch`.
* `persistence_timeout` - (Optional) The timeout period of a persistent connection. Unit: seconds.
* `established_timeout` - (Optional) The timeout period of an established TCP connection. Unit: seconds.
* `healthy_threshold` - (Optional) The number of consecutive successful health checks before a backend server is declared healthy. Valid values: `1` to `10`.
* `unhealthy_threshold` - (Optional) The number of consecutive failed health checks before a backend server is declared unhealthy. Valid values: `1` to `10`.
* `health_check_connect_timeout` - (Optional) The amount of time to wait for a response from a health check. Unit: seconds. Valid values: `1` to `300`.
* `health_check_interval` - (Optional) The interval between health checks. Unit: seconds. Valid values: `1` to `50`.
* `health_check_connect_port` - (Optional) The port used for health checks. Valid values: `1` to `65535`.
* `health_check_domain` - (Optional) The domain name used for health checks.
* `health_check_http_code` - (Optional) The HTTP status codes to use for health checks. Valid values: `http_2xx`, `http_3xx`, `http_4xx`, `http_5xx`, and combinations joined by commas.
* `health_check_type` - (Optional) The health check type. Valid values: `tcp`, `http`.
* `health_check_uri` - (Optional) The URI used for health checks.
* `eip_transmit` - (Optional) Whether to transmit traffic through an EIP. Valid values: `on`, `off`.
* `status` - (Optional) The status of the listener. Valid values: `Running`, `Stopped`. Setting `Running` starts the listener, setting `Stopped` stops it.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the tcp listener. The value formats as `<load_balancer_id>:<listener_port>`.
* `protocol` - The listener protocol. The value is `tcp`.

## Import

ENS Load Balancer TCP Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_load_balancer_tcp_listener.example <load_balancer_id>:<listener_port>
```
