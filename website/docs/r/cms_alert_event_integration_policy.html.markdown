---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_event_integration_policy"
description: |-
  Provides a Alicloud Cms Alert Event Integration Policy resource.
---

# alicloud_cms_alert_event_integration_policy

Provides a Cms Alert Event Integration Policy resource.

For information about Cms Alert Event Integration Policy and how to use it, see [What is Alert Event Integration Policy](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateAlertEventIntegrationPolicy).

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}"
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = "${var.name}"
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_event_integration_policy" "default" {
  alert_event_integration_policy_name = "${var.name}"
  workspace                           = alicloud_cms_workspace.default.id
  type                                = "SYS_EVENT"
  description                         = "terraform-example"
  integration_setting                 = "{\"key\":\"value\"}"
  enable                              = true

  filter_setting {
    conditions {
      field = "name"
      value = "test"
      op    = "eq"
    }
    expression = "test-expr"
    relation   = "and"
  }

  transformer_setting {
    source    = "name"
    target    = "name"
    type      = "replace"
    value     = "test"
    variable  = "name"
    label_key = "key"
    reg_exp   = ".*"
    mapping = {
      key = "value"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `alert_event_integration_policy_name` - (Required) The name of the alert event integration policy.
* `workspace` - (Required, ForceNew) The workspace ID to which the alert event integration policy belongs.
* `description` - (Optional) The description of the alert event integration policy.
* `type` - (Optional) The type of the alert event integration policy. Valid values: `SYS_EVENT`, `CMS_ALERT`, `K8S_EVENT`, `CUSTOM`, `PROMETHEUS`, `GRAFANA`, `ZABBIX`, `SKYWALKING`, `OPEN_FALCON`, `NAGIOS`.
* `integration_setting` - (Optional) The integration setting of the alert event integration policy, in JSON format.
* `enable` - (Optional) Whether to enable the alert event integration policy. Default to `true`.
* `filter_setting` - (Optional) The filter setting of the alert event integration policy. See [`filter_setting`](#filter_setting) below.
* `transformer_setting` - (Optional) The transformer setting of the alert event integration policy. See [`transformer_setting`](#transformer_setting) below.

### `filter_setting`

The filter_setting supports the following:

* `conditions` - (Optional) The list of filter conditions. See [`conditions`](#filter_setting-conditions) below.
* `expression` - (Optional) The filter expression.
* `relation` - (Optional) The relation between filter conditions.

### `filter_setting-conditions`

The filter_setting-conditions supports the following:

* `field` - (Optional) The field name of the filter condition.
* `value` - (Optional) The field value of the filter condition.
* `op` - (Optional) The comparison operator of the filter condition.

### `transformer_setting`

The transformer_setting supports the following:

* `filter_setting` - (Optional) The filter setting of the transformer. See [`filter_setting`](#transformer_setting-filter_setting) below.
* `label_key` - (Optional) The label key of the transformer.
* `mapping` - (Optional) The mapping table of the transformer.
* `reg_exp` - (Optional) The text extract regular expression of the transformer.
* `source` - (Optional) The source path of the transformer.
* `target` - (Optional) The transform target of the transformer.
* `type` - (Optional) The action type of the transformer.
* `value` - (Optional) The value of the transformer.
* `variable` - (Optional) The variable name of the transformer.

### `transformer_setting-filter_setting`

The transformer_setting-filter_setting supports the following:

* `conditions` - (Optional) The list of filter conditions. See [`conditions`](#transformer_setting-filter_setting-conditions) below.
* `expression` - (Optional) The filter expression.
* `relation` - (Optional) The relation between filter conditions.

### `transformer_setting-filter_setting-conditions`

The transformer_setting-filter_setting-conditions supports the following:

* `field` - (Optional) The field name of the filter condition.
* `value` - (Optional) The field value of the filter condition.
* `op` - (Optional) The comparison operator of the filter condition.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alert event integration policy.
* `region_id` - The region ID of the resource.
* `create_time` - The creation time of the resource.
* `update_time` - The update time of the resource.
* `user_id` - The user ID of the resource.

## Import

Cms Alert Event Integration Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_alert_event_integration_policy.example <id>
```
