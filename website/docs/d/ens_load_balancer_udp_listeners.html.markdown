---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer_udp_listeners"
sidebar_current: "docs-alicloud-datasource-ens-load-balancer-udp-listeners"
description: |-
  Provides a list of ENS Load Balancer Udp Listener owned by an Alibaba Cloud account.
---

# alicloud_ens_load_balancer_udp_listeners

This data source provides ENS Load Balancer Udp Listener available to the user. [What is Load Balancer Udp Listener](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-createloadbalancerudplistener)

-> **NOTE:** Available since v1.291.0.

## Example Usage

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
  load_balancer_id    = alicloud_ens_load_balancer.default.id
  listener_port       = 53
  backend_server_port = 53
  description         = "example-udp-listener"
  status              = "Stopped"
}

data "alicloud_ens_load_balancer_udp_listeners" "default" {
  load_balancer_id = alicloud_ens_load_balancer.default.id
  ids              = [alicloud_ens_load_balancer_udp_listener.default.id]
}

output "udp_listener_id" {
  value = data.alicloud_ens_load_balancer_udp_listeners.default.listeners.0.id
}
```

## Argument Reference

The following arguments are supported:
* `load_balancer_id` - (Required) The ID of the load balancer instance.
* `ids` - (Optional, Computed) A list of Load Balancer Udp Listener IDs. The value is formulated as `<load_balancer_id>:<listener_port>`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Load Balancer Udp Listener IDs.
* `listeners` - A list of Load Balancer Udp Listener Entries. Each element contains the following attributes:
  * `id` - The ID of the listener. The value is formulated as `<load_balancer_id>:<listener_port>`.
  * `listener_port` - The frontend port used by the load balancer instance.
  * `load_balancer_id` - The ID of the load balancer instance.
  * `protocol` - The protocol of the listener. The value is `udp`.
  * `description` - The description of the listener.
  * `status` - The status of the listener.
  * `backend_server_port` - **NOTE:** This field is only available when `enable_details` is `true`. The port used by the backend server of the load balancer instance.
  * `eip_transmit` - **NOTE:** This field is only available when `enable_details` is `true`. Whether EIP transparent transmission is enabled.
  * `established_timeout` - **NOTE:** This field is only available when `enable_details` is `true`. The timeout period of the connection. Unit: seconds.
  * `health_check_connect_port` - **NOTE:** This field is only available when `enable_details` is `true`. The port used for health checks.
  * `health_check_connect_timeout` - **NOTE:** This field is only available when `enable_details` is `true`. The amount of time to wait for a response from the health check. Unit: seconds.
  * `health_check_exp` - **NOTE:** This field is only available when `enable_details` is `true`. The expected response string for the UDP listener health check.
  * `health_check_interval` - **NOTE:** This field is only available when `enable_details` is `true`. The interval between two consecutive health checks. Unit: seconds.
  * `health_check_req` - **NOTE:** This field is only available when `enable_details` is `true`. The request string for the UDP listener health check.
  * `healthy_threshold` - **NOTE:** This field is only available when `enable_details` is `true`. The number of consecutive successful health checks that must occur before a backend server is declared healthy.
  * `scheduler` - **NOTE:** This field is only available when `enable_details` is `true`. The scheduling algorithm.
  * `unhealthy_threshold` - **NOTE:** This field is only available when `enable_details` is `true`. The number of consecutive failed health checks that must occur before a backend server is declared unhealthy.
