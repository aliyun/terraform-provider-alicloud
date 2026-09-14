---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_domain"
sidebar_current: "docs-alicloud-datasource-api-gateway-domain"
description: |-
  Provides an Alicloud Api Gateway Domain data source.
---

# alicloud_api_gateway_domain

Provides an API Gateway domain data source that queries the binding information of a custom domain name attached to an API Gateway group.

-> **NOTE:** Available since v1.229.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_api_gateway_domain" "default" {
  group_id    = "example_group_id"
  domain_name = "demo.example.com"
}

output "domain_id" {
  value = data.alicloud_api_gateway_domain.default.id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The ID of the API Gateway group that the custom domain binds to.
* `domain_name` - (Required, ForceNew) The custom domain name to query.
* `output_file` - (Optional) File path where to write the result in JSON format.

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
