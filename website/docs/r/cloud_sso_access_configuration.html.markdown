---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_access_configuration"
sidebar_current: "docs-alicloud-resource-cloud-sso-access-configuration"
description: |-
  Provides a Alicloud Cloud SSO Access Configuration resource.
---

# alicloud\_cloud\_sso\_access\_configuration

Provides a Cloud SSO Access Configuration resource.

For information about Cloud SSO Access Configuration and how to use it, see [What is Access Configuration](https://www.alibabacloud.com/help/en/doc-detail/266737.html).

-> **NOTE:** Available in v1.145.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example-value"
}
data "alicloud_cloud_sso_directories" "default" {}
resource "alicloud_cloud_sso_directory" "default" {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}
resource "alicloud_cloud_sso_access_configuration" "default" {
  access_configuration_name = var.name
  directory_id              = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
  permission_policies {
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
    permission_policy_type     = "Inline"
    permission_policy_name     = var.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `access_configuration_name` - (Required, ForceNew) The AccessConfigurationName of the Access Configuration. The name of the resource. The name can be up to `32` characters long and can contain letters, digits, and hyphens (-).
* `description` - (Optional) The Description of the  Access Configuration. The description can be up to `1024` characters long.
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `permission_policies` - (Optional) The Policy List. See the following `Block permission_policies`.
* `relay_state` - (Optional) The RelayState of the Access Configuration, Cloud SSO users use this access configuration to access the RD account, the initial access page address. Must be the Alibaba Cloud console page, the default is the console home page.
* `session_duration` - (Optional, Computed) The SessionDuration of the Access Configuration. Valid Value: `900` to `43200`. Unit: Seconds.
* `force_remove_permission_policies` - (Optional) This parameter is used to force deletion `permission_policies`. Valid Value: `true` and `false`.

* **NOTE:** The `permission_policies` will be removed automatically when the resource is deleted, please operate with caution. If there are left more permission policies in the access configuration, please remove them before deleting the access configuration.

#### Block permission_policies

The permission_policies supports the following: 

* `permission_policy_document` - (Optional, Sensitive) The Content of Policy.
* `permission_policy_name` - (Required) The Policy Name of policy. The name of the resource. 
* `permission_policy_type` - (Required) The Policy Type of policy. Valid values: `System`, `Inline`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Access Configuration. The value formats as `<directory_id>:<access_configuration_id>`.
* `access_configuration_id` - The AccessConfigurationId of the Access Configuration.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Access Configuration.
* `update` - (Defaults to 5 mins) Used when update the Access Configuration.
* `delete` - (Defaults to 5 mins) Used when delete the Access Configuration.


## Import

Cloud SSO Access Configuration can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_sso_access_configuration.example <directory_id>:<access_configuration_id>
```
