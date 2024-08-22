---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_vpc_attachment"
description: |-
  Provides a Alicloud CEN Transit Router VPC Attachment resource.
---

# alicloud_cen_transit_router_vpc_attachment

Provides a CEN Transit Router VPC Attachment resource that associate the VPC with the CEN instance. [What is Cen Transit Router VPC Attachment](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitroutervpcattachment)

-> **NOTE:** Available since v1.126.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cen_transit_router_vpc_attachment&exampleId=509f440a-f327-0f3f-84dd-e67b0264a3307907872b&activeTab=example&spm=docs.r.cen_transit_router_vpc_attachment.0.509f440af3" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

locals {
  master_zone = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[0]
  slave_zone  = data.alicloud_cen_transit_router_available_resources.default.resources[0].slave_zones[1]
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "example_master" {
  vswitch_name = var.name
  cidr_block   = "192.168.1.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.master_zone
}

resource "alicloud_vswitch" "example_slave" {
  vswitch_name = var.name
  cidr_block   = "192.168.2.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.slave_zone
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.example.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "example" {
  cen_id            = alicloud_cen_instance.example.id
  transit_router_id = alicloud_cen_transit_router.example.transit_router_id
  vpc_id            = alicloud_vpc.example.id
  zone_mappings {
    zone_id    = local.master_zone
    vswitch_id = alicloud_vswitch.example_master.id
  }
  zone_mappings {
    zone_id    = local.slave_zone
    vswitch_id = alicloud_vswitch.example_slave.id
  }
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `resource_type` - (Optional, ForceNew) The resource type of the transit router vpc attachment. Default value: `VPC`. Valid values: `VPC`.
* `payment_type` - (Optional, ForceNew, Available since v1.168.0) The payment type of the resource. Default value: `PayAsYouGo`. Valid values: `PayAsYouGo`.
* `vpc_owner_id` - (Optional, ForceNew) The owner id of vpc.
* `auto_publish_route_enabled` - (Optional, Bool, Available since v1.204.0) Whether the transit router is automatically published to the VPC instance. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `transit_router_attachment_name` - (Optional) The name of the transit router vbr attachment.
* `transit_router_attachment_description` - (Optional) The description of the transit router vbr attachment.
* `zone_mappings` - (Required, Set) The list of zone mapping of the VPC. See [`zone_mappings`](#zone_mappings) below. **NOTE:** From version 1.184.0, `zone_mappings` can be modified.
-> **NOTE:** The Zone of CEN has MasterZone and SlaveZone, first zone_id of zone_mapping need be MasterZone. We have a API to describeZones[API](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-listtransitrouteravailableresource)
* `tags` - (Optional, Available since v1.193.1) A mapping of tags to assign to the resource.
* `dry_run` - (Optional, Bool) The dry run.
* `route_table_association_enabled` - (Optional, Bool, Deprecated since v1.192.0) Whether to enabled route table association. **NOTE:** "Field `route_table_association_enabled` has been deprecated from provider version 1.192.0. Please use the resource `alicloud_cen_transit_router_route_table_association` instead, [how to use alicloud_cen_transit_router_route_table_association](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_route_table_association)."
* `route_table_propagation_enabled` - (Optional, Bool, Deprecated since v1.192.0) Whether to enabled route table propagation. **NOTE:** "Field `route_table_propagation_enabled` has been deprecated from provider version 1.192.0. Please use the resource `alicloud_cen_transit_router_route_table_propagation` instead, [how to use alicloud_cen_transit_router_route_table_propagation](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_route_table_propagation)."

### `zone_mappings`

The zone_mappings supports the following:

* `vswitch_id` - (Optional, ForceNew) The VSwitch id of attachment.
* `zone_id` - (Optional, ForceNew) The zone Id of VSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router VPC Attachment. It formats as `<cen_id>:<transit_router_attachment_id>`.
* `transit_router_attachment_id` - The ID of the Transit Router Attachment.
* `status` - The associating status of the network.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when creating the cen transit router vpc attachment (until it reaches the initial `Attached` status).
* `update` - (Defaults to 3 mins) Used when update the cen transit router vpc attachment.
* `delete` - (Defaults to 3 mins) Used when delete the cen transit router vpc attachment.

## Import

CEN Transit Router VPC Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vpc_attachment.example <cen_id>:<transit_router_attachment_id>
```
