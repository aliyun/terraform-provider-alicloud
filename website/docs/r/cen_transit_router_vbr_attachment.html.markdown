---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vbr_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_vbr_attachment"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router VBR Attachment resource.
---

# alicloud_cen_transit_router_vbr_attachment

Provides a Cloud Enterprise Network (CEN) Transit Router VBR Attachment resource.

For information about Cloud Enterprise Network (CEN) Transit Router VBR Attachment and how to use it, see [What is Transit Router VBR Attachment](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitroutervbrattachment)

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_vbr_attachment&exampleId=5ae00d58-3d74-bbad-72d1-fc88f25aca0d0b654a78&activeTab=example&spm=docs.r.cen_transit_router_vbr_attachment.0.5ae00d583d&intl_lang=EN_US" target="_blank">
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

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = 2420
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  cen_id                                = alicloud_cen_instance.default.id
  vbr_id                                = alicloud_express_connect_virtual_border_router.default.id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `vbr_id` - (Required, ForceNew) The ID of the VBR.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `resource_type` - (Optional, ForceNew) The resource type of the transit router vbr attachment. Default value: `VBR`. Valid values: `VBR`.
* `vbr_owner_id` - (Optional, ForceNew) The owner id of the vbr.
* `auto_publish_route_enabled` - (Optional, Bool) Specifies whether to enable the Enterprise Edition transit router to automatically advertise routes to the VBR. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `transit_router_attachment_name` - (Optional) The name of the transit router vbr attachment.
* `transit_router_attachment_description` - (Optional) The description of the transit router vbr attachment.
* `tags` - (Optional, Available since v1.193.1) A mapping of tags to assign to the resource.
* `dry_run` - (Optional, Bool) Specifies whether to perform a dry run. Default value: `false`. Valid values: `true`, `false`.
* `route_table_association_enabled` - (Optional, Bool, Deprecated since v1.233.1) Whether to enabled route table association. **NOTE:** "Field `route_table_association_enabled` has been deprecated from provider version 1.233.1. Please use the resource `alicloud_cen_transit_router_route_table_association` instead, [how to use alicloud_cen_transit_router_route_table_association](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_route_table_association)."
* `route_table_propagation_enabled` - (Optional, Bool, Deprecated since v1.233.1) Whether to enabled route table propagation. **NOTE:** "Field `route_table_propagation_enabled` has been deprecated from provider version 1.233.1. Please use the resource `alicloud_cen_transit_router_route_table_propagation` instead, [how to use alicloud_cen_transit_router_route_table_propagation](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_route_table_propagation)."

->**NOTE:** Ensure that the vbr is not used in Express Connect.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router VBR Attachment. It formats as `<cen_id>:<transit_router_attachment_id>`.
* `transit_router_attachment_id` - The ID of the VBR connection.
* `status` - The status of the Transit Router VBR Attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Transit Router VBR Attachment.
* `update` - (Defaults to 10 mins) Used when update the Transit Router VBR Attachment.
* `delete` - (Defaults to 10 mins) Used when delete the Transit Router VBR Attachment.

## Import

Cloud Enterprise Network (CEN) Transit Router VBR Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vbr_attachment.example <cen_id>:<transit_router_attachment_id>
```
