---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_havip"
sidebar_current: "docs-alicloud-resource-havip"
description: |-
  Provides a Alicloud HaVip resource.
---

# alicloud_havip

Provides a HaVip resource, see [What is HAVIP](https://www.alibabacloud.com/help/zh/vpc/developer-reference/api-createhavip).

-> **NOTE:** Terraform will auto build havip instance  while it uses `alicloud_havip` to build a havip resource.

-> **NOTE:** Available since v1.18.0.

-> **DEPRECATED:**  This resource has been renamed to [alicloud_vpc_ha_vip](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/vpc_ha_vip) from version 1.205.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_havip" "example" {
  vswitch_id  = alicloud_vswitch.example.id
  description = var.name
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_havip&spm=docs.r.havip.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The vswitch_id of the HaVip, the field can't be changed.
* `ip_address` - (Optional, ForceNew) The ip address of the HaVip. If not filled, the default will be assigned one from the vswitch.
* `description` - (Optional) The description of the HaVip instance.
* `havip_name` - (Optional, Deprecated) The name of the HaVip instance.
* `associated_instances` - (Optional) The ID of the instance with which the HAVIP is associated.
* `ha_vip_id` - (Optional) The ID of the HAVIP.
* `resource_group_id` - (Optional) The ID of the resource group to which the HAVIP belongs.
* `associated_eip_addresses` - (Optional) The elastic IP address (EIP) associated with the HAVIP.
* `associated_instance_type` - (Optional) The type of the instance with which the HAVIP is associated. Valid values:
  - `EcsInstance`: an ECS instance.
  - `NetworkInterface`: an ENI.
* `ha_vip_name` - (Optional) The name of the HAVIP.
* `tags` - (Optional) The list of tags.
* `vpc_id` - (Optional) The ID of the VPC to which the HAVIP belongs.
* `create_time` - (Optional) The time when the HAVIP was created.
* `master_instance_id` - (Optional) The ID of the active instance that is associated with the HAVIP.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HaVip instance id.
* `status` - (Available since v1.120.0) The status of the HaVip instance.

## Timeouts

-> **NOTE:** Available since v1.120.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the HaVip instance.
* `update` - (Defaults to 5 mins) Used when updating the HaVip instance.
* `delete` - (Defaults to 5 mins) Used when deleting the HaVip instance.

## Import

The havip can be imported using the id, e.g.

```shell
$ terraform import alicloud_havip.foo havip-abc123456
```




