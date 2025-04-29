---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_data_flow"
sidebar_current: "docs-alicloud-resource-nas-data-flow"
description: |-
  Provides a Alicloud File Storage (NAS) Data Flow resource.
---

# alicloud_nas_data_flow

Provides a File Storage (NAS) Data Flow resource.

For information about File Storage (NAS) Data Flow and how to use it, see [What is Data Flow](https://www.alibabacloud.com/help/en/doc-detail/27530.html).

-> **NOTE:** Available since v1.153.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_data_flow&exampleId=6ff8d26a-7906-f0c9-065c-7bfc0d61b9d2412e48cb&activeTab=example&spm=docs.r.nas_data_flow.0.6ff8d26a79&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_nas_zones" "example" {
  file_system_type = "cpfs"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_nas_zones.example.zones[1].zone_id
}

resource "alicloud_nas_file_system" "example" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  description      = "terraform-example"
  zone_id          = data.alicloud_nas_zones.example.zones[1].zone_id
  vpc_id           = alicloud_vpc.example.id
  vswitch_id       = alicloud_vswitch.example.id
}

resource "alicloud_nas_mount_target" "example" {
  file_system_id = alicloud_nas_file_system.example.id
  vswitch_id     = alicloud_vswitch.example.id
}
resource "random_integer" "example" {
  max = 99999
  min = 10000
}
resource "alicloud_oss_bucket" "example" {
  bucket = "example-value-${random_integer.example.result}"
  acl    = "private"
  tags = {
    cpfs-dataflow = "true"
  }
}

resource "alicloud_nas_fileset" "example" {
  file_system_id   = alicloud_nas_mount_target.example.file_system_id
  description      = "terraform-example"
  file_system_path = "/example_path/"
}

resource "alicloud_nas_data_flow" "example" {
  fset_id              = alicloud_nas_fileset.example.fileset_id
  description          = "terraform-example"
  file_system_id       = alicloud_nas_file_system.example.id
  source_security_type = "SSL"
  source_storage       = join("", ["oss://", alicloud_oss_bucket.example.bucket])
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
* `source_security_type` - (Optional, ForceNew) The security protection type of the source storage. If the source storage must be accessed through security protection, specify the security protection type of the source storage. Value:
  - `NONE` (default): Indicates that the source storage does not need to be accessed through security protection.
  - `SSL`: Protects access through SSL certificates.
* `source_storage` - (Required, ForceNew) The access path of the source store. Format: `<storage type>://<path>`. Among them:
  - storage type: currently only OSS is supported.
  - path: the bucket name of OSS.
    - Only lowercase letters, numbers, and dashes (-) are supported and must start and end with lowercase letters or numbers.
    - `8` to `128` English characters in length.
    - Use UTF-8 coding.
    - Cannot start with `http://` and `https://`.
* `status` - (Optional) The status of the Data flow. Valid values: `Running`, `Stopped`.
* `throughput` - (Required) The maximum transmission bandwidth of data flow, unit: `MB/s`. Valid values: `1200`, `1500`, `600`. **NOTE:** The transmission bandwidth of data flow must be less than the IO bandwidth of the file system.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Data Flow. The value formats as `<file_system_id>:<data_flow_id>`.
* `data_flow_id` - The ID of the Data flow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Data Flow.
* `update` - (Defaults to 10 mins) Used when update the Data Flow.
* `delete` - (Defaults to 10 mins) Used when delete the Data Flow.

## Import

File Storage (NAS) Data Flow can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_data_flow.example <file_system_id>:<data_flow_id>
```