---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_domain"
description: |-
  Provides a Alicloud Api Gateway Domain resource.
---

# alicloud_api_gateway_domain

Provides an Api Gateway custom domain resource that binds a custom domain to an API group. After binding, the custom domain resolves to the system-assigned subdomain (`sub_domain`) of the group so that APIs in the group can be accessed through the custom domain.

For information about Api Gateway Domain and how to use it, see [SetDomain](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-setdomain).

-> **NOTE:** Available since v2.3.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_api_gateway_group" "default" {
  name = var.name
}

resource "alicloud_api_gateway_domain" "default" {
  group_id           = alicloud_api_gateway_group.default.id
  domain_name        = "api.demo.com"
  custom_domain_type = "INTERNET"
  bind_stage_name    = "RELEASE"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The ID of the API group that the custom domain binds to. The API group ID is system-generated and globally unique.
* `domain_name` - (Required, ForceNew) The custom domain to bind to the API group.
* `bind_stage_name` - (Optional) The stage the custom domain can access. Valid values: `TEST`, `PRE`, `RELEASE`. When not set or set to an empty string, the domain can access all stages.
* `custom_domain_type` - (Optional) The type of the custom domain. Valid values: `INTERNET` (public) and `INTRANET` (intranet/VPC). Defaults to `INTERNET`.
* `is_http_redirect_to_https` - (Optional) Whether to redirect HTTP requests to HTTPS.
* `is_force` - (Optional) Whether to force binding the domain when a conflict exists. When `true`, the domain is force-bound to the current group and unbound from any conflicting group. Use with caution.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format `<group_id>:<domain_name>`.
* `sub_domain` - The system-assigned subdomain of the group, used for testing API calls.
* `domain_binding_status` - The binding status of the custom domain. Valid values: `BINDING` (normal), `BOUND` (not yet in effect).
* `domain_legal_status` - The legal status of the domain. Valid values: `NORMAL`, `ABNORMAL`.
* `domain_websocket_status` - The WebSocket status of the domain.
* `domain_remark` - The remarks of the domain, such as the reason for an abnormal status.
* `domain_cname_status` - The CNAME resolution status. Valid values: `RESOLVED`, `UNRESOLVED`.
* `certificate_id` - The ID of the bound SSL certificate.
* `certificate_name` - The name of the bound SSL certificate.
* `certificate_valid_start` - The start time of the certificate validity period (timestamp in milliseconds).
* `certificate_valid_end` - The end time of the certificate validity period (timestamp in milliseconds).

## Import

Api Gateway Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_domain.example <group_id>:<domain_name>
```
