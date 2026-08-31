---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_report_template"
description: |-
  Provides a Alicloud Cloud Config (Config) Report Template resource.
---

# alicloud_config_report_template

Provides a Cloud Config (Config) Report Template resource.

Config Compliance Report Template.

For information about Cloud Config (Config) Report Template and how to use it, see [What is Report Template](https://next.api.alibabacloud.com/document/Config/2020-09-07/CreateReportTemplate).

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_report_template&exampleId=3611415d-f121-bf1f-975f-f3b97ec697015ef7b94c&activeTab=example&spm=docs.r.config_report_template.0.3611415df1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_config_report_template" "default" {
  report_granularity = "AllInOne"
  report_scope {
    key        = "RuleId"
    value      = "cr-xxx"
    match_type = "In"
  }
  report_file_formats         = "excel"
  report_template_name        = "example-name"
  report_template_description = "example-desc"
  subscription_frequency      = " "
  report_language             = "en-US"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_config_report_template&spm=docs.r.config_report_template.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `report_file_formats` - (Optional) Report Format
* `report_granularity` - (Optional) Report Aggregation Granularity
* `report_language` - (Optional) This property does not have a description in the spec, please add it before generating code.
* `report_scope` - (Optional, ForceNew, List) Report range, yes and logic between multiple sets of k-v pairs. See [`report_scope`](#report_scope) below.
* `report_template_description` - (Optional) Report Template Description
* `report_template_name` - (Required) Report Template Name
* `subscription_frequency` - (Optional) Report subscription frequency. If this field is not empty, it is a Cron expression in Quartz format triggered by the subscription notification.

The format is: Seconds, time, day, month, week. The following are examples of commonly used Cron expressions:
  - Execute at 0 o'clock every day: 0 0 0 * *?
  - Every Monday at 15: 30: 0 30 15? * MON
  - Execute at 2 o'clock on the 1st of each month: 0 0 2 1 *?

Among them:
  -"*" Indicates any value
  - What-? Used for day and week fields, indicating that no specific value is specified
  - MON means Monday

-> **NOTE:**  The trigger time is UTC +8, and the settings of the cron expression can be converted according to the time zone.

-> **NOTE:**  It can only be triggered according to the cron expression time as much as possible. The cron expression limits the same template to trigger at most one notification per day.


### `report_scope`

The report_scope supports the following:
* `key` - (Optional) Key for reporting scope, currently supported:
  - AggregatorId
  - CompliancePackId
  - RuleId
* `match_type` - (Optional) The matching logic. Currently, only In is supported.
* `value` - (Optional) The value of the report range. Each k-v pair is an OR logic. For example, multiple rule IDs can be separated by commas (,).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Report Template.
* `delete` - (Defaults to 5 mins) Used when delete the Report Template.
* `update` - (Defaults to 5 mins) Used when update the Report Template.

## Import

Cloud Config (Config) Report Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_report_template.example <id>
```