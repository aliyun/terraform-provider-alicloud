---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_directory"
sidebar_current: "docs-alicloud-resource-cloud-sso-directory"
description: |-
  Provides a Alicloud Cloud SSO Directory resource.
---

# alicloud\_cloud\_sso\_directory

Provides a Cloud SSO Directory resource.

For information about Cloud SSO Directory and how to use it, see [What is Directory](https://www.alibabacloud.com/help/doc-detail/263624.html).

-> **NOTE:** Available in v1.135.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_sso_directory" "default" {
  directory_name = "example-value"
}

```

## Argument Reference

The following arguments are supported:

* `directory_name` - (Optional, Sensitive) The name of the CloudSSO directory. The length is 2-64 characters, and it can contain lowercase letters, numbers, and dashes (-). It cannot start or end with a dash and cannot have two consecutive dashes. Need to be globally unique, and capitalization is not supported. Cannot start with `d-`.
* `mfa_authentication_status` - (Optional) The mfa authentication status. Valid values: `Enabled` or `Disabled`. Default to `Enabled`.
* `scim_synchronization_status` - (Optional) The scim synchronization status. Valid values: `Enabled` or `Disabled`. Default to `Disabled`.
* `saml_identity_provider_configuration` - (Optional) The saml identity provider configuration.
  * `encoded_metadata_document` - (Optional, Sensitive) Base64 encoded IdP metadata document. **NOTE:** If the IdP Metadata has been uploaded, no update will be made if this parameter is not specified, otherwise the update will be made according to the parameter content. If IdP Metadata has not been uploaded, and the parameter `sso_status` is `Enabled`, this parameter must be provided. If the IdP Metadata has not been uploaded, and the parameter `sso_status` is `Disabled`, this parameter can be omitted, and the IdP Metadata will remain empty.
  * `sso_status` - (Optional) SAML SSO login enabled status. Valid values: `Enabled` or `Disabled`. Default to `Disabled`.

-> **NOTE:** The `saml_identity_provider_configuration` will be removed automatically when the resource is deleted, please operate with caution. If there are left more configuration in the directory, please remove them before deleting the directory.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Directory.

## Import

Cloud SSO Directory can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_sso_directory.example <id>
```
