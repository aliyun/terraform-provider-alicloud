---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_route_table_aggregation"
sidebar_current: "docs-alicloud-resource-cen-transit-route-table-aggregation"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Route Table Aggregation resource.
---

# alicloud_cen_transit_route_table_aggregation

Provides a Cloud Enterprise Network (CEN) Transit Route Table Aggregation resource.

For information about Cloud Enterprise Network (CEN) Transit Route Table Aggregation and how to use it, see [What is Transit Route Table Aggregation](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitroutetableaggregation).

-> **NOTE:** Available since v1.202.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_route_table_aggregation&exampleId=4f81bb83-d48e-254f-cb61-7ddd046ecf93c9213bab&activeTab=example&spm=docs.r.cen_transit_route_table_aggregation.0.4f81bb83d4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

resource "alicloud_cen_transit_route_table_aggregation" "example" {
  transit_route_table_id                      = alicloud_cen_transit_router_route_table.example.transit_router_route_table_id
  transit_route_table_aggregation_cidr        = "10.0.0.0/8"
  transit_route_table_aggregation_scope       = "VPC"
  transit_route_table_aggregation_name        = "tf_example"
  transit_route_table_aggregation_description = "tf_example"
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Transit Route Table Aggregation.
* `delete` - (Defaults to 3 mins) Used when delete the Transit Route Table Aggregation.

## Import

Cloud Enterprise Network (CEN) Transit Route Table Aggregation can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_route_table_aggregation.example <transit_route_table_id>:<transit_route_table_aggregation_cidr>
```
