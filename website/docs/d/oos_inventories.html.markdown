---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_inventories"
sidebar_current: "docs-alicloud-datasource-oos-inventories"
description: |-
  Provides the configuration inventory of the managed instances in Operation Orchestration Service.
---

# alicloud\_oos\_inventories

This data source provides the configuration inventory of the managed instances, or the aggregated inventory information, in Operation Orchestration Service (OOS).

-> **NOTE:** Available since v2.2.0.

## Example Usage

```terraform
data "alicloud_oos_inventories" "default" {
  max_results = 10
}

output "inventory_entities" {
  value = data.alicloud_oos_inventories.default.entities
}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Optional, String) The ID of the region. Defaults to the region of the provider.
* `next_token` - (Optional, String) The token that starts the query. Leave it unset to start from the first page.
* `max_results` - (Optional, Int) The number of entries per page. Valid values: `1` to `100`. Defaults to `50`.
* `filter` - (Optional, List) The filter rules of the components. A maximum of 5 rules are supported. Each element contains the following attributes:
  * `name` - (Optional, String) The name of the component property, for example `ACS:InstanceInformation.InstanceId`.
  * `operator` - (Required, String) The comparison operator of the property value. Valid values: `Equal`, `NotEqual`, `BeginWith`, `LessThan` and `GreaterThan`.
  * `value` - (Optional, List) The property values. A maximum of 20 values are supported.
* `aggregator` - (Optional, List) The aggregators. Valid values: `ACS:Application.Name` and `ACS:Application.Version`. A maximum of 5 aggregators are supported.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source, which is derived from the query condition.
* `next_token_result` - (String) The token that is used to retrieve the next page of results.
* `entities` - (List) The inventory entities. Each element is the JSON representation of one entity. The keys of an entity are the component names, which vary with the queried component. In the aggregation mode, the `Content` of an entity contains the `Count`.
