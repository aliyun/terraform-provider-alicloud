---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer_http_listener"
description: |-
  Provides a Alicloud ENS Load Balancer HTTP Listener resource.
---

# alicloud_ens_load_balancer_http_listener

Provides a ENS Load Balancer HTTP Listener resource.

HTTP listener for Edge Load Balancer (ELB) in Edge Node Service (ENS).

For information about ENS Load Balancer HTTP Listener and how to use it, see [CreateLoadBalancerHTTPListener](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-createloadbalancerhttplistener).

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_network" "network" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "192.168.0.0/16"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "switch" {
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = var.name
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.network.id
}

resource "alicloud_ens_load_balancer" "default" {
  load_balancer_name = var.name
  payment_type       = "PayAsYouGo"
  ens_region_id      = "cn-chenzhou-telecom_unicom_cmcc"
  load_balancer_spec = "elb.s1.small"
  vswitch_id         = alicloud_ens_vswitch.switch.id
  network_id         = alicloud_ens_network.network.id
}

resource "alicloud_ens_load_balancer_http_listener" "default" {
  load_balancer_id = alicloud_ens_load_balancer.default.id
  listener_port    = 80
  health_check     = "on"
  health_check_uri = "/"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The ID of the load balancer instance.
* `listener_port` - (Required, ForceNew, Int) The port on which the listener forwards requests. Valid values: `1` to `65535`.
* `backend_server_port` - (Optional, ForceNew, Int) The port used by the backend server. Valid values: `1` to `65535`.
* `description` - (Optional) The name of the listener. The name must be 1 to 80 characters in length and cannot start with `http://` or `https://`.
* `scheduler` - (Optional) The scheduling algorithm. Valid values: `wrr` (Weighted Round Robin, default), `wlc` (Weighted Least Connections), `rr` (Round Robin), `sch` (Source IP Hash), `qch` (QUIC Connection ID Hash), `iqch` (iQUIC CID Hash).
* `health_check` - (Optional) Specifies whether to enable health checks. Valid values: `on`, `off` (default).
* `health_check_domain` - (Optional) The domain name used for health checks.
* `health_check_uri` - (Optional) The URI used for health checks. Must start with `/` and be 1 to 80 characters in length.
* `healthy_threshold` - (Optional, Int) The number of consecutive successful health checks required before the backend server is considered healthy. Valid values: `2` to `10`. Default: `3`.
* `unhealthy_threshold` - (Optional, Int) The number of consecutive failed health checks required before the backend server is considered unhealthy. Valid values: `2` to `10`. Default: `3`.
* `health_check_timeout` - (Optional, Int) The timeout period of a health check response in seconds. Valid values: `1` to `300`. Default: `5`.
* `health_check_connect_port` - (Optional, Computed, Int) The port used for health checks. Valid values: `1` to `65535`.
* `health_check_interval` - (Optional, Int) The interval between two consecutive health checks in seconds. Valid values: `1` to `50`. Default: `2`.
* `health_check_http_code` - (Optional) The HTTP status code that indicates a successful health check. Valid values: `http_2xx` (default), `http_3xx`, `http_4xx`, `http_5xx`.
* `health_check_method` - (Optional) The HTTP method used for health checks. Valid values: `head` (default), `get`.
* `idle_timeout` - (Optional, Int) The timeout period of an idle connection in seconds. Valid values: `1` to `60`. Default: `15`.
* `request_timeout` - (Optional, Int) The request timeout period in seconds. A 504 error is returned if the backend server does not respond within this period. Valid values: `1` to `180`. Default: `60`.
* `x_forwarded_for` - (Optional) Specifies whether to use the X-Forwarded-For header to retrieve the real IP address of the client. Valid values: `on`, `off` (default).
* `listener_forward` - (Optional, ForceNew) Specifies whether to enable HTTP-to-HTTPS forwarding. Valid values: `on`, `off` (default).
* `forward_port` - (Optional, ForceNew, Int) The port to which HTTP requests are forwarded for HTTPS.
* `status` - (Optional) The status of the listener. Set it to `running` to start the listener or `stopped` to stop it. If not specified, the listener is created in the `stopped` state.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format `<load_balancer_id>:<listener_port>`.
* `status` - The status of the listener.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when creating the listener.
* `update` - (Defaults to 5 mins) Used when updating the listener.
* `delete` - (Defaults to 15 mins) Used when deleting the listener.

## Import

ENS Load Balancer HTTP Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_load_balancer_http_listener.example <load_balancer_id>:<listener_port>
```
