---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_diagnostic_report"
description: |-
  Provides a Alicloud ECS Diagnostic Report resource.
---

# alicloud_ecs_diagnostic_report

Provides a ECS Diagnostic Report resource.



For information about ECS Diagnostic Report and how to use it, see [What is Diagnostic Report](https://next.api.alibabacloud.com/document/Ecs/2014-05-26/CreateDiagnosticReport).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_diagnostic_report" "example" {
  resource_id   = "i-bp1io7xxxutup5l2y4vt"
  metric_set_id = "dms-instancedefault"
  additional_options = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:
* `additional_options` - (Optional, ForceNew, Map) json extension
* `end_time` - (Optional, ForceNew) End time. Only takes effect for diagnostic indicators that do not need to execute cloud assistant instructions in GuestOS.
* `metric_set_id` - (Optional, ForceNew) The ID of the diagnostic indicator set. If it is empty, the default diagnostic metric collection dms-instancedefault of the ECS instance is used.
* `resource_id` - (Required, ForceNew) The resource ID.
* `start_time` - (Optional, ForceNew) Start time. Only takes effect for diagnostic indicators that do not need to execute cloud assistant instructions in GuestOS.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Import

ECS Diagnostic Report can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_diagnostic_report.example <diagnostic_report_id>
```
