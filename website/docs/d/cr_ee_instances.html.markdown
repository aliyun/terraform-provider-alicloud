---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_instances"
sidebar_current: "docs-alicloud-datasource-cr-ee-instances"
description: |-
  Provides a list of Cr Instance owned by an Alibaba Cloud account.
---

# alicloud_cr_ee_instances

This data source provides Cr Instance available to the user.[What is Instance](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}


resource "alicloud_cr_ee_instance" "default" {
  default_oss_bucket = "true"
  instance_name      = "terraform-eco-526"
  renewal_status     = "ManualRenewal"
  image_scanner      = "DISABLE"
  period             = 1
  payment_type       = "Subscription"
  instance_type      = "Economy"
}

data "alicloud_cr_ee_instances" "default" {
  ids           = ["${alicloud_cr_ee_instance.default.id}"]
  name_regex    = alicloud_cr_ee_instance.default.instance_name
  instance_name = "terraform-eco-526"
}

output "alicloud_cr_ee_instance_example_id" {
  value = data.alicloud_cr_ee_instances.default.instances.0.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_name` - (ForceNew, Optional) InstanceName
* `resource_group_id` - (ForceNew, Optional) The ID of the resource group
* `ids` - (Optional, Computed) A list of Instance IDs.
* `name_regex` - (Optional) A regex string to filter results by Group Metric Rule name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Instance IDs.
* `names` - A list of name of Instances.
* `instances` - A list of Instance Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the resource.
  * `end_time` - **NOTE:** This field is only available when `enable_details` is `true`. Expiration Time.
  * `instance_endpoints` - **NOTE:** This field is only available when `enable_details` is `true`. Instance Network Access Endpoint List.
    * `domains` - Domain List.
      * `domain` - Domain.
      * `type` - Domain Type.
    * `enable` - enable.
    * `endpoint_type` - Network Access Endpoint Type.
  * `instance_id` - Instance ID.
  * `instance_issue` - Instance issue.
  * `instance_name` - InstanceName.
  * `modified_time` - Last modification time.
  * `payment_type` - **NOTE:** This field is only available when `enable_details` is `true`. Payment type, value:.
  * `region_id` - RegionId.
  * `renew_period` - **NOTE:** This field is only available when `enable_details` is `true`. Automatic renewal cycle, in months.
  * `renewal_status` - **NOTE:** This field is only available when `enable_details` is `true`. Automatic renewal status, value:.
  * `resource_group_id` - The ID of the resource group.
  * `status` - **NOTE:** This field is only available when `enable_details` is `true`. Instance Status.
  * `id` - The ID of the resource supplied above.
