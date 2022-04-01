---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_rule"
sidebar_current: "docs-alicloud-resource-slb-rule"
description: |-
  Provides a Load Banlancer Forwarding Rule Resource and add it to one Listener.
---

# alicloud\_slb\_rule

A forwarding rule is configured in `HTTP`/`HTTPS` listener and it used to listen a list of backend servers which in one specified virtual backend server group.
You can add forwarding rules to a listener to forward requests based on the domain names or the URL in the request.

-> **NOTE:** One virtual backend server group can be attached in multiple forwarding rules.

-> **NOTE:** At least one "Domain" or "Url" must be specified when creating a new rule.

-> **NOTE:** Having the same 'Domain' and 'Url' rule can not be created repeatedly in the one listener.

-> **NOTE:** Rule only be created in the `HTTP` or `HTTPS` listener.

-> **NOTE:** Only rule's virtual server group can be modified.

## Example Usage

```
variable "name" {
  default = "slbrulebasicconfig"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/16"
  zone_id           = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id = alicloud_vswitch.default.id
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 22
  frontend_port             = 22
  protocol                  = "http"
  bandwidth                 = 5
  health_check_connect_port = "20"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  servers {
    server_ids = alicloud_instance.default.*.id
    port       = 80
    weight     = 100
  }
}

resource "alicloud_slb_rule" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  frontend_port             = alicloud_slb_listener.default.frontend_port
  name                      = var.name
  domain                    = "*.aliyun.com"
  url                       = "/image"
  server_group_id           = alicloud_slb_server_group.default.id
  cookie                    = "23ffsa"
  cookie_timeout            = 100
  health_check_http_code    = "http_2xx"
  health_check_interval     = 10
  health_check_uri          = "/test"
  health_check_connect_port = 80
  health_check_timeout      = 30
  healthy_threshold         = 3
  unhealthy_threshold       = 5
  sticky_session            = "on"
  sticky_session_type       = "server"
  listener_sync             = "off"
  scheduler                 = "rr"
  health_check_domain       = "test"
  health_check              = "on"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch the new forwarding rule.
* `name` - (Optional) Name of the forwarding rule. Our plugin provides a default name: "tf-slb-rule".
* `frontend_port` - (Required, ForceNew) The listener frontend port which is used to launch the new forwarding rule. Valid range: [1-65535].
* `domain` - (Optional, ForceNew) Domain name of the forwarding rule. It can contain letters a-z, numbers 0-9, hyphens (-), and periods (.),
and wildcard characters. The following two domain name formats are supported:
   - Standard domain name: www.test.com
   - Wildcard domain name: *.test.com. wildcard (*) must be the first character in the format of (*.)
* `url` - (Optional, ForceNew) Domain of the forwarding rule. It must be 2-80 characters in length. Only letters a-z, numbers 0-9,
and characters '-' '/' '?' '%' '#' and '&' are allowed. URLs must be started with the character '/', but cannot be '/' alone.
* `server_group_id` - (Required) ID of a virtual server group that will be forwarded.
* `scheduler` - (Optional, Available in v1.51.0+) Scheduling algorithm, Valid values are `wrr`, `rr` and `wlc`.  Default to "wrr". This parameter is required  and takes effect only when ListenerSync is set to off.
* `sticky_session` - (Optional, Available in v1.51.0+) Whether to enable session persistence, Valid values are `on` and `off`. Default to `off`. This parameter is required  and takes effect only when ListenerSync is set to off.                                                                                                                                                                                                                                                 
* `sticky_session_type` - (Optional, Available in v1.51.0+) Mode for handling the cookie. If `sticky_session` is "on", it is mandatory. Otherwise, it will be ignored. Valid values are `insert` and `server`. `insert` means it is inserted from Server Load Balancer; `server` means the Server Load Balancer learns from the backend server.
* `cookie_timeout` - (Optional, Available in v1.51.0+) Cookie timeout. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "insert". Otherwise, it will be ignored. Valid value range: [1-86400] in seconds.
* `cookie` - (Optional, Available in v1.51.0+) The cookie configured on the server. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "server". Otherwise, it will be ignored. Valid value：String in line with RFC 2965, with length being 1- 200. It only contains characters such as ASCII codes, English letters and digits instead of the comma, semicolon or spacing, and it cannot start with $.
* `health_check` - (Optional, Available in v1.51.0+) Whether to enable health check. Valid values are`on` and `off`. TCP and UDP listener's HealthCheck is always on, so it will be ignore when launching TCP or UDP listener. This parameter is required  and takes effect only when ListenerSync is set to off.
* `health_check_domain` - (Optional, Available in v1.51.0+) Domain name used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and only characters such as letters, digits, ‘-‘ and ‘.’ are allowed. When it is not set or empty,  Server Load Balancer uses the private network IP address of each backend server as Domain used for health check.
* `health_check_uri` - (Optional, Available in v1.51.0+) URI used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and it must start with /. Only characters such as letters, digits, ‘-’, ‘/’, ‘.’, ‘%’, ‘?’, #’ and ‘&’ are allowed.
* `health_check_connect_port` - (Optional, Available in v1.51.0+) Port used for health check. Valid value range: [1-65535]. Default to "None" means the backend server port is used.
* `healthy_threshold` - (Optional, Available in v1.51.0+) Threshold determining the result of the health check is success. It is required when `health_check` is on. Valid value range: [1-10] in seconds. Default to 3.
* `unhealthy_threshold` - (Optional, Available in v1.51.0+) Threshold determining the result of the health check is fail. It is required when `health_check` is on. Valid value range: [1-10] in seconds. Default to 3.
* `health_check_timeout` - (Optional, Available in v1.51.0+) Maximum timeout of each health check response. It is required when `health_check` is on. Valid value range: [1-300] in seconds. Default to 5. Note: If `health_check_timeout` < `health_check_interval`, its will be replaced by `health_check_interval`.
* `health_check_interval` - (Optional, Available in v1.51.0+) Time interval of health checks. It is required when `health_check` is on. Valid value range: [1-50] in seconds. Default to 2.
* `health_check_http_code` - (Optional, Available in v1.51.0+) Regular health check HTTP status code. Multiple codes are segmented by “,”. It is required when `health_check` is on. Default to `http_2xx`.  Valid values are: `http_2xx`,  `http_3xx`, `http_4xx` and `http_5xx`.
* `listener_sync` - (Optional, Available in v1.51.0+) Indicates whether a forwarding rule inherits the settings of a health check , session persistence, and scheduling algorithm from a listener. Default to on.
* `delete_protection_validation` - (Optional, Available in 1.63.0+) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the forwarding rule.

### Timeouts

-> **NOTE:** Available in v1.163.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the forwarding rule.
* `update` - (Defaults to 1 mins) Used when update the forwarding rule.
* `delete` - (Defaults to 1 mins) Used when delete the forwarding rule.
                                                                                             
## Import

Load balancer forwarding rule can be imported using the id, e.g.

```
$ terraform import alicloud_slb_rule.example rule-abc123456
```
