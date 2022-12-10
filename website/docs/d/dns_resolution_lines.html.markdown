---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domains"
sidebar_current: "docs-alicloud-datasource-dns-domains"
description: |-
    Provides a list of domains available to the user.
---

# alicloud\_dns\_resolution\_lines

This data source provides a list of DNS Resolution Lines in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.60.0.

## Example Usage

```
data "alicloud_dns_resolution_lines" "resolution_lines_ds" {
  line_codes = [ "cn_unicom_shanxi" ]
  output_file       = "support_lines.txt"
}

output "first_line_code" {
  value = "${data.alicloud_dns_resolution_lines.resolution_lines_ds.lines.0.line_code}"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Optional) Domain Name. 
* `line_codes` - (Optional) A list of lines codes.
* `line_display_names` - (Optional) A list of line display names.
* `user_client_ip` - (Optional) The ip of user client.
* `lang` - (Optional) language.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `line_codes` - A list of lines codes.
* `line_display_names` - A list of line display names.
* `lines` - A list of cloud resolution line. Each element contains the following attributes:
  * `line_codes` - Line code.
  * `line_display_name` - Line display name.
  * `line_name` - Line name.
