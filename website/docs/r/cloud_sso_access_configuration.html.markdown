---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_access_configuration"
sidebar_current: "docs-alicloud-resource-cloud-sso-access-configuration"
description: |-
  Provides a Alicloud Cloud SSO Access Configuration resource.
---

# alicloud_cloud_sso_access_configuration

Provides a Cloud SSO Access Configuration resource.

For information about Cloud SSO Access Configuration and how to use it, see [What is Access Configuration](https://www.alibabacloud.com/help/en/cloudsso/latest/api-cloudsso-2021-05-15-createaccessconfiguration).

-> **NOTE:** Available since v1.145.0.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cloud_sso_access_configuration&exampleId=cbcf9104-25c1-674a-bcb5-961dc0d7147d55e53399&activeTab=example&spm=docs.r.cloud_sso_access_configuration.0.cbcf910425&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cloud_sso_user" "default" {
  directory_id = local.directory_id
  user_name    = var.name
}

resource "alicloud_cloud_sso_access_configuration" "default" {
  access_configuration_name = var.name
  directory_id              = local.directory_id
  permission_policies {
    permission_policy_type     = "Inline"
    permission_policy_name     = var.name
    permission_policy_document = <<EOF
    {
        "Statement":[
      {
        "Action":"ecs:Get*",
        "Effect":"Allow",
        "Resource":[
            "*"
        ]
      }
      ],
        "Version": "1"
    }
    EOF
  }
}
```

## Argument Reference

The following arguments are supported:

* `access_configuration_name` - (Required, ForceNew) The AccessConfigurationName of the Access Configuration. The name of the resource. The name can be up to `32` characters long and can contain letters, digits, and hyphens (-).
* `description` - (Optional) The Description of the  Access Configuration. The description can be up to `1024` characters long.
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `permission_policies` - (Optional) The Policy List. See [`permission_policies`](#permission_policies) below.
* `relay_state` - (Optional) The RelayState of the Access Configuration, Cloud SSO users use this access configuration to access the RD account, the initial access page address. Must be the Alibaba Cloud console page, the default is the console home page.
* `session_duration` - (Optional) The SessionDuration of the Access Configuration. Valid Value: `900` to `43200`. Unit: Seconds.
* `force_remove_permission_policies` - (Optional) This parameter is used to force deletion `permission_policies`. Valid Value: `true` and `false`.

* **NOTE:** The `permission_policies` will be removed automatically when the resource is deleted, please operate with caution. If there are left more permission policies in the access configuration, please remove them before deleting the access configuration.

### `permission_policies`

The permission_policies supports the following: 

* `permission_policy_document` - (Optional, Sensitive) The Content of Policy.
* `permission_policy_name` - (Required) The Policy Name of policy. The name of the resource. 
* `permission_policy_type` - (Required) The Policy Type of policy. Valid values: `System`, `Inline`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Access Configuration. The value formats as `<directory_id>:<access_configuration_id>`.
* `access_configuration_id` - The AccessConfigurationId of the Access Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Access Configuration.
* `update` - (Defaults to 5 mins) Used when update the Access Configuration.
* `delete` - (Defaults to 5 mins) Used when delete the Access Configuration.


## Import

Cloud SSO Access Configuration can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_access_configuration.example <directory_id>:<access_configuration_id>
```
