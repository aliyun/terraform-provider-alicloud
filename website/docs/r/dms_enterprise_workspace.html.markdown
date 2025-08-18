---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_workspace"
description: |-
  Provides a Alicloud DMS Enterprise Workspace resource.
---

# alicloud_dms_enterprise_workspace

Provides a DMS Enterprise Workspace resource.



For information about DMS Enterprise Workspace and how to use it, see [What is Workspace](https://next.api.alibabacloud.com/document/dms-enterprise/2018-11-01/CreateWorkspace).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_vpc" "vpc_create" {
  is_default  = false
  description = "example vpc"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_dms_enterprise_workspace" "default" {
  description    = var.name
  workspace_name = "${var.name}-${random_integer.default.result}"
  vpc_id         = alicloud_vpc.vpc_create.id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Required) The description of the Workspace.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `workspace_name` - (Required) The name of the Workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Workspace.
* `delete` - (Defaults to 5 mins) Used when delete the Workspace.
* `update` - (Defaults to 5 mins) Used when update the Workspace.

## Import

DMS Enterprise Workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_dms_enterprise_workspace.example <id>
```
