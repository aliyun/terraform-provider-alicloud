---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregators"
sidebar_current: "docs-alicloud-datasource-config-aggregators"
description: |-
  Provides a list of Config Aggregators to the user.
---

# alicloud\_config\_aggregators

This data source provides the Config Aggregators of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_config_aggregators" "example" {
  ids        = ["ca-3ce2626622af0005****"]
  name_regex = "the_resource_name"
}

output "first_config_aggregator_id" {
  value = data.alicloud_config_aggregators.example.aggregators.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of aggregator ids.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by aggregator name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid Values:  `0`: creating `1`: normal `2`: deleting.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Aggregator names.
* `aggregators` - A list of config aggregators. Each element contains the following attributes:
	* `account_id` - The Aliyun Uid.
	* `aggregator_accounts` - Account information in aggregator.
		* `account_id` - Aggregator account uid.
		* `account_name` - Aggregator account name.
		* `account_type` - Aggregator account source type.
	* `aggregator_id` - The id of aggregator.
	* `aggregator_name` - The name of aggregator.
	* `aggregator_type` - The type of aggregator.
	* `description` - The description of aggregator.
	* `id` - The id of the aggregator.
	* `status` - The status of the resource.
