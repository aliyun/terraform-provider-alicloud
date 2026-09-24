---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_event_integration_policies"
description: |-
  Provides a list of Cms Alert Event Integration Policy to the user.
---

<!-- markdownlint-configure-file {"MD007": {"indent": 4}} -->

# alicloud_cms_alert_event_integration_policies

This data source provides the Cms Alert Event Integration Policies of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_alert_event_integration_policies" "default" {
  workspace = alicloud_cms_workspace.default.id
  ids       = ["example-id"]
}

output "first_policy_id" {
  value = data.alicloud_cms_alert_event_integration_policies.default.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of alert event integration policy IDs.
* `name_regex` - (Optional) A regex string to filter results by the alert event integration policy name.
* `workspace` - (Optional) The workspace ID to which the alert event integration policies belong.
* `enable` - (Optional) Whether to filter policies by the enable status.
* `output_file` - (Optional) Save the result to a file.

## Attributes Reference

The following attributes are exported:

* `names` - A list of alert event integration policy names.
* `policies` - A list of alert event integration policies. Each element contains the following attributes:
    * `id` - The ID of the alert event integration policy.
    * `alert_event_integration_policy_name` - The name of the alert event integration policy.
    * `description` - The description of the alert event integration policy.
    * `type` - The type of the alert event integration policy.
    * `workspace` - The workspace ID of the alert event integration policy.
    * `integration_setting` - The integration setting of the alert event integration policy.
    * `enable` - Whether the alert event integration policy is enabled.
    * `region_id` - The region ID of the resource.
    * `create_time` - The creation time of the resource.
    * `update_time` - The update time of the resource.
    * `user_id` - The user ID of the resource.
    * `filter_setting` - The filter setting of the alert event integration policy.
        * `conditions` - The list of filter conditions.
            * `field` - The field name of the filter condition.
            * `value` - The field value of the filter condition.
            * `op` - The comparison operator of the filter condition.
        * `expression` - The filter expression.
        * `relation` - The relation between filter conditions.
    * `transformer_setting` - The transformer setting of the alert event integration policy.
        * `filter_setting` - The filter setting of the transformer.
            * `conditions` - The list of filter conditions.
                * `field` - The field name of the filter condition.
                * `value` - The field value of the filter condition.
                * `op` - The comparison operator of the filter condition.
            * `expression` - The filter expression.
            * `relation` - The relation between filter conditions.
        * `label_key` - The label key of the transformer.
        * `mapping` - The mapping table of the transformer.
        * `reg_exp` - The text extract regular expression of the transformer.
        * `source` - The source path of the transformer.
        * `target` - The transform target of the transformer.
        * `type` - The action type of the transformer.
        * `value` - The value of the transformer.
        * `variable` - The variable name of the transformer.
* `next_token` - The next token for pagination.
