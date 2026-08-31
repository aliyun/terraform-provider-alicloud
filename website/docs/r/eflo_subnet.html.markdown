---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_subnet"
sidebar_current: "docs-alicloud-resource-eflo-subnet"
description: |-
  Provides a Alicloud Eflo Subnet resource.
---

# alicloud_eflo_subnet

Provides a Eflo Subnet resource.

For information about Eflo Subnet and how to use it, see [What is Subnet](https://www.alibabacloud.com/help/en/pai/user-guide/overview-of-intelligent-computing-lingjun).

-> **NOTE:** Available since v1.204.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_subnet&exampleId=19e968f8-f184-41fc-f07f-b70cac556a686e99beff&activeTab=example&spm=docs.r.eflo_subnet.0.19e968f8f1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}
data "alicloud_zones" "default" {}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eflo_vpd" "default" {
  cidr              = "10.0.0.0/8"
  vpd_name          = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

resource "alicloud_eflo_subnet" "default" {
  subnet_name = var.name
  zone_id     = data.alicloud_zones.default.zones.0.id
  cidr        = "10.0.0.0/16"
  vpd_id      = alicloud_eflo_vpd.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eflo_subnet&spm=docs.r.eflo_subnet.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cidr` - (Required, ForceNew) CIDR network segment.
* `subnet_name` - (Required) The Subnet name.
* `type` - (Optional, ForceNew) Eflo subnet usage type. optional value:
  - General type is not filled in
  - OOB:OOB type
  - LB: LB type
* `zone_id` - (Required, ForceNew) The zone ID  of the resource.
* `vpd_id` - (Required, ForceNew) The Eflo VPD ID.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<vpd_id>:<subnet_id>`.
* `create_time` - The creation time of the resource.
* `gmt_modified` - Modification time.
* `message` - Error message.
* `resource_group_id` - Resource Group ID.
* `status` - The status of the resource.
* `subnet_id` - The id of the subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Subnet.
* `delete` - (Defaults to 5 mins) Used when delete the Subnet.
* `update` - (Defaults to 5 mins) Used when update the Subnet.

## Import

Eflo Subnet can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_subnet.example <vpd_id>:<subnet_id>
```