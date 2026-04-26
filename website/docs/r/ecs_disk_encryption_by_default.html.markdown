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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_disk_encryption_by_default&exampleId=7721a500-9b55-82d6-81b4-ff1d0791ba7dbfcdae14&activeTab=example&spm=docs.r.ecs_disk_encryption_by_default.0.7721a5009b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ecs_disk_encryption_by_default" "default" {
  encrypted = true
}
```

### Deleting `alicloud_ecs_disk_encryption_by_default` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ecs_disk_encryption_by_default`. Terraform will remove this resource from the state file, however resources may remain.


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_disk_encryption_by_default&spm=docs.r.ecs_disk_encryption_by_default.example&intl_lang=EN_US)


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
