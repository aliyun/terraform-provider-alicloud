---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_compliance_pack"
sidebar_current: "docs-alicloud-resource-config-aggregate-compliance-pack"
description: |-
  Provides a Alicloud Cloud Config Aggregate Compliance Pack resource.
---

# alicloud_config_aggregate_compliance_pack

Provides a Cloud Config Aggregate Compliance Pack resource.

For information about Cloud Config Aggregate Compliance Pack and how to use it, see [What is Aggregate Compliance Pack](https://www.alibabacloud.com/help/en/cloud-config/latest/api-config-2020-09-07-createaggregatecompliancepack).

-> **NOTE:** Available since v1.124.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}
data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}

resource "alicloud_config_aggregator" "default" {
  aggregator_accounts {
    account_id   = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
    account_name = data.alicloud_resource_manager_accounts.default.accounts.0.display_name
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

resource "alicloud_config_aggregate_compliance_pack" "default" {
  aggregate_compliance_pack_name = var.name
  aggregator_id                  = alicloud_config_aggregator.default.id
  description                    = var.name
  risk_level                     = 1
  config_rule_ids {
    config_rule_id = alicloud_config_aggregate_config_rule.default.config_rule_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `aggregate_compliance_pack_name` - (Required)The name of compliance package name. **NOTE:** the `aggregate_compliance_pack_name` supports modification since V1.145.0.
* `aggregator_id` - (Required, ForceNew)The ID of aggregator.
* `compliance_pack_template_id` - (Optional from v1.141.0, ForceNew)The Template ID of compliance package.
* `config_rules` - (Optional, Deprecated from v1.141.0) A list of Config Rules. See [`config_rules`](#config_rules) below. 
* `config_rule_ids` - (Optional, Available since v1.141.0) A list of Config Rule IDs. See [`config_rule_ids`](#config_rule_ids) below. 
* `description` - (Required) The description of compliance package.
* `risk_level` - (Required) The Risk Level. Valid values: `1`: critical `2`: warning `3`: info.

### `config_rules`

The config_rules supports the following: 

* `config_rule_parameters` - (Optional) A list of parameter rules. See [`config_rule_parameters`](#config_rules-config_rule_parameters) below. 
* `managed_rule_identifier` - (Required) The Managed Rule Identifier.

### `config_rule_ids`

The config_rule_ids supports the following:

* `config_rule_id` - (Optional) The rule ID of Aggregate Config Rule.

### `config_rules-config_rule_parameters`

The config_rule_parameters supports the following: 

* `parameter_name` - (Optional) The Parameter Name.
* `parameter_value` - (Optional) The Parameter Value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Aggregate Compliance Pack. The value is formatted `<aggregator_id>:<aggregator_compliance_pack_id>`.
* `status` - The status of the resource. The valid values: `CREATING`, `ACTIVE`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Aggregate Compliance Pack.
* `update` - (Defaults to 1 mins) Used when update the Aggregate Compliance Pack.

## Import

Cloud Config Aggregate Compliance Pack can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_aggregate_compliance_pack.example <aggregator_id>:<aggregator_compliance_pack_id>
```
