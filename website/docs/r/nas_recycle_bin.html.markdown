---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_recycle_bin"
sidebar_current: "docs-alicloud-resource-nas-recycle-bin"
description: |-
  Provides a Alicloud File Storage (NAS) Recycle Bin resource.
---

# alicloud\_nas\_recycle\_bin

Provides a File Storage (NAS) Recycle Bin resource.

For information about File Storage (NAS) Recycle Bin and how to use it, see [What is Recycle Bin](https://www.alibabacloud.com/help/en/doc-detail/264185.html).

-> **NOTE:** Available in v1.155.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_recycle_bin&exampleId=3d880a4f-bcc5-f243-1c79-905439ec95b411821d47&activeTab=example&spm=docs.r.nas_recycle_bin.0.3d880a4fbc&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_nas_zones" "example" {
  file_system_type = "standard"
}

resource "alicloud_nas_file_system" "example" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = "terraform-example"
  encrypt_type  = "1"
  zone_id       = data.alicloud_nas_zones.example.zones[0].zone_id
}

resource "alicloud_nas_recycle_bin" "example" {
  file_system_id = alicloud_nas_file_system.example.id
  reserved_days  = 3
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_nas_recycle_bin&spm=docs.r.nas_recycle_bin.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system for which you want to enable the recycle bin feature.
* `reserved_days` - (Optional, Computed) The period for which the files in the recycle bin are retained. Unit: days. Valid values: `1` to `180`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Recycle Bin. Its value is same as `file_system_id`.
* `status` - The status of the recycle bin.

## Import

File Storage (NAS) Recycle Bin can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_recycle_bin.example <file_system_id>
```