---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_access_assignment"
sidebar_current: "docs-alicloud-resource-cloud-sso-access-assignment"
description: |-
  Provides a Alicloud Cloud SSO Access Assignment resource.
---

# alicloud\_cloud\_sso\_access\_assignment

Provides a Cloud SSO Access Assignment resource.

For information about Cloud SSO Access Assignment and how to use it, see [What is Access Assignment](https://www.alibabacloud.com/help/en/doc-detail/265996.htm).

-> **NOTE:** When you configure access assignment for the first time, access configuration will be automatically deployed.

-> **NOTE:** Available in v1.145.0+.

-> **NOTE:** Cloud SSO Only Support `cn-shanghai` And `us-west-1` Region

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

resource alicloud_cloud_sso_user default {
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id[""])[0]
  user_name    = var.name
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

resource alicloud_cloud_sso_access_assignment default {
  directory_id            = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id[""])[0]
  access_configuration_id = alicloud_cloud_sso_access_configuration.default.access_configuration_id
  target_type             = "RD-Account"
  target_id               = data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id
  principal_type          = "User"
  principal_id            = alicloud_cloud_sso_user.default.user_id
  deprovision_strategy    = "DeprovisionForLastAccessAssignmentOnAccount"
  depends_on              = [alicloud_cloud_sso_access_configuration_provisioning.default]
}
```

## Argument Reference

The following arguments are supported:

* `access_configuration_id` - (Required, ForceNew) The Access configuration ID.
* `deprovision_strategy` - (Optional) The deprovision strategy. Valid values: `DeprovisionForLastAccessAssignmentOnAccount` and `None`. Default Value: `DeprovisionForLastAccessAssignmentOnAccount`. **NOTE:** When `deprovision_strategy` is `DeprovisionForLastAccessAssignmentOnAccount`, and the access assignment to be deleted is the last access assignment for the same account and the same AC, this option is used for the undeployment operation。
* `directory_id` - (Required, ForceNew) The ID of the Directory.
* `principal_id` - (Required, ForceNew)The ID of the access assignment.
* `principal_type` - (Required, ForceNew) The identity type of the access assignment, which can be a user or a user group. Valid values: `Group`, `User`.
* `target_id` - (Required, ForceNew) The ID of the target to create the resource range.
* `target_type` - (Required, ForceNew) The type of the resource range target to be accessed. Valid values: `RD-Account`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Access Assignment. The value formats as `<directory_id>:<access_configuration_id>:<target_type>:<target_id>:<principal_type>:<principal_id>`. 

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Cloud SSO Access Assignment.
* `delete` - (Defaults to 5 mins) Used when delete the Cloud SSO Access Assignment.

## Import

Cloud SSO Access Assignment can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_sso_access_assignment.example <directory_id>:<access_configuration_id>:<target_type>:<target_id>:<principal_type>:<principal_id>
```

