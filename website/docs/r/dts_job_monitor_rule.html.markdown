---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_job_monitor_rule"
sidebar_current: "docs-alicloud-resource-dts-job-monitor-rule"
description: |-
  Provides a Alicloud DTS Job Monitor Rule resource.
---

# alicloud\_dts\_job\_monitor\_rule

Provides a DTS Job Monitor Rule resource.

For information about DTS Job Monitor Rule and how to use it, see [What is Job Monitor Rule](https://www.aliyun.com/product/dts).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dts_job_monitor_rule" "example" {
  dts_job_id = "example_value"
  type       = "delay"
}
```

## Argument Reference

The following arguments are supported:

* `dts_job_id` - (Required, ForceNew) Migration, synchronization or subscription task ID can be by calling the [DescribeDtsJobs] get.
* `type` - (Required, ForceNew)  Monitoring rules of type, valid values: `delay`, `error`. **delay**: delay alarm. **error**: abnormal alarm.
* `delay_rule_time` - (Optional, Computed) Trigger delay alarm threshold, which is measured in seconds.
* `phone` - (Optional, Computed) The alarm is triggered after notification of the contact phone number, A plurality of phone numbers between them with a comma (,) to separate.
* `state` - (Optional, Computed) Whether to enable monitoring rules, valid values: `Y`, `N`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Job Monitor Rule. Its value is same as `dts_job_id`.

## Import

DTS Job Monitor Rule can be imported using the id, e.g.

```
$ terraform import alicloud_dts_job_monitor_rule.example <dts_job_id>
```
