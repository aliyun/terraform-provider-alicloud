---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table"
sidebar_current: "docs-alicloud-resource-cen-transit_router_route_table"
description: |-
  Provides a Alicloud CEN transit router route table resource.
---

# alicloud_cen_transit_router_route_table

Provides a CEN transit router route table resource.[What is Cen Transit Router Route Table](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitrouterroutetable)

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

resource "alicloud_cen_transit_router_route_table" "example" {
  transit_router_id = alicloud_cen_transit_router.example.transit_router_id
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_id` - (Required, ForceNew) The ID of the transit router.
* `transit_router_route_table_name` - (Optional) The name of the transit router route table.
* `transit_router_route_table_description` - (Optional) The description of the transit router route table.
* `tags` - (Optional, Available in v1.201.0+) A mapping of tags to assign to the resource.
* `dry_run` - (Optional) The dry run.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_route_table_id>`.
* `status` - The associating status of the Transit Router.
* `transit_router_route_table_id` - The id of the transit router route table.
* `transit_router_route_table_type` - The type of the transit router route table. Valid values: `Custom`, `System`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when creating the cen transit router route table (until it reaches the initial `Active` status).
* `update` - (Defaults to 3 mins) Used when update the cen transit router route table.
* `delete` - (Defaults to 3 mins) Used when delete the cen transit router route table.

## Import

CEN transit router route table  can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_route_table.default tr-*********:vtb-********
```
