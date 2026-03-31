---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disk_default_kms_key_id"
description: |-
  Provides a Alicloud ECS Disk Default KMS Key ID resource.
---

# alicloud_ecs_disk_default_kms_key_id

Provides an ECS Disk Default KMS Key ID resource to configure the default KMS key used for account-level disk encryption.

For information about ECS Disk Default KMS Key ID and how to use it, see [What is Disk Default KMS Key ID](https://www.alibabacloud.com/help/en/doc-detail/59643.htm).

-> **NOTE:** Available since v1.274.0.

-> **NOTE:** This resource manages the default KMS key for account-level disk encryption in the current region.

-> **NOTE:** You need to enable ECS disk encryption by default first using `alicloud_ecs_disk_encryption_by_default`.

-> **NOTE:** You need to have KMS (Key Management Service) enabled and appropriate permissions configured.

## Example Usage

Basic Usage

```terraform
# Use existing VPC and VSwitch
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-i"
}

data "alicloud_vswitches" "default2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-j"
}

# Create KMS instance
resource "alicloud_kms_instance" "example" {
  product_version = "3"
  vpc_id          = data.alicloud_vpcs.default.ids.0
  zone_ids = [
    "cn-hangzhou-i",
    "cn-hangzhou-j"
  ]
  vswitch_ids = [
    data.alicloud_vswitches.default.ids.0,
    data.alicloud_vswitches.default2.ids.0
  ]
  vpc_num      = "1"
  key_num      = "1000"
  secret_num   = "100"
  spec         = "1000"
  payment_type = "PayAsYouGo"

  timeouts {
    delete = "20m"
  }
}

# Create a KMS key in the instance
resource "alicloud_kms_key" "example" {
  description            = "KMS key for ECS disk encryption"
  pending_window_in_days = 7
  key_usage              = "ENCRYPT/DECRYPT"
  key_spec               = "Aliyun_AES_256"
  dkms_instance_id       = alicloud_kms_instance.example.id

  timeouts {
    delete = "20m"
  }
}

# Enable ECS disk encryption by default first
resource "alicloud_ecs_disk_encryption_by_default" "example" {
  enabled = true
}

# Configure the default KMS key for disk encryption
resource "alicloud_ecs_disk_default_kms_key_id" "example" {
  kms_key_id = alicloud_kms_key.example.id
}
```

## Argument Reference

The following arguments are supported:

* `kms_key_id` - (Required) The KMS key ID used for ECS disk encryption by default. You can use the KMS key ID or the alias.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is the region ID where the resource is located.

## Import

ECS Disk Default KMS Key ID can be imported using the region id, e.g.

```shell
$ terraform import alicloud_ecs_disk_default_kms_key_id.example cn-hangzhou
```
