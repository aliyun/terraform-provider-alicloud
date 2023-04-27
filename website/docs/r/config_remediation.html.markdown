---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_remediation"
sidebar_current: "docs-alicloud-resource-config-remediation"
description: |-
  Provides a Alicloud Config Remediation resource.
---

# alicloud_config_remediation

Provides a Config Remediation resource.

For information about Config Remediation and how to use it, see [What is Remediation](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available in v1.204.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_config_remediation" "default" {
  config_rule_id          = alicloud_config_rule.prerequirement-rule.config_rule_id
  remediation_template_id = "ACS-TAG-TagResources"
  remediation_source_type = "ALIYUN"
  invoke_type             = "MANUAL_EXECUTION"
  params                  = "{\"regionId\":\"{regionId}\",\"tags\":\"{\\\"terraform\\\":\\\"terraform\\\"}\",\"resourceType\":\"{resourceType}\",\"resourceIds\":\"{resourceId}\"}"
  remediation_type        = "OOS"
}
```

## Argument Reference

The following arguments are supported:
* `config_rule_id` - (Required, ForceNew) Rule ID.
* `invoke_type` - (Required) Execution type, valid values: `Manual`, `Automatic`.
* `params` - (Required, JsonString) Remediation parameter.
* `remediation_source_type` - (ForceNew, Computed, Optional) Remediation resource type, valid values: `ALIYUN` , `CUSTOMER`.
* `remediation_template_id` - (Required) Remediation template ID.
* `remediation_type` - (Required, ForceNew) Remediation type, valid values: `OOS`, `FC`.

The following arguments will be discarded. Please use new fields as soon as possible:



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `remediation_id` - Remediation ID.
* `remediation_source_type` - Remediation resource type, valid values: `ALIYUN` , `CUSTOMER`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Remediation.
* `delete` - (Defaults to 5 mins) Used when delete the Remediation.
* `update` - (Defaults to 5 mins) Used when update the Remediation.

## Import

Config Remediation can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_remediation.example <id>
```