---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_server_group"
description: |-
  Provides a Alicloud NLB Server Group resource.
---

# alicloud_nlb_server_group

Provides a NLB Server Group resource. 

For information about NLB Server Group and how to use it, see [What is Server Group](https://www.alibabacloud.com/help/en/server-load-balancer/latest/api-nlb-2022-04-30-createservergroup).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_vpc" "defaultEyoN1M" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_resource_manager_resource_group" "rg1" {
  display_name        = "nlbrg1"
  resource_group_name = "${var.name}1"
}

resource "alicloud_resource_manager_resource_group" "rg2" {
  display_name        = "nlbrg2"
  resource_group_name = "${var.name}2"
}


resource "alicloud_nlb_server_group" "default" {
  preserve_client_ip_enabled = true
  connection_drain_enabled   = true
  any_port_enabled           = true
  protocol                   = "TCP"
  server_group_type          = "Instance"
  server_group_name          = var.name
  vpc_id                     = alicloud_vpc.defaultEyoN1M.id
  health_check {
    health_check_enabled         = true
    health_check_type            = "HTTP"
    health_check_connect_port    = 1
    healthy_threshold            = 2
    unhealthy_threshold          = 2
    health_check_connect_timeout = 1
    health_check_interval        = 5
    health_check_domain          = "$SERVER_IP"
    health_check_url             = "/rdktest"
    health_check_http_code       = ["http_2xx"]
    http_check_method            = "HEAD"
  }
  resource_group_id        = alicloud_resource_manager_resource_group.rg1.id
  scheduler                = "Wrr"
  connection_drain_timeout = 10
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, ForceNew, Computed, Available since v1.186.0) Protocol version. Value:
  - **Ipv4**:IPv4 type.
  - **DualStack**: Double Stack type.
* `any_port_enabled` - (Optional, ForceNew, Computed) Whether to enable full port forwarding.
* `connection_drain_enabled` - (Optional, Computed) Whether to open the connection gracefully interrupted. Value:
  - **true**: on.
  - **false**: closed.
* `connection_drain_timeout` - (Optional, Computed, Available since v1.186.0) Set the connection elegant interrupt timeout. Unit: seconds.Valid values: **10** ~ **900**.
* `health_check` - (Optional, ForceNew, Computed, Available since v1.186.0) Health check configuration information. See [`health_check`](#health_check) below.
* `preserve_client_ip_enabled` - (Optional, Computed, Available since v1.186.0) Whether to enable the client address retention function. Value:
  - **true**: on.
  - **false**: closed.
-> **NOTE:**  special instructions: When **AddressIPVersion** is of the **ipv4** type, the default value is **true * *. **Addrestipversion** can only be **false** when the value of **ipv6** is **ipv6**, and can be **true** when supported by the underlying layer * *.
* `protocol` - (Optional, ForceNew, Computed, Available since v1.186.0) The backend Forwarding Protocol. Valid values: **TCP**, **UDP**, or **TCPSSL**.
* `resource_group_id` - (Optional, Computed, Available since v1.186.0) The ID of the resource group to which the server group belongs.
* `scheduler` - (Optional, Computed, Available since v1.186.0) Scheduling algorithm. Value:
  - **Wrr**: Weighted polling. The higher the weight of the backend server, the higher the probability of being polled.
  - **Rr**: polls external requests are distributed to backend servers in sequence according to the access order. sch: Source IP hash: The same source address is scheduled to the same backend server.
  - **Tch**: Quadruple hash, based on the consistent hash of the Quad (source IP, Destination IP, source port, and destination port), the same stream is scheduled to the same backend server.
  - **Qch**: a QUIC ID hash that allows you to hash requests with the same QUIC ID to the same backend server.
  - **Sch**: Source IP hashing is used. Requests from the same source IP address are forwarded to the same backend server.
* `server_group_name` - (Required, Available since v1.186.0) The name of the server group. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
* `server_group_type` - (Optional, ForceNew, Computed, Available since v1.186.0) Server group type. Value:
  - **Instance**: The server type. You can add **Ecs**, **Ens**, and **Eci** instances to the server group.
  - **Ip**: Ip address type. You can add Ip addresses to a server group of this type.
* `tags` - (Optional, Map, Available since v1.186.0) Label.
* `vpc_id` - (Required, ForceNew, Available since v1.186.0) The ID of the VPC to which the server group belongs.

The following arguments will be discarded. Please use new fields as soon as possible:
* `connection_drain` - (Deprecated since v1.210.0). Field 'connection_drain' has been deprecated from provider version 1.210.0. New field 'connection_drain_enabled' instead.

### `health_check`

The health_check supports the following:
* `health_check_connect_port` - (Optional, Computed) The port of the backend server for health checks.Valid values: **0** ~ **65535**. **0** indicates that the port of the backend server is used for health check.
* `health_check_connect_timeout` - (Optional, Computed) Maximum timeout for health check responses. Unit: seconds.Valid values: **1** ~ **300**.
* `health_check_domain` - (Optional, Computed) The domain name used for health check. Value:
  - **$SERVER_IP**: uses the intranet IP of the backend server.
  - **domain**: Specify a specific domain name. The length is limited to 1 to 80 characters. Only lowercase letters, numbers, dashes (-), and half-width periods (.) can be used.
-> **NOTE:**  This parameter takes effect only when **HealthCheckType** is **HTTP.
* `health_check_enabled` - (Optional, Computed) Whether to enable health check. Valid values:
  - **true**: on.
  - **false**: closed.
* `health_check_http_code` - (Optional, Computed) Health status return code. Multiple status codes are separated by commas (,).
Valid values: **http\_2xx**, **http\_3xx**, **http\_4xx**, and **http\_5xx * *.
-> **NOTE:**  This parameter takes effect only when **HealthCheckType** is **HTTP.
* `health_check_interval` - (Optional, Computed) Time interval of health examination. Unit: seconds. Valid values: **5** ~ **5000**.
* `health_check_type` - (Optional, Computed) Health check protocol. Valid values: **TCP** or **HTTP**.
* `health_check_url` - (Optional, Computed) Health check path.
-> **NOTE:**  This parameter takes effect only when **HealthCheckType** is **HTTP**.
* `healthy_threshold` - (Optional, Computed) After the health check is successful, the health check status of the backend server is determined from **failed** to **successful**. Valid values: **2** to **10**.
* `http_check_method` - (Optional, Computed) The health check method. Valid values: **GET** or **HEAD**.
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