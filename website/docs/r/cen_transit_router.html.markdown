---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router"
sidebar_current: "docs-alicloud-resource-cen-transit_router"
description: |-
Provides a Alicloud CEN transit router resource.
---

# alicloud\_cen_transit_router

Provides a CEN transit router resource that associate the transitRouter with the CEN instance.[What is Cen Transit Router](https://help.aliyun.com/document_detail/261169.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

Basic Usage

```terraform
# Create a new tr-attachment and use it to attach one transit router to a new CEN
variable "name" {
  default = "tf-testAccCenTransitRouter"
}

resource "alicloud_cen_instance" "default" {
  name        = var.name
  description = "terraform01"
}

resource "alicloud_cen_transit_router" "default" {
  name       = var.name
  cen_id     = alicloud_cen_instance.default.id
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `type` - (Optional, Available in 1.126.0+) The Type of the Transit Router. Valid values: `Enterprise`, `Basic`.
* `transit_router_name` - (Optional) The name of the transit router. 
* `transit_router_description` - (Optional) The description of the transit router.
* `dry_run` - (Optional,ForceNew) The dry run.


## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<cen_id>:<transit_router_id>`.
* `status` - The associating status of the Transit Router.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router.default cen-*****:tr-*******
```
