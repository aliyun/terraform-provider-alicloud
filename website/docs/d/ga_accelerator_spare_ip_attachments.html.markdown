---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_accelerator_spare_ip_attachments"
sidebar_current: "docs-alicloud-datasource-ga-accelerator-spare-ip-attachments"
description: |-
  Provides a list of Ga Accelerator Spare Ip Attachments to the user.
---

# alicloud\_ga\_accelerator\_spare\_ip\_attachments

This data source provides the Ga Accelerator Spare Ip Attachments of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.167.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_accelerator_spare_ip_attachments" "ids" {
  accelerator_id = "example_value"
  ids            = ["example_value-1", "example_value-2"]
}
output "ga_accelerator_spare_ip_attachment_id_1" {
  value = data.alicloud_ga_accelerator_spare_ip_attachments.ids.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the global acceleration instance.
* `ids` - (Optional, ForceNew, Computed)  A list of Accelerator Spare Ip Attachment IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the standby CNAME IP address. Valid values: `active`, `inuse`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `attachments` - A list of Ga Accelerator Spare Ip Attachments. Each element contains the following attributes:
	* `accelerator_id` - The ID of the global acceleration instance.
	* `id` - The ID of the Accelerator Spare Ip Attachment.
	* `spare_ip` - The standby IP address of CNAME. When the acceleration area is abnormal, the traffic is switched to the standby IP address.
	* `status` - The status of the standby CNAME IP address. Valid values: `active`, `inuse`.