---
subcategory: "Database File System (DBFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dbfs_service_linked_role"
sidebar_current: "docs-alicloud-resource-dbfs-service-linked-role"
description: |-
    Provides a resource to create the Dbfs service-linked roles(SLR).
---

# alicloud_dbfs_service_linked_role

Using this data source can create Dbfs service-linked roles(SLR). Dbfs may need to access another Alibaba Cloud service to implement a specific feature. In this case, Dbfs must assume a specific service-linked role, which is a Resource Access Management (RAM) role, to obtain permissions to access another Alibaba Cloud service. 

For information about Dbfs service-linked roles(SLR) and how to use it, see [What is service-linked roles](https://www.alibabacloud.com/help/en/resource-management/resource-group/developer-reference/api-resourcemanager-2020-03-31-createservicelinkedrole-rg).

-> **NOTE:** Available since v1.157.0.


## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dbfs_service_linked_role&exampleId=222941ef-fd7b-75d2-cecf-c2b606b58881e5ebff02&activeTab=example&spm=docs.r.dbfs_service_linked_role.0.222941effd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

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

```shell
$ terraform import alicloud_dbfs_service_linked_role.example <product_name>
```
