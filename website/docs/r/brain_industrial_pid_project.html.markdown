---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_pid_project"
sidebar_current: "docs-alicloud-resource-brain-industrial-pid-project"
description: |-
  Provides a Alicloud Brain Industrial Pid Project resource.
---

# alicloud\_brain\_industrial\_pid\_project

Provides a Brain Industrial Pid Project resource.

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_brain_industrial_pid_project" "example" {
  pid_organization_id = "3e74e684-cbb5-xxxx"
  pid_project_name    = "tf-testAcc"
}

```

## Argument Reference

The following arguments are supported:

* `pid_organization_id` - (Required) The ID of Pid Organization.
* `pid_project_desc` - (Optional) The description of Pid Project.
* `pid_project_name` - (Required) The name of Pid Project.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Pid Project.

## Import

Brain Industrial Pid Project can be imported using the id, e.g.

```
$ terraform import alicloud_brain_industrial_pid_project.example <id>
```
