---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_service_linked_role"
sidebar_current: "docs-alicloud-resource-dbfs-service-linked-role"
description: |-
    Provides a resource to create the Dbfs service-linked roles(SLR).
---

# alicloud\_dbfs\_service\_linked\_role

Using this data source can create Dbfs service-linked roles(SLR). Dbfs may need to access another Alibaba Cloud service to implement a specific feature. In this case, Dbfs must assume a specific service-linked role, which is a Resource Access Management (RAM) role, to obtain permissions to access another Alibaba Cloud service. 

For information about Dbfs service-linked roles(SLR) and how to use it, see [What is service-linked roles](https://www.alibabacloud.com/help/doc-detail/181425.htm).

-> **NOTE:** Available in v1.157.0+.


## Example Usage

```terraform
resource "alicloud_dbfs_service_linked_role" "service_linked_role" {
  product_name = "AliyunServiceRoleForDbfs"
}
```

## Argument Reference

The following arguments are supported:

* `product_name` - (Required, ForceNew) The product name for SLR. Dbfs can automatically create the following service-linked roles: `AliyunServiceRoleForDbfs`.

## Attributes Reference

* `id` - The ID of the Resource. The value is same as `product_name`.
* `status` - The status of the service Associated role. Valid Values: `true`: Created. `false`: not created.

## Import

Dbfs service-linked roles(SLR) can be imported using the id, e.g.

```
$ terraform import alicloud_dbfs_service_linked_role.example <product_name>
```
