---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_application_group"
sidebar_current: "docs-alicloud-resource-oos-application-group"
description: |-
  Provides a Alicloud OOS Application Group resource.
---

# alicloud\_oos\_application\_group

Provides a OOS Application Group resource.

For information about OOS Application Group and how to use it, see [What is Application Group](https://www.alibabacloud.com/help/en/doc-detail/120556.html).

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = "example_value"
  description       = "example_value"
  tags = {
    Created = "TF"
  }
}

resource "alicloud_oos_application_group" "default" {
  application_group_name = var.name
  application_name       = alicloud_oos_application.default.id
  deploy_region_id       = "example_value"
  description            = "example_value"
  import_tag_key         = "example_value"
  import_tag_value       = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `application_group_name` - (Required, ForceNew) The name of the Application group.
* `application_name` - (Required, ForceNew) The name of the Application.
* `deploy_region_id` - (Required, ForceNew) The region ID of the deployment.
* `description` - (Optional, ForceNew) Application group description information.
* `import_tag_key` - (Optional, ForceNew) The tag key must be passed in at the same time as the tag value (import_tag_value) or none, not just one. If both `import_tag_key` and `import_tag_value` are left empty, the default is app-{ApplicationName} (application name).
* `import_tag_value` - (Optional, ForceNew) The tag value must be passed in at the same time as the tag key (import_tag_key) or none, not just one. If both `import_tag_key` and `import_tag_value` are left empty, the default is application group name.
.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Application Group. The value formats as `<application_name>:<application_group_name>`.

## Import

OOS Application Group can be imported using the id, e.g.

```
$ terraform import alicloud_oos_application_group.example <application_name>:<application_group_name>
```