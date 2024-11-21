---
subcategory: "Data Security Center (SDDP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_rule"
sidebar_current: "docs-alicloud-resource-sddp-rule"
description: |-
  Provides a Alicloud Data Security Center Rule resource.
---

# alicloud_sddp_rule

Provides a Data Security Center Rule resource.

For information about Data Security Center Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/en/data-security-center/latest/api-sddp-2019-01-03-createrule).

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sddp_rule&exampleId=3ff14e92-ed84-7505-2c85-ab6ba64d03964c1879c9&activeTab=example&spm=docs.r.sddp_rule.0.3ff14e92ed&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example-name"
}

resource "alicloud_sddp_rule" "default" {
  rule_name     = var.name
  category      = "2"
  content       = <<EOF
  [
    {
      "rule": [
        {
          "operator": "contains",
          "target": "content",
          "value": "tf-testACCContent"
        }
      ],
      "ruleRelation": "AND"
    }
  ]
  EOF
  risk_level_id = "4"
  product_code  = "OSS"
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required) The name of the sensitive data detection rule. **NOTE:** From version 1.222.0, `rule_name` can be modified.
* `category` - (Required, Int) The content type of the sensitive data detection rule. Valid values:
  - `0`: Keyword.
  - `2`: Regular expression.
**NOTE:** From version 1.222.0, `category` can be modified.
* `content` - (Required) The content of the sensitive data detection rule. **NOTE:** From version 1.222.0, `content` can be modified.
* `content_category` - (Optional, ForceNew) The type of the content in the sensitive data detection rule. **NOTE:** From version 1.222.0, `content_category` cannot be modified.
* `risk_level_id` - (Optional) The sensitivity level of the sensitive data that hits the sensitive data detection rule. Valid values:
  - `2`: S1, which indicates the low sensitivity level.
  - `3`: S2, which indicates the medium sensitivity level.
  - `4`: S3, which indicates the high sensitivity level.
  - `5`: S4, which indicates the highest sensitivity level.
* `rule_type` - (Optional, Int) The type of the sensitive data detection rule. Valid values:
  - `1`: Sensitive data detection rule.
  - `2`: Audit rule.
  - `3`: Anomalous event detection rule.
* `product_code` - (Optional) The name of the service to which data in the column of the table belongs. Valid values: `OSS`, `RDS`, `ODPS`(MaxCompute).
* `product_id` - (Optional) The ID of the service to which the data asset belongs. Valid values:
  - `1`:MaxCompute.
  - `2`:OSS.
  - `5`:RDS.
* `warn_level` - (Optional, Int) The risk level of the alert that is triggered. Valid values:
  - `1`: Low warn Level.
  - `2`: Medium Risk Level.
  - `3`: High Risk Level.
* `stat_express` - (Optional, ForceNew) The statistical expression. **NOTE:** From version 1.222.0, `stat_express` cannot be modified.
* `target` - (Optional, ForceNew) The code of the service to which the sensitive data detection rule is applied. **NOTE:** From version 1.222.0, `target` cannot be modified.
* `status` - (Optional) Sensitive Specifies whether to enable the sensitive data detection rule. Valid values:
  - `0`: Disable.
  - `1`: Enable.
* `description` - (Optional, ForceNew) The description of the rule. **NOTE:** From version 1.222.0, `description` cannot be modified.
* `lang` - (Optional) The language of the content within the request and response. Default value: `zh`. Valid values:
  - `zh`: Chinese.
  - `en`: English.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule.
* `custom_type` - The type of the sensitive data detection rule. **NOTE:** From version 1.222.0, `custom_type` cannot be specified when create Rule.

## Import

Data Security Center Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_sddp_rule.example <id>
```
