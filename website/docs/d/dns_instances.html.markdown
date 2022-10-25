---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_instances"
sidebar_current: "docs-alicloud-datasource-dns-instances"
description: |-
    Provides a list of instances available to the DNS.
---

# alicloud\_dns\_instances

-> **DEPRECATED:**  This resource has been renamed to [alicloud_alidns_instances](https://www.terraform.io/docs/providers/alicloud/d/alidns_instances) from version 1.95.0. 

This data source provides a list of DNS instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:**  Available in 1.84.0+.

## Example Usage

```terraform
data "alicloud_dns_instances" "example" {
  ids = ["dns-cn-oew1npk****"]
}
output "first_instance_id" {
  value = "${data.alicloud_dns_instances.example.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance IDs. 
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - Id of the instance.
  * `domain_numbers` - Number of domain names bound.
  * `dns_security` - DNS security level.
  * `instance_id` - Id of the instance resource.
  * `version_code` - Paid package version.
  * `version_name` - Paid package version name.
