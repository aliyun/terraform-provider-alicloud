---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_route_table_aggregation"
sidebar_current: "docs-alicloud-resource-cen-transit-route-table-aggregation"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Route Table Aggregation resource.
---

# alicloud\_cen\_transit\_route\_table\_aggregation

Provides a Cloud Enterprise Network (CEN) Transit Route Table Aggregation resource.

For information about Cloud Enterprise Network (CEN) Transit Route Table Aggregation and how to use it, see [What is Transit Route Table Aggregation](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-createtransitroutetableaggregation).

-> **NOTE:** Available in v1.202.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_instance" "default" {
  cen_instance_name = "tf-example"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_route_table_aggregation" "default" {
  transit_route_table_id                      = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
  transit_route_table_aggregation_cidr        = "10.0.0.0/8"
  transit_route_table_aggregation_scope       = "VPC"
  transit_route_table_aggregation_name        = "tf-example-name"
  transit_route_table_aggregation_description = "tf-example-description"
}
```

## Argument Reference

The following arguments are supported:

* `transit_route_table_id` - (Required, ForceNew) The ID of the route table of the Enterprise Edition transit router.
* `transit_route_table_aggregation_cidr` - (Required, ForceNew) The destination CIDR block of the aggregate route. CIDR blocks that start with `0` or `100.64`. Multicast CIDR blocks, including `224.0.0.1` to `239.255.255.254`.
* `transit_route_table_aggregation_scope` - (Required, ForceNew) The scope of networks that you want to advertise the aggregate route. Valid Value: `VPC`.
* `transit_route_table_aggregation_name` - (Optional, ForceNew) The name of the aggregate route.
* `transit_route_table_aggregation_description` - (Optional, ForceNew) The description of the aggregate route.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Route Table Aggregation. It formats as `<transit_route_table_id>:<transit_route_table_aggregation_cidr>`.
* `status` - The status of the Transit Route Table Aggregation.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Transit Route Table Aggregation.
* `delete` - (Defaults to 3 mins) Used when delete the Transit Route Table Aggregation.

## Import

Cloud Enterprise Network (CEN) Transit Route Table Aggregation can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_route_table_aggregation.example <transit_route_table_id>:<transit_route_table_aggregation_cidr>
```
