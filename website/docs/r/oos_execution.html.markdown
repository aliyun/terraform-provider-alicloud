---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_execution"
sidebar_current: "docs-alicloud-resource-oos-execution"
description: |-
  Provides a OOS Execution resource.
---

# alicloud_oos_execution

Provides a OOS Execution resource. For information about Alicloud OOS Execution and how to use it, see [What is Resource Alicloud OOS Execution](https://www.alibabacloud.com/help/doc-detail/120771.htm).

-> **NOTE:** Available since v1.93.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oos_execution&exampleId=b1fee8b3-109b-a25f-6099-35fc9e29f4e9f073eb47&activeTab=example&spm=docs.r.oos_execution.0.b1fee8b310&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oos_template" "default" {
  content       = <<EOF
  {
    "FormatVersion": "OOS-2019-06-01",
    "Description": "Update Describe instances of given status",
    "Parameters":{
      "Status":{
        "Type": "String",
        "Description": "(Required) The status of the Ecs instance."
      }
    },
    "Tasks": [
      {
        "Properties" :{
          "Parameters":{
            "Status": "{{ Status }}"
          },
          "API": "DescribeInstances",
          "Service": "Ecs"
        },
        "Name": "foo",
        "Action": "ACS::ExecuteApi"
      }]
  }
  EOF
  template_name = "tf-example-name-${random_integer.default.result}"
  version_name  = "example"
  tags = {
    "Created" = "TF",
    "For"     = "acceptance Test"
  }
}

resource "alicloud_oos_execution" "example" {
  template_name = alicloud_oos_template.default.template_name
  description   = "From TF Test"
  parameters    = <<EOF
				{"Status":"Running"}
		  	EOF
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of OOS Execution.
* `loop_mode` - (Optional, ForceNew) The loop mode of OOS Execution.
* `mode` - (Optional, ForceNew) The mode of OOS Execution. Valid: `Automatic`, `Debug`. Default to `Automatic`.
* `parameters` - (Optional, ForceNew) The parameters required by the template. Default to `{}`.
* `parent_execution_id` - (Optional, ForceNew) The id of parent execution.
* `safety_check` - (Optional, ForceNew) The mode of safety check.
* `template_name` - (Required, ForceNew) The name of execution template.
* `template_version` - (Optional, ForceNew) The version of execution template.
* `template_content` - (Optional, ForceNew, Available in v1.114.0+) The content of template. When the user selects an existing template to create and execute a task, it is not necessary to pass in this field.
                    
## Attributes Reference

The following attributes are exported:

* `id` - The id of OOS Execution.
* `counters` - The counters of OOS Execution.
* `create_date` - The time when the execution was created.
* `end_date` - The time when the execution was ended.
* `executed_by` - The user who execute the template.
* `is_parent` - Whether to include subtasks.
* `outputs` - The outputs of OOS Execution.
* `ram_role` - The role that executes the current template.
* `start_date` - The time when the execution was started.
* `status` - The status of OOS Execution.
* `status_message` - The message of status.
* `template_id` - The id of template.
* `update_date` - The time when the execution was updated.

## Timeouts


The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when creating the alicloud_oos_execution.

## Import

OOS Execution can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_execution.example exec-ef6xxxx
```
