---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_instances"
sidebar_current: "docs-alicloud-datasource-kms-instances"
description: |-
  Provides a list of Kms Instance owned by an Alibaba Cloud account.
---

# alicloud_kms_instances

This data source provides Kms Instance available to the user.[What is Instance](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available since v1.242.0.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
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
  vswitch_ids     = ["${alicloud_vswitch.vswitch.id}"]
  bind_vpcs {
    vpc_id       = alicloud_vswitch.shareVswitch.vpc_id
    region_id    = "cn-hangzhou"
    vswitch_id   = alicloud_vswitch.shareVswitch.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  bind_vpcs {
    vpc_id       = alicloud_vswitch.share-vswitch2.vpc_id
    region_id    = "cn-hangzhou"
    vswitch_id   = alicloud_vswitch.share-vswitch2.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  bind_vpcs {
    vpc_id       = alicloud_vswitch.share-vsw3.vpc_id
    region_id    = "cn-hangzhou"
    vswitch_id   = alicloud_vswitch.share-vsw3.id
    vpc_owner_id = data.alicloud_account.current.id
  }
  log          = "0"
  period       = "1"
  log_storage  = "0"
  payment_type = "Subscription"
}

data "alicloud_kms_instances" "default" {
  ids = ["${alicloud_kms_instance.default.id}"]
}

output "alicloud_kms_instance_example_id" {
  value = data.alicloud_kms_instances.default.instances.0.instance_id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Instance IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Instance IDs.
* `instances` - A list of Instance Entries. Each element contains the following attributes:
  * `instance_id` - The first ID of the resource
