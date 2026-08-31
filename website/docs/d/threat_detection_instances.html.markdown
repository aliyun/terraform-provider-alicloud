---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_instances"
sidebar_current: "docs-alicloud-datasource-threat_detection-instances"
description: |-
  Provides a list of Threat Detection Instance owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_instances

This data source provides Threat Detection Instance available to the user.[What is Instance](https://www.alibabacloud.com/help/en/security-center/latest/what-is-security-center)

-> **NOTE:** Available in 1.199.0+

## Example Usage

```terraform
data "alicloud_threat_detection_instances" "default" {
  ids = ["${alicloud_threat_detection_instance.default.id}"]
}

output "alicloud_threat_detection_instance_example_id" {
  value = data.alicloud_threat_detection_instances.default.instances.0.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (ForceNew,Optional) The first ID of the resource
* `ids` - (Optional, ForceNew, Computed) A list of Instance IDs.
* `renew_status` - (Optional, ForceNew) The renewal status of the specified instance. Valid values: `AutoRenewal`, `ManualRenewal`, `NotRenewal`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Instance IDs.
* `instances` - A list of Instance Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the resource
  * `instance_id` - The first ID of the resource
  * `payment_type` - The payment type of the resource.
  * `status` - The status of the resource.
  * `id` - ID of the instance.
