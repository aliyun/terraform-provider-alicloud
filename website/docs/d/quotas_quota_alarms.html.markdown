---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_quota_alarms"
sidebar_current: "docs-alicloud-datasource-quotas-quota-alarms"
description: |-
  Provides a list of Quotas Quota Alarms to the user.
---

# alicloud\_quotas\_quota\_alarms

This data source provides the Quotas Quota Alarms of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_quotas_quota_alarms" "example" {
  ids        = ["5VR90-421F886-81E9-xxx"]
  name_regex = "tf-testAcc"
}

output "first_quotas_quota_alarm_id" {
  value = data.alicloud_quotas_quota_alarms.example.alarms.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Quota Alarm IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Quota Alarm name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `product_code` - (Optional, ForceNew) The Product Code.
* `quota_action_code` - (Optional, ForceNew) The Quota Action Code.
* `quota_alarm_name` - (Optional, ForceNew) The name of Quota Alarm.
* `quota_dimensions` - (Optional, ForceNew) The Quota Dimensions.

#### Block quota_dimensions

The dimensions supports the following: 

* `key` - (Optional, ForceNew) The key of quota_dimensions.
* `value` - (Optional, ForceNew) The value of quota_dimensions.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Quota Alarm names.
* `alarms` - A list of Quotas Quota Alarms. Each element contains the following attributes:
	* `alarm_id` - The first ID of the resource.
	* `id` - The ID of the Quota Alarm.
	* `product_code` - The Product Code.
	* `quota_action_code` - The Quota Action Code.
	* `quota_alarm_name` - The name of Quota Alarm.
	* `quota_dimensions` - The Quota Dimensions.
		* `key` - The key of quota_dimensions.
		* `value` - The value of quota_dimensions.
	* `threshold` - The threshold of Quota Alarm.
	* `threshold_percent` - The threshold percent of Quota Alarm.
	* `web_hook` - The WebHook of Quota Alarm.
