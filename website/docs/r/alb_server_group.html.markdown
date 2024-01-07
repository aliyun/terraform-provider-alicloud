---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_server_group"
sidebar_current: "docs-alicloud-resource-alb-server-group"
description: |-
    Provides a Alicloud ALB Server Group resource.
---

# alicloud_alb_server_group

Provides a ALB Server Group resource.

For information about ALB Server Group and how to use it,
see [What is Server Group](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createservergroup).

-> **NOTE:** Available since v1.131.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_zones" "example" {
  available_resource_creation = "Instance"
}
data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "example" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
}
data "alicloud_resource_manager_resource_groups" "example" {
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/16"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name        = var.name
  description = var.name
  vpc_id      = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  instance_name     = var.name
  image_id          = data.alicloud_images.example.images.0.id
  instance_type     = data.alicloud_instance_types.example.instance_types.0.id
  security_groups   = [alicloud_security_group.example.id]
  vswitch_id        = alicloud_vswitch.example.id
}

resource "alicloud_alb_server_group" "example" {
  protocol          = "HTTP"
  vpc_id            = alicloud_vpc.example.id
  server_group_name = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.example.groups.0.id
  health_check_config {
    health_check_connect_port = "46325"
    health_check_enabled      = true
    health_check_host         = "tf-example.com"
    health_check_codes        = ["http_2xx", "http_3xx", "http_4xx"]
    health_check_http_version = "HTTP1.1"
    health_check_interval     = "2"
    health_check_method       = "HEAD"
    health_check_path         = "/tf-example"
    health_check_protocol     = "HTTP"
    health_check_timeout      = 5
    healthy_threshold         = 3
    unhealthy_threshold       = 3
  }
  sticky_session_config {
    sticky_session_enabled = true
    cookie                 = "tf-example"
    sticky_session_type    = "Server"
  }
  tags = {
    Created = "TF"
  }
  servers {
    description = var.name
    port        = 80
    server_id   = alicloud_instance.example.id
    server_ip   = alicloud_instance.example.private_ip
    server_type = "Ecs"
    weight      = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `health_check_config` - (Optional) The configuration of health checks. See [`health_check_config`](#health_check_config) below for details.
* `protocol` - (Optional, ForceNew) The server protocol. Valid values: ` HTTPS`, `HTTP`. While `server_group_type` is `Fc` this parameter will not take effect.
* `resource_group_id` - (Optional) The ID of the resource group.
* `scheduler` - (Optional) The scheduling algorithm. Valid values: ` Sch`, ` Wlc`, `Wrr`.
* `server_group_name` - (Optional) The name of the resource.
* `servers` - (Optional) The backend server. See [`servers`](#servers) below for details.
* `sticky_session_config` - (Optional, ForceNew) The configuration of the sticky session. See [`sticky_session_config`](#sticky_session_config) below for details.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC that you want to access. While `server_group_type` is `Fc` this parameter will not take effect.
* `server_group_type` - (Optional, ForceNew, Available in v1.193.0+) The type of the server group. Valid values:
  - `Instance` (default): allows you add servers by specifying Ecs, Ens, or Eci.
  - `Ip`: allows you to add servers by specifying IP addresses.
  - `Fc`: allows you to add servers by specifying functions of Function Compute.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### `sticky_session_config`

The sticky_session_config supports the following:

* `cookie` - (Optional) the cookie that is configured on the server. **NOTE:** This parameter exists if the `StickySession`
  parameter is set to `On` and the `StickySessionType` parameter is set to `server`.
* `cookie_timeout` - (Optional) The timeout period of a cookie. The timeout period of a cookie. Unit: seconds. Valid values: `1`
  to `86400`. Default value: `1000`.
* `sticky_session_enabled` - (Optional) Indicates whether sticky session is enabled. Values: `true` and `false`. Default
  value: `false`.  **NOTE:** This parameter exists if the `StickySession` parameter is set to `On`.
* `sticky_session_type` - (Optional) The method that is used to handle a cookie. Values: `Server` and `Insert`.

### `servers`

The servers supports the following:

* `server_type` - (Required) The type of the server. The type of the server. Valid values: 
  - Ecs: an ECS instance.
  - Eni: an ENI.
  - Eci: an elastic container instance.
  - Ip(Available in v1.194.0+): an IP address.
  - fc(Available in v1.194.0+): a function.
* `weight` - (Optional) The weight of the server. Valid values: `0` to `100`. Default value: `100`. If the value is set to `0`, no
  requests are forwarded to the server. **Note:** You do not need to set this parameter if you set `server_type` to `Fc`.
* `description` - (Optional) The description of the server.
* `port` - (Optional) The port that is used by the server. Valid values: `1` to `65535`. **Note:** This parameter is required if the `server_type` parameter is set to `Ecs`, `Eni`, `Eci`, or `Ip`. You do not need to configure this parameter if you set `server_type` to `Fc`.
* `server_id` - (Required) The ID of the backend server.
  - If `server_group_type` is set to `Instance`, set the parameter to the ID of an Elastic Compute Service (ECS) instance, an elastic network interface (ENI), or an elastic container instance. These backend servers are specified by Ecs, Eni, or Eci.
  - If `server_group_type` is set to `Ip`, set the parameter to an IP address specified in the server group.
  - If `server_group_type` is set to `Fc`, set the parameter to the Alibaba Cloud Resource Name (ARN) of a function specified in the server group.
* `server_ip` - (Optional) The IP address of an Elastic Compute Service (ECS) instance, an elastic network interface (ENI), or an elastic container instance. **Note:** If `server_group_type` is set to `Fc`, you do not need to configure parameters, otherwise this attribute is required. If `server_group_type` is set to `Ip`, the value of this property is the same as the `server_id` value.
* `remote_ip_enabled` - (Optional, Available in v1.194.0+) Specifies whether to enable the remote IP address feature. You can specify up to 40 servers in each call. **Note:** If `server_type` is set to `Ip`, this parameter is available.
* `status` - (Optional) The status of the backend server. Valid values:
  - `Adding`: The backend server is being added.
  - `Available`: The backend server is added.
  - `Configuring`: The backend server is being configured.
  - `Removing`: The backend server is being removed.

### `health_check_config`

The health_check_config supports the following:

* `unhealthy_threshold` - (Optional) The number of consecutive health checks that a healthy backend server must consecutively fail before it is declared unhealthy. In this case, the health check state is changed from success to fail. Valid values: `2` to `10`. Default value: `3`.
* `health_check_codes` - (Optional) The status code for a successful health check.  Multiple status codes can be specified as a list. Valid values: `http_2xx`, `http_3xx`, `http_4xx`, and `http_5xx`. Default value: `http_2xx`. **NOTE:** This
  parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
* `health_check_path` - (Optional) The forwarding rule path of health checks. **NOTE:** This parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
* `healthy_threshold` - (Optional) The number of health checks that an unhealthy backend server must pass consecutively before it is declared healthy. In this case, the health check state is changed from fail to success. Valid values: 2 to 10. Default value: 3.
* `health_check_http_version` - (Optional) HTTP protocol version. Valid values: `HTTP1.0` and `HTTP1.1`. Default value: `HTTP1.1`. **NOTE:** This parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
* `health_check_interval` - (Optional) The time interval between two consecutive health checks. Unit: seconds. Valid values: `1` to `50`. Default value: `2`.
* `health_check_method` - (Optional) Health check method. Valid values: `GET` and `HEAD`. Default: `GET`. **NOTE:** This parameter exists if the `HealthCheckProtocol` parameter is set to `HTTP`.
* `health_check_protocol` - (Optional) Health check protocol. Valid values: `HTTP` and `TCP`, `HTTPS`.
* `health_check_timeout` - (Optional) The timeout period of a health check response. If a backend Elastic Compute Service (ECS) instance does not send an expected response within the specified period of time, the ECS instance is considered unhealthy. Unit: seconds. Valid values: 1 to 300. Default value: 5. **NOTE:** If the value of the `HealthCHeckTimeout` parameter is smaller than that of the `HealthCheckInterval` parameter, the value of the `HealthCHeckTimeout` parameter is ignored and the value of the `HealthCheckInterval` parameter is regarded as the timeout period.
* `health_check_connect_port` - (Optional) The port of the backend server that is used for health checks. Valid values: `0` to `65535`. Default value: `0`. A value of 0 indicates that a backend server port is used for health checks.
* `health_check_enabled` - (Optional) Indicates whether health checks are enabled. Valid values: `true`, `false`. Default value: `true`.
* `health_check_host` - (Optional) The domain name that is used for health checks.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Server Group.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Instance.
* `update` - (Defaults to 6 mins) Used when update the Instance.
* `delete` - (Defaults to 6 mins) Used when delete the Instance.

## Import

ALB Server Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_server_group.example <id>
```
