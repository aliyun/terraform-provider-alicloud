---
subcategory: "GWLB"
layout: "alicloud"
page_title: "Alicloud: alicloud_gwlb_server_group"
description: |-
  Provides a Alicloud GWLB Server Group resource.
---

# alicloud_gwlb_server_group

Provides a GWLB Server Group resource.



For information about GWLB Server Group and how to use it, see [What is Server Group](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gwlb_server_group&exampleId=d9955b88-af24-a5e2-ca7c-64cfecfcafb4335c91ce&activeTab=example&spm=docs.r.gwlb_server_group.0.d9955b88af&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc" "defaultEaxcvb" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultc3uVID" {
  vpc_id       = alicloud_vpc.defaultEaxcvb.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%s3", var.name)
}

resource "alicloud_security_group" "default" {
  name        = "tf-example"
  description = "New security group"
  vpc_id      = alicloud_vpc.defaultEaxcvb.id
}

resource "alicloud_instance" "default5DqP8f" {
  vswitch_id                 = alicloud_vswitch.defaultc3uVID.id
  image_id                   = "aliyun_2_1903_x64_20G_alibase_20231221.vhd"
  instance_type              = "ecs.g6.large"
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  instance_name              = format("%s4", var.name)
  description                = "tf-example-ecs"
  security_groups            = [alicloud_security_group.default.id]
  availability_zone          = alicloud_vswitch.defaultc3uVID.zone_id
  instance_charge_type       = "PostPaid"
}

resource "alicloud_gwlb_server_group" "default" {
  protocol          = "GENEVE"
  server_group_type = "Instance"
  vpc_id            = alicloud_vpc.defaultEaxcvb.id
  dry_run           = "false"
  server_group_name = "tf-exampleacccn-wulanchabugwlbservergroup24005"
  servers {
    server_id   = alicloud_instance.default5DqP8f.id
    server_type = "Ecs"
  }

  scheduler = "5TCH"
}
```

## Argument Reference

The following arguments are supported:
* `connection_drain_config` - (Optional, Computed, List) Connected graceful interrupt configuration. See [`connection_drain_config`](#connection_drain_config) below.
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. 
* `health_check_config` - (Optional, Computed, List) Health check configurations. See [`health_check_config`](#health_check_config) below.
* `protocol` - (Optional, ForceNew, Computed) Backend Protocol. Value:

  - *GENEVE (default)**.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `scheduler` - (Optional, Computed) Scheduling algorithm. Value:
  - **5TCH (default)**: quintuple hash, which is based on the consistent hash of the quintuple (source IP, Destination IP, source port, destination port, and protocol). The same flow is scheduled to the same backend server.
  - `3TCH`: a three-tuple hash, which is based on the consistent hash of three tuples (source IP address, destination IP address, and protocol). The same flow is dispatched to the same backend server.
  - `2TCH`: Binary Group hash, which is based on the consistent hash of the binary group (source IP and destination IP). The same flow is scheduled to the same backend server.
* `server_group_name` - (Optional) The server group name.

  It must be 2 to 128 characters in length, start with an uppercase letter or a Chinese character, and can contain digits, half-width periods (.), underscores (_), and dashes (-).
* `server_group_type` - (Optional, ForceNew, Computed) The server group type. Value:
  - **Instance (default)**: The instance type. You can add Ecs, Eni, and Eci instances to the server group.
  - `Ip`: The Ip address type. You can directly add backend servers of the Ip address type to the server group.
* `servers` - (Optional, Set) List of servers. See [`servers`](#servers) below.
* `tags` - (Optional, Map) List of resource tags.
* `vpc_id` - (Required, ForceNew) The VPC instance ID.

-> **NOTE:**  If the value of ServerGroupType is Instance, only servers in the VPC can be added to the server group.


### `connection_drain_config`

The connection_drain_config supports the following:
* `connection_drain_enabled` - (Optional, Computed) Whether to open the connection graceful interrupt. Value:
  - `true`: enabled.
  - `false`: Closed.
* `connection_drain_timeout` - (Optional, Computed, Int) Connection Grace interrupt timeout.

  Unit: seconds.

  Value range: 1~3600.

### `health_check_config`

The health_check_config supports the following:
* `health_check_connect_port` - (Optional, Computed, Int) The port of the backend server used for health check.

  Value range: **1 to 65535**.

  Default value: `80`.
* `health_check_connect_timeout` - (Optional, Computed, Int) The maximum timeout period for health check responses.

  Unit: seconds.

  Value range: **1 to 300**.

  Default value: `5`.
* `health_check_domain` - (Optional, Computed) The domain name used for health checks. Value:
  - **$SERVER_IP (default)**: Use the internal IP address of the backend server.
  - `domain`: Specify a specific domain name. The length is limited to 1 to 80 characters. Only uppercase and lowercase letters, digits, dashes (-), and periods (.) can be used.

-> **NOTE:**  This parameter takes effect only when the HealthCheckProtocol is HTTP.

* `health_check_enabled` - (Optional, Computed) Whether to enable health check. Value:
  - **true (default)**: enabled.
  - `false`: Closed.
* `health_check_http_code` - (Optional, Set) Health status return code list.
* `health_check_interval` - (Optional, Computed, Int) The time interval of the health check.

  Unit: seconds.

  Value range: **1~50**.

  Default value: `10`.
* `health_check_path` - (Optional, Computed) Health check path.

  It can be 1 to 80 characters in length and can only use upper and lower case letters, digits, dashes (-), forward slashes (/), half-width periods (.), percent signs (%), and half-width question marks (?), Pound sign (#) and and(&) and extended character set_;~! ()*[]@$^: ',+ =

  Must start with a forward slash (/).

-> **NOTE:**  This parameter takes effect only when the HealthCheckProtocol is HTTP.

* `health_check_protocol` - (Optional, Computed) Health check protocol, value:
  - `TCP` (default): Sends a SYN handshake packet to check whether the server port is alive.
  - `HTTP`: Sends a GET request to simulate the access behavior of the browser to check whether the server application is healthy.
* `healthy_threshold` - (Optional, Computed, Int) After the number of consecutive successful health checks, the health check status of the backend server is determined as successful from failed.

  Value range: **2 to 10**.

  Default value: `2`.
* `unhealthy_threshold` - (Optional, Computed, Int) The number of consecutive failed health checks that determine the health check status of the backend server from success to failure.

  Value range: **2 to 10**.

  Default value: `2`.

### `servers`

The servers supports the following:
* `server_id` - (Required) The ID of the backend server.
* `server_ip` - (Optional) Server ip.
* `server_type` - (Required) Backend server type. Valid values:
  - `Ecs`: ECS instance.
  - `Eni`: ENI instance.
  - `Eci`: ECI elastic container.
  - `Ip`: Ip address.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the server group.
* `servers` - List of servers.
  * `server_group_id` - The server group ID.
  * `port` - The port used by the backend server.
  * `status` - The status of the backend server. Value:
  - `Adding`: Adding.
  - `Available`: Normal Available status.
  - `Draining`: connection gracefully interrupting.
  - `Removing`: Removing.
  - `Replacing`: Replacing.
* `status` - Server group status. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Server Group.
* `delete` - (Defaults to 5 mins) Used when delete the Server Group.
* `update` - (Defaults to 5 mins) Used when update the Server Group.

## Import

GWLB Server Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_gwlb_server_group.example <id>
```