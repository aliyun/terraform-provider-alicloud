---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_config_rule"
sidebar_current: "docs-alicloud-resource-config-aggregate-config-rule"
description: |-
  Provides a Alicloud Cloud Config Aggregate Config Rule resource.
---

# alicloud\_config\_aggregate\_config\_rule

Provides a Cloud Config Aggregate Config Rule resource.

For information about Cloud Config Aggregate Config Rule and how to use it, see [What is Aggregate Config Rule](https://help.aliyun.com/).

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_config_aggregator" "example" {
  aggregator_accounts {
    account_id   = "140278452670****"
    account_name = "test-2"
    account_type = "ResourceDirectory"
  }
  aggregator_name = "tf-testaccaggregator"
  description     = "tf-testaccaggregator"
}

resource "alicloud_config_aggregate_config_rule" "example" {
  aggregate_config_rule_name      = "tf-testaccconfig1234"
  aggregator_id                   = alicloud_config_aggregator.example.aggregator_id
  config_rule_trigger_types       = "ConfigurationItemChangeNotification"
  source_owner                    = "ALIYUN"
  source_identifier               = "ecs-cpu-min-count-limit"
  risk_level                      = 1
  resource_types_scope            = ["ACS::ECS::Instance"]
  input_parameters                = {
        cpuCount = "4",
  }
}

```

## Argument Reference

The following arguments are supported:

* `aggregate_config_rule_name` - (Required, ForceNew) The name of the rule.
* `aggregator_id` - (Required, ForceNew) The Aggregator Id.
* `config_rule_trigger_types` - (Required) The config rule trigger types. Valid values `ConfigurationItemChangeNotification`, `ScheduledNotification`.
* `description` - (Optional) The description of the rule.
* `exclude_resource_ids_scope` - (Optional) Exclude ResourceId List.
* `input_parameters` - (Optional) The settings map of the input parameters for the rule.
* `source_identifier`- (Required, ForceNew) The name of the custom rule or managed rule.
* `source_owner`- (Required, ForceNew) The source owner of the Config Rule. Valid values `ALIYUN` and `CUSTOM_FC`.
* `maximum_execution_frequency` - (Optional) The frequency of the compliance evaluations. Valid values:  `One_Hour`, `Three_Hours`, `Six_Hours`, `Twelve_Hours`, `TwentyFour_Hours`.
* `region_ids_scope` - (Optional) The region ids scope.
* `resource_group_ids_scope` - (Optional) The resource group ids scope.
* `resource_types_scope` - (Required) The types of the resources to be evaluated against the rule.
* `risk_level` - (Required) The risk level of the resources that are not compliant with the rule. Valid values:  `1`: critical `2`: warning `3`: info.
* `tag_key_scope` - (Optional) The tag key scope.
* `tag_value_scope` - (Optional) The tag value scope.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Aggregate Config Rule. The value is formatted `<aggregator_id>:<config_rule_id>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Aggregate Config Rule.
* `update` - (Defaults to 6 mins) Used when update the Aggregate Config Rule.

## Import

Cloud Config Aggregate Config Rule can be imported using the id, e.g.

```
$ terraform import alicloud_config_aggregate_config_rule.example <aggregator_id>:<config_rule_id>
```