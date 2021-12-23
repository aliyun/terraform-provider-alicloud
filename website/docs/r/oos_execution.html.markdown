---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_execution"
sidebar_current: "docs-alicloud-resource-oos-execution"
description: |-
  Provides a OOS Execution resource.
---

# alicloud\_oos\_execution

Provides a OOS Execution resource. For information about Alicloud OOS Execution and how to use it, see [What is Resource Alicloud OOS Execution](https://www.alibabacloud.com/help/doc-detail/120771.htm).

-> **NOTE:** Available in 1.93.0+.

## Example Usage

```terraform
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
  template_name = "test-name"
  version_name  = "test"
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

### Timeouts


The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when creating the alicloud_oos_execution (until it reaches the initial `Running` status).

## Import

OOS Execution can be imported using the id, e.g.

```
$ terraform import alicloud_oos_execution.example exec-ef6xxxx
```
