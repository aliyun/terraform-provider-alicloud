---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_instance_security_group_attachment"
description: |-
  Provides a Alicloud ENS Instance Security Group Attachment resource.
---

# alicloud_ens_instance_security_group_attachment

Provides a ENS Instance Security Group Attachment resource. Unbind instance and security group.

For information about ENS Instance Security Group Attachment and how to use it, see [What is Instance Security Group Attachment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.216.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_instance_security_group_attachment&exampleId=44e9b8d3-bb3b-b021-7483-94c140151d2d715f9997&activeTab=example&spm=docs.r.ens_instance_security_group_attachment.0.44e9b8d3bb&intl_lang=EN_US" target="_blank">
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