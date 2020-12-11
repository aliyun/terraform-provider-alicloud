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
For information about Alicloud Config Rule and how to use it, see [What is Alicloud Config Rule](https://www.alibabacloud.com/help/en/doc-detail/127388.htm).

-> **NOTE:** Available in v1.99.0+.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-northeast-1`.

-> **NOTE:** If you use custom rules, you need to create your own rule functions in advance. Please refer to the link for [Create a custom rule.](https://www.alibabacloud.com/help/en/doc-detail/127405.htm)

## Example Usage

```terraform
# Audit ECS instances under VPC using preset rules
resource "alicloud_config_rule" "example" {
  rule_name                       = "instances-in-vpc"
  source_identifier               = "ecs-instances-in-vpc"
  source_owner                    = "ALIYUN"
  scope_compliance_resource_types = ["ACS::ECS::Instance"]
  description                     = "ecs instances in vpc"
  input_parameters = {
    vpcIds = "vpc-uf6gksw4ctjd******"
  }
  risk_level                         = 1
  scope_compliance_resource_id       = "i-uf6j6rl141ps******"
  source_detail_message_type         = "ConfigurationItemChangeNotification"
  source_maximum_execution_frequency = "Twelve_Hours"
}

```
## Argument Reference

The following arguments are supported:

* `rule_name` - (Required, ForceNew) The name of the Config Rule. 
* `description` - (Optional) The description of the Config Rule.
* `multi_account` - (Optional, ForceNew) Whether the enterprise management account is a member account to create or modify rules. Valid values: `true`: Enterprise management accounts create or modify rules for all member accounts in the resource directory. `false`:The enterprise management account creates or modifies rules for this account. Default value is `false`.
* `member_id` - (Optional, ForceNew) The ID of the member account to which the rule to be created or modified belongs. The default is empty. When `multi_account` is set to true, this parameter is valid.
* `risk_level` - (Required) The risk level of the Config Rule. Valid values: `1`: Critical ,`2`: Warning , `3`: Info.
* `source_owner` - (Required, ForceNew) The source owner of the Config Rule. Values: `CUSTOM_FC`: Custom rules, `ALIYUN`: Trusteeship rules.
* `source_detail_message_type` - (Required) Trigger mechanism of rules. Valid values: `ConfigurationItemChangeNotification`,`OversizedConfigurationItemChangeNotification` and `ScheduledNotification`.
* `source_identifier` - (Required, ForceNew) The name of the custom rule or managed rules. Using managed rules, refer to [List of Managed rules.](https://www.alibabacloud.com/help/en/doc-detail/127404.htm)
* `input_parameters` - (Optional) Threshold value for managed rule triggering. 
* `source_maximum_execution_frequency` - (Optional) Rule execution cycle. Valid values: `One_Hour`, `Three_Hours`, `Six_Hours`, `Twelve_Hours` and `TwentyFour_Hours`.
* `scope_compliance_resource_types` - (Required) Resource types to be evaluated. [Alibaba Cloud services that support Cloud Config.](https://www.alibabacloud.com/help/en/doc-detail/127411.htm)
* `scope_compliance_resource_id` - (Optional) The ID of the resource to be evaluated. If not set, all resources are evaluated.

-> **NOTE:** When you use the personal version to configure auditing, please ignore `multi_account` and `member_id`.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of the Config Rule.  

## Import

Alicloud Config Rule can be imported using the id, e.g.

```
$ terraform import alicloud_config_rule.this cr-ed4bad756057********
```
