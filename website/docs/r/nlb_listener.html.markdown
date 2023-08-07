---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_listener"
sidebar_current: "docs-alicloud-resource-nlb-listener"
description: |-
  Provides a Alicloud NLB Listener resource.
---

# alicloud\_nlb\_listener

Provides a NLB Listener resource.

For information about NLB Listener and how to use it, see [What is Listener](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createlistener-nl).

-> **NOTE:** Available since v1.191.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_nlb_server_group" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  server_group_name = var.name
  server_group_type = "Instance"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  scheduler         = "Wrr"
  protocol          = "TCP"
  health_check {
    health_check_url             = "/test/index.html"
    health_check_domain          = "tf-testAcc.com"
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
  connection_drain           = true
  connection_drain_timeout   = 60
  preserve_client_ip_enabled = true
  tags = {
    Created = "TF"
  }
  address_ip_version = "Ipv4"
}

data "alicloud_nlb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nlb_zones.default.zones.0.id
}
data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nlb_zones.default.zones.1.id
}
locals {
  zone_id_1    = data.alicloud_nlb_zones.default.zones.0.id
  vswitch_id_1 = data.alicloud_vswitches.default_1.ids[0]
  zone_id_2    = data.alicloud_nlb_zones.default.zones.1.id
  vswitch_id_2 = data.alicloud_vswitches.default_2.ids[0]
}
resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  tags = {
    Created = "tfTestAcc0"
    For     = "Tftestacc 0"
  }
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_mappings {
    vswitch_id = local.vswitch_id_1
    zone_id    = local.zone_id_1
  }
  zone_mappings {
    vswitch_id = local.vswitch_id_2
    zone_id    = local.zone_id_2
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
  sec_sensor_enabled     = "true"
  cps                    = "10000"
  mss                    = "0"
}
```

## Argument Reference

The following arguments are supported:

* `listener_description` - (Optional) Custom listener name. The length is limited to 2 to 256 characters, supports Chinese and English letters, and can include numbers, commas (,), half-width periods (.), half-width semicolons (;), forward slashes (/), at(@), underscores (_), and dashes (-).
* `listener_port` - (Required, ForceNew) Listening port. Valid values: 0 ~ 65535. `0`: indicates that full port listening is used. When set to `0`, you must configure `StartPort` and `EndPort`.
* `listener_protocol` - (Required, ForceNew) The listening protocol. Valid values: `TCP`, `UDP`, or `TCPSSL`.
* `load_balancer_id` - (Required, ForceNew) The ID of the network-based server load balancer instance.
* `server_group_id` - (Required) The ID of the server group.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Running`, `Stopped`.
* `end_port` - (Optional, ForceNew) Full port listening end port. Valid values: `0` ~ `65535`. The value of the end port is less than the start port.
* `start_port` - (Optional, ForceNew) Full Port listens to the starting port. Valid values: `0` ~ `65535`.
* `alpn_enabled` - (Optional, Computed) Specifies whether to enable Application-Layer Protocol Negotiation (ALPN).
* `ca_certificate_ids` - (Optional) The list of certificate authority (CA) certificates. This parameter takes effect only for listeners that use SSL over TCP. **Note:** Only one CA certificate is supported.
* `sec_sensor_enabled` - (Optional, Computed) Specifies whether to enable fine-grained monitoring.
* `certificate_ids` - (Optional) The list of server certificates. This parameter takes effect only for listeners that use SSL over TCP. **Note:** Only one server certificate is supported.
* `idle_timeout` - (Optional, Computed) The timeout period of an idle connection. Unit: seconds. Valid values: `1` to `900`. Default value: `900`.
* `security_policy_id` - (Optional, Computed) The ID of the security policy. System security policies and custom security policies are supported. 
  System security policies valid values: `tls_cipher_policy_1_0` (default), `tls_cipher_policy_1_1,` `tls_cipher_policy_1_2`, `tls_cipher_policy_1_2_strict`, and `tls_cipher_policy_1_2_strict_with_1_3`.
  Custom security policies can be created by resource `alicloud_nlb_security_policy`.
* `alpn_policy` - (Optional) The ALPN policy.
* `proxy_protocol_enabled` - (Optional, Computed) Specifies whether to use the Proxy protocol to pass client IP addresses to backend servers.
* `ca_enabled` - (Optional, Computed) Specifies whether to enable mutual authentication.
* `mss` - (Optional) The maximum size of a TCP segment. Unit: bytes. Valid values: 0 to 1500. 0 specifies that the maximum segment size remains unchanged. **Note:** This parameter is supported only by listeners that use SSL over TCP.
* `cps` - (Optional) The maximum number of connections that can be created per second on the NLB instance. Valid values: 0 to 1000000. 0 specifies that the number of connections is unlimited.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Listener.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Listener.
* `delete` - (Defaults to 1 mins) Used when delete the Listener.
* `update` - (Defaults to 1 mins) Used when update the Listener.

## Import

NLB Listener can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_listener.example <id>
```