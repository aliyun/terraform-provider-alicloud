---
subcategory: "APIG"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_http_api"
description: |-
  Provides a Alicloud APIG Http Api resource.
---

# alicloud_apig_http_api

Provides a APIG Http Api resource.



For information about APIG Http Api and how to use it, see [What is Http Api](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Http Api.
* `delete` - (Defaults to 5 mins) Used when delete the Http Api.
* `update` - (Defaults to 5 mins) Used when update the Http Api.

## Import

APIG Http Api can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_http_api.example <id>
```