---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_compliance_pack"
sidebar_current: "docs-alicloud-resource-config-aggregate-compliance-pack"
description: |-
  Provides a Alicloud Cloud Config Aggregate Compliance Pack resource.
---

# alicloud\_config\_aggregate\_compliance\_pack

Provides a Cloud Config Aggregate Compliance Pack resource.

For information about Cloud Config Aggregate Compliance Pack and how to use it, see [What is Aggregate Compliance Pack](https://help.aliyun.com/).

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_name"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_instances" "default" {}

resource "alicloud_config_aggregator" "default" {
  aggregator_accounts {
    account_id   = "140278452670****"
    account_name = "test-2"
    account_type = "ResourceDirectory"
  }
  aggregator_name = "tf-testaccaggregator"
  description     = "tf-testaccaggregator"
}


resource "alicloud_config_aggregate_config_rule" "default" {
  aggregator_id              = alicloud_config_aggregator.default.id
  aggregate_config_rule_name = var.name
  source_owner               = "ALIYUN"
  source_identifier          = "ecs-cpu-min-count-limit"
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  resource_types_scope       = ["ACS::ECS::Instance"]
  risk_level                 = 1
  description                = var.name
  exclude_resource_ids_scope = data.alicloud_instances.default.ids.0
  input_parameters = {
    cpuCount = "4",
  }
  region_ids_scope         = "cn-hangzhou"
  resource_group_ids_scope = data.alicloud_resource_manager_resource_groups.default.ids.0
  tag_key_scope            = "tFTest"
  tag_value_scope          = "forTF 123"
}

resource "alicloud_config_aggregate_compliance_pack" "default" {
  aggregate_compliance_pack_name = "tf-testaccConfig1234"
  aggregator_id                  = alicloud_config_aggregator.default.id
  description                    = "tf-testaccConfig1234"
  risk_level                     = 1
  config_rule_ids {
    config_rule_id = alicloud_config_aggregate_config_rule.default.config_rule_id
  }
}

```

## Argument Reference

The following arguments are supported:

* `aggregate_compliance_pack_name` - (Required, ForceNew)The name of compliance package name.
* `aggregator_id` - (Required, ForceNew)The ID of aggregator.
* `compliance_pack_template_id` - (Optional form v1.141.0, ForceNew)The Template ID of compliance package.
* `config_rules` - (Optional, Computed, Deprecated form v1.141.0) A list of Config Rules.
* `config_rule_ids` - (Optional, Computed, Available in v1.141.0) A list of Config Rule IDs.
* `description` - (Required) Teh description of compliance package.
* `risk_level` - (Required) The Risk Level. Valid values: `1`, `2`, `3`.

#### Block config_rules

The config_rules supports the followingsss: 

* `config_rule_parameters` - (Optional) A list of parameter rules.
* `managed_rule_identifier` - (Required) The Managed Rule Identifier.

#### Block config_rule_ids

The config_rule_ids supports the followings:

* `config_rule_id` - (Optional) The rule ID of Aggregate Config Rule.

#### Block config_rule_parameters

The config_rule_parameters supports the following: 

* `parameter_name` - (Optional) The Parameter Name.
* `parameter_value` - (Optional) The Parameter Value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Aggregate Compliance Pack. The value is formatted `<aggregator_id>:<aggregator_compliance_pack_id>`.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Aggregate Compliance Pack.
* `update` - (Defaults to 1 mins) Used when update the Aggregate Compliance Pack.

## Import

Cloud Config Aggregate Compliance Pack can be imported using the id, e.g.

```
$ terraform import alicloud_config_aggregate_compliance_pack.example <aggregator_id>:<aggregator_compliance_pack_id>
```