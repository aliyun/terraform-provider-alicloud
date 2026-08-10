---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer_udp_listener"
description: |-
  Provides a Alicloud ENS Load Balancer UDP Listener resource.
---

# alicloud_ens_load_balancer_udp_listener

Provides a ENS Load Balancer UDP Listener resource.

Load-balanced UDP listener.

For information about ENS Load Balancer UDP Listener and how to use it, see [What is Load Balancer UDP Listener](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateLoadBalancerUDPListener).

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
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "default8QXHtu" {
  network_name  = "example用例-exampleudp监听"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultN8wZgT" {
  cidr_block    = "10.0.6.0/24"
  vswitch_name  = "example用例-exampleudp监听"
  ens_region_id = alicloud_ens_network.default8QXHtu.ens_region_id
  network_id    = alicloud_ens_network.default8QXHtu.id
}

resource "alicloud_ens_load_balancer" "defaultgNxO1j" {
  load_balancer_name = "example用例-exampleudp监听"
  vswitch_id         = alicloud_ens_vswitch.defaultN8wZgT.id
  payment_type       = "PayAsYouGo"
  ens_region_id      = alicloud_ens_vswitch.defaultN8wZgT.ens_region_id
  network_id         = alicloud_ens_vswitch.defaultN8wZgT.network_id
  load_balancer_spec = "elb.s1.small"
}


resource "alicloud_ens_load_balancer_udp_listener" "default" {
  listener_port                = "53"
  health_check_interval        = "1"
  description                  = "example1"
  unhealthy_threshold          = "2"
  scheduler                    = "rr"
  health_check_connect_timeout = "1"
  load_balancer_id             = alicloud_ens_load_balancer.defaultgNxO1j.id
  backend_server_port          = "53"
  health_check_connect_port    = "53"
  health_check_req             = "hello"
  healthy_threshold            = "2"
  health_check_exp             = "rep"
  eip_transmit                 = "on"
  status                       = "Stopped"
  established_timeout          = "100"
}
```

## Argument Reference

The following arguments are supported:
* `backend_server_port` - (Optional, ForceNew, Int) The port used by the backend of the SLB instance. Valid values: `1` to 65535.
* `description` - (Optional) Sets the description of the listener.
* `eip_transmit` - (Optional) Whether to enable EIP transparent transmission. Value:
  - `on`: on.
  - `off` (default): turned off.
* `established_timeout` - (Optional, Int) The connection timeout time. Unit: seconds.
Value range: 10~900.
* `health_check_connect_port` - (Optional, Int) The port used for health check. Valid values: `1` to 65535. If this parameter is not set, the backend service port (BackendServerPort) is used.
* `health_check_connect_timeout` - (Optional, Int) The amount of time to wait to receive a response from the health check. If the backend ECS instances do not respond within the specified time, the health check fails.
  - Default value: 5 seconds.
  - Value: `1` ~ 300.
  - Unit: seconds.

-> **NOTE:** - is only valid when the HealthCheck value is on.
- If the value of HealthCheckConnectTimeout is less than the value of HealthCheckInterval, HealthCheckConnectTimeout is invalid and the timeout is the value of HealthCheckInterval.
* `health_check_exp` - (Optional) The response string of the UDP listener health check, which can contain only letters and numbers. The maximum length is 64 characters.
* `health_check_interval` - (Optional, Int) The interval between health checks. Value: `1` to `50`, in seconds.

-> **NOTE:**  is valid only when the HealthCheck value is on.

* `health_check_req` - (Optional) The request string of the UDP listener health check, which can contain only letters and numbers. The maximum length is 64 characters.
* `healthy_threshold` - (Optional, Int) After the number of consecutive successful health checks, the health check status of the backend server is determined from fail (the backend server is unreachable) to success (the backend server is reachable). Value: `2` ~ 10.

-> **NOTE:**  is valid only when the HealthCheck value is on.

* `listener_port` - (Required, ForceNew, Int) The port used by the front end of the Server Load Balancer instance.
* `load_balancer_id` - (Required, ForceNew) The ID of the load balancing instance.
* `scheduler` - (Optional) Scheduling algorithm. Value:
  - `wrr`: The higher the weight value, the higher the number of times (probability) that the backend server is polled.
  - `wlc`: In addition to polling based on the weight value set for each backend server, the actual load of the backend server (that is, the number of connections) is also considered. When the weight value is the same, the number of times (probability) that a backend server with a smaller number of current connections is polled is higher.
  - `rr`: distributes external requests to backend servers in the order they are accessed.
  - `sch`: a consistent Hash based on the source IP address. The same source IP address is dispatched to the same backend server.
  - `qch`: Consistent Hash based on QUIC Connection ID. The same QUIC Connection ID is dispatched to the same backend server.
  - `iqch`: performs consistent Hash on specific three bytes of iQUIC CID, and dispatches the same second to fourth bytes to the same backend server.
* `status` - (Optional, Computed) The current status of the listener. Value:
  - `Running`: Normal operation.
  - `Stopped`: The run is Stopped.
  - `Starting`: Starting.
  - `Configuring`: Configuring.
  - `Stopping`: The running is being stopped.
* `unhealthy_threshold` - (Optional, Int) The number of consecutive health check failures that determine the health check status of the backend server from success (the backend server is reachable) to fail (the backend server is unreachable). Value: `2` ~ 10.

-> **NOTE:**  is valid only when the HealthCheck value is on.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<load_balancer_id>:<listener_port>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer UDP Listener.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer UDP Listener.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer UDP Listener.

## Import

ENS Load Balancer UDP Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_load_balancer_udp_listener.example <load_balancer_id>:<listener_port>
```