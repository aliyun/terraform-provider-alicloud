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

For information about Cloud Config Compliance Pack and how to use it, see [What is Compliance Pack](https://help.aliyun.com/).

-> **NOTE:** Available in v1.124.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_config_compliance_pack" "example" {
 compliance_pack_name        = "tf-testaccConfig1234"
  compliance_pack_template_id = "ct-3d20ff4e06a30027f76e"
  description                 = "tf-testaccConfig1234"
  risk_level                  = "1"
  config_rules {
    managed_rule_identifier = "ecs-snapshot-retention-days"
    config_rule_parameters {
      parameter_name  = "days"
      parameter_value = "7"
    }
  }
  config_rules {
    managed_rule_identifier = "ecs-instance-expired-check"
    config_rule_parameters {
      parameter_name  = "days"
      parameter_value = "60"
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `compliance_pack_name` - (Required, ForceNew) The Compliance Package Name.
* `compliance_pack_template_id` - (Optional, ForceNew) Compliance Package Template Id.
* `config_rules` - (Required) A list of Config Rules.
* `description` - (Required) The Description of compliance pack.
* `risk_level` - (Required) The Risk Level. Valid values:  `1`: critical, `2`: warning, `3`: info.

#### Block config_rules

The config_rules supports the following:

* `config_rule_parameters` - (Optional) A list of Config Rule Parameters.
* `managed_rule_identifier` - (Required) The Managed Rule Identifier.
* `config_rule_id` - (Optional, Computed, Available in v1.138.0+) The ID of the config rule.

#### Block config_rule_parameters

The config_rule_parameters supports the following:

* `parameter_name` - (Optional) The parameter name.
* `parameter_value` - (Optional) The parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Compliance Pack.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Compliance Pack.
* `update` - (Defaults to 2 mins) Used when update the Compliance Pack.

## Import

Cloud Config Compliance Pack can be imported using the id, e.g.

```
$ terraform import alicloud_config_compliance_pack.example <id>
```