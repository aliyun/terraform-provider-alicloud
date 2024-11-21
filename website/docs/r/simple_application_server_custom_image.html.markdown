---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_custom_image"
sidebar_current: "docs-alicloud-resource-simple-application-server-custom-image"
description: |-
  Provides a Alicloud Simple Application Server Custom Image resource.
---

# alicloud_simple_application_server_custom_image

Provides a Simple Application Server Custom Image resource.

For information about Simple Application Server Custom Image and how to use it, see [What is Custom Image](https://www.alibabacloud.com/help/en/doc-detail/333535.htm).

-> **NOTE:** Available since v1.143.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_simple_application_server_custom_image&exampleId=68f8c52d-bfa6-7037-76d1-2f176a7d4a96bbd03de9&activeTab=example&spm=docs.r.simple_application_server_custom_image.0.68f8c52dbf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_simple_application_server_images" "default" {}
data "alicloud_simple_application_server_plans" "default" {}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

data "alicloud_simple_application_server_disks" "default" {
  instance_id = alicloud_simple_application_server_instance.default.id
}

resource "alicloud_simple_application_server_snapshot" "default" {
  disk_id       = data.alicloud_simple_application_server_disks.default.ids.0
  snapshot_name = var.name
}

resource "alicloud_simple_application_server_custom_image" "default" {
  custom_image_name  = var.name
  instance_id        = alicloud_simple_application_server_instance.default.id
  system_snapshot_id = alicloud_simple_application_server_snapshot.default.id
  status             = "Share"
  description        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `custom_image_name` - (Required, ForceNew) The name of the resource. The name must be `2` to `128` characters in length. It must start with a letter or a number. It can contain letters, digits, colons (:), underscores (_) and hyphens (-).
* `description` - (Optional, ForceNew) The description of the Custom Image.
* `instance_id` - (Required) The ID of the instance.
* `system_snapshot_id` - (Required) The ID of the system snapshot.
* `status` - (Optional) The Shared status of the Custom Image. Valid values: `Share`, `UnShare`.

 **NOTE:** The `status` will be automatically change to `UnShare` when the resource is deleted, please operate with caution.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Image.

## Import

Simple Application Server Custom Image can be imported using the id, e.g.

```shell
$ terraform import alicloud_simple_application_server_custom_image.example <id>
```