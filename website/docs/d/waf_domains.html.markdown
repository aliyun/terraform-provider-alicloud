---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_domains"
sidebar_current: "docs-alicloud-datasource-waf-domains"
description: |-
  Provides a datasource to retrieve domain names.
---

# alicloud\_waf\_domains

Provides a WAF datasource to retrieve domains.

For information about WAF and how to use it, see [What is Alibaba Cloud WAF](https://www.alibabacloud.com/help/doc-detail/28517.htm).

-> **NOTE:** Available in 1.86.0+ .

## Example Usage

```
data "alicloud_waf_domains" "default" {
  instance_id = "waf-cf-xxxxx"
}
```
## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of WAF domain names. Each item is domain name.
* `instance_id` - (Required) The Id of waf instance to which waf domain belongs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of WAF domain names. Each item is domain name.
* `domains` - A list of Domains. Each element contains the following attributes:
  * `domain` - Name of the domain.
```
