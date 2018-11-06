---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_ca_certificates"
sidebar_current: "docs-alicloud-datasource-slb-ca-certificates"
description: |-
    Provides a list of slb CA certificates.
---
# alicloud\_slb_ca_certificates

This data source provides the CA certificate list.

## Example Usage

```
data "alicloud_slb_ca_certificates" "sample_ds" {
}

output "first_slb_ca_certificate_id" {
  value = "${data.alicloud_slb_ca_certificates.sample_ds.certificates.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ca certificates IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by ca certificate name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `certificates` - A list of SLB ca certificates. Each element contains the following attributes:
  * `id` - CA certificate ID.
  * `name` - CA certificate name.
  * `fingerprint` - CA certificate fingerprint.
  * `common_name` - CA certificate common name.
  * `expired_time` - CA certificate expired time.
  * `expired_timestamp` - CA certificate expired timestamp.
  * `created_time` - CA certificate created time.
  * `created_timestamp` - CA certificate created timestamp.
  * `resource_group_id` - The resource group Id of CA certificate.
  * `region_id` - The region Id of CA certificate.

