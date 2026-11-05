---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_requester_qos_info"
description: |-
  Provides a Alicloud OSS Bucket Requester QoS Info resource.
---

# alicloud_oss_bucket_requester_qos_info

Provides a OSS Bucket Requester QoS Info resource. The requester QoS info attaches a single QoS (quality of service) quota configuration to a specific requester (RAM user) of a bucket, capping the requester's bandwidth and queries-per-second across the extranet, intranet and total dimensions.

For information about OSS Bucket Requester QoS Info and how to use it, see [Bucket requester QoS info](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketrequesterqosinfo).

-> **NOTE:** Available since v2.1.0.

-> **NOTE:** A bucket can carry at most one QoS configuration per requester. Creating or updating the resource replaces the whole configuration, so all nine quota fields are required.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_requester_qos_info&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

resource "alicloud_ram_user" "CreateRequester" {
  name         = "${var.name}-${random_integer.default.result}"
  display_name = var.name
}

resource "alicloud_oss_bucket_requester_qos_info" "default" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket

  # The requester is the RAM user the QoS limits apply to.
  qos_requester = alicloud_ram_user.CreateRequester.id

  extranet_download_bandwidth  = 1048576
  extranet_qps                 = 100
  extranet_upload_bandwidth    = 1048576
  intranet_download_bandwidth  = 1048576
  intranet_qps                 = 100
  intranet_upload_bandwidth    = 1048576
  total_download_bandwidth     = 2097152
  total_qps                    = 200
  total_upload_bandwidth       = 2097152
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket to which the requester QoS info is attached.
* `qos_requester` - (Required, ForceNew) The requester the QoS configuration applies to, identified by its RAM user id.
* `extranet_download_bandwidth` - (Required) The maximum download bandwidth over the extranet, in bit/s.
* `extranet_qps` - (Required) The maximum queries per second over the extranet.
* `extranet_upload_bandwidth` - (Required) The maximum upload bandwidth over the extranet, in bit/s.
* `intranet_download_bandwidth` - (Required) The maximum download bandwidth over the intranet, in bit/s.
* `intranet_qps` - (Required) The maximum queries per second over the intranet.
* `intranet_upload_bandwidth` - (Required) The maximum upload bandwidth over the intranet, in bit/s.
* `total_download_bandwidth` - (Required) The maximum total download bandwidth, in bit/s.
* `total_qps` - (Required) The maximum total queries per second.
* `total_upload_bandwidth` - (Required) The maximum total upload bandwidth, in bit/s.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above. The value is formulated as `<bucket>:<qos_requester>`.

## Import

OSS Bucket Requester QoS Info can be imported using the id, which consists of bucket and qos_requester, e.g.

```shell
$ terraform import alicloud_oss_bucket_requester_qos_info.example <bucket>:<qos_requester>
```
