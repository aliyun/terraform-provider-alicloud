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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sls_logtail_pipeline_config&exampleId=b1a0c32d-b81a-159e-faad-c86ce111ecccec86be0c&activeTab=example&spm=docs.r.sls_logtail_pipeline_config.0.b1a0c32db8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "example" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = "terraform logtail pipeline config example"
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "example-store"
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
}

resource "alicloud_sls_logtail_pipeline_config" "example" {
  project     = alicloud_log_project.example.project_name
  config_name = "${var.name}-${random_integer.default.result}"

  inputs = [{
    Type                     = "input_file"
    FilePaths                = "[\\\"/home/*.log\\\"]"
    EnableContainerDiscovery = false
    MaxDirSearchDepth        = 0
    FileEncoding             = "utf8"
  }]

  processors = [{
    Type      = "processor_parse_regex_native"
    SourceKey = "content"
    Regex     = ".*"
    Keys      = "[\\\"key1\\\",\\\"key2\\\"]"
  }]

  flushers = [{
    Type          = "flusher_sls"
    Logstore      = alicloud_log_store.example.logstore_name
    TelemetryType = "logs"
    Region        = "cn-shanghai"
    Endpoint      = "cn-shanghai-intranet.log.aliyuncs.com"
  }]

  aggregators = [{
    Type           = "aggregator_default"
    MaxSizeBytes   = 1048576
    MaxTimeSeconds = 3
  }]
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_sls_logtail_pipeline_config&spm=docs.r.sls_logtail_pipeline_config.example&intl_lang=EN_US)

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