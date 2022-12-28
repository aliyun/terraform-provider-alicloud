---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_association"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-association"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association resource.
---

# alicloud\_cen\_transit\_router\_multicast\_domain\_association

Provides a Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association resource.

For information about Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association and how to use it, see [What is Transit Router Multicast Domain Association](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-associatetransitroutermulticastdomain).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_instance" "default" {
  cen_instance_name = "tf-example"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id            = alicloud_cen_instance.default.id
  support_multicast = true
}

resource "alicloud_cen_transit_router_multicast_domain" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_transit_router.default.cen_id
  transit_router_id = alicloud_cen_transit_router_multicast_domain.default.transit_router_id
  vpc_id            = "your_vpc_id"
  zone_mappings {
    zone_id    = "your_zone_id"
    vswitch_id = "your_vswitch_id"
  }
}

resource "alicloud_cen_transit_router_multicast_domain_association" "default" {
  transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.default.id
  transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
  vswitch_id                         = "your_vswitch_id"
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain.
* `transit_router_attachment_id` - (Required, ForceNew) The ID of the VPC connection.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Multicast Domain Association. It formats as `<transit_router_multicast_domain_id>:<transit_router_attachment_id>:<vswitch_id>`.
* `status` - The status of the Transit Router Multicast Domain Association.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Transit Router Multicast Domain Association.
* `delete` - (Defaults to 3 mins) Used when delete the Transit Router Multicast Domain Association.

## Import

Cloud Enterprise Network (CEN) Transit Router Multicast Domain Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_multicast_domain_association.example <transit_router_multicast_domain_id>:<transit_router_attachment_id>:<vswitch_id>
```
