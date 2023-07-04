---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Multicast Domain resource.
---

# alicloud_cen_transit_router_multicast_domain

Provides a Cloud Enterprise Network (CEN) Transit Router Multicast Domain resource.

For information about Cloud Enterprise Network (CEN) Transit Router Multicast Domain and how to use it, see [What is Transit Router Multicast Domain](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-createtransitroutermulticastdomain).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = "tf_example"
  cen_id              = alicloud_cen_instance.example.id
  support_multicast   = true
}

resource "alicloud_cen_transit_router_multicast_domain" "example" {
  transit_router_id                           = alicloud_cen_transit_router.example.transit_router_id
  transit_router_multicast_domain_name        = "tf_example"
  transit_router_multicast_domain_description = "tf_example"
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_multicast_domain_name` - (Optional) The name of the multicast domain. The name must be 0 to 128 characters in length, and can contain letters, digits, commas (,), periods (.), semicolons (;), forward slashes (/), at signs (@), underscores (_), and hyphens (-).
* `transit_router_multicast_domain_description` - (Optional) The description of the multicast domain. The description must be 0 to 256 characters in length, and can contain letters, digits, commas (,), periods (.), semicolons (;), forward slashes (/), at signs (@), underscores (_), and hyphens (-).
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Multicast Domain.
* `status` - The status of the Transit Router Multicast Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Transit Router Multicast Domain.
* `update` - (Defaults to 3 mins) Used when update the Transit Router Multicast Domain.
* `delete` - (Defaults to 3 mins) Used when delete the Transit Router Multicast Domain.

## Import

Cloud Enterprise Network (CEN) Transit Router Multicast Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_multicast_domain.example <id>
```
