---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_group_snapshot_policy"
description: |-
  Provides a Alicloud ECS Security Group Snapshot Policy resource.
---

# alicloud_security_group_snapshot_policy

Provides a ECS Security Group Snapshot Policy resource.

For information about ECS Security Group Snapshot Policy and how to use it, see [Alibaba Cloud ECS](https://www.alibabacloud.com/help/en/ecs/).

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_security_group_snapshot_policy" "default" {
  security_group_snapshot_policy_name = "terraform-example"
  status                              = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_snapshot_policy_name` - (Optional) The name of the security group snapshot policy.
* `status` - (Optional) The status of the security group snapshot policy. Valid values: `enable` and `disable`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. It is the same as `security_group_snapshot_policy_id`.
* `security_group_snapshot_policy_id` - The ID of the security group snapshot policy.
* `region_id` - The region ID of the security group snapshot policy.

## Import

ECS Security Group Snapshot Policy can be imported using the id, e.g.

```shell
terraform import alicloud_security_group_snapshot_policy.default <id>
```
