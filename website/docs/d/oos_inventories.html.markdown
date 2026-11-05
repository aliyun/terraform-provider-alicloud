---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_inventories"
description: |-
  Provides a list of Oos Inventories to the user.
---

# alicloud_oos_inventories

This data source provides the Oos Inventories of the current Alibaba Cloud user.

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_oos_inventories" "default" {
}

output "inventory_entities" {
  value = data.alicloud_oos_inventories.default.inventories
}
```

## Argument Reference

The following arguments are supported:

* `next_token` - (Optional, String) The token that starts the query. Leave it unset to start from the first page.
* `max_results` - (Optional, Int) The number of entries per page. Valid values: `1` to `100`. Defaults to `50`.
* `aggregator` - (Optional, List) The aggregators. You can use one or more aggregators to query instance aggregate information. Such as: `ACS:InstanceInformation.PlatformName`. Valid values:

  | TypeName | Aggregatable Attribute |
  | --- | --- |
  | ACS:InstanceInformation | PlatformName, PlatformType, PlatformVersion |
  | ACS:Application | Name, Version |
  | ACS:WindowsRole | Name |
  | ACS:Service | Name |
  | ACS:WindowsRegistry | KeyPath |
  | ACS:WindowsUpdate | HotFixId |

* `filter` - (Optional, Set) The component filter rules. See [`filter`](#filter) below.

### `filter`

The filter supports the following:

* `name` - (Optional, String) The property name of the component, for example `ACS:InstanceInformation.InstanceId`.
* `operator` - (Optional, String) The comparison operator for the component property query value. Valid values: `Equal`, `NotEqual`, `BeginWith`, `LessThan` and `GreaterThan`.
* `value` - (Optional, List) The list of component property query values.

## Attributes Reference

The following attributes are exported:

* `next_token_result` - The token that is used to retrieve the next page of results.
* `inventories` - The inventory entities. Each element is the JSON representation of one entity. The keys of an entity are the component names, which vary with the queried component. In the aggregation mode, the `Content` of an entity contains the `Count`.
