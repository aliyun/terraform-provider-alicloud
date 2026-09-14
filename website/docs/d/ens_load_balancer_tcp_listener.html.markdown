---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer_tcp_listener"
description: |-
  Provides a data source to query the ENS Load Balancer TCP Listener.
---

# alicloud_ens_load_balancer_tcp_listener

This data source provides the ENS Load Balancer TCP Listener resource.

For information about ENS Load Balancer TCP Listener and how to use it, see [What is ENS Load Balancer TCP Listener](https://www.alibabacloud.com/help/en/ens/).

-> **NOTE:** Available since v1.288.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_ens_load_balancer_tcp_listener" "default" {
  load_balancer_id = alicloud_ens_load_balancer_tcp_listener.default.load_balancer_id
  listener_port    = 80
}

resource "alicloud_ens_load_balancer_tcp_listener" "default" {
  load_balancer_id    = "lb-xxxxx"
  listener_port       = 80
  backend_server_port = 8080
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The ID of the ENS Load Balancer instance.
* `listener_port` - (Required, ForceNew) The front-end port used by the ENS Load Balancer instance.

## Attributes Reference

The following attributes are exported in addition to the `load_balancer_id` and `listener_port` arguments above:

* `id` - The resource ID. The value formats as `<load_balancer_id>:<listener_port>`.
* `backend_server_port` - The back-end port used by the ENS Load Balancer instance.
* `description` - The description of the listener.
* `scheduler` - The scheduling algorithm.
* `persistence_timeout` - The timeout period of a persistent connection.
* `established_timeout` - The timeout period of an established TCP connection.
* `healthy_threshold` - The number of consecutive successful health checks before a backend server is declared healthy.
* `unhealthy_threshold` - The number of consecutive failed health checks before a backend server is declared unhealthy.
* `health_check_connect_timeout` - The amount of time to wait for a response from a health check.
* `health_check_interval` - The interval between health checks.
* `health_check_connect_port` - The port used for health checks.
* `health_check_domain` - The domain name used for health checks.
* `health_check_http_code` - The HTTP status codes to use for health checks.
* `health_check_type` - The health check type.
* `health_check_uri` - The URI used for health checks.
* `eip_transmit` - Whether to transmit traffic through an EIP.
* `status` - The status of the listener.
* `protocol` - The listener protocol.
