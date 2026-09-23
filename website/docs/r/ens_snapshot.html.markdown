---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_snapshot"
description: |-
  Provides a Alicloud ENS Snapshot resource.
---

# alicloud_ens_snapshot

Provides a ENS Snapshot resource. Snapshot. When you use it for the first time, please contact the product classmates to add a resource whitelist.

For information about ENS Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createsnapshot).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_snapshot&exampleId=6f6f2605-0c5c-6362-2efb-c4e990f469f511443b4a&activeTab=example&spm=docs.r.ens_snapshot.0.6f6f26050c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_disk" "disk" {
  category      = "cloud_efficiency"
  size          = "20"
  payment_type  = "PayAsYouGo"
  ens_region_id = "ch-zurich-1"
}

resource "alicloud_ens_snapshot" "default" {
  description   = var.name
  ens_region_id = "ch-zurich-1"
  snapshot_name = var.name

  disk_id = alicloud_ens_disk.disk.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ens_snapshot&spm=docs.r.ens_snapshot.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Snapshot Description Information.
* `disk_id` - (Required, ForceNew) Cloud Disk ID.
* `ens_region_id` - (Required, ForceNew) The node ID of ENS.
* `snapshot_name` - (Optional) Name of the snapshot instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Instance creation timeIt is expressed in accordance with the ISO8601 standard and uses UTC +0 time in the format of yyyy-MM-ddTHH:mm:ssZ.Example value: 2020-08-20 T14:52:28Z.
* `status` - Snapshot Status. Valid values: creating, available, deleting, error.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Snapshot.
* `delete` - (Defaults to 5 mins) Used when delete the Snapshot.
* `update` - (Defaults to 5 mins) Used when update the Snapshot.

## Import

ENS Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_snapshot.example <id>
```