---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_rule"
sidebar_current: "docs-alicloud-resource-slb-rule"
description: |-
  Provides a Alicloud Load Balancer Forwarding Rule Resource and add it to one Listener.
---

# alicloud_slb_rule

Provides a Lindorm Instance resource.

For information about Load Balancer Forwarding Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/en/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-dir-forwarding-rules).

-> **NOTE:** Available since v1.6.0.

A forwarding rule is configured in `HTTP`/`HTTPS` listener and it used to listen a list of backend servers which in one specified virtual backend server group.
You can add forwarding rules to a listener to forward requests based on the domain names or the URL in the request.

-> **NOTE:** One virtual backend server group can be attached in multiple forwarding rules.

-> **NOTE:** At least one "Domain" or "Url" must be specified when creating a new rule.

-> **NOTE:** Having the same 'Domain' and 'Url' rule can not be created repeatedly in the one listener.

-> **NOTE:** Rule only be created in the `HTTP` or `HTTPS` listener.

-> **NOTE:** Only rule's virtual server group can be modified.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_rule&exampleId=3924f6b8-5e44-cfa9-6e0b-653b370141419fb6fd66&activeTab=example&spm=docs.r.slb_rule.0.3924f6b85e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "slb_rule_name" {
  default = "terraform-example"
}

data "alicloud_zones" "rule" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "rule" {
  availability_zone = data.alicloud_zones.rule.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "rule" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "rule" {
  vpc_name   = var.slb_rule_name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "rule" {
  vpc_id       = alicloud_vpc.rule.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.rule.zones[0].id
  vswitch_name = var.slb_rule_name
}

resource "alicloud_security_group" "rule" {
  name   = var.slb_rule_name
  vpc_id = alicloud_vpc.rule.id
}

resource "alicloud_instance" "rule" {
  image_id                   = data.alicloud_images.rule.images[0].id
  instance_type              = data.alicloud_instance_types.rule.instance_types[0].id
  security_groups            = alicloud_security_group.rule.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.rule.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.rule.id
  instance_name              = var.slb_rule_name
}

resource "alicloud_slb_load_balancer" "rule" {
  load_balancer_name   = var.slb_rule_name
  vswitch_id           = alicloud_vswitch.rule.id
  instance_charge_type = "PayByCLCU"
}

resource "alicloud_slb_listener" "rule" {
  load_balancer_id          = alicloud_slb_load_balancer.rule.id
  backend_port              = 22
  frontend_port             = 22
  protocol                  = "http"
  bandwidth                 = 5
  health_check_connect_port = "20"
}

resource "alicloud_slb_server_group" "rule" {
  load_balancer_id = alicloud_slb_load_balancer.rule.id
  name             = var.slb_rule_name
}

resource "alicloud_slb_rule" "rule" {
  load_balancer_id          = alicloud_slb_load_balancer.rule.id
  frontend_port             = alicloud_slb_listener.rule.frontend_port
  name                      = var.slb_rule_name
  domain                    = "*.aliyun.com"
  url                       = "/image"
  server_group_id           = alicloud_slb_server_group.rule.id
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
* `frontend_port` - (Required, Int, ForceNew) The listener frontend port which is used to launch the new forwarding rule. Valid values: [1-65535].
* `domain` - (Optional, ForceNew) Domain name of the forwarding rule. It can contain letters a-z, numbers 0-9, hyphens (-), and periods (.),
and wildcard characters. The following two domain name formats are supported:
  - Standard domain name: www.test.com
  - Wildcard domain name: *.test.com. wildcard (*) must be the first character in the format of (*.)
* `url` - (Optional, ForceNew) Domain of the forwarding rule. It must be 2-80 characters in length. Only letters a-z, numbers 0-9, and characters '-' '/' '?' '%' '#' and '&' are allowed. URLs must be started with the character '/', but cannot be '/' alone.
* `server_group_id` - (Required) ID of a virtual server group that will be forwarded.
* `scheduler` - (Optional, Available since v1.51.0) Scheduling algorithm. Valid values: `wrr`, `rr` and `wlc`. Default value: `wrr`. **NOTE:** `scheduler` is required and takes effect only when `listener_sync` is set to `off`.
* `sticky_session` - (Optional, Available since v1.51.0) Whether to enable session persistence. Valid values: `on` and `off`. Default value: `off`. **NOTE:** `sticky_session` is required and takes effect only when `listener_sync` is set to `off`.
* `sticky_session_type` - (Optional, Available since v1.51.0) Mode for handling the cookie. If `sticky_session` is `on`, it is mandatory. Otherwise, it will be ignored. Valid values: `insert` and `server`. `insert` means it is inserted from Server Load Balancer; `server` means the Server Load Balancer learns from the backend server.
* `cookie_timeout` - (Optional, Int, Available since v1.51.0) Cookie timeout. It is mandatory when `sticky_session` is `on` and `sticky_session_type` is `insert`. Otherwise, it will be ignored. Valid values: [1-86400] in seconds.
* `cookie` - (Optional, Available since v1.51.0) The cookie configured on the server. It is mandatory when `sticky_session` is `on` and `sticky_session_type` is `server`. Otherwise, it will be ignored. Valid value：String in line with RFC 2965, with length being `1` - `200`. It only contains characters such as ASCII codes, English letters and digits instead of the comma, semicolon or spacing, and it cannot start with $.
* `health_check` - (Optional, Available since v1.51.0) Whether to enable health check. Valid values: `on` and `off`. `TCP` and `UDP` listener's `health_check` is always `on`, so it will be ignore when launching `TCP` or `UDP` listener. **NOTE:** `health_check` is required and takes effect only when `listener_sync` is set to `off`.
* `health_check_domain` - (Optional, Available since v1.51.0) Domain name used for health check. When it used to launch TCP listener, `health_check_type` must be `http`. Its length is limited to 1-80 and only characters such as letters, digits, ‘-‘ and ‘.’ are allowed. When it is not set or empty, Server Load Balancer uses the private network IP address of each backend server as Domain used for health check.
* `health_check_uri` - (Optional, Available since v1.51.0) URI used for health check. When it used to launch TCP listener, `health_check_type` must be `http`. Its length is limited to 1-80 and it must start with /. Only characters such as letters, digits, ‘-’, ‘/’, ‘.’, ‘%’, ‘?’, #’ and ‘&’ are allowed.
* `health_check_connect_port` - (Optional, Int, Available since v1.51.0) Port used for health check. Valid values: [1-65535]. Default value: `None` means the backend server port is used.
* `healthy_threshold` - (Optional, Available since v1.51.0) Threshold determining the result of the health check is success. It is required when `health_check` is `on`. Valid values: [1-10] in seconds. Default value: `3`.
* `unhealthy_threshold` - (Optional, Int, Available since v1.51.0) Threshold determining the result of the health check is fail. It is required when `health_check` is `on`. Valid values: [1-10] in seconds. Default value: `3`.
* `health_check_timeout` - (Optional, Int, Available since v1.51.0) Maximum timeout of each health check response. It is required when `health_check` is `on`. Valid values: [1-300] in seconds. Default value: `5`. Note: If `health_check_timeout` < `health_check_interval`, its will be replaced by `health_check_interval`.
* `health_check_interval` - (Optional, Int, Int, Available since v1.51.0) Time interval of health checks. It is required when `health_check` is `on`. Valid values: [1-50] in seconds. Default value: `2`.
* `health_check_http_code` - (Optional, Available since v1.51.0) Regular health check HTTP status code. Multiple codes are segmented by “,”. It is required when `health_check` is `on`. Default value: `http_2xx`. Valid values: `http_2xx`, `http_3xx`, `http_4xx` and `http_5xx`.
* `listener_sync` - (Optional, Available since v1.51.0) Indicates whether a forwarding rule inherits the settings of a health check , session persistence, and scheduling algorithm from a listener. Default value: `on`. Valid values: `on` and `off`.
* `delete_protection_validation` - (Optional, Bool, Available since v1.63.0) Checking DeleteProtection of SLB instance before deleting. If `true`, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule.

## Timeouts

-> **NOTE:** Available since v1.163.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the forwarding rule.
* `update` - (Defaults to 1 mins) Used when update the forwarding rule.
* `delete` - (Defaults to 1 mins) Used when delete the forwarding rule.
                                                                                             
## Import

Load balancer forwarding rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_rule.example <id>
```
