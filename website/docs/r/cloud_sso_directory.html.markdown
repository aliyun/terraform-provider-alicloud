---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_directory"
sidebar_current: "docs-alicloud-resource-cloud-sso-directory"
description: |-
  Provides a Alicloud Cloud SSO Directory resource.
---

# alicloud_cloud_sso_directory

Provides a Cloud SSO Directory resource.

For information about Cloud SSO Directory and how to use it, see [What is Directory](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-createdirectory).

-> **NOTE:** Available since v1.135.0.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_directory&exampleId=df1d7142-f49f-ab3b-ba90-79a4ddbcb4422c032e32&activeTab=example&spm=docs.r.cloud_sso_directory.0.df1d7142f4&intl_lang=EN_US" target="_blank">
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
data "alicloud_cloud_sso_directories" "default" {}

resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `directory_name` - (Optional, Sensitive) The name of the CloudSSO directory. The length is 2-64 characters, and it can contain lowercase letters, numbers, and dashes (-). It cannot start or end with a dash and cannot have two consecutive dashes. Need to be globally unique, and capitalization is not supported. Cannot start with `d-`.
* `mfa_authentication_status` - (Optional) The mfa authentication status. Valid values: `Enabled` or `Disabled`. Default to `Enabled`.
* `scim_synchronization_status` - (Optional) The scim synchronization status. Valid values: `Enabled` or `Disabled`. Default to `Disabled`.
* `saml_identity_provider_configuration` - (Optional, ForceNew) The saml identity provider configuration. See [`saml_identity_provider_configuration`](#saml_identity_provider_configuration) below.

-> **NOTE:** The `saml_identity_provider_configuration` will be removed automatically when the resource is deleted, please operate with caution. If there are left more configuration in the directory, please remove them before deleting the directory.

### `saml_identity_provider_configuration`

The saml_identity_provider_configuration supports the following:

* `encoded_metadata_document` - (Optional, Sensitive) Base64 encoded IdP metadata document. **NOTE:** If the IdP Metadata has been uploaded, no update will be made if this parameter is not specified, otherwise the update will be made according to the parameter content. If IdP Metadata has not been uploaded, and the parameter `sso_status` is `Enabled`, this parameter must be provided. If the IdP Metadata has not been uploaded, and the parameter `sso_status` is `Disabled`, this parameter can be omitted, and the IdP Metadata will remain empty.
* `sso_status` - (Optional) SAML SSO login enabled status. Valid values: `Enabled` or `Disabled`. Default to `Disabled`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Directory.

## Import

Cloud SSO Directory can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_directory.example <id>
```
