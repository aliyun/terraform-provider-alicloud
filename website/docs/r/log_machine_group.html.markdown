---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_machine_group"
sidebar_current: "docs-alicloud-resource-log-machine-group"
description: |-
  Provides a Alicloud log tail machine group resource.
---

# alicloud\_log\_machine\_group

Log Service manages all the ECS instances whose logs need to be collected by using the Logtail client in the form of machine groups.
 [Refer to details](https://www.alibabacloud.com/help/doc-detail/28966.htm)

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}
resource "alicloud_log_machine_group" "example" {
  project       = "${alicloud_log_project.example.name}"
  name          = "tf-machine-group"
  identify_type = "ip"
  topic         = "terraform"
  identify_list = ["10.0.0.1", "10.0.0.2"]
}
```
## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the machine group belongs.
* `name` - (Required, ForceNew) The machine group name, which is unique in the same project.
* `identify_type` - (Optional) The machine identification type, including IP and user-defined identity. Valid values are "ip" and "userdefined". Default to "ip".
* `identify_list`- (Required) The specific machine identification, which can be an IP address or user-defined identity.
* `topic` - (Optional) The topic of a machine group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log machine group. It formats of `<project>:<name>`.
* `project` - The project name.
* `name` - The machine group name.
* `identify_type` - The machine identification type.
* `identify_list` - The machine identification.
* `topic` - The machine group topic.

## Import

Log machine group can be imported using the id, e.g.

```
$ terraform import alicloud_log_machine_group.example tf-log:tf-machine-group
```
