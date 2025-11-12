---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vpd"
description: |-
  Provides a Alicloud Eflo Vpd resource.
---

# alicloud_eflo_vpd

Provides a Eflo Vpd resource.

Lingjun Network Segment.

For information about Eflo Vpd and how to use it, see [What is Vpd](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateVpd).

-> **NOTE:** Available since v1.201.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_vpd&exampleId=befafc06-85e5-7365-f8bc-16b6760e59a82cb4f3b4&activeTab=example&spm=docs.r.eflo_vpd.0.befafc0685&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_eflo_vpd" "default" {
  cidr              = "10.0.0.0/8"
  vpd_name          = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:
* `cidr` - (Required, ForceNew) The CIDR block of the VPD.
* `resource_group_id` - (Optional) The Resource group ID. **NOTE:** From version 1.263.0, `resource_group_id` can be modified.
* `secondary_cidr_blocks` - (Optional, List, Available since v1.263.0) The additional CIDR block.
* `tags` - (Optional, Map, Available since v1.263.0) The tag of the resource.
* `vpd_name` - (Required) The name of the VPD instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the activation code was created.
* `gmt_modified` - The time when the O&M task was modified.
* `region_id` - (Available since v1.263.0) The region ID.
* `status` - The current state of the instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpd.
* `delete` - (Defaults to 5 mins) Used when delete the Vpd.
* `update` - (Defaults to 5 mins) Used when update the Vpd.

## Import

Eflo Vpd can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_vpd.example <id>
```
