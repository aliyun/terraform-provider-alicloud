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
resource "alicloud_config_aggregator" "example" {
  aggregator_accounts {
    account_id   = "140278452670****"
    account_name = "test-2"
    account_type = "ResourceDirectory"
  }
  aggregator_name = "tf-testaccaggregator"
  description     = "tf-testaccaggregator"
}

resource "alicloud_config_aggregate_compliance_pack" "example" {
  aggregate_compliance_pack_name = "tf-testaccConfig1234"
  aggregator_id                  = alicloud_config_aggregators.example.id
  compliance_pack_template_id    = "ct-3d20ff4e06a30027f76e"
  description                    = "tf-testaccConfig1234"
  risk_level                     = 1
  config_rules {
    managed_rule_identifier = "ecs-instance-expired-check"
    config_rule_parameters {
      parameter_name  = "days"
      parameter_value = "60"
    }
  }
  config_rules {
    managed_rule_identifier = "ecs-snapshot-retention-days"
    config_rule_parameters {
      parameter_name  = "days"
      parameter_value = "7"
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `aggregate_compliance_pack_name` - (Required, ForceNew)The name of compliance package name.
* `aggregator_id` - (Required, ForceNew)The ID of aggregator.
* `compliance_pack_template_id` - (Required, ForceNew)The Template ID of compliance package.
* `config_rules` - (Required) A list of  compliance package rules.
* `description` - (Required) Teh description of compliance package.
* `risk_level` - (Required) The Risk Level. Valid values: `1`, `2`, `3`.

#### Block config_rules

The config_rules supports the following: 

* `config_rule_parameters` - (Optional) A list of parameter rules.
* `managed_rule_identifier` - (Required) The Managed Rule Identifier.

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