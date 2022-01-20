---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_data_flow"
sidebar_current: "docs-alicloud-resource-nas-data-flow"
description: |-
  Provides a Alicloud Network Attached Storage (NAS) Data Flow resource.
---

# alicloud\_nas\_data\_flow

Provides a Network Attached Storage (NAS) Data Flow resource.

For information about Network Attached Storage (NAS) Data Flow and how to use it, see [What is Data Flow](https://www.alibabacloud.com/help/en/doc-detail/27530.html).

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id    = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  description      = "tf-testacc"
  zone_id          = local.zone_id
  vpc_id           = data.alicloud_vpcs.default.ids.0
  vswitch_id       = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_nas_mount_target" "default" {
  file_system_id = alicloud_nas_file_system.default.id
  vswitch_id     = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_oss_bucket" "default" {
  bucket = "example_value"
  acl    = "private"
  tags = {
    cpfs-dataflow = "true"
  }
}

resource "alicloud_nas_fileset" "default" {
  depends_on       = ["alicloud_nas_mount_target.default"]
  file_system_id   = alicloud_nas_file_system.default.id
  description      = "example_value"
  file_system_path = "/example_path/"
}

resource "alicloud_nas_data_flow" "default" {
  fset_id              = alicloud_nas_fileset.default.fileset_id
  description          = "example_value"
  file_system_id       = alicloud_nas_file_system.default.id
  source_security_type = "SSL"
  source_storage       = join("", ["oss://", alicloud_oss_bucket.default.bucket])
  throughput           = 600
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The Description of the data flow. Restrictions:
  - `2` ~ `128` English or Chinese characters in length.
  - Must start with uppercase or lowercase letters or Chinese, and cannot start with `http://` and `https://`.
  - Can contain numbers, semicolons (:), underscores (_), or dashes (-).
* `dry_run` - (Optional) The dry run.
* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `fset_id` - (Required, ForceNew) The ID of the Fileset.
* `source_security_type` - (Optional, Computed, ForceNew) The security protection type of the source storage. If the source storage must be accessed through security protection, specify the security protection type of the source storage. Value:
  - `NONE` (default): Indicates that the source storage does not need to be accessed through security protection.
  - `SSL`: Protects access through SSL certificates.
* `source_storage` - (Required, ForceNew) The access path of the source store. Format: `<storage type>://<path>`. Among them:
  - storage type: currently only OSS is supported.
  - path: the bucket name of OSS.
    - Only lowercase letters, numbers, and dashes (-) are supported and must start and end with lowercase letters or numbers.
    - `8` to `128` English characters in length.
    - Use UTF-8 coding.
    - Cannot start with `http://` and `https://`.
* `status` - (Optional, Computed) The status of the Data flow. Valid values: `Running`, `Stopped`.
* `throughput` - (Required) The maximum transmission bandwidth of data flow, unit: `MB/s`. Valid values: `1200`, `1500`, `600`. **NOTE:** The transmission bandwidth of data flow must be less than the IO bandwidth of the file system.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Data Flow. The value formats as `<file_system_id>:<data_flow_id>`.
* `data_flow_id` - The ID of the Data flow.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Data Flow.
* `update` - (Defaults to 10 mins) Used when update the Data Flow.
* `delete` - (Defaults to 10 mins) Used when delete the Data Flow.

## Import

Network Attached Storage (NAS) Data Flow can be imported using the id, e.g.

```
$ terraform import alicloud_nas_data_flow.example <file_system_id>:<data_flow_id>
```