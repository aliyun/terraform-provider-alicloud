---
subcategory: "Cloud Phone"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_phone_instance_group"
description: |-
  Provides a Alicloud Cloud Phone Instance Group resource.
---

# alicloud_cloud_phone_instance_group

Provides a Cloud Phone Instance Group resource.



For information about Cloud Phone Instance Group and how to use it, see [What is Instance Group](https://next.api.alibabacloud.com/document/eds-aic/2023-09-30/CreateAndroidInstanceGroup).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "region_id" {
  default = "cn-hangzhou"
}

resource "alicloud_cloud_phone_policy" "defaultjZ1gi0" {
  lock_resolution   = "off"
  resolution_width  = "720"
  camera_redirect   = "on"
  policy_group_name = "defaultPolicyGroup"
  resolution_height = "1280"
  clipboard         = "readwrite"
  net_redirect_policy {
    net_redirect = "off"
    custom_proxy = "off"
  }
}

resource "alicloud_ecd_simple_office_site" "defaultH2a5KS" {
  office_site_name = "InitOfficeSite"
  cidr_block       = "172.16.0.0/12"
}


resource "alicloud_cloud_phone_instance_group" "default" {
  instance_group_spec = "acp.basic.small"
  instance_group_name = "AutoCreateGroupName"
  period              = "1"
  number_of_instances = "1"
  charge_type         = "PostPaid"
  image_id            = "imgc-075cllfeuazh03tg9"
  period_unit         = "Hour"
  auto_renew          = false
  amount              = "1"
  auto_pay            = false
  gpu_acceleration    = false
  policy_group_id     = alicloud_cloud_phone_policy.defaultjZ1gi0.id
  office_site_id      = alicloud_ecd_simple_office_site.defaultH2a5KS.id
}
```

## Argument Reference

The following arguments are supported:
* `amount` - (Optional, Int) The number of instance groups. The default value is 1 and the maximum value is 100.
* `auto_pay` - (Optional) Whether to pay automatically. The default is false.
* `auto_renew` - (Optional) Whether to enable automatic renewal. The default is false.
* `charge_type` - (Optional, ForceNew) The billing type.
* `gpu_acceleration` - (Optional) Whether to enable GPU acceleration. The default value is false.
  - true: On.
  - false: closed.
* `image_id` - (Required, ForceNew) The image ID. 
* `instance_group_name` - (Optional) The instance group name

-> **NOTE:** >

-> **NOTE:** - The instance group name must be no more than 30 characters in length. Start with an uppercase/lowercase letter or Chinese. It cannot start with http:// or https://. Only Chinese, English, numbers, half-width colons (:), underscores (_), periods (.), or hyphens (-) are supported.

* `instance_group_spec` - (Required, ForceNew) Instance group specifications. 
* `number_of_instances` - (Optional, ForceNew, Int) The number of instances in the instance group. The maximum value is 100.
* `office_site_id` - (Optional) The network ID.
  - Create a shared network instance: Network ID is optional. Enter the network ID whose type is **Shared Network** on the [cloud mobile phone console> Network](https://wya.wuying.aliyun.com/network) page. If the console does not have a shared network, you can fill it in. A shared network is automatically created when the instance group is created.
  - Create a VPC instance: the network ID is required. Enter the network ID of `VPC` on the [cloud mobile phone console> Network](https://wya.wuying.aliyun.com/network) page. If the console does not have a VPC network, you need to create a network first.
* `period` - (Optional, Int) The duration of the resource purchase. The unit is specified by PeriodUnit.
* `period_unit` - (Optional) The unit of time for purchasing resources.
* `policy_group_id` - (Optional) The policy ID. You can query the list of policies by calling [ListPolicyGroups](~~ ListPolicyGroups ~~).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - Instance group status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance Group.
* `delete` - (Defaults to 5 mins) Used when delete the Instance Group.
* `update` - (Defaults to 5 mins) Used when update the Instance Group.

## Import

Cloud Phone Instance Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_phone_instance_group.example <id>
```