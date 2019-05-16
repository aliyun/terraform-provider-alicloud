---
layout: "alicloud"
page_title: "Alicloud: alicloud_forward_entry"
sidebar_current: "docs-alicloud-resource-vpc"
description: |-
  Provides a Alicloud forward resource.
---

# alicloud\_forward

Provides a forward resource.

## Example Usage

Basic Usage

```
variable "name" {
	default = "ForwardEntryConfig"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_forward_entry" "foo"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

resource "alicloud_forward_entry" "foo1"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "443"
	ip_protocol = "udp"
	internal_ip = "172.16.0.4"
	internal_port = "8080"
}
resource "alicloud_forward_entry" "foo2"{
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "99"
	ip_protocol = "udp"
	internal_ip = "172.16.0.5"
	internal_port = "8082"
}


```
## Argument Reference

The following arguments are supported:

* `forward_table_id` - (Required, ForceNew) The value can get from `alicloud_nat_gateway` Attributes "forward_table_ids".
* `external_ip` - (Required, ForceNew) The external ip address, the ip must along bandwidth package public ip which `alicloud_nat_gateway` argument `bandwidth_packages`.
* `external_port` - (Required) The external port, valid value is 1~65535|any.
* `ip_protocol` - (Required) The ip protocal, valid value is tcp|udp|any.
* `internal_ip` - (Required) The internal ip, must a private ip.
* `internal_port` - (Required) The internal port, valid value is 1~65535|any.