---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disk_default_kms_key"
description: |-
  Provides a Alicloud Ecs Disk Default Kms Key resource.
---

# alicloud_ecs_disk_default_kms_key

Provides a Ecs Disk Default Kms Key resource.

The encryption key used by default for cloud storage encryption.

For information about Ecs Disk Default Kms Key and how to use it, see [What is Disk Default Kms Key](https://next.api.alibabacloud.com/document/Ecs/2014-05-26/ModifyDiskDefaultKMSKeyId).

-> **NOTE:** Available since v1.278.0.

-> **NOTE:** Destroying this resource will reset the default CMK to the account's AliCloud-managed default CMK for Ecs.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_disk_default_kms_key&exampleId=694c9cea-74cf-ec2e-6bbf-f3934ddad03d973049b0&activeTab=example&spm=docs.r.ecs_disk_default_kms_key.0.694c9cea74&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_kms_keys" "default" {
  filters = "[{\"Key\":\"KeyState\",\"Values\":[\"Enabled\"]}]"
}

resource "alicloud_ecs_disk_default_kms_key" "default" {
  kms_key_id = data.alicloud_kms_keys.default.ids.0
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_disk_default_kms_key&spm=docs.r.ecs_disk_default_kms_key.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `kms_key_id` - (Required) The ID of the KMS key.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Disk Default Kms Key.
* `delete` - (Defaults to 5 mins) Used when delete the Disk Default Kms Key.
* `update` - (Defaults to 5 mins) Used when update the Disk Default Kms Key.

## Import

Ecs Disk Default Kms Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_disk_default_kms_key.example <region_id>
```
