---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_service_linked_role"
sidebar_current: "docs-alicloud-resource-rds-service-linked-role"
description: |-
  Provides a Alicloud RDS Service Linked Role.
---

# alicloud\_rds\_service\_linked\_role

Provides a RDS Service Linked Role.

For information about RDS Service Linked Role and how to use it, see [What is Service Linked Role.](https://www.alibabacloud.com/help/en/doc-detail/171226.htm).

-> **NOTE:** Available in v1.189.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_rds_service_linked_role" "default" {
  service_name = "AliyunServiceRoleForRdsPgsqlOnEcs"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required, ForceNew) The product name for SLR. RDS can automatically create the following service-linked roles: `AliyunServiceRoleForRdsPgsqlOnEcs`, `AliyunServiceRoleForRDSProxyOnEcs`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Service Linked Role. The value formats as `<service_name>`.
* `role_name` - The name of the role.
* `role_id` - The ID of the role.
* `arn` - The Alibaba Cloud Resource Name (ARN) of the role.

## Import

RDS Service Linked Role can be imported using the id, e.g.

```
$ terraform import alicloud_rds_service_linked_role.default <service_name>
```