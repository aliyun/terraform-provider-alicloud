---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Multicast Domain resource.
---

# alicloud\_cen\_transit\_router\_multicast\_domain

Provides a Cloud Enterprise Network (CEN) Transit Router Multicast Domain resource.

For information about Cloud Enterprise Network (CEN) Transit Router Multicast Domain and how to use it, see [What is Transit Router Multicast Domain](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-createtransitroutermulticastdomain).

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
  transit_router_id                           = alicloud_cen_transit_router.default.transit_router_id
  transit_router_multicast_domain_name        = "tf-example-name"
  transit_router_multicast_domain_description = "tf-example-description"
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Transit Router Multicast Domain.
* `update` - (Defaults to 3 mins) Used when update the Transit Router Multicast Domain.
* `delete` - (Defaults to 3 mins) Used when delete the Transit Router Multicast Domain.

## Import

Cloud Enterprise Network (CEN) Transit Router Multicast Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_multicast_domain.example <id>
```
