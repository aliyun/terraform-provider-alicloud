---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_remediation"
description: |-
  Provides a Alicloud Cloud Config (Config) Aggregate Remediation resource.
---

# alicloud_config_aggregate_remediation

Provides a Cloud Config (Config) Aggregate Remediation resource.

Rule remediation in multi-account scenarios.

For information about Cloud Config (Config) Aggregate Remediation and how to use it, see [What is Aggregate Remediation](https://next.api.alibabacloud.com/document/Config/2020-09-07/CreateAggregateRemediation).

-> **NOTE:** Available since v1.267.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_config_aggregator" "create-agg" {
  aggregator_name = "rd"
  description     = "rd"
  aggregator_type = "RD"
}

resource "alicloud_config_aggregate_config_rule" "create-rule" {
  source_owner               = "ALIYUN"
  source_identifier          = "required-tags"
  aggregate_config_rule_name = "agg-rule-name"
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  risk_level                 = "1"
  resource_types_scope       = ["ACS::OSS::Bucket"]
  aggregator_id              = alicloud_config_aggregator.create-agg.id
  input_parameters = {
    tag1Key   = "aaa"
    tag1Value = "bbb"
  }
}


resource "alicloud_config_aggregate_remediation" "default" {
  config_rule_id          = alicloud_config_aggregate_config_rule.create-rule.config_rule_id
  remediation_template_id = "ACS-TAG-TagResources"
  remediation_source_type = "ALIYUN"
  invoke_type             = "MANUAL_EXECUTION"
  remediation_type        = "OOS"
  aggregator_id           = alicloud_config_aggregator.create-agg.id
  remediation_origin_params = jsonencode({
    properties = [
      {
        name          = "regionId"
        type          = "String"
        value         = "{regionId}"
        allowedValues = []
        description   = "region ID"
      },
      {
        name          = "tags"
        type          = "Json"
        value         = "{\"aaa\":\"bbb\"}"
        allowedValues = []
        description   = "resource tags (for example,{\"k1\":\"v1\",\"k2\":\"v2\"})."
      },
      {
        name          = "resourceType"
        type          = "String"
        value         = "{resourceType}"
        allowedValues = []
        description   = "resource type"
      },
      {
        name          = "resourceIds"
        type          = "ARRAY"
        value         = "[{\"resources\":[]}]"
        allowedValues = []
        description   = "Resource ID List"
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:
* `aggregator_id` - (Required, ForceNew) The account Group ID.
For more information about how to obtain the account group ID, see [ListAggregators](~~ 255797 ~~).
* `config_rule_id` - (Required, ForceNew) The rule ID.
For more information about how to obtain the rule ID, see [ListAggregateConfigRules].
* `invoke_type` - (Required) Correction of execution mode. Value:
  - NON_EXECUTION: Not executed.
  - AUTO_EXECUTION: Automatically executed.
  - MANUAL_EXECUTION: Execute manually.
  - NOT_CONFIG: Not set.
* `remediation_origin_params` - (Required) Correct the parameters of the settings.
For more information about how to obtain the parameters of remediation settings, see the parameter 'Template definition' in [ListRemediationTemplates](~~ 416781 ~~) '.
* `remediation_source_type` - (Optional, ForceNew) The source of the template to perform the correction. Value:
  - ALIYUN (default): Official website template.
  - CUSTOM: CUSTOM template.
  - NONE: NONE.
* `remediation_template_id` - (Required) The ID of the correction template.
* `remediation_type` - (Required, ForceNew) Remediation type. Value:
  - OOS: Operation and maintenance orchestration (Template correction).
  - FC: Function Compute (custom correction).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<aggregator_id>:<remediation_id>`.
* `remediation_id` - Multi-account remediation ID

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Aggregate Remediation.
* `delete` - (Defaults to 5 mins) Used when delete the Aggregate Remediation.
* `update` - (Defaults to 5 mins) Used when update the Aggregate Remediation.

## Import

Cloud Config (Config) Aggregate Remediation can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_aggregate_remediation.example <aggregator_id>:<remediation_id>
```