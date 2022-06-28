---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_ca_certificates"
sidebar_current: "docs-alicloud-datasource-slb-ca-certificates"
description: |-
    Provides a list of slb CA certificates.
---
# alicloud\_slb_ca_certificates

This data source provides the CA certificate list.

## Example Usage

```terraform
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
* `resource_group_id` - (Optional, ForceNew, Available in 1.60.0+) The Id of resource group which ca certificates belongs.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB ca certificates IDs.
* `names` - A list of SLB ca certificates names.
* `certificates` - A list of SLB ca certificates. Each element contains the following attributes:
  * `id` - CA certificate ID.
  * `ca_certificate_name` - (Available in v1.123.1+) CA certificate name.
  * `fingerprint` - CA certificate fingerprint.
  * `common_name` - CA certificate common name.
  * `expired_time` - CA certificate expired time.
  * `expired_timestamp` - CA certificate expired timestamp.
  * `created_timestamp` - CA certificate created timestamp.
  * `resource_group_id` - The resource group Id of CA certificate.
  * `tags` - (Available in v1.66.0+) A mapping of tags to assign to the resource.
  * `ca_certificate_id` - (Available in v1.123.1+) CA certificate ID.
  * `name` - (Deprecated from v1.123.1) Deprecated and replace by `ca_certificate_name`.
  * `created_time` - (Removed from v1.123.1) CA certificate created time.
  * `region_id` - (Removed from v1.123.1) The region Id of CA certificate.
