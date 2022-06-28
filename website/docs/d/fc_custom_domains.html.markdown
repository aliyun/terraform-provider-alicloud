---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_custom_domains"
sidebar_current: "docs-alicloud-datasource-fc-services"
description: |-
    Provides a list of FC custom domains to the user.
---

# alicloud\_fc_custom_domains

This data source provides the Function Compute custom domains of the current Alibaba Cloud user.

-> **NOTE:** Available in 1.98.0+

## Example Usage

```terraform
data "alicloud_fc_custom_domains" "fc_domains" {
  name_regex = "sample_fc_custom_domain"
}

output "first_fc_custom_domain_name" {
  value = data.alicloud_fc_custom_domains.fc_domains_ds.domains.0.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `ids` (Optional) - A list of functions ids.
* `name_regex` - (Optional) A regex string to filter results by Function Compute custom domain name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of custom domain ids.
* `names` - A list of custom domain names.
* `domains` - A list of custom domains, including the following attributes:
  * `id` - The custom domain id, same as domain name.
  * `domain_name` - The custom domain name.
  * `protocol` - The custom domain protocol.
  * `account_id` - The account id.
  * `api_version` - The API version of the Function Compute service.
  * `created_time` - The created time of the custom domain.
  * `last_modified_time` - The last modified time of the custom domain.
  * `route_config` - The configuration of domain route, mapping the path and Function Compute function.
    * `path` - The path that requests are routed from.
    * `service_name` - The name of the Function Compute service that requests are routed to. 
    * `function_name` - The name of the Function Compute function that requests are routed to.
    * `qualifier` - The version or alias of the Function Compute service that requests are routed to. For example, qualifier v1 indicates that the requests are routed to the version 1 Function Compute service.
    * `methods` - The requests of the specified HTTP methos are routed from. Valid method: GET, POST, DELETE, HEAD, PUT and PATCH. For example, "GET, HEAD" methods indicate that only requests from GET and HEAD methods are routed.
  * `cert_config` - The configuration of HTTPS certificate.
    * `cert_name` - The name of the certificate.
    * `private_key` - Private key of the HTTPS certificates, follow the 'pem' format.
    * `certificate` - Certificate data of the HTTPS certificates, follow the 'pem'.
      
