---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_cen_tr_firewall"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Cen Tr Firewall resource.
---

# alicloud_cloud_firewall_vpc_cen_tr_firewall

Provides a Cloud Firewall Vpc Cen Tr Firewall resource. VPC firewall Cloud Enterprise Network Enterprise Edition.

For information about Cloud Firewall Vpc Cen Tr Firewall and how to use it, see [What is Vpc Cen Tr Firewall](https://www.alibabacloud.com/help/en/cloud-firewall/latest/vpc-firewall-limits).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "zone2" {
  default = "cn-hangzhou-i"
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
  zone_id      = var.zone1
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id       = alicloud_vpc.vpc1.id
  cidr_block   = "192.168.1.128/26"
  vswitch_name = var.name
  zone_id      = var.zone2
}

resource "alicloud_vpc" "vpc2" {
  description = "created by terraform"
  cidr_block  = "192.168.2.0/24"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vpc2vsw1" {
  cidr_block   = "192.168.2.0/25"
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.vpc2.id
  zone_id      = var.zone1
}

resource "alicloud_vswitch" "vpc2vsw2" {
  cidr_block   = "192.168.2.128/26"
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.vpc2.id
  zone_id      = var.zone2
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = var.zone1
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
    zone_id    = var.zone2
  }
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_transit_router.tr.cen_id
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  vpc_id = alicloud_vpc.vpc2.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = var.zone1
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = var.zone2
  }
  cen_id = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.cen_id
}


resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "default" {
  firewall_description      = "VpcCenTrFirewall created by terraform"
  region_no                 = "cn-hangzhou"
  route_mode                = "managed"
  cen_id                    = alicloud_cen_instance.cen.cen_id
  firewall_vpc_cidr         = "192.168.3.0/24"
  transit_router_id         = alicloud_cen_transit_router.tr.id
  tr_attachment_master_cidr = "192.168.3.192/26"
  firewall_name             = var.name
  firewall_subnet_cidr      = "192.168.3.0/25"
  tr_attachment_slave_cidr  = "192.168.3.128/26"
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
* `route_mode` - (Required, ForceNew) The routing pattern. Value: managed: indicates automatic mode.
* `tr_attachment_master_cidr` - (Required, ForceNew) Required in automatic mode, the primary CIDR of network used to connect to the TR in the firewall VPC.
* `tr_attachment_slave_cidr` - (Required, ForceNew) Required in automatic mode, the the secondary CIDR of the subnet in the firewall VPC used to connect to TR.
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