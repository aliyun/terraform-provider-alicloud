---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_ingresses"
sidebar_current: "docs-alicloud-datasource-sae-ingresses"
description: |-
  Provides a list of Sae Ingresses to the user.
---

# alicloud\_sae\_ingresses

This data source provides the Sae Ingresses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.137.0+.

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

variable "desc" {
  default = "example_value"
}
variable "namespace_id" {
  default = "cn-hangzhou:yourname"
}

resource "alicloud_sae_namespace" "default" {
  namespace_id          = var.namespace_id
  namespace_name        = var.name
  namespace_description = var.desc
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
data "alicloud_sae_ingresses" "default" {
  ids = [alicloud_sae_ingress.default.id]
}
output "sae_ingress_id" {
  value = data.alicloud_sae_ingresses.default.IngressList.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Ingress IDs.
* `namespace_id` - (Required, ForceNew) The Id of Namespace.It can contain 2 to 32 lowercase characters.The value is in format `{RegionId}:{namespace}`
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ingresses` - A list of Sae Ingresses. Each element contains the following attributes:
	* `cert_id` - Cert Id.
	* `default_rule` - Default Rule.
	* `description` - Description.
	* `id` - The ID of the Ingress.
	* `ingress_id` - The first ID of the resource.
	* `listener_port` - SLB listening port.
	* `namespace_id` - The Id of Namespace.It can contain 2 to 32 characters.The value is in format {RegionId}:{namespace}.
	* `slb_id` - SLB ID.
