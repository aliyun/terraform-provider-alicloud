---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_domain"
sidebar_current: "docs-alicloud-resource-api-gateway-domain"
description: |-
  Provides an Alicloud Api Gateway Domain resource.
---

# alicloud_api_gateway_domain

Provides an API Gateway domain resource that binds a custom domain name to an API Gateway group.

-> **NOTE:** Available since v1.229.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_api_gateway_group" "default" {
  name        = var.name
  description = "demo group"
}

resource "alicloud_api_gateway_domain" "default" {
  group_id    = alicloud_api_gateway_group.default.id
  domain_name = "demo.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The ID of the API Gateway group that the custom domain binds to.
* `domain_name` - (Required, ForceNew) The custom domain name to bind to the API Gateway group.
* `ssl_ocsp_enable` - (Optional) Whether to enable SSL OCSP stapling for the bound certificate of the domain.
* `ssl_ocsp_cache_enable` - (Optional) Whether to enable SSL OCSP cache for the domain.
* `client_cert_s_dn_pass_through` - (Optional) Whether to pass through the client certificate Subject DN.
* `wss_enable` - (Optional) Whether to enable WebSocket for the domain. Valid values: `ON`, `OFF`.

-> **NOTE:** `ssl_ocsp_enable`, `ssl_ocsp_cache_enable`, `client_cert_s_dn_pass_through` and `wss_enable` are write-only attributes: they are sent to the backing API when the resource is updated, and are not returned by the read API, so their values are preserved in Terraform state across refresh.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value formats as `<group_id>:<domain_name>`.
* `sub_domain` - The system allocated sub domain of the API Gateway group.
* `domain_binding_status` - The binding status of the domain.
* `domain_remark` - The remark of the domain.
* `domain_web_socket_status` - The WebSocket status of the domain.
* `domain_legal_status` - The legal status of the domain.
* `domain_cname_status` - The CNAME status of the domain.
* `certificate_id` - The ID of the bound SSL certificate.
* `certificate_name` - The name of the bound SSL certificate.
* `certificate_valid_start` - The valid start timestamp of the bound SSL certificate.
* `certificate_valid_end` - The valid end timestamp of the bound SSL certificate.

## Import

API Gateway domain can be imported using the id, which consists of group_id and domain_name, e.g.

```shell
$ terraform import alicloud_api_gateway_domain.example <group_id>:<domain_name>
```
