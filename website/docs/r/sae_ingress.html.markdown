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

For information about Serverless App Engine (SAE) Ingress and how to use it, see [What is Ingress](https://www.alibabacloud.com/help/en/sae/latest/createingress).

-> **NOTE:** Available since v1.137.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_slb" "default" {
  name          = var.name
  specification = "slb.s2.small"
  vswitch_id    = data.alicloud_vswitches.default.ids.0
}

variable "namespace_id" {
  default = "cn-hangzhou:yourname"
}

resource "alicloud_sae_namespace" "default" {
  namespace_id          = var.namespace_id
  namespace_name        = var.name
  namespace_description = var.name
}

resource "alicloud_sae_application" "default" {
  app_description = "your_app_description"
  app_name        = "your_app_name"
  namespace_id    = "your_namespace_id"
  package_url     = "your_package_url"
  package_type    = "your_package_url"
  jdk             = "jdk_specifications"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  replicas        = "your_replicas"
  cpu             = "cpu_specifications"
  memory          = "memory_specifications"

}

resource "alicloud_sae_ingress" "default" {
  slb_id        = alicloud_slb.default.id
  namespace_id  = alicloud_sae_namespace.default.id
  listener_port = "your_listener_port"
  rules {
    app_id         = alicloud_sae_application.default.id
    container_port = "your_container_port"
    domain         = "your_domain"
    app_name       = "your_name"
    path           = "your_path"
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
* `load_balance_type` - (Optional, Computed, Available since v1.207.0) The type of the SLB instance. Default value: `clb`. Valid values: `clb`, `alb`.
* `listener_protocol` - (Optional, Computed, Available since v1.207.0) The protocol that is used to forward requests. Default value: `HTTP`. Valid values: `HTTP`, `HTTPS`.
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
