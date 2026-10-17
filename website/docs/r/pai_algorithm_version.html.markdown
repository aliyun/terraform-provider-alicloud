---
subcategory: "PaiStudio"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_algorithm_version"
sidebar_current: "docs-alicloud-resource-pai-algorithm-version"
description: |-
  Provides a PaiStudio Algorithm Version resource.
---

# alicloud_pai_algorithm_version

Provides a PaiStudio Algorithm Version resource.

For information about PaiStudio Algorithm Version and how to use it, see
[What is AlgorithmVersion](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paistudio-2022-01-12-createalgorithmversion).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_pai_algorithm_version" "example" {
  algorithm_id      = "algo-xxxxxxxxxxxxxxxx"
  algorithm_version = "v1.0.0"
  algorithm_spec = jsonencode({
    JobType = "TensorFlow"
    Image   = "registry.cn-hangzhou.aliyuncs.com/my-namespace/tf:1.0"
    Command = "python train.py"
  })
}
```

## Argument Reference

The following arguments are supported:

* `algorithm_id` - (Required, ForceNew) The ID of the algorithm that the version belongs to. It is the parent algorithm registered in PaiStudio.
* `algorithm_version` - (Required, ForceNew) The version name of the algorithm version, e.g. `v1.0.0`.
* `algorithm_spec` - (Required) The algorithm spec, a JSON-encoded object matching the PaiStudio `AlgorithmSpec` schema (JobType, Image, Command, InputChannels, OutputChannels, HyperParameters, etc.). Use `jsonencode(...)` to build it from a Terraform object.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. It is formatted as `<algorithm_id>:<algorithm_version>`.
* `tenant_id` - The tenant ID of the algorithm.
* `user_id` - The user ID of the algorithm.
* `algorithm_provider` - The algorithm provider.
* `algorithm_name` - The name of the algorithm.
* `gmt_modified_time` - The time when the algorithm version was last modified.
* `region_id` - The region ID of the resource.
* `create_time` - The time when the algorithm version was created.

## Import

PaiStudio Algorithm Version can be imported using the composite id, e.g.

```shell
$ terraform import alicloud_pai_algorithm_version.example <algorithm_id>:<algorithm_version>
```
