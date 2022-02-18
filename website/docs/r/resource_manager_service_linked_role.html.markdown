---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_service_linked_role"
sidebar_current: "docs-alicloud-resource-resource-manager-service-linked-role"
description: |-
  Provides a Alicloud Resource Manager Service Linked Role.
---

# alicloud\_resource\_manager\_service\_linked\_role

Provides a Resource Manager Service Linked Role.

For information about Resource Manager Service Linked Role and how to use it, see [What is Service Linked Role.](https://www.alibabacloud.com/help/en/doc-detail/171226.htm).

-> **NOTE:** Available in v1.157.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_service_linked_role" "default" {
  service_name = "ops.elasticsearch.aliyuncs.com"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, ForceNew) The service name. For more information about the service name, see [Cloud services that support service linked roles](https://www.alibabacloud.com/help/en/doc-detail/160674.htm)
* `custom_suffix` - (Optional, ForceNew) The suffix of the role name. Only a few service linked roles support custom suffixes. The role name (including its suffix) must be 1 to 64 characters in length and can contain letters, digits, periods (.), and hyphens (-). For example, if the suffix is Example, the role name is ServiceLinkedRoleName_Example.
* `description` - (Optional, ForceNew) The description of the service linked role.  This parameter must be specified for only the service linked roles that support custom suffixes. Otherwise, the preset value is used and cannot be modified. The description must be 1 to 1,024 characters in length.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Service Linked Role. The value formats as `<service_name>:<role_name>`.
* `role_name` - The name of the role.
* `role_id` - The ID of the role.
* `arn` - The Alibaba Cloud Resource Name (ARN) of the role.

## Import

Resource Manager Service Linked Role can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_service_linked_role.default <service_name>:<role_name>
```