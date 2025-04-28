---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_package"
sidebar_current: "docs-alicloud-resource-cen-bandwidth-package"
description: |-
  Provides a Alicloud CEN bandwidth package resource.
---

# alicloud_cen_bandwidth_package

Provides a CEN bandwidth package resource. The CEN bandwidth package is an abstracted object that includes an interconnection bandwidth and interconnection areas. To buy a bandwidth package, you must specify the areas to connect. An area consists of one or more Alibaba Cloud regions. The areas in CEN include Mainland China, Asia Pacific, North America, and Europe.

For information about CEN and how to use it, see [Manage bandwidth packages](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createcenbandwidthpackage).

-> **NOTE:** Available since v1.18.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_bandwidth_package&exampleId=af9941dc-7b32-fad5-25ab-d125a19110be0c3a64f1&activeTab=example&spm=docs.r.cen_bandwidth_package.0.af9941dc7b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cen_bandwidth_package" "example" {
  bandwidth                  = 5
  cen_bandwidth_package_name = "tf_example"
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}
```

### Deleting `alicloud_cen_bandwidth_package` or removing it from your configuration

The `alicloud_cen_bandwidth_package` resource allows you to manage `payment_type = "PrePaid"` bandwidth package, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your statefile and management, but will not destroy the Bandwidth Package.
You can resume managing the subscription bandwidth package via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The bandwidth in Mbps of the bandwidth package. Cannot be less than 2Mbps.
* `geographic_region_ids` - (Required, Deprecated from 1.98.0+) Field `geographic_region_ids` has been deprecated from version 1.97.0. Use `geographic_region_a_id` and `geographic_region_b_id` instead.
* `name` - (Optional, Deprecated from 1.98.0+) Field `name` has been deprecated from version 1.97.0. Use `cen_bandwidth_package_name` and instead.
* `description` - (Optional) The description of the bandwidth package. Default to null.
* `charge_type` - (Optional, Deprecated from 1.98.0+) Field `charge_type` has been deprecated from version 1.97.0. Use `payment_type` and instead.
* `period` - (Optional) The purchase period in month. Valid value: `1`, `2`, `3`, `6`, `12`.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `geographic_region_a_id` - (Optional, ForceNew, Available since v1.98.0) The area A to which the network instance belongs. Valid values: `China` | `North-America` | `Asia-Pacific` | `Europe` | `Australia`.
* `geographic_region_b_id` - (Optional, ForceNew, Available since v1.98.0) The area B to which the network instance belongs. Valid values: `China` | `North-America` | `Asia-Pacific` | `Europe` | `Australia`.
* `payment_type` - (Optional, ForceNew, Available since v1.98.0) The billing method. Valid value: `PostPaid` | `PrePaid`. Default to `PrePaid`. If set to PrePaid, the bandwidth package can't be deleted before expired time.
* `cen_bandwidth_package_name` - (Optional, Available since v1.98.0) The name of the bandwidth package. Defaults to null.

->**NOTE:** PrePaid mode will deduct fees from the account directly and the bandwidth package can't be deleted before expired time. 

->**NOTE:** The PostPaid mode is only for test. Please open a ticket if you need.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the bandwidth package.
* `expired_time` - The time of the bandwidth package to expire.
* `status` - The association status of the bandwidth package.

## Timeouts

-> **NOTE:** Available since v1.98.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when creating the CEN bandwidth package. (until it reaches the initial `Idle` status).
* `delete` - (Defaults to 11 mins) Used when delete the CEN bandwidth package.

## Import

CEN bandwidth package can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_bandwidth_package.example cenbwp-abc123456
```

