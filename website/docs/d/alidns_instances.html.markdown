---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_instances"
sidebar_current: "docs-alicloud-datasource-alidns-instances"
description: |-
    Provides a list of instances available to the Alidns.
---

# alicloud\_alidns\_instances

This data source provides a list of Alidns instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:**  Available in 1.95.0+.

## Example Usage

```terraform
data "alicloud_alidns_instances" "example" {
  ids = ["dns-cn-oew1npk****"]
}
output "first_instance_id" {
  value = "${data.alicloud_alidns_instances.example.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs.
* `lang` - (Optional) Language.
* `user_client_ip` - (Optional) The IP address of the client. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `domain_type` - (Optional, Available in 1.124.1+) The type of domain.

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
  * `domain` - (Available in 1.124.1+) The domain name.
  * `payment_type` - (Available in 1.124.1+) The payment type of alidns instance.
