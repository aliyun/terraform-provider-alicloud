---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_server_group"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Server Group resource.
---

# alicloud_alb_server_group

Provides a Application Load Balancer (ALB) Server Group resource.



For information about Application Load Balancer (ALB) Server Group and how to use it, see [What is Server Group](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createservergroup).

-> **NOTE:** Available since v1.131.0.

## Example Usage

Basic Usage

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
* `connection_drain_config` - (Optional, List, Available since v1.242.0) Elegant interrupt configuration. See [`connection_drain_config`](#connection_drain_config) below.
* `cross_zone_enabled` - (Optional, Computed, Available since v1.242.0) Indicates whether cross-zone load balancing is enabled for the server group. Valid values:

  *   `true` (default)

  *   `false`

-> **NOTE:**  

  *   Basic ALB instances do not support server groups that have cross-zone load balancing disabled. Only Standard and WAF-enabled ALB instances support server groups that have cross-zone load balancing.

  *   Cross-zone load balancing can be disabled for server groups of the server and IP type, but not for server groups of the Function Compute type.

  *   When cross-zone load balancing is disabled, session persistence cannot be enabled.
* `health_check_config` - (Required, List) The configuration of health checks See [`health_check_config`](#health_check_config) below.
* `health_check_template_id` - (Optional, Available since v1.242.0) The template ID.
* `protocol` - (Optional, ForceNew, Computed) The backend protocol. Valid values:

  *   `HTTP`: allows you to associate an HTTPS, HTTP, or QUIC listener with the server group. This is the default value.

  *   `HTTPS`: allows you to associate HTTPS listeners with backend servers.

  *   `gRPC`: allows you to associate an HTTPS or QUIC listener with the server group.

-> **NOTE:**   You do not need to specify a backend protocol if you set `ServerGroupType` to `Fc`.

* `resource_group_id` - (Optional, Computed) The ID of the resource group to which you want to transfer the cloud resource.

-> **NOTE:**   You can use resource groups to manage resources within your Alibaba Cloud account by group. This helps you resolve issues such as resource grouping and permission management for your Alibaba Cloud account. For more information, see [What is resource management?](https://www.alibabacloud.com/help/en/doc-detail/94475.html)

* `scheduler` - (Optional, Computed) The scheduling algorithm. Valid values:

  *   `Wrr` (default): The weighted round-robin algorithm is used. Backend servers that have higher weights receive more requests than those that have lower weights.

  *   `Wlc`: The weighted least connections algorithm is used. Requests are distributed based on the weights and the number of connections to backend servers. If two backend servers have the same weight, the backend server that has fewer connections is expected to receive more requests.

  *   `Sch`: The consistent hashing algorithm is used. Requests from the same source IP address are distributed to the same backend server.

-> **NOTE:**  This parameter takes effect when the `ServerGroupType` parameter is set to `Instance` or `Ip`.

* `server_group_name` - (Required) The name of the server group. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-). The name must start with a letter.
* `server_group_type` - (Optional, ForceNew, Computed, Available since v1.193.0) The type of server group. Valid values:

  - `Instance` (default): allows you to add servers by specifying `Ecs`, `Eni`, or `Eci`.
  - `Ip`: allows you to add servers by specifying IP addresses.
  - `Fc`: allows you to add servers by specifying functions of Function Compute.
* `servers` - (Optional, Set) List of servers. See [`servers`](#servers) below.
* `slow_start_config` - (Optional, List, Available since v1.242.0) Slow start configuration. See [`slow_start_config`](#slow_start_config) below.
* `sticky_session_config` - (Optional, List) The configuration of the sticky session See [`sticky_session_config`](#sticky_session_config) below.
* `tags` - (Optional, Map) The tag of the resource
* `uch_config` - (Optional, List, Available since v1.242.0) Url consistency hash parameter configuration See [`uch_config`](#uch_config) below.
* `upstream_keepalive_enabled` - (Optional, Available since v1.242.0) Specifies whether to enable persistent TCP connections.
* `vpc_id` - (Optional, ForceNew) The ID of the virtual private cloud (VPC). You can add only servers that are deployed in the specified VPC to the server group.

-> **NOTE:**   This parameter takes effect when the `ServerGroupType` parameter is set to `Instance` or `Ip`.


### `connection_drain_config`

The connection_drain_config supports the following:
* `connection_drain_enabled` - (Optional, Computed, Available since v1.242.0) Specifies whether to enable connection draining. Valid values:

  - `true`
  - `false` (default)
* `connection_drain_timeout` - (Optional, Computed, Int, Available since v1.242.0) The timeout period of connection draining.

  Valid values: `0` to `900`.

  Default value: `300`.

### `health_check_config`

The health_check_config supports the following:
* `health_check_codes` - (Optional, Computed, List) The status code for a successful health check
* `health_check_connect_port` - (Optional, Int) The backend port that is used for health checks.

  Valid values: `0` to `65535`.

  If you set the value to `0`, the backend port is used for health checks.

-> **NOTE:**   This parameter takes effect only if you set `HealthCheckEnabled` to `true`.

* `health_check_enabled` - (Required) Specifies whether to enable the health check feature. Valid values:

  *   `true`

  *   `false`

-> **NOTE:**   If the `ServerGroupType` parameter is set to `Instance` or `Ip`, the health check feature is enabled by default. If the `ServerGroupType` parameter is set to `Fc`, the health check feature is disabled by default.

* `health_check_host` - (Optional, Computed) The domain name that is used for health checks.

  *   **Backend Server Internal IP** (default): Use the internal IP address of backend servers as the health check domain name.

  *   **Custom Domain Name**: Enter a domain name.

    *   The domain name must be 1 to 80 characters in length.
    *   The domain name can contain lowercase letters, digits, hyphens (-), and periods (.).
    *   The domain name must contain at least one period (.) but cannot start or end with a period (.).
    *   The rightmost domain label of the domain name can contain only letters, and cannot contain digits or hyphens (-).
    *   The domain name cannot start or end with a hyphen (-).

-> **NOTE:**   This parameter takes effect only if `HealthCheckProtocol` is set to `HTTP`, `HTTPS`, or `gRPC`.

* `health_check_http_version` - (Optional, Computed) The HTTP version that is used for health checks. Valid values:

  *   **HTTP1.0**

  *   **HTTP1.1**

-> **NOTE:**   This parameter takes effect only if you set `HealthCheckEnabled` to true and `HealthCheckProtocol` to `HTTP` or `HTTPS`.

* `health_check_interval` - (Optional, Computed, Int) The interval at which health checks are performed. Unit: seconds.

  Valid values: `1` to `50`.

-> **NOTE:**   This parameter takes effect only if you set `HealthCheckEnabled` to `true`.

* `health_check_method` - (Optional, Computed) The HTTP method that is used for health checks. Valid values:

  *   `GET`: If the length of a response exceeds 8 KB, the response is truncated. However, the health check result is not affected.

  *   `POST`: gRPC health checks use the POST method by default.

  *   `HEAD`: HTTP and HTTPS health checks use the HEAD method by default.

-> **NOTE:**   This parameter takes effect only if you set `HealthCheckEnabled` to true and `HealthCheckProtocol` to `HTTP`, `HTTPS`, or `gRPC`.

* `health_check_path` - (Optional, Computed) The URL that is used for health checks.

  The URL must be 1 to 80 characters in length, and can contain letters, digits, and the following special characters: `- / . % ? # & =`. It can also contain the following extended characters: `_ ; ~ ! ( ) * [ ] @ $ ^ : ' , +`. The URL must start with a forward slash (`/`).

-> **NOTE:**   This parameter takes effect only if you set `HealthCheckEnabled` to `true` and `HealthCheckProtocol` to `HTTP` or `HTTPS`.

* `health_check_protocol` - (Optional, Computed) The protocol that is used for health checks. Valid values:

  - `HTTP`: HTTP health checks simulate browser behaviors by sending HEAD or GET requests to probe the availability of backend servers.
  - `HTTPS`: HTTPS health checks simulate browser behaviors by sending HEAD or GET requests to probe the availability of backend servers. HTTPS provides higher security than HTTP because HTTPS supports data encryption.
  - `TCP`: TCP health checks send TCP SYN packets to a backend server to probe the availability of backend servers.
  - `gRPC`: gRPC health checks send POST or GET requests to a backend server to check whether the backend server is healthy.
* `health_check_timeout` - (Optional, Computed, Int) The timeout period of a health check response. If a backend ECS instance does not respond within the specified timeout period, the ECS instance fails the health check. Unit: seconds.

  Valid values: `1` to `300`.

-> **NOTE:**   This parameter takes effect only if you set `HealthCheckEnabled` to `true`.

* `healthy_threshold` - (Optional, Computed, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. In this case, the health check status of the backend server changes from `fail` to `success`.

  Valid values: `2` to `10`.

  Default value: `3`.
* `unhealthy_threshold` - (Optional, Computed, Int) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. In this case, the health check status of the backend server changes from `success` to `fail`.

  Valid values: `2` to `10`.

  Default value: `3`.

### `servers`

The servers supports the following:
* `description` - (Optional) The description of the backend server. The description must be 2 to 256 characters in length, and cannot start with http:// or https://.
* `port` - (Optional, Int) The port that is used by the backend server. Valid values: `1` to `65535`. You can specify at most 200 servers in each call.

-> **NOTE:**   This parameter is required if you set `ServerType` to `Ecs`, `Eni`, `Eci`, or `Ip`. You do not need to set this parameter if `ServerType` is set to `Fc`.

* `remote_ip_enabled` - (Optional, Available since v1.194.0) Specifies whether to enable the remote IP feature. You can specify at most 200 servers in each call. Default values:

  *   `true`: enables the feature.

  *   `false`: disables the feature.

-> **NOTE:**   This parameter takes effect only when `ServerType` is set to `Ip`.

* `server_id` - (Required) The ID of the backend server. You can specify at most 200 servers in each call.

  *   If the server group is of the `Instance` type, set ServerId to the ID of a resource of the `Ecs`, `Eni`, or `Eci` type.

  *   If the server group is of the `Ip` type, set ServerId to IP addresses.

-> **NOTE:**   You cannot perform this operation on a server group of the Function Compute type. You can call the [ListServerGroups](https://www.alibabacloud.com/help/en/doc-detail/213627.html) operation to query the type of server groups.

* `server_ip` - (Optional) The IP address of the backend server. You can specify at most 200 servers in each call.

-> **NOTE:**   You do not need to set this parameter if you set `ServerType` to `Fc`.

* `server_type` - (Required, Available since v1.194.0) The type of the backend server. You can specify at most 200 servers in each call. Default values:

  - `Ecs`: Elastic Compute Service (ECS) instance
  - `Eni`: elastic network interface (ENI)
  - `Eci`: elastic container instance
  - `Ip`: IP address
  - `Fc`: Function Compute
* `weight` - (Optional, Computed, Int) The weight of the backend server. Valid values: `0` to `100`. Default value: `100`. If the value is set to `0`, no requests are forwarded to the server. You can specify at most 200 servers in each call.

-> **NOTE:**   You do not need to set this parameter if you set `ServerType` to `Fc`.


### `slow_start_config`

The slow_start_config supports the following:
* `slow_start_duration` - (Optional, Computed, Int, Available since v1.242.0) The duration of a slow start.

  Valid values: 30 to 900.

  Default value: 30.
* `slow_start_enabled` - (Optional, Computed, Available since v1.242.0) Indicates whether slow starts are enabled. Valid values:

  - `true`
  - `false`

### `sticky_session_config`

The sticky_session_config supports the following:
* `cookie` - (Optional, Computed) The cookie to be configured on the server.

  The cookie must be 1 to 200 characters in length and can contain only ASCII characters and digits. It cannot contain commas (,), semicolons (;), or space characters. It cannot start with a dollar sign ($).

-> **NOTE:**  This parameter takes effect when the `StickySessionEnabled` parameter is set to `true` and the `StickySessionType` parameter is set to `Server`.

* `cookie_timeout` - (Optional, Computed, Int) The maximum amount of time to wait before the session cookie expires. Unit: seconds.

  Valid values: `1` to `86400`.

  Default value: `1000`.

-> **NOTE:**   This parameter takes effect only when `StickySessionEnabled` is set to `true` and `StickySessionType` is set to `Insert`.

* `sticky_session_enabled` - (Optional) Specifies whether to enable session persistence. Valid values:

  *   `true`

  *   `false`

-> **NOTE:**   This parameter takes effect when the `ServerGroupType` parameter is set to `Instance` or `Ip`.

* `sticky_session_type` - (Optional, Computed) The method that is used to handle a cookie. Valid values:

  *   `Insert`: inserts a cookie.

  ALB inserts a cookie (SERVERID) into the first HTTP or HTTPS response packet that is sent to a client. The next request from the client contains this cookie and the listener forwards this request to the recorded backend server.

  *   `Server`: rewrites a cookie.

  When ALB detects a user-defined cookie, it overwrites the original cookie with the user-defined cookie. Subsequent requests to ALB carry this user-defined cookie, and ALB determines the destination servers of the requests based on the cookies.

-> **NOTE:**  This parameter takes effect when the `StickySessionEnabled` parameter is set to `true` for the server group.


### `uch_config`

The uch_config supports the following:
* `type` - (Optional, Available since v1.242.0) The parameter type. Only QueryString can be filled.
* `value` - (Optional, Available since v1.242.0) Consistency hash parameter value

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `servers` - List of servers.
  * `server_group_id` - The ID of the server group.
  * `status` - The addition status of the backend server. Value:
  - `Adding`: Adding.
  - `Available`: normal availability.
  - `Configuring`: The configuration is under configuration.
  - `Removing`: Removing.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Server Group.
* `delete` - (Defaults to 5 mins) Used when delete the Server Group.
* `update` - (Defaults to 5 mins) Used when update the Server Group.

## Import

Application Load Balancer (ALB) Server Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_server_group.example <id>
```