---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router"
sidebar_current: "docs-alicloud-resource-cen-transit_router"
description: |-
  Provides a Alicloud CEN transit router resource.
---

# alicloud_cen_transit_router

Provides a CEN transit router resource that associate the transitRouter with the CEN instance.[What is Cen Transit Router](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-createtransitrouter)

-> **NOTE:** Available since v1.126.0.

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
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `transit_router_name` - (Optional) The name of the transit router.
* `transit_router_description` - (Optional) The description of the transit router.
* `support_multicast` - (Optional, ForceNew, Available in v1.195.0+) Specifies whether to enable the multicast feature for the Enterprise Edition transit router. Valid values: `false`, `true`. Default Value: `false`. The multicast feature is supported only in specific regions. You can call [ListTransitRouterAvailableResource](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-listtransitrouteravailableresource) to query the regions that support multicast.
* `dry_run` - (Optional) The dry run.
* `tags` - (Optional, Available in v1.193.1+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<cen_id>:<transit_router_id>`.
* `status` - The associating status of the Transit Router.
* `type` - The Type of the Transit Router. Valid values: `Enterprise`, `Basic`.
* `transit_router_id` -  The transit router id of the transit router.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when creating the cen transit router (until it reaches the initial `Active` status).
* `update` - (Defaults to 3 mins) Used when update the cen transit router.
* `delete` - (Defaults to 3 mins) Used when delete the cen transit router.

## Import

CEN instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router.default cen-*****:tr-*******
```
