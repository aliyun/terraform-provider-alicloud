---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domain_txt_guid"
sidebar_current: "docs-alicloud-datasource-dns-domain-txt-guid"
description: |-
    Provides the generation of txt records to realize the retrieval and verification of domain names.
---

# alicloud\_dns\_domain\_txt\_guid

Provides the generation of txt records to realize the retrieval and verification of domain names.

-> **NOTE:** Available in v1.80.0+.

## Example Usage

```
data "alicloud_dns_domain_txt_guid" "this" {
  domain_name = "test111.abc"
  type = "ADD_SUB_DOMAIN"
}

output "rr" {
  value = "${data.alicloud_dns_domain_txt_guid.this.rr}"
}

output "value" {
  value = "${data.alicloud_dns_domain_txt_guid.this.value}"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required) Verified domain name.
* `type` - (Required) Txt verification function. Value:`ADD_SUB_DOMAIN`, `RETRIEVAL`.
* `lang` - (Optional) User language.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `rr` - Host record.
* `value` - Record the value.

  
