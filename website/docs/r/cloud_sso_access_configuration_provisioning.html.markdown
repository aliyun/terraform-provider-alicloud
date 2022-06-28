---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_access_configuration_provisioning"
sidebar_current: "docs-alicloud-resource-cloud-sso-access-configuration-provisioning"
description: |-
  Provides a Alicloud Cloud SSO Access Configuration Provisioning resource.
---

# alicloud\_cloud\_sso\_access\_configuration\_provisioning

Provides a Cloud SSO Access Configuration Provisioning resource.

For information about Cloud SSO Access Configuration Provisioning and how to use it, see [What is Access Configuration Provisioning](https://www.alibabacloud.com/help/en/doc-detail/266737.html).

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage


```terraform
variable name {
  default = "example-name"
}

data alicloud_cloud_sso_directories default {}

data alicloud_resource_manager_resource_directories default {}

resource alicloud_cloud_sso_directory default {
  count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name = var.name
}

resource alicloud_cloud_sso_access_configuration default {
  access_configuration_name = var.name
  directory_id              = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id[""])[0]
}

resource alicloud_cloud_sso_access_configuration_provisioning default {
  directory_id            = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id[""])[0]
  access_configuration_id = alicloud_cloud_sso_access_configuration.default.access_configuration_id
  target_type             = "RD-Account"
  target_id               = data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id
}
```

## Argument Reference

The following arguments are supported:

* `access_configuration_id` - (Required, ForceNew) The Access configuration ID.
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `target_id` - (Required, ForceNew) The ID of the target to create the resource range.
* `target_type` - (Required, ForceNew) The type of the resource range target to be accessed. Valid values: `RD-Account`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Access Assignment. The value formats as `<directory_id>:<access_configuration_id>:<target_type>:<target_id>`.
* `status` - The status of the resource. Valid values: `Provisioned`, `ReprovisionRequired` and `DeprovisionFailed`.


## Import

Cloud SSO Access Configuration Provisioning can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_sso_access_assignment.example <directory_id>:<access_configuration_id>:<target_type>:<target_id>
```
