---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer_udp_listeners"
sidebar_current: "docs-alicloud-datasource-ens-load-balancer-udp-listeners"
description: |-
  Provides a list of ENS Load Balancer UDP Listener owned by an Alibaba Cloud account.
---

# alicloud_ens_load_balancer_udp_listeners

This data source provides ENS Load Balancer UDP Listener available to the user.[What is Load Balancer UDP Listener](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateLoadBalancerUDPListener)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}

variable "ens_region_id" {
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "default8QXHtu" {
  network_name  = "测试用例-测试udp监听"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultN8wZgT" {
  cidr_block    = "10.0.6.0/24"
  vswitch_name  = "测试用例-测试udp监听"
  ens_region_id = alicloud_ens_network.default8QXHtu.ens_region_id
  network_id    = alicloud_ens_network.default8QXHtu.id
}

resource "alicloud_ens_load_balancer" "defaultgNxO1j" {
  load_balancer_name = "测试用例-测试udp监听"
  vswitch_id         = alicloud_ens_vswitch.defaultN8wZgT.id
  payment_type       = "PayAsYouGo"
  ens_region_id      = alicloud_ens_vswitch.defaultN8wZgT.ens_region_id
  network_id         = alicloud_ens_vswitch.defaultN8wZgT.network_id
  load_balancer_spec = "elb.s1.small"
}


resource "alicloud_ens_load_balancer_udp_listener" "default" {
  listener_port                = "53"
  health_check_interval        = "1"
  description                  = "test1"
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

data "alicloud_ens_load_balancer_udp_listeners" "default" {
  ids              = ["${alicloud_ens_load_balancer_udp_listener.default.id}"]
  load_balancer_id = alicloud_ens_load_balancer.defaultgNxO1j.id
}

output "alicloud_ens_load_balancer_udp_listener_example_id" {
  value = data.alicloud_ens_load_balancer_udp_listeners.default.listeners.0.id
}
```

## Argument Reference

The following arguments are supported:
* `load_balancer_id` - (Required, ForceNew) The ID of the load balancing instance.
* `ids` - (Optional, Computed) A list of Load Balancer UDP Listener IDs. The value is formulated as `<load_balancer_id>:<listener_port>`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Load Balancer UDP Listener IDs.
* `listeners` - A list of Load Balancer UDP Listener Entries. Each element contains the following attributes:
    * `backend_server_port` - **NOTE:** This field is only available when `enable_details` is `true`. The port used by the backend of the SLB instance.
    * `description` - Sets the description of the listener.
    * `eip_transmit` - **NOTE:** This field is only available when `enable_details` is `true`. Whether to enable EIP transparent transmission.
    * `established_timeout` - **NOTE:** This field is only available when `enable_details` is `true`. The connection timeout time.
    * `health_check_connect_port` - **NOTE:** This field is only available when `enable_details` is `true`. The port used for health check.
    * `health_check_connect_timeout` - **NOTE:** This field is only available when `enable_details` is `true`. The amount of time to wait to receive a response from the health check.
    * `health_check_exp` - **NOTE:** This field is only available when `enable_details` is `true`. The response string of the UDP listener health check, which can contain only letters and numbers.
    * `health_check_interval` - **NOTE:** This field is only available when `enable_details` is `true`. The interval between health checks.
    * `health_check_req` - **NOTE:** This field is only available when `enable_details` is `true`. The request string of the UDP listener health check, which can contain only letters and numbers.
    * `healthy_threshold` - **NOTE:** This field is only available when `enable_details` is `true`. After the number of consecutive successful health checks, the health check status of the backend server is determined from fail (the backend server is unreachable) to success (the backend server is reachable).
    * `listener_port` - The port used by the front end of the Server Load Balancer instance.
    * `load_balancer_id` - The ID of the load balancing instance.
    * `protocol` - Protocol.
    * `scheduler` - **NOTE:** This field is only available when `enable_details` is `true`. Scheduling algorithm.
    * `status` - The current status of the listener.
    * `unhealthy_threshold` - **NOTE:** This field is only available when `enable_details` is `true`. The number of consecutive health check failures that determine the health check status of the backend server from success (the backend server is reachable) to fail (the backend server is unreachable).
    * `id` - The ID of the resource supplied above.
