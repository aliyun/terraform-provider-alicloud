---
subcategory: "Serverless Workflow"
layout: "alicloud"
page_title: "Alicloud: alicloud_fnf_schedule"
sidebar_current: "docs-alicloud-resource-fnf-schedule"
description: |-
  Provides a Alicloud Serverless Workflow Schedule resource.
---

# alicloud\_fnf\_schedule

Provides a Serverless Workflow Schedule resource.

For information about Serverless Workflow Schedule and how to use it, see [What is Schedule](https://www.alibabacloud.com/help/en/doc-detail/168934.htm).

-> **NOTE:** Available in v1.105.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_fnf_flow" "example" {
  definition  = <<EOF
  version: v1beta1
  type: flow
  steps:
    - type: pass
      name: helloworld
  EOF  
  description = "tf-testaccFnFFlow983041"
  name        = "tf-testAccSchedule"
  type        = "FDL"
}

resource "alicloud_fnf_schedule" "example" {
  cron_expression = "30 9 * * * *"
  description     = "tf-testaccFnFSchedule983041"
  enable          = "true"
  flow_name       = alicloud_fnf_flow.example.name
  payload         = "{\"tf-test\": \"test success\"}"
  schedule_name   = "tf-testaccFnFSchedule983041"
}
```

## Argument Reference

The following arguments are supported:

* `cron_expression` - (Required) The CRON expression of the time-based schedule to be created.
* `description` - (Optional) The description of the time-based schedule to be created.
* `enable` - (Optional) Specifies whether to enable the time-based schedule you want to create. Valid values: `false`, `true`.
* `flow_name` - (Required, ForceNew) The name of the flow bound to the time-based schedule you want to create.
* `payload` - (Optional) The trigger message of the time-based schedule to be created. It must be in JSON object format.
* `schedule_name` - (Required, ForceNew) The name of the time-based schedule to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Schedule. The value is formatted `<schedule_name>:<flow_name>`.
* `last_modified_time` - The time when the time-based schedule was last updated.
* `schedule_id` - The ID of the time-based schedule.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Schedule.
* `delete` - (Defaults to 1 mins) Used when delete the Schedule.
* `update` - (Defaults to 1 mins) Used when update the Schedule.

## Import

Serverless Workflow Schedule can be imported using the id, e.g.

```
$ terraform import alicloud_fnf_schedule.example <schedule_name>:<flow_name>
```
