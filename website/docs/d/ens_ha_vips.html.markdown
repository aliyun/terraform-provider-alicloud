---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_ha_vips"
sidebar_current: "docs-alicloud-datasource-ens-ha-vips"
description: |-
  Provides a list of ENS Ha Vip owned by an Alibaba Cloud account.
---

# alicloud_ens_ha_vips

This data source provides ENS Ha Vip available to the user.[What is Ha Vip](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateHaVip)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "ens_region_id" {
  default = "cn-hangzhou-58"
}

resource "alicloud_ens_network" "default4wYgcV" {
  network_name  = "tf-example"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultcW3Eib" {
  cidr_block    = "10.0.9.0/24"
  vswitch_name  = "tf-example"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default4wYgcV.id
}


resource "alicloud_ens_ha_vip" "default" {
  description = "desc1"
  vswitch_id  = alicloud_ens_vswitch.defaultcW3Eib.id
  amount      = 1
  ip_address  = "10.0.9.5"
  ha_vip_name = "tf-example"
}

data "alicloud_ens_ha_vips" "default" {
  ids         = ["${alicloud_ens_ha_vip.default.id}"]
  name_regex  = alicloud_ens_ha_vip.default.ha_vip_name
  ha_vip_name = "tf-example"
  ip_address  = "10.0.9.5"
  vswitch_id  = alicloud_ens_vswitch.defaultcW3Eib.id
}

output "alicloud_ens_ha_vip_example_id" {
  value = data.alicloud_ens_ha_vips.default.vips.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ens_region_id` - (Optional) The node ID of ENS.
* `ha_vip_id` - (Optional) The ID of the HaVip.
* `ha_vip_name` - (Optional) The name of the HaVip instance.
* `ip_address` - (Optional) The IP address of the HaVip.
* `network_id` - (Optional) The network ID.
* `status` - (Optional) The HaVip status. Value:
  - Creating: Creating.
  - Available: Available.
  - InUse: In use.
  - Deleting: Deleting.
* `vswitch_id` - (Optional) The vSwitch ID.
* `ids` - (Optional, Computed) A list of HaVip IDs.
* `name_regex` - (Optional) A regex string to filter results by HaVip name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Ha Vip IDs.
* `names` - A list of name of Ha Vips.
* `vips` - A list of Ha Vip Entries. Each element contains the following attributes:
    * `create_time` - The creation time of the resource.
    * `description` - The description of the HaVip instance.
    * `ens_region_id` - The node ID of ENS.
    * `ha_vip_id` - The first ID of the resource.
    * `ha_vip_name` - The name of the HaVip instance.
    * `ip_address` - The IP address of the AVIP.
    * `network_id` - The network ID.
    * `status` - The HaVip status.
    * `vswitch_id` - The vSwitch ID.
    * `id` - The ID of the resource supplied above.
