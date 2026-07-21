---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_eip_instance_attachment"
description: |-
  Provides a Alicloud Ens Eip Instance Attachment resource.
---

# alicloud_ens_eip_instance_attachment

Provides a Ens Eip Instance Attachment resource.

Bind an EIP to an instance.

For information about Ens Eip Instance Attachment and how to use it, see [What is Eip Instance Attachment](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-associateenseipaddress).

-> **NOTE:** Available since v1.227.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_eip_instance_attachment&exampleId=d5589880-8fc2-8f8a-194f-9cb17dc318885c9c72c9&activeTab=example&spm=docs.r.ens_eip_instance_attachment.0.d55898808f&intl_lang=EN_US" target="_blank">
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

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_instance" "defaultXKjq1W" {
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy      = "Concentrate"
  schedule_area_level      = "Region"
  image_id                 = "centos_6_08_64_20G_alibase_20171208"
  payment_type             = "Subscription"
  instance_type            = "ens.sn1.stiny"
  password                 = "12345678abcABC"
  status                   = "Running"
  amount                   = "1"
  internet_charge_type     = "95BandwidthByMonth"
  instance_name            = var.name
  auto_use_coupon          = "true"
  instance_charge_strategy = "PriceHighPriority"
  ens_region_id            = var.ens_region_id
  period_unit              = "Month"
}

resource "alicloud_ens_eip" "defaultsGsN4e" {
  bandwidth            = "5"
  eip_name             = var.name
  ens_region_id        = var.ens_region_id
  internet_charge_type = "95BandwidthByMonth"
  payment_type         = "PayAsYouGo"
}

resource "alicloud_ens_eip_instance_attachment" "default" {
  instance_id   = alicloud_ens_instance.defaultXKjq1W.id
  allocation_id = alicloud_ens_eip.defaultsGsN4e.id
  instance_type = "EnsInstance"
  standby       = "false"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ens_eip_instance_attachment&spm=docs.r.ens_eip_instance_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `allocation_id` - (Required, ForceNew) The first ID of the resource
* `instance_id` - (Required, ForceNew) Instance ID
* `instance_type` - (Optional, ForceNew, Computed) The type of the EIP instance. Value:
  - `Nat`:NAT gateway.
  - `SlbInstance`: Server Load Balancer (ELB).
  - `NetworkInterface`: Secondary ENI.
  - `EnsInstance` (default): The ENS instance.
* `standby` - (Optional, ForceNew) Indicates whether the EIP is a backup EIP. Value:
  - true: Spare.
  - false: not standby.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<allocation_id>:<instance_id>:<instance_type>`.
* `status` - The status of the EIP.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Eip Instance Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Eip Instance Attachment.

## Import

Ens Eip Instance Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_eip_instance_attachment.example <allocation_id>:<instance_id>:<instance_type>
```