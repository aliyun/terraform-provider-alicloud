---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_rule"
sidebar_current: "docs-alicloud-resource-config-rule"
description: |-
  Provides a Alicloud Config Rule resource.
---

# alicloud\_config\_rule

Provides a a Alicloud Config Rule resource. Cloud Config checks the validity of resources based on rules. You can create rules to evaluate resources as needed.
For information about Alicloud Config Rule and how to use it, see [What is Alicloud Config Rule](https://www.alibabacloud.com/help/doc-detail/154216.html).

-> **NOTE:** Available in v1.99.0+.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-southeast-1`.

-> **NOTE:** If you use custom rules, you need to create your own rule functions in advance. Please refer to the link for [Create a custom rule.](https://www.alibabacloud.com/help/en/doc-detail/127405.htm)

## Example Usage

```terraform
# Audit ECS instances under VPC using preset rules
resource "alicloud_config_rule" "example" {
  rule_name            = "instances-in-vpc"
  source_identifier    = "ecs-instances-in-vpc"
  source_owner         = "ALIYUN"
  resource_types_scope = ["ACS::ECS::Instance"]
  description          = "ecs instances in vpc"
  input_parameters = {
    vpcIds = "vpc-uf6gksw4ctjd******"
  }
  risk_level                = 1
  config_rule_trigger_types = "ConfigurationItemChangeNotification"
}

```
## Argument Reference

The following arguments are supported:

* `maximum_execution_frequency` - (Optional, Available in v1.124.1+) The frequency of the compliance evaluations, it is required if the ConfigRuleTriggerTypes value is ScheduledNotification. Valid values: `One_Hour`, `Three_Hours`, `Six_Hours`, `Twelve_Hours`, `TwentyFour_Hours`.
* `resource_types_scope` - (Optional, Available in v1.124.1+) Resource types to be evaluated. [Alibaba Cloud services that support Cloud Config.](https://www.alibabacloud.com/help/en/doc-detail/127411.htm)
* `config_rule_trigger_types` - (Optional, Available in v1.124.1+) The trigger type of the rule. Valid values: `ConfigurationItemChangeNotification`: The rule is triggered upon configuration changes. `ScheduledNotification`: The rule is triggered as scheduled.
* `exclude_resource_ids_scope` - (Optional, Available in v1.124.1+) The rule monitors excluded resource IDs, multiple of which are separated by commas, only applies to rules created based on managed rules, custom rule this field is empty.
* `region_ids_scope` - (Optional, Available in v1.124.1+) The rule monitors region IDs, separated by commas, only applies to rules created based on managed rules.
* `resource_group_ids_scope` - (Optional, Available in v1.124.1+) The rule monitors resource group IDs, separated by commas, only applies to rules created based on managed rules.
* `tag_key_scope` - (Optional, Available in v1.124.1+) The rule monitors the tag key, only applies to rules created based on managed rules.
* `tag_value_scope` - (Optional, Available in v1.124.1+) The rule monitors the tag value, use with the `tag_key_scope` options. only applies to rules created based on managed rules.
* `rule_name` - (Required, ForceNew) The name of the Config Rule. 
* `description` - (Optional) The description of the Config Rule.
* `risk_level` - (Required) The risk level of the Config Rule. Valid values: `1`: Critical ,`2`: Warning , `3`: Info.
* `source_owner` - (Required, ForceNew) Specifies whether you or Alibaba Cloud owns and manages the rule. Valid values: `CUSTOM_FC`: The rule is a custom rule and you own the rule. `ALIYUN`: The rule is a managed rule and Alibaba Cloud owns the rule.
* `source_identifier` - (Required, ForceNew) The identifier of the rule. For a managed rule, the value is the identifier of the managed rule. For a custom rule, the value is the ARN of the custom rule. Using managed rules, refer to [List of Managed rules.](https://www.alibabacloud.com/help/en/doc-detail/127404.htm)
* `input_parameters` - (Optional) Threshold value for managed rule triggering. 
* `source_detail_message_type` - (Optional, Deprecated) Field `source_detail_message_type` has been deprecated from provider version 1.124.1. New field `config_rule_trigger_types` instead.
* `source_maximum_execution_frequency` - (Optional, Deprecated) Field `source_maximum_execution_frequency` has been deprecated from provider version 1.124.1. New field `maximum_execution_frequency` instead.
* `scope_compliance_resource_types` - (Optional, Deprecated) Field `scope_compliance_resource_types` has been deprecated from provider version 1.124.1. New field `resource_types_scope` instead.
* `multi_account` - (Optional, Removed) Field `multi_account` has been removed from provider version 1.124.1. 
* `member_id` - (Optional, Removed) Field `member_id` has been removed from provider version 1.124.1. 
* `scope_compliance_resource_id` - (Optional, Removed) Field `scope_compliance_resource_id` has been removed from provider version 1.124.1. 
* `status` - (Optional, Available in v1.145.0+) The rule status. The valid values: `ACTIVE`, `INACTIVE`.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of the Config Rule.  

### Timeouts

-> **NOTE:** Available in v1.124.1.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Config Rule.
* `update` - (Defaults to 10 mins) Used when update the Config Rule.

## Import

Alicloud Config Rule can be imported using the id, e.g.

```
$ terraform import alicloud_config_rule.this cr-ed4bad756057********
```
