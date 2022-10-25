---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eips"
sidebar_current: "docs-alicloud-datasource-eips"
description: |-
    Provides a list of EIP owned by an Alibaba Cloud account.
---

# alicloud\_eips

-> **DEPRECATED:**  This datasource has been deprecated from version `1.126.0`. Please use new datasource [alicloud_eip_addresses](https://www.terraform.io/docs/providers/alicloud/d/eip_addresses).

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
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.
* `in_use` - (Deprecated) Deprecated since the version 1.8.0 of this provider.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew, Available in 1.58.0+) The Id of resource group which the eips belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of EIP IDs.
* `names` - (Optional) A list of EIP names.
* `eips` - A list of EIPs. Each element contains the following attributes:
  * `id` - ID of the EIP.
  * `status` - EIP status. Possible values are: `Associating`, `Unassociating`, `InUse` and `Available`.
  * `ip_address` - Public IP Address of the the EIP.
  * `bandwidth` - EIP internet max bandwidth in Mbps.
  * `internet_charge_type` - EIP internet charge type.
  * `instance_id` - The ID of the instance that is being bound.
  * `instance_type` - The instance type of that the EIP is bound.
  * `creation_time` - Time of creation.
  * `resource_group_id` - The Id of resource group which the eips belongs.
  * `deletion_protection` - (Optional, Available in v1.124.4+) Whether enable the deletion protection or not.
