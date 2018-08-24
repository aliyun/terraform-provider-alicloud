---
layout: "alicloud"
page_title: "Alicloud: alicloud_eips"
sidebar_current: "docs-alicloud-datasource-eips"
description: |-
    Provides a list of EIP owned by an Alibaba Cloud account.
---

# alicloud\_eips

This data source provides a list of EIPs (Elastic IP address) owned by an Alibaba Cloud account.

## Example Usage

```
data "alicloud_eips" "eips_ds" {
}

output "first_eip_id" {
  value = "${data.alicloud_eips.eips_ds.eips.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of EIP IDs.
* `ip_addresses` - (Optional) A list of EIP public IP addresses.
* `in_use` - (Deprecated) Deprecated since the version 1.8.0 of this provider.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `eips` - A list of EIPs. Each element contains the following attributes:
  * `id` - ID of the EIP.
  * `status` - EIP status. Possible values are: `Associating`, `Unassociating`, `InUse` and `Available`.
  * `ip_address` - Public IP Address of the the EIP.
  * `bandwidth` - EIP internet max bandwidth in Mbps.
  * `internet_charge_type` - EIP internet charge type.
  * `instance_id` - The ID of the instance that is being bound.
  * `instance_type` - The instance type of that the EIP is bound.
  * `creation_time` - Time of creation.
