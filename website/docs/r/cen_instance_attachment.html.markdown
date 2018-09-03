---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance_attachment"
sidebar_current: "docs-alicloud-cen-instance-attachment"
description: |-
  Provides a Alicloud CEN child instance attachment resource.
---

# alicloud\_cen_instance_attachment

Provides a CEN child instance attachment resource.

## Example Usage

Basic Usage

```
# Create a new instance-attachment and use it to attach one child instance to a new CEN
resource "alicloud_cen" "cen" {
	name = "terraform-01"
	description = "terraform01"
}

resource "alicloud_vpc" "vpc" {
	name = "terraform-01"
	cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_instance_attachment" "foo" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) The ID of the CEN.
* `child_instance_id` - (Required) The ID of the child instance to attach.
* `child_instance_type` - (Required) The type of the child instance to attach. Valid value: VPC | VBR.
* `child_instance_region_id` - (Required) The region ID of the child instance to attach.

~>**NOTE:** Ensure that the child instance is not used in Express Connect.

## Attributes Reference

The following attributes are exported:

* `cen_id` - (Required) The ID of the CEN.
* `child_instance_id` - (Required) The ID of the child instance to attach.
* `child_instance_type` - (Required) The type of the child instance to attach. Valid value: VPC | VBR.
* `child_instance_region_id` - (Required) The region ID of the child instance to attach.
