---
subcategory: "API Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_log_config"
sidebar_current: "docs-alicloud-resource-api-gateway-log-config"
description: |-
  Provides a Alicloud Api Gateway Log Config resource.
---

# alicloud\_api\_gateway\_log\_config

Provides a Api Gateway Log Config resource.

For information about Api Gateway Log Config and how to use it, see [What is Log Config](https://help.aliyun.com/document_detail/400392.html).

-> **NOTE:** Available in v1.185.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_api_gateway_log_config" "default" {
  sls_project   = "example_value"
  sls_log_store = "example_value"
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Api Gateway Log Config.
* `update` - (Defaults to 3 mins) Used when update the Api Gateway Log Config.
* `delete` - (Defaults to 3 mins) Used when delete the Api Gateway Log Config.

## Import

Api Gateway Log Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_log_config.example <log_type>
```