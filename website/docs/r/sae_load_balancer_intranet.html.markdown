---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_load_balancer_intranet"
sidebar_current: "docs-alicloud-resource-sae-load-balancer-intranet"
description: |-
  Provides an Alicloud Serverless App Engine (SAE) Application Load Balancer Attachment resource.
---

# alicloud_sae_load_balancer_intranet

Provides an Alicloud Serverless App Engine (SAE) Application Load Balancer Attachment resource.

For information about Serverless App Engine (SAE) Load Balancer Intranet Attachment and how to use it, see [alicloud_sae_load_balancer_intranet](https://www.alibabacloud.com/help/en/sae/latest/bindslb).

-> **NOTE:** Available since v1.165.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}
data "alicloud_regions" "default" {
  current = true
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_sae_namespace" "default" {
  namespace_id              = "${data.alicloud_regions.default.regions.0.id}:example${random_integer.default.result}"
  namespace_name            = var.name
  namespace_description     = var.name
  enable_micro_registration = false
}

resource "alicloud_sae_application" "default" {
  app_description   = var.name
  app_name          = "${var.name}-${random_integer.default.result}"
  namespace_id      = alicloud_sae_namespace.default.id
  image_url         = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type      = "Image"
  jdk               = "Open JDK 8"
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
  vswitch_id        = alicloud_vswitch.default.id
  timezone          = "Asia/Beijing"
  replicas          = "5"
  cpu               = "500"
  memory            = "2048"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "intranet"
}

resource "alicloud_sae_load_balancer_intranet" "default" {
  app_id          = alicloud_sae_application.default.id
  intranet_slb_id = alicloud_slb_load_balancer.default.id
  intranet {
    protocol    = "TCP"
    port        = 80
    target_port = 8080
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The target application ID that needs to be bound to the SLB.
* `intranet_slb_id` - (Optional) The intranet SLB ID.
* `intranet` - (Required) The bound private network SLB. See [`intranet`](#intranet) below.

### `intranet`

The intranet supports the following:

* `protocol` - (Optional) The Network protocol. Valid values: `TCP` ,`HTTP`,`HTTPS`.
* `https_cert_id` - (Optional) The SSL certificate. `https_cert_id` is required when HTTPS is selected
* `target_port` - (Optional) The Container port.
* `port` - (Optional) The SLB Port.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is the same as the application ID.
* `intranet_ip` - Use designated private network SLBs that have been purchased to support non-shared instances.

## Import

The resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_load_balancer_intranet.example <id>
```
