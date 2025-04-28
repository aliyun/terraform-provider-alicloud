---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_directory"
description: |-
  Provides a Alicloud Cloud SSO Directory resource.
---

# alicloud_cloud_sso_directory

Provides a Cloud SSO Directory resource.



For information about Cloud SSO Directory and how to use it, see [What is Directory](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-createdirectory).

-> **NOTE:** Available since v1.135.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_directory&exampleId=eaf9a8c8-5a0b-4a66-08d2-b8382e878edd34409e2e&activeTab=example&spm=docs.r.cloud_sso_directory.0.eaf9a8c85a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_cloud_sso_directory" "default" {
  directory_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `directory_global_access_status` - (Optional, Computed, Available since v1.248.0) Directory Global Acceleration activation status
* `directory_name` - (Optional) DirectoryName
* `login_preference` - (Optional, List, Available since v1.248.0) Login preferences See [`login_preference`](#login_preference) below.
* `mfa_authentication_status` - (Optional, Computed) MFA Authentication Status
* `mfa_authentication_setting_info` - (Optional, List, Available since v1.248.0) Global MFA verification configuration. See [`mfa_authentication_setting_info`](#mfa_authentication_setting_info) below.
* `password_policy` - (Optional, Computed, List, Available since v1.248.0) Password policy See [`password_policy`](#password_policy) below.
* `saml_identity_provider_configuration` - (Optional, Computed, List) Identity Provider (IDP) See [`saml_identity_provider_configuration`](#saml_identity_provider_configuration) below.
* `scim_synchronization_status` - (Optional, Computed) SCIM Synchronization Status
* `saml_service_provider` - (Optional, Computed, List, Available since v1.248.0) SP information. See [`saml_service_provider`](#saml_service_provider) below.
* `user_provisioning_configuration` - (Optional, List, Available since v1.248.0) User Provisioning configuration See [`user_provisioning_configuration`](#user_provisioning_configuration) below.

### `login_preference`

The login_preference supports the following:
* `allow_user_to_get_credentials` - (Optional, Computed, Available since v1.248.0) Whether the user can obtain the program access credential in the portal after logging in.
* `login_network_masks` - (Optional, Available since v1.248.0) IP address whitelist

### `mfa_authentication_setting_info`

The mfa_authentication_setting_info supports the following:
* `mfa_authentication_advance_settings` - (Optional, Computed, Available since v1.248.0) Global MFA validation policy
* `operation_for_risk_login` - (Optional, Computed, Available since v1.248.0) MFA verification policy for abnormal logon.

### `password_policy`

The password_policy supports the following:
* `max_login_attempts` - (Optional, Computed, Int, Available since v1.248.0) Number of password retries.
* `max_password_age` - (Optional, Computed, Int, Available since v1.248.0) Password validity period.
* `min_password_different_chars` - (Optional, Computed, Int, Available since v1.248.0) The minimum number of different characters in a password.
* `min_password_length` - (Optional, Computed, Int, Available since v1.248.0) Minimum password length.
* `password_not_contain_username` - (Optional, Computed, Available since v1.248.0) Whether the user name is not allowed in the password.
* `password_reuse_prevention` - (Optional, Computed, Int, Available since v1.248.0) Historical password check policy.

### `saml_identity_provider_configuration`

The saml_identity_provider_configuration supports the following:
* `binding_type` - (Optional, ForceNew, Computed, Available since v1.248.0) The Binding method for initiating a SAML request.
* `encoded_metadata_document` - (Optional, Computed) EncodedMetadataDocument
* `entity_id` - (Optional, ForceNew, Available since v1.248.0) EntityId
* `login_url` - (Optional, ForceNew, Available since v1.248.0) LoginUrl
* `sso_status` - (Optional, Computed) SSOStatus
* `want_request_signed` - (Optional, ForceNew, Available since v1.248.0) SP Request whether the signature is required

### `saml_service_provider`

The saml_service_provider supports the following:
* `authn_sign_algo` - (Optional, Computed, Available since v1.248.0) Signature algorithms supported by AuthNRequest
* `certificate_type` - (Optional, Computed, Available since v1.248.0) Type of certificate used for signing in the SSO process
* `support_encrypted_assertion` - (Optional, Computed, Available since v1.248.0) Whether IdP-side encryption of Assertion is supported.

### `user_provisioning_configuration`

The user_provisioning_configuration supports the following:
* `default_landing_page` - (Optional, Available since v1.248.0) The duration of the Session after the user logs in.
* `session_duration` - (Optional, Available since v1.248.0) The duration of the Session after the user logs in.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - CreateTime
* `password_policy` - Password policy
  * `hard_expire` - Whether to restrict login after Password Expiration
  * `max_password_length` - Maximum password length.
  * `require_lower_case_chars` - Whether lowercase letters are required in the password.
  * `require_numbers` - Whether numbers are required in the password.
  * `require_symbols` - Whether symbols are required in the password.
  * `require_upper_case_chars` - Whether uppercase letters are required in the password.
* `saml_identity_provider_configuration` - Identity Provider (IDP)
  * `certificate_ids` - Certificate ID list
  * `create_time` - CreateTime
  * `update_time` - UpdateTime
* `saml_service_provider` - SP information.
  * `acs_url` - ACS URL of SP.
  * `encoded_metadata_document` - SP metadata document (Base64 encoding).
  * `entity_id` - SP identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Directory.
* `delete` - (Defaults to 5 mins) Used when delete the Directory.
* `update` - (Defaults to 5 mins) Used when update the Directory.

## Import

Cloud SSO Directory can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_directory.example <id>
```