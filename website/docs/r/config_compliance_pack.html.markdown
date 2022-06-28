---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_compliance_pack"
sidebar_current: "docs-alicloud-resource-config-compliance-pack"
description: |-
  Provides a Alicloud Cloud Config Compliance Pack resource.
---

# alicloud\_config\_compliance\_pack

Provides a Cloud Config Compliance Pack resource.

For information about Cloud Config Compliance Pack and how to use it, see [What is Compliance Pack](https://www.alibabacloud.com/help/en/doc-detail/194753.html).

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_name"
}

data "alicloud_instances" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_config_rule" "default" {
  rule_name                  = var.name
  description                = var.name
  source_identifier          = "ecs-instances-in-vpc"
  source_owner               = "ALIYUN"
  resource_types_scope       = ["ACS::ECS::Instance"]
  risk_level                 = 1
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  tag_key_scope              = "tfTest"
  tag_value_scope            = "tfTest 123"
  resource_group_ids_scope   = data.alicloud_resource_manager_resource_groups.default.ids.0
  exclude_resource_ids_scope = data.alicloud_instances.default.instances[0].id
  region_ids_scope           = "cn-hangzhou"
  input_parameters = {
    vpcIds = data.alicloud_instances.default.instances[0].vpc_id
  }
}

resource "alicloud_config_compliance_pack" "default" {
  compliance_pack_name = "tf-testaccConfig1234"
  description          = "tf-testaccConfig1234"
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
* `config_rules` - (Optional form v1.141.0, Computed, Deprecated form v1.141.0) A list of Config Rules.
* `config_rule_ids` - (Optional, Computed, Available in v1.141.0) A list of Config Rule IDs.
* `description` - (Required) The Description of compliance pack.
* `risk_level` - (Required) The Risk Level. Valid values:  `1`: critical, `2`: warning, `3`: info.

#### Block config_rules

The config_rules supports the following: 

* `config_rule_parameters` - (Optional) A list of Config Rule Parameters.
* `managed_rule_identifier` - (Required) The Managed Rule Identifier.

#### Block config_rule_ids

The config_rule_ids supports the following:

* `config_rule_id` - (Optional) The rule ID of Config Rule.

#### Block config_rule_parameters

The config_rule_parameters supports the following: 

* `parameter_name` - (Optional) The parameter name.
* `parameter_value` - (Optional) The parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Compliance Pack.
* `status` -  The status of the resource. The valid values: `CREATING`, `ACTIVE`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Compliance Pack.
* `update` - (Defaults to 2 mins) Used when update the Compliance Pack.

## Import

Cloud Config Compliance Pack can be imported using the id, e.g.

```
$ terraform import alicloud_config_compliance_pack.example <id>
```
