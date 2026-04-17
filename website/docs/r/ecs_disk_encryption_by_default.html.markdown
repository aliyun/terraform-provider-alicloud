---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disk_encryption_by_default"
description: |-
  Provides a Alicloud ECS Disk Encryption By Default resource.
---

# alicloud_ecs_disk_encryption_by_default

Provides an ECS Disk Encryption By Default resource to enable or disable account-level default disk encryption.

For information about ECS Disk Encryption By Default and how to use it, see [What is Disk Encryption By Default](https://www.alibabacloud.com/help/en/doc-detail/59643.htm).

-> **NOTE:** Available since v1.274.0.

-> **NOTE:** This resource manages account-level disk encryption settings for the current region. Once enabled, all new disks created in the region will be encrypted by default.

-> **NOTE:** You need to have KMS (Key Management Service) enabled before using this resource.

## Example Usage

Basic Usage

```terraform
# Enable ECS disk encryption by default
resource "alicloud_ecs_disk_encryption_by_default" "default" {
  enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Optional) Whether to enable ECS disk encryption by default. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is the region ID where the resource is located.

## Import

ECS Disk Encryption By Default can be imported using the region id, e.g.

```shell
$ terraform import alicloud_ecs_disk_encryption_by_default.example cn-hangzhou
```