---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerators"
sidebar_current: "docs-alicloud-datasource-ga-basic-accelerators"
description: |-
  Provides a list of Global Accelerator (GA) Basic Accelerators to the user.
---

# alicloud_ga_basic_accelerators

This data source provides the Global Accelerator (GA) Basic Accelerators of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_basic_accelerators" "ids" {
  ids = ["example_id"]
}

output "ga_basic_accelerator_id_1" {
  value = data.alicloud_ga_basic_accelerators.ids.accelerators.0.id
}

data "alicloud_ga_basic_accelerators" "nameRegex" {
  name_regex = "tf-example"
}

output "ga_basic_accelerator_id_2" {
  value = data.alicloud_ga_basic_accelerators.nameRegex.accelerators.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Global Accelerator Basic Accelerator IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Global Accelerator Basic Accelerator name.
* `accelerator_id` - (Optional, ForceNew) The ID of the Global Accelerator Basic Accelerator instance.
* `status` - (Optional, ForceNew) The status of the Global Accelerator Basic Accelerator instance. Valid Value: `init`, `active`, `configuring`, `binding`, `unbinding`, `deleting`, `finacialLocked`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Global Accelerator Basic Accelerator names.
* `accelerators` - A list of Global Accelerator Basic Accelerators. Each element contains the following attributes:
  * `id` - The id of the Global Accelerator Basic Accelerator.
  * `basic_accelerator_id` - The id of the Global Accelerator Basic Accelerator instance.
  * `basic_accelerator_name` - The name of the Global Accelerator Basic Accelerator instance.
  * `basic_endpoint_group_id` - The ID of the endpoint group that is associated with the Global Accelerator Basic Accelerator instance.
  * `basic_ip_set_id` - The ID of the acceleration region.
  * `bandwidth_billing_type` - The bandwidth billing method.
  * `instance_charge_type` - The billing method of the Global Accelerator Basic Accelerator instance. Only `PREPAY` is returned, which indicates the subscription billing method.
  * `description` - The description of the Global Accelerator Basic Accelerator instance.
  * `region_id` - The ID of the region where the Global Accelerator Basic Accelerator instance is deployed.
  * `create_time` - The timestamp that indicates when the Global Accelerator Basic Accelerator instance was created.
  * `expired_time` - The timestamp that indicates when the Global Accelerator Basic Accelerator instance was expired.
  * `status` - The status of the Global Accelerator Basic Accelerator instance.
  * `basic_bandwidth_package` - The details about the basic bandwidth plan that is associated with the Global Accelerator Basic Accelerator instance.
    * `instance_id` - The ID of the basic bandwidth plan.
    * `bandwidth` - The bandwidth value of the basic bandwidth plan. Unit: Mbit/s.
    * `bandwidth_type` - The type of the bandwidth that is provided by the basic bandwidth plan.
  * `cross_domain_bandwidth_package` - The details about the cross-region acceleration bandwidth plan that is associated with the Global Accelerator Basic Accelerator instance. **NOTE:** This array is returned only for Global Accelerator Basic Accelerator instances that are created on the International site.
    * `instance_id` - The ID of the cross-region acceleration bandwidth plan.
    * `bandwidth` - The bandwidth value of the cross-region acceleration bandwidth plan. Unit: Mbit/s.
