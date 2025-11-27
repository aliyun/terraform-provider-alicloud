---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_instance"
description: |-
  Provides a Alicloud KMS Instance resource.
---

# alicloud_kms_instance

Provides a KMS Instance resource.

For information about KMS Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/zh/key-management-service/latest/kms-instance-management).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Create a subscription kms instance

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_instance&exampleId=50f16f43-baa1-2b6c-6234-cf91838823c60b3149b2&activeTab=example&spm=docs.r.kms_instance.0.50f16f43ba&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = var.region
}
variable "region" {
  default = "cn-hangzhou"
}
variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "current" {}
resource "alicloud_vpc" "vpc-amp-instance-example" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id     = alicloud_vpc.vpc-amp-instance-example.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  vpc_id     = alicloud_vpc.vpc-amp-instance-example.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s3", var.name)
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s5", var.name)
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s7", var.name)
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_kms_instance" "default" {
  vpc_num         = "7"
  key_num         = "1000"
  secret_num      = "0"
  spec            = "1000"
  renew_status    = "ManualRenewal"
  product_version = "3"
  renew_period    = "3"
  vpc_id          = alicloud_vswitch.vswitch.vpc_id
  zone_ids        = ["cn-hangzhou-k", "cn-hangzhou-j"]
  vswitch_ids     = [alicloud_vswitch.vswitch-j.id]
  bind_vpcs {
    vpc_id       = alicloud_vswitch.shareVswitch.vpc_id
    region_id    = var.region
    vswitch_id   = alicloud_vswitch.shareVswitch.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  bind_vpcs {
    vpc_id       = alicloud_vswitch.share-vswitch2.vpc_id
    region_id    = var.region
    vswitch_id   = alicloud_vswitch.share-vswitch2.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  bind_vpcs {
    vpc_id       = alicloud_vswitch.share-vsw3.vpc_id
    region_id    = var.region
    vswitch_id   = alicloud_vswitch.share-vsw3.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  log          = "0"
  period       = "1"
  log_storage  = "0"
  payment_type = "Subscription"
}
```
Create a pay-as-you-go kms instance

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_instance&exampleId=20570abd-25c8-8f37-382f-db85438021960c9477b3&activeTab=example&spm=docs.r.kms_instance.1.20570abd25&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = var.region
}
variable "region" {
  default = "cn-hangzhou"
}
variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "current" {}
resource "alicloud_vpc" "vpc-amp-instance-example" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id     = alicloud_vpc.vpc-amp-instance-example.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  vpc_id     = alicloud_vpc.vpc-amp-instance-example.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s3", var.name)
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s5", var.name)
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s7", var.name)
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_kms_instance" "default" {
  payment_type                = "PayAsYouGo"
  product_version             = 3
  vpc_id                      = alicloud_vswitch.vswitch.vpc_id
  zone_ids                    = [alicloud_vswitch.vswitch.zone_id, alicloud_vswitch.vswitch-j.zone_id]
  vswitch_ids                 = [alicloud_vswitch.vswitch.id]
  force_delete_without_backup = true
  bind_vpcs {
    vpc_id       = alicloud_vswitch.shareVswitch.vpc_id
    region_id    = var.region
    vswitch_id   = alicloud_vswitch.shareVswitch.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  bind_vpcs {
    vpc_id       = alicloud_vswitch.share-vswitch2.vpc_id
    region_id    = var.region
    vswitch_id   = alicloud_vswitch.share-vswitch2.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  bind_vpcs {
    vpc_id       = alicloud_vswitch.share-vsw3.vpc_id
    region_id    = var.region
    vswitch_id   = alicloud_vswitch.share-vsw3.id
    vpc_owner_id = data.alicloud_account.current.id
  }
}
```

## Argument Reference

The following arguments are supported:
* `bind_vpcs` - (Optional, Set) Aucillary VPCs used to access this KMS instance See [`bind_vpcs`](#bind_vpcs) below.
* `force_delete_without_backup` - (Optional, Available since v1.223.2) Whether to force deletion even without backup.

-> **NOTE:** This parameter only takes effect when deletion is triggered.

* `instance_name` - (Optional, Computed) The name of the resource
* `key_num` - (Optional, Int) Maximum number of stored keys. The attribute is valid when the attribute `payment_type` is `Subscription`.
* `log` - (Optional, Computed) Instance Audit Log Switch. This attribute was limited to Subscription (prepaid) payment type before v1.264.0. As of v1.264.0, it is also supported for PayAsYouGo (postpaid) instances.
* `log_storage` - (Optional, Computed, Int) Instance log capacity. This attribute was limited to Subscription (prepaid) payment type before v1.264.0. As of v1.264.0, it is also supported for PayAsYouGo (postpaid) instances.
* `payment_type` - (Optional, ForceNew, Computed) The billing method. Valid values:

  - Subscription: the subscription billing method.
  - PayAsYouGo: the pay-as-you-go billing method.
* `period` - (Optional, Int) The subscription duration. Unit: month. The value must be an integral multiple of 12.

-> **NOTE:**   This parameter is required if you create a subscription instance.


-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `product_version` - (Optional, Computed) KMS Instance commodity type (software/hardware)
* `renew_period` - (Optional, Int) The auto-renewal period. Unit: month.

-> **NOTE:**   This parameter is required if the `RenewalStatus` parameter is set to `AutoRenewal`.

* `renew_status` - (Optional, Computed) The renewal status of the specified instance. Valid values:

  - AutoRenewal: The instance is automatically renewed.
  - ManualRenewal: The instance is manually renewed.
  - NotRenewal: The instance is not renewed.
* `renewal_period_unit` - (Optional, Available since v1.257.0) Automatic renewal period unit, value:
  - M: Month.
  - Y: Year.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `secret_num` - (Optional, Int) Maximum number of Secrets. The attribute is valid when the attribute `payment_type` is `Subscription`.
* `spec` - (Optional, Int) The computation performance level of the KMS instance. The attribute is valid when the attribute `payment_type` is `Subscription`.
* `tags` - (Optional, Map, Available since v1.259.0) The tag of the resource
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) that is associated with the KMS instance.
* `vpc_num` - (Optional, Int) The number of managed accesses. The maximum number of VPCs that can access this KMS instance. The attribute is valid when the attribute `payment_type` is `Subscription`.
* `vswitch_ids` - (Required, ForceNew, Set) Instance bind vswitches
* `zone_ids` - (Required, ForceNew, Set) zone id

### `bind_vpcs`

The bind_vpcs supports the following:
* `region_id` - (Optional) region id
* `vswitch_id` - (Optional) vswitch id
* `vpc_id` - (Optional) VPC ID
* `vpc_owner_id` - (Optional) VPC owner root user ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `ca_certificate_chain_pem` - KMS instance certificate chain in PEM format.
* `create_time` - The creation time of the resource.
* `end_date` - (Available since v1.233.1) Instance expiration time.
* `instance_name` - The name of the resource.
* `status` - Instance status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when create the Instance.
* `delete` - (Defaults to 10 mins) Used when delete the Instance.
* `update` - (Defaults to 60 mins) Used when update the Instance.

## Import

KMS Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_instance.example <id>
```