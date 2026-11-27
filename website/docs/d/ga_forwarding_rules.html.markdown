---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_forwarding_rules"
sidebar_current: "docs-alicloud-datasource-ga-forwarding-rules"
description: |-
  Provides a list of Global Accelerator (GA) Forwarding Rules to the user.
---

# alicloud_ga_forwarding_rules

This data source provides the Global Accelerator (GA) Forwarding Rules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_forwarding_rules" "example" {
  accelerator_id = "example_value"
  listener_id    = "example_value"
  ids            = ["example_value"]
}

output "first_ga_forwarding_rule_id" {
  value = data.alicloud_ga_forwarding_rules.example.forwarding_rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required) The ID of the Global Accelerator instance.
* `listener_id` - (Required) The ID of the listener.
* `ids` - (Optional, Computed)  A list of Forwarding Rule IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional)  The status of the acceleration region. Valid values: `active`, `configuring`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `forwarding_rules` - A list of Ga Forwarding Rules.

---

### `forwarding_rules`

The forwarding_rules supports the following attributes:

* `priority` -  Forwarding policy priority.
* `forwarding_rule_id` - Forwarding Policy ID.
* `forwarding_rule_name` - Forwarding policy name. The length of the name is 2-128 English or Chinese characters.
* `listener_id` - The ID of the listener.
* `rule_conditions` -  Forward action.
* `rule_actions` - The IP protocol used by the GA instance.
* `id` -  The resource ID in terraform of Forwarding Rule.
* `forwarding_rule_status` -  Forwarding Policy Status.

---

### `forwarding_rules.rule_conditions`

The rule_conditions supports the following attributes:

* `rule_condition_type` - Forwarding condition type.
* `path_config` - Path configuration information.
* `host_config` - Domain name configuration information.

---

### `forwarding_rules.rule_conditions.path_config`

The path_config supports the following attributes:

* `values` - The length of the path is 1-128 characters.

---

### `forwarding_rules.rule_conditions.host_config`

The host_config supports the following attributes:

* `values` - The domain name is 3-128 characters long.

---

### `forwarding_rules.rule_actions`

The rule_actions supports the following attributes:

* `order` - Forwarding priority.
* `rule_action_type` - Forward action type.
* `forward_group_config` - Forwarding configuration.

---

### `forwarding_rules.rule_actions.forward_group_config`

The forward_group_config supports the following attributes:

* `server_group_tuples` - Terminal node group configuration.

---

### `forwarding_rules.rule_actions.forward_group_config.server_group_tuples`

The server_group_tuples supports the following attributes:

* `endpoint_group_id` - Terminal node group ID.
