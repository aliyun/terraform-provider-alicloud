---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_server_group"
description: |-
  Provides a Alicloud Network Load Balancer (NLB) Server Group resource.
---

# alicloud_nlb_server_group

Provides a Network Load Balancer (NLB) Server Group resource.

For information about Network Load Balancer (NLB) Server Group and how to use it, see [What is Server Group](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createservergroup-nlb).

-> **NOTE:** Available since v1.186.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_server_group&exampleId=17064786-ee99-c7b2-1199-9e6b9966c75ef373983e&activeTab=example&spm=docs.r.nlb_server_group.0.17064786ee&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_nlb_server_group" "default" {
  resource_group_id        = data.alicloud_resource_manager_resource_groups.default.ids.0
  server_group_name        = var.name
  server_group_type        = "Instance"
  vpc_id                   = alicloud_vpc.default.id
  scheduler                = "Wrr"
  protocol                 = "TCP"
  connection_drain_enabled = true
  connection_drain_timeout = 60
  address_ip_version       = "Ipv4"
  health_check {
    health_check_enabled         = true
    health_check_type            = "TCP"
    health_check_connect_port    = 0
    healthy_threshold            = 2
    unhealthy_threshold          = 2
    health_check_connect_timeout = 5
    health_check_interval        = 10
    http_check_method            = "GET"
    health_check_http_code       = ["http_2xx", "http_3xx", "http_4xx"]
  }
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, ForceNew, Computed) The protocol version. Valid values:

  - `ipv4` (default): IPv4
  - `DualStack`: dual stack
* `any_port_enabled` - (Optional, ForceNew, Computed) Specifies whether to enable all-port forwarding. Valid values:

  - `true`
  - `false` (default)
* `connection_drain_enabled` - (Optional, Computed) Specifies whether to enable connection draining. Valid values:

  - `true`
  - `false` (default)
* `connection_drain_timeout` - (Optional, Computed, Int) The timeout period of connection draining. Unit: seconds. Valid values: `10` to `900`.
* `health_check` - (Optional, Computed, List) Health check configuration information. See [`health_check`](#health_check) below.
* `preserve_client_ip_enabled` - (Optional, Computed) Specifies whether to enable client IP preservation. Valid values:

  - `true` (default)
  - `false`
  -> **NOTE:** If you set the value to true and `protocol` to TCP, the server group cannot be associated with TCPSSL listeners.
* `protocol` - (Optional, ForceNew, Computed) The protocol used to forward requests to the backend servers. Valid values:

  - `TCP` (default)
  - `UDP`
  - `TCPSSL`
* `resource_group_id` - (Optional, Computed) The ID of the new resource group.
You can log on to the [Resource Management console](https://resourcemanager.console.aliyun.com/resource-groups) to view resource group IDs.
* `scheduler` - (Optional, Computed) The scheduling algorithm. Valid values:

  - **Wrr:** The weighted round-robin algorithm is used. Backend servers with higher weights receive more requests than backend servers with lower weights. This is the default value.
  - **rr:** The round-robin algorithm is used. Requests are forwarded to backend servers in sequence.
  - **sch:** Source IP hashing is used. Requests from the same source IP address are forwarded to the same backend server.
  - **tch:** Four-element hashing is used. It specifies consistent hashing that is based on four factors: source IP address, destination IP address, source port, and destination port. Requests that contain the same information based on the four factors are forwarded to the same backend server.
  - `qch`: QUIC ID hashing. Requests that contain the same QUIC ID are forwarded to the same backend server.
* `server_group_name` - (Required) The new name of the server group.
The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-). The name must start with a letter.
* `server_group_type` - (Optional, ForceNew, Computed) The type of server group. Valid values:

  - `Instance`: allows you to add servers of the `Ecs`, `Eni`, or `Eci` type. This is the default value.
  - `Ip`: allows you to add servers by specifying IP addresses.
* `tags` - (Optional, Map) Label.
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) to which the server group belongs.

-> **NOTE:**  If `ServerGroupType` is set to `Instance`, only servers in the specified VPC can be added to the server group.


The following arguments will be discarded. Please use new fields as soon as possible:
* `connection_drain` - (Deprecated since v1.231.0). Field 'connection_drain' has been deprecated from provider version 1.231.0. New field 'connection_drain_enabled' instead.

### `health_check`

The health_check supports the following:
* `health_check_connect_port` - (Optional, Computed, Int) The port that you want to use for health checks on backend servers.
Valid values: `0` to `65535`.
Default value: `0`. If you set the value to 0, the port of the backend server is used for health checks.
* `health_check_connect_timeout` - (Optional, Computed, Int) The maximum timeout period of a health check. Unit: seconds. Valid values: `1` to `300`. Default value: `5`.
* `health_check_domain` - (Optional, Computed) The domain name that you want to use for health checks. Valid values:
  - `$SERVER_IP`: the private IP address of a backend server.
  - `domain`: a specified domain name. The domain name must be 1 to 80 characters in length, and can contain lowercase letters, digits, hyphens (-), and periods (.).

-> **NOTE:**  This parameter takes effect only when `HealthCheckType` is set to `HTTP`.

* `health_check_enabled` - (Optional, Computed) Specifies whether to enable the health check feature. Valid values:

  - `true` (default)
  - `false`
* `health_check_exp` - (Optional, Available since v1.243.0) health check response character string. The value contains a maximum of 512 characters
* `health_check_http_code` - (Optional, Computed, List) The HTTP status codes to return for health checks. Separate multiple HTTP status codes with commas (,). Valid values: `http\_2xx` (default), `http\_3xx`, `http\_4xx`, and `http\_5xx`.

-> **NOTE:**  This parameter takes effect only when `HealthCheckType` is set to `HTTP`.

* `health_check_interval` - (Optional, Computed, Int) The interval at which health checks are performed. Unit: seconds.
Valid values: `5` to `50`.
Default value: `10`.
* `health_check_req` - (Optional, Available since v1.243.0) UDP healthy check request string, the value is a character string of 512 characters
* `health_check_type` - (Optional, Computed) The protocol that you want to use for health checks. Valid values: `TCP` (default) and `HTTP`.
* `health_check_url` - (Optional, Computed) The path to which health check requests are sent.

The path must be 1 to 80 characters in length, and can contain only letters, digits, and the following special characters: `- / . % ? # & =`. It can also contain the following extended characters: `_ ; ~ ! ( ) * [ ] @ $ ^ : ' , +`. The path must start with a forward slash (/).

-> **NOTE:**  This parameter takes effect only when `HealthCheckType` is set to `HTTP`.

* `healthy_threshold` - (Optional, Computed, Int) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. In this case, the health status changes from `fail` to `success`.
Valid values: `2` to `10`.
Default value: `2`.
* `http_check_method` - (Optional) The HTTP method that is used for health checks. Valid values: `GET` (default) and `HEAD`.

-> **NOTE:**  This parameter takes effect only when `HealthCheckType` is set to `HTTP`.

* `unhealthy_threshold` - (Optional, Computed, Int) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. In this case, the health status changes from `success` to `fail`.
Valid values: `2` to `10`.
Default value: `2`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The ID of the region where the NLB instance is deployed.
* `status` - Server group status. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Server Group.
* `delete` - (Defaults to 5 mins) Used when delete the Server Group.
* `update` - (Defaults to 5 mins) Used when update the Server Group.

## Import

Network Load Balancer (NLB) Server Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_server_group.example <id>
```