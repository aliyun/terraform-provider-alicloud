---
subcategory: "APIG"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_http_api"
description: |-
  Provides a Alicloud APIG Http Api resource.
---

# alicloud_apig_http_api

Provides a APIG Http Api resource.



For information about APIG Http Api and how to use it, see [What is Http Api](https://next.api.aliyun.com/api/APIG/2024-03-27/CreateHttpApi).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_apig_http_api&exampleId=663376aa-68d1-c514-bf85-c36ddfc4ff787f19c9d2&activeTab=example&spm=docs.r.apig_http_api.0.663376aa68&intl_lang=EN_US" target="_blank">
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

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}

data "alicloud_resource_manager_resource_groups" "default" {}


resource "alicloud_apig_http_api" "default" {
  http_api_name = var.name
  protocols     = ["${var.protocol}"]
  base_path     = "/v1"
  description   = "zhiwei_pop_examplecase"
  type          = "Rest"
}
```

## Argument Reference

The following arguments are supported:
* `base_path` - (Optional) API path
* `description` - (Optional) Description of API
* `http_api_name` - (Required, ForceNew) The name of the resource
* `protocols` - (Required, List) API protocol
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `type` - (Optional, ForceNew) API type

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Http Api.
* `delete` - (Defaults to 5 mins) Used when delete the Http Api.
* `update` - (Defaults to 5 mins) Used when update the Http Api.

## Import

APIG Http Api can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_http_api.example <id>
```