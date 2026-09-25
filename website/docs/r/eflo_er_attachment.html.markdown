---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_er_attachment"
description: |-
  Provides a Alicloud Eflo Er Attachment resource.
---

# alicloud_eflo_er_attachment

Provides a Eflo Er Attachment resource.

Attaches a Lingjun network instance (VPD or VCC) to a Lingjun HUB.

For information about Eflo Er Attachment and how to use it, see [What is Er Attachment](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateErAttachment).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** The network instance connection can only be deleted after all routing strategies (`alicloud_eflo_er_route_map`) that reference the attached instance are deleted.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

variable "zone_id" {
  default = "cn-wulanchabu-b"
}

resource "alicloud_eflo_vpd" "default" {
  cidr     = "10.0.0.0/8"
  vpd_name = "tf-vpd-${var.name}"
}

resource "alicloud_eflo_er" "default" {
  er_name        = "tf-er-${var.name}"
  master_zone_id = var.zone_id
}

resource "alicloud_eflo_er_attachment" "default" {
  er_id                  = alicloud_eflo_er.default.id
  er_attachment_name     = "tf-attachment-${var.name}"
  instance_id            = alicloud_eflo_vpd.default.id
  instance_type          = "VPD"
  auto_receive_all_route = false
}
```

## Argument Reference

The following arguments are supported:
* `auto_receive_all_route` - (Required, ForceNew) Specifies whether to automatically receive all routes from all instances that are attached to the Lingjun HUB.
* `er_attachment_name` - (Required) The name of the network instance connection.
* `er_id` - (Required, ForceNew) The ID of the Lingjun HUB.
* `instance_id` - (Required, ForceNew) The ID of the network instance. Valid values are Lingjun VPDs and Lingjun VCCs.
* `instance_type` - (Required, ForceNew) The type of the network instance. Valid values: `VPD`, `VCC`.
* `resource_tenant_id` - (Optional, Computed, ForceNew) The tenant ID of the network instance. This parameter is required for cross-account instances.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<er_id>:<er_attachment_id>`.
* `across` - Whether the attached network instance belongs to another account.
* `create_time` - The creation time of the network instance connection.
* `er_attachment_id` - The ID of the network instance connection.
* `gmt_modified` - The time when the network instance connection was modified.
* `instance_name` - The name of the network instance.
* `message` - The status message of the network instance connection.
* `region_id` - The region ID of the resource.
* `resource_group_id` - The ID of the resource group.
* `status` - The status of the network instance connection.
* `tenant_id` - The tenant ID of the network instance connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Er Attachment.
* `delete` - (Defaults to 20 mins) Used when delete the Er Attachment.
* `update` - (Defaults to 10 mins) Used when update the Er Attachment.

## Import

Eflo Er Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_er_attachment.example <er_id>:<er_attachment_id>
```
