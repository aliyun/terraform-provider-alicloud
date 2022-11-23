---
subcategory: "DAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_das_switch_das_pro"
sidebar_current: "docs-alicloud-resource-das-switch-das-pro"
description: |-
  Provides a Alicloud DAS Switch Das Pro resource.
---

# alicloud\_das\_switch\_das\_pro

Provides a DAS Switch Das Pro resource.

For information about DAS Switch Das Pro and how to use it, see [What is Switch Das Pro](https://www.alibabacloud.com/help/en/database-autonomy-service/latest/enabledaspro).

-> **NOTE:** Available in v1.193.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_das_switch_das_pro" "default" {
  instance_id   = "your_sql_id"
  sql_retention = 30
  user_id       = "your_account_id"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the database instance.
* `sql_retention` - (Optional, Computed, ForceNew) The storage duration of SQL Explorer data. Valid values: `30`, `180`, `365`, `1095`, `1825`. Unit: days. Default value: `30`.
* `user_id` - (Optional, Computed, ForceNew) The ID of the Alibaba Cloud account that is used to create the database instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Switch Das Pro. Its value is same as `instance_id`.
* `status` - Whether the database instance has DAS professional.

## Import

DAS Switch Das Pro can be imported using the id, e.g.

```shell
$ terraform import alicloud_das_switch_das_pro.example <id>
```
