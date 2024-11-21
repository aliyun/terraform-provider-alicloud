---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_listener"
sidebar_current: "docs-alicloud-resource-slb-listener"
description: |-
  Provides a Alicloud Classic Load Balancer (SLB) Listener resource.
---

# alicloud_slb_listener

Provides a Classic Load Balancer (SLB) Load Balancer Listener resource.

For information about Classic Load Balancer (SLB) and how to use it, see [What is Classic Load Balancer](https://www.alibabacloud.com/help/doc-detail/27539.htm).

For information about listener and how to use it, please see the following:

* [Configure a HTTP Classic Load Balancer (SLB) Listener](https://www.alibabacloud.com/help/doc-detail/27592.htm).
* [Configure a HTTPS Classic Load Balancer (SLB) Listener](https://www.alibabacloud.com/help/doc-detail/27593.htm).
* [Configure a TCP Classic Load Balancer (SLB) Listener](https://www.alibabacloud.com/help/doc-detail/27594.htm).
* [Configure a UDP Classic Load Balancer (SLB) Listener](https://www.alibabacloud.com/help/doc-detail/27595.htm).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_listener&exampleId=d81ccb9a-e8c9-58fc-3b3c-9413b8fc14c23294f58f&activeTab=example&spm=docs.r.slb_listener.0.d81ccb9ae8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_slb_load_balancer" "listener" {
  load_balancer_name   = "${var.name}-${random_integer.default.result}"
  internet_charge_type = "PayByTraffic"
  address_type         = "internet"
  instance_charge_type = "PayByCLCU"
}

resource "alicloud_slb_listener" "listener" {
  load_balancer_id          = alicloud_slb_load_balancer.listener.id
  backend_port              = 80
  frontend_port             = 80
  protocol                  = "http"
  bandwidth                 = 10
  sticky_session            = "on"
  sticky_session_type       = "insert"
  cookie_timeout            = 86400
  cookie                    = "tfslblistenercookie"
  health_check              = "on"
  health_check_domain       = "ali.com"
  health_check_uri          = "/cons"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx,http_3xx"
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status      = "on"
  acl_type        = "white"
  acl_id          = alicloud_slb_acl.listener.id
  request_timeout = 80
  idle_timeout    = 30
}

resource "alicloud_slb_acl" "listener" {
  name       = "${var.name}-${random_integer.default.result}"
  ip_version = "ipv4"
}

resource "alicloud_slb_acl_entry_attachment" "first" {
  acl_id  = alicloud_slb_acl.listener.id
  entry   = "10.10.10.0/24"
  comment = "first"
}

resource "alicloud_slb_acl_entry_attachment" "second" {
  acl_id  = alicloud_slb_acl.listener.id
  entry   = "168.10.10.0/24"
  comment = "second"
}
```

## Argument Reference

The following arguments are supported:

* `HTTP Listener` - See [`HTTP Listener`](#HTTP Listener) below.
* `HTTPS Listener` -  See [`HTTPS Listener`](#HTTPS Listener) below.
* `TCP Listener` -  See [`TCP Listener`](#TCP Listener) below.
* `UDP Listener` -  See [`UDP Listener`](#UDP Listener) below.

### HTTP Listener

The HTTP Listener supports the following:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new listener.
* `protocol` - (Required, ForceNew) The protocol to listen on. Valid values: `http`.
* `frontend_port` - (Required, Int, ForceNew) The frontend port that is used by the CLB instance. Valid values: `1` to `65535`.
* `backend_port` - (Optional, Int, ForceNew) The backend port that is used by the CLB instance. Valid values: `1` to `65535`. **NOTE:** If `server_group_id` is not set, `backend_port` is required.
* `bandwidth` - (Optional, Int) The maximum bandwidth of the listener. Unit: Mbit/s. Valid values:
  - `-1`: If you set `bandwidth` to `-1`, the bandwidth of the listener is unlimited.
  - `1` to `1000`: The sum of the maximum bandwidth that you specify for all listeners of the CLB instance cannot exceed the maximum bandwidth of the CLB instance.
-> **NOTE:** Currently, this `bandwidth` is available on `Domestic Site Account`.
* `scheduler` - (Optional) The scheduling algorithm. Default value: `wrr`. Valid values:
  - `wrr`: Backend servers with higher weights receive more requests than those with lower weights.
  - `rr`: Requests are distributed to backend servers in sequence.
* `server_group_id` - (Optional) The ID of the vServer group. It's the ID of resource `alicloud_slb_server_group`.
* `acl_status` - (Optional) Specifies whether to enable access control. Default value: `off`. Valid values: `on`, `off`.
* `acl_type` - (Optional) The type of the network ACL. Valid values: `black`, `white`. **NOTE:** If `acl_status` is set to `on`, `acl_type` is required. Otherwise, it will be ignored.
* `acl_id` - (Optional) The ID of the network ACL that is associated with the listener. **NOTE:** If `acl_status` is set to `on`, `acl_id` is required. Otherwise, it will be ignored.
* `sticky_session` - (Optional) Specifies whether to enable session persistence. Default value: `off`. Valid values: `on`, `off`.
* `sticky_session_type` - (Optional) The method that is used to handle a cookie. Valid values: `insert`, `server`. **NOTE:** If `sticky_session` is set to `on`, `sticky_session_type` is required. Otherwise, it will be ignored.
* `cookie_timeout` - (Optional, Int) The timeout period of a cookie. Unit: seconds. Valid values: `1` to `86400`. **NOTE:** If `sticky_session` is set to `on`, and `sticky_session_type` is set to `insert`, `cookie_timeout` is required. Otherwise, it will be ignored.
* `cookie` - (Optional) The cookie that is configured on the server. The `cookie` must be `1` to `200` characters in length and can contain only ASCII characters and digits. It cannot contain commas (,), semicolons (;), or space characters. It cannot start with a dollar sign ($). **NOTE:** If `sticky_session` is set to `on`, and `sticky_session_type` is set to `server`, `cookie` is required. Otherwise, it will be ignored.
* `health_check` - (Optional) Specifies whether to enable the health check feature. Default value: `on`. Valid values: `on`, `off`. **NOTE:** `TCP` and `UDP` listener's HealthCheck is always on, so it will be ignored when launching `TCP` or `UDP` listener.
* `health_check_method` - (Optional, Available since v1.70.0) The health check method used in HTTP health checks. Valid values: `head`, `get`. **NOTE:** `health_check_method` takes effect only if `health_check` is set to `on`.
* `health_check_domain` - (Optional) The domain name that is used for health checks. **NOTE:** `health_check_domain` takes effect only if `health_check` is set to `on`.
* `health_check_uri` - (Optional) The URI that is used for health checks. The `health_check_uri` must be `1` to `80` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), percent signs (%), question marks (?), number signs (#), and ampersands (&). The URI must start with a forward slash (/) but cannot be a single forward slash (/).
**NOTE:** `health_check_uri` takes effect only if `health_check` is set to `on`.
* `health_check_connect_port` - (Optional, Int) The backend port that is used for health checks. Valid values: `0` to `65535`. **NOTE:** `health_check_connect_port` takes effect only if `health_check` is set to `on`.
* `healthy_threshold` - (Optional, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. Default value: `3`. Valid values: `2` to `10`. **NOTE:** `healthy_threshold` takes effect only if `health_check` is set to `on`.
* `unhealthy_threshold` - (Optional, Int) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. Default value: `3`. Valid values: `2` to `10`. **NOTE:** `unhealthy_threshold` takes effect only if `health_check` is set to `on`.
* `health_check_timeout` - (Optional, Int) The timeout period of a health check response. Unit: seconds. Default value: `5`. Valid values: `1` to `300`. **NOTE:** If `health_check_timeout` < `health_check_interval`, `health_check_timeout` will be replaced by `health_check_interval`. `health_check_timeout` takes effect only if `health_check` is set to `on`.
* `health_check_interval` - (Optional, Int) The interval between two consecutive health checks. Unit: seconds. Default value: `2`. Valid values: `1` to `50`. **NOTE:** `health_check_interval` takes effect only if `health_check` is set to `on`.
* `health_check_http_code` - (Optional) The HTTP status code for a successful health check. Separate multiple HTTP status codes with commas (`,`). Default value: `http_2xx`. Valid values: `http_2xx`, `http_3xx`, `http_4xx` and `http_5xx`. **NOTE:** `health_check_http_code` takes effect only if `health_check` is set to `on`.
* `gzip` - (Optional, Bool, Available since v1.13.0) Specifies whether to enable GZIP compression to compress specific types of files. Default value: `true`. Valid values: `true`, `false`.
* `idle_timeout` - (Optional, Int) The timeout period of an idle connection. Unit: seconds. Default value: `15`. Valid values: `1` to `60`.
* `request_timeout` - (Optional, Int) The timeout period of a request. Unit: seconds. Default value: `60`. Valid values: `1` to `180`.
* `forward_port` - (Optional, ForceNew, Int, Available since v1.40.0) The listening port that is used to redirect HTTP requests to HTTPS.
* `listener_forward` - (Optional, ForceNew, Available since v1.40.0) Specifies whether to enable HTTP-to-HTTPS redirection. Default value: `off`. Valid values: `on`, `off`.
* `x_forwarded_for` - (Optional, Set, Available since v1.13.0) Whether to set additional HTTP Header field "X-Forwarded-For". See [`x_forwarded_for`](#x_forwarded_for) below.
* `description` - (Optional, Available since v1.69.0) The name of the listener. The name must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), and underscores (_).
* `delete_protection_validation` - (Optional, Bool, Available since v1.63.0) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default value: `false`.

### HTTPS Listener

The HTTPS Listener supports the following:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new listener.
* `protocol` - (Required, ForceNew) The protocol to listen on. Valid values: `https`.
* `frontend_port` - (Required, Int, ForceNew) The frontend port that is used by the CLB instance. Valid values: `1` to `65535`.
* `bandwidth` - (Required, Int) The maximum bandwidth of the listener. Unit: Mbit/s. Valid values:
  - `-1`: For a pay-by-data-transfer Internet-facing CLB instance, if you set `bandwidth` to `-1`, the bandwidth of the listener is unlimited.
  - `1` to `1000`: The sum of the maximum bandwidth that you specify for all listeners of the CLB instance cannot exceed the maximum bandwidth of the CLB instance.
-> **NOTE:** Currently, this `bandwidth` is available on `Domestic Site Account`.
* `backend_port` - (Optional, Int, ForceNew) The backend port that is used by the CLB instance. Valid values: `1` to `65535`. **NOTE:** If `server_group_id` is not set, `backend_port` is required.
* `scheduler` - (Optional) The scheduling algorithm. Default value: `wrr`. Valid values:
  - `wrr`: Backend servers with higher weights receive more requests than those with lower weights.
  - `rr`: Requests are distributed to backend servers in sequence.
* `server_group_id` - (Optional) The ID of the vServer group. It's the ID of resource `alicloud_slb_server_group`.
* `acl_status` - (Optional) Specifies whether to enable access control. Default value: `off`. Valid values: `on`, `off`.
* `acl_type` - (Optional) The type of the network ACL. Valid values: `black`, `white`. **NOTE:** If `acl_status` is set to `on`, `acl_type` is required. Otherwise, it will be ignored.
* `acl_id` - (Optional) The ID of the network ACL that is associated with the listener. **NOTE:** If `acl_status` is set to `on`, `acl_id` is required. Otherwise, it will be ignored.
* `sticky_session` - (Optional) Specifies whether to enable session persistence. Default value: `off`. Valid values: `on`, `off`.
* `sticky_session_type` - (Optional) The method that is used to handle a cookie. Valid values: `insert`, `server`. **NOTE:** If `sticky_session` is set to `on`, `sticky_session_type` is required. Otherwise, it will be ignored.
* `cookie_timeout` - (Optional, Int) The timeout period of a cookie. Unit: seconds. Valid values: `1` to `86400`. **NOTE:** If `sticky_session` is set to `on`, and `sticky_session_type` is set to `insert`, `cookie_timeout` is required. Otherwise, it will be ignored.
* `cookie` - (Optional) The cookie that is configured on the server. The `cookie` must be `1` to `200` characters in length and can contain only ASCII characters and digits. It cannot contain commas (,), semicolons (;), or space characters. It cannot start with a dollar sign ($). **NOTE:** If `sticky_session` is set to `on`, and `sticky_session_type` is set to `server`, `cookie` is required. Otherwise, it will be ignored.
* `health_check` - (Optional) Specifies whether to enable the health check feature. Default value: `on`. Valid values: `on`, `off`. **NOTE:** `TCP` and `UDP` listener's HealthCheck is always on, so it will be ignored when launching `TCP` or `UDP` listener.
* `health_check_method` - (Optional, Available since v1.70.0) The health check method used in HTTP health checks. Valid values: `head`, `get`. **NOTE:** `health_check_method` takes effect only if `health_check` is set to `on`.
* `health_check_domain` - (Optional) The domain name that is used for health checks. **NOTE:** `health_check_domain` takes effect only if `health_check` is set to `on`.
* `health_check_uri` - (Optional) The URI that is used for health checks. The `health_check_uri` must be `1` to `80` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), percent signs (%), question marks (?), number signs (#), and ampersands (&). The URI must start with a forward slash (/) but cannot be a single forward slash (/).
**NOTE:** `health_check_uri` takes effect only if `health_check` is set to `on`.
* `health_check_connect_port` - (Optional, Int) The backend port that is used for health checks. Valid values: `0` to `65535`. **NOTE:** `health_check_connect_port` takes effect only if `health_check` is set to `on`.
* `healthy_threshold` - (Optional, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. Default value: `3`. Valid values: `2` to `10`. **NOTE:** `healthy_threshold` takes effect only if `health_check` is set to `on`.
* `unhealthy_threshold` - (Optional, Int) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. Default value: `3`. Valid values: `2` to `10`. **NOTE:** `unhealthy_threshold` takes effect only if `health_check` is set to `on`.
* `health_check_timeout` - (Optional, Int) The timeout period of a health check response. Unit: seconds. Default value: `5`. Valid values: `1` to `300`. **NOTE:** If `health_check_timeout` < `health_check_interval`, `health_check_timeout` will be replaced by `health_check_interval`. `health_check_timeout` takes effect only if `health_check` is set to `on`.
* `health_check_interval` - (Optional, Int) The interval between two consecutive health checks. Unit: seconds. Default value: `2`. Valid values: `1` to `50`. **NOTE:** `health_check_interval` takes effect only if `health_check` is set to `on`.
* `health_check_http_code` - (Optional) The HTTP status code for a successful health check. Separate multiple HTTP status codes with commas (`,`). Default value: `http_2xx`. Valid values: `http_2xx`, `http_3xx`, `http_4xx` and `http_5xx`. **NOTE:** `health_check_http_code` takes effect only if `health_check` is set to `on`.
* `server_certificate_id` - (Optional, Available since v1.59.0) The ID of the server certificate. **NOTE:** `server_certificate_id` is also required when the value of the `ssl_certificate_id` is Empty.
* `ca_certificate_id` - (Optional, Available since v1.104) The ID of the certification authority (CA) certificate.
* `gzip` - (Optional, Bool, Available since v1.13.0) Specifies whether to enable GZIP compression to compress specific types of files. Default value: `true`. Valid values: `true`, `false`.
* `idle_timeout` - (Optional, Int) The timeout period of an idle connection. Unit: seconds. Default value: `15`. Valid values: `1` to `60`.
* `request_timeout` - (Optional, Int) The timeout period of a request. Unit: seconds. Default value: `60`. Valid values: `1` to `180`.
* `enable_http2` - (Optional) Specifies whether to enable HTTP/2. Default value: `on`. Valid values: `on`, `off`.
* `tls_cipher_policy` - (Optional) The Transport Layer Security (TLS) security policy. Default value: `tls_cipher_policy_1_0`. Valid values: `tls_cipher_policy_1_0`, `tls_cipher_policy_1_1`, `tls_cipher_policy_1_2`, `tls_cipher_policy_1_2_strict`, `tls_cipher_policy_1_2_strict_with_1_3`. **NOTE:** From version 1.229.1, `tls_cipher_policy` can be set to `tls_cipher_policy_1_2_strict_with_1_3`.
* `x_forwarded_for` - (Optional, Set, Available since v1.13.0) Whether to set additional HTTP Header field "X-Forwarded-For". See [`x_forwarded_for`](#x_forwarded_for) below.
* `description` - (Optional, Available since v1.69.0) The name of the listener. The name must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), and underscores (_).
* `delete_protection_validation` - (Optional, Bool, Available since v1.63.0) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default value: `false`.
* `ssl_certificate_id` - (Deprecated since v1.59.0) The ID of the server certificate. **NOTE:** Field `ssl_certificate_id` has been deprecated from provider version 1.59.0. New field `server_certificate_id` instead.

### TCP Listener

The TCP Listener supports the following:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new listener.
* `protocol` - (Required, ForceNew) The protocol to listen on. Valid values: `tcp`.
* `frontend_port` - (Required, Int, ForceNew) The frontend port that is used by the CLB instance. Valid values: `1` to `65535`.
* `bandwidth` - (Required, Int) The maximum bandwidth of the listener. Unit: Mbit/s. Valid values:
  - `-1`: For a pay-by-data-transfer Internet-facing CLB instance, if you set `bandwidth` to `-1`, the bandwidth of the listener is unlimited.
  - `1` to `1000`: The sum of the maximum bandwidth that you specify for all listeners of the CLB instance cannot exceed the maximum bandwidth of the CLB instance.
-> **NOTE:** Currently, this `bandwidth` is available on `Domestic Site Account`.
* `backend_port` - (Optional, Int, ForceNew) The backend port that is used by the CLB instance. Valid values: `1` to `65535`. **NOTE:** If `server_group_id` is not set, `backend_port` is required.
* `scheduler` - (Optional) The scheduling algorithm. Default value: `wrr`. Valid values:
  - `wrr`: Backend servers with higher weights receive more requests than those with lower weights.
  - `rr`: Requests are distributed to backend servers in sequence.
  - `sch`: Specifies consistent hashing that is based on source IP addresses. Requests from the same source IP address are distributed to the same backend server.
  - `tch`: Specifies consistent hashing that is based on four factors: source IP address, destination IP address, source port, and destination port. Requests that contain the same information based on the four factors are distributed to the same backend server.
**NOTE:** Only high-performance CLB instances support the `sch` and `tch` consistent hashing algorithms.
* `server_group_id` - (Optional) The ID of the vServer group. It's the ID of resource `alicloud_slb_server_group`.
* `master_slave_server_group_id` - (Optional) The ID of the primary/secondary server group. **NOTE:** You cannot set both `server_group_id` and `master_slave_server_group_id`.
* `acl_status` - (Optional) Specifies whether to enable access control. Default value: `off`. Valid values: `on`, `off`.
* `acl_type` - (Optional) The type of the network ACL. Valid values: `black`, `white`. **NOTE:** If `acl_status` is set to `on`, `acl_type` is required. Otherwise, it will be ignored.
* `acl_id` - (Optional) The ID of the network ACL that is associated with the listener. **NOTE:** If `acl_status` is set to `on`, `acl_id` is required. Otherwise, it will be ignored.
* `persistence_timeout` - (Optional, Int) The timeout period of session persistence. Unit: seconds. Default value: `0`. Valid values: `0` to `3600`.
* `health_check_type` - (Optional) The type of health checks. Default value: `tcp`. Valid values: `tcp`, `http`.
* `health_check_domain` - (Optional) The domain name that is used for health checks.
* `health_check_uri` - (Optional) The URI that is used for health checks. The `health_check_uri` must be `1` to `80` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), percent signs (%), question marks (?), number signs (#), and ampersands (&). The URI must start with a forward slash (/) but cannot be a single forward slash (/).
**NOTE:** You can set `health_check_uri` when the `TCP` listener requires `HTTP` health checks. If you do not set `health_check_uri`, `TCP` health checks will be performed.
* `health_check_connect_port` - (Optional, Int) The backend port that is used for health checks. Valid values: `0` to `65535`. **NOTE:** If `health_check_connect_port` is not set, the backend port specified by `backend_port` is used for health checks.
* `healthy_threshold` - (Optional, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. Default value: `3`. Valid values: `2` to `10`.
* `unhealthy_threshold` - (Optional, Int) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. Default value: `3`. Valid values: `2` to `10`.
* `health_check_timeout` - (Optional, Int) The maximum timeout period of a health check response. Unit: seconds. Default value: `5`. Valid values: `1` to `300`.
* `health_check_interval` - (Optional, Int) The interval between two consecutive health checks. Unit: seconds. Default value: `2`. Valid values: `1` to `50`.
* `health_check_http_code` - (Optional) The HTTP status code for a successful health check. Separate multiple HTTP status codes with commas (`,`). Default value: `http_2xx`. Valid values: `http_2xx`, `http_3xx`, `http_4xx` and `http_5xx`.
* `established_timeout` - (Optional, Int) The timeout period of a connection. Unit: seconds. Default value: `900`. Valid values: `10` to `900`.
* `proxy_protocol_v2_enabled` - (Optional, Bool, Available since v1.187.0) Specifies whether to use the Proxy protocol to pass client IP addresses to backend servers. Default value: `false`. Valid values: `true`, `false`.
* `description` - (Optional, Available since v1.69.0) The name of the listener. The name must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), and underscores (_).
* `delete_protection_validation` - (Optional, Bool, Available since v1.63.0) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default value: `false`.

### UDP Listener

The UDP Listener supports the following:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new listener.
* `protocol` - (Required, ForceNew) The protocol to listen on. Valid values: `udp`.
* `frontend_port` - (Required, Int, ForceNew) The frontend port that is used by the CLB instance. Valid values: `1` to `65535`.
* `bandwidth` - (Required, Int) The maximum bandwidth of the listener. Unit: Mbit/s. Valid values:
  - `-1`: For a pay-by-data-transfer Internet-facing CLB instance, if you set `bandwidth` to `-1`, the bandwidth of the listener is unlimited.
  - `1` to `1000`: The sum of the maximum bandwidth that you specify for all listeners of the CLB instance cannot exceed the maximum bandwidth of the CLB instance.
-> **NOTE:** Currently, this `bandwidth` is available on `Domestic Site Account`.
* `backend_port` - (Optional, Int, ForceNew) The backend port that is used by the CLB instance. Valid values: `1` to `65535`. **NOTE:** If `server_group_id` is not set, `backend_port` is required.
* `scheduler` - (Optional) The scheduling algorithm. Default value: `wrr`. Valid values:
  - `wrr`: Backend servers with higher weights receive more requests than those with lower weights.
  - `rr`: Requests are distributed to backend servers in sequence.
  - `sch`: Specifies consistent hashing that is based on source IP addresses. Requests from the same source IP address are distributed to the same backend server.
  - `tch`: Specifies consistent hashing that is based on four factors: source IP address, destination IP address, source port, and destination port. Requests that contain the same information based on the four factors are distributed to the same backend server.
  - `qch`: Specifies consistent hashing that is based on QUIC connection IDs. Requests that contain the same QUIC connection ID are distributed to the same backend server.
**NOTE:** Only high-performance CLB instances support the `sch`, `tch`, and `qch` consistent hashing algorithms.
* `server_group_id` - (Optional) The ID of the vServer group. It's the ID of resource `alicloud_slb_server_group`.
* `master_slave_server_group_id` - (Optional) The ID of the primary/secondary server group. **NOTE:** You cannot set both `server_group_id` and `master_slave_server_group_id`.
* `acl_status` - (Optional) Specifies whether to enable access control. Default value: `off`. Valid values: `on`, `off`.
* `acl_type` - (Optional) The type of the network ACL. Valid values: `black`, `white`. **NOTE:** If `acl_status` is set to `on`, `acl_type` is required. Otherwise, it will be ignored.
* `acl_id` - (Optional) The ID of the network ACL that is associated with the listener. **NOTE:** If `acl_status` is set to `on`, `acl_id` is required. Otherwise, it will be ignored.
* `health_check_connect_port` - (Optional, Int) The backend port that is used for health checks. Valid values: `0` to `65535`. **NOTE:** If `health_check_connect_port` is not set, the backend port specified by `backend_port` is used for health checks.
* `healthy_threshold` - (Optional, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. Default value: `3`. Valid values: `2` to `10`.
* `unhealthy_threshold` - (Optional, Int) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. Default value: `3`. Valid values: `2` to `10`.
* `health_check_timeout` - (Optional, Int) The maximum timeout period of a health check response. Unit: seconds. Default value: `5`. Valid values: `1` to `300`.
* `health_check_interval` - (Optional, Int) The interval between two consecutive health checks. Unit: seconds. Default value: `2`. Valid values: `1` to `50`.
* `proxy_protocol_v2_enabled` - (Optional, Bool, Available since v1.187.0) Specifies whether to use the Proxy protocol to pass client IP addresses to backend servers. Default value: `false`. Valid values: `true`, `false`.
* `description` - (Optional, Available since v1.69.0) The name of the listener. The name must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), and underscores (_).
* `delete_protection_validation` - (Optional, Bool, Available since v1.63.0) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default value: `false`.

-> **NOTE:** Once enable the http redirect to https function, any parameters excepted forward_port, listener_forward, load_balancer_id, frontend_port, protocol will be ignored. More info, please refer to [Redirect http to https](https://www.alibabacloud.com/help/doc-detail/89151.htm?spm=a2c63.p38356.b99.186.42f66384mpjUTB).

-> **NOTE:** Advantanced feature such as `tls_cipher_policy`, can not be updated when load balancer instance is "Shared-Performance". More info, please refer to [Configure a HTTPS Listener](https://www.alibabacloud.com/help/doc-detail/27593.htm).

### `x_forwarded_for`

The x_forwarded_for mapping supports the following:

* `retrive_slb_ip` - (Optional, Bool) Indicates whether the SLB-IP header is used to retrieve the virtual IP address (VIP) requested by the client. Default value: `false`. Valid values: `true`, `false`.
* `retrive_slb_id` - (Optional, Bool) Indicates whether the SLB-ID header is used to retrieve the ID of the CLB instance. Default value: `false`. Valid values: `true`, `false`.
* `retrive_slb_proto` - (Optional, Bool) Specifies whether to use the X-Forwarded-Proto header to retrieve the listener protocol. Default value: `false`. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Listener. It formats as `<load_balancer_id>:<protocol>:<frontend_port>`.
-> **NOTE:** Before provider version 1.57.1, it formats as `<load_balancer_id>:<frontend_port>`.
* `x_forwarded_for` - Whether to set additional HTTP Header field "X-Forwarded-For".
  * `retrive_client_ip` - Whether to retrieve the client ip.

## Import

Classic Load Balancer (SLB) Load Balancer Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_listener.example <load_balancer_id>:<protocol>:<frontend_port>
```

**NOTE:** Before provider version 1.57.1, Classic Load Balancer (SLB) Load Balancer Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_listener.example <load_balancer_id>:<frontend_port>
```
