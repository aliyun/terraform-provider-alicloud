---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_scim_server_credentials"
sidebar_current: "docs-alicloud-datasource-cloud-sso-scim-server-credentials"
description: |-
  Provides a list of Cloud Sso Scim Server Credentials to the user.
---

# alicloud\_cloud\_sso\_scim\_server\_credentials

This data source provides the Cloud Sso Scim Server Credentials of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region


## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_sso_scim_server_credentials" "ids" {
  directory_id = "example_value"
  ids          = ["example_value-1", "example_value-2"]
}
output "cloud_sso_scim_server_credential_id_1" {
  value = data.alicloud_cloud_sso_scim_server_credentials.ids.credentials.0.id
}
```

## Argument Reference

The following arguments are supported:

* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `ids` - (Optional, ForceNew, Computed)  A list of SCIM Server Credential IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The Status of the resource. Valid values: `Disabled`, `Enabled`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `credentials` - A list of Cloud Sso Scim Server Credentials. Each element contains the following attributes:
	* `create_time` - The CreateTime of the resource.
	* `credential_id` - The CredentialId of the resource.
	* `credential_secret` - The CredentialSecret of the resource.
	* `credential_type` - The CredentialType of the resource.
	* `directory_id` - The ID of the Directory.
	* `expire_time` - The ExpireTime of the resource.
	* `id` - The ID of the SCIM Server Credential.
	* `status` - The Status of the resource. Valid values: `Disabled`, `Enabled`.
