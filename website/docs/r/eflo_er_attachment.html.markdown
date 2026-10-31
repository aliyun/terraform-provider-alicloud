---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_er_attachment"
description: |-
  Provides a Alicloud Eflo Er Attachment resource.
---

# alicloud_eflo_er_attachment

Provides a Eflo Er Attachment resource.

Network instance connection.

For information about Eflo Er Attachment and how to use it, see [What is Er Attachment](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateErAttachment).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** The service does not always return `across`, `create_time` and `gmt_modified`; these attributes stay empty when the response omits them.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

data "alicloud_account" "current" {}

resource "alicloud_eflo_er" "default" {
  er_name        = var.name
  master_zone_id = "cn-wulanchabu-b"
}

resource "alicloud_eflo_vpd" "default" {
  cidr     = "10.0.0.0/8"
  vpd_name = var.name
}

resource "alicloud_eflo_er_attachment" "default" {
  instance_id            = alicloud_eflo_vpd.default.id
  resource_tenant_id     = data.alicloud_account.current.id
  er_id                  = alicloud_eflo_er.default.id
  er_attachment_name     = var.name
  instance_type          = "VPD"
  auto_receive_all_route = true
}
```

## Argument Reference

The following arguments are supported:
* `auto_receive_all_route` - (Required, ForceNew, Bool) Automatically receive all routes.
* `er_attachment_name` - (Required, String) Network instance connection name.
* `er_id` - (Required, ForceNew, String) Lingjun HUB.
* `instance_id` - (Required, ForceNew, String) Instance ID.
* `instance_type` - (Required, ForceNew, String) Instance type. Valid values: `VPD`, `VCC`.
* `resource_tenant_id` - (Optional, Computed, ForceNew, String) Tenant ID to which the resource belongs. The value only causes a replacement when the configured value changes; when it is left unset the value returned by the service is used.

## Attributes Reference

The following attributes are exported:
* `id` - (Computed, String) The ID of the resource supplied above. The value is formulated as `<er_id>:<er_attachment_id>`.
* `across` - (Computed, Bool) Cross account resources.
* `create_time` - (Computed, String) Creation time.
* `er_attachment_id` - (Computed, String) The connection ID of the network instance.
* `gmt_modified` - (Computed, String) Modification time.
* `instance_name` - (Computed, String) Instance name.
* `message` - (Computed, String) Prompt message.
* `status` - (Computed, String) Status.
* `tenant_id` - (Computed, String) Tenant ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - Defaults to `5 mins`. Used when create the Er Attachment.
* `delete` - Defaults to `5 mins`. Used when delete the Er Attachment.
* `update` - Defaults to `5 mins`. Used when update the Er Attachment.

## Import

Eflo Er Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_er_attachment.example <er_id>:<er_attachment_id>
```