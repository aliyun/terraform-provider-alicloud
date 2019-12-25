---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_instance_attachments"
sidebar_current: "docs-alicloud-datasource-ots-instance-attachments"
description: |-
    Provides a list of ots instance attachments to the user.
---

# alicloud\_ots\_instance\_attachments

This data source provides the ots instance attachments of the current Alibaba Cloud user.

## Example Usage

```
data "alicloud_ots_instance_attachments" "attachments_ds" {
  instance_name = "sample-instance"
  name_regex    = "testvpc"
  output_file   = "attachments.txt"
}

output "first_ots_attachment_id" {
  value = data.alicloud_ots_instance_attachments.attachments_ds.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of OTS instance.
* `name_regex` - (Optional) A regex string to filter results by vpc name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of vpc names.
* `vpc_ids` - A list of vpc ids.
* `attachments` - A list of instance attachments. Each element contains the following attributes:
  * `id` - The resource ID, the value is same as "instance_name".
  * `domain` - The domain of the instance attachment.
  * `endpoint` - The access endpoint of the instance attachment.
  * `region` - The region of the instance attachment.
  * `instance_name` - The instance name.
  * `vpc_name` - The name of attaching VPC to instance.
  * `vpc_id` - The ID of attaching VPC to instance.
	
