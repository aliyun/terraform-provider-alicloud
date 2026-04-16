---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_eni_sg_attachment"
sidebar_current: "docs-alicloud-resource-ecs-eni-sg-attachment"
description: |-
  Provides an Alicloud ECS ENI Security Group attachment resource.
---

# alicloud_ecs_eni_sg_attachment

Provides an ECS ENI security group attachment resource.

This resource merges `attach_security_group_ids` into the current ENI security group relationship and preserves existing groups that are not explicitly attached.
Because ECS API may not preserve list order in responses, this resource treats security group relationships as a set for state convergence.

When this resource is destroyed, it restores the ENI security groups to the original list captured at creation time.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_eni_sg_attachment" "default" {
  network_interface_id = "eni-xxxxxxxxxxxx"

  # Current ENI security groups: [sg-1, sg-2, sg-3]
  # attach_security_group_ids: [sg-4, sg-3]
  # Effective relationship after apply contains: sg-1, sg-2, sg-3, sg-4
  # (order is not guaranteed by ECS API)
  attach_security_group_ids = ["sg-4", "sg-3"]
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_id` - (Required, ForceNew) The ID of the network interface.
* `attach_security_group_ids` - (Required, ForceNew, List) Security group IDs to merge into the network interface's existing security group relationship.
  **NOTE:** The final relationship is guaranteed, but returned order from ECS may differ.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, same as `network_interface_id`.
* `original_security_group_ids` - The original security group IDs captured before attachment apply.
* `effective_security_group_ids` - The effective security group IDs currently observed on the network interface.
