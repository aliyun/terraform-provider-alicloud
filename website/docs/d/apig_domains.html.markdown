---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_domains"
sidebar_current: "docs-alicloud-datasource-apig-domains"
description: |-
  Provides a list of Apig Domain owned by an Alibaba Cloud account.
---

# alicloud_apig_domains

This data source provides Apig Domain available to the user.[What is Domain](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateDomain)

-> **NOTE:** Available since v1.284.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_apig_domain" "default" {
  domain_name  = "example-domain.example.com"
  gateway_type = "API"
  protocol     = "HTTP"
}

data "alicloud_apig_domains" "default" {
  ids               = ["${alicloud_apig_domain.default.id}"]
  name_regex        = alicloud_apig_domain.default.domain_name
  resource_group_id = ""
}

output "alicloud_apig_domain_example_id" {
  value = data.alicloud_apig_domains.default.domains.0.id
}
```

## Argument Reference

The following arguments are supported:
* `resource_group_id` - (ForceNew, Optional) The ID of the resource group
* `ids` - (Optional, ForceNew, Computed) A list of Domain IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `enable_details` - (Optional, ForceNew) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Domain IDs.
* `names` - A list of name of Domains.
* `domains` - A list of Domain Entries. Each element contains the following attributes:
    * `ca_cert_identifier` - **NOTE:** This field is only available when `enable_details` is `true`. CA certificate identifier.
    * `cert_identifier` - tls cert identifier.
    * `client_ca_cert` - client CA certificate.
    * `domain_id` - domain id.
    * `domain_name` - domain name.
    * `domain_scope` - domain scope.
    * `force_https` - Set the HTTPS protocol type and whether to enable forced HTTPS redirection.
    * `http2_option` - **NOTE:** This field is only available when `enable_details` is `true`. Whether to enable http2 settings.
    * `m_tls_enabled` - Whether to enable mTLS mutual authentication.
    * `protocol` - Protocol, HTTP/HTTPS.
    * `resource_group_id` - The ID of the resource group.
    * `tls_cipher_suites_config` - **NOTE:** This field is only available when `enable_details` is `true`. TlsCipherSuitesConfig.
        * `config_type` - config type, Default or Custom.
        * `tls_cipher_suite` - tls Cipher Suite.
            * `name` - cipher suite name.
            * `support_versions` - support versions.
    * `tls_max` - **NOTE:** This field is only available when `enable_details` is `true`. The maximum version of the TLS protocol.
    * `tls_min` - **NOTE:** This field is only available when `enable_details` is `true`. The minimum version of the TLS protocol.
    * `id` - The ID of the resource supplied above.
