---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_cen_tr_firewall"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Cen Tr Firewall resource.
---

# alicloud_cloud_firewall_vpc_cen_tr_firewall

Provides a Cloud Firewall Vpc Cen Tr Firewall resource.

VPC Firewall for Cloud Enterprise Network Enterprise Edition.

For information about Cloud Firewall Vpc Cen Tr Firewall and how to use it, see [What is Vpc Cen Tr Firewall](https://www.alibabacloud.com/help/en/cloud-firewall/cloudfirewall/developer-reference/api-cloudfw-2017-12-07-createtrfirewallv2).

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

variable "region" {
  default = "cn-hangzhou"
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
  firewall_description      = "VpcCenTrFirewall created by terraform"
  region_no                 = var.region
  firewall_vpc_cidr         = var.firewall_vpc_cidr
  transit_router_id         = alicloud_cen_transit_router.tr.transit_router_id
  route_mode                = "managed"

  depends_on = [time_sleep.wait_10_minutes]
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_vpc_cen_tr_firewall&spm=docs.r.cloud_firewall_vpc_cen_tr_firewall.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `cen_id` - (Required, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance.
* `firewall_description` - (Optional, ForceNew) The description of the firewall.
* `firewall_name` - (Required) The name of the Cloud Firewall instance.
* `firewall_vpc_cidr` - (Optional, ForceNew) The CIDR block of the firewall VPC in automatic mode.
* `firewall_vpc_id` - (Optional, ForceNew, Available since v2.2.0) The ID of the VPC where the firewall ENI is created in manual mode.
* `firewall_vswitch_id` - (Optional, ForceNew, Available since v2.2.0) The ID of the vSwitch where the firewall ENI is created in manual mode.
* `region_no` - (Required, ForceNew) The region ID of transit router instance I.
* `route_mode` - (Required, ForceNew) The routing mode. Valid values:
  - `managed`: automatic mode. You only need to specify `firewall_vpc_cidr`. The `firewall_subnet_cidr`, `tr_attachment_slave_cidr`, and `tr_attachment_master_cidr` are deprecated and do not need to be specified.
  - `manual`: (Available since v2.2.0) manual mode for custom VPCs. You must specify `firewall_vpc_id`, `firewall_vswitch_id`, `tr_attachment_master_zone`, and `tr_attachment_slave_zone`.
  - `semi-managed`: (Available since v2.2.0) manual mode for managed VPCs. You only need to specify `firewall_vpc_cidr`.

  -> **NOTE:**  If this parameter is not specified, VPC border firewalls in all routing modes are queried.

* `tr_attachment_master_zone` - (Optional, ForceNew) The primary zone of the vSwitch.
* `tr_attachment_slave_zone` - (Optional, ForceNew) The secondary zone of the vSwitch.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router instance.
* `firewall_subnet_cidr` - (Optional, ForceNew, Deprecated since v2.2.0) Field `firewall_subnet_cidr` has been deprecated from provider version 2.2.0.
* `tr_attachment_master_cidr` - (Optional, ForceNew, Deprecated since v2.2.0) Field `tr_attachment_master_cidr` has been deprecated from provider version 2.2.0.
* `tr_attachment_slave_cidr` - (Optional, ForceNew, Deprecated since v2.2.0) Field `tr_attachment_slave_cidr` has been deprecated from provider version 2.2.0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `firewall_eni_id` - The ID of the firewall ENI.
* `firewall_eni_vpc_id` - The ID of the VPC where the firewall ENI resides.
* `firewall_vpc_attachment_id` - The attachment ID used in the firewall VPC to connect to the transit router (TR) in automatic mode.
* `status` - The status of the firewall.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 41 mins) Used when create the Vpc Cen Tr Firewall.
* `delete` - (Defaults to 46 mins) Used when delete the Vpc Cen Tr Firewall.
* `update` - (Defaults to 5 mins) Used when update the Vpc Cen Tr Firewall.

## Import

Cloud Firewall Vpc Cen Tr Firewall can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_vpc_cen_tr_firewall.example <firewall_id>
```
