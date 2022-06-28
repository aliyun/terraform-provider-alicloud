---
subcategory: "Data Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_rules"
sidebar_current: "docs-alicloud-datasource-sddp-rules"
description: |-
  Provides a list of Sddp Rules to the user.
---

# alicloud\_sddp\_rules

This data source provides the Sddp Rules of the current Alibaba Cloud user.

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

data "alicloud_sddp_rules" "default" {
  ids = [alicloud_sddp_rule.default.id]
}
output "sddp_rule_id" {
  value = data.alicloud_sddp_rules.default.id
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Optional, ForceNew) Sensitive Data Identification Rules for the Type of. Valid values:
  * `0`: Keyword. 
  * `2`: Regular expression.
* `content_category` - (Optional, ForceNew) The Content Classification.
* `custom_type` - (Optional, ForceNew)  Sensitive Data Identification Rules of Type. Valid values: 
  * `0`: the Built-in.
  * `1`: The User-Defined.
* `ids` - (Optional, ForceNew, Computed)  A list of Rule IDs.
* `lang` - (Optional, ForceNew) The Request and Receive the Language of the Message Type. Valid values: 
  * `zh`: Chinese.
  * `en`: English.
* `name` - (Optional, ForceNew) The name of rule.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `product_id` - (Optional, ForceNew) Product ID. Valid values:
  * `1`:MaxCompute.
  * `2`:OSS.
  * `5`:RDS.
* `risk_level_id` - (Optional) Sensitive Data Identification Rules of Risk Level ID. Valid values:
  * `2`:S1, Weak Risk Level.
  * `3`:S2, Medium Risk Level. 
  * `4`:S3 High Risk Level. 
  * `5`:S4, the Highest Risk Level.
* `rule_type` - (Optional, ForceNew) Rule Type.
* `status` - (Optional, ForceNew) Sensitive Data Identification Rules Detection State of.
* `warn_level` - (Optional) The Level of Risk. Valid values:
  * `1`: Weak warn Level.
  * `2`: Medium Risk Level.
  * `3`: High Risk Level.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Rule names.
* `rules` - A list of Sddp Rules. Each element contains the following attributes:
	* `category` - Sensitive Data Identification Rules for the Type of.
	* `category_name` - Sensitive Data Identification Rules Belongs Type Name.
	* `content` - Sensitive Data Identification Rules the Content.
	* `content_category` - The Content Classification.
	* `create_time` - Sensitive Data Identification Rules the Creation Time of the Number of Milliseconds.
	* `custom_type` - Sensitive Data Identification Rules of Type. 0: the Built-in 1: The User-Defined.
	* `description` - Sensitive Data Identification a Description of the Rule Information.
	* `display_name` - Sensitive Data Identification Rules, Founder of Account Display Name.
	* `gmt_modified` - Sensitive Data Identification Rules to the Modified Time of the Number of Milliseconds.
	* `group_id` - Group ID.
	* `id` - The ID of the Rule.
	* `login_name` - Sensitive Data Identification Rules, Founder Of Account Login.
	* `major_key` - The Primary Key.
	* `product_code` - Product Code.
	* `product_id` - Product ID.
	* `risk_level_id` - Sensitive Data Identification Rules of Risk Level ID. Valid values:1:S1, Weak Risk Level. 2:S2, Medium Risk Level. 3:S3 High Risk Level. 4:S4, the Highest Risk Level.
	* `risk_level_name` - Sensitive Data Identification Rules the Risk Level of. S1: Weak Risk Level S2: Moderate Risk Level S3: High Risk Level S4: the Highest Risk Level.
	* `rule_id` - The first ID of the resource.
	* `stat_express` - Triggered the Alarm Conditions.
	* `status` - Sensitive Data Identification Rules Detection State of.
	* `target` - The Target.
	* `user_id` - The User ID.
	* `warn_level` - The Level of Risk.
