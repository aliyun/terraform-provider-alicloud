---
subcategory: "EIP Bandwidth Plan (CBWP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_common_bandwidth_packages"
sidebar_current: "docs-alicloud-datasource-common-bandwidth-packages"
description: |-
    Provides a list of Common Bandwidth Packages owned by an Alibaba Cloud account.
---

# alicloud\_common\_bandwidth\_packages

This data source provides a list of Common Bandwidth Packages owned by an Alibaba Cloud account.

-> **NOTE:** Available in 1.36.0+.

## Example Usage

```terraform
data "alicloud_common_bandwidth_packages" "foo" {
  name_regex = "^tf-testAcc.*"
  ids        = ["${alicloud_common_bandwidth_package.foo.id}"]
}

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth   = "2"
  name        = "tf-testAccCommonBandwidthPackage"
  description = "tf-testAcc-CommonBandwidthPackage"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Common Bandwidth Packages IDs.
* `name_regex` - (Optional) A regex string to filter results by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew, Available in 1.58.0+) The Id of resource group which the common bandwidth package belongs.
* `bandwidth_package_name` - (Optional, ForceNew, Available in 1.120.0+) The name of bandwidth package.
* `dry_run` - (Optional, ForceNew, Available in 1.120.0+) Specifies whether to precheck only the request.
* `include_reservation_data` - (Optional, ForceNew, Available in 1.120.0+) Specifies whether to return data of orders that have not taken effect.
* `status` - (Optional, ForceNew, Available in 1.120.0+) The status of bandwidth package. Valid values: `Available` and `Pending`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of Common Bandwidth Packages IDs.
* `names` - A list of Common Bandwidth Packages names.
* `packages` - A list of Common Bandwidth Packages. Each element contains the following attributes:
  * `id` - ID of the Common Bandwidth Package.
  * `bandwidth` - The peak bandwidth of the Internet Shared Bandwidth instance.
  * `status` - Status of the Common Bandwidth Package.
  * `name` - Name of the Common Bandwidth Package.
  * `description` - The description of the Common Bandwidth Package instance.
  * `business_status` - The business status of the Common Bandwidth Package instance.
  * `isp` - ISP of the Common Bandwidth Package.
  * `creation_time` - (Deprecated from v1.120.0) Time of creation.
  * `public_ip_addresses` - Public ip addresses that in the Common Bandwidth Pakcage.
  * `resource_group_id` - The Id of resource group which the common bandwidth package belongs.
  * `bandwidth_package_id` - The resource ID of bandwidth package.
  * `bandwidth_package_name` - The name of bandwidth package.
  * `deletion_protection` - The deletion protection of bandwidth package.
  * `expired_time` - The expired time of bandwidth package.
  * `has_reservation_data` - Is has reservation data.
  * `internet_charge_type` - The internet charge type of bandwidth package.
  * `payment_type` - The payment type of bandwidth package.
  * `ratio` - The ratio of bandwidth package.
  * `reservation_active_time` - The active time of reservation.
  * `reservation_bandwidth` - The bandwidth of reservation.
  * `reservation_internet_charge_type` - The charge type of reservation internet.
  * `reservation_order_type` - The type of reservation order.
  * `service_managed` - The service managed.
  
## Public ip addresses Block
  
  The public ip addresses mapping supports the following:
  
  * `ip_address`   - The address of the EIP.
  * `allocation_id` - The ID of the EIP instance.
  * `bandwidth_package_ip_relation_status` - The IP relation status of bandwidth package.
