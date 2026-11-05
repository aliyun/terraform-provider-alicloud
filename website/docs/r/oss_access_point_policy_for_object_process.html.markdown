---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_access_point_policy_for_object_process"
description: |-
  Provides a Alicloud OSS Access Point Policy For Object Process resource.
---

# alicloud_oss_access_point_policy_for_object_process

Provides a OSS Access Point Policy For Object Process resource. The policy is a JSON-formatted authorization policy attached to an OSS Object FC access point (对象 FC 接入点), which controls who can access the bucket through the Object FC access point.

For information about OSS Object FC Access Point and how to use it, see [Manage Object FC access points](https://www.alibabacloud.com/help/en/oss/user-guide/object-fc-access-point-management).

-> **NOTE:** Available since v2.1.0.

-> **NOTE:** The Object FC access point is a separate prerequisite that must be created before the policy. Terraform does not yet manage the Object FC access point itself; create it via the OSS console, ossutil, or the `CreateAccessPointForObjectProcess` API first, then reference its name and bucket here.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_access_point_policy_for_object_process&activeTab=example&spm=docs.r.oss_access_point_policy_for_object_process.0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

variable "access_point_name" {
  default = "fc-ap-01"
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

data "alicloud_caller_identity" "current" {}

# The Object FC access point named by var.access_point_name must already exist
# on the bucket above; create it with ossutil or the OSS console first.

resource "alicloud_oss_access_point_policy_for_object_process" "default" {
  access_point_name = var.access_point_name
  bucket            = alicloud_oss_bucket.CreateBucket.bucket
  policy = jsonencode({
    "Version" : "1",
    "Statement" : [{
      "Action" : ["oss:GetObject"],
      "Effect" : "Allow",
      "Principal" : [data.alicloud_caller_identity.current.account_id],
      "Resource" : ["acs:oss:cn-hangzhou:${data.alicloud_caller_identity.current.account_id}:accesspointforobjectprocess/${var.access_point_name}/object/*"]
    }]
  })
}
```

## Argument Reference

The following arguments are supported:
* `access_point_name` - (Required, ForceNew) The name of the Object FC access point.
* `bucket` - (Required, ForceNew) The name of the bucket to which the Object FC access point belongs.
* `policy` - (Required) The JSON-formatted access point policy document. The policy is compared semantically (key ordering and whitespace differences are ignored).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<bucket>:<access_point_name>`.

## Import

OSS Access Point Policy For Object Process can be imported using the id, which consists of bucket and access_point_name, e.g.

```shell
$ terraform import alicloud_oss_access_point_policy_for_object_process.example <bucket>:<access_point_name>
```
