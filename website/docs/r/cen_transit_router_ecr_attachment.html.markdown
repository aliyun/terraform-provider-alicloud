---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_ecr_attachment"
description: |-
  Provides a Alicloud CEN Transit Router Ecr Attachment resource.
---

# alicloud_cen_transit_router_ecr_attachment

Provides a CEN Transit Router Ecr Attachment resource.



For information about CEN Transit Router Ecr Attachment and how to use it, see [What is Transit Router Ecr Attachment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "asn" {
  default = "4200000667"
}

resource "alicloud_express_connect_router_express_connect_router" "defaultO8Hcfx" {
  alibaba_side_asn = var.asn
  ecr_name         = var.name
}

resource "alicloud_cen_instance" "defaultQKBiay" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultQa94Y1" {
  cen_id              = alicloud_cen_instance.defaultQKBiay.id
  transit_router_name = var.name
}

data "alicloud_account" "current" {
}

resource "alicloud_express_connect_router_tr_association" "defaultedPu6c" {
  association_region_id   = "cn-hangzhou"
  ecr_id                  = alicloud_express_connect_router_express_connect_router.defaultO8Hcfx.id
  cen_id                  = alicloud_cen_instance.defaultQKBiay.id
  transit_router_id       = alicloud_cen_transit_router.defaultQa94Y1.transit_router_id
  transit_router_owner_id = data.alicloud_account.current.id
}

resource "alicloud_cen_transit_router_ecr_attachment" "default" {
  ecr_id                                = alicloud_express_connect_router_express_connect_router.defaultO8Hcfx.id
  cen_id                                = alicloud_express_connect_router_tr_association.defaultedPu6c.cen_id
  transit_router_ecr_attachment_name    = var.name
  transit_router_attachment_description = var.name
  transit_router_id                     = alicloud_cen_transit_router.defaultQa94Y1.transit_router_id
  ecr_owner_id                          = data.alicloud_account.current.id
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (Optional, ForceNew) CenId
* `ecr_id` - (Required, ForceNew) EcrId
* `ecr_owner_id` - (Optional, ForceNew, Int) EcrOwnerId
* `tags` - (Optional, Map) The tag of the resource
* `transit_router_attachment_description` - (Optional) TransitRouterAttachmentDescription
* `transit_router_ecr_attachment_name` - (Optional) TransitRouterAttachmentName
* `transit_router_id` - (Optional, ForceNew) TransitRouterId

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Ecr Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Ecr Attachment.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Ecr Attachment.

## Import

CEN Transit Router Ecr Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_ecr_attachment.example <id>
```