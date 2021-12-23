---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_config_rules"
sidebar_current: "docs-alicloud-datasource-config-aggregate-config-rules"
description: |-
  Provides a list of Config Aggregate Config Rules to the user.
---

# alicloud\_config\_aggregate\_config\_rules

This data source provides the Config Aggregate Config Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_config_aggregate_config_rules" "example" {
  aggregator_id = "ca-3a9b626622af001d****"
  ids           = ["cr-5154626622af0034****"]
  name_regex    = "the_resource_name"
}

output "first_config_aggregate_config_rule_id" {
  value = data.alicloud_config_aggregate_config_rules.example.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, ForceNew) The ID of aggregator.
* `aggregate_config_rule_name` - (Optional, ForceNew) The config rule name.
* `status` - (Optional, ForceNew) The state of the config rule, valid values: `ACTIVE`, `DELETING`, `EVALUATING` and `INACTIVE`. 
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Aggregate Config Rule IDs.
* `risk_level` - Optional, ForceNew) The Risk Level. Valid values `1`: critical, `2`: warning, `3`: info.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Aggregate Config Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Aggregate Config Rule names.
* `rules` - A list of Config Aggregate Config Rules. Each element contains the following attributes:
	* `account_id` - The Aliyun User ID.
	* `aggregate_config_rule_name` - The name of the rule.
	* `aggregator_id` - The ID of Aggregator.
	* `compliance` -The Compliance information.
		* `count` - The Count.
		* `compliance_type` - The Compliance Type.
	* `compliance_pack_id` - The ID of Compliance Package.
	* `config_rule_arn` - The config rule arn.
	* `config_rule_id` - The ID of the rule.
	* `status` - The status of the rule. 
	* `event_source` - Event source of the Config Rule. 
	* `description` - The description of the rule.
	* `id` - The ID of the Aggregate Config Rule.
	* `config_rule_trigger_types` - The trigger types of config rules.
	* `exclude_resource_ids_scope` - The id of the resources to be evaluated against the rule.
    * `source_identifier`- The identifier of the managed rule or the arn of the custom function.
    * `source_owner`- The source owner of the Config Rule.
	* `maximum_execution_frequency` - The frequency of the compliance evaluations.
	* `region_ids_scope` - The scope of resource region ids.
	* `resource_group_ids_scope` - The scope of resource group ids.
	* `tag_key_scope` - The scope of tay key.
	* `tag_value_scope` - The scope of tay value.
	* `input_parameters` - The settings of the input parameters for the rule.
	* `modified_timestamp` - The timestamp when the rule was last modified.
	* `risk_level` - The risk level of the resources that are not compliant with the rule. Valid values: `1`: critical, `2`: warning, `3`: info.
