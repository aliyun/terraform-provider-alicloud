---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_source"
description: |-
  Provides a Alicloud APIG Source resource.
---

# alicloud_apig_source

Provides a APIG Source resource.

Service source of a cloud-native API gateway, such as an ACK cluster or an MSE Nacos instance.

For information about APIG Source and how to use it, see [What is Source](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateSource).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "cluster_id" {
  description = "The ID of an existing ACK cluster used as the gateway service source"
  type        = string
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "10.0.0.0/24"
}

resource "alicloud_apig_gateway" "default" {
  gateway_name    = var.name
  spec            = "apigw.small.x1"
  gateway_edition = "Professional"
  gateway_type    = "API"
  payment_type    = "PayAsYouGo"
  vpc {
    vpc_id = alicloud_vpc.default.id
  }
  vswitch {
    vswitch_id = alicloud_vswitch.default.id
  }
  network_access_config {
    type = "Internet"
  }
  zone_config {
    select_option = "Auto"
  }
  log_config {
    sls {
      enable = false
    }
  }
}

resource "alicloud_apig_source" "default" {
  type       = "K8S"
  gateway_id = alicloud_apig_gateway.default.id
  k8s_source_info {
    cluster_id = var.cluster_id
  }
}
```

## Argument Reference

The following arguments are supported:
* `gateway_id` - (Required, ForceNew) The ID of the gateway instance.
* `k8s_source_info` - (Optional, ForceNew, Set) The ACK cluster source information. Required when `type` is `K8S`. See [`k8s_source_info`](#k8s_source_info) below.
* `nacos_source_info` - (Optional, ForceNew, Set) The MSE Nacos source information. Required when `type` is `MSE_NACOS`. See [`nacos_source_info`](#nacos_source_info) below.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `type` - (Required, ForceNew) The type of the source. Valid values: `K8S` (Container Service for Kubernetes), `MSE_NACOS` (MSE Nacos).

### `k8s_source_info`

The k8s_source_info supports the following:
* `cluster_id` - (Optional, ForceNew) The ID of the ACK cluster.

### `nacos_source_info`

The nacos_source_info supports the following:
* `instance_id` - (Optional, ForceNew) The ID of the MSE Nacos instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `association_reason` - The reason for the association status.
* `association_status` - The association status of the source.
* `create_time` - The creation timestamp of the source.
* `nacos_source_info` - The MSE Nacos source information.
  * `address` - The access address of the Nacos instance.
  * `cluster_id` - The ID of the Nacos cluster.
* `source_name` - The name of the source.
* `update_time` - The update timestamp of the source.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax/operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Source.
* `delete` - (Defaults to 5 mins) Used when delete the Source.
* `update` - (Defaults to 6 mins) Used when update the Source.

## Import

APIG Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_source.example <source_id>
```
