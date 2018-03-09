---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_listener"
sidebar_current: "docs-alicloud-resource-slb-listener"
description: |-
  Provides an Application Load Banlancer resource.
---

# alicloud\_slb\_listener

Provides an Application Load Balancer Listener resource.

## Example Usage

```
# Create a new load balancer and listeners
resource "alicloud_slb" "instance" {
  name                 = "test-slb-tf"
  internet             = true
  internet_charge_type = "paybybandwidth"
  bandwidth            = 25
}

resource "alicloud_slb_listener" "http" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = 80
  frontend_port = 80
  bandwidth = 10
  protocol = "http"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 86400
}
resource "alicloud_slb_listener" "tcp" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new listener.
* `frontend_port` - (Required, ForceNew) Port used by the Server Load Balancer instance frontend. Valid value range: [1-65535].
* `backend_port` - (Required, ForceNew) Port used by the Server Load Balancer instance backend. Valid value range: [1-65535].
* `protocol` - (Required, ForceNew) The protocol to listen on. Valid values are [`http`, `https`, `tcp`, `udp`].
* `bandwidth` - (Required) Bandwidth peak of Listener. For the public network instance charged per traffic consumed, the Bandwidth on Listener can be set to -1, indicating the bandwidth peak is unlimited. Valid values are [-1, 1-1000] in Mbps.
* `scheduler` - (Optinal) Scheduling algorithm, Valid values are `wrr` and `wlc`.  Default to "wrr".
* `sticky_session` - (Optinal) Whether to enable session persistence, Valid values are `on` and `off`. Default to `off`.
* `sticky_session_type` - (Optinal) Mode for handling the cookie. If `sticky_session` is "on", it is mandatory. Otherwise, it will be ignored. Valid values are `insert` and `server`. `insert` means it is inserted from Server Load Balancer; `server` means the Server Load Balancer learns from the backend server.
* `cookie_timeout` - (Optinal) Cookie timeout. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "insert". Otherwise, it will be ignored. Valid value range: [1-86400] in seconds.
* `cookie` - (Optinal) The cookie configured on the server. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "server". Otherwise, it will be ignored. Valid value：String in line with RFC 2965, with length being 1- 200. It only contains characters such as ASCII codes, English letters and digits instead of the comma, semicolon or spacing, and it cannot start with $.
* `persistence_timeout` - (Optinal) Timeout of connection persistence. Valid value range: [0-3600] in seconds. Default to 0 and means closing it.
* `health_check` - (Optinal) Whether to enable health check. Valid values are`on` and `off`. TCP and UDP listener's HealthCheck is always on, so it will be ignore when launching TCP or UDP listener.
* `health_check_type` - (Optinal) Type of health check. Valid values are: `tcp` and `http`. Default to `tcp` . TCP supports TCP and HTTP health check mode, you can select the particular mode depending on your application.
* `health_check_domain` - (Optinal) Domain name used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and only characters such as letters, digits, ‘-‘ and ‘.’ are allowed. When it is not set or empty,  Server Load Balancer uses the private network IP address of each backend server as Domain used for health check.
* `health_check_uri` - (Optinal) URI used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and it must start with /. Only characters such as letters, digits, ‘-’, ‘/’, ‘.’, ‘%’, ‘?’, #’ and ‘&’ are allowed.
* `health_check_connect_port` - (Optinal) Port used for health check. Valid value range: [1-65535]. Default to "None" means the backend server port is used.
* `healthy_threshold` - (Optinal) Threshold determining the result of the health check is success. It is required when `health_check` is on. Valid value range: [1-10] in seconds. Default to 3.
* `unhealthy_threshold` - (Optinal) Threshold determining the result of the health check is fail. It is required when `health_check` is on. Valid value range: [1-10] in seconds. Default to 3.
* `health_check_timeout` - (Optinal) Maximum timeout of each health check response. It is required when `health_check` is on. Valid value range: [1-300] in seconds. Default to 5. Note: If `health_check_timeout` < `health_check_interval`, its will be replaced by `health_check_interval`.
* `health_check_interval` - (Optinal) Time interval of health checks. It is required when `health_check` is on. Valid value range: [1-50] in seconds. Default to 2.
* `health_check_http_code` - (Optinal) Regular health check HTTP status code. Multiple codes are segmented by “,”. It is required when `health_check` is on. Default to `http_2xx`.  Valid values are: `http_2xx`,  `http_3xx`, `http_4xx` and `http_5xx`.
* `ssl_certificate_id` - (Optinal) Security certificate ID.
* `server_group_id` - (Optinal) ID of a virtual server group that will be forwarded.

## Listener fields and protocol mapping

load balance support 4 protocal to listen on, they are `http`,`https`,`tcp`,`udp`, the every listener support which portocal following:

listener parameter | support protocol | value range |
------------- | ------------- | ------------- | 
backend_port | http & https & tcp & udp | 1-65535 | 
frontend_port | http & https & tcp & udp | 1-65535 |
protocol | http & https & tcp & udp |
bandwidth | http & https & tcp & udp | -1 / 1-1000 |
scheduler | http & https & tcp & udp | wrr or wlc |
sticky_session | http & https | on or off |
sticky_session_type | http & https | insert or server | 
cookie_timeout | http & https | 1-86400  | 
cookie | http & https |   | 
persistence_timeout | tcp & udp | 0-3600 | 
health_check | http & https | on or off | 
health_check_type | tcp | tcp or http | 
health_check_domain | http & https & tcp | 
health_check_uri | http & https & tcp |  | 
health_check_connect_port | http & https & tcp & udp | 1-65535 or -520 | 
healthy_threshold | http & https & tcp & udp | 1-10 | 
unhealthy_threshold | http & https & tcp & udp | 1-10 | 
health_check_timeout | http & https & tcp & udp | 1-300 |
health_check_interval | http & https & tcp & udp | 1-50 |
health_check_http_code | http & https & tcp | http_2xx,http_3xx,http_4xx,http_5xx | 
ssl_certificate_id | https |  |  


The listener mapping supports the following:

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer listener. It is consist of `load_balancer_id` and `frontend_port`: <load_balancer_id>:<frontend_port>.
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
* `health_check_uri` - URI used for health check.
* `health_check_connect_port` - Port used for health check.
* `healthy_threshold` - Threshold determining the result of the health check is success.
* `unhealthy_threshold` - Threshold determining the result of the health check is fail.
* `health_check_timeout` - Maximum timeout of each health check response.
* `health_check_interval` - Time interval of health checks.
* `health_check_http_code` - Regular health check HTTP status code.
* `ssl_certificate_id` - (Optinal) Security certificate ID.
* `server_group_id` - (Optinal) ID of a virtual server group that will be forwarded.

## Import

Load balancer listener can be imported using the id, e.g.

```
$ terraform import alicloud_slb_listener.example "lb-abc123456:22"
```
