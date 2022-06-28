---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_listener"
sidebar_current: "docs-alicloud-resource-slb-listener"
description: |-
  Provides an Application Load Banlancer resource.
---

# alicloud\_slb\_listener

Provides an Application Load Balancer Listener resource.

For information about slb and how to use it, see [What is Server Load Balancer](https://www.alibabacloud.com/help/doc-detail/27539.htm).

For information about listener and how to use it, to see the following:

* [Configure a HTTP Listener](https://www.alibabacloud.com/help/doc-detail/27592.htm).
* [Configure a HTTPS Listener](https://www.alibabacloud.com/help/doc-detail/27593.htm).
* [Configure a TCP Listener](https://www.alibabacloud.com/help/doc-detail/27594.htm).
* [Configure a UDP Listener](https://www.alibabacloud.com/help/doc-detail/27595.htm).

## Example Usage

```
variable "name" {
  default = "testcreatehttplistener"
}

variable "ip_version" {
  default = "ipv4"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name   = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet             = true
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 80
  frontend_port             = 80
  protocol                  = "http"
  bandwidth                 = 10
  sticky_session            = "on"
  sticky_session_type       = "insert"
  cookie_timeout            = 86400
  cookie                    = "testslblistenercookie"
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
  acl_id          = alicloud_slb_acl.default.id
  request_timeout = 80
  idle_timeout    = 30
}

resource "alicloud_slb_acl" "default" {
  name       = var.name
  ip_version = var.ip_version
  entry_list {
    entry   = "10.10.10.0/24"
    comment = "first"
  }
  entry_list {
    entry   = "168.10.10.0/24"
    comment = "second"
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new listener.
* `frontend_port` - (Required, ForceNew) Port used by the Server Load Balancer instance frontend. Valid value range: [1-65535].
* `backend_port` - (Optional, ForceNew) Port used by the Server Load Balancer instance backend. Valid value range: [1-65535].
* `protocol` - (Required, ForceNew) The protocol to listen on. Valid values are [`http`, `https`, `tcp`, `udp`].
* `bandwidth` - (Optional, Computed) Bandwidth peak of Listener. For the public network instance charged per traffic consumed, the Bandwidth on Listener can be set to -1, indicating the bandwidth peak is unlimited. Valid values are [-1, 1-1000] in Mbps.
* `description` - (Optional, Available in 1.69.0+) The description of slb listener. This description can have a string of 1 to 80 characters. Default value: null.
* `scheduler` - (Optional) Scheduling algorithm,  Valid values: `wrr`, `rr`, `wlc`, `sch`, `tcp`, `qch`. Default to `wrr`. 
  Only when `protocol` is `tcp` or `udp`, `scheduler` can be set to `sch`. Only when instance is guaranteed-performance instance and `protocol` is `tcp` or `udp`, `scheduler` can be set to `tch`. Only when instance is guaranteed-performance instance and `protocol` is `udp`, `scheduler` can be set to `qch`.
* `sticky_session` - (Optional) Whether to enable session persistence, Valid values are `on` and `off`. Default to `off`.
* `sticky_session_type` - (Optional) Mode for handling the cookie. If `sticky_session` is "on", it is mandatory. Otherwise, it will be ignored. Valid values are `insert` and `server`. `insert` means it is inserted from Server Load Balancer; `server` means the Server Load Balancer learns from the backend server.
* `cookie_timeout` - (Optional) Cookie timeout. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "insert". Otherwise, it will be ignored. Valid value range: [1-86400] in seconds.
* `cookie` - (Optional) The cookie configured on the server. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "server". Otherwise, it will be ignored. Valid value：String in line with RFC 2965, with length being 1- 200. It only contains characters such as ASCII codes, English letters and digits instead of the comma, semicolon or spacing, and it cannot start with $.
* `persistence_timeout` - (Optional) Timeout of connection persistence. Valid value range: [0-3600] in seconds. Default to 0 and means closing it.
* `health_check` - (Optional) Whether to enable health check. Valid values are`on` and `off`. TCP and UDP listener's HealthCheck is always on, so it will be ignore when launching TCP or UDP listener.
* `health_check_type` - (Optional) Type of health check. Valid values are: `tcp` and `http`. Default to `tcp` . TCP supports TCP and HTTP health check mode, you can select the particular mode depending on your application.
* `health_check_domain` - (Optional) Domain name used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and only characters such as letters, digits, ‘-‘ and ‘.’ are allowed. When it is not set or empty,  Server Load Balancer uses the private network IP address of each backend server as Domain used for health check.
* `health_check_uri` - (Optional) URI used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and it must start with /. Only characters such as letters, digits, ‘-’, ‘/’, ‘.’, ‘%’, ‘?’, #’ and ‘&’ are allowed.
* `health_check_connect_port` - (Optional) The port that is used for health checks. Valid value range: [0-65535]. Default to `0` means that the port on a backend server is used for health checks.
* `healthy_threshold` - (Optional) The number of health checks that an unhealthy backend server must consecutively pass before it can be declared healthy. In this case, the health check state is changed from fail to success. It is required when `health_check` is on. Valid value range: [2-10] in seconds. Default to 3. **NOTE:** This parameter takes effect only if the `health_check` parameter is set to `on`.
* `unhealthy_threshold` - (Optional) The number of health checks that a healthy backend server must consecutively fail before it can be declared unhealthy. In this case, the health check state is changed from success to fail. It is required when `health_check` is on. Valid value range: [2-10] in seconds. Default to 3. **NOTE:** This parameter takes effect only if the `health_check` parameter is set to `on`.
* `health_check_timeout` - (Optional) Maximum timeout of each health check response. It is required when `health_check` is on. Valid value range: [1-300] in seconds. Default to 5. Note: If `health_check_timeout` < `health_check_interval`, its will be replaced by `health_check_interval`.
* `health_check_interval` - (Optional) Time interval of health checks. It is required when `health_check` is on. Valid value range: [1-50] in seconds. Default to 2.
* `health_check_http_code` - (Optional) Regular health check HTTP status code. Multiple codes are segmented by “,”. It is required when `health_check` is on. Default to `http_2xx`.  Valid values are: `http_2xx`,  `http_3xx`, `http_4xx` and `http_5xx`.
* `health_check_method` - (Optional, Available in 1.70.0+) HealthCheckMethod used for health check.Valid values: ["head", "get"] `http` and `https` support regions ap-northeast-1, ap-southeast-1, ap-southeast-2, ap-southeast-3, us-east-1, us-west-1, eu-central-1, ap-south-1, me-east-1, cn-huhehaote, cn-zhangjiakou, ap-southeast-5, cn-shenzhen, cn-hongkong, cn-qingdao, cn-chengdu, eu-west-1, cn-hangzhou", cn-beijing, cn-shanghai.This function does not support the TCP protocol .
* `ssl_certificate_id` - (Deprecated) SLB Server certificate ID. It has been deprecated from 1.59.0 and using `server_certificate_id` instead. 
* `server_certificate_id` - (Optional, Available in 1.59.0+) SLB Server certificate ID. It is required when `protocol` is `https`. The `server_certificate_id` is also required when the value of the `ssl_certificate_id`  is Empty.
* `ca_certificate_id` - (Optional, Available in 1.104) SLB CA certificate ID. Only when `protocol` is `https` can be specified.
* `gzip` - (Optional) Whether to enable "Gzip Compression". If enabled, files of specific file types will be compressed, otherwise, no files will be compressed. Default to true. Available in v1.13.0+.
* `x_forwarded_for` - (Optional) Whether to set additional HTTP Header field "X-Forwarded-For" (documented below). Available in v1.13.0+. The details see Block `x_forwarded_for`.
* `acl_status` - (Optional) Whether to enable "acl(access control list)", the acl is specified by `acl_id`. Valid values are `on` and `off`. Default to `off`.
* `acl_type` - (Optional) Mode for handling the acl specified by acl_id. If `acl_status` is "on", it is mandatory. Otherwise, it will be ignored. Valid values are `white` and `black`. `white` means the Listener can only be accessed by client ip belongs to the acl; `black` means the Listener can not be accessed by client ip belongs to the acl.
* `acl_id` - (Optional) the id of access control list to be apply on the listener, is the id of resource alicloud_slb_acl. If `acl_status` is "on", it is mandatory. Otherwise, it will be ignored.
* `established_timeout` - (Optional) Timeout of tcp listener established connection idle timeout. Valid value range: [10-900] in seconds. Default to 900.
* `idle_timeout` - (Optional) Timeout of http or https listener established connection idle timeout. Valid value range: [1-60] in seconds. Default to 15.
* `request_timeout` - (Optional) Timeout of http or https listener request (which does not get response from backend) timeout. Valid value range: [1-180] in seconds. Default to 60.
* `enable_http2` - (Optional) Whether to enable https listener support http2 or not. Valid values are `on` and `off`. Default to `on`.
* `tls_cipher_policy` - (Optional)  Https listener TLS cipher policy. Valid values are `tls_cipher_policy_1_0`, `tls_cipher_policy_1_1`, `tls_cipher_policy_1_2`, `tls_cipher_policy_1_2_strict`. Default to `tls_cipher_policy_1_0`. Currently the `tls_cipher_policy` can not be updated when load balancer instance is "Shared-Performance".
* `server_group_id` - (Optional) the id of server group to be apply on the listener, is the id of resource `alicloud_slb_server_group`.
* `listener_forward` - (Optional, ForceNew, Available in 1.40.0+) Whether to enable http redirect to https, Valid values are `on` and `off`. Default to `off`.
* `master_slave_server_group_id` - (Optional) The ID of the master slave server group.
* `forward_port` - (Optional, ForceNew, Available in 1.40.0+) The port that http redirect to https.
* `delete_protection_validation` - (Optional, Available in 1.63.0+) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

-> **NOTE:** Once enable the http redirect to https function, any parameters excepted forward_port,listener_forward,load_balancer_id,frontend_port,protocol will be ignored. More info, please refer to [Redirect http to https](https://www.alibabacloud.com/help/doc-detail/89151.htm?spm=a2c63.p38356.b99.186.42f66384mpjUTB).

-> **NOTE:** Advantanced feature such as `tls_cipher_policy`, can not be updated when load balancer instance is "Shared-Performance". More info, please refer to [Configure a HTTPS Listener](https://www.alibabacloud.com/help/doc-detail/27593.htm).

### Block x_forwarded_for

The x_forwarded_for mapping supports the following:

* `retrive_slb_ip` - (Optional) Whether to use the XForwardedFor_SLBIP header to obtain the public IP address of the SLB instance. Default to false.
* `retrive_slb_id` - (Optional) Whether to use the XForwardedFor header to obtain the ID of the SLB instance. Default to false.
* `retrive_slb_proto` - (Optional) Whether to use the XForwardedFor_proto header to obtain the protocol used by the listener. Default to false.

## Listener fields and protocol mapping

load balance support 4 protocol to listen on, they are `http`,`https`,`tcp`,`udp`, the every listener support which portocal following:

listener parameter | support protocol | value range |
------------- | ------------- | ------------- | 
backend_port | http & https & tcp & udp | 1-65535 | 
frontend_port | http & https & tcp & udp | 1-65535 |
protocol | http & https & tcp & udp |
bandwidth | http & https & tcp & udp | -1 / 1-1000 |
scheduler | http & https & tcp & udp | wrr, rr, wlc, tch, qch |
sticky_session | http & https | on or off |
sticky_session_type | http & https | insert or server | 
cookie_timeout | http & https | 1-86400  | 
cookie | http & https |   | 
persistence_timeout | tcp & udp | 0-3600 | 
health_check | http & https | on or off | 
health_check_type | tcp | tcp or http | 
health_check_domain | http & https & tcp | 
health_check_method | http & https & tcp | 
health_check_uri | http & https & tcp |  | 
health_check_connect_port | http & https & tcp & udp | 1-65535 or -520 | 
healthy_threshold | http & https & tcp & udp | 1-10 | 
unhealthy_threshold | http & https & tcp & udp | 1-10 | 
health_check_timeout | http & https & tcp & udp | 1-300 |
health_check_interval | http & https & tcp & udp | 1-50 |
health_check_http_code | http & https & tcp | http_2xx,http_3xx,http_4xx,http_5xx | 
server_certificate_id | https |  |
gzip | http & https | true or false  |
x_forwarded_for | http & https |  |
acl_status | http & https & tcp & udp | on or off |
acl_type   | http & https & tcp & udp | white or black |
acl_id     | http & https & tcp & udp | the id of resource alicloud_slb_acl|
established_timeout | tcp       | 10-900|
idle_timeout |http & https      | 1-60  |
request_timeout |http & https   | 1-180 |
enable_http2    |https          | on or off |
tls_cipher_policy |https        |  tls_cipher_policy_1_0, tls_cipher_policy_1_1, tls_cipher_policy_1_2, tls_cipher_policy_1_2_strict |
server_group_id    | http & https & tcp & udp | the id of resource alicloud_slb_server_group |

The listener mapping supports the following:

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer listener. Its format as `<load_balancer_id>:<protocol>:<frontend_port>`. Before verson 1.57.1, the foramt as `<load_balancer_id>:<frontend_port>`.
* `load_balancer_id` - The Load Balancer ID which is used to launch a new listener.
* `frontend_port` - Port used by the Server Load Balancer instance frontend.
* `backend_port` - Port used by the Server Load Balancer instance backend.
* `protocol` - The protocol to listen on.
* `bandwidth` - Bandwidth peak of Listener.
* `scheduler` - Scheduling algorithm.
* `sticky_session` - Whether to enable session persistence.
* `sticky_session_type` - Mode for handling the cookie.
* `cookie_timeout` - Cookie timeout.
* `cookie` - The cookie configured on the server.
* `persistence_timeout` - Timeout of connection persistence.
* `health_check` - Whether to enable health check.
* `health_check_type` - Type of health check.
* `health_check_domain` - Domain name used for health check.
* `health_check_method` - HealthCheckMethod used for health check.
* `health_check_uri` - URI used for health check.
* `health_check_connect_port` - Port used for health check.
* `healthy_threshold` - Threshold determining the result of the health check is success.
* `unhealthy_threshold` - Threshold determining the result of the health check is fail.
* `health_check_timeout` - Maximum timeout of each health check response.
* `health_check_interval` - Time interval of health checks.
* `health_check_http_code` - Regular health check HTTP status code.
* `server_certificate_id` - (Optional) Security certificate ID.

## Import

Load balancer listener can be imported using the id, e.g.

```
$ terraform import alicloud_slb_listener.example "lb-abc123456:tcp:22"
```
