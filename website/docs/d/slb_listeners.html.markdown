---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_listeners"
sidebar_current: "docs-alicloud-datasource-slb-listeners"
description: |-
    Provides a list of server load balancer listeners to the user.
---

# alicloud\_slb_listeners

This data source provides the listeners related to a server load balancer of the current Alibaba Cloud user.

## Example Usage

```
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "tf-testAccSlbListenertcp"
}
resource "alicloud_slb_listener" "tcp" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = "22"
  frontend_port             = "22"
  protocol                  = "tcp"
  bandwidth                 = "10"
  health_check_type         = "tcp"
  persistence_timeout       = 3600
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx"
  health_check_connect_port = 20
  health_check_uri          = "/console"
  established_timeout       = 600
}

data "alicloud_slb_listeners" "sample_ds" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
}

output "first_slb_listener_protocol" {
  value = data.alicloud_slb_listeners.sample_ds.slb_listeners[0].protocol
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) ID of the SLB with listeners.
* `protocol` - (Optional) Filter listeners by the specified protocol. Valid values: `http`, `https`, `tcp` and `udp`.
* `frontend_port` - (Optional) Filter listeners by the specified frontend port.
* `description_regex` - (Optional, Available in 1.69.0+) A regex string to filter results by SLB listener description.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `slb_listeners` - A list of SLB listeners. Each element contains the following attributes:
  * `frontend_port` - Frontend port used to receive incoming traffic and distribute it to the backend servers.
  * `backend_port` - Port opened on the backend server to receive requests.
  * `protocol` - Listener protocol. Possible values: `http`, `https`, `tcp` and `udp`.
  * `status` - Listener status.
  * `security_status` - Security status. Only available when the protocol is `https`.
  * `bandwidth` - Peak bandwidth. If the value is set to -1, the listener is not limited by bandwidth.
  * `scheduler` - Algorithm used to distribute traffic. Possible values: `wrr` (weighted round robin), `wlc` (weighted least connection) and `rr` (round robin).
  * `server_group_id` - ID of the linked VServer group.
  * `master_slave_server_group_id` - ID of the active/standby server group.
  * `persistence_timeout` - Timeout value of the TCP connection in seconds. If the value is 0, the session persistence function is disabled. Only available when the protocol is `tcp`.
  * `established_timeout` - Connection timeout in seconds for the Layer 4 TCP listener. Only available when the protocol is `tcp`.
  * `sticky_session` - Indicate whether session persistence is enabled or not. If enabled, all session requests from the same client are sent to the same backend server. Possible values are `on` and `off`. Only available when the protocol is `http` or `https`.
  * `sticky_session_type` - Method used to handle the cookie. Possible values are `insert` (cookie added to the response) and `server` (cookie set by the backend server). Only available when the protocol is `http` or `https` and sticky_session is `on`.
  * `cookie_timeout` - Cookie timeout in seconds. Only available when the sticky_session_type is `insert`.
  * `cookie` - Cookie configured by the backend server. Only available when the sticky_session_type is `server`.
  * `health_check` - Indicate whether health check is enabled of not. Possible values are `on` and `off`.
  * `health_check_type` - Health check method. Possible values are `tcp` and `http`. Only available when the protocol is `tcp`.
  * `health_check_domain` - Domain name used for health check. The SLB sends HTTP head requests to the backend server, the domain is useful when the backend server verifies the host field in the requests. Only available when the protocol is `http`, `https` or `tcp` (in this case health_check_type must be `http`).
  * `health_check_uri` - URI used for health check. Only available when the protocol is `http`, `https` or `tcp` (in this case health_check_type must be `http`).
  * `health_check_connect_port` - Port used for health check.
  * `health_check_connect_timeout` - Amount of time in seconds to wait for the response for a health check.
  * `healthy_threshold` - Number of consecutive successes of health check performed on the same ECS instance (from failure to success).
  * `unhealthy_threshold` - Number of consecutive failures of health check performed on the same ECS instance (from success to failure).
  * `health_check_timeout` - Amount of time in seconds to wait for the response from a health check. If an ECS instance sends no response within the specified timeout period, the health check fails. Only available when the protocol is `http` or `https`.
  * `health_check_interval` - Time interval between two consecutive health checks.
  * `health_check_http_code` - HTTP status codes indicating that the health check is normal. It can contain several comma-separated values such as "http_2xx,http_3xx". Only available when the protocol is `http`, `https` or `tcp` (in this case health_check_type must be `http`).
  * `gzip` - Indicate whether Gzip compression is enabled or not. Possible values are `on` and `off`. Only available when the protocol is `http` or `https`.
  * `ssl_certificate_id` - ID of the server certificate. Only available when the protocol is `https`.
  * `ca_certificate_id` - ID of the CA certificate (only required when two-way authentication is used). Only available when the protocol is `https`.
  * `x_forwarded_for` - Indicate whether the HTTP header field "X-Forwarded-For" is added or not; it allows the backend server to know about the user's IP address. Possible values are `on` and `off`. Only available when the protocol is `http` or `https`.
  * `x_forwarded_for_slb_ip` - Indicate whether the HTTP header field "X-Forwarded-For_SLBIP" is added or not; it allows the backend server to know about the SLB IP address. Possible values are `on` and `off`. Only available when the protocol is `http` or `https`.
  * `x_forwarded_for_slb_id` - Indicate whether the HTTP header field "X-Forwarded-For_SLBID" is added or not; it allows the backend server to know about the SLB ID. Possible values are `on` and `off`. Only available when the protocol is `http` or `https`.
  * `x_forwarded_for_slb_proto` - Indicate whether the HTTP header field "X-Forwarded-For_proto" is added or not; it allows the backend server to know about the user's protocol. Possible values are `on` and `off`. Only available when the protocol is `http` or `https`.
  * `idle_timeout` - Timeout of http or https listener established connection idle timeout. Valid value range: [1-60] in seconds. Default to 15.
  * `request_timeout` - Timeout of http or https listener request (which does not get response from backend) timeout. Valid value range: [1-180] in seconds. Default to 60.
  * `enable_http2` -  Whether to enable https listener support http2 or not. Valid values are `on` and `off`. Default to `on`.
  * `tls_cipher_policy` - Https listener TLS cipher policy. Valid values are `tls_cipher_policy_1_0`, `tls_cipher_policy_1_1`, `tls_cipher_policy_1_2`, `tls_cipher_policy_1_2_strict`. Default to `tls_cipher_policy_1_0`.
  * `description` - The description of slb listener.
