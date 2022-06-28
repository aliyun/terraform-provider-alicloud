---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_project"
sidebar_current: "docs-alicloud-resource-log-project"
description: |-
  Provides a Alicloud log project resource.
---

# alicloud\_log\_project

The project is the resource management unit in Log Service and is used to isolate and control resources.
You can manage all the logs and the related log sources of an application by using projects. [Refer to details](https://www.alibabacloud.com/help/doc-detail/48873.htm).

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
  tags = {"test":"test"}
}
```

## Module Support

You can use the existing [sls module](https://registry.terraform.io/modules/terraform-alicloud-modules/sls/alicloud) 
to create SLS project, store and store index one-click, like ECS instances.

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the log project. It is the only in one Alicloud account.
* `description` - (Optional) Description of the log project.
* `tags` - (Optional) Log project tags.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It sames as its name.
* `name` - Log project name.
* `description` - Log project description.
* `tags` - Log project tags.

### Timeouts

-> **NOTE:** Available in 1.126.0+

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the log project.

## Import

Log project can be imported using the id or name, e.g.

```
$ terraform import alicloud_log_project.example tf-log
```
