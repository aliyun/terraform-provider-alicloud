---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_disk"
description: |-
  Provides a Alicloud Simple Application Server Disk resource.
---

# alicloud_simple_application_server_disk

Provides a Simple Application Server Disk resource.



For information about Simple Application Server Disk and how to use it, see [What is Disk](https://next.api.alibabacloud.com/document/SWAS-OPEN/2020-06-01/CreateDisk).

-> **NOTE:** Available since v1.273.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_simple_application_server_instance" "defaultV70JQf" {
  instance_name     = "examplewujie"
  status            = "Running"
  plan_id           = "swas.s1.c2m2s50b3"
  image_id          = "21e9617bd4754f77a090d2fbc94916a4"
  period            = "1"
  data_disk_size    = "0"
  password          = "@3612568Wj"
  payment_type      = "Subscription"
  auto_renew        = true
  auto_renew_period = "1"
}


resource "alicloud_simple_application_server_disk" "default" {
  disk_size   = "20"
  instance_id = alicloud_simple_application_server_instance.defaultV70JQf.id
  remark      = "example"
}
```

### Deleting `alicloud_simple_application_server_disk` or removing it from your configuration

Terraform cannot destroy resource `alicloud_simple_application_server_disk`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `disk_size` - (Required, ForceNew, Int) disk size
* `instance_id` - (Required, ForceNew) instance id
* `remark` - (Optional) Note information.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - The creation time of the resource.
* `disk_name` - The name of the resource.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Disk.
* `update` - (Defaults to 5 mins) Used when update the Disk.

## Import

Simple Application Server Disk can be imported using the id, e.g.

```shell
$ terraform import alicloud_simple_application_server_disk.example <disk_id>
```