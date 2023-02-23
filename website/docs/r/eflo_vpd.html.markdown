---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vpd"
sidebar_current: "docs-alicloud-resource-eflo-vpd"
description: |-
  Provides a Alicloud Eflo Vpd resource.
---

# alicloud_eflo_vpd

Provides a Eflo Vpd resource.

For information about Eflo Vpd and how to use it, see [What is Vpd](https://help.aliyun.com/document_detail/604976.html).

-> **NOTE:** Available in v1.201.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_eflo_vpd" "default" {
  cidr     = "10.0.0.0/8"
  vpd_name = "RMC-Terraform-Test"
}
```

## Argument Reference

The following arguments are supported:
* `cidr` - (ForceNew,Required) CIDR network segment
* `resource_group_id` - (ForceNew,Optional) The Resource group id
* `vpd_name` - (Required) The Name of the VPD.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource
* `gmt_modified` - Modification time
* `status` - The Vpd status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpd.
* `delete` - (Defaults to 5 mins) Used when delete the Vpd.
* `update` - (Defaults to 5 mins) Used when update the Vpd.

## Import

Eflo Vpd can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_vpd.example <id>
```