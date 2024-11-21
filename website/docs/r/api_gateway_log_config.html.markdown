---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_log_config"
sidebar_current: "docs-alicloud-resource-api-gateway-log-config"
description: |-
  Provides a Alicloud Api Gateway Log Config resource.
---

# alicloud_api_gateway_log_config

Provides a Api Gateway Log Config resource.

For information about Api Gateway Log Config and how to use it, see [What is Log Config](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-createlogconfig).

-> **NOTE:** Available since v1.185.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_log_config&exampleId=5087d780-a384-eb1f-f937-0ec18e6dd4646168907d&activeTab=example&spm=docs.r.api_gateway_log_config.0.5087d780a3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = var.name
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "${var.name}-${random_integer.default.result}"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_api_gateway_log_config" "example" {
  sls_project   = alicloud_log_project.example.project_name
  sls_log_store = alicloud_log_store.example.logstore_name
  log_type      = "PROVIDER"
}
```

## Argument Reference

The following arguments are supported:

* `sls_project` - (Required) The name of the Project.
* `sls_log_store` - (Required) The name of the Log Store.
* `log_type` - (Required, ForceNew) The type the of log. Valid values: `PROVIDER`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Log Config. Its value is same as `log_type`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Api Gateway Log Config.
* `update` - (Defaults to 3 mins) Used when update the Api Gateway Log Config.
* `delete` - (Defaults to 3 mins) Used when delete the Api Gateway Log Config.

## Import

Api Gateway Log Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_log_config.example <log_type>
```