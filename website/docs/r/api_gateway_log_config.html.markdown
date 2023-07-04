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

```terraform
data "alicloud_api_gateway_log_configs" "default" {
  log_type = "PROVIDER"
}
locals {
  count = length(data.alicloud_api_gateway_log_configs.default.configs) > 0 ? 0 : 1
}

resource "random_integer" "default" {
  count = local.count
  max   = 99999
  min   = 10000
}

resource "alicloud_log_project" "example" {
  count       = local.count
  name        = "terraform-example-${random_integer.default[0].result}"
  description = "terraform-example"
}

resource "alicloud_log_store" "example" {
  count                 = local.count
  project               = alicloud_log_project.example[0].name
  name                  = "terraform-example"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_api_gateway_log_config" "example" {
  count         = local.count
  sls_project   = alicloud_log_project.example[0].name
  sls_log_store = alicloud_log_store.example[0].name
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