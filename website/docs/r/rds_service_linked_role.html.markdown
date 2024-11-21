---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_service_linked_role"
sidebar_current: "docs-alicloud-resource-rds-service-linked-role"
description: |-
  Provides a Alicloud RDS Service Linked Role.
---

# alicloud_rds_service_linked_role

Provides a RDS Service Linked Role.

For information about RDS Service Linked Role and how to use it, see [What is Service Linked Role.](https://www.alibabacloud.com/help/en/doc-detail/171226.htm).

-> **NOTE:** Available since v1.189.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_service_linked_role&exampleId=0a7707d6-9760-35fd-28e2-89ac134ecea66e699b4b&activeTab=example&spm=docs.r.rds_service_linked_role.0.0a7707d697&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

```shell
$ terraform import alicloud_rds_service_linked_role.default <service_name>
```