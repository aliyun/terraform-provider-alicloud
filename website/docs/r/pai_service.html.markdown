---
subcategory: "PAI"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_service"
description: |-
  Provides a Alicloud PAI Service resource.
---

# alicloud_pai_service

Provides a PAI Service resource.

Eas service instance.

For information about PAI Service and how to use it, see [What is Service](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.238.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_service" "default" {
  labels = {
    "0" = jsonencode({ "LabelKey" : "examplekey", "LabelValue" : "examplevalue" })
  }
  develop        = "false"
  service_config = jsonencode({ "metadata" : { "cpu" : 1, "gpu" : 0, "instance" : 1, "memory" : 2000, "name" : "tfexample", "rpc" : { "keepalive" : 70000 } }, "model_path" : "http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz", "processor_entry" : "libecho.so", "processor_path" : "http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz", "processor_type" : "cpp" })
}
```

## Argument Reference

The following arguments are supported:
* `develop` - (Optional) Whether to enter the development mode.
* `labels` - (Optional, Map) Service Tag.
* `service_config` - (Required, JsonString) Service configuration information. Please refer to https://www.alibabacloud.com/help/en/pai/user-guide/parameters-of-model-services
* `status` - (Optional, Computed) Service Current Status, valid values `Running`, `Stopped`.
* `workspace_id` - (Optional) Workspace id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time of the service
* `region_id` - The region ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service.
* `delete` - (Defaults to 5 mins) Used when delete the Service.
* `update` - (Defaults to 5 mins) Used when update the Service.

## Import

PAI Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_service.example <id>
```