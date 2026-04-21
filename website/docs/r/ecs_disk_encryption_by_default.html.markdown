---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disk_encryption_by_default"
description: |-
  Provides a Alicloud Ecs Disk Encryption By Default resource.
---

# alicloud_ecs_disk_encryption_by_default

Provides a Ecs Disk Encryption By Default resource.

Default encryption configuration capability for cloud storage.

For information about Ecs Disk Encryption By Default and how to use it, see [What is Disk Encryption By Default](https://next.api.alibabacloud.com/document/Ecs/2014-05-26/EnableDiskEncryptionByDefault).

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_disk_encryption_by_default" "default" {
  encrypted = true
}
```

### Deleting `alicloud_ecs_disk_encryption_by_default` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ecs_disk_encryption_by_default`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `encrypted` - (Optional, Bool) Indicates whether account-level default encryption of EBS resources is enabled in the region. Valid values:
  - `true`: Enable.
  - `false`: Disable.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Disk Encryption By Default.
* `update` - (Defaults to 5 mins) Used when update the Disk Encryption By Default.

## Import

Ecs Disk Encryption By Default can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_disk_encryption_by_default.example <region_id>
```
