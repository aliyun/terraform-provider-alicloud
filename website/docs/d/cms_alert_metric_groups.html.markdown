---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_metric_groups"
sidebar_current: "docs-alicloud-datasource-cms-alert-metric-groups"
description: |-
  Provides a list of Cms Alert Metric Group available to the user.
---

# alicloud_cms_alert_metric_groups

This data source provides the predefined alert metric groups of CloudMonitor 2.0. Alert metric groups are read-only, system-defined entries that group alert metrics and carry shared filters and parameters; they cannot be created, modified or deleted by the user.

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_cms_alert_metric_groups" "default" {
  include_details = true
}

output "alicloud_cms_alert_metric_group_example_id" {
  value = data.alicloud_cms_alert_metric_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Alert Metric Group IDs.
* `datasource_type` - (Optional, ForceNew) The data source type, used to filter alert metric groups that support a certain data source type. Example: `arms_metrics`.
* `include_details` - (Optional, ForceNew, Bool) Whether to include detail attributes such as `filters` and `params`. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Alert Metric Group IDs.
* `groups` - A list of Alert Metric Group Entries.

### `groups`

The groups supports the following attributes:
* `id` - The ID of the alert metric group.
* `alert_metric_group_id` - The unique ID of the alert metric group. The `.` separator can be used to identify the hierarchy structure.
* `datasource_types` - The data source types supported by the alert metric group, in JSON array string format. Example: `["arms_metrics"]`.
* `display_name_cn` - The display name of the alert metric group in Chinese.
* `display_name_en` - The display name of the alert metric group in English.
* `description_cn` - The description of the alert metric group in Chinese.
* `description_en` - The description of the alert metric group in English.
* `order_index` - The ordering index of the alert metric group. A larger value indicates a higher priority within the same level.
* `filters` - The shared filter conditions of the metrics in the group. Metrics inherit these filters automatically instead of defining them repeatedly at the metric level. **NOTE:** This field is only available when `include_details` is `true`.
* `params` - The shared parameters of the metrics in the group. All metrics inherit them automatically. Only read-only text parameters can be defined at the group level. **NOTE:** This field is only available when `include_details` is `true`.

### `groups.filters`

The filters supports the following attributes:
* `dim` - The filter dimension.
* `display_name_cn` - The display name of the filter in Chinese.
* `display_name_en` - The display name of the filter in English.
* `hidden` - Indicates whether the filter is hidden. A hidden filter is not displayed in the frontend interaction, but its value can still be uploaded when rendering PromQL.
* `opt` - The filter operator.
* `label_disabled` - Indicates whether the filter is excluded from the label filter of PromQL.
* `dim_disabled` - Indicates whether the filter is excluded from the group by clause of PromQL.
* `supported_opts` - The supported operators of the filter.

### `groups.filters.supported_opts`

The supported_opts supports the following attributes:
* `display_name_cn` - The display name of the operator in Chinese.
* `display_name_en` - The display name of the operator in English.
* `value` - The value of the operator.

### `groups.params`

The params supports the following attributes:
* `max_width` - The maximum width of the input box. Only valid for `SELECT_PARAM` and `INPUT_PARAM`.
* `min_width` - The minimum width of the input box. Only valid for `SELECT_PARAM` and `INPUT_PARAM`.
* `name` - The name of the parameter.
* `placeholder_cn` - The Chinese placeholder displayed in the frontend. Only valid for `INPUT_PARAM`.
* `placeholder_en` - The English placeholder displayed in the frontend. Only valid for `INPUT_PARAM`.
* `type` - The type of the parameter. Valid values: `TEXT_PARAM` (read-only text parameter defined by the backend, no user input control is displayed), `INPUT_PARAM` (input box parameter) and `SELECT_PARAM` (select box parameter).
* `value` - The value of the parameter.
* `values` - The optional value list of the dropdown. Only valid for `SELECT_PARAM`.

### `groups.params.values`

The values supports the following attributes:
* `label_cn` - The Chinese display name of the option.
* `label_en` - The English display name of the option.
* `value` - The value of the option.
