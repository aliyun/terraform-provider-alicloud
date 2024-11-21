---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_nat_firewall"
description: |-
  Provides a Alicloud Cloud Firewall Nat Firewall resource.
---

# alicloud_cloud_firewall_nat_firewall

Provides a Cloud Firewall Nat Firewall resource. 

For information about Cloud Firewall Nat Firewall and how to use it, see [What is Nat Firewall](https://www.alibabacloud.com/help/zh/cloud-firewall/developer-reference/api-cloudfw-2017-12-07-createsecurityproxy).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_nat_firewall&exampleId=19d83e5d-c8c8-684a-59b0-b0b5f4d5458dab48e0cc&activeTab=example&spm=docs.r.cloud_firewall_nat_firewall.0.19d83e5dc8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shenzhen"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultikZ0gD" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultp4O7qi" {
  vpc_id       = alicloud_vpc.defaultikZ0gD.id
  cidr_block   = "172.16.6.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_nat_gateway" "default2iRZpC" {
  eip_bind_mode    = "MULTI_BINDED"
  vpc_id           = alicloud_vpc.defaultikZ0gD.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.defaultp4O7qi.id
  nat_type         = "Enhanced"
  network_type     = "internet"
}

resource "alicloud_eip_address" "defaultyiRwgs" {
  address_name = var.name
}

resource "alicloud_eip_association" "defaults2MTuO" {
  instance_id   = alicloud_nat_gateway.default2iRZpC.id
  allocation_id = alicloud_eip_address.defaultyiRwgs.id
  mode          = "NAT"
  instance_type = "Nat"
}

resource "alicloud_snat_entry" "defaultAKE43g" {
  snat_ip           = alicloud_eip_address.defaultyiRwgs.ip_address
  snat_table_id     = alicloud_nat_gateway.default2iRZpC.snat_table_ids
  source_vswitch_id = alicloud_vswitch.defaultp4O7qi.id
}

resource "alicloud_cloud_firewall_nat_firewall" "default" {
  nat_gateway_id = alicloud_nat_gateway.default2iRZpC.id
  nat_route_entry_list {
    nexthop_type     = "NatGateway"
    route_table_id   = alicloud_vpc.defaultikZ0gD.route_table_id
    nexthop_id       = alicloud_nat_gateway.default2iRZpC.id
    destination_cidr = "0.0.0.0/0"
  }

  firewall_switch = "close"
  vswitch_auto    = "true"
  status          = "closed"
  region_no       = "cn-shenzhen"
  lang            = "zh"
  proxy_name      = "CFW-example"
  vswitch_id      = alicloud_snat_entry.defaultAKE43g.source_vswitch_id
  strict_mode     = "0"
  vpc_id          = alicloud_vpc.defaultikZ0gD.id
  vswitch_cidr    = "172.16.5.0/24"
}
```

## Argument Reference

The following arguments are supported:
* `firewall_switch` - (Optional) Safety protection switch. Value:-**open**: open-**close**: close.
* `lang` - (Optional) Lang.
* `nat_gateway_id` - (Required, ForceNew) NAT gateway ID.
* `nat_route_entry_list` - (Required, ForceNew) The list of routes to be switched by the NAT gateway. See [`nat_route_entry_list`](#nat_route_entry_list) below.
* `proxy_name` - (Required, ForceNew) NAT firewall name.
* `region_no` - (Required, ForceNew) Region.
* `status` - (Optional, Computed) The status of the resource.
* `strict_mode` - (Optional, ForceNew) Whether strict mode is enabled 1-Enable strict mode 0-Disable strict mode.
* `vswitch_id` - (Optional) The switch ID. Required for switch manual mode.
* `vpc_id` - (Required, ForceNew) The ID of the VPC instance.
* `vswitch_auto` - (Optional) Whether to use switch automatic mode. Value: **true**: Use automatic mode: **false**: Use manual mode.
* `vswitch_cidr` - (Optional) The network segment of the virtual switch. Required for Switch automatic mode.

### `nat_route_entry_list`

The nat_route_entry_list supports the following:
* `destination_cidr` - (Required, ForceNew) The destination network segment of the default route.
* `nexthop_id` - (Required, ForceNew) The next hop address of the original NAT gateway.
* `nexthop_type` - (Required, ForceNew) The network type of the next hop. Value: NatGateway : NAT Gateway.
* `route_table_id` - (Required, ForceNew) The route table where the default route of the NAT gateway is located.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Nat Firewall.
* `delete` - (Defaults to 5 mins) Used when delete the Nat Firewall.
* `update` - (Defaults to 5 mins) Used when update the Nat Firewall.

## Import

Cloud Firewall Nat Firewall can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_nat_firewall.example <id>
```