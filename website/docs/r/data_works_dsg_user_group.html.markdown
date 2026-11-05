---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_dsg_user_group"
description: |-
  Provides a Alicloud Data Works Data Security Guard (DSG) User Group resource.
---

# alicloud_data_works_dsg_user_group

Provides a Data Works Data Security Guard (DSG) User Group resource.

For information about Data Works Data Security Guard (DSG) User Group and how to use it, see [DsgUserGroupAddOrUpdate](https://next.api.alibabacloud.com/document/DataWorks/2020-05-18/DsgUserGroupAddOrUpdate).

-> **NOTE:** Available since v2.3.0.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_data_works_dsg_user_group" "default" {
  name            = var.name
  owner           = "user@example.com"
  user_group_type = 1
  accounts        = ["user@example.com"]
}
```

The following example creates a MaxCompute role user group (`user_group_type = 3`), which requires `project_name`.

```terraform
resource "alicloud_data_works_dsg_user_group" "maxcompute" {
  name            = var.name
  owner           = "user@example.com"
  user_group_type = 3
  project_name    = "mc_project"
}
```

## Argument Reference

The following arguments are supported:

* `accounts` - (Optional, Set of String) The accounts that belong to the user group.
* `name` - (Required, ForceNew, String) The name of the user group. The parameter is immutable after resource creation.
* `owner` - (Required, String) The owner of the user group.
* `project_name` - (Optional, ForceNew, String) The MaxCompute project name. **Note: This parameter is required when `user_group_type` is set to `3` (MaxCompute role).** The parameter is write-only: it is not returned by the query API, so the configured value is preserved in state.
* `user_group_type` - (Required, ForceNew, Int) The type of the user group. Valid values: `1`, `2`, `3`. A value of `3` indicates a MaxCompute role group and requires `project_name`. The parameter is immutable after resource creation.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed, String) The ID of the user group.
* `gmt_create` - (Computed, String) The creation time of the user group.
* `gmt_modified` - (Computed, String) The last modification time of the user group.

## Import

Data Works Data Security Guard (DSG) User Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_dsg_user_group.example <user_group_id>
```
