---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_rules"
sidebar_current: "docs-alicloud-datasource-config-rules"
description: |-
    Provides a list of Config Rules to the user.
---

# alicloud\_config\_rules

This data source provides the Config Rules of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.99.0+.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-northeast-1`.

## Example Usage

```terraform
data "alicloud_config_rules" "example" {
  ids        = ["cr-ed4bad756057********"]
  name_regex = "tftest"
}

output "first_config_rule_id" {
  value = data.alicloud_config_rules.example.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Config Rule IDs.
* `config_rule_state` - (Optional, ForceNew) The state of the config rule, valid values: `ACTIVE`, `DELETING`, `DELETING_RESULTS`, `EVALUATING` and `INACTIVE`. 
* `multi_account` - (Optional, ForceNew) Whether the enterprise management account queries the rule details of member accounts.
* `member_id` - (Optional, ForceNew) The ID of the member account to which the rule to be queried belongs. The default is empty. When `multi_account` is set to true, this parameter is valid.
* `risk_level` - (Optional, ForceNew) The risk level of Config Rule. Valid values: `1`: Critical ,`2`: Warning , `3`: Info.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

-> **NOTE:** When you use the personal version to configure auditing, please ignore `multi_account` and `member_id`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Config Rule IDs.
* `names` - A list of Config Rule names.
* `rules` - A list of Config Rules. Each element contains the following attributes:
    * `id` - The ID of the Config Rule.
    * `account_id`- The ID of the Alicloud account.
    * `config_rule_arn`- The ARN of the Config Rule.
    * `config_rule_id`- The ID of the Config Rule.
    * `config_rule_state`- The state of the Config Rule.
    * `create_timestamp`- The timestamp of the Config Rule created.
    * `description`- The description of the Config Rule.
    * `input_parameters`- The input paramrters of the Config Rule.
    * `modified_timestamp`- the timestamp of the Config Rule modified.
    * `risk_level`- The risk level of the Config Rule.
    * `rule_name`- The name of the Config Rule.
    * `source_details`- The source details of the Config Rule.
        * `event_source` - Event source of the Config Rule.
        * `maximum_execution_frequency` - Rule execution cycle.
        * `message_type` - Rule trigger mechanism.
    * `source_identifier`- The name of the custom rule or managed rule.
    * `source_owner`- The source owner of the Config Rule.
