---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_secret"
sidebar_current: "docs-alicloud-resource-kms-secret"
description: |-
  Provides a Alicloud KMS Secret resource.
---

# alicloud_kms_secret

Provides a KMS Secret resource.

For information about KMS Secret and how to use it, see [What is Secret](https://www.alibabacloud.com/help/en/kms/developer-reference/api-createsecret).

-> **NOTE:** Available since v1.76.0.

-> **NOTE:** Starting from KMS 3.0, secrets are created inside a dedicated KMS instance. When you create a secret in KMS 3.0, you must specify `dkms_instance_id` and `encryption_key_id`. The examples below show how to provide them. If your account still uses the legacy shared KMS and does not require a dedicated instance, these two parameters can be omitted, but this legacy behavior is no longer recommended.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_secret&exampleId=9d3c7df0-2204-7685-a48a-90c820e68b66169212e7&activeTab=example&spm=docs.r.kms_secret.0.9d3c7df022&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

# Replace the IDs with your own dedicated KMS instance and key.
data "alicloud_kms_instances" "default" {
  ids = ["kst-bjj62d8f5e0sgtx8h****"]
}

data "alicloud_kms_keys" "default" {
  ids = ["key-gzz63ff0db5hg3qje****"]
}

resource "alicloud_kms_secret" "default" {
  secret_name                   = var.name
  secret_data                   = "Secret data"
  version_id                    = "v1"
  dkms_instance_id              = data.alicloud_kms_instances.default.instances.0.instance_id
  encryption_key_id             = data.alicloud_kms_keys.default.keys.0.id
  force_delete_without_recovery = true
}
```

### Create a generic secret

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_kms_instances" "default" {
  ids = ["kst-bjj62d8f5e0sgtx8h****"]
}

data "alicloud_kms_keys" "default" {
  ids = ["key-gzz63ff0db5hg3qje****"]
}

resource "alicloud_kms_secret" "generic" {
  secret_name                   = var.name
  secret_type                   = "Generic"
  secret_data                   = "Secret data"
  version_id                    = "v1"
  secret_data_type              = "text"
  dkms_instance_id              = data.alicloud_kms_instances.default.instances.0.instance_id
  encryption_key_id             = data.alicloud_kms_keys.default.keys.0.id
  force_delete_without_recovery = true
}
```

### Create an RDS secret

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_kms_instances" "default" {
  ids = ["kst-bjj62d8f5e0sgtx8h****"]
}

data "alicloud_kms_keys" "default" {
  ids = ["key-gzz63ff0db5hg3qje****"]
}

resource "alicloud_kms_secret" "rds" {
  secret_name = var.name
  secret_type = "Rds"
  secret_data = jsonencode({
    Accounts = [
      {
        AccountName     = "rds_user"
        AccountPassword = "YourPassword12345!"
      }
    ]
  })
  version_id        = "v1"
  dkms_instance_id  = data.alicloud_kms_instances.default.instances.0.instance_id
  encryption_key_id = data.alicloud_kms_keys.default.keys.0.id
  extended_config = jsonencode({
    SecretSubType = "SingleUser"
    DBInstanceId  = "rm-bp1b3dd3a506e****"
  })
  force_delete_without_recovery = true
}
```

### Create a Redis/Tair secret

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_kms_instances" "default" {
  ids = ["kst-bjj62d8f5e0sgtx8h****"]
}

data "alicloud_kms_keys" "default" {
  ids = ["key-gzz63ff0db5hg3qje****"]
}

resource "alicloud_kms_secret" "redis" {
  secret_name       = var.name
  secret_type       = "Redis"
  secret_data       = "$Auto"
  version_id        = "v1"
  dkms_instance_id  = data.alicloud_kms_instances.default.instances.0.instance_id
  encryption_key_id = data.alicloud_kms_keys.default.keys.0.id
  extended_config = jsonencode({
    SecretSubType    = "DoubleUsers"
    AccountName      = "redis_user"
    CloneAccountName = "redis_user_clone"
    AccountPrivilege = "ROLE_READ_ONLY"
    InstanceId       = "r-bp1b3dd3a506e****"
    RegionId         = "cn-hangzhou"
  })
  force_delete_without_recovery = true
}
```

### Create a RAM credentials secret

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_kms_instances" "default" {
  ids = ["kst-bjj62d8f5e0sgtx8h****"]
}

data "alicloud_kms_keys" "default" {
  ids = ["key-gzz63ff0db5hg3qje****"]
}

resource "alicloud_kms_secret" "ram" {
  secret_name = "$Auto"
  secret_type = "RAMCredentials"
  secret_data = jsonencode({
    AccessKeys = [
      {
        AccessKeyId     = "LTAI********************"
        AccessKeySecret = "YourAccessKeySecret"
      }
    ]
  })
  version_id        = "v1"
  dkms_instance_id  = data.alicloud_kms_instances.default.instances.0.instance_id
  encryption_key_id = data.alicloud_kms_keys.default.keys.0.id
  extended_config = jsonencode({
    SecretSubType = "RamUserAccessKey"
    UserName      = "ram_user"
  })
  force_delete_without_recovery = true
}
```

### Create an ECS secret

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_kms_instances" "default" {
  ids = ["kst-bjj62d8f5e0sgtx8h****"]
}

data "alicloud_kms_keys" "default" {
  ids = ["key-gzz63ff0db5hg3qje****"]
}

resource "alicloud_kms_secret" "ecs" {
  secret_name = "acs/ecs/${var.name}"
  secret_type = "ECS"
  secret_data = jsonencode({
    UserName = "root"
    Password = "YourPassword12345!"
  })
  version_id        = "v1"
  dkms_instance_id  = data.alicloud_kms_instances.default.instances.0.instance_id
  encryption_key_id = data.alicloud_kms_keys.default.keys.0.id
  extended_config = jsonencode({
    SecretSubType = "Password"
    RegionId      = "cn-hangzhou"
    InstanceId    = "i-bp1b3dd3a506e****"
  })
  force_delete_without_recovery = true
}
```

### Create a PolarDB secret

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_kms_instances" "default" {
  ids = ["kst-bjj62d8f5e0sgtx8h****"]
}

data "alicloud_kms_keys" "default" {
  ids = ["key-gzz63ff0db5hg3qje****"]
}

resource "alicloud_kms_secret" "polardb" {
  secret_name       = var.name
  secret_type       = "PolarDB"
  secret_data       = "$Auto"
  version_id        = "v1"
  dkms_instance_id  = data.alicloud_kms_instances.default.instances.0.instance_id
  encryption_key_id = data.alicloud_kms_keys.default.keys.0.id
  extended_config = jsonencode({
    SecretSubType    = "DoubleUsers"
    RegionId         = "cn-hangzhou"
    DBClusterId      = "pc-bp1b3dd3a506e****"
    DBType           = "MySQL"
    AccountName      = "polardb_user"
    CloneAccountName = "polardb_user_clone"
    AccountType      = "Normal"
    AccountPrivilege = "ReadOnly"
    DBName           = "testdb"
  })
  force_delete_without_recovery = true
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_kms_secret&spm=docs.r.kms_secret.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, ForceNew) The name of the secret. For more information, see [How to use it](https://www.alibabacloud.com/help/en/key-management-service/latest/kms-createsecret).
* `secret_data` - (Required) The data of the secret. **NOTE:** From version 1.204.1, `secret_data` updating diff will be ignored when `secret_type` is not `Generic`. The expected format of `secret_data` depends on `secret_type`:
  - `Generic`: any custom string.
  - `Rds`: a JSON string in the form `{"Accounts":[{"AccountName":"<rds-account-name>","AccountPassword":"<rds-account-password>"}]}`.
  - `RAMCredentials`: a JSON string in the form `{"AccessKeys":[{"AccessKeyId":"<access-key-id>","AccessKeySecret":"<access-key-secret>"}]}`.
  - `ECS`: a JSON string. When `extended_config.SecretSubType` is `Password`, use `{"UserName":"<user-name>","Password":"<password>"}`. When `SecretSubType` is `SSHKey`, use `{"UserName":"<user-name>","PublicKey":"<public-key>","PrivateKey":"<private-key>"}`.
  - `Redis` and `PolarDB`: use the literal string `$Auto`.
* `version_id` - (Required) The version number of the initial version.
* `secret_type` - (Optional, ForceNew, Available since v1.204.1) The type of the secret. Valid values:
  - `Generic`: Generic secret.
  - `Rds`: ApsaraDB RDS secret.
  - `Redis`: (Available since v1.253.0) ApsaraDB for Redis secret.
  - `RAMCredentials`: RAM secret.
  - `ECS`: ECS secret.
  - `PolarDB`: (Available since v1.253.0) PolarDB secret.
* `secret_data_type` - (Optional) The type of the secret value. Default value: `text`. Valid values: `text`, `binary`.
* `encryption_key_id` - (Optional, ForceNew) The ID of the KMS key used to encrypt the secret value. The key and the secret must belong to the same KMS instance, and the key must be a symmetric key. **NOTE:** In KMS 3.0, this parameter is required.
* `dkms_instance_id` - (Optional, ForceNew, Available since v1.183.0) The ID of the KMS instance. **NOTE:** In KMS 3.0, this parameter is required.
* `extended_config` - (Optional, ForceNew, Available since v1.204.1) The extended configuration of the secret. This parameter is required when `secret_type` is `Rds`, `Redis`, `RAMCredentials`, `ECS` or `PolarDB`. The value is a JSON string. For more information, see [How to use it](https://www.alibabacloud.com/help/en/key-management-service/latest/kms-createsecret).
* `enable_automatic_rotation` - (Optional, Bool, Available since v1.124.0) Specifies whether to enable automatic rotation. Default value: `false`. Valid values: `true`, `false`.
* `rotation_interval` - (Optional, Available since v1.124.0) The interval for automatic rotation. For more information, see [How to use it](https://www.alibabacloud.com/help/en/key-management-service/latest/kms-createsecret).
* `policy` - (Optional, Available since v1.224.0) The content of the secret policy. The value is in the JSON format. The value can be up to 32,768 bytes in length. For more information, see [How to use it](https://www.alibabacloud.com/help/en/kms/developer-reference/api-setsecretpolicy).
* `description` - (Optional) The description of the secret.
* `force_delete_without_recovery` - (Optional, Bool) Specifies whether to immediately delete a secret. Default value: `false`. Valid values: `true`, `false`.
* `recovery_window_in_days` - (Optional, Int) Specifies the recovery period of the secret if you do not forcibly delete it. Unit: Days. Default value: `30`. Valid values: `7` to `30`. **NOTE:**  If `force_delete_without_recovery` is set to `true`, `recovery_window_in_days` will be ignored.
* `version_stages` - (Optional, List) The stage label that is used to mark the new version.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Secret.
* `arn` - The ARN of the secret.
* `create_time` - (Available since v1.224.0) The time when the secret is created.
* `planned_delete_time` - The time when the secret is scheduled to be deleted.

## Timeouts

-> **NOTE:** Available since v1.103.2.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Secret.
* `update` - (Defaults to 5 mins) Used when update the Secret.
* `delete` - (Defaults to 5 mins) Used when delete the Secret. 

## Import

KMS Secret can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_secret.example <id>
```
