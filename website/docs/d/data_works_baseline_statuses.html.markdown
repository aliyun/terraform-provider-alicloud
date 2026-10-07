---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_baseline_statuses"
sidebar_current: "docs-alicloud-datasource-data-works-baseline-statuses"
description: |-
  Provides a list of Data Works Baseline Statuses to the user.
---

# alicloud_data_works_baseline_statuses

This data source provides the Data Works Baseline Statuses available to the current Alibaba Cloud user.

-> **NOTE:** Available since v1.232.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_baseline_statuses" "example" {
  bizdate       = "2026-09-16T00:00:00Z"
  status        = "ERROR,SAFE"
  finish_status = "UNFINISH"
}

output "first_baseline_id" {
  value = data.alicloud_data_works_baseline_statuses.example.statuses.0.baseline_id
}
```

Filter by IDs

```terraform
data "alicloud_data_works_baseline_statuses" "ids" {
  bizdate = "2026-09-16T00:00:00Z"
  ids     = ["123456"]
}
```

## Argument Reference

The following arguments are supported:

* `bizdate` - (Required, ForceNew) The business date in UTC format. Format: `yyyy-MM-dd'T'HH:mm:ssZ`.
* `baseline_types` - (Optional, ForceNew) The baseline type. Valid values: `DAILY`, `HOURLY`. Multiple values are separated by commas.
* `owner` - (Optional, ForceNew) The Alibaba Cloud UID of the baseline owner.
* `priority` - (Optional, ForceNew) The baseline priority. Valid values: `1`, `3`, `5`, `7`, `8`. Multiple values are separated by commas.
* `status` - (Optional, ForceNew) The baseline status. Valid values: `ERROR`, `SAFE`, `DANGEROUS`, `OVER`. Multiple values are separated by commas.
* `finish_status` - (Optional, ForceNew) The completion status. Valid values: `UNFINISH`, `FINISH`. Multiple values are separated by commas.
* `search_text` - (Optional, ForceNew) The search keyword by baseline name or ID.
* `topic_id` - (Optional, ForceNew) The ID of the related event.
* `ids` - (Optional, ForceNew, Computed) A list of Baseline IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Baseline IDs.
* `statuses` - A list of Data Works Baseline Statuses. Each element contains the following attributes:

### `statuses`

The following attributes are exported for each baseline status:

* `id` - The Baseline ID.
* `baseline_id` - The Baseline ID.
* `baseline_name` - The name of the baseline.
* `baseline_type` - The type of the baseline. Valid values: `DAILY`, `HOURLY`.
* `buffer` - The buffer time in seconds.
* `status` - The status of the baseline. Valid values: `ERROR`, `SAFE`, `DANGEROUS`, `OVER`.
* `owner` - The UID of the baseline owner. Multiple UIDs are separated by commas.
* `priority` - The priority of the baseline. Valid values: `1`, `3`, `5`, `7`, `8`.
* `finish_status` - The completion status. Valid values: `UNFINISH`, `FINISH`.
* `finish_time` - The completion timestamp. Returned only when `finish_status` is `FINISH`.
* `project_id` - The workspace ID.
* `bizdate` - The business date timestamp.
* `exp_time` - The warning time.
* `in_group_id` - The cycle number. `1` for daily baselines and `1` to `24` for hourly baselines.
* `sla_time` - The actual completion time.
* `end_cast` - The estimated completion time.
* `region_id` - The region ID of the baseline.
