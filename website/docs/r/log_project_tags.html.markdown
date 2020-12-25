---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_project_tags"
sidebar_current: "docs-alicloud-resource-log-project-tags"
description: |-
Provides a Alicloud log project tags resource.
---

# alicloud\_log\_project\_tags
Project tags is a resource used to tag a project. You can use this function to tag a project.
[Refer to details](https://www.alibabacloud.com/help/doc-detail/48873.htm).

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "alicloud_log_project_tags" "default" {
    project_name = "tf-log-tags"
    tags = {"name1":"aliyun"}
}

```


## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The project name. It is the only in one Alicloud account.
* `tags` - (Required, ForceNew) Label of project.

## Attributes Reference

The following attributes are exported:

* `id` - The id of project tags.


## Import

Log project can be imported using the id or name, e.g.

```
$ terraform import alicloud_log_project.example tf-log:tf-log-tags
```
