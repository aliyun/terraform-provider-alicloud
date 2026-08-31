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

For information about PAI Service and how to use it, see [What is Service](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eas-2021-07-01-createservice).

-> **NOTE:** Field `labels` has been removed since version 1.245.0. Please use new field `tags`.

-> **NOTE:** Available since v1.238.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pai_service&exampleId=caf7b11d-148a-9e8f-649f-42377d0e81610fcd8aba&activeTab=example&spm=docs.r.pai_service.0.caf7b11d14&intl_lang=EN_US" target="_blank">
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

resource "alicloud_pai_service" "default" {
  develop        = "false"
  service_config = jsonencode({ "metadata" : { "cpu" : 1, "gpu" : 0, "instance" : 1, "memory" : 2000, "name" : "tfexample", "rpc" : { "keepalive" : 70000 } }, "model_path" : "http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz", "processor_entry" : "libecho.so", "processor_path" : "http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz", "processor_type" : "cpp" })
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_pai_service&spm=docs.r.pai_service.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `develop` - (Optional) Whether to enter the development mode.
* `service_config` - (Required, JsonString) Service configuration information. Please refer to https://www.alibabacloud.com/help/en/pai/user-guide/parameters-of-model-services
* `status` - (Optional, Computed) Service Current Status.
* `tags` - (Optional, Map, Available since v1.245.0) The tag of the resource.
* `workspace_id` - (Optional) Workspace id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time of the service
* `region_id` - The region ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service.
* `delete` - (Defaults to 5 mins) Used when delete the Service.
* `update` - (Defaults to 16 mins) Used when update the Service.

## Import

PAI Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_service.example <id>
```