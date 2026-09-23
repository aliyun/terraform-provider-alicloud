---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance_attachments"
sidebar_current: "docs-alicloud-datasource-cen-instance-attachments"
description: |-
    Provides a list of CEN(Cloud Enterprise Network) Instance Attachments by an Alibaba Cloud Account.
---

# alicloud\_cen\_instance\_attachments

This data source provides Cen Instance Attachments of the current Alibaba Cloud User.

-> **NOTE:** Available in v1.97.0+.

## Example Usage

```terraform
data "alicloud_cen_instance_attachments" "example" {
  instance_id = "cen-o40h17ll9w********"
}

output "the_first_attachmented_instance_id" {
  value = "${data.alicloud_cen_instance_attachments.example.attachments.0.child_instance_id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the CEN instance.
* `child_instance_region_id` - (Optional, ForceNew) The region to which the network to be queried belongs.
* `child_instance_type` - (Optional, ForceNew) The type of the associated network. Valid values: `VPC`, `VBR` and `CCN`.
* `status` - (Optional, ForceNew) The status of the Cen Child Instance Attachment. Valid value: `Attaching`, `Attached` and `Aetaching`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN Instance Attachment IDs.
* `attachments` - A list of CEN Instance Attachments. Each element contains the following attributes:
  * `id` - The ID of the CEN Instance Attachment.
  * `instance_id` - The ID of the CEN instance.
  * `child_instance_attach_time` - The time when the network is associated with the CEN instance.
  * `child_instance_id` - The ID of the network.
  * `child_instance_owner_id` - The ID of the account to which the network belongs.
  * `child_instance_region_id` - The ID of the region to which the network belongs.
  * `status` - The status of the network.
  * `child_instance_type` - The type of the associated network.
