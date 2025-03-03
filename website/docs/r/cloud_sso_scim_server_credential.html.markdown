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
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_cloud_sso_directories" "default" {
}

resource "alicloud_cloud_sso_scim_server_credential" "default" {
  directory_id           = data.alicloud_cloud_sso_directories.default.directories.0.id
  credential_secret_file = "./credential_secret_file.txt"
}
```

## Argument Reference

The following arguments are supported:

* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `status` - (Optional) The status of the SCIM Server Credential. Valid values: `Enabled`, `Disabled`.
* `credential_secret_file` - (Optional, ForceNew, Available since v1.245.0) The name of file that can save Credential ID and Credential Secret. Strongly suggest you to specified it when you creating credential, otherwise, you wouldn't get its secret ever.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of SCIM Server Credential. It formats as `<directory_id>:<credential_id>`.
* `credential_id` - The ID of the SCIM credential.
* `credential_type` - (Available since v1.245.0) The type of the SCIM credential.
* `create_time` - (Available since v1.245.0) The time when the SCIM credential was created.
* `expire_time` - (Available since v1.245.0) The time when the SCIM credential expires.

## Import

Cloud SSO SCIM Server Credential can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_scim_server_credential.example <directory_id>:<credential_id>
```
