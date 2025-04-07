---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_ingress"
sidebar_current: "docs-alicloud-resource-sae-ingress"
description: |-
  Provides a Alicloud Serverless App Engine (SAE) Ingress resource.
---

# alicloud_sae_ingress

Provides a Serverless App Engine (SAE) Ingress resource.

For information about Serverless App Engine (SAE) Ingress and how to use it, see [What is Ingress](https://next.api.aliyun.com/api/sae/2019-05-06/CreateIngress).

-> **NOTE:** Available since v1.137.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sae_ingress&exampleId=5e63c4a7-8223-82c6-ffe6-8714e3e43c9adfb5667e&activeTab=example&spm=docs.r.sae_ingress.0.5e63c4a782&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  image_url         = "registry-vpc.${data.alicloud_regions.default.regions.0.id}.aliyuncs.com/sae-demo-image/consumer:1.0"
  package_type      = "Image"
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

resource "alicloud_sae_ingress" "default" {
  slb_id        = alicloud_slb_load_balancer.default.id
  namespace_id  = alicloud_sae_namespace.default.id
  listener_port = "80"
  rules {
    app_id         = alicloud_sae_application.default.id
    container_port = "443"
    domain         = "www.alicloud.com"
    app_name       = alicloud_sae_application.default.app_name
    path           = "/"
  }
  default_rule {
    app_id         = alicloud_sae_application.default.id
    container_port = "443"
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace_id` - (Required, ForceNew) The ID of Namespace. It can contain 2 to 32 lowercase characters.The value is in format `{RegionId}:{namespace}`.
* `slb_id` - (Required, ForceNew) SLB ID.
* `listener_port` - (Required, Int) SLB listening port.
* `cert_id` - (Optional) The certificate ID of the HTTPS listener. The `cert_id` takes effect only when `load_balance_type` is set to `clb`.
* `cert_ids` - (Optional, Available since v1.207.0) The certificate IDs of the HTTPS listener, and multiple certificate IDs are separated by commas. The `cert_ids` takes effect only when `load_balance_type` is set to `alb`.
* `load_balance_type` - (Optional, Available since v1.207.0) The type of the SLB instance. Default value: `clb`. Valid values: `clb`, `alb`.
* `listener_protocol` - (Optional, Available since v1.207.0) The protocol that is used to forward requests. Default value: `HTTP`. Valid values: `HTTP`, `HTTPS`.
* `description` - (Optional) Description.
* `rules` - (Required, Set) Forwarding rules. Forward traffic to the specified application according to the domain name and path. See [`rules`](#rules) below.
* `default_rule` - (Optional, Set) Default Rule. See [`default_rule`](#default_rule) below.

### `default_rule`

The default_rule supports the following:

* `app_id` - (Optional) Target application ID.
* `app_name` - (Optional) Target application name.
* `container_port` - (Optional, Int) Application backend port.

### `rules`

The rules supports the following:

* `app_id` - (Required) Target application ID.
* `app_name` - (Required) Target application name.
* `container_port` - (Required, Int) Application backend port.
* `domain` - (Required) Application domain name.
* `path` - (Required) URL path.
* `rewrite_path` - (Optional, Available since v1.207.0) The rewrite path.
* `backend_protocol` - (Optional, Available since v1.207.0) The backend protocol.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ingress.

## Import

Serverless App Engine (SAE) Ingress can be imported using the id, e.g.

```shell
$ terraform import alicloud_sae_ingress.example <id>
```
