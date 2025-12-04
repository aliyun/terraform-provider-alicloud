---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_config_rule"
sidebar_current: "docs-alicloud-resource-config-aggregate-config-rule"
description: |-
  Provides a Alicloud Cloud Config Aggregate Config Rule resource.
---

# alicloud_config_aggregate_config_rule

Provides a Cloud Config Aggregate Config Rule resource.

For information about Cloud Config Aggregate Config Rule and how to use it, see [What is Aggregate Config Rule](https://www.alibabacloud.com/help/en/cloud-config/latest/api-config-2020-09-07-createaggregateconfigrule).

-> **NOTE:** Available since v1.124.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_aggregate_config_rule&exampleId=0b57245e-ee44-f056-7a35-37474e3366b6c5230193&activeTab=example&spm=docs.r.config_aggregate_config_rule.0.0b57245eee&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}

locals {
  last = length(data.alicloud_resource_manager_accounts.default.accounts) - 1
}

resource "alicloud_config_aggregator" "default" {
  aggregator_accounts {
    account_id   = data.alicloud_resource_manager_accounts.default.accounts[local.last].account_id
    account_name = data.alicloud_resource_manager_accounts.default.accounts[local.last].display_name
    account_type = "ResourceDirectory"
  }
  aggregator_name = var.name
  description     = var.name
  aggregator_type = "CUSTOM"
}
resource "alicloud_config_aggregate_config_rule" "default" {
  aggregate_config_rule_name = "contains-tag"
  aggregator_id              = alicloud_config_aggregator.default.id
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  source_owner               = "ALIYUN"
  source_identifier          = "contains-tag"
  risk_level                 = 1
  resource_types_scope       = ["ACS::ECS::Instance"]
  input_parameters = {
    key   = "example"
    value = "example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `aggregate_config_rule_name` - (Required, ForceNew) The name of the rule.
* `aggregator_id` - (Required, ForceNew) The Aggregator Id.
* `config_rule_trigger_types` - (Required) The trigger type of the rule. Valid values:
  - `ConfigurationItemChangeNotification`: The rule is triggered by configuration changes.
  - `ScheduledNotification`: The rule is periodically triggered.
-> **NOTE:** Separate multiple trigger types with commas (,).
* `description` - (Optional) The description of the rule.
* `exclude_resource_ids_scope` - (Optional) The rule monitors excluded resource IDs, multiple of which are separated by commas, only applies to rules created based on managed rules, , custom rule this field is empty.
* `input_parameters` - (Optional) The settings map of the input parameters for the rule.
* `source_identifier`- (Required, ForceNew) The identifier of the rule. For a managed rule, the value is the identifier of the managed rule. For a custom rule, the value is the ARN of the custom rule. Using managed rules, refer to [List of Managed rules.](https://www.alibabacloud.com/help/en/doc-detail/127404.htm)
* `source_owner`- (Required, ForceNew) Specifies whether you or Alibaba Cloud owns and manages the rule. Valid values: `CUSTOM_FC`: The rule is a custom rule and you own the rule. `ALIYUN`: The rule is a managed rule and Alibaba Cloud owns the rule.
* `maximum_execution_frequency` - (Optional) The frequency of the compliance evaluations. Valid values:  `One_Hour`, `Three_Hours`, `Six_Hours`, `Twelve_Hours`, `TwentyFour_Hours`. System default value is `TwentyFour_Hours` and valid when the `config_rule_trigger_types` is `ScheduledNotification`.
* `region_ids_scope` - (Optional) The rule monitors region IDs, separated by commas, only applies to rules created based on managed rules.
* `resource_group_ids_scope` - (Optional) The rule monitors resource group IDs, separated by commas, only applies to rules created based on managed rules.
* `resource_types_scope` - (Required) Resource types to be evaluated. [Alibaba Cloud services that support Cloud Config.](https://www.alibabacloud.com/help/en/doc-detail/127411.htm)
* `risk_level` - (Required) The risk level of the resources that are not compliant with the rule. Valid values:  `1`: critical `2`: warning `3`: info.
* `tag_key_scope` - (Optional) The rule monitors the tag key, only applies to rules created based on managed rules.
* `tag_value_scope` - (Optional) The rule monitors the tag value, use with the `tag_key_scope` options. only applies to rules created based on managed rules.
* `status` - (Optional, Available since v1.145.0) The rule status. The valid values: `ACTIVE`, `INACTIVE`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Aggregate Config Rule. The value is formatted `<aggregator_id>:<config_rule_id>`.
* `config_rule_id` - (Available since v1.141.0) The rule ID of Aggregate Config Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Aggregate Config Rule.
* `update` - (Defaults to 10 mins) Used when update the Aggregate Config Rule.

## Import

Cloud Config Aggregate Config Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_aggregate_config_rule.example "<aggregator_id>:<config_rule_id>"
```
