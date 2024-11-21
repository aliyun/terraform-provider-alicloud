---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_server_group"
sidebar_current: "docs-alicloud-resource-alb-server-group"
description: |-
  Provides a Alicloud ALB Server Group resource.
---

# alicloud_alb_server_group

Provides an ALB Server Group resource.

For information about ALB Server Group and how to use it, see [What is Server Group](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createservergroup).

-> **NOTE:** Available since v1.131.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_server_group&exampleId=b90e8c32-455e-254a-a908-adb3d68538ff301b884c&activeTab=example&spm=docs.r.alb_server_group.0.b90e8c3245&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "example" {
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
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
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
  sticky_session_config {
    sticky_session_enabled = true
    cookie                 = "tf-example"
    sticky_session_type    = "Server"
  }
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
  servers {
    description = var.name
    port        = 80
    server_id   = alicloud_instance.example.id
    server_ip   = alicloud_instance.example.private_ip
    server_type = "Ecs"
    weight      = 10
  }
  tags = {
    Created = "TF"
  }
}
```

## Argument Reference

The following arguments are supported:

* `server_group_name` - (Required) The name of the server group.
* `server_group_type` - (Optional, ForceNew, Available since v1.193.0) The type of the server group. Default value: `Instance`. Valid values:
  - `Instance`: allows you add servers by specifying Ecs, Ens, or Eci.
  - `Ip`: allows you to add servers by specifying IP addresses.
  - `Fc`: allows you to add servers by specifying functions of Function Compute.
* `protocol` - (Optional, ForceNew) The server protocol. Valid values: ` HTTP`, `HTTPS`, `gRPC`. While `server_group_type` is `Fc` this parameter will not take effect. From version 1.215.0, `protocol` can be set to `gRPC`.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC that you want to access. **NOTE:** This parameter takes effect when the `server_group_type` parameter is set to `Instance` or `Ip`.
* `scheduler` - (Optional) The scheduling algorithm. Valid values: ` Sch`, ` Wlc`, `Wrr`. **NOTE:** This parameter takes effect when the `server_group_type` parameter is set to `Instance` or `Ip`.
* `resource_group_id` - (Optional) The ID of the resource group.
* `sticky_session_config` - (Optional, Set) The configuration of session persistence. See [`sticky_session_config`](#sticky_session_config) below.
* `health_check_config` - (Required, Set) The configuration of health checks. See [`health_check_config`](#health_check_config) below.
* `servers` - (Optional, Set) The backend servers. See [`servers`](#servers) below.
* `dry_run` - (Optional, Bool) The dry run.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### `sticky_session_config`

The sticky_session_config supports the following:

* `sticky_session_enabled` - (Optional, Bool) Specifies whether to enable session persistence. Default value: `false`. Valid values: `true`, `false`. **NOTE:** This parameter takes effect when the `server_group_type` parameter is set to `Instance` or `Ip`.
* `sticky_session_type` - (Optional) The method that is used to handle a cookie. Valid values: `Server`, `Insert`.
* `cookie` - (Optional) The cookie to be configured on the server. **NOTE:** This parameter takes effect when the `sticky_session_enabled` parameter is set to `true` and the `sticky_session_type` parameter is set to `Server`.
* `cookie_timeout` - (Optional, Int) The timeout period of a cookie. Unit: seconds. Default value: `1000`. Valid values: `1` to `86400`. **NOTE:** This parameter takes effect when the `sticky_session_enabled` parameter is set to `true` and the `sticky_session_type` parameter is set to `Insert`.

### `health_check_config`

The health_check_config supports the following:

* `health_check_enabled` - (Required, Bool) Specifies whether to enable the health check feature. Valid values: `true`, `false`.
* `health_check_connect_port` - (Optional, Int) The backend port that is used for health checks. Default value: `0`. Valid values: `0` to `65535`. A value of 0 indicates that a backend server port is used for health checks.
* `health_check_host` - (Optional) The domain name that is used for health checks.
* `health_check_http_version` - (Optional) The version of the HTTP protocol. Default value: `HTTP1.1`. Valid values: `HTTP1.0` and `HTTP1.1`. **NOTE:** This parameter takes effect only when `health_check_protocol` is set to `HTTP` or `HTTPS`.
* `health_check_interval` - (Optional, Int) The interval at which health checks are performed. Unit: seconds. Default value: `2`. Valid values: `1` to `50`.
* `health_check_method` - (Optional) The HTTP method that is used for health checks. Default value: `GET`. Valid values: `GET`, `POST`, `HEAD`. **NOTE:** This parameter takes effect only when `health_check_protocol` is set to `HTTP`, `HTTPS`, or `gRPC`. From version 1.215.0, `health_check_method` can be set to `POST`.
* `health_check_path` - (Optional) The path that is used for health checks. **NOTE:** This parameter takes effect only when `health_check_protocol` is set to `HTTP` or `HTTPS`.
* `health_check_protocol` - (Optional) The protocol that is used for health checks. Valid values: `HTTP`, `HTTPS`, `TCP` and `gRPC`.
* `health_check_timeout` - (Optional, Int) The timeout period for a health check response. If a backend Elastic Compute Service (ECS) instance does not send an expected response within the specified period of time, the ECS instance is considered unhealthy. Unit: seconds. Default value: `5`. Valid values: `1` to `300`. **NOTE:** If the value of `health_check_timeout` is smaller than the value of `health_check_interval`, the value of `health_check_timeout` is ignored and the value of `health_check_interval` is used.
* `healthy_threshold` - (Optional, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. Default value: `3`. Valid values: `2` to `10`.
* `unhealthy_threshold` - (Optional, Int) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. Default value: `3`. Valid values: `2` to `10`.
* `health_check_codes` - (Optional, List) The HTTP status codes that are used to indicate whether the backend server passes the health check. Valid values:
  - If `health_check_protocol` is set to `HTTP` or `HTTPS`. Valid values: `http_2xx`, `http_3xx`, `http_4xx`, and `http_5xx`. Default value: `http_2xx`.
  - If `health_check_protocol` is set to `gRPC`. Valid values: `0` to `99`. Default value: `0`.

### `servers`

The servers supports the following:

* `server_id` - (Required) The ID of the backend server.
  - If `server_group_type` is set to `Instance`, set the parameter to the ID of an Elastic Compute Service (ECS) instance, an elastic network interface (ENI), or an elastic container instance. These backend servers are specified by Ecs, Eni, or Eci.
  - If `server_group_type` is set to `Ip`, set the parameter to an IP address specified in the server group.
  - If `server_group_type` is set to `Fc`, set the parameter to the Alibaba Cloud Resource Name (ARN) of a function specified in the server group.
* `server_type` - (Required) The type of the server. The type of the server. Valid values:
  - `Ecs`: an ECS instance.
  - `Eni`: an ENI.
  - `Eci`: an elastic container instance.
  - `Ip`(Available since v1.194.0): an IP address.
  - `Fc`(Available since v1.194.0): a function.
* `server_ip` - (Optional) The IP address of an Elastic Compute Service (ECS) instance, an elastic network interface (ENI), or an elastic container instance. **Note:** If `server_group_type` is set to `Fc`, you do not need to configure parameters, otherwise this attribute is required. If `server_group_type` is set to `Ip`, the value of this property is the same as the `server_id` value.
* `port` - (Optional, Int) The port used by the backend server. Valid values: `1` to `65535`. **Note:** This parameter is required if the `server_type` parameter is set to `Ecs`, `Eni`, `Eci`, or `Ip`. You do not need to configure this parameter if you set `server_type` to `Fc`.
* `remote_ip_enabled` - (Optional, Bool, Available since v1.194.0) Specifies whether to enable the remote IP address feature. You can specify up to 40 servers in each call. **Note:** If `server_type` is set to `Ip`, this parameter is available.
* `weight` - (Optional, Int) The weight of the server. Default value: `100`. Valid values: `0` to `100`. If the value is set to `0`, no requests are forwarded to the server. **Note:** You do not need to set this parameter if you set `server_type` to `Fc`.
* `description` - (Optional) The description of the backend server.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Server Group.
* `status` - The status of the Server Group.
* `servers` - The backend servers.
  * `status` - The status of the backend server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Server Group.
* `update` - (Defaults to 6 mins) Used when update the Server Group.
* `delete` - (Defaults to 6 mins) Used when delete the Server Group.

## Import

ALB Server Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_server_group.example <id>
```
