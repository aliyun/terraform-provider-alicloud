---
subcategory: "Live"
layout: "alicloud"
page_title: "Alicloud: alicloud_live_private_line"
description: |-
  Provides a Alicloud Live Private Line resource.
---

# alicloud_live_private_line

Provides a Live Private Line resource.

Binds a Global Accelerator (GA) instance to a live streaming link to enable private-line acceleration for push or pull streams, at the granularity of a single stream.

For information about Live Private Line and how to use it, see [What is LivePrivateLine](https://next.api.alibabacloud.com/document/live/2016-11-01/CreateLivePrivateLine).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** `alicloud_live_private_line` has no in-place update API. Changing any argument forces a destroy-and-recreate of the binding.

-> **NOTE:** `DeleteLivePrivateLine` only **unbinds** the GA instance from the live link; it does **not** release the GA instance. When `reuse = "no"` creates a new GA instance, `terraform destroy` leaves that GA instance running and billable — release it in the Global Accelerator console or with `alicloud_ga_accelerator` if you also manage it through Terraform.

## Example Usage

### Create a new acceleration channel (reuse = "no")

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "domain_name" {
  default = "demo.aliyundoc.com"
}

variable "app_name" {
  default = "live"
}

variable "stream_name" {
  default = "exampleStream"
}

resource "alicloud_live_private_line" "default" {
  domain_name       = var.domain_name
  app_name          = var.app_name
  stream_name       = var.stream_name
  video_center      = "cn-shanghai"
  acceleration_type = "play"
  acceleration_area = "ap-southeast-1"
  reuse             = "no"
  max_bandwidth     = 20
}
```

### Reuse an existing GA instance (reuse = "yes")

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "ga_instance_id" {
  description = "An existing Global Accelerator instance id to bind."
  default     = "ga-bp1iovsdpf01ym9su****"
}

resource "alicloud_live_private_line" "default" {
  domain_name       = "demo.aliyundoc.com"
  app_name          = "live"
  stream_name       = "exampleStream"
  video_center      = "cn-shanghai"
  acceleration_type = "publish"
  acceleration_area = "ap-southeast-1"
  reuse             = "yes"
  instance_id       = var.ga_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) The main pull-stream domain.
* `app_name` - (Required, ForceNew) The application name, matching the AppName in the push/pull-stream URL.
* `stream_name` - (Required, ForceNew) The live stream name, matching the StreamName in the push/pull-stream URL.
* `video_center` - (Required, ForceNew) The live center. Valid values: `cn-beijing`, `cn-shanghai`, `cn-shenzhen`, `cn-qingdao`, `ap-northeast-1`, `ap-southeast-5`, `eu-central-1`, `ap-southeast-1`, `ap-south-1`.
* `acceleration_type` - (Required, ForceNew) The acceleration type. Valid values: `play` (pull-stream acceleration), `publish` (push-stream acceleration).
* `acceleration_area` - (Required, ForceNew) The acceleration channel (GA endpoint region), e.g. `ap-southeast-1`.
* `reuse` - (Required, ForceNew) Whether to reuse an existing acceleration channel. Valid values: `yes` (reuse an existing GA instance), `no` (create a new acceleration channel). Write-only: not returned by the Describe entry, so it cannot be reconstructed after import and must be re-set in the configuration.
* `max_bandwidth` - (Optional, ForceNew) The acceleration bandwidth in Mbps. Required when `reuse = "no"`. This field is write-only: it is not returned by the Describe entry, so refresh keeps the configured value and import leaves it unset.
* `instance_id` - (Required when `reuse = "yes"`, Computed otherwise, ForceNew) The GA instance id. When `reuse = "yes"` it references the GA instance to bind; when `reuse = "no"` the API creates a GA instance and this attribute is populated from the binding detail.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id, in the format `<domain_name>:<app_name>:<stream_name>:<acceleration_type>`. There is no single server-side id; the four identity fields compose it.
* `status` - The instance status. `active`: available; `inactive`: unavailable.
* `binding_status` - Whether the GA instance is bound to the live link. `yes`: bound; `no`: not bound.
* `ip` - The acceleration IP address.
* `uid` - The UID.

## Timeouts

-> **NOTE:** `terraform-plugin-framework` does not support user-configurable resource timeouts yet. The resource uses fixed timeouts: `Create` 10 minutes, `Delete` 5 minutes.

## Import

Live Private Line can be imported using the composite id, e.g.

```shell
$ terraform import alicloud_live_private_line.default demo.aliyundoc.com:live:exampleStream:play
```

The id is `<domain_name>:<app_name>:<stream_name>:<acceleration_type>`. Write-only fields (`max_bandwidth`, `reuse`) are not reconstructable from the Describe entry and are left unset after import; re-set them in the configuration before the next apply.
