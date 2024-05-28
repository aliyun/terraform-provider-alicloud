---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_secret"
description: |-
  Provides a Alicloud KMS Secret resource.
---

# alicloud_kms_secret

Provides a KMS Secret resource. 

For information about KMS Secret and how to use it, see [What is Secret](https://www.alibabacloud.com/help/en/).

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


resource "alicloud_kms_secret" "default" {
  version_id                = "v1"
  secret_name               = var.name
  description               = "test"
  secret_type               = "ECS"
  enable_automatic_rotation = true
  secret_data_type          = "text"
  secret_data               = "{\"UserName\":\"ecs-user1423\",\"Password\":\"ecs-user1423\"}"
  policy                    = "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1117600963847258:*\"]},\"Sid\":\"kms default secret policy\"}]}"
  encryption_key_id         = "key-phzz661646aeygbqrfqg5j"
  dkms_instance_id          = "kst-phzz61dbabacquz826y29f"
  rotation_interval         = "604800s"
  extended_config           = "{\"SecretSubType\":\"Password\",\"InstanceId\":\"i-bp1fr0a100hd0yf6ksq3\",\"CustomData\":{\"Key1\":\"v1\",\"fds\":\"fdsf\"},\"CommandId\":\"cmd-ACS-KMS-RotateECSSecret-For-Linux.sh\",\"RegionId\":\"cn-hangzhou\"}"
}
```

## Argument Reference

The following arguments are supported:
* `certificate_id` - (Optional) The ID of the certificate.
-> **NOTE:**  KeyId, SecretName, and cerificateid must and can only specify one of the parameters.
* `custom_data` - (Optional) Expand the custom data in the configuration.
* `description` - (Optional, Available since v1.76.0) Description information for the credential.
* `dkms_instance_id` - (Optional, ForceNew, Available since v1.76.0) The ID of the dedicated KMS instance.
* `enable_automatic_rotation` - (Optional, Available since v1.76.0) Whether to enable automatic rotation, value:
  - true: Turns on automatic rotation.
  - false (default): does not turn on automatic rotation.
-> **NOTE:**  This parameter is valid when the value of SecretType is Rds(RDS credentials), RAMCredentials(RAM credentials), or ECS(ECS credentials). When the value of SecretType is Generic, automatic rotation is not supported. You can manually rotate by using the [PutSecretValue](~~ 154503 ~~) operation.
* `encryption_key_id` - (Optional, ForceNew, Available since v1.76.0) The identifier of the KMS user Master key used to encrypt the protection credential value.
* `extended_config` - (Optional, ForceNew, Available since v1.76.0) Expanding configuration of credentials.
-> **NOTE:**  only managed RDS credentials, managed RAM credentials, or managed ECS credentials return this parameter.
* `extended_config_custom_data` - (Optional, Map) Expand the custom data in the configuration.
-> **NOTE:** - If this parameter is specified, the existing expansion configuration of the credential will be updated.
-> **NOTE:** - Normal credentials do not support Setting this parameter.
* `key_id` - (Optional) The key ID. The globally unique identifier of The Master Key (CMK).
-> **NOTE:**  KeyId, SecretName, and cerificateid must and can only specify one of the parameters.
* `move_to_version` - (Optional) Use the specified version status to mark the version specified by this parameter.
-> **NOTE:**  - RemoveFromVersion and moveoverversion specify at least one of the parameters.
-> **NOTE:** - This parameter must be specified when the value of VersionStage is ACSCurrent or acspreviewed.
* `policy` - (Optional) Policy.
* `policy_name` - (Optional) PolicyName.
* `recovery_window_in_days` - (Optional, Available since v1.76.0) Delete credentials as recoverable and specify the recoverable window (number of days). Default value: 30.
* `remove_from_version` - (Optional) Removes the specified version status from the version specified by this parameter.
-> **NOTE:**  RemoveFromVersion and moveoverversion specify at least one of the parameters.
* `rotation_interval` - (Optional, ForceNew, Available since v1.76.0) The period of automatic rotation of credentials.
The format is 'integer[unit]', where 'integer' represents the length of time and 'unit' represents the unit of time. 'unit' value: s (seconds). For example, a 7-day cycle is 604800s.
-> **NOTE:**  This parameter is returned when automatic rotation is turned on.
* `secret_data` - (Required, ForceNew, Available since v1.76.0) The credential value. The length does not exceed 30720 bytes (30KB). KMS uses the specified key to encrypt it and stores it in the initial version.
  - When the value of SecretType is Generic, you can customize the credential value.
  - When the value of SecretType is Rds(RDS credentials), the credential value format is '{"Accounts":[{"AccountName":"","AccountPassword":""}]}'. Where 'AccountName' is the account name of the RDS instance, and 'AccountPassword' is the account password of the RDS instance.
  - When the value of SecretType is RAMCredentials(RAM credentials), the credential value format is '{"AccessKeys":[{"AccessKeyId":"","AccessKeySecret":""}]}'. Where 'AccessKeyId' is the access key ID, and 'accesskeysecret 'is the access key content. You must specify all the accesskeys of the RAM user.
  - When the value of SecretType is ECS(ECS credentials), the credential value format is:
  - When the value of SecretSubType in ExtendedConfig parameter is Password: '{"UserName": "","Password": ""}'. Where 'UserName' is the username used to log on to the ECS instance, and 'Password' is the password used to log on to the ECS instance.
  - When the value of SecretSubType in ExtendedConfig parameter is SSHKey: '{"UserName": "", "PublicKey": "", "PrivateKey": ""}'. Where 'PublicKey' is the SSH format public key used to log on to the ECS instance, and 'PrivateKey' is the private key used to log on to the ECS instance.
* `secret_data_type` - (Optional, ForceNew, Computed, Available since v1.76.0) Credential value type. Value:
  - text (default): text type
  - binary: binary type
-> **NOTE:**  When the value of SecretType is Rds, RAMCredentials, or ECS, the value of SecretDataType can only be text.
* `secret_name` - (Required, ForceNew, Available since v1.76.0) The credential name.
* `secret_type` - (Optional, ForceNew, Available since v1.76.0) The credential type. Value:
  - Generic: Normal credentials.
  - Rds: managed RDS credentials.
  - RAMCredentials: managed RAM credentials.
  - ECS: the managed ECS credentials.
* `tags` - (Optional, ForceNew, Map, Available since v1.76.0) The resource label of the credential. If FetchTags is false or is not specified, this parameter is not returned.
* `version_id` - (Required, ForceNew, Available since v1.76.0) The credential version number.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the credential was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Secret.
* `delete` - (Defaults to 5 mins) Used when delete the Secret.
* `update` - (Defaults to 5 mins) Used when update the Secret.

## Import

KMS Secret can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_secret.example <id>
```