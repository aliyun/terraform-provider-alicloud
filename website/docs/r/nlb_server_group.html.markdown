---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_server_group"
description: |-
  Provides a Alicloud NLB Server Group resource.
---

# alicloud_nlb_server_group

Provides a NLB Server Group resource. 

For information about NLB Server Group and how to use it, see [What is Server Group](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createservergroup-nlb).

-> **NOTE:** Available since v1.186.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_nlb_server_group&exampleId=208efd1e-a7d0-6cb4-1cb9-6ebebde32e91673e4e8b&activeTab=example&spm=docs.r.nlb_server_group.0.208efd1ea7" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

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
  connection_drain         = true
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
* `address_ip_version` - (Optional, ForceNew, Computed) The protocol version. Valid values: `Ipv4` (default), `DualStack`.
* `any_port_enabled` - (Optional, ForceNew, Computed, Available since v1.214.0) Full port forwarding.
* `connection_drain_enabled` - (Optional, Computed, Available since v1.214.0) Specifies whether to enable connection draining.
* `connection_drain_timeout` - (Optional, Computed) Set the connection elegant interrupt timeout. Unit: seconds. Valid values: **10** ~ **900**.
* `health_check` - (Optional, ForceNew, Computed) Health check configuration information. See [`health_check`](#health_check) below.
* `preserve_client_ip_enabled` - (Optional, Computed) Indicates whether client address retention is enabled. Special instructions: When **AddressIPVersion** is of the **ipv4** type, the default value is **true**. **Addrestipversion** can only be **false** when the value of **ipv6** is **ipv6**, and can be **true** when supported by the underlying layer.
* `protocol` - (Optional, ForceNew, Computed) The backend protocol. Valid values: `TCP` (default), `UDP`, and `TCPSSL`.
* `resource_group_id` - (Optional, Computed)  The ID of the resource group to which the security group belongs.
* `scheduler` - (Optional, Computed) The routing algorithm. Valid values:
  - `Wrr` (default): The Weighted Round Robin algorithm is used. Backend servers with higher weights receive more requests than backend servers with lower weights.
  - `Rr`: The round-robin algorithm is used. Requests are forwarded to backend servers in sequence.
  - `Sch`: Source IP hashing is used. Requests from the same source IP address are forwarded to the same backend server.
  - `Tch`: Four-element hashing is used. It specifies consistent hashing that is based on four factors: source IP address, destination IP address, source port, and destination port. Requests that contain the same information based on the four factors are forwarded to the same backend server.
  - `Qch`: QUIC ID hashing is used. Requests that contain the same QUIC ID are forwarded to the same backend server.
* `server_group_name` - (Required) The name of the server group. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
* `server_group_type` - (Optional, ForceNew, Computed) The type of the server group. Valid values: 
  - `Instance` (default): allows you to specify `Ecs`, `Ens`, or `Eci`.
  - `Ip`: allows you to specify IP addresses.
* `tags` - (Optional, Map) Label.
* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the server group belongs.

The following arguments will be discarded. Please use new fields as soon as possible:
* `connection_drain` - (Deprecated since v1.214.0). Field 'connection_drain' has been deprecated from provider version 1.214.0. New field 'connection_drain_enabled' instead.

### `health_check`

The health_check supports the following:
* `health_check_connect_port` - (Optional, Computed) The port of the backend server for health checks. Valid values: **0** ~ **65535**. **0** indicates that the port of the backend server is used for health check.
* `health_check_connect_timeout` - (Optional, Computed) Maximum timeout for health check responses. Unit: seconds. Valid values: **1** ~ **300**.
* `health_check_domain` - (Optional, Computed) The domain name used for health check. Valid values:
  - **$SERVER_IP**: uses the intranet IP of the backend server.
  - **domain**: Specify a specific domain name. The length is limited to 1 to 80 characters. Only lowercase letters, numbers, dashes (-), and half-width periods (.) can be used.
-> **NOTE:**  This parameter takes effect only when **HealthCheckType** is **HTTP**.
* `health_check_enabled` - (Optional, Computed) Whether to enable health check. Valid values:
  - **true**: on.
  - **false**: closed.
* `health_check_http_code` - (Optional, Computed) Health status return code. Multiple status codes are separated by commas (,). Valid values: **http\_2xx**, **http\_3xx**, **http\_4xx**, and **http\_5xx**.
-> **NOTE:**  This parameter takes effect only when **HealthCheckType** is **HTTP**.
* `health_check_interval` - (Optional, Computed) Time interval of health examination. Unit: seconds.  Valid values: **5** ~ **50**.
* `health_check_type` - (Optional, Computed) Health check protocol. Valid values: **TCP** or **HTTP**.
* `health_check_url` - (Optional, Computed) Health check path.
-> **NOTE:**  This parameter takes effect only when **HealthCheckType** is **HTTP**.
* `healthy_threshold` - (Optional, Computed) After the health check is successful, the health check status of the backend server is determined from **failed** to **successful**.  Valid values: **2** to **10**.
* `http_check_method` - (Optional) The health check method. Valid values: **GET** or **HEAD**.
-> **NOTE:**  This parameter takes effect only when **HealthCheckType** is **HTTP**.
* `unhealthy_threshold` - (Optional, Computed) After the health check fails for many times in a row, the health check status of the backend server is determined from **Success** to **Failure**. Valid values: **2** to **10**.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - Server group status. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Server Group.
* `delete` - (Defaults to 5 mins) Used when delete the Server Group.
* `update` - (Defaults to 5 mins) Used when update the Server Group.

## Import

NLB Server Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_server_group.example <id>
```