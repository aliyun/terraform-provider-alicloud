---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_cen_tr_firewalls"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-vpc-cen-tr-firewalls"
description: |-
  Provides a list of Cloud Firewall Vpc Cen Tr Firewall owned by an Alibaba Cloud account.
---

# alicloud_cloud_firewall_vpc_cen_tr_firewalls

This data source provides Cloud Firewall Vpc Cen Tr Firewall available to the user.[What is Vpc Cen Tr Firewall](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/CreateTrFirewallV2)

-> **NOTE:** Available since v1.243.0.

## Example Usage

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

resource "alicloud_cen_instance" "cen" {
  description       = "terraform example"
  cen_instance_name = "Cen_Terraform_example01"
}

resource "alicloud_cen_transit_router" "tr" {
  support_multicast          = false
  transit_router_name        = "CEN_TR_Terraform"
  transit_router_description = "tr-created-by-terraform"
  cen_id                     = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  description = "created by terraform"
  cidr_block  = "192.168.1.0/24"
  vpc_name    = "vpc1-Terraform"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  cidr_block   = "192.168.1.0/25"
  vswitch_name = "vpc1-vsw1"
  vpc_id       = alicloud_vpc.vpc1.id
  zone_id      = var.zone1
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id       = alicloud_vpc.vpc1.id
  cidr_block   = "192.168.1.128/26"
  vswitch_name = "vpc1-vsw2"
  zone_id      = var.zone2
}

resource "alicloud_vpc" "vpc2" {
  description = "created by terraform"
  cidr_block  = "192.168.2.0/24"
  vpc_name    = "vpc2-Terraform"
}

resource "alicloud_vswitch" "vpc2vsw1" {
  cidr_block   = "192.168.2.0/25"
  vswitch_name = "vpc2-vsw1"
  vpc_id       = alicloud_vpc.vpc2.id
  zone_id      = var.zone1
}

resource "alicloud_vswitch" "vpc2vsw2" {
  cidr_block   = "192.168.2.128/26"
  vswitch_name = "vpc2-vsw2"
  vpc_id       = alicloud_vpc.vpc2.id
  zone_id      = var.zone2
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  auto_publish_route_enabled = false
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
  }
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  auto_publish_route_enabled = false
  vpc_id                     = alicloud_vpc.vpc2.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = alicloud_vswitch.vpc2vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = alicloud_vswitch.vpc2vsw2.zone_id
  }
  cen_id = alicloud_cen_instance.cen.id
}


resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "default" {
  firewall_description      = "VpcCenTrFirewall created by terraform"
  region_no                 = var.region
  route_mode                = "managed"
  cen_id                    = alicloud_cen_instance.cen.id
  firewall_vpc_cidr         = var.firewall_vpc_cidr
  transit_router_id         = alicloud_cen_transit_router.tr.transit_router_id
  tr_attachment_master_cidr = var.tr_attachment_master_cidr
  firewall_name             = var.firewall_name
  firewall_subnet_cidr      = var.firewall_subnet_cidr
  tr_attachment_slave_cidr  = var.tr_attachment_slave_cidr
}

data "alicloud_cloud_firewall_vpc_cen_tr_firewalls" "default" {
  ids               = ["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}"]
  cen_id            = alicloud_cen_instance.cen.id
  firewall_name     = var.firewall_name
  region_no         = var.region
  route_mode        = "managed"
  transit_router_id = alicloud_cen_transit_router.tr.transit_router_id
}

output "alicloud_cloud_firewall_vpc_cen_tr_firewall_example_id" {
  value = data.alicloud_cloud_firewall_vpc_cen_tr_firewalls.default.firewalls.0.id
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (ForceNew, Optional) The ID of the CEN instance.
* `current_page` - (ForceNew, Optional) The page number of the pagination query. The default value is 1.
* `firewall_id` - (ForceNew, Optional) Firewall ID
* `firewall_name` - (ForceNew, Optional) The name of Cloud Firewall.
* `firewall_switch_status` - (ForceNew, Optional) The status of the VPC boundary firewall. Value:-**opened**: opened-**closed**: closed-**notconfigured**: indicates that the VPC boundary firewall has not been configured yet.-**configured**: indicates that the VPC boundary firewall has been configured.-**creating**: indicates that a VPC boundary firewall is being created.-**opening**: indicates that the VPC border firewall is being enabled.-**deleting**: indicates that the VPC boundary firewall is being deleted.> If this parameter is not set, the VPC boundary firewall in all states is queried.
* `page_number` - (ForceNew, Optional) Current page number.
* `page_size` - (ForceNew, Optional) The maximum number of pieces of data per page that are displayed during a paged query. The default value is 10.
* `region_no` - (ForceNew, Optional) The region ID of the transit router instance.
* `route_mode` - (ForceNew, Optional) The routing pattern. Value: managed: indicates automatic mode
* `transit_router_id` - (ForceNew, Optional) The ID of the transit router instance.
* `ids` - (Optional, ForceNew, Computed) A list of Vpc Cen Tr Firewall IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Vpc Cen Tr Firewall IDs.
* `firewalls` - A list of Vpc Cen Tr Firewall Entries. Each element contains the following attributes:
  * `cen_id` - The ID of the CEN instance.
  * `cen_name` - The name of the CEN instance.
  * `firewall_id` - Firewall ID
  * `firewall_name` - The name of Cloud Firewall.
  * `firewall_switch_status` - The status of the VPC boundary firewall. Value:-**opened**: opened-**closed**: closed-**notconfigured**: indicates that the VPC boundary firewall has not been configured yet.-**configured**: indicates that the VPC boundary firewall has been configured.-**creating**: indicates that a VPC boundary firewall is being created.-**opening**: indicates that the VPC border firewall is being enabled.-**deleting**: indicates that the VPC boundary firewall is being deleted.> If this parameter is not set, the VPC boundary firewall in all states is queried.
  * `ips_config` - IPS configuration information.
    * `basic_rules` - Basic rule switch. Value:-**1**: On-**0**: Closed state.
    * `enable_all_patch` - Virtual patch switch. Value:-**1**: On-**0**: Closed state.
    * `run_mode` - IPS defense mode. Value:-**1**: Intercept mode-**0**: Observation mode.
  * `precheck_status` - Whether the wall can be opened automatically. Value:-**passed**: can automatically open the wall-**failed**: The wall cannot be opened automatically-**unknown**: unknown status
  * `region_no` - The region ID of the transit router instance.
  * `region_status` - Geographically open. Value:-**enable**: enabled, indicating that the VPC border firewall can be configured for the region.-**disable**: Not enabled, indicating that the VPC boundary firewall is not allowed for the region.
  * `result_code` - The operation result code of creating the VPC boundary firewall. Value:-**RegionDisable**: indicates that the region where the network instance is located is not supported by the VPC border firewall. You cannot create a VPC border firewall.-**Empty string**, indicating that the network instance can create a VPC firewall.
  * `route_mode` - The routing pattern. Value: managed: indicates automatic mode
  * `transit_router_id` - The ID of the transit router instance.
  * `id` - The ID of the resource supplied above.
