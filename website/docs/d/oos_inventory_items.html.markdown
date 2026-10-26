---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_inventory_items"
sidebar_current: "docs-alicloud-datasource-oos-inventory-items"
description: |-
  Provides a list of configuration inventory entries collected by OOS.
---

# alicloud_oos_inventory_items

This data source provides a list of configuration inventory entries collected by CloudOps Orchestration Service (OOS) in an Alibaba Cloud account. It calls the [SearchInventory](https://www.alibabacloud.com/help/en/operation-orchestration-service/developer-reference/api-oos-2019-06-01-searchinventory) API.

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** The request is sent to the OOS endpoint of the region configured on the provider. `region_id` therefore has to match the provider region; a different value is rejected instead of returning the inventory of the endpoint region. To query another region, configure a provider alias for that region and use that provider in this data source.

## Example Usage

```terraform
data "alicloud_oos_inventory_items" "example" {
  max_results = 10

  filter = [
    {
      name     = "ACS:InstanceInformation.InstanceId"
      operator = "Equal"
      value    = ["i-bp1cpoxxx"]
    }
  ]
}

output "inventory_items" {
  value = jsondecode(data.alicloud_oos_inventory_items.example.entities)
}
```

The following example aggregates the collected applications by name and version:

```terraform
data "alicloud_oos_inventory_items" "aggregated" {
  filter = [
    {
      name     = "ACS:InstanceInformation.InstanceId"
      operator = "Equal"
      value    = ["i-bp1cpoxxx"]
    }
  ]

  aggregator = ["ACS:Application.Name", "ACS:Application.Version"]
}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Optional, String) The ID of the region where the configuration inventory is queried. Defaults to the region of the provider. It must be the region configured on the provider: the request is sent to the OOS endpoint of the provider region, so a different value is rejected instead of returning the inventory of the endpoint region. To query another region, use a provider alias configured for that region.
* `filter` - (Optional, List) Component filters. Up to 5 filter groups are supported. Every entry must set `name`, `operator` and `value`. Each entry supports the following attributes:
  * `name` - (Required, String) The component attribute name, for example `ACS:InstanceInformation.InstanceId`. The attribute names of each component can be discovered with the [GetInventorySchema](https://www.alibabacloud.com/help/en/operation-orchestration-service/developer-reference/api-oos-2019-06-01-getinventoryschema) API.
  * `operator` - (Required, String) The comparison operator of the filter value. Valid values: `Equal`, `NotEqual`, `BeginWith`, `LessThan` and `GreaterThan`.
  * `value` - (Required, List) The component attribute values to match. Up to 20 values are supported.
* `aggregator` - (Optional, List) The aggregators used to query inventory aggregation information. Up to 5 aggregators are supported. Each value has the form `ACS:<component>.<property>`. The aggregatable properties depend on the component and on the inventory data the account has reported; run the [GetInventorySchema](https://www.alibabacloud.com/help/en/operation-orchestration-service/developer-reference/api-oos-2019-06-01-getinventoryschema) API with `Aggregator` set to `true` to list the properties a component supports for aggregation. Commonly used values include `ACS:InstanceInformation.PlatformName`, `ACS:InstanceInformation.PlatformType`, `ACS:InstanceInformation.PlatformVersion`, `ACS:Application.Name` and `ACS:Application.Version`; other components such as `ACS:WindowsRole`, `ACS:Service`, `ACS:WindowsRegistry` or `ACS:WindowsUpdate` may offer further aggregatable properties depending on the inventory types the account collects. The provider does not restrict the values; an unsupported component/property combination is rejected by the API with an error such as `InvalidParameter`.
* `max_results` - (Optional, Int) The number of entries to return on each page. Valid values: 1 to 100. Defaults to `50`. All pages are collected, so this only controls the request page size.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (String) The ID of the data source.
* `entities` - (String) The inventory entries, encoded as a JSON string. Each array element is one component record and is keyed by the component name, for example `ACS:InstanceInformation` or `ACS:Application`. The component record carries the component content, which is a dynamic structure whose attribute names depend on the component; use `jsondecode` to read it in Terraform. The value is the JSON encoding of the `Entities` array returned by `SearchInventory`, with entries collected over multiple pages concatenated in page order. It mirrors the response exactly: a query without any entry yields `[]`, and an aggregation query can return a component envelope with an empty `Content` list.
