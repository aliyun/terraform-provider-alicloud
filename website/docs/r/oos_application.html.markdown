---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_application"
sidebar_current: "docs-alicloud-resource-oos-application"
description: |-
  Provides a Alicloud OOS Application resource.
---

# alicloud\_oos\_application

Provides a OOS Application resource.

For information about OOS Application and how to use it, see [What is Application](https://www.alibabacloud.com/help/en/doc-detail/120556.html).

-> **NOTE:** Available in v1.145.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tftestacc"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = var.name
  description       = var.name
  tags = {
    Created = "TF"
  }
}

```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required) The name of the application.
* `description` - (Optional) Application group description information.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional) The tag of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Application. The value is formate as <application_name>.

## Import

OOS Application can be imported using the id, e.g.

```
$ terraform import alicloud_oos_application.example <id>
```