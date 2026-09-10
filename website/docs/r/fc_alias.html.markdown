---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_alias"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides an Alicloud Function Compute Alias resource. 
---

# alicloud_fc_alias

Creates a Function Compute service alias. Creates an alias that points to the specified Function Compute service version. 
 For the detailed information, please refer to the [developer guide](https://www.alibabacloud.com/help/en/fc/developer-reference/api-createalias).

-> **NOTE:** Available since v1.104.0.

-> **DEPRECATED:** This resource (alicloud_fc_alias) represents the legacy Function Compute API. For new deployments and advanced features, please use the `alicloud_fcv3_alias` resource which provides access to the latest FCV3 capabilities. The legacy FC alias supports only basic functionality and does not include newer features available in FCV3.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_fc_service" "default" {
  name        = "example-value-${random_integer.default.result}"
  description = "example-value"
  publish     = "true"
}

resource "alicloud_fc_alias" "example" {
  alias_name      = "example-value"
  description     = "example-value"
  service_name    = alicloud_fc_service.default.name
  service_version = "1"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_fc_alias&spm=docs.r.fc_alias.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, ForceNew) Name for the alias you are creating. 
* `description` - (Optional) Description of the alias.
* `service_name` - (Required, ForceNew) The Function Compute service name.
* `service_version` - (Required) The Function Compute service version for which you are creating the alias. Pattern: (LATEST|[0-9]+).
* `routing_config` - (Optional) The Function Compute alias' route configuration settings. See [`routing_config`](#routing_config) below.

### `routing_config`

The routing_config supports the following:

* `additional_version_weights` - (Optional) A map that defines the proportion of events that should be sent to different versions of a Function Compute service.


## Import

Function Compute alias can be imported using the id, e.g.

```shell
$ terraform import alicloud_fc_alias.example my_alias_id
```
