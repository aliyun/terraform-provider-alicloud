---
subcategory: "NAT Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_snat_entry"
description: |-
  Provides a Alicloud snat resource.
---

# alicloud_snat_entry

Provides a NAT Gateway Snat Entry resource.



For information about NAT Gateway Snat Entry and how to use it, see [What is Snat Entry](https://www.alibabacloud.com/help/en/nat-gateway/developer-reference/api-vpc-2016-04-28-createsnatentry-natgws).

-> **NOTE:** Available since v1.119.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_snat_entry&exampleId=03cb9142-7fb8-29c9-bf5d-497c893c5b759b577a6a&activeTab=example&spm=docs.r.snat_entry.0.03cb91427f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
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
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_nat_gateway" "default" {
  vpc_id           = alicloud_vpc.default.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.default.id
  nat_type         = "Enhanced"
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

resource "alicloud_eip_association" "default" {
  allocation_id = alicloud_eip_address.default.id
  instance_id   = alicloud_nat_gateway.default.id
}

resource "alicloud_snat_entry" "default" {
  snat_table_id     = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = alicloud_vswitch.default.id
  snat_ip           = alicloud_eip_address.default.ip_address
}
```

## Argument Reference

The following arguments are supported:
* `eip_affinity` - (Optional, Int, Available since v1.241.0) Specifies whether to enable EIP affinity. Default value: `0`. Valid values:
  - `0`: Disable.
  - `1`: Enable.
* `snat_entry_name` - (Optional, Available since v1.71.2) The name of the SNAT entry. The name must be `2` to `128` characters in length. It must start with a letter but cannot start with `http://` or `https://`.
* `snat_ip` - (Required) The IP of a SNAT entry. Separate multiple EIP or NAT IP addresses with commas (,). **NOTE:** From version 1.241.0, `snat_ip` can be modified.
* `snat_table_id` - (Required, ForceNew) The ID of the SNAT table.
* `source_cidr` - (Optional, ForceNew, Available since v1.71.1) The source CIDR block specified in the SNAT entry.
* `source_vswitch_id` - (Optional, ForceNew) The ID of the vSwitch.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Snat Entry. It formats as `<snat_table_id>:<snat_entry_id>`.
-> **NOTE:** Before provider version 1.37.0, it formats as `<snat_entry_id>`
* `snat_entry_id` - The id of the snat entry on the server.
* `status` - (Available since v1.119.1) The ID of the SNAT entry.

## Timeouts

-> **NOTE:** Available since v1.119.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Snat Entry.
* `delete` - (Defaults to 5 mins) Used when delete the Snat Entry.
* `update` - (Defaults to 5 mins) Used when update the Snat Entry.

## Import

NAT Gateway Snat Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_snat_entry.example <snat_table_id>:<snat_entry_id>
```

**NOTE:** Before provider version 1.37.0, NAT Gateway Snat Entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_snat_entry.example <snat_entry_id>
```
