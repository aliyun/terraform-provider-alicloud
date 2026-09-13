---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_compliance_pack"
sidebar_current: "docs-alicloud-resource-config-compliance-pack"
description: |-
  Provides a Alicloud Cloud Config Compliance Pack resource.
---

# alicloud_config_compliance_pack

Provides a Cloud Config Compliance Pack resource.

For information about Cloud Config Compliance Pack and how to use it, see [What is Compliance Pack](https://www.alibabacloud.com/help/en/cloud-config/latest/api-config-2020-09-07-createcompliancepack).

-> **NOTE:** Available since v1.124.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_compliance_pack&exampleId=c021cc95-b828-dec8-ff1c-e468119dc9065d5bbe9c&activeTab=example&spm=docs.r.config_compliance_pack.0.c021cc95b8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example-config-name"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_config_rule" "rule1" {
  description                 = var.name
  source_owner                = "ALIYUN"
  source_identifier           = "ram-user-ak-create-date-expired-check"
  risk_level                  = 1
  maximum_execution_frequency = "TwentyFour_Hours"
  region_ids_scope            = data.alicloud_regions.default.regions.0.id
  config_rule_trigger_types   = "ScheduledNotification"
  resource_types_scope        = ["ACS::RAM::User"]
  rule_name                   = "ciscompliancecheck_ram-user-ak-create-date-expired-check"
  input_parameters = {
    days = "90"
  }
}

resource "alicloud_config_rule" "rule2" {
  description               = var.name
  source_owner              = "ALIYUN"
  source_identifier         = "adb-cluster-maintain-time-check"
  risk_level                = 2
  region_ids_scope          = data.alicloud_regions.default.regions.0.id
  config_rule_trigger_types = "ScheduledNotification"
  resource_types_scope      = ["ACS::ADB::DBCluster"]
  rule_name                 = "governance-evaluation-adb-cluster-maintain-time-check"
  input_parameters = {
    maintainTimes = "02:00-04:00,06:00-08:00,12:00-13:00"
  }
}

resource "alicloud_config_compliance_pack" "default" {
  compliance_pack_name = var.name
  description          = "CloudGovernanceCenter evaluation"
  risk_level           = "2"
  config_rule_ids {
    config_rule_id = alicloud_config_rule.rule1.id
  }
  config_rule_ids {
    config_rule_id = alicloud_config_rule.rule2.id
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_config_compliance_pack&spm=docs.r.config_compliance_pack.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `compliance_pack_name` - (Required) The Compliance Package Name. **NOTE:** From version 1.146.0, `compliance_pack_name` can be modified.
* `description` - (Required) The Description of compliance pack.
* `risk_level` - (Required, Int) The Risk Level. Valid values:
  - `1`: critical.
  - `2`: warning.
  - `3`: info.
* `compliance_pack_template_id` - (Optional, ForceNew) Compliance Package Template Id.
* `template_content` - (Optional, ForceNew) The template content of the compliance pack.
* `config_rule_ids` - (Optional, Set, Available since v1.141.0) A list of Config Rule IDs. See [`config_rule_ids`](#config_rule_ids) below.
* `config_rules` - (Optional, Set, Deprecated since v1.141.0) A list of Config Rules. See [`config_rules`](#config_rules) below. **NOTE:** Field `config_rules` has been deprecated from provider version 1.141.0. New field `config_rule_ids` instead.
* `tag_key_scope` - (Optional) The tag key scope of the compliance pack.
* `tag_value_scope` - (Optional) The tag value scope of the compliance pack.
* `resource_ids_scope` - (Optional) The resource IDs scope of the compliance pack.
* `exclude_resource_ids_scope` - (Optional) The exclude resource IDs scope of the compliance pack.
* `resource_group_ids_scope` - (Optional) The resource group IDs scope of the compliance pack.
* `exclude_resource_group_ids_scope` - (Optional) The exclude resource group IDs scope of the compliance pack.
* `region_ids_scope` - (Optional) The region IDs scope of the compliance pack.
* `exclude_region_ids_scope` - (Optional) The exclude region IDs scope of the compliance pack.
* `tags_scope` - (Optional, List) The tags scope of the compliance pack. See [`tags_scope`](#tags_scope) below.
* `exclude_tags_scope` - (Optional, List) The exclude tags scope of the compliance pack. See [`exclude_tags_scope`](#exclude_tags_scope) below.

### `config_rule_ids`

The config_rule_ids supports the following:

* `config_rule_id` - (Optional) The rule ID of Config Rule.

### `config_rules`

The config_rules supports the following:

* `managed_rule_identifier` - (Required) The Managed Rule Identifier.
* `config_rule_parameters` - (Optional, Set) A list of Config Rule Parameters. See [`config_rule_parameters`](#config_rules-config_rule_parameters) below.
* `config_rule_name` - (Optional) The name of the Config Rule.
* `description` - (Optional) The description of the Config Rule.
* `risk_level` - (Optional, Int) The risk level of the Config Rule. Valid values: `1`, `2`, `3`.

### `config_rules-config_rule_parameters`

The config_rule_parameters supports the following:

* `parameter_name` - (Optional) The parameter name.
* `parameter_value` - (Optional) The parameter value.

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

* `id` - The resource ID in terraform of Compliance Pack.
* `status` -  The status of the Compliance Pack.
* `template_content` - The template content of the compliance pack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Compliance Pack.
* `update` - (Defaults to 2 mins) Used when update the Compliance Pack.
* `delete` - (Defaults to 1 mins) Used when delete the Compliance Pack.

## Import

Cloud Config Compliance Pack can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_compliance_pack.example <id>
```
