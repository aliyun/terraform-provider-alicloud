---
subcategory: "Serverless Workflow (FnF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fnf_schedule"
sidebar_current: "docs-alicloud-resource-fnf-schedule"
description: |-
  Provides a Alicloud Serverless Workflow Schedule resource.
---

# alicloud_fnf_schedule

Provides a Serverless Workflow Schedule resource.

For information about Serverless Workflow Schedule and how to use it, see [What is Schedule](https://www.alibabacloud.com/help/en/doc-detail/168934.htm).

-> **NOTE:** Available since v1.105.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fnf_schedule&exampleId=d3753d0e-4558-f0b5-2acb-5d75c171d1be1a93fe3d&activeTab=example&spm=docs.r.fnf_schedule.0.d3753d0e45&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_fnf_flow" "example" {
  definition  = <<EOF
  version: v1beta1
  type: flow
  steps:
    - type: pass
      name: helloworld
  EOF  
  description = "tf-exampleFnFFlow983041"
  name        = "tf-exampleSchedule"
  type        = "FDL"
}

resource "alicloud_fnf_schedule" "example" {
  cron_expression = "30 9 * * * *"
  description     = "tf-exampleFnFSchedule983041"
  enable          = "true"
  flow_name       = alicloud_fnf_flow.example.name
  payload         = "{\"tf-example\": \"example success\"}"
  schedule_name   = "tf-exampleFnFSchedule983041"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_fnf_schedule&spm=docs.r.fnf_schedule.example&intl_lang=EN_US)

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Schedule.
* `delete` - (Defaults to 1 mins) Used when delete the Schedule.
* `update` - (Defaults to 1 mins) Used when update the Schedule.

## Import

Serverless Workflow Schedule can be imported using the id, e.g.

```shell
$ terraform import alicloud_fnf_schedule.example <schedule_name>:<flow_name>
```
