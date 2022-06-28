---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_certificates"
sidebar_current: "docs-alicloud-datasource-waf-certificates"
description: |-
  Provides a list of Web Application Firewall Certificates to the user.
---

# alicloud\_waf\_certificates

This data source provides the Waf Certificates of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_waf_certificates" "default" {
  ids         = ["your_certificate_id"]
  instance_id = "your_instance_id"
  domain      = "your_domain_name"
}
output "waf_certificate" {
  value = data.alicloud_waf_certificates.default.certificates.0
}

```

## Argument Reference

The following arguments are supported:

* `domain` - (Optional, ForceNew) WAF domain name.
* `ids` - (Optional, ForceNew, Computed)  A list of Certificate IDs.
* `instance_id` - (Required, ForceNew) WAF instance ID.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Certificate name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Certificate names.
* `certificates` - A list of Waf Certificates. Each element contains the following attributes:
    * `certificate_id` - Certificate recording ID.
    * `certificate_name` - Your certificate name.
    * `common_name` - Certificate bound to the domain name.
    * `domain` - The domain that you want to add to WAF.
    * `id` - The ID of the Certificate.
    * `instance_id` - WAF instance ID.
	
