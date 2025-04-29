---
subcategory: "GWLB"
layout: "alicloud"
page_title: "Alicloud: alicloud_gwlb_server_group"
description: |-
  Provides a Alicloud GWLB Server Group resource.
---

# alicloud_gwlb_server_group

Provides a GWLB Server Group resource.



For information about GWLB Server Group and how to use it, see [What is Server Group](https://www.alibabacloud.com/help/en/slb/gateway-based-load-balancing-gwlb/developer-reference/api-gwlb-2024-04-15-createservergroup).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gwlb_server_group&exampleId=b49a0f92-a93f-a2f7-1e84-f5a882211cc13bb17352&activeTab=example&spm=docs.r.gwlb_server_group.0.b49a0f92a9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultEaxcvb" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "tf-gwlb-vpc"
}

resource "alicloud_vswitch" "defaultc3uVID" {
  vpc_id       = alicloud_vpc.defaultEaxcvb.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "tf-example-vsw1"
}

resource "alicloud_security_group" "default7NNxRl" {
  description         = "sg"
  security_group_name = "sg_name"
  vpc_id              = alicloud_vpc.defaultEaxcvb.id
  security_group_type = "normal"
}

resource "alicloud_instance" "defaultH6McvC" {
  vswitch_id = alicloud_vswitch.defaultc3uVID.id
  image_id   = "aliyun_2_1903_x64_20G_alibase_20231221.vhd"

  instance_type        = "ecs.g6.large"
  system_disk_category = "cloud_efficiency"

  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  instance_name              = format("%s4", var.name)
  description                = "tf-example-ecs"
  security_groups            = [alicloud_security_group.default7NNxRl.id]

  availability_zone    = alicloud_vswitch.defaultc3uVID.zone_id
  instance_charge_type = "PostPaid"
}

resource "alicloud_gwlb_server_group" "default" {
  dry_run = "false"
  servers {
    server_id   = alicloud_instance.defaultH6McvC.id
    server_type = "Ecs"
  }

  scheduler = "5TCH"
  protocol  = "GENEVE"
  connection_drain_config {
    connection_drain_enabled = "true"
    connection_drain_timeout = "1"
  }

  vpc_id            = alicloud_vpc.defaultEaxcvb.id
  server_group_type = "Instance"
  server_group_name = var.name
  health_check_config {
    health_check_connect_port    = "80"
    health_check_enabled         = "true"
    health_check_protocol        = "HTTP"
    health_check_connect_timeout = "5"
    health_check_domain          = "www.domain.com"
    health_check_http_code = [
      "http_2xx",
      "http_3xx",
      "http_4xx"
    ]
    health_check_interval = "10"
    health_check_path     = "/health-check"
    healthy_threshold     = "2"
    unhealthy_threshold   = "2"
  }

  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```

## Argument Reference

The following arguments are supported:
* `connection_drain_config` - (Optional, Computed, List) Connected graceful interrupt configuration. See [`connection_drain_config`](#connection_drain_config) below.
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:

  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error code is returned. If the request passes the dry run, the `DryRunOperation` error code is returned.
  - `false` (default): performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `health_check_config` - (Optional, Computed, List) Health check configurations. See [`health_check_config`](#health_check_config) below.
* `protocol` - (Optional, ForceNew, Computed) The backend protocol. Valid values:

  - `GENEVE`(default)
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `scheduler` - (Optional, Computed) The scheduling algorithm. Valid values:

  - `5TCH` (default): specifies consistent hashing that is based on the following factors: source IP address, destination IP address, source port, protocol, and destination port. Requests that contain the same information based on the preceding factors are forwarded to the same backend server.
  - `3TCH`: specifies consistent hashing that is based on the following factors: source IP address, destination IP address, and protocol. Requests that contain the same information based on the preceding factors are forwarded to the same backend server.
  - `2TCH`: specifies consistent hashing that is based on the following factors: source IP address and destination IP address. Requests that contain the same information based on the preceding factors are forwarded to the same backend server.
* `server_group_name` - (Optional) The server group name.

  The name must be 2 to 128 characters in length, and can contain digits, periods (.), underscores (\_), and hyphens (-). It must start with a letter.
* `server_group_type` - (Optional, ForceNew, Computed) The type of server group. Valid values:

  - `Instance` (default): allows you to specify servers of the `Ecs`, `Eni`, or `Eci` type.
  - `Ip`: allows you to add servers of by specifying IP addresses.
* `servers` - (Optional, Set) The backend servers that you want to remove.

-> **NOTE:**  You can remove at most 200 backend servers in each call.
 See [`servers`](#servers) below.
* `tags` - (Optional, Map) The tag keys.

  You can specify at most 20 tags in each call.
* `vpc_id` - (Required, ForceNew) The VPC ID.

-> **NOTE:**  If `ServerGroupType` is set to `Instance`, only servers in the specified VPC can be added to the server group.


### `connection_drain_config`

The connection_drain_config supports the following:
* `connection_drain_enabled` - (Optional, Computed, Available since v1.236.0) Indicates whether connection draining is enabled. Valid values:

  - `true`
  - `false`
* `connection_drain_timeout` - (Optional, Computed, Int, Available since v1.236.0) The timeout period of connection draining.

  Unit: seconds

  Valid values: `1` to `3600`.

  Default value: `300`.

### `health_check_config`

The health_check_config supports the following:
* `health_check_connect_port` - (Optional, Computed, Int, Available since v1.236.0) The backend server port that is used for health checks.

  Valid values: `1` to `65535`.

  Default value: `80`.
* `health_check_connect_timeout` - (Optional, Computed, Int, Available since v1.236.0) The maximum timeout period of a health check response.

  Unit: seconds

  Valid values: `1` to `300`.

  Default value: `5`.
* `health_check_domain` - (Optional, Computed, Available since v1.236.0) The domain name that you want to use for health checks. Valid values:

  *   **$SERVER_IP** (default): the private IP address of a backend server.

  *   `domain`: a domain name. The domain name must be 1 to 80 characters in length, and can contain letters, digits, hyphens (-), and periods (.).

-> **NOTE:**  This parameter takes effect only if you set `HealthCheckProtocol` to `HTTP`.

* `health_check_enabled` - (Optional, Computed, Available since v1.236.0) Specifies whether to enable the health check feature. Valid values:

  - `true` (default)
  - `false`
* `health_check_http_code` - (Optional, Set, Available since v1.236.0) The HTTP status codes that the system returns for health checks.
* `health_check_interval` - (Optional, Computed, Int, Available since v1.236.0) The interval at which health checks are performed.

  Unit: seconds

  Valid values: `1` to `50`.

  Default value: `10`.
* `health_check_path` - (Optional, Computed, Available since v1.236.0) The URL that is used for health checks.

  The URL must be 1 to 80 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), percent signs (%), question marks (?), number signs (#), and ampersands (&). The URL can also contain the following extended characters: \_ ; ~ ! ( ) \* \[ ] @ $ ^ : ' , + =

  The URL must start with a forward slash (/).

-> **NOTE:**  This parameter takes effect only if you set `HealthCheckProtocol` to `HTTP`.

* `health_check_protocol` - (Optional, Computed, Available since v1.236.0) The protocol that is used for health checks. Valid values:

  - `TCP`: TCP health checks send TCP SYN packets to a backend server to check whether the port of the backend server is reachable.
  - `HTTP`: HTTP health checks simulate a process that uses a web browser to access resources by sending HEAD or GET requests to an instance. These requests are used to check whether the instance is healthy.
* `healthy_threshold` - (Optional, Computed, Int, Available since v1.236.0) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. In this case, the health status changes from `fail` to `success`.

  Valid values: `2` to `10`.

  Default value: `2`.
* `unhealthy_threshold` - (Optional, Computed, Int, Available since v1.236.0) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. In this case, the health status changes from `success` to `fail`.

  Valid values: `2` to `10`.

  Default value: `2`.

### `servers`

The servers supports the following:
* `server_id` - (Required) The backend server ID.

  - If the server group is of the `Instance` type, set this parameter to the IDs of servers of the `Ecs`, `Eni`, or `Eci` type.
  - If the server group is of the `Ip` type, set ServerId to IP addresses.
* `server_ip` - (Optional) The IP address of the backend server.
* `server_type` - (Required) The type of the backend server. Valid values:

  - `Ecs`: Elastic Compute Service (ECS) instance
  - `Eni`: elastic network interface (ENI)
  - `Eci`: elastic container instance
  - `Ip`: IP address

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the resource was created. The time follows the ISO 8601 standard in the **yyyy-MM-ddTHH:mm:ssZ** format. The time is displayed in UTC.
* `servers` - The backend servers that you want to remove.
  * `port` - (Optional, Computed, Int) The port that is used by the backend server.
  * `server_group_id` - The server group ID.
  * `status` - Indicates the status of the backend server. Valid values:
  - `Adding`: The backend server is being added.
  - `Available`: The backend server is available.
  - `Draining`: The backend server is in connection draining.
  - `Removing`: The backend server is being removed.
  - `Replacing`: The backend server is being replaced.
* `status` - Indicates the status of the backend server. 


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Server Group.
* `delete` - (Defaults to 5 mins) Used when delete the Server Group.
* `update` - (Defaults to 5 mins) Used when update the Server Group.

## Import

GWLB Server Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_gwlb_server_group.example <id>
```