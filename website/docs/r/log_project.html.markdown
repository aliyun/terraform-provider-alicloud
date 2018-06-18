---
layout: "alicloud"
page_title: "Alicloud: alicloud_log_project"
sidebar_current: "docs-alicloud-resource-log-project"
description: |-
  Provides a Alicloud log project resource.
---

# alicloud\_log\_project

The project is the resource management unit in Log Service and is used to isolate and control resources. You can manage all the logs and the related log sources of an application by using projects.

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name       = "tf-log"
  description = "created by terraform"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the log project. It is the only in one Alicloud account.
* `description` - (ForceNew) Description of the log project. At present, it is not modified by terraform.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It sames as its name.
* `name` - Log project name.
* `description` - Log project description.

## Import

Log project can be imported using the id or name, e.g.

```
$ terraform import alicloud_log_project.example tf-log
```
