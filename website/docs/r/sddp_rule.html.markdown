---
subcategory: "Data Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_rule"
sidebar_current: "docs-alicloud-resource-sddp-rule"
description: |-
  Provides a Alicloud Data Security Center Rule resource.
---

# alicloud\_sddp\_rule

Provides a Data Security Center Rule resource.

For information about Data Security Center Rule and how to use it, see [What is Rule](https://help.aliyun.com/product/88674.html).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sddp_rule" "default" {
  category      = "0"
  content       = "content"
  rule_name     = "rule_name"
  risk_level_id = "4"
  product_code  = "ODPS"
}

```

## Argument Reference

The following arguments are supported:

* `category` - (Required, ForceNew) Sensitive Data Identification Rules for the Type of. Valid values:
  * `0`: Keyword.
  * `2`: Regular expression.
* `content` - (Required, ForceNew) Sensitive Data Identification Rules the Content.
* `content_category` - (Optional) The Content Classification.
* `custom_type` - (Optional, ForceNew) Sensitive Data Identification Rules of Type. Valid values: 
  * `0`: the Built-in.
  * `1`: The User-Defined.
* `description` - (Optional) Sensitive Data Identification a Description of the Rule Information.
* `lang` - (Optional) The Request and Receive the Language of the Message Type. Valid values:
  * `zh`: Chinese.
  * `en`: English.
* `product_code` - (Optional) Product Code. Valid values: `OSS`,`RDS`,`ODPS`(MaxCompute).
* `product_id` - (Optional) Product ID. Valid values: 
  * `1`:MaxCompute
  * `2`:OSS
  * `5`:RDS.
* `risk_level_id` - (Optional) Sensitive Data Identification Rules of Risk Level ID. Valid values: 
  * `2`:S1, Weak Risk Level. 
  * `3`:S2, Medium Risk Level. 
  * `4`:S3 High Risk Level. 
  * `5`:S4, the Highest Risk Level.
* `rule_name` - (Required, ForceNew) Sensitive Data Identification Name of the Rule.
* `rule_type` - (Optional) Rule Type.
* `stat_express` - (Optional) Triggered the Alarm Conditions.
* `status` - (Optional, Computed) Sensitive Data Identification Rules Detection State of.
* `target` - (Optional) The Target of rule.
* `warn_level` - (Optional) The Level of Risk. Valid values: 
  * `1`: Weak warn Level. 
  * `2`: Medium Risk Level. 
  * `3`: High Risk Level.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule.

## Import

Data Security Center Rule can be imported using the id, e.g.

```
$ terraform import alicloud_sddp_rule.example <id>
```
