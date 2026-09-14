---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_load_balancer_nlb"
sidebar_current: "docs-alicloud-resource-sae-load-balancer-nlb"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Load Balancer NLB resource.
---

# alicloud_sae_load_balancer_nlb

Provides a Serverless App Engine (SAE) Load Balancer NLB resource.

Bind an NLB (Network Load Balancer) instance to an SAE application to configure NLB access. The resource supports binding an existing NLB instance by its ID or letting SAE create and manage the NLB on your behalf.

For information about SAE NLB and how to use it, see [Bind an NLB to an application](https://www.alibabacloud.com/help/en/sae/developer-reference/api-sae-2019-05-06-bindnlb).

-> **NOTE:** Available since v1.293.0.

## Example Usage

Bind an existing NLB to an SAE application

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_sae_namespace" "default" {
  namespace_id          = "${data.alicloud_zones.default.zones.0.id}:example"
  namespace_name        = var.name
  namespace_description = var.name
}

resource "alicloud_sae_application" "default" {
  app_name     = var.name
  namespace_id = alicloud_sae_namespace.default.namespace_id
  image_url    = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type = "Image"
  cpu          = "500"
  memory       = "2048"
  replicas     = "2"
  vswitch_id   = data.alicloud_vswitches.default.ids.0
  vpc_id       = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  address_type       = "Internet"
  vpc_id             = data.alicloud_vpcs.default.ids.0
  zone_mappings {
    vswitch_id = data.alicloud_vswitches.default.vswitches.0.id
    zone_id    = data.alicloud_zones.default.zones.0.id
  }
}

resource "alicloud_sae_load_balancer_nlb" "default" {
  app_id       = alicloud_sae_application.default.id
  nlb_id       = alicloud_nlb_load_balancer.default.id
  address_type = "Internet"
  listeners {
    port        = 80
    target_port = 8080
    protocol    = "TCP"
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the target SAE application that the NLB is bound to. This is the primary key of the resource.
* `nlb_id` - (Optional, ForceNew) The ID of an existing NLB instance to bind to the application. If omitted, SAE will create and manage the NLB instance.
* `address_type` - (Optional, ForceNew) The network address type of the NLB. Valid values: `Internet`, `Intranet`.
* `listeners` - (Optional) A list of NLB listener configurations. See [`listeners`](#listeners) below.
* `zone_mappings` - (Optional) A list of zone mapping configurations. See [`zone_mappings`](#zone_mappings) below.

### `listeners`

The listeners supports the following:
* `port` - (Optional) The NLB listening port.
* `target_port` - (Optional) The container port that the traffic is forwarded to.
* `protocol` - (Optional) The protocol of the listener. Valid values: `TCP`, `UDP`, `TCPSSL`.
* `cert_ids` - (Optional) The certificate IDs for the TCPSSL protocol.

### `zone_mappings`

The zone_mappings supports the following:
* `vswitch_id` - (Optional) The vSwitch ID.
* `zone_id` - (Optional) The zone ID.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which is the app ID.
* `dns_name` - The DNS name of the NLB instance.
* `created_by_sae` - Whether the NLB instance was created by SAE.

## Import

SAE Load Balancer NLB can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_load_balancer_nlb.example <app_id>
```
