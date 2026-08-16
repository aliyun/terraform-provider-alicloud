---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_resource_group"
sidebar_current: "docs-alicloud-resource-oss-bucket-resource-group"
description: |-
  Provides a OSS bucket resource group configuration resource.
---

# alicloud\_oss\_bucket\_resource\_group

Provides an independent resource group configuration resource for OSS bucket.

For information about OSS resource group and how to use it, see [Use resource groups](https://www.alibabacloud.com/help/oss/user-guide/configure-a-resource-group).

-> **NOTE:** Available since v1.217.3.

## Example Usage

Set bucket resource group configuration

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket" {
  bucket = "example-src-${random_integer.default.result}"
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_oss_bucket_resource_group" "bucket-resource-group" {
  bucket            = alicloud_oss_bucket.bucket.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket.
* `resource_group_id` - (Required) The ID of the resource group to which the bucket belongs.


## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.

## Import

Oss Bucket Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_resource_group.example
```

