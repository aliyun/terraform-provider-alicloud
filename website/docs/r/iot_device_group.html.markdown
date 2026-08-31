---
subcategory: "Internet of Things (Iot)"
layout: "alicloud"
page_title: "Alicloud: alicloud_iot_device_group"
sidebar_current: "docs-alicloud-resource-iot-device-group"
description: |-
  Provides a Alicloud Iot Device Group resource.
---

# alicloud_iot_device_group

Provides a Iot Device Group resource.

For information about Iot Device Group and how to use it, see [What is Device Group](https://www.alibabacloud.com/help/product/30520.htm).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_iot_device_group&exampleId=a177a88b-84e0-5a56-f6d4-4cf61be5f76f18dbe55c&activeTab=example&spm=docs.r.iot_device_group.0.a177a88b84&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
resource "alicloud_iot_device_group" "example" {
  group_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_iot_device_group&spm=docs.r.iot_device_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `group_desc` - (Optional) The GroupDesc of the device group.
* `group_name` - (Required, ForceNew) The GroupName of the device group.
* `iot_instance_id` - (Optional) The id of the Iot Instance.
* `super_group_id` - (Optional, ForceNew) The id of the SuperGroup.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Device Group.

## Import

Iot Device Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_iot_device_group.example <id>
```
