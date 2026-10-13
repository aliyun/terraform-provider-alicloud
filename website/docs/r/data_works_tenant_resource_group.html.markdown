---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_tenant_resource_group"
description: |-
  Provides a Alicloud Data Works Tenant Resource Group resource.
---

# alicloud_data_works_tenant_resource_group

Provides a Data Works Tenant Resource Group resource.

For information about Data Works Tenant Resource Group and how to use it, see [What is Tenant Resource Group](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2020-05-18-createresourcegroup).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "default" {
  vpc_name   = format("%s_vpc", var.name)
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = "cn-hangzhou-h"
  vswitch_name = format("%s_vsw", var.name)
}

resource "alicloud_data_works_tenant_resource_group" "default" {
  tenant_resource_group_name        = var.name
  tenant_resource_group_description = "terraform example"
  payment_type                      = "PostPaid"
  vpc_id                            = alicloud_vpc.default.id
  vswitch_id                        = alicloud_vswitch.default.id
  tags = {
    Created = "TF"
    For     = "TenantRG"
  }
}
```

## Argument Reference

The following arguments are supported:

* `tenant_resource_group_name` - (Required, ForceNew) The name of the tenant resource group. The name can be a maximum of 128 characters in length and can contain letters, digits, and underscores (\_). The name must start with a letter.
* `tenant_resource_group_description` - (Optional) The description of the tenant resource group. The description can be a maximum of 128 characters in length.
* `payment_type` - (Required, ForceNew) The billing method of the resource group. Valid values: `PostPaid` (pay-as-you-go) and `PrePaid` (subscription).
* `vpc_id` - (Required, ForceNew) The ID of the VPC with which the resource group is associated by default.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch with which the resource group is associated by default.
* `aliyun_resource_group_id` - (Optional, ForceNew) The ID of the Alibaba Cloud resource group to which the resource group belongs.
* `payment_duration` - (Optional, ForceNew) The subscription duration. This parameter is required only when `payment_type` is set to `PrePaid`.
* `payment_duration_unit` - (Optional, ForceNew) The unit of the subscription duration. Valid values: `Month` and `Year`. This parameter is required only when `payment_type` is set to `PrePaid`.
* `auto_renew_enabled` - (Optional, ForceNew) Specifies whether to enable auto-renewal. This parameter takes effect only when `payment_type` is set to `PrePaid`.
* `spec` - (Optional, ForceNew) The specifications of the resource group, in CU. This parameter is required only when `payment_type` is set to `PrePaid`.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the tenant resource group.
* `status` - The status of the resource group.
* `create_time` - The creation time of the resource group.
* `resource_group_type` - The type of the resource group.
* `order_instance_id` - The order instance ID of the resource group.
* `create_user` - The user who created the resource group.
* `spec_amount` - The amount of the resource group specification.
* `spec_standard` - The standard of the resource group specification.

## Import

Data Works Tenant Resource Group can be imported using the `id`, e.g.

```shell
$ terraform import alicloud_data_works_tenant_resource_group.example <id>
```
