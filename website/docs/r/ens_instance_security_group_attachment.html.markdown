---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_instance_security_group_attachment"
description: |-
  Provides a Alicloud ENS Instance Security Group Attachment resource.
---

# alicloud_ens_instance_security_group_attachment

Provides a ENS Instance Security Group Attachment resource.

Bind instance and security group.

For information about ENS Instance Security Group Attachment and how to use it, see [What is Instance Security Group Attachment](https://next.api.alibabacloud.com/document/Ens/2017-11-10/JoinSecurityGroup).

-> **NOTE:** Available since v1.216.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ens_instance" "default" {
  system_disk {
    size = "20"
  }
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password                   = "12345678ABCabc"
  amount                     = "1"
  period                     = "1"
  internet_max_bandwidth_out = "10"
  public_ip_identification   = true
  ens_region_id              = "cn-chenzhou-telecom_unicom_cmcc"
  period_unit                = "Month"
}

resource "alicloud_ens_security_group" "default" {
  description         = "InstanceSecurityGroupAttachment_Description"
  security_group_name = var.name

}


resource "alicloud_ens_instance_security_group_attachment" "default" {
  instance_id       = alicloud_ens_instance.default.id
  security_group_id = alicloud_ens_security_group.default.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Optional, ForceNew, Computed) Instance ID.
* `security_group_id` - (Required, ForceNew) Security group ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<security_group_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance Security Group Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Instance Security Group Attachment.

## Import

ENS Instance Security Group Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_instance_security_group_attachment.example <instance_id>:<security_group_id>
```