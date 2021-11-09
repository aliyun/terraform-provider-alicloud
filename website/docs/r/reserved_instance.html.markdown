---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_reserved_instance"
sidebar_current: "docs-alicloud-resource-reserved-instance"
description: |-
  Provides an ECS Reserved Instance resource.
---

# alicloud\_reserved\_instance

Provides an Reserved Instance resource.

-> **NOTE:** Available in 1.65.0+

## Example Usage

```terraform
resource "alicloud_reserved_instance" "default" {
  instance_type   = "ecs.g6.large"
  instance_amount = "1"
  period_unit     = "Year"
  offering_type   = "All Upfront"
  name            = name
  description     = "ReservedInstance"
  zone_id         = "cn-hangzhou-h"
  scope           = "Zone"
  period          = "1"
}
```

## Argument Reference

The following arguments are supported:

* `offering_type` - (Required, ForceNew) Payment type of the RI. Optional values: `No Upfront`: No upfront payment is required., `Partial Upfront`: A portion of upfront payment is required.`All Upfront`: Full upfront payment is required.
* `zone_id` - (Optional, ForceNew) ID of the zone to which the RI belongs. When Scope is set to Zone, this parameter is required. For information about the zone list, see [DescribeZones](https://www.alibabacloud.com/help/doc-detail/25610.html).
* `scope` - (Optional, ForceNew) Scope of the RI. Optional values: `Region`: region-level, `Zone`: zone-level. Default is `Region`.
* `instance_type` - (Optional, ForceNew) Instance type of the RI. For more information, see [Instance type families](https://www.alibabacloud.com/help/doc-detail/25378.html).
* `instance_amount` - (Optional, ForceNew) Number of instances allocated to an RI (An RI is a coupon that includes one or more allocated instances.).
* `Period` - (Optional, ForceNew) Term of the RI. Unit: years. Optional values: 1 and 3.
* `period_unit` - (Optional, ForceNew) Term unit. Optional value: Year.
* `resource_group_id` - (Optional, ForceNew) Resource group ID.
* `description` - (Optional) Description of the RI. 2 to 256 English or Chinese characters. It cannot start with http:// or https://.
* `name` - (Optional) Name of the RI. The name must be a string of 2 to 128 characters in length and can contain letters, numbers, colons (:), underscores (_), and hyphens. It must start with a letter. It cannot start with http:// or https://.
* `platform` - (Optional, ForceNew) The operating system type of the image used by the instance. Optional values: `Windows`, `Linux`. Default is `Linux`.

### Removing alicloud_reserved_instance from your configuration
 
The alicloud_reserved_instance resource allows you to manage your ReservedInstance, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the ReservedInstance.
 

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the ReservedInstance.

## Import

reservedInstance can be imported using id, e.g.

```
$ terraform import alicloud_reserved_instance.default ecsri-uf6df4xm0h3licit****
```

