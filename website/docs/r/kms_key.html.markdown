---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_key"
description: |-
  Provides a Alicloud KMS Key resource.
---

# alicloud_kms_key

Provides a KMS Key resource. 

For information about KMS Key and how to use it, see [What is Key](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.223.0.

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

resource "alicloud_vpc" "create-vpc" {
  vpc_name = var.name
}

resource "alicloud_vswitch" "vswitch-k" {
  vpc_id     = alicloud_vpc.create-vpc.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch-j" {
  vpc_id     = alicloud_vpc.create-vpc.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_kms_instance" "create-instance" {
  vpc_num         = "1"
  key_num         = "1000"
  product_type    = "kms_ddi_public_cn"
  secret_num      = "100"
  product_version = "3"
  renew_status    = "AutoRenewal"
  vpc_id          = alicloud_vpc.create-vpc.id
  vswitch_ids     = ["${alicloud_vswitch.vswitch-j.id}"]
  zone_ids        = ["${alicloud_vswitch.vswitch-j.zone_id}", "${alicloud_vswitch.vswitch-k.zone_id}"]
  spec            = "1000"
  renew_period    = "1"
}


resource "alicloud_kms_key" "default" {
  protection_level          = "SOFTWARE"
  description               = "key description example"
  key_spec                  = "Aliyun_AES_256"
  key_usage                 = "ENCRYPT/DECRYPT"
  rotation_interval         = "604800s"
  dkms_instance_id          = alicloud_kms_instance.create-instance.id
  origin                    = "Aliyun_KMS"
  policy                    = "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1192853035118460:*\"]},\"Sid\":\"kms default key policy\"}]}"
  enable_automatic_rotation = true
  automatic_rotation        = "Enabled"
  status                    = "Enabled"
}
```

### Deleting `alicloud_kms_key` or removing it from your configuration

Terraform cannot destroy resource `alicloud_kms_key`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `automatic_rotation` - (Optional, Available since v1.85.0) Whether automatic key rotation is enabled. Valid values:
  - Enabled: Auto rotation is on.
  - Disabled: Auto rotation is not turned on.
  - Suspended: Automatic rotation is Suspended.
* `certificate_id` - (Optional) The ID of the certificate.
-> **NOTE:**  Only one of the KeyId, SecretName, and CertificateId parameters must be specified.
* `description` - (Optional, Available since v1.85.0) The description of the key.
* `dkms_instance_id` - (Optional, ForceNew, Available since v1.85.0) The ID of the KMS instance.
* `enable_automatic_rotation` - (Optional) Whether to enable automatic key rotation. Value:
  - true: On
  - false (default): Not enabled

This parameter value is valid only when the key management type to which the key belongs supports automatic rotation. For details, see [Key Rotation](~~ 2358146 ~~).
.
* `enable_deletion_protection` - (Optional) Whether to enable deletion protection, value:
  - true: Enable deletion protection.
  - false (default): Delete protection is turned off.
* `key_spec` - (Optional, ForceNew, Available since v1.85.0) The specifications of the key.
* `key_usage` - (Optional, ForceNew, Available since v1.85.0) The purpose of the key.
* `origin` - (Optional, ForceNew, Available since v1.85.0) Key material source. .
* `pending_window_in_days` - (Optional, Available since v1.85.0) Key pre-deletion cycle. During this period of time, you can revoke and delete the key that is in the pending state. After the pre-delete time, you cannot revoke the deletion. Value range: 7~366. Unit: Days.
* `policy` - (Optional) Policy.
* `policy_name` - (Optional) Policy Name.
* `protection_level` - (Optional, ForceNew, Available since v1.85.0) The protection level of the key.
* `rotation_interval` - (Optional, Available since v1.85.0) The period of automatic key rotation. The unit is seconds, and the format is an integer value followed by the character s. For example, a 7-day cycle is 604800s.  This value is returned only if the AutomaticRotation parameter value is Enabled or Suspended.
* `secret_name` - (Optional) Credential name.
-> **NOTE:**  KeyId, SecretName, and cerificateid must and can only specify one of the parameters.

.
* `status` - (Optional, Computed, Available since v1.85.0) The state of the key. For more information, see [the impact of the status of the user's master key on API calls](~~ 44211 ~~).
* `tags` - (Optional, ForceNew, Map, Available since v1.85.0) Label.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Key.
* `update` - (Defaults to 5 mins) Used when update the Key.

## Import

KMS Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_key.example <id>
```