---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_triggers"
sidebar_current: "docs-alicloud-datasource-fc-triggers"
description: |-
    Provides a list of FC triggers to the user.
---

# alicloud\_fc_triggers

This data source provides the Function Compute triggers of the current Alibaba Cloud user.

## Example Usage

```
data "alicloud_fc_triggers" "fc_triggers_ds" {
  service_name  = "sample_service"
  function_name = "sample_function"
  name_regex    = "sample_fc_trigger"
}

output "first_fc_trigger_name" {
  value = "${data.alicloud_fc_triggers.fc_triggers_ds.triggers.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - FC service name.
* `function_name` - FC function name.
* `name_regex` - (Optional) A regex string to filter results by FC trigger name.
* `ids` (Optional, Available in 1.53.0+) - A list of FC triggers ids.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of FC triggers ids.
* `names` - A list of FC triggers names.
* `triggers` - A list of FC triggers. Each element contains the following attributes:
  * `id` - FC trigger ID.
  * `name` - FC trigger name.
  * `source_arn` - Event source resource address. See [Create a trigger](https://www.alibabacloud.com/help/doc-detail/53102.htm) for more details.
  * `type` - Type of the trigger. Valid values: `oss`, `log`, `timer`, `http`, `mns_topic`, `cdn_events` and `eventbridge`.
  * `invocation_role` - RAM role arn attached to the Function Compute trigger. Role used by the event source to call the function. The value format is "acs:ram::$account-id:role/$role-name". See [Create a trigger](https://www.alibabacloud.com/help/doc-detail/53102.htm) for more details.
  * `config` - JSON-encoded trigger configuration. See [Configure triggers and events](https://www.alibabacloud.com/help/doc-detail/70140.htm) for more details.
  * `creation_time` - FC trigger creation time.
  * `last_modification_time` - FC trigger last modification time.
