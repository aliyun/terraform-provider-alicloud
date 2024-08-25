---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_listener"
description: |-
  Provides a Alicloud NLB Listener resource.
---

# alicloud_nlb_listener

Provides a NLB Listener resource. 

For information about NLB Listener and how to use it, see [What is Listener](https://www.alibabacloud.com/help/en/server-load-balancer/latest/api-nlb-2022-04-30-createlistener).

-> **NOTE:** Available since v1.191.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_nlb_listener&exampleId=719b36d3-b3c5-f8e2-9897-0b9972458db784c383ce&activeTab=example&spm=docs.r.nlb_listener.0.719b36d3b3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_nlb_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "default1" {
  vswitch_name = var.name
  cidr_block   = "10.4.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  vpc_id             = alicloud_vpc.default.id
  tags = {
    Created = "TF",
    For     = "example",
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default1.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
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

resource "alicloud_nlb_listener" "default" {
  listener_protocol      = "TCP"
  listener_port          = "80"
  listener_description   = var.name
  load_balancer_id       = alicloud_nlb_load_balancer.default.id
  server_group_id        = alicloud_nlb_server_group.default.id
  idle_timeout           = "900"
  proxy_protocol_enabled = "true"
  cps                    = "10000"
  mss                    = "0"
}
```

## Argument Reference

The following arguments are supported:
* `alpn_enabled` - (Optional) Whether ALPN is turned on. Value:
  - **true**: on.
  - **false**: closed.
* `alpn_policy` - (Optional) ALPN policy. Value:
  - **HTTP1Only**
  - **HTTP2Only**
  - **HTTP2Preferred**
  - **HTTP2Optional**.
* `ca_certificate_ids` - (Optional) CA certificate list information. Currently, only one CA certificate can be added.
-> **NOTE:**  This parameter only takes effect for TCPSSL listeners.
* `ca_enabled` - (Optional) Whether to start two-way authentication. Value:
  - **true**: start.
  - **false**: closed.
* `certificate_ids` - (Optional) Server certificate list information. Currently, only one server certificate can be added.
-> **NOTE:**  This parameter only takes effect for TCPSSL listeners.
* `cps` - (Optional) The new connection speed limit for a network-based load balancing instance per second. Valid values: **0** ~ **1000000**. **0** indicates unlimited speed.
* `end_port` - (Optional, ForceNew) Full port listening end port. Valid values: **0** ~ **65535 * *. The value of the end port is less than the start port.
* `idle_timeout` - (Optional) Connection idle timeout time. Unit: seconds. Valid values: **1** ~ **900**.
* `listener_description` - (Optional) Custom listener name.The length is limited to 2 to 256 characters, supports Chinese and English letters, and can include numbers, commas (,), half-width periods (.), half-width semicolons (;), forward slashes (/), at(@), underscores (_), and dashes (-).
* `listener_port` - (Required, ForceNew) Listening port. Valid values: **0** ~ **65535 * *. **0**: indicates that full port listening is used. When set to **0**, you must configure **StartPort** and **EndPort**.
* `listener_protocol` - (Required, ForceNew) The listening protocol. Valid values: **TCP**, **UDP**, or **TCPSSL**.
* `load_balancer_id` - (Required, ForceNew) The ID of the network-based server load balancer instance.
* `mss` - (Optional) The maximum segment size of the TCP message. Unit: Bytes. Valid values: **0** ~ **1500**. **0** indicates that the MSS value of the TCP message is not modified.
-> **NOTE:**  only TCP and TCPSSL listeners support this field value.
* `proxy_protocol_enabled` - (Optional) Whether to enable the Proxy Protocol to carry the source address of the client to the backend server. Value:
  - **true**: on.
  - **false**: closed.
* `sec_sensor_enabled` - (Optional) Whether to turn on the second-level monitoring function. Value:
  - **true**: on.
  - **false**: closed.
* `security_policy_id` - (Optional) Security policy ID. Support system security policies and custom security policies. Valid values: **tls_cipher_policy_1_0**, **tls_cipher_policy_1_1**, **tls_cipher_policy_1_2**, **tls_cipher_policy_1_2_strict**, or **tls_cipher_policy_1_2_strict_with_1_3**.
-> **NOTE:**  This parameter only takes effect for TCPSSL listeners.
* `server_group_id` - (Required) The ID of the server group.
* `start_port` - (Optional, ForceNew) Full Port listens to the starting port. Valid values: **0** ~ **65535**.
* `status` - (Optional, Computed) The status of the resource.
* `tags` - (Optional, Map, Available since v1.217.1) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Listener.
* `delete` - (Defaults to 5 mins) Used when delete the Listener.
* `update` - (Defaults to 5 mins) Used when update the Listener.

## Import

NLB Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_listener.example <id>
```