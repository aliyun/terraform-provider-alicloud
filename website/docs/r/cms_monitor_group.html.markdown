---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_monitor_group"
sidebar_current: "docs-alicloud-resource-cms-monitor-group"
description: |-
  Provides a Alicloud Cloud Monitor Service Monitor Group resource.
---

# alicloud_cms_monitor_group

Provides a Cloud Monitor Service Monitor Group resource.

For information about Cloud Monitor Service Monitor Group and how to use it, see [What is Monitor Group](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createmonitorgroup).

-> **NOTE:** Available since v1.113.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_monitor_group" "example" {
  monitor_group_name = "tf-example-accmonitorgroup"
}

resource "alicloud_cms_monitor_group" "default2" {
  contact_groups      = ["your_contact_groups"]
  resource_group_id   = "your_resource_group_id"
  resource_group_name = "resource_group_name"
  tags = {
    Created = "TF"
    For     = "Acceptance-example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `contact_groups` - (Optional) The alert group to which alert notifications will be sent.
* `monitor_group_name` - (Optional) The name of the application group.
* `resource_group_id` - (Optional, Available since v1.141.0) The ID of the resource group.
* `resource_group_name` - (Optional, Available since v1.141.0) The name of the resource group.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Monitor Group.

## Import

Cloud Monitor Service Monitor Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_monitor_group.example <id>
```
