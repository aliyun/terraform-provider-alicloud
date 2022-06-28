---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_server_certificates"
sidebar_current: "docs-alicloud-datasource-slb-server-certificates"
description: |-
    Provides a list of slb server certificates.
---
# alicloud\_slb_server_certificates

This data source provides the server certificate list.

## Example Usage

```
data "alicloud_slb_server_certificates" "sample_ds" {
}

output "first_slb_server_certificate_id" {
  value = "${data.alicloud_slb_server_certificates.sample_ds.certificates.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of server certificates IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by server certificate name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew, Available in 1.58.0+) The Id of resource group which the slb server certificates belongs.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB server certificates IDs.
* `names` - A list of SLB server certificates names.
* `certificates` - A list of SLB server certificates. Each element contains the following attributes:
  * `id` - Server certificate ID.
  * `name` - Server certificate name.
  * `fingerprint` - Server certificate fingerprint.
  * `common_name` - Server certificate common name.
  * `subject_alternative_names` - Server certificate subject alternative name list.
  * `expired_time` - Server certificate expired time.
  * `expired_timestamp` - Server certificate expired timestamp.
  * `created_time` - Server certificate created time.
  * `created_timestamp` - Server certificate created timestamp.
  * `alicloud_certificate_id` - Id of server certificate issued by alibaba cloud.
  * `alicloud_certificate_name`- Name of server certificate issued by alibaba cloud.
  * `is_alicloud_certificate`- Is server certificate issued by alibaba cloud or not.
  * `resource_group_id` - The Id of resource group which the slb server certificates belongs.
  * `tags` - (Available in v1.66.0+) A mapping of tags to assign to the resource.
