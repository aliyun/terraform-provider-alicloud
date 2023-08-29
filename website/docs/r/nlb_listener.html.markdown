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

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswtich" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = "${var.name}1"
  cidr_block   = "192.168.10.0/24"
}

resource "alicloud_vswitch" "vswtich2" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.1.id
  vswitch_name = "${var.name}2"
  cidr_block   = "192.168.20.0/24"
}

resource "alicloud_nlb_load_balancer" "nlb" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich2.id
    zone_id    = alicloud_vswitch.vswtich2.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich.id
    zone_id    = alicloud_vswitch.vswtich.zone_id
  }
  ipv6_address_type  = "Intranet"
  load_balancer_type = "Network"
  vpc_id             = alicloud_vpc.vpc.id
  address_type       = "Internet"
  address_ip_version = "Ipv4"
}

resource "alicloud_nlb_server_group" "sg1" {
  scheduler         = "Wrr"
  server_group_type = "Instance"
  vpc_id            = alicloud_vpc.vpc.id
  any_port_enabled  = true
  protocol          = "TCPSSL"
  server_group_name = "${var.name}4"
  resource_group_id = alicloud_nlb_load_balancer.nlb.resource_group_id
}

resource "alicloud_nlb_server_group" "sg2" {
  scheduler         = "Wrr"
  server_group_type = "Instance"
  vpc_id            = alicloud_vpc.vpc.id
  any_port_enabled  = true
  protocol          = "TCPSSL"
  server_group_name = "${var.name}5"
  resource_group_id = alicloud_nlb_load_balancer.nlb.resource_group_id
}


resource "alicloud_nlb_listener" "default" {
  ca_enabled             = true
  proxy_protocol_enabled = true
  sec_sensor_enabled     = true
  alpn_enabled           = true
  load_balancer_id       = alicloud_nlb_load_balancer.nlb.id
  server_group_id        = alicloud_nlb_server_group.sg1.id
  listener_protocol      = "TCPSSL"
  certificate_ids        = ["10793882-cn-hangzhou"]
  listener_description   = "test"
  security_policy_id     = "tls_cipher_policy_1_0"
  alpn_policy            = "HTTP1Only"
  start_port             = 244
  end_port               = 566
  idle_timeout           = 1
  ca_certificate_ids     = ["1ee31045-0d28-64c5-9ccf-6b50b7c1bd22"]
  status                 = "Running"
}
```

## Argument Reference

The following arguments are supported:
* `alpn_enabled` - (Optional, Available since v1.191.0) Whether ALPN is turned on. Value:
  - **true**: on.
  - **false**: closed.
* `alpn_policy` - (Optional, Available since v1.191.0) ALPN policy. Value:
  - **HTTP1Only**
  - **HTTP2Only**
  - **HTTP2Preferred**
  - **HTTP2Optional**.
* `ca_certificate_ids` - (Optional, Available since v1.191.0) The list of certificate authority (CA) certificates. This parameter takes effect only for listeners that use SSL over TCP. Note: Only one CA certificate is supported.
* `ca_enabled` - (Optional, Available since v1.191.0) Whether to start two-way authentication. Value:
  - **true**: start.
  - **false**: closed.
* `certificate_ids` - (Optional, Available since v1.191.0) The list of server certificates. This parameter takes effect only for listeners that use SSL over TCP. Note: Only one server certificate is supported.
* `cps` - (Optional, Available since v1.191.0) The maximum number of connections that can be created per second on the NLB instance. Valid values: 0 to 1000000. 0 specifies that the number of connections is unlimited.
* `end_port` - (Optional, ForceNew, Available since v1.191.0) Full port listening end port. Valid values: **0** ~ **65535**. The value of the end port is less than the start port.
* `idle_timeout` - (Optional, Computed, Available since v1.191.0) Connection idle timeout time. Unit: seconds. Valid values: **0** ~ **900**.
* `listener_description` - (Optional, Available since v1.191.0) Custom listener name.The length is limited to 2 to 256 characters, supports Chinese and English letters, and can include numbers, commas (,), half-width periods (.), half-width semicolons (;), forward slashes (/), at(@), underscores (_), and dashes (-).
* `listener_port` - (Required, ForceNew, Available since v1.191.0) Listening port. Valid values: **0** ~ **65535**. **0**: indicates that full port listening is used. When set to **0**, you must configure **StartPort** and **EndPort**.
* `listener_protocol` - (Required, ForceNew, Available since v1.191.0) The listening protocol. Valid values: **TCP**, **UDP**, or **TCPSSL**.
* `load_balancer_id` - (Required, ForceNew, Available since v1.191.0) The ID of the network-based server load balancer instance.
* `mss` - (Optional, Available since v1.191.0) The maximum segment size of the TCP message. Unit: Bytes. Valid values: **0** ~ **1500**. **0** indicates that the MSS value of the TCP message is not modified.
-> **NOTE:**  only TCP and TCPSSL listeners support this field value.
* `proxy_protocol_enabled` - (Optional, Available since v1.191.0) Whether to enable the Proxy Protocol to carry the source address of the client to the backend server. Value:
  - **true**: on.
  - **false**: closed.
* `sec_sensor_enabled` - (Optional, Available since v1.191.0) Whether to turn on the second-level monitoring function. Value:
  - **true**: on.
  - **false**: closed.
* `security_policy_id` - (Optional, Available since v1.191.0) The ID of the security policy. System security policies and custom security policies are supported. System security policies valid values: `tls_cipher_policy_1_0` (default), `tls_cipher_policy_1_1`, `tls_cipher_policy_1_2`, `tls_cipher_policy_1_2_strict`, and `tls_cipher_policy_1_2_strict_with_1_3`. Custom security policies can be created by resource `alicloud_nlb_security_policy`.
* `server_group_id` - (Required, Available since v1.191.0) The ID of the server group.
* `start_port` - (Optional, ForceNew, Available since v1.191.0) Full Port listens to the starting port. Valid values: **0** ~ **65535**.
* `status` - (Optional, Computed, Available since v1.191.0) The status of the resource.
* `tags` - (Optional, Map) The tag of the resource.

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