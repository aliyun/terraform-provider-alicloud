---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_machine_group"
description: |-
  Provides a Alicloud Log Service (SLS) Machine Group resource.
---

# alicloud_sls_machine_group

Provides a Log Service (SLS) Machine Group resource.



For information about Log Service (SLS) Machine Group and how to use it, see [What is Machine Group](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateMachineGroup).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-nanjing"
}

variable "project_name" {
  default = "project-for-machine-group-terraform"
}

resource "alicloud_log_project" "defaultyJqrue" {
  description = "for terraform example"
  name        = var.project_name
}


resource "alicloud_sls_machine_group" "default" {
  group_name            = "group1"
  project_name          = var.project_name
  machine_identify_type = "ip"
  group_attribute {
    group_topic   = "example"
    external_name = "example"
  }
  machine_list = ["192.168.1.1"]
}
```

## Argument Reference

The following arguments are supported:
* `group_attribute` - (Optional, ForceNew, List) Properties of machine groups. For details, please refer to the groupAttribute parameter description in the following table. See [`group_attribute`](#group_attribute) below.
* `group_name` - (Required, ForceNew) Machine Group name
* `group_type` - (Optional, ForceNew) Machine group type, optional value is empty.
* `machine_identify_type` - (Required, ForceNew) Machine identification type.
  - ip: ip address Machine Group.
  - userdefined: user-defined identity Machine Group.
* `machine_list` - (Required, ForceNew, List) The identification information of the machine group.
  - If machineidentifiytype is configured to ip, enter the ip address of the server.
  - If machineidentifiytype is configured to userdefined, enter a custom identifier here.
* `project_name` - (Required, ForceNew) Project name

### `group_attribute`

The group_attribute supports the following:
* `external_name` - (Optional, ForceNew) The external management system identification on which the machine group depends.
* `group_topic` - (Optional, ForceNew) The log topic of the machine group.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_name>:<group_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Machine Group.
* `delete` - (Defaults to 5 mins) Used when delete the Machine Group.

## Import

Log Service (SLS) Machine Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_machine_group.example <project_name>:<group_name>
```