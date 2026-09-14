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

-> **NOTE:** Available since v1.240.0.

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
  access_point_name = "${var.name}-ap"
  bucket            = alicloud_oss_bucket.CreateBucket.bucket
  network_origin    = "internet"
}

resource "alicloud_oss_access_point_policy" "default" {
  access_point_name = alicloud_oss_access_point.CreateAccessPoint.access_point_name
  bucket            = alicloud_oss_bucket.CreateBucket.bucket
  policy            = jsonencode({ "Version" : "1", "Statement" : [{ "Action" : ["oss:PutObject", "oss:GetObject"], "Effect" : "Allow", "Principal" : ["1234567890"], "Resource" : ["acs:oss:*:1234567890:*/*"] }] })
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Point Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Access Point Policy.
* `update` - (Defaults to 5 mins) Used when update the Access Point Policy.

## Import

OSS Access Point Policy can be imported using the id, which consists of bucket and access_point_name, e.g.

```shell
$ terraform import alicloud_oss_access_point_policy.example <bucket>:<access_point_name>
```
