---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachments"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-vpc-attachments"
description: |-
  Provides a list of CEN Transit Router VPC Attachments to the user.
---

# alicloud_cen_transit_router_vpc_attachments

This data source provides the CEN Transit Router VPC Attachments of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}

data "alicloud_vswitches" "default_master" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.1
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id                                = alicloud_cen_instance.default.id
  vpc_id                                = data.alicloud_vpcs.default.ids.0
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
  zone_mappings {
    vswitch_id = data.alicloud_vswitches.default_master.vswitches.0.id
    zone_id    = data.alicloud_vswitches.default_master.vswitches.0.zone_id
  }
  zone_mappings {
    vswitch_id = data.alicloud_vswitches.default.vswitches.0.id
    zone_id    = data.alicloud_vswitches.default.vswitches.0.zone_id
  }
}

data "alicloud_cen_transit_router_vpc_attachments" "ids" {
  ids    = [alicloud_cen_transit_router_vpc_attachment.default.id]
  cen_id = alicloud_cen_instance.default.id
}

output "cen_transit_router_vpc_attachments_id_0" {
  value = data.alicloud_cen_transit_router_vpc_attachments.ids.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Transit Router VPC Attachment IDs.
* `name_regex` - (Optional, ForceNew, Available since v1.224.0) A regex string to filter results by Transit Router VPC Attachment name.
* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `vpc_id` - (Optional, ForceNew, Available since v1.224.0) The ID of the VPC.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `transit_router_attachment_id` - (Optional, ForceNew, Available since v1.224.0) The ID of the Transit Router VPC Attachment.
* `status` - (Optional, ForceNew) The status of the Transit Router VPC Attachment. Valid Values: `Attached`, `Attaching`, `Detaching`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Transit Router VPC Attachment names.
* `attachments` - A list of Transit Router VPC Attachments. Each element contains the following attributes:
  * `id` - The resource ID in terraform of Transit Router VPC Attachment. It formats as `<cen_id>:<transit_router_attachment_id>`.
  * `cen_id` - (Available since v1.224.0) The ID of the CEN instance.
  * `transit_router_attachment_id` - The ID of the Transit Router VPC Attachment.
  * `vpc_id` - The ID of the VPC.
  * `transit_router_id` - (Available since v1.224.0) The ID of the transit router.
  * `resource_type` - The resource type of the Transit Router VPC Attachment.
  * `payment_type` - The payment type of the resource.
  * `vpc_owner_id` - The Owner ID of the VPC.
  * `auto_publish_route_enabled` - (Available since v1.224.0) Whether the transit router is automatically published to the VPC instance.
  * `transit_router_attachment_name` - The name of the Transit Router VPC Attachment.
  * `transit_router_attachment_description` - The description of the Transit Router VPC Attachment.
  * `status` - The status of the Transit Router VPC Attachment.
  * `zone_mappings` - The list of zone mapping of the VPC.
    * `vswitch_id` - The ID of the vSwitch.
    * `zone_id` - The ID of the zone.
