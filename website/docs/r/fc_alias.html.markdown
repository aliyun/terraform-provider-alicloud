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


## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_alias&exampleId=ece83230-3dff-a95b-59d1-50065f03e40acf0b9512&activeTab=example&spm=docs.r.fc_alias.0.ece832303d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
