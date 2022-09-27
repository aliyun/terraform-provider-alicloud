---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_server_group"
sidebar_current: "docs-alicloud-resource-nlb-server-group"
description: |-
  Provides a Alicloud NLB Server Group resource.
---

# alicloud\_nlb\_server\_group

Provides a NLB Server Group resource.

For information about NLB Server Group and how to use it, see [What is Server Group](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createservergroup-nlb).

-> **NOTE:** Available in v1.186.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
resource "alicloud_nlb_server_group" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  server_group_name = var.name
  server_group_type = "Instance"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  scheduler         = "Wrr"
  protocol          = "TCP"
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
  connection_drain         = true
  connection_drain_timeout = 60
  tags = {
    Created = "TF"
  }
  address_ip_version = "Ipv4"
}
```

## Argument Reference

The following arguments are supported:

* `address_ip_version` - (Optional, Computed, ForceNew) The protocol version. Valid values: `Ipv4` (default), `DualStack`.
* `connection_drain` - (Optional, Computed) Specifies whether to enable connection draining.
* `connection_drain_timeout` - (Optional, Computed) The timeout period of connection draining. Unit: seconds. Valid values: 10 to 900.
* `health_check` - (Required) HealthCheck. See the following `Block health_check`.
* `protocol` - (Optional, Computed, ForceNew) The backend protocol. Valid values: `TCP` (default), `UDP`, and `TCPSSL`.
* `resource_group_id` - (Optional, Computed, ForceNew) The ID of the resource group to which the security group belongs.
* `scheduler` - (Optional, Computed) The routing algorithm. Valid values:
  - `Wrr` (default): The Weighted Round Robin algorithm is used. Backend servers with higher weights receive more requests than backend servers with lower weights.
  - `Rr`: The round-robin algorithm is used. Requests are forwarded to backend servers in sequence.
  - `Sch`: Source IP hashing is used. Requests from the same source IP address are forwarded to the same backend server.
  - `Tch`: Four-element hashing is used. It specifies consistent hashing that is based on four factors: source IP address, destination IP address, source port, and destination port. Requests that contain the same information based on the four factors are forwarded to the same backend server.
  - `Qch`: QUIC ID hashing is used. Requests that contain the same QUIC ID are forwarded to the same backend server.
* `server_group_name` - (Required) The name of the server group. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
* `server_group_type` - (Optional, Computed, ForceNew) The type of the server group. Valid values:
  - `Instance` (default): allows you to specify `Ecs`, `Ens`, or `Eci`.
  - `Ip`: allows you to specify IP addresses.
* `vpc_id` - (Required, ForceNew) The id of the vpc.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `preserve_client_ip_enabled` - (Optional, Computed) Indicates whether client address retention is enabled.

#### Block health_check

The health_check supports the following: 

* `health_check_connect_port` - (Optional) The backend port that is used for health checks. Valid values: 0 to 65535. Default value: 0. If you set the value to 0, the port of a backend server is used for health checks.
* `health_check_connect_timeout` - (Optional) The maximum timeout period of a health check response. Unit: seconds. Valid values: 1 to 300. Default value: 5.
* `health_check_domain` - (Optional) The domain name that is used for health checks. Valid values:
  - `$SERVER_IP`: the private IP address of a backend server.
  - `domain`: a specified domain name. The domain name must be 1 to 80 characters in length, and can contain letters, digits, hyphens (-), and periods (.).
* `health_check_enabled` - (Optional) Specifies whether to enable health checks.
* `health_check_interval` - (Optional) The interval between two consecutive health checks. Unit: seconds. Valid values: 5 to 5000. Default value: 10.
* `health_check_type` - (Optional) The protocol that is used for health checks. Valid values: `TCP` (default) and `HTTP`.
* `health_check_url` - (Optional) The path to which health check requests are sent. The path must be 1 to 80 characters in length, and can contain only letters, digits, and the following special characters: `- / . % ? # & =`. It can also contain the following extended characters: `_ ; ~ ! ( ) * [ ] @ $ ^ : ' , +`. The path must start with a forward slash (/). **Note:** This parameter takes effect only if `health_check_type` is set to `http`.
* `healthy_threshold` - (Optional) The number of times that an unhealthy backend server must consecutively pass health checks before it is declared healthy. In this case, the health status is changed from fail to success. Valid values: 2 to 10. Default value: 2.
* `unhealthy_threshold` - (Optional) The number of times that a healthy backend server must consecutively fail health checks before it is declared unhealthy. In this case, the health status is changed from success to fail. Valid values: 2 to 10. Default value: 2.
* `http_check_method` - (Optional) The HTTP method that is used for health checks. Valid values: `GET` and `HEAD`. **Note:** This parameter takes effect only if `health_check_type` is set to `http`.
* `health_check_http_code` - (Optional) The HTTP status codes to return to health checks. Separate multiple HTTP status codes with commas (,). Valid values: http_2xx (default), http_3xx, http_4xx, and http_5xx. **Note:** This parameter takes effect only if `health_check_type` is set to `http`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Server Group.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Server Group.
* `update` - (Defaults to 1 mins) Used when update the Server Group.
* `delete` - (Defaults to 1 mins) Used when delete the Server Group.

## Import

NLB Server Group can be imported using the id, e.g.

```
$ terraform import alicloud_nlb_server_group.example <id>
```