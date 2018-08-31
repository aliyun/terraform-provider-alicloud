---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidthpackage_attachment"
sidebar_current: "docs-alicloud-resource-cen-bandwidthpackage-attachment"
description: |-
  Provides a Alicloud CEN bandwidth package attachment resource.
---

# alicloud\_cen_bandwidthpackage_attachment

Provides a CEN bandwidth package attachment resource.

## Example Usage

Basic Usage

```
# Create a new bandwidthpackage-attachment and use it to attach one bandwidth package to a new CEN
resource "alicloud_cen" "cen" {
     name = "terraform-01"
     description = "terraform01"
}

resource "alicloud_cen_bandwidthpackage" "bwp" {
    bandwidth = 20
    geographic_region_id = [
		"China",
		"Asia-Pacific"]
}

resource "alicloud_cen_bandwidthpackage_attachment" "foo" {
    cen_id = "${alicloud_cen.cen.id}"
    cen_bandwidthpackage_id = "${alicloud_cen_bandwidthpackage.bwp.id}"
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) The ID of the CEN.
* `cen_bandwidthpackage_id` - (Required) The ID of the bandwidth package.

## Attributes Reference

The following attributes are exported:

* `cen_id` - (Required) The ID of the CEN.
* `cen_bandwidthpackage_id` - (Required) The ID of the bandwidth package.
