---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_cen_tr_firewall"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Cen Tr Firewall resource.
---

# alicloud_cloud_firewall_vpc_cen_tr_firewall

Provides a Cloud Firewall Vpc Cen Tr Firewall resource.

VPC firewall Cloud Enterprise Network Enterprise Edition.

For information about Cloud Firewall Vpc Cen Tr Firewall and how to use it, see [What is Vpc Cen Tr Firewall](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_vpc_cen_tr_firewall&exampleId=37a1e48f-5e2b-3aa6-5c58-4d9edb62f71ef8f83c50&activeTab=example&spm=docs.r.cloud_firewall_vpc_cen_tr_firewall.0.37a1e48f5e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "description" {
  default = "Created by Terraform"
}

variable "firewall_name" {
  default = "tf-example"
}

variable "tr_attachment_master_cidr" {
  default = "192.168.3.192/26"
}

variable "firewall_subnet_cidr" {
  default = "192.168.3.0/25"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "tr_attachment_slave_cidr" {
  default = "192.168.3.128/26"
}

variable "firewall_vpc_cidr" {
  default = "192.168.3.0/24"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "firewall_name_update" {
  default = "tf-example-1"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_cen_instance" "cen" {
  description       = "terraform example"
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "tr" {
  transit_router_name        = var.name
  transit_router_description = "tr-created-by-terraform"
  cen_id                     = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  description = "created by terraform"
  cidr_block  = "192.168.1.0/24"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vpc1vsw1" {
  cidr_block   = "192.168.1.0/25"
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.vpc1.id
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id       = alicloud_vpc.vpc1.id
  cidr_block   = "192.168.1.128/26"
  vswitch_name = var.name
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.vpc1.id
  route_table_name = var.name
  description      = var.name
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
  }
  zone_mappings {
    zone_id    = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
  }
  vpc_id            = alicloud_vpc.vpc1.id
  cen_id            = alicloud_cen_instance.cen.id
  transit_router_id = alicloud_cen_transit_router.tr.transit_router_id
  depends_on        = [alicloud_route_table.foo]
}

resource "time_sleep" "wait_10_minutes" {
  depends_on = [alicloud_cen_transit_router_vpc_attachment.tr-vpc1]

  create_duration = "10m"
}

resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "default" {
  cen_id                    = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.cen_id
  firewall_name             = var.name
  firewall_subnet_cidr      = var.firewall_subnet_cidr
  tr_attachment_slave_cidr  = var.tr_attachment_slave_cidr
  firewall_description      = "VpcCenTrFirewall created by terraform"
  region_no                 = var.region
  tr_attachment_master_cidr = var.tr_attachment_master_cidr
  firewall_vpc_cidr         = var.firewall_vpc_cidr
  transit_router_id         = alicloud_cen_transit_router.tr.transit_router_id
  route_mode                = "managed"

  depends_on = [time_sleep.wait_10_minutes]
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `firewall_description` - (Optional, ForceNew) Firewall description.
* `firewall_name` - (Required) The name of Cloud Firewall.
* `firewall_subnet_cidr` - (Required, ForceNew) Required in automatic mode, the CIDR of subnet used to store the firewall ENI in the firewall VPC.
* `firewall_vpc_cidr` - (Required, ForceNew) Required in automatic mode,  th CIDR of firewall VPC.
* `region_no` - (Required, ForceNew) The region ID of the transit router instance.
* `route_mode` - (Required, ForceNew) The routing pattern. Value: managed: indicates automatic mode

* `tr_attachment_master_cidr` - (Required, ForceNew) Required in automatic mode, the primary CIDR of network used to connect to the TR in the firewall VPC.
* `tr_attachment_master_zone` - (Optional) The primary zone of the switch.

* `tr_attachment_slave_cidr` - (Required, ForceNew) Required in automatic mode, the the secondary CIDR of the subnet in the firewall VPC used to connect to TR.
* `tr_attachment_slave_zone` - (Optional) Switch standby area.

* `transit_router_id` - (Required, ForceNew) The ID of the transit router instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - Firewall status. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Cen Tr Firewall.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Cen Tr Firewall.
* `update` - (Defaults to 5 mins) Used when update the Vpc Cen Tr Firewall.

## Import

Cloud Firewall Vpc Cen Tr Firewall can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_vpc_cen_tr_firewall.example <id>
```