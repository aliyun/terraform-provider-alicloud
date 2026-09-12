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
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_aggregate_compliance_pack&exampleId=77f1d8f8-6c27-15f2-ad02-1e43d5f1c29fb33c6349&activeTab=example&spm=docs.r.config_aggregate_compliance_pack.0.77f1d8f86c&intl_lang=EN_US" target="_blank">
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_config_aggregate_compliance_pack&spm=docs.r.config_aggregate_compliance_pack.example&intl_lang=EN_US)

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
* `template_content` - (Optional, ForceNew) The template content of compliance package. It is a JSON string that defines the compliance pack rules. If not specified, the template content is derived from `compliance_pack_template_id`.
* `config_rule_ids` - (Optional, Set, Available since v1.141.0) A list of Config Rule IDs. See [`config_rule_ids`](#config_rule_ids) below.
* `config_rules` - (Optional, Set, Deprecated since v1.141.0) A list of Config Rules. See [`config_rules`](#config_rules) below. **NOTE:** Field `config_rules` has been deprecated from provider version 1.141.0. New field `config_rule_ids` instead.
* `tag_key_scope` - (Optional) The tag key scope. The compliance pack only applies to resources with the specified tag key.
* `tag_value_scope` - (Optional) The tag value scope. The compliance pack only applies to resources with the specified tag value.
* `resource_ids_scope` - (Optional) The resource IDs scope. The compliance pack only applies to the specified resource IDs.
* `exclude_resource_ids_scope` - (Optional) The exclude resource IDs scope. The compliance pack does not apply to the specified resource IDs.
* `resource_group_ids_scope` - (Optional) The resource group IDs scope. The compliance pack only applies to resources in the specified resource groups.
* `exclude_resource_group_ids_scope` - (Optional) The exclude resource group IDs scope. The compliance pack does not apply to resources in the specified resource groups.
* `region_ids_scope` - (Optional) The region IDs scope. The compliance pack only applies to resources in the specified regions.
* `exclude_region_ids_scope` - (Optional) The exclude region IDs scope. The compliance pack does not apply to resources in the specified regions.
* `tags_scope` - (Optional) The tags scope. The compliance pack only applies to resources with the specified tags. See [`tags_scope`](#tags_scope) below.
* `exclude_tags_scope` - (Optional) The exclude tags scope. The compliance pack does not apply to resources with the specified tags. See [`exclude_tags_scope`](#exclude_tags_scope) below.

### `config_rule_ids`

The config_rule_ids supports the following:

* `config_rule_id` - (Optional) The rule ID of Aggregate Config Rule.

### `config_rules`

The config_rules supports the following:

* `managed_rule_identifier` - (Required) The Managed Rule Identifier.
* `config_rule_name` - (Optional) The name of the config rule.
* `description` - (Optional) The description of the config rule.
* `risk_level` - (Optional, Int) The risk level of the config rule. Valid values: `1`, `2`, `3`.
* `config_rule_parameters` - (Optional, Set) A list of parameter rules. See [`config_rule_parameters`](#config_rules-config_rule_parameters) below.

### `config_rules-config_rule_parameters`

The config_rule_parameters supports the following:

* `parameter_name` - (Optional) The Parameter Name.
* `parameter_value` - (Optional) The Parameter Value.

### `tags_scope`

The tags_scope supports the following:

* `tag_key` - (Optional) The tag key.
* `tag_value` - (Optional) The tag value.

### `exclude_tags_scope`

The exclude_tags_scope supports the following:

* `tag_key` - (Optional) The tag key.
* `tag_value` - (Optional) The tag value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Aggregate Compliance Pack. It formats as `<aggregator_id>:<aggregator_compliance_pack_id>`.
* `aggregator_compliance_pack_id` - The ID of the compliance package.
* `status` - The status of the Aggregate Compliance Pack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Aggregate Compliance Pack.
* `update` - (Defaults to 1 mins) Used when update the Aggregate Compliance Pack.
* `delete` - (Defaults to 1 mins) Used when delete the Aggregate Compliance Pack.

## Import

Cloud Config Aggregate Compliance Pack can be imported using the id, which consists of aggregator_id and aggregator_compliance_pack_id, e.g.

```shell
$ terraform import alicloud_config_aggregate_compliance_pack.example <aggregator_id>:<aggregator_compliance_pack_id>
```
