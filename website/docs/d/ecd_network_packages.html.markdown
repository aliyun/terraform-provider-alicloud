---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_network_packages"
sidebar_current: "docs-alicloud-datasource-ecd-network-packages"
description: |-
  Provides a list of Ecd Network Packages to the user.
---

# alicloud\_ecd\_network\_packages

This data source provides the Ecd Network Packages of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = "example_value"
}

resource "alicloud_ecd_network_package" "default" {
  bandwidth      = "10"
  office_site_id = alicloud_ecd_simple_office_site.default.id
}

data "alicloud_ecd_network_packages" "default" {
  ids = [alicloud_ecd_network_package.default.id]
}
output "ecd_network_package_id_1" {
  value = data.alicloud_ecd_network_packages.default.packages.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Network Package IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of network package. Valid values: `Creating`, `InUse`, `Releasing`,`Released`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `packages` - A list of Ecd Network Packages. Each element contains the following attributes:
	* `bandwidth` - The bandwidth of package.
	* `create_time` - The creation time of network package.
	* `expired_time` - The expired time of package.
	* `id` - The ID of the Network Package.
	* `internet_charge_type` - The internet charge type  of  package.
	* `network_package_id` - The ID of network package.
	* `office_site_id` - The ID of office site.
	* `office_site_name` - The name of office site.
	* `status` - The status of network package. Valid values: `Creating`, `InUse`, `Releasing`,`Released`.
	* `eip_addresses` - The public IP address list of the network packet.
