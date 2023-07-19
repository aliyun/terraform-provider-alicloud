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

```terraform
variable "name" {
  default = "tf-example-config"
}
data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_config_rule" "default" {
  description               = "If the ACL policy of the OSS bucket denies read access from the Internet, the configuration is considered compliant."
  source_owner              = "ALIYUN"
  source_identifier         = "oss-bucket-public-read-prohibited"
  risk_level                = 1
  tag_key_scope             = "For"
  tag_value_scope           = "example"
  region_ids_scope          = data.alicloud_regions.default.regions.0.id
  config_rule_trigger_types = "ConfigurationItemChangeNotification"
  resource_types_scope      = ["ACS::OSS::Bucket"]
  rule_name                 = "oss-bucket-public-read-prohibited"
}

resource "alicloud_config_compliance_pack" "default" {
  compliance_pack_name = var.name
  description          = var.name
  risk_level           = "1"
  config_rule_ids {
    config_rule_id = alicloud_config_rule.default.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `compliance_pack_name` - (Required) The Compliance Package Name. . **NOTE:** the `compliance_pack_name` supports modification since V1.146.0.
* `compliance_pack_template_id` - (Optional, ForceNew) Compliance Package Template Id.
* `config_rules` - (Optional form v1.141.0, Deprecated from v1.141.0) A list of Config Rules. See [`config_rules`](#config_rules) below. 
* `config_rule_ids` - (Optional, Available since v1.141.0) A list of Config Rule IDs. See [`config_rule_ids`](#config_rule_ids) below. 
* `description` - (Required) The Description of compliance pack.
* `risk_level` - (Required) The Risk Level. Valid values:  `1`: critical, `2`: warning, `3`: info.

### `config_rules`

The config_rules supports the following: 

* `config_rule_parameters` - (Optional) A list of Config Rule Parameters. See [`config_rule_parameters`](#config_rules-config_rule_parameters) below. 
* `managed_rule_identifier` - (Required) The Managed Rule Identifier.

### `config_rule_ids`

The config_rule_ids supports the following:

* `config_rule_id` - (Optional) The rule ID of Config Rule.

### `config_rules-config_rule_parameters`

The config_rule_parameters supports the following: 

* `parameter_name` - (Optional) The parameter name.
* `parameter_value` - (Optional) The parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Compliance Pack.
* `status` -  The status of the resource. The valid values: `CREATING`, `ACTIVE`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Compliance Pack.
* `update` - (Defaults to 2 mins) Used when update the Compliance Pack.

## Import

Cloud Config Compliance Pack can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_compliance_pack.example <id>
```
