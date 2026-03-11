---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_logtail_pipeline_config"
description: |-
  Provides a Alicloud Log Service (SLS) Logtail Pipeline Config resource.
---

# alicloud_sls_logtail_pipeline_config

Provides a Log Service (SLS) Logtail Pipeline Config resource.

Logtail Pipeline Collection Configuration.

For information about Log Service (SLS) Logtail Pipeline Config and how to use it, see [What is Logtail Pipeline Config](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateLogtailPipelineConfig).

-> **NOTE:** Available since v1.273.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}


resource "alicloud_sls_logtail_pipeline_config" "default" {
  project     = "terraform-logstore-example-578"
  config_name = "pl-auto-example"
  inputs {
  }
  flushers {
  }
}
```

## Argument Reference

The following arguments are supported:
* `aggregators` - (Optional, List) This property does not have a description in the spec, please add it before generating code. See [`aggregators`](#aggregators) below.
* `config_name` - (Required, ForceNew) The name of the resource
* `flushers` - (Required, List) This property does not have a description in the spec, please add it before generating code. See [`flushers`](#flushers) below.
* `globals` - (Optional, Map) This property does not have a description in the spec, please add it before generating code.
* `inputs` - (Required, List) The creation time of the resource See [`inputs`](#inputs) below.
* `log_sample` - (Optional) This property does not have a description in the spec, please add it before generating code.
* `processors` - (Optional, List) This property does not have a description in the spec, please add it before generating code. See [`processors`](#processors) below.
* `project` - (Required, ForceNew) The first ID of the resource
* `task` - (Optional, Map) This property does not have a description in the spec, please add it before generating code.

### `aggregators`

The aggregators supports the following:

### `flushers`

The flushers supports the following:

### `inputs`

The inputs supports the following:

### `processors`

The processors supports the following:

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<project>:<config_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Logtail Pipeline Config.
* `delete` - (Defaults to 5 mins) Used when delete the Logtail Pipeline Config.
* `update` - (Defaults to 5 mins) Used when update the Logtail Pipeline Config.

## Import

Log Service (SLS) Logtail Pipeline Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_logtail_pipeline_config.example <project>:<config_name>
```