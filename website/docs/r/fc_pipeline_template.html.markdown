---
subcategory: "FC"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_pipeline_template"
description: |-
  Provides a FC Pipeline Template resource.
---

# alicloud_fc_pipeline_template

Provides a FC Pipeline Template resource.

A pipeline template is a reusable definition of a pipeline's execution graph
(context data and ordered tasks). Templates are managed by the FC serverless
deployment API and referenced by pipelines either by name or by embedding the
template spec.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_fc_pipeline_template" "default" {
  name        = var.name
  description = "terraform-example pipeline template"
  spec        = jsonencode({
    context = {
      data = {
        region = "cn-hangzhou"
      }
    }
  })
  labels = {
    team = "devops"
  }
}
```

For example, the `tasks` array of `spec` carries an ordered execution graph
where each task may declare a `name`, a `taskTemplate` it references, a
task-local `context`, and a `runAfters` list of `{ name = "..." }` entries
naming tasks that must finish first. See the [PipelineTemplate
spec](https://api.aliyun.com/) reference for the full schema.

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew, String) The name of the pipeline template. It must be unique within the account and cannot be changed after creation. 1 to 40 characters; must start with a lowercase letter, contain only lowercase letters, digits and hyphens, and not end with a hyphen.
* `description` - (Optional, String) The description of the pipeline template.
* `spec` - (Optional, String) The pipeline template spec, encoded as a JSON object. Use `jsonencode()` in HCL. The object has the following top-level fields:
  - `context` - (Optional) An object with a single `data` field, an arbitrary JSON object describing the execution context shared by all tasks.
  - `tasks` - (Optional) An array of task execution objects. Each task has:
    - `name` - (Required) The task name, unique within the template.
    - `taskTemplate` - (Optional) The name of the task template the task references.
    - `context` - (Optional) A task-local context object with the same shape as the template `context`.
    - `runAfters` - (Optional) An array of `{ name = "..." }` objects naming tasks that must complete before this task starts.
  - `contextSchema` - (Read-only) The schema of the context, managed by the server.
* `labels` - (Optional, Map of String) The labels of the pipeline template.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, same as `name`.
* `uid` - (Computed, String) The UID of the account that owns the template.
* `kind` - (Computed, String) The kind of the resource.
* `created_time` - (Computed, String) The creation time of the template.
* `deletion_time` - (Computed, String) The deletion time of the template.
* `generation` - (Computed, Int) The generation of the spec, incremented when the spec changes.
* `resource_version` - (Computed, Int) The resource version, incremented on every update.

## Import

FC Pipeline Template can be imported using the id, e.g. the template `name`:

```shell
$ terraform import alicloud_fc_pipeline_template.example <name>
```
