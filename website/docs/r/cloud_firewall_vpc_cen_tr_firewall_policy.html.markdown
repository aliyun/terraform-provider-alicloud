---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_cen_tr_firewall_policy"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Cen Tr Firewall Policy resource.
---

# alicloud_cloud_firewall_vpc_cen_tr_firewall_policy

Provides a Cloud Firewall Vpc Cen Tr Firewall Policy resource.

VPC border firewall Cloud Enterprise Network Enterprise Edition drainage policy.

For information about Cloud Firewall Vpc Cen Tr Firewall Policy and how to use it, see [What is Vpc Cen Tr Firewall Policy](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/CreateTrFirewallV2RoutePolicy).

-> **NOTE:** Available since v1.268.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "zone4" {
  default = "cn-hangzhou-k"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

variable "zone3" {
  default = "cn-hangzhou-j"
}

resource "alicloud_cen_instance" "cen" {
  description       = "yqc-example"
  cen_instance_name = "yqc-example-CenInstance"
}

resource "alicloud_cen_transit_router" "tr" {
  cen_id              = alicloud_cen_instance.cen.id
  transit_router_name = "yqc-example-TransitRouter"
}

resource "alicloud_express_connect_router_express_connect_router" "ExpressConnectRouter" {
  ecr_name         = "yqc-example-ecr"
  alibaba_side_asn = "64514"
  description      = "22222"
}

resource "alicloud_express_connect_router_tr_association" "ExpressConnectRouterTrAssociation" {
  association_region_id = var.region
  ecr_id                = alicloud_express_connect_router_express_connect_router.ExpressConnectRouter.id
  cen_id                = alicloud_cen_instance.cen.id
  transit_router_id     = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_cen_transit_router_ecr_attachment" "ExpressConnectRouterTrAssociation" {
  ecr_id                                = alicloud_express_connect_router_express_connect_router.ExpressConnectRouter.id
  cen_id                                = alicloud_cen_instance.cen.id
  transit_router_ecr_attachment_name    = "yqc-example-TransitRouterEcrAttachmentName"
  transit_router_attachment_description = "yqc-example-TransitRouterAttachmentDescription"
  transit_router_id                     = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_vpc" "vpc1" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-example-01"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  vpc_id     = alicloud_vpc.vpc1.id
  cidr_block = "172.16.1.0/24"
  zone_id    = var.zone1
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id     = alicloud_vpc.vpc1.id
  cidr_block = "172.16.2.0/24"
  zone_id    = var.zone2
}

resource "alicloud_vpc" "vpc2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-example-02"
}

resource "alicloud_vswitch" "vpc2vsw1" {
  vpc_id     = alicloud_vpc.vpc2.id
  zone_id    = var.zone1
  cidr_block = "172.16.3.0/24"
}

resource "alicloud_vswitch" "vpc2vsw2" {
  vpc_id     = alicloud_vpc.vpc2.id
  cidr_block = "172.16.4.0/24"
  zone_id    = var.zone2
}

resource "alicloud_vpc" "vpc3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-example-03"
}

resource "alicloud_vswitch" "vpc3vsw1" {
  vpc_id     = alicloud_vpc.vpc3.id
  zone_id    = var.zone1
  cidr_block = "172.17.1.0/24"
}

resource "alicloud_vswitch" "vpc3vsw2" {
  vpc_id     = alicloud_vpc.vpc3.id
  cidr_block = "172.17.2.0/24"
  zone_id    = var.zone2
}

resource "alicloud_vpc" "vpc4" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-example-04"
}

resource "alicloud_vswitch" "vpc4vsw1" {
  vpc_id     = alicloud_vpc.vpc4.id
  zone_id    = var.zone1
  cidr_block = "172.16.8.0/24"
}

resource "alicloud_vswitch" "vpc4vsw2" {
  vpc_id     = alicloud_vpc.vpc4.id
  zone_id    = var.zone2
  cidr_block = "172.16.9.0/24"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  vpc_id = alicloud_vpc.vpc1.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
  }
  cen_id                             = alicloud_cen_instance.cen.id
  transit_router_id                  = alicloud_cen_transit_router.tr.transit_router_id
  auto_publish_route_enabled         = true
  transit_router_vpc_attachment_name = "TransitRouterVpcAttachmentName-1"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  auto_publish_route_enabled = true
  vpc_id                     = alicloud_vpc.vpc2.id
  cen_id                     = alicloud_cen_instance.cen.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = alicloud_vswitch.vpc2vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = alicloud_vswitch.vpc2vsw2.zone_id
  }
  transit_router_id                  = alicloud_cen_transit_router.tr.transit_router_id
  transit_router_vpc_attachment_name = "TransitRouterVpcAttachmentName-2"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc3" {
  auto_publish_route_enabled = true
  vpc_id                     = alicloud_vpc.vpc3.id
  cen_id                     = alicloud_cen_instance.cen.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc3vsw1.id
    zone_id    = alicloud_vswitch.vpc3vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc3vsw2.id
    zone_id    = alicloud_vswitch.vpc3vsw2.zone_id
  }
  transit_router_id                  = alicloud_cen_transit_router.tr.transit_router_id
  transit_router_vpc_attachment_name = "TransitRouterVpcAttachmentName-3"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc4" {
  vpc_id = alicloud_vpc.vpc4.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc4vsw1.id
    zone_id    = alicloud_vswitch.vpc4vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc4vsw2.id
    zone_id    = alicloud_vswitch.vpc4vsw2.zone_id
  }
  cen_id                                = alicloud_cen_instance.cen.id
  transit_router_id                     = alicloud_cen_transit_router.tr.transit_router_id
  transit_router_vpc_attachment_name    = "TransitRouterVpcAttachmentName-4"
  auto_publish_route_enabled            = true
  transit_router_attachment_description = "TransitRouterAttachmentDescription4"
}

resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "VpcCenTrFirewall" {
  route_mode                = "managed"
  region_no                 = var.region
  firewall_description      = "VpcCenTrFirewall created by terraform"
  tr_attachment_master_zone = var.zone1
  firewall_name             = "yqc-example-Firewall"
  tr_attachment_master_cidr = "10.0.2.0/24"
  firewall_subnet_cidr      = "10.0.1.0/24"
  cen_id                    = alicloud_cen_instance.cen.id
  tr_attachment_slave_cidr  = "10.0.3.0/24"
  tr_attachment_slave_zone  = var.zone2
  firewall_vpc_cidr         = "10.0.0.0/16"
  transit_router_id         = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_cen_transit_router_route_table" "TransitRouterRouteTable" {
  transit_router_route_table_description = "111"
  transit_router_route_table_name        = "222"
  transit_router_id                      = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation1" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation1" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation2" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc2.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation2" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc2.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation3" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc3.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation3" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc3.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation4" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc4.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation4" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc4.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation5" {
  transit_router_attachment_id  = alicloud_cen_transit_router_ecr_attachment.ExpressConnectRouterTrAssociation.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation5" {
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
  transit_router_attachment_id  = alicloud_cen_transit_router_ecr_attachment.ExpressConnectRouterTrAssociation.id
}


resource "alicloud_cloud_firewall_vpc_cen_tr_firewall_policy" "default" {
  src_candidate_list {
    candidate_id   = alicloud_express_connect_router_express_connect_router.ExpressConnectRouter.id
    candidate_type = "ECR"
  }
  src_candidate_list {
    candidate_id   = alicloud_vpc.vpc1.id
    candidate_type = "VPC"
  }
  policy_type        = "fullmesh"
  policy_description = "111111"
  firewall_id        = alicloud_cloud_firewall_vpc_cen_tr_firewall.VpcCenTrFirewall.id
  policy_name        = "222222"
}
```

## Argument Reference

The following arguments are supported:
* `dest_candidate_list` - (Optional, Set) List of Secondary Traffic Redirection instances. See [`dest_candidate_list`](#dest_candidate_list) below.
* `firewall_id` - (Required, ForceNew) The ID of the VPC firewall instance.
* `policy_description` - (Required, ForceNew) Traffic Redirection description.
* `policy_name` - (Required, ForceNew) Traffic Redirection Template Name.
* `policy_type` - (Required, ForceNew) VPC boundary firewall cloud enterprise edition drainage scenario type. Value:
  - `fullmesh`: Multipoint Interconnection
  - `One_to_one`: Point to Point
  - `End_to_end`: point to multipoint
* `should_recover` - (Optional) Whether to restore the drainage configuration. Value:
  - true: Route Rollback
  - false: Route withdrawal

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `src_candidate_list` - (Required, Set) List of Primary Traffic Redirection instances. See [`src_candidate_list`](#src_candidate_list) below.
* `status` - (Optional, Computed) The policy state. Value:
  - creating: creating
  - deleting: deleting
  - opening: opening
  - opened: opened
  - closing: closing
  - closed: closed

### `dest_candidate_list`

The dest_candidate_list supports the following:
* `candidate_id` - (Optional) The ID of the Traffic Redirection instance.
* `candidate_type` - (Optional) The Traffic Redirection instance type.

### `src_candidate_list`

The src_candidate_list supports the following:
* `candidate_id` - (Optional) The ID of the Traffic Redirection instance.
* `candidate_type` - (Optional) The Traffic Redirection instance type.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<firewall_id>:<tr_firewall_route_policy_id>`.
* `tr_firewall_route_policy_id` - The ID of the firewall routing policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 16 mins) Used when create the Vpc Cen Tr Firewall Policy.
* `delete` - (Defaults to 16 mins) Used when delete the Vpc Cen Tr Firewall Policy.
* `update` - (Defaults to 16 mins) Used when update the Vpc Cen Tr Firewall Policy.

## Import

Cloud Firewall Vpc Cen Tr Firewall Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_vpc_cen_tr_firewall_policy.example <firewall_id>:<tr_firewall_route_policy_id>
```