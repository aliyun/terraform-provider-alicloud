---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_attachments"
sidebar_current: "docs-alicloud-datasource-slb-attachments"
description: |-
    Provides a list of server load balancer attachments to the user.
---

# alicloud\_slb_attachments

This data source provides the server load balancer attachments of the current Alibaba Cloud user.

## Example Usage

```
data "alicloud_slb_attachments" "sample_ds" {
  load_balancer_id = "${alicloud_slb_load_balancer.sample_slb.id}"
}

output "first_slb_attachment_instance_id" {
  value = "${data.alicloud_slb_attachments.sample_ds.slb_attachments.0.instance_id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - ID of the SLB with attachments.
* `instance_ids` - (Optional) List of attached ECS instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `slb_attachments` - A list of SLB attachments. Each element contains the following attributes:
  * `instance_id` - ID of the attached ECS instance.
  * `weight` - Weight associated to the ECS instance.
