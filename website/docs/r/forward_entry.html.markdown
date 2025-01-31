---
subcategory: "NAT Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_forward_entry"
sidebar_current: "docs-alicloud-resource-vpc"
description: |-
  Provides a Alicloud forward resource.
---

# alicloud_forward_entry

Provides a forward resource.

-> **NOTE:** Available since v1.40.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_forward_entry&exampleId=207dc119-7890-c1c7-0c3b-d22ae8ec57c985b9329e&activeTab=example&spm=docs.r.forward_entry.0.207dc11978&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "forward-entry-example-name"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_nat_gateway" "default" {
  vpc_id               = alicloud_vpc.default.id
  internet_charge_type = "PayByLcu"
  nat_gateway_name     = var.name
  nat_type             = "Enhanced"
  vswitch_id           = alicloud_vswitch.default.id
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

resource "alicloud_eip_association" "default" {
  allocation_id = alicloud_eip_address.default.id
  instance_id   = alicloud_nat_gateway.default.id
}

resource "alicloud_forward_entry" "default" {
  forward_table_id = alicloud_nat_gateway.default.forward_table_ids
  external_ip      = alicloud_eip_address.default.ip_address
  external_port    = "80"
  ip_protocol      = "tcp"
  internal_ip      = "172.16.0.3"
  internal_port    = "8080"
}
```

## Argument Reference

The following arguments are supported:

* `forward_table_id` - (Required, ForceNew) The value can get from `alicloud_nat_gateway` Attributes "forward_table_ids".
* `name` - (Optional, Deprecated) Field `name` has been deprecated from provider version 1.119.1. New field `forward_entry_name` instead.
* `forward_entry_name` - (Optional, Available since 1.119.1+) The name of forward entry.
* `external_ip` - (Required) The external ip address, the ip must along bandwidth package public ip which `alicloud_nat_gateway` argument `bandwidth_packages`.
* `external_port` - (Required) The external port, valid value is 1~65535|any.
* `ip_protocol` - (Required) The ip protocol, valid value is tcp|udp|any.
* `internal_ip` - (Required) The internal ip, must a private ip.
* `internal_port` - (Required) The internal port, valid value is 1~65535|any.
* `port_break` - (Optional, Available since 1.119.1+) Specifies whether to remove limits on the port range. Default value is `false`.

-> **NOTE:** A SNAT entry and a DNAT entry may use the same public IP address. If you want to specify a port number greater than 1024 in this case, set `port_break` to true.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the forward entry. The value formats as `<forward_table_id>:<forward_entry_id>`
* `forward_entry_id` - The id of the forward entry on the server.
* `status` - (Available since 1.119.1+) The status of forward entry.

### Timeouts
-> **NOTE:** Available since 1.119.1+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the forward entry.
* `update` - (Defaults to 10 mins) Used when update the forward entry. 
* `delete` - (Defaults to 10 mins) Used when delete the forward entry. 

## Import

Forward Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_forward_entry.foo ftb-1aece3:fwd-232ce2
```
