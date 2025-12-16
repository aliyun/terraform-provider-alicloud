---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_access_configuration_provisioning"
sidebar_current: "docs-alicloud-resource-cloud-sso-access-configuration-provisioning"
description: |-
  Provides a Alicloud Cloud SSO Access Configuration Provisioning resource.
---

# alicloud_cloud_sso_access_configuration_provisioning

Provides a Cloud SSO Access Configuration Provisioning resource.

For information about Cloud SSO Access Configuration Provisioning and how to use it, see [What is Access Configuration Provisioning](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-addpermissionpolicytoaccessconfiguration).

-> **NOTE:** Available since v1.148.0.

## Example Usage

Basic Usage


<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_access_configuration_provisioning&exampleId=0c81affa-1707-8255-a2c8-ba5f44a4da78a28fd0fc&activeTab=example&spm=docs.r.cloud_sso_access_configuration_provisioning.0.0c81affa17&intl_lang=EN_US" target="_blank">
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
data "alicloud_resource_manager_resource_directories" "default" {}

resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}

locals {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cloud_sso_user" "default" {
  directory_id = local.directory_id
  user_name    = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cloud_sso_access_configuration" "default" {
  access_configuration_name = "${var.name}-${random_integer.default.result}"
  directory_id              = local.directory_id
}

resource "alicloud_cloud_sso_access_configuration_provisioning" "default" {
  directory_id            = local.directory_id
  access_configuration_id = alicloud_cloud_sso_access_configuration.default.access_configuration_id
  target_type             = "RD-Account"
  target_id               = data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_sso_access_configuration_provisioning&spm=docs.r.cloud_sso_access_configuration_provisioning.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `access_configuration_id` - (Required, ForceNew) The Access configuration ID.
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `target_id` - (Required, ForceNew) The ID of the target to create the resource range.
* `target_type` - (Required, ForceNew) The type of the resource range target to be accessed. Valid values: `RD-Account`.
* `status` - (Optional) The status of the resource. Valid values: `Provisioned`, `ReprovisionRequired` and `DeprovisionFailed`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Access Assignment. The value formats as `<directory_id>:<access_configuration_id>:<target_type>:<target_id>`.


## Import

Cloud SSO Access Configuration Provisioning can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_access_assignment.example <directory_id>:<access_configuration_id>:<target_type>:<target_id>
```
