---
subcategory: "Pai Studio"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_resource_group"
description: |-
  Provides a PAI Studio Resource Group resource.
---

# alicloud_pai_resource_group

Provides a PAI Studio Resource Group resource.

For information about PAI Studio Resource Group and how to use it, see [What is Resource Group](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paistudio-2022-01-12-createresourcegroup).

-> **NOTE:** Available since v1.249.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

resource "alicloud_pai_resource_group" "default" {
  name                        = var.name
  description                 = "tf-test-resource-group"
  resource_type               = "Ecs"
  computing_resource_provider = "Ecs"
  version                     = "1.0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the Resource Group.
* `description` - (Optional) The description of the Resource Group.
* `resource_type` - (Optional, ForceNew) The compute resource type of the Resource Group. Valid values: `Ecs` (general compute) and `Lingjun` (Lingjun intelligent compute).
* `computing_resource_provider` - (Optional, ForceNew) The computing resource provider of the Resource Group. Valid values: `Ecs` and `Lingjun`. **NOTE:** The API does not return this attribute for Ecs resource groups on read; the provider preserves the configured value in state, so it is **ignored during `terraform import`**.
* `version` - (Optional, ForceNew) The version of the Resource Group. It takes effect only when `resource_type` is `Ecs`.
* `unbind` - (Optional) Whether to unbind the Resource Group from the currently connected VPC. Write-only: the API does not return this attribute on read, so the value is preserved from configuration.
* `tags` - (Optional, ForceNew) The tags of the Resource Group. At most 20 entries.

The `user_vpc` block supports:

* `vpc_id` - (Optional) The ID of the VPC.
* `switch_id` - (Optional) The ID of the vSwitch.
* `security_group_id` - (Optional) The ID of the security group.
* `role_arn` - (Optional) The ARN of the role used for OSS access authorization.
* `default_route` - (Optional) The default route. Valid values: `eth0` (public gateway) and `eth1` (user ENI).
* `extended_cidrs` - (Optional) The extended CIDR blocks.

The `default_forward_info` block under `user_vpc` supports:

* `eip_allocation_id` - (Optional) The EIP allocation ID.
* `nat_gateway_id` - (Optional) The NAT gateway ID.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource Group (ResourceGroupID).
* `status` - The status of the Resource Group.
* `cluster_id` - The ID of the cluster to which the Resource Group belongs.
* `creator_id` - The user ID of the Resource Group creator.
* `workspace_id` - The ID of the workspace.
* `support_rdma` - Whether the Resource Group supports RDMA.
* `gmt_created_time` - The creation time of the Resource Group, in ISO 8601 format (UTC).
* `gmt_modified_time` - The modification time of the Resource Group, in ISO 8601 format (UTC).

## Import

PAI Studio Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_resource_group.example <id>
```

-> **NOTE:** `computing_resource_provider` is ignored during import because the API does not return it for Ecs resource groups.
