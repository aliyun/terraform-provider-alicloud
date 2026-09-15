---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_secret"
description: |-
  Provides a Alicloud APIG Secret resource.
---

# alicloud_apig_secret

Provides a APIG Secret resource.

API gateway secret.

For information about APIG Secret and how to use it, see [What is Secret](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateSecret).

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-j"
}

resource "alicloud_kms_instance" "default" {
  product_version = "3"
  vpc_num         = "1"
  key_num         = "1000"
  secret_num      = "1000"
  spec            = "1000"
  vpc_id          = data.alicloud_vpcs.default.ids.0
  vswitch_ids = [
    data.alicloud_vswitches.default.ids.0
  ]
  zone_ids = [
    "cn-hangzhou-k",
    "cn-hangzhou-j"
  ]
}

resource "alicloud_kms_key" "default" {
  dkms_instance_id       = alicloud_kms_instance.default.id
  pending_window_in_days = 7
}

resource "alicloud_kms_secret" "default" {
  secret_data                   = var.name
  secret_name                   = var.name
  version_id                    = "v1"
  dkms_instance_id              = alicloud_kms_key.default.dkms_instance_id
  encryption_key_id             = alicloud_kms_key.default.id
  force_delete_without_recovery = true
}

resource "alicloud_apig_secret" "default" {
  gateway_type  = "API"
  name          = var.name
  secret_source = "KMS"
  secret_data   = alicloud_kms_secret.default.secret_data
  kms_config {
    kms_instance_id = alicloud_kms_secret.default.dkms_instance_id
    kms_key_id      = alicloud_kms_secret.default.encryption_key_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The secret description.
* `gateway_type` - (Required, ForceNew) The gateway type.
* `kms_config` - (Required, ForceNew) The KMS config of the secret. See [`kms_config`](#kms_config) below.
* `name` - (Required, ForceNew) The secret name.
* `secret_data` - (Required) The KMS credential value.
* `secret_source` - (Required, ForceNew) The source of the key.

### `kms_config`

The kms_config supports the following.

* `kms_instance_id` - (Required, ForceNew) The KMS instance ID.
* `kms_key_id` - (Required, ForceNew) The KMS key ID.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Secret.
* `reference_count` - The reference count of the secret.
* `status` - The status of the secret.
* `create_timestamp` - The creation timestamp of the secret.
* `update_timestamp` - The update timestamp of the secret.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Secret.
* `update` - (Defaults to 5 mins) Used when update the Secret.
* `delete` - (Defaults to 5 mins) Used when delete the Secret.

## Import

APIG Secret can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_secret.example <secret_id>
```
