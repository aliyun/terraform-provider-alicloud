---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_access_point_policy"
description: |-
  Provides a Alicloud OSS Access Point Policy resource.
---

# alicloud_oss_access_point_policy

Provides a OSS Access Point Policy resource. The access point policy is a JSON-formatted authorization policy attached to an OSS access point, which controls who can access the bucket through this access point.

For information about OSS Access Point Policy and how to use it, see [Manage access point policies](https://www.alibabacloud.com/help/en/oss/user-guide/access-points).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_oss_access_point" "CreateAccessPoint" {
  access_point_name = var.name
  bucket            = alicloud_oss_bucket.CreateBucket.bucket
  network_origin    = "vpc"
  vpc_configuration {
    vpc_id = "vpc-abcexample"
  }
  public_access_block_configuration {
    block_public_access = false
  }
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_oss_access_point_policy" "default" {
  access_point_name = alicloud_oss_access_point.CreateAccessPoint.access_point_name
  bucket            = alicloud_oss_bucket.CreateBucket.bucket
  policy = jsonencode({
    "Version" : "1",
    "Statement" : [{
      "Action" : ["oss:PutObject", "oss:GetObject"],
      "Effect" : "Allow",
      "Principal" : [data.alicloud_caller_identity.current.account_id],
      "Resource" : ["acs:oss:cn-hangzhou:${data.alicloud_caller_identity.current.account_id}:accesspoint/${alicloud_oss_access_point.CreateAccessPoint.access_point_name}/object/*"]
    }]
  })
}
```

## Argument Reference

The following arguments are supported:
* `access_point_name` - (Required, ForceNew) The name of the access point.
* `bucket` - (Required, ForceNew) The name of the bucket to which the access point belongs.
* `policy` - (Required) The JSON-formatted access point policy document. The policy is compared semantically (key ordering and whitespace differences are ignored).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<bucket>:<access_point_name>`.

## Import

OSS Access Point Policy can be imported using the id, which consists of bucket and access_point_name, e.g.

```shell
$ terraform import alicloud_oss_access_point_policy.example <bucket>:<access_point_name>
```
