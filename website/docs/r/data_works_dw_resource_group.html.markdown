---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_dw_resource_group"
description: |-
  Provides a Alicloud Data Works Dw Resource Group resource.
---

# alicloud_data_works_dw_resource_group

Provides a Data Works Dw Resource Group resource.



For information about Data Works Dw Resource Group and how to use it, see [What is Dw Resource Group](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_data_works_project" "defaultZImuCO" {
  description  = "default_pj002"
  project_name = var.name
  display_name = "default_pj002"
}

resource "alicloud_vpc" "defaulte4zhaL" {
  description = "default_resgv2_vpc001"
  vpc_name    = format("%s1", var.name)
  cidr_block  = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default675v38" {
  description  = "default_resg_vsw001"
  vpc_id       = alicloud_vpc.defaulte4zhaL.id
  zone_id      = "cn-beijing-g"
  vswitch_name = format("%s2", var.name)
  cidr_block   = "172.16.0.0/24"
}


resource "alicloud_data_works_dw_resource_group" "default" {
  payment_type          = "PostPaid"
  default_vpc_id        = alicloud_vpc.defaulte4zhaL.id
  remark                = "openapi_example"
  resource_group_name   = "openapi_pop2_example_resg00002"
  default_vswitch_id    = alicloud_vswitch.default675v38.id
  payment_duration_unit = "Month"
  specification         = "500"
  payment_duration      = "1"
}
```

## Argument Reference

The following arguments are supported:
* `auto_renew` - (Optional) Whether to automatically renew. The default value is false.
* `default_vpc_id` - (Required, ForceNew) The VPC ID of the default network resource.
* `default_vswitch_id` - (Required, ForceNew) The ID of the vswitch bound to the network resource by default.
* `payment_duration` - (Optional, Int) When the payment method is Prepaid, the unit is PaymentDurationUnit.
* `payment_duration_unit` - (Optional) When the payment method is Prepaid, the payment duration unit, Month is Month, Year is Year, and there is no other optional value.
* `payment_type` - (Optional, ForceNew, Computed) The billing type of the resource group. PrePaid is Subscription, and PostPaid is Pay-As-You-Go.
* `remark` - (Required) Resource Group Comments
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `resource_group_name` - (Optional) Resource group name
* `specification` - (Optional, Int) Package year and package month resource group specifications, unit CU
* `tags` - (Optional, Computed, Map) The tag of the resource

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Resource group creation time
* `status` - Resource group status:

  Normal: Normal (in operation/service);

  Stop: Freeze (expired);

  Deleted: Deleted (released/destroyed);

  Creating: Creating (starting);

  CreateFailed: creation failed (startup failed);

  Updating: Updating (expanding/shrinking/changing);

  UpdateFailed: update failed (expansion failed/upgrade failed);

  Deleting: Deleting (releasing/destroying);

  DeleteFailed: delete failed (release failed/destroy failed);

  Timeout: Timeout.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dw Resource Group.
* `delete` - (Defaults to 5 mins) Used when delete the Dw Resource Group.
* `update` - (Defaults to 5 mins) Used when update the Dw Resource Group.

## Import

Data Works Dw Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_dw_resource_group.example <id>
```