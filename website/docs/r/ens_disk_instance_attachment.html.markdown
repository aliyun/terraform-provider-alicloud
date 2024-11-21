---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_disk_instance_attachment"
description: |-
  Provides a Alicloud ENS Disk Instance Attachment resource.
---

# alicloud_ens_disk_instance_attachment

Provides a ENS Disk Instance Attachment resource. Disk instance mount.

For information about ENS Disk Instance Attachment and how to use it, see [What is Disk Instance Attachment](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-attachdisk).

-> **NOTE:** Available since v1.216.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_disk_instance_attachment&exampleId=abdcd0b9-2982-4e75-7e83-82240eccf33e9464c87f&activeTab=example&spm=docs.r.ens_disk_instance_attachment.0.abdcd0b929&intl_lang=EN_US" target="_blank">
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

resource "alicloud_ens_disk" "default" {
  size          = "20"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  payment_type  = "PayAsYouGo"
  category      = "cloud_efficiency"
}

resource "alicloud_ens_instance" "default" {
  system_disk {
    size = "20"
  }
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password                   = "12345678ABCabc"
  amount                     = "1"
  internet_max_bandwidth_out = "10"
  unique_suffix              = true
  public_ip_identification   = true
  ens_region_id              = "cn-chenzhou-telecom_unicom_cmcc"
  schedule_area_level        = "Region"
  period_unit                = "Month"
  period                     = "1"
  timeouts {
    delete = "50m"
  }
}


resource "alicloud_ens_disk_instance_attachment" "default" {
  instance_id          = alicloud_ens_instance.default.id
  delete_with_instance = "false"
  disk_id              = alicloud_ens_disk.default.id
}
```

## Argument Reference

The following arguments are supported:
* `delete_with_instance` - (Optional) Whether the cloud disk to be mounted is released with the instance  Value: true: When the instance is released, the cloud disk is released together with the instance. false: When the instance is released, the cloud disk is retained and is not released together with the instance. Empty means false by default.
* `disk_id` - (Required, ForceNew) The ID of the cloud disk to be mounted. The Cloud Disk (DiskId) and the instance (InstanceId) must be on the same node.
* `instance_id` - (Required, ForceNew) Instance ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<disk_id>:<instance_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Disk Instance Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Disk Instance Attachment.

## Import

ENS Disk Instance Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_disk_instance_attachment.example <disk_id>:<instance_id>
```