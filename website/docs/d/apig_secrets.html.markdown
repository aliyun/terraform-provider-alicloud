---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_secrets"
description: |-
  Provides a list of Apig Secrets to the user.
---

# alicloud_apig_secrets

This data source provides the Apig Secrets of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.294.0.

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
  description   = var.name
  kms_config {
    kms_instance_id = alicloud_kms_secret.default.dkms_instance_id
    kms_key_id      = alicloud_kms_secret.default.encryption_key_id
  }
}

data "alicloud_apig_secrets" "ids" {
  ids = [alicloud_apig_secret.default.id]
}

output "apig_secrets_id_0" {
  value = data.alicloud_apig_secrets.ids.secrets.0.id
}
```

## Argument Reference

The following attributes are exported:

* `ids` - (Optional, List) A list of Secret IDs.
* `name_regex` - (Optional) A regex string to filter results by Secret name.
* `gateway_type` - (Optional) The gateway type associated with the secret. Valid values: `AI`, `API`.
* `name_like` - (Optional) The secret name. Fuzzy match is supported.
* `status` - (Optional) The current status of the secret. Valid values: `ENABLE`, `DISABLE`, `DELETED`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Secret names.
* `secrets` - A list of Secrets. Each element contains the following attributes:
  * `id` - The ID of the resource supplied above.
  * `secret_id` - The secret ID.
  * `gateway_type` - The gateway type.
  * `name` - The secret name.
  * `description` - The secret description.
  * `secret_source` - The secret source.
  * `reference_count` - The reference count of the secret.
  * `status` - The status of the secret.
  * `create_timestamp` - The creation timestamp of the secret.
  * `update_timestamp` - The update timestamp of the secret.
  * `kms_config` - The KMS config of the secret.
    * `kms_instance_id` - The KMS instance ID.
    * `kms_key_id` - The KMS key ID.
    * `kms_secret_arn` - The KMS secret arn.
    * `version_id` - The version ID.
