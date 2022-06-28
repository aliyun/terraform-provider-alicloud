---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_rules"
sidebar_current: "docs-alicloud-datasource-config-rules"
description: |-
    Provides a list of Config Rules to the user.
---

# alicloud\_config\_rules

This data source provides the Config Rules of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.99.0+.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-southeast-1`.

## Example Usage

```terraform
data "alicloud_config_rules" "example" {
  ids        = ["cr-ed4bad756057********"]
  name_regex = "tftest"
}

output "first_config_rule_id" {
  value = data.alicloud_config_rules.example.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Config Rule IDs.
* `status` - (Optional, ForceNew, Available in 1.124.1+) The status of the config rule, valid values: `ACTIVE`, `DELETING`, `EVALUATING` and `INACTIVE`. 
* `rule_name` - (Optional, ForceNew, Available in 1.124.1+) The name of config rule.
* `multi_account` - (Optional, ForceNew,Removed) Field `multi_account` has been removed from provider version 1.146.0. Please Use the Resource `alicloud_config_aggregate_config_rule`.
* `member_id` - (Optional, ForceNew,Removed) Field `multi_account` has been removed from provider version 1.146.0. Please Use the Resource `alicloud_config_aggregate_config_rule`.
* `risk_level` - (Optional, ForceNew) The risk level of Config Rule. Valid values: `1`: Critical ,`2`: Warning , `3`: Info.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by rule name.
* `message_type` - (Optional, ForceNew,  Available in v1.104.0+, Remove) Field `message_type` has been removed from provider version 1.124.1.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `config_rule_state` - (Optional, ForceNew, Deprecated) Field `config_rule_state` has been deprecated from provider version 1.124.1. New field `status` instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Config Rule IDs.
* `names` - A list of Config Rule names.
* `rules` - A list of Config Rules. Each element contains the following attributes:
    * `id` - The ID of the Config Rule.
    * `account_id`- The ID of the Alicloud account.
    * `config_rule_arn`- The ARN of the Config Rule.
    * `config_rule_id`- The ID of the Config Rule.
    * `config_rule_state`- The state of the Config Rule.
    * `description`- The description of the Config Rule.
    * `input_parameters`- The input parameters of the Config Rule.
    * `modified_timestamp`- the timestamp of the Config Rule modified.
    * `risk_level`- The risk level of the Config Rule.
    * `rule_name`- The name of the Config Rule.
    * `event_source` - Event source of the Config Rule.
    * `scope_compliance_resource_id` - (Removed)  Field `scope_compliance_resource_id` has been removed from provider version 1.124.1. Please use 'exclude_resource_ids_scope' instead. 
    * `scope_compliance_resource_types` - The types of the resources to be evaluated against the rule.
    * `source_detail_message_type` - Rule trigger mechanism.
    * `source_maximum_execution_frequency` - Rule execution cycle. 
    * `source_identifier`- The identifier of the managed rule or the arn of the custom function.
    * `source_owner`- The source owner of the Config Rule.
    * `compliance` - The information about the compliance evaluations based on the rule.
        * `compliance_type` - The compliance evaluation result of the target resources.
        * `count` - The number of resources with the specified compliance evaluation result.
    * `config_rule_trigger_types` - (Available in 1.124.1+) A list of trigger types of config rule.
    * `exclude_resource_ids_scope` - (Available in 1.124.1+) The scope of exclude of resource ids.
    * `maximum_execution_frequency` - (Available in 1.124.1+) The frequency of maximum execution.
    * `region_ids_scope` - (Available in 1.124.1+) The scope of region ids.
    * `resource_group_ids_scope` - (Available in 1.124.1+) The scope of resource group ids.
    * `resource_types_scope` - (Available in 1.124.1+) The scope of resource types.
    * `status` - (Available in 1.124.1+) The status of config rule.
    * `tag_key_scope` - (Available in 1.124.1+) The scope of tag key.
    * `tag_value_scope` - (Available in 1.124.1+) The scope of tag value.
