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

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cloud_sso_access_configuration&exampleId=2b1babd0-0321-8746-6959-5cb65fd1515e7aa117e2&activeTab=example&spm=docs.r.cloud_sso_access_configuration.0.2b1babd003&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cloud_sso_access_configuration" "default" {
  directory_id              = data.alicloud_cloud_sso_directories.default.directories.0.id
  access_configuration_name = var.name
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

* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `access_configuration_name` - (Required, ForceNew) The name of the access configuration. The name can be up to `32` characters long and can contain letters, digits, and hyphens (-).
* `session_duration` - (Optional, Int) The SessionDuration of the Access Configuration. Unit: Seconds. Valid values: `900` to `43200`.
* `relay_state` - (Optional) The RelayState of the Access Configuration, Cloud SSO users use this access configuration to access the RD account, the initial access page address. Must be the Alibaba Cloud console page, the default is the console home page.
* `description` - (Optional) The description of the access configuration. The description can be up to `1024` characters in length.
* `permission_policies` - (Optional, Set) The Policy List. See [`permission_policies`](#permission_policies) below.
* `force_remove_permission_policies` - (Optional, Bool) This parameter is used to force deletion `permission_policies`. Valid Value: `true`, `false`.

* **NOTE:** The `permission_policies` will be removed automatically when the resource is deleted, please operate with caution. If there are left more permission policies in the access configuration, please remove them before deleting the access configuration.

### `permission_policies`

The permission_policies supports the following: 

* `permission_policy_type` - (Required) The type of the policy. Valid values: `System`, `Inline`.
* `permission_policy_name` - (Required) The name of the policy.
* `permission_policy_document` - (Optional) The configurations of the inline policy. **NOTE:** If `permission_policy_type` is set to `Inline`, `permission_policy_document` is required.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Access Configuration. It formats as `<directory_id>:<access_configuration_id>`.
* `access_configuration_id` - The ID of the Access Configuration.

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
