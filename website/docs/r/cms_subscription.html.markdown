---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_subscription"
sidebar_current: "docs-alicloud-resource-cms-subscription"
description: |-
  Provides a Alicloud Cloud Monitor Service Subscription resource.
---

# alicloud_cms_subscription

Provides a Cloud Monitor Service Subscription resource.

For information about Cloud Monitor Service Subscription and how to use it, see [What is Subscription](https://www.alibabacloud.com/help/en/cms/developer-reference/api-cms-2024-03-30-createsubscription).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  description       = var.name
  enable            = true
  filter_setting {
    conditions {
      field = "product"
      value = "ECS"
      op    = "EQ"
    }
    expression = "product"
    relation   = "AND"
  }
}
```

## Argument Reference

The following arguments are supported:

* `subscription_name` - (Optional) The name of the subscription.
* `description` - (Optional) The description of the subscription.
* `notify_strategy_id` - (Optional) The ID of the notify strategy.
* `workspace` - (Optional) The workspace of the subscription.
* `enable` - (Optional, Available since v1.248.0) Whether to enable the subscription. Valid values: `true`, `false`. Default to `true`.
* `filter_setting` - (Optional) The filter setting of the subscription. See [`filter_setting`](#filter_setting) below.
* `pushing_setting` - (Optional) The pushing setting of the subscription. See [`pushing_setting`](#pushing_setting) below.

### filter_setting

The `filter_setting` supports the following:

* `conditions` - (Optional) A list of filter conditions. See [`conditions`](#conditions) below.
* `expression` - (Optional) The filter expression.
* `relation` - (Optional) The relation between conditions.

#### conditions

The `conditions` supports the following:

* `field` - (Required) The filter field.
* `value` - (Required) The filter value.
* `op` - (Required) The filter operator. Valid values: `IN`, `EQ`.

### pushing_setting

The `pushing_setting` supports the following:

* `alert_action_ids` - (Optional) A list of alert action plan IDs.
* `restore_action_ids` - (Optional) A list of restore action plan IDs.
* `template_uuid` - (Optional) The template UUID.
* `response_plan_id` - (Optional) The response plan ID.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the subscription.
* `subscription_id` - The ID of the subscription.
* `user_id` - The user ID.
* `create_time` - The creation time of the subscription.
* `update_time` - The update time of the subscription.
* `region_id` - The region ID of the subscription.
* `subscription_type` - The type of the subscription.

## Import

Cloud Monitor Service Subscription can be imported using the id, e.g.

```terraform
terraform import alicloud_cms_subscription.example <id>
```
