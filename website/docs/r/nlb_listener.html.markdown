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
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_listener&exampleId=719b36d3-b3c5-f8e2-9897-0b9972458db784c383ce&activeTab=example&spm=docs.r.nlb_listener.0.719b36d3b3&intl_lang=EN_US" target="_blank">
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
* `alpn_enabled` - (Optional, Computed) Specifies whether to enable Application-Layer Protocol Negotiation (ALPN). Valid values:
  - `true`
  - `false` (default)

-> **NOTE:**  Effective only for TCPSSL listener

* `alpn_policy` - (Optional) The ALPN policy. Valid values:
  - `HTTP1Only`: uses only HTTP 1.x. The priority of HTTP 1.1 is higher than the priority of HTTP 1.0.
  - `HTTP2Only`: uses only HTTP 2.0.
  - `HTTP2Optional`: preferentially uses HTTP 1.x over HTTP 2.0. The priority of HTTP 1.1 is higher than the priority of HTTP 1.0, and the priority of HTTP 1.0 is higher than the priority of HTTP 2.0.
  - `HTTP2Preferred`: preferentially uses HTTP 2.0 over HTTP 1.x. The priority of HTTP 2.0 is higher than the priority of HTTP 1.1, and the priority of HTTP 1.1 is higher than the priority of HTTP 1.0.

-> **NOTE:**  This parameter is required if AlpnEnabled is set to true.

-> **NOTE:**  Effective only for TCPSSL listener.

* `ca_certificate_ids` - (Optional, List) The list of certificate authority (CA) certificates. This parameter takes effect only for listeners that use SSL over TCP. 

-> **NOTE:**  Only one CA certificate is supported.

* `ca_enabled` - (Optional, Computed) Specifies whether to enable mutual authentication. Valid values:
  - `true` : yes
  - `false` (default): no
* `certificate_ids` - (Optional, List) The list of server certificates. This parameter takes effect only for listeners that use SSL over TCP. 

-> **NOTE:**  This parameter takes effect only for TCPSSL listeners.

* `cps` - (Optional, Int) The maximum number of connections that can be created per second on the NLB instance. Valid values: `0` to `1000000`. `0` specifies that the number of connections is unlimited.
* `end_port` - (Optional, ForceNew) The last port in the listener port range. Valid values: `0` to `65535`. The number of the last port must be greater than the number of the first port.

-> **NOTE:**  This parameter is required when `ListenerPort` is set to `0`.

* `idle_timeout` - (Optional, Int) The timeout period of idle connections. Unit: seconds. Valid values: `1` to `900`. Default value: `900`.
* `listener_description` - (Optional) Enter a name for the listener.

  The description must be 2 to 256 characters in length, and can contain letters, digits, commas (,), periods (.), semicolons (;), forward slashes (/), at signs (@), underscores (\_), and hyphens (-).
* `listener_port` - (Required, ForceNew, Int) The listener port. Valid values: `0` to `65535`.

  If you set the value to `0`, the listener listens by port range. If you set the value to `0`, you must specify `StartPort` and `EndPort`.
* `listener_protocol` - (Required, ForceNew) The listening protocol. Valid values: `TCP`, `UDP`, and `TCPSSL`.
* `load_balancer_id` - (Required, ForceNew) The ID of the Network Load Balancer (NLB) instance.
* `mss` - (Optional, Int) The maximum size of a TCP segment. Unit: bytes. Valid values: `0` to `1500`. `0` specifies that the maximum segment size remains unchanged.

-> **NOTE:**  This parameter is supported only by TCP listeners and listeners that use SSL over TCP.

* `proxy_protocol_enabled` - (Optional, Computed) Specifies whether to use the Proxy protocol to pass client IP addresses to backend servers. Valid values:
  - `true`
  - `false` (default)
* `sec_sensor_enabled` - (Optional, Computed) Specifies whether to enable fine-grained monitoring. Valid values:
  - `true`
  - `false` (default)

-> **NOTE:**  Before enabling this function, ensure that the HdMonitor storage has been configured in the region. Otherwise, create listener may fails.

* `security_policy_id` - (Optional, Computed) The security policy ID. System security policies and custom security policies are supported.

  Valid values: `tls_cipher_policy\_1\_0` (default), `tls_cipher_policy\_1\_1`, `tls_cipher_policy\_1\_2`, `tls_cipher_policy\_1\_2\_strict`, and `tls_cipher_policy\_1\_2\_strict_with\_1\_3`.

-> **NOTE:**  This parameter takes effect only for listeners that use SSL over TCP.

* `server_group_id` - (Required) The ID of the server group.
* `start_port` - (Optional, ForceNew) The first port in the listener port range. Valid values: `0` to `65535`.

-> **NOTE:**  This parameter is required when `ListenerPort` is set to `0`.

* `status` - (Optional, Computed) The status of the resource. Valid values: `Running`, `Stopped`. When you want to enable this instance, you can set the property value to `Running`; 
* `tags` - (Optional, Map) The tag of the resource

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