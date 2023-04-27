---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_rule"
sidebar_current: "docs-alicloud-resource-config-rule"
description: |-
  Provides a Alicloud Config Rule resource.
---

# alicloud_config_rule

Provides a Config Rule resource.

For information about Config Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available in v1.204.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_config_rule" "default" {
  description                = "关联的资源类型下实体资源均已有指定标签，存在没有指定标签的资源则视为“不合规”。"
  source_owner               = "ALIYUN"
  source_identifier          = "required-tags"
  risk_level                 = 1
  tag_value_scope            = "test"
  tag_key_scope              = "test"
  exclude_resource_ids_scope = "test"
  region_ids_scope           = "cn-hangzhou"
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  resource_group_ids_scope   = "rg-acfmvoh45rhcfly"
  resource_types_scope = [
  "ACS::RDS::DBInstance"]
  rule_name = "tf-cicd-rule-by-required-tags"
  input_parameters {
  }
}
```

## Argument Reference

The following arguments are supported:
* `config_rule_trigger_types` - (Required) The trigger type of the rule. Valid values:  `ConfigurationItemChangeNotification`: The rule is triggered upon configuration changes. `ScheduledNotification`: The rule is triggered as scheduled.
* `description` - (Optional) The description of the rule.
* `exclude_resource_ids_scope` - (Optional) The rule monitors excluded resource IDs, multiple of which are separated by commas, only applies to rules created based on managed rules, , custom rule this field is empty.
* `input_parameters` - (Optional) The settings of the input parameters for the rule.See the following `Block InputParameters`.
* `maximum_execution_frequency` - (Computed, Optional) The frequency of the compliance evaluations, it is required if the ConfigRuleTriggerTypes value is ScheduledNotification. Valid values:  `One_Hour`, `Three_Hours`, `Six_Hours`, `Twelve_Hours`, `TwentyFour_Hours`.
* `region_ids_scope` - (Optional) The rule monitors region IDs, separated by commas, only applies to rules created based on managed rules.
* `resource_group_ids_scope` - (Optional) The rule monitors resource group IDs, separated by commas, only applies to rules created based on managed rules.
* `resource_types_scope` - (Required) The types of the resources to be evaluated against the rule.
* `risk_level` - (Required) The risk level of the resources that are not compliant with the rule. Valid values:  `1`: critical `2`: warning `3`: info
* `rule_name` - (Required, ForceNew) The name of the rule.
* `source_identifier` - (Required, ForceNew) The identifier of the rule.  For a managed rule, the value is the name of the managed rule. For a custom rule, the value is the ARN of the custom rule.
* `source_owner` - (Required, ForceNew) Specifies whether you or Alibaba Cloud owns and manages the rule. Valid values:  `CUSTOM_FC`: The rule is a custom rule and you own the rule. `ALIYUN`: The rule is a managed rule and Alibaba Cloud owns the rule
* `status` - (Computed, Optional) The status of the rule. Valid values: ACTIVE: The rule is monitoring the configurations of target resources. DELETING_RESULTS: The compliance evaluation result returned by the rule is being deleted. EVALUATING: The rule is triggered and is evaluating whether the configurations of target resources are compliant. INACTIVE: The rule is disabled from monitoring the configurations of target resources.
* `tag_key_scope` - (Optional) The rule monitors the tag key, only applies to rules created based on managed rules.
* `tag_value_scope` - (Optional) The rule monitors the tag value, only applies to rules created based on managed rules.

The following arguments will be discarded. Please use new fields as soon as possible:
* `source_detail_message_type` - (Deprecated) Field 'source_detail_message_type' has been deprecated from provider version 1.124.1. New field 'config_rule_trigger_types' instead.
* `source_maximum_execution_frequency` - (Deprecated) Field 'source_maximum_execution_frequency' has been deprecated from provider version 1.124.1. New field 'maximum_execution_frequency' instead.
* `scope_compliance_resource_types` - (Deprecated) Field 'scope_compliance_resource_types' has been deprecated from provider version 1.124.1. New field 'resource_types_scope' instead.

#### Block InputParameters

The InputParameters supports the following:



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `account_id` - The ID of Alicloud account.
* `compliance` - compliance information.
  * `compliance_type` - The type of compliance. Valid values: `COMPLIANT`, `NON_COMPLIANT`, `NOT_APPLICABLE`, `INSUFFICIENT_DATA`.
  * `count` - The count of compliance.
* `compliance_pack_id` - Compliance Package ID.
* `config_rule_arn` - config rule arn.
* `config_rule_id` - The ID of the rule.
* `create_time` - The timestamp when the rule was created.
* `event_source` - The event source of the rule.
* `maximum_execution_frequency` - The frequency of the compliance evaluations, it is required if the ConfigRuleTriggerTypes value is ScheduledNotification. Valid values:  `One_Hour`, `Three_Hours`, `Six_Hours`, `Twelve_Hours`, `TwentyFour_Hours`.
* `modified_timestamp` - The timestamp when the rule was last modified.
* `status` - The status of the rule. Valid values: ACTIVE: The rule is monitoring the configurations of target resources. DELETING_RESULTS: The compliance evaluation result returned by the rule is being deleted. EVALUATING: The rule is triggered and is evaluating whether the configurations of target resources are compliant. INACTIVE: The rule is disabled from monitoring the configurations of target resources.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Rule.
* `update` - (Defaults to 5 mins) Used when update the Rule.

## Import

Config Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_rule.example <id>
```