---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_cloud_gtm_address"
description: |-
  Provides a Alicloud Alidns Cloud Gtm Address resource.
---

# alicloud_alidns_cloud_gtm_address

Provides a Alidns Cloud Gtm Address resource.



For information about Alidns Cloud Gtm Address and how to use it, see [What is Cloud Gtm Address](https://next.api.alibabacloud.com/document/Alidns/2015-01-09/CreateCloudGtmAddress).

-> **NOTE:** Available since v1.275.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}


resource "alicloud_alidns_cloud_gtm_address" "default" {
  type             = "domain"
  health_judgement = "all_ok"
  address          = "www.tianxuan.top"
  enable_status    = "enable"
  available_mode   = "auto"
  name             = "resource-example-addr-3"
}
```

## Argument Reference

The following arguments are supported:
* `address` - (Required) This property does not have a description in the spec, please add it before generating code.
* `available_mode` - (Required) This property does not have a description in the spec, please add it before generating code.
* `enable_status` - (Required) This property does not have a description in the spec, please add it before generating code.
* `health_judgement` - (Required) This property does not have a description in the spec, please add it before generating code.
* `health_tasks` - (Optional, List) This property does not have a description in the spec, please add it before generating code. See [`health_tasks`](#health_tasks) below.
* `manual_available_status` - (Optional) This property does not have a description in the spec, please add it before generating code.
* `name` - (Required) The name of the resource
* `remark` - (Optional) This property does not have a description in the spec, please add it before generating code.
* `type` - (Required, ForceNew) This property does not have a description in the spec, please add it before generating code.

### `health_tasks`

The health_tasks supports the following:
* `port` - (Optional, Int) This property does not have a description in the spec, please add it before generating code.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - The creation time of the resource.
* `health_tasks` - 当前属性没有在镇元上录入属性描述�
  * `template_id` - 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cloud Gtm Address.
* `delete` - (Defaults to 5 mins) Used when delete the Cloud Gtm Address.
* `update` - (Defaults to 5 mins) Used when update the Cloud Gtm Address.

## Import

Alidns Cloud Gtm Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_cloud_gtm_address.example <address_id>
```