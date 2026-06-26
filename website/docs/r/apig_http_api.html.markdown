---
subcategory: "Cloud Native API Gateway (APIG)"
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
* `ai_protocols` - (Optional, List, Available since v1.284.0) AI protocols

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `base_path` - (Optional) API base path, which must start with a forward slash (/).

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `deploy_configs` - (Optional, List, Available since v1.284.0) API deployment configurations. Currently, only AI APIs support deployment configurations, and only a single deployment configuration can be specified.

-> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `description` - (Optional) API description.
* `enable_auth` - (Optional, Available since v1.284.0) Whether to enable authentication.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `http_api_name` - (Required, ForceNew) Perform an exact search by name.
* `model_category` - (Optional, Available since v1.284.0) AI model category

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `protocols` - (Required, List) List of API access protocols.
* `resource_group_id` - (Optional, Computed) Target resource group ID.

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `type` - (Optional, ForceNew) The type of the HTTP API. Multiple types are supported and must be separated by commas (,).  
  - Http  
  - Rest  
  - LLM  
  - WebSocket  
  - HttpIngress  

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
$ terraform import alicloud_apig_http_api.example <http_api_id>
```