---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_server_group"
sidebar_current: "docs-alicloud-resource-alb-server-group"
description: |- 
    Provides a Alicloud ALB Server Group resource.
---

# alicloud\_alb\_server\_group

Provides a ALB Server Group resource.

For information about ALB Server Group and how to use it,
see [What is Server Group](https://www.alibabacloud.com/help/doc-detail/213627.htm).

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
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
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
}
data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = alicloud_vpc.default.id
  server_group_name = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  health_check_config {
    health_check_connect_port = "46325"
    health_check_enabled      = true
    health_check_host         = "tf-testAcc.com"
    health_check_codes        = ["http_2xx", "http_3xx", "http_4xx"]
    health_check_http_version = "HTTP1.1"
    health_check_interval     = "2"
    health_check_method       = "HEAD"
    health_check_path         = "/tf-testAcc"
    health_check_protocol     = "HTTP"
    health_check_timeout      = 5
    healthy_threshold         = 3
    unhealthy_threshold       = 3
  }
  sticky_session_config {
    sticky_session_enabled = true
    cookie                 = "tf-testAcc"
    sticky_session_type    = "Server"
  }
  tags = {
    Created = "TF"
  }
  servers {
    description = var.name
    port        = 80
    server_id   = alicloud_instance.default.id
    server_ip   = alicloud_instance.default.private_ip
    server_type = "Ecs"
    weight      = 10
  }
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `health_check_config` - (Optional, ForceNew) The configuration of health checks.
* `protocol` - (Optional, Computed, ForceNew) The server protocol. Valid values: ` HTTPS`, `HTTP`.
* `resource_group_id` - (Optional) The ID of the resource group.
* `scheduler` - (Optional, Computed) The scheduling algorithm. Valid values: ` Sch`, ` Wlc`, `Wrr`.
* `server_group_name` - (Optional) The name of the resource.
* `servers` - (Optional) The backend server.
* `sticky_session_config` - (Optional, ForceNew) The configuration of the sticky session.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC that you want to access.

#### Block sticky_session_config

The sticky_session_config supports the following:

* `cookie` - the cookie that is configured on the server. **NOTE:** This parameter exists if the `StickySession`
  parameter is set to `On` and the `StickySessionType` parameter is set to `server`.
* `cookie_timeout` - The timeout period of a cookie. The timeout period of a cookie. Unit: seconds. Valid values: `1`
  to `86400`. Default value: `1000`.
* `sticky_session_enabled` - Indicates whether sticky session is enabled. Values: `true` and `false`. Default
  value: `false`.  **NOTE:** This parameter exists if the `StickySession` parameter is set to `On`.
* `sticky_session_type` - The method that is used to handle a cookie. Values: `Server` and `Insert`.

#### Block servers

The servers supports the following:

* `server_type` - The type of the server. The type of the server. Valid values: `Ecs`, `Eni` and `Eci`.
* `weight` - The weight of the server. Valid values: `0` to `100`. Default value: `100`. If the value is set to `0`, no
  requests are forwarded to the server.
* `description` - The description of the server.
* `port` - The port that is used by the server. Valid values: `1` to `65535`.
* `server_id` - The ID of the ECS instance, ENI instance or ECI instance.
* `server_ip` - The IP address of the ENI instance when it is in the inclusive ENI mode.

#### Block health_check_config

The health_check_config supports the following:

* `unhealthy_threshold` - The number of consecutive health checks that a healthy backend server must consecutively fail
  before it is declared unhealthy. In this case, the health check state is changed from success to fail. Valid
  values: `2` to `10`. Default value: `3`.
* `health_check_codes` - The status code for a successful health check. Multiple status codes can be specified as a
  list. Valid values: `http_2xx`, `http_3xx`, `http_4xx`, and `http_5xx`. Default value: `http_2xx`. **NOTE:** This
  parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
* `health_check_path` - The forwarding rule path of health checks. **NOTE:** This parameter exists if
  the `HealthCheckProtocol` parameter is set to `HTTP`.
* `healthy_threshold` - The number of health checks that an unhealthy backend server must pass consecutively before it
  is declared healthy. In this case, the health check state is changed from fail to success. Valid values: 2 to 10.
  Default value: 3.
* `health_check_http_version` - HTTP protocol version. Valid values: `HTTP1.0` and `HTTP1.1`. Default value: `HTTP1.1`
  . **NOTE:** This parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
* `health_check_interval` - The time interval between two consecutive health checks. Unit: seconds. Valid values: `1`
  to `50`. Default value: `2`.
* `health_check_method` - Health check method. Valid values: `GET` and `HEAD`. Default: `GET`. **NOTE:** This parameter
  exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
* `health_check_protocol` - Health check protocol. Valid values: `HTTP` and `TCP`.
* `health_check_timeout` - The timeout period of a health check response. If a backend Elastic Compute Service (ECS)
  instance does not send an expected response within the specified period of time, the ECS instance is considered
  unhealthy. Unit: seconds. Valid values: 1 to 300. Default value: 5. **NOTE:** If the value of the `HealthCHeckTimeout`
  parameter is smaller than that of the `HealthCheckInterval` parameter, the value of the `HealthCHeckTimeout` parameter
  is ignored and the value of the `HealthCheckInterval` parameter is regarded as the timeout period.
* `health_check_connect_port` - The port of the backend server that is used for health checks. Valid values: `0`
  to `65535`. Default value: `0`. A value of 0 indicates that a backend server port is used for health checks.
* `health_check_enabled` - Indicates whether health checks are enabled. Valid values: `true`, `false`. Default
  value: `true`.
* `health_check_host` - The domain name that is used for health checks.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Server Group.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Instance.
* `update` - (Defaults to 6 mins) Used when update the Instance.
* `delete` - (Defaults to 6 mins) Used when delete the Instance.

## Import

ALB Server Group can be imported using the id, e.g.

```
$ terraform import alicloud_alb_server_group.example <id>
```
