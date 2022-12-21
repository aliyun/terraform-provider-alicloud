---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_grant_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit-router-grant-attachment"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Grant Attachment resource.
---

# alicloud\_cen\_transit\_router\_grant\_attachment

Provides a Cloud Enterprise Network (CEN) Transit Router Grant Attachment resource.

For information about Cloud Enterprise Network (CEN) Transit Router Grant Attachment and how to use it, see [What is Transit Router Grant Attachment](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/grantinstancetotransitrouter).

-> **NOTE:** Available in v1.187.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = "test for transit router grant attachment"
}

resource "alicloud_cen_transit_router_grant_attachment" "default" {
  cen_id        = alicloud_cen_instance.default.id
  cen_owner_id  = "your_cen_owner_id"
  instance_id   = data.alicloud_vpcs.default.ids.0
  instance_type = "VPC"
  order_type    = "PayByCenOwner"
}
```

## Argument Reference

The following arguments are supported:
* `cen_id`- (Required, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance to which the transit router belongs.
* `cen_owner_id` - (Required, ForceNew) The ID of the Alibaba Cloud account to which the CEN instance belongs.
* `instance_id` - (Required, ForceNew) The ID of the network instance.
* `instance_type` - (Required, ForceNew) The type of the network instance. Valid values: `VPC`, `ExpressConnect`, `VPN`.
* `order_type` - (Optional, Computed, ForceNew) The entity that pays the fees of the network instance. Valid values: `PayByResourceOwner`, `PayByCenOwner`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Transit Router Grant Attachment. The value formats as `<instance_type>:<instance_id>:<cen_owner_id>:<cen_id>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the cen transit router grant attachment.
* `delete` - (Defaults to 1 mins) Used when delete the cen transit router grant attachment.


## Import

Cloud Enterprise Network (CEN) Transit Router Grant Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_grant_attachment.example <instance_type>:<instance_id>:<cen_owner_id>:<cen_id>
```