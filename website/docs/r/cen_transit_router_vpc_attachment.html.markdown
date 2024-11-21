---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachment"
description: |-
  Provides a Alicloud CEN Transit Router Vpc Attachment resource.
---

# alicloud_cen_transit_router_vpc_attachment

Provides a CEN Transit Router VPC Attachment resource that associate the VPC with the CEN instance. [What is Cen Transit Router VPC Attachment](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitroutervpcattachment)

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_vpc_attachment&exampleId=509f440a-f327-0f3f-84dd-e67b0264a3307907872b&activeTab=example&spm=docs.r.cen_transit_router_vpc_attachment.0.509f440af3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `auto_publish_route_enabled` - (Optional) Specifies whether to enable the Enterprise Edition transit router to automatically advertise routes to VPCs. Valid values:
  - **false:** (default)
  - `true`

* `cen_id` - (Optional) The ID of the Cloud Enterprise Network (CEN) instance.

* `dry_run` - (Optional) Whether to perform PreCheck on this request, including permissions and instance status verification. Value:
  - `false` (default): A normal request is sent, and a VPC connection is directly created after the check is passed.
  - `true`: The check request is sent, only verification is performed, and no VPC connection is created. Check items include whether required parameters and request format are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
* `force_delete` - (Optional, Available since v1.230.1) Whether to forcibly delete the VPC connection. The value is:
  - `false` (default): before deleting the VPC connection, check whether there are related resource dependencies, such as Association forwarding and route learning. If related dependencies exist, deletion is not allowed and the corresponding error is returned.
  - `true`: When you delete a VPC connection, all related dependencies are deleted by default.
* `payment_type` - (Optional, ForceNew, Computed) The billing method. The default value is `PayAsYouGo`, which specifies the pay-as-you-go billing method.

* `tags` - (Optional, Map) The tag of the resource
* `transit_router_attachment_description` - (Optional) The description of the VPC connection.

  The description must be 2 to 256 characters in length. The description must start with a letter but cannot start with `http://` or `https://`.

* `transit_router_id` - (Optional, ForceNew) The ID of the Enterprise Edition transit router.

* `transit_router_vpc_attachment_name` - (Optional, Available since v1.230.1) The name of the VPC connection.

  The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (\_), and hyphens (-). It must start with a letter.

* `transit_router_vpc_attachment_options` - (Optional, Map, Available since v1.230.1) TransitRouterVpcAttachmentOptions
* `vpc_id` - (Required, ForceNew) The VPC ID.

* `vpc_owner_id` - (Optional, ForceNew, Computed, Int) VpcOwnerId
* `resource_type` - (Optional, ForceNew) The resource type of the transit router vpc attachment. Default value: `VPC`. Valid values: `VPC`.
* `zone_mappings` - (Required, Set) ZoneMappingss See [`zone_mappings`](#zone_mappings) below.

The following arguments will be discarded. Please use new fields as soon as possible:
* `transit_router_attachment_name` - (Deprecated since v1.230.1). Field 'transit_router_attachment_name' has been deprecated from provider version 1.230.1. New field 'transit_router_vpc_attachment_name' instead.
* `route_table_association_enabled` - (Optional, Bool, Deprecated since v1.192.0) Whether to enabled route table association. **NOTE:** "Field `route_table_association_enabled` has been deprecated from provider version 1.192.0. Please use the resource `alicloud_cen_transit_router_route_table_association` instead, [how to use alicloud_cen_transit_router_route_table_association](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_route_table_association)."
* `route_table_propagation_enabled` - (Optional, Bool, Deprecated since v1.192.0) Whether to enabled route table propagation. **NOTE:** "Field `route_table_propagation_enabled` has been deprecated from provider version 1.192.0. Please use the resource `alicloud_cen_transit_router_route_table_propagation` instead, [how to use alicloud_cen_transit_router_route_table_propagation](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_route_table_propagation)."

### `zone_mappings`

The zone_mappings supports the following:
* `vswitch_id` - (Required) The ID of the vSwitch that you want to add to the VPC connection.  You can specify at most 10 vSwitches in each call.
  - If the VPC connection belongs to the current Alibaba Cloud account, you can call the [DescribeVSwitches](https://www.alibabacloud.com/help/en/doc-detail/35748.html) operation to query the IDs of the vSwitches and zones of the VPC.
  - If the VPC connection belongs to another Alibaba Cloud account, you can call the [ListGrantVSwitchesToCen](https://www.alibabacloud.com/help/en/doc-detail/427599.html) operation to query the IDs of the vSwitches and zones of the VPC.
* `zone_id` - (Required) The ID of the zone that supports Enterprise Edition transit routers.  You can call the [DescribeZones](https://www.alibabacloud.com/help/en/doc-detail/36064.html) operation to query the most recent zone list.  You can specify at most 10 zones in each call.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `transit_router_attachment_id` - The ID of the Transit Router Attachment.
* `create_time` - The creation time of the resource
* `status` - Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Vpc Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Vpc Attachment.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Vpc Attachment.

## Import

CEN Transit Router Vpc Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vpc_attachment.example <id>
```