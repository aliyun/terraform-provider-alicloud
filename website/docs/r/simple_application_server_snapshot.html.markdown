---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_snapshot"
sidebar_current: "docs-alicloud-resource-simple-application-server-snapshot"
description: |-
  Provides a Alicloud Simple Application Server Snapshot resource.
---

# alicloud_simple_application_server_snapshot

Provides a Simple Application Server Snapshot resource.

For information about Simple Application Server Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/doc-detail/190452.htm).

-> **NOTE:** Available since v1.143.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_simple_application_server_snapshot&exampleId=a1423c32-da0a-2905-0ae4-e77eeb1588c419f26a93&activeTab=example&spm=docs.r.simple_application_server_snapshot.0.a1423c32da&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf_example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_simple_application_server_images" "default" {
  platform = "Linux"
}
data "alicloud_simple_application_server_plans" "default" {
  platform = "Linux"
}

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
  snapshot_name = "${var.name}-${random_integer.default.result}"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_simple_application_server_snapshot&spm=docs.r.simple_application_server_snapshot.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, ForceNew) The ID of the disk.
* `snapshot_name` - (Required, ForceNew) The name of the snapshot. The name must be `2` to `50` characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), periods (.),and hyphens (-).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot.
* `status` - The status of the snapshot. Valid values: `Progressing`, `Accomplished` and `Failed`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Snapshot.

## Import

Simple Application Server Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_simple_application_server_snapshot.example <id>
```