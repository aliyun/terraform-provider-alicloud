---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_quality_check_results"
sidebar_current: "docs-alicloud-datasource-data-works-quality-check-results"
description: |-
  Provides a list of Data Works Quality Check Results to the user.
---

# alicloud\_data\_works\_quality\_check\_results

This data source provides the Data Works Quality Check Results of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_quality_check_results" "example" {
  project_name = "your_project_name"
  start_date   = "2024-01-01 00:00:00"
  end_date     = "2024-01-02 00:00:00"
  entity_id    = "your_entity_id"
}

output "quality_check_result_id" {
  value = data.alicloud_data_works_quality_check_results.example.results.0.id
}
```

Query quality check results by rule id

```terraform
data "alicloud_data_works_quality_check_results" "example" {
  project_name = "your_project_name"
  start_date   = "2024-01-01 00:00:00"
  end_date     = "2024-01-02 00:00:00"
  rule_id      = "your_rule_id"
}
```

## Argument Reference

The following arguments are supported:

* `end_date` - (Required, ForceNew) The end of the business date range. The value must be in the `yyyy-MM-dd HH:mm:ss` format.
* `entity_id` - (Optional, ForceNew) The ID of the partition expression. It is required when querying results by entity. Conflict with `rule_id`.
* `ids` - (Optional, ForceNew, Computed) A list of Quality Check Result IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `page_number` - (Optional, ForceNew) The page number. Default to 1.
* `page_size` - (Optional, ForceNew) The number of entries per page. Default to 10. Maximum value is 20.
* `project_name` - (Required, ForceNew) The name of the engine or data source.
* `rule_id` - (Optional, ForceNew) The ID of the monitoring rule. It is required when querying results by rule. Conflict with `entity_id`.
* `start_date` - (Required, ForceNew) The start of the business date range. The value must be in the `yyyy-MM-dd HH:mm:ss` format.

-> **NOTE:** The time range to query cannot exceed 7 days. At least one of `entity_id` or `rule_id` should be specified.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `results` - A list of Data Works Quality Check Results. Each element contains the following attributes:

### `results`

* `id` - The ID of the Quality Check Result.
* `qualityt_check_result_id` - The primary key ID of the Quality Check Result.
* `entity_id` - The ID of the partition expression.
* `rule_id` - The ID of the monitoring rule.
* `project_name` - The name of the project.
* `table_name` - The name of the table.
* `op` - The operator.
* `task_id` - The ID of the task.
* `checker_name` - The name of the checker.
* `checker_id` - The ID of the checker.
* `checker_type` - The type of the checker.
* `begin_time` - The begin time.
* `end_time` - The end time.
* `expect_value` - The comparison value of the fixed value.
* `upper_value` - The historical maximum value.
* `lower_value` - The historical minimum value.
* `warning_threshold` - The orange warning threshold.
* `critical_threshold` - The red warning threshold.
* `trend` - The trend value.
* `check_result` - The quality check result value.
* `fixed_check` - Whether it is a fixed value check.
* `is_prediction` - Whether it is a prediction.
* `block_type` - Whether it is a strong rule.
* `match_expression` - The partition expression.
* `actual_expression` - The actual partition.
* `where_condition` - The where condition.
* `property` - The rule field.
* `method_name` - The collection method.
* `date_type` - The date type.
* `rule_name` - The name of the rule.
* `template_name` - The name of the template.
* `template_id` - The ID of the template.
* `comment` - The comment of the rule.
* `external_id` - The external node ID.
* `external_type` - The external trigger type.
