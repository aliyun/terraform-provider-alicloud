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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_aggregate_compliance_pack&exampleId=12808e93-aeee-b159-88cd-c757c2d5ffa6498d3279&activeTab=example&spm=docs.r.config_aggregate_compliance_pack.0.12808e93ae&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
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
  description                = var.name
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

* `aggregator_id` - (Required, ForceNew) The ID of aggregator.
* `aggregate_compliance_pack_name` - (Required) The name of compliance package name. **NOTE:** From version 1.145.0, `aggregate_compliance_pack_name` can be modified.
* `description` - (Required) The description of compliance package.
* `risk_level` - (Required, Int) The Risk Level. Valid values:
  - `1`: critical.
  - `2`: warning.
  - `3`: info.
* `compliance_pack_template_id` - (Optional, ForceNew, Available since v1.141.0) The Template ID of compliance package.
* `config_rule_ids` - (Optional, Set, Available since v1.141.0) A list of Config Rule IDs. See [`config_rule_ids`](#config_rule_ids) below.
* `config_rules` - (Optional, Set, Deprecated since v1.141.0) A list of Config Rules. See [`config_rules`](#config_rules) below. **NOTE:** Field `config_rules` has been deprecated from provider version 1.141.0. New field `config_rule_ids` instead.

### `config_rule_ids`

The config_rule_ids supports the following:

* `config_rule_id` - (Optional) The rule ID of Aggregate Config Rule.

### `config_rules`

The config_rules supports the following: 

* `managed_rule_identifier` - (Required) The Managed Rule Identifier.
* `config_rule_parameters` - (Optional, Set) A list of parameter rules. See [`config_rule_parameters`](#config_rules-config_rule_parameters) below.

### `config_rules-config_rule_parameters`

The config_rule_parameters supports the following: 

* `parameter_name` - (Optional) The Parameter Name.
* `parameter_value` - (Optional) The Parameter Value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Aggregate Compliance Pack. It formats as `<aggregator_id>:<aggregator_compliance_pack_id>`.
* `aggregator_compliance_pack_id` - The ID of the compliance package.
* `status` - The status of the Aggregate Compliance Pack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Aggregate Compliance Pack.
* `update` - (Defaults to 1 mins) Used when update the Aggregate Compliance Pack.
* `delete` - (Defaults to 1 mins) Used when delete the Aggregate Compliance Pack.

## Import

Cloud Config Aggregate Compliance Pack can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_aggregate_compliance_pack.example <aggregator_id>:<aggregator_compliance_pack_id>
```
