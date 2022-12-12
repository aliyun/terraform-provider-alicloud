---
subcategory: "VPC"
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

```terraform
resource "alicloud_havip" "foo" {
  vswitch_id  = "vsw-fakeid"
  description = "test_havip"
}
```
## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The vswitch_id of the HaVip, the field can't be changed.
* `ip_address` - (Optional, ForceNew) The ip address of the HaVip. If not filled, the default will be assigned one from the vswitch.
* `description` - (Optional) The description of the HaVip instance.
* `havip_name` - (Optional, Available in v1.120.0+) The name of the HaVip instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HaVip instance id.
* `status` - (Available in v1.120.0+) The status of the HaVip instance.

#### Timeouts

-> **NOTE:** Available in v1.120.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the HaVip instance.
* `update` - (Defaults to 5 mins) Used when updating the HaVip instance.
* `delete` - (Defaults to 5 mins) Used when deleting the HaVip instance.

## Import

The havip can be imported using the id, e.g.

```shell
$ terraform import alicloud_havip.foo havip-abc123456
```




