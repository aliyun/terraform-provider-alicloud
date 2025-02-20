---
subcategory: "Cloud Phone"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_phone_instance"
description: |-
  Provides a Alicloud Cloud Phone Instance resource.
---

# alicloud_cloud_phone_instance

Provides a Cloud Phone Instance resource.

cloud phone instance.

For information about Cloud Phone Instance and how to use it, see [What is Instance](https://next.api.alibabacloud.com/document/eds-aic/2023-09-30/DescribeAndroidInstances).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_phone_instance&exampleId=59386538-3f9b-0797-826f-161d7aeb05dd2cf4aaa2&activeTab=example&spm=docs.r.cloud_phone_instance.0.593865383f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cloud_phone_policy" "defaultjZ1gi0" {
}

resource "alicloud_cloud_phone_instance_group" "defaultYHMlTO" {
  instance_group_spec = "acp.basic.small"
  policy_group_id     = alicloud_cloud_phone_policy.defaultjZ1gi0.id
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
}


resource "alicloud_cloud_phone_instance" "default" {
  android_instance_group_id = alicloud_cloud_phone_instance_group.defaultYHMlTO.id
  android_instance_name     = "CreateInstanceName"
}
```

### Deleting `alicloud_cloud_phone_instance` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_phone_instance`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `android_instance_group_id` - (Optional, ForceNew) The ID of the instance group to which the instance belongs
* `android_instance_name` - (Optional) The instance name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Cloud Phone Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_phone_instance.example <id>
```