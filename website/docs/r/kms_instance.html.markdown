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

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "tf-exampleacc-kms-instance"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "default" {
  count      = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
  cidr_block = "172.16.0.0/12"
  vpc_name   = "tf-exampleacc-kms-instance"
}

data "alicloud_vswitches" "vswitch" {
  vpc_id  = local.vpc_id
  zone_id = "cn-hangzhou-k"
}

data "alicloud_vswitches" "vswitch-j" {
  vpc_id  = local.vpc_id
  zone_id = "cn-hangzhou-j"
}

locals {
  vpc_id  = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids.0 : concat(alicloud_vpc.default[*].id, [""])[0]
  vsw_id  = length(data.alicloud_vswitches.vswitch.ids) > 0 ? data.alicloud_vswitches.vswitch.ids.0 : concat(alicloud_vswitch.vswitch[*].id, [""])[0]
  vswj_id = length(data.alicloud_vswitches.vswitch-j.ids) > 0 ? data.alicloud_vswitches.vswitch-j.ids.0 : concat(alicloud_vswitch.vswitch-j[*].id, [""])[0]
}

resource "alicloud_vswitch" "vswitch" {
  count      = length(data.alicloud_vswitches.vswitch.ids) > 0 ? 0 : 1
  vpc_id     = local.vpc_id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vswitch-j" {
  count      = length(data.alicloud_vswitches.vswitch-j.ids) > 0 ? 0 : 1
  vpc_id     = local.vpc_id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "shareVPC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s1", var.name)
}

resource "alicloud_vswitch" "shareVswitch" {
  vpc_id     = alicloud_vpc.shareVPC.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s2", var.name)
}

resource "alicloud_vswitch" "share-vswitch2" {
  vpc_id     = alicloud_vpc.share-VPC2.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vpc" "share-VPC3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = format("%s3", var.name)
}

resource "alicloud_vswitch" "share-vsw3" {
  vpc_id     = alicloud_vpc.share-VPC3.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_kms_instance" "default" {
  product_version = "3"
  vpc_id          = local.vpc_id
  zone_ids = [
    "cn-hangzhou-k",
    "cn-hangzhou-j"
  ]
  vswitch_ids = [
    "${local.vsw_id}"
  ]
  vpc_num    = "1"
  key_num    = "1000"
  secret_num = "0"
  spec       = "1000"
}

# Save Instance's CA certificate chain to a local file
# resource "local_file" "ca_certificate_chain_pem" {
#   content  = alicloud_kms_instance.default.ca_certificate_chain_pem
#   filename = "ca.pem"
# }
```

### Deleting `alicloud_kms_instance` or removing it from your configuration

Terraform cannot destroy resource `alicloud_kms_instance`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `bind_vpcs` - (Optional) Aucillary VPCs used to access this KMS instance. See [`bind_vpcs`](#bind_vpcs) below.
* `key_num` - (Required) Maximum number of stored keys.
* `product_version` - (Optional) KMS Instance commodity type (software/hardware). Currently, only version 3 is supported.
* `renew_period` - (Optional) Automatic renewal period, in months.
* `renew_status` - (Optional) Renewal options (manual renewal, automatic renewal, no renewal).
* `secret_num` - (Required) Maximum number of Secrets.
* `spec` - (Required) The computation performance level of the KMS instance.
* `vpc_id` - (Required, ForceNew) Instance VPC id.
* `vpc_num` - (Required) The number of managed accesses. The maximum number of VPCs that can access this KMS instance.
* `vswitch_ids` - (Required, ForceNew) Instance bind vswitches.
* `zone_ids` - (Required, ForceNew) zone id.

### `bind_vpcs`

The bind_vpcs supports the following:
* `region_id` - (Optional) region id.
* `vswitch_id` - (Optional) vswitch id.
* `vpc_id` - (Optional) VPC ID.
* `vpc_owner_id` - (Optional) VPC owner root user ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `ca_certificate_chain_pem` - KMS instance certificate chain in PEM format.
* `create_time` - The creation time of the resource.
* `instance_name` - The name of the resource.
* `status` - Instance status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when create the Instance.
* `update` - (Defaults to 60 mins) Used when update the Instance.

## Import

KMS Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_instance.example <id>
```