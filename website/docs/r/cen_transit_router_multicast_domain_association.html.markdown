---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_association"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-association"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association resource.
---

# alicloud_cen_transit_router_multicast_domain_association

Provides a Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association resource.

For information about Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association and how to use it, see [What is Transit Router Multicast Domain Association](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-associatetransitroutermulticastdomain).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_multicast_domain_association&exampleId=013ff97b-ac69-fc5e-14e2-7be4981f42d60cfe1096&activeTab=example&spm=docs.r.cen_transit_router_multicast_domain_association.0.013ff97bac&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

locals {
  zone = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "192.168.1.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.zone
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.example.id
  support_multicast   = true
}

resource "alicloud_cen_transit_router_multicast_domain" "example" {
  transit_router_id                    = alicloud_cen_transit_router.example.transit_router_id
  transit_router_multicast_domain_name = var.name
}

resource "alicloud_cen_transit_router_vpc_attachment" "example" {
  cen_id            = alicloud_cen_transit_router.example.cen_id
  transit_router_id = alicloud_cen_transit_router_multicast_domain.example.transit_router_id
  vpc_id            = alicloud_vpc.example.id
  zone_mappings {
    zone_id    = local.zone
    vswitch_id = alicloud_vswitch.example.id
  }
}

resource "alicloud_cen_transit_router_multicast_domain_association" "example" {
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.example.id
  transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.example.transit_router_attachment_id
  vswitch_id                         = alicloud_vswitch.example.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_transit_router_multicast_domain_association&spm=docs.r.cen_transit_router_multicast_domain_association.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain.
* `transit_router_attachment_id` - (Required, ForceNew) The ID of the VPC connection.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Multicast Domain Association. It formats as `<transit_router_multicast_domain_id>:<transit_router_attachment_id>:<vswitch_id>`.
* `status` - The status of the Transit Router Multicast Domain Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Transit Router Multicast Domain Association.
* `delete` - (Defaults to 3 mins) Used when delete the Transit Router Multicast Domain Association.

## Import

Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_multicast_domain_association.example <transit_router_multicast_domain_id>:<transit_router_attachment_id>:<vswitch_id>
```
