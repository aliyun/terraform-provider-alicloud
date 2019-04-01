---
layout: "alicloud"
page_title: "Alicloud: alicloud_havip"
sidebar_current: "docs-alicloud-resource-havip"
description: |-
  Provides a Alicloud HaVip resource.
---

# alicloud\_havip

Provides a HaVip resource.

-> **NOTE:** Terraform will auto build havip instance  while it uses `alicloud_havip` to build a havip resource.

## Example Usage

Basic Usage

```
resource "alicloud_havip" "foo" {
    vswitch_id = "vsw-fakeid"
    description = "test_havip"
}
```
## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The vswitch_id of the HaVip, the field can't be changed.
* `ip_address` - (Optional, ForceNew) The ip address of the HaVip. If not filled, the default will be assigned one from the vswitch.
* `description` - (Optional) The description of the HaVip instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HaVip instance id.

## Import

The havip can be imported using the id, e.g.

```
$ terraform import alicloud_havip.foo havip-abc123456
```




