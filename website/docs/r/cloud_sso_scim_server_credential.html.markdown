---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_scim_server_credential"
sidebar_current: "docs-alicloud-resource-cloud-sso-scim-server-credential"
description: |-
  Provides a Alicloud Cloud SSO SCIM Server Credential resource.
---

# alicloud_cloud_sso_scim_server_credential

Provides a Cloud SSO SCIM Server Credential resource.

For information about Cloud SSO SCIM Server Credential and how to use it, see [What is Cloud SSO SCIM Server Credential](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-createscimservercredential).

-> **NOTE:** Available since v1.138.0.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region


## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_scim_server_credential&exampleId=ec2683a7-924c-f9a5-509a-3100f10a75c5bd9f9fe8&activeTab=example&spm=docs.r.cloud_sso_scim_server_credential.0.ec2683a792&intl_lang=EN_US" target="_blank">
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

locals {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}

resource "alicloud_cloud_sso_scim_server_credential" "default" {
  directory_id = local.directory_id
}
```

## Argument Reference

The following arguments are supported:

* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `status` - (Optional) The Status of the resource. Valid values: `Disabled`, `Enabled`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of SCIM Server Credential. The value formats as `<directory_id>:<credential_id>`.
* `credential_id` - The CredentialId of the resource.

## Import

Cloud SSO SCIM Server Credential can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_scim_server_credential.example <directory_id>:<credential_id>
```
