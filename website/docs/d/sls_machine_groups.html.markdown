---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_machine_groups"
sidebar_current: "docs-alicloud-datasource-sls-machine-groups"
description: |-
  Provides a list of Sls Machine Group owned by an Alibaba Cloud account.
---

# alicloud_sls_machine_groups

This data source provides Sls Machine Group available to the user.[What is Machine Group](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateMachineGroup)

-> **NOTE:** Available since v1.259.0.

## Example Usage

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

data "alicloud_sls_machine_groups" "default" {
  ids          = ["${alicloud_sls_machine_group.default.id}"]
  group_name   = "group1"
  project_name = var.project_name
}

output "alicloud_sls_machine_group_example_id" {
  value = data.alicloud_sls_machine_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:
* `group_name` - (ForceNew, Optional) Machine Group name
* `project_name` - (Required, ForceNew) Project name
* `ids` - (Optional, ForceNew, Computed) A list of Machine Group IDs. The value is formulated as `<project_name>:<group_name>`.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Machine Group IDs.
* `groups` - A list of Machine Group Entries. Each element contains the following attributes:
  * `group_name` - Machine Group name
  * `id` - The ID of the resource supplied above.
