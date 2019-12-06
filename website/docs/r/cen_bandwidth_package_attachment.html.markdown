---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_package_attachment"
sidebar_current: "docs-alicloud-resource-cen-bandwidth-package-attachment"
description: |-
  Provides a Alicloud CEN bandwidth package attachment resource.
---

# alicloud\_cen_bandwidth_package_attachment

Provides a CEN bandwidth package attachment resource. The resource can be used to bind a bandwidth package to a specified CEN instance.

## Example Usage

Basic Usage

```
# Create a new bandwidth package attachment and use it to attach a bandwidth package to a new CEN
resource "alicloud_cen_instance" "cen" {
  name        = "tf-testAccCenBandwidthPackageAttachmentConfig"
  description = "tf-testAccCenBandwidthPackageAttachmentDescription"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
  bandwidth = 20
  geographic_region_ids = [
    "China",
  "Asia-Pacific"]
}

resource "alicloud_cen_bandwidth_package_attachment" "foo" {
  instance_id          = "${alicloud_cen_instance.cen.id}"
  bandwidth_package_id = "${alicloud_cen_bandwidth_package.bwp.id}"
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the CEN.
* `bandwidth_package_id` - (Required, ForceNew) The ID of the bandwidth package.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, the same as bandwidth_package_id.

## Import

CEN bandwidth package attachment resource can be imported using the id, e.g.

```
$terraform import alicloud_cen_bandwidth_package_attachment.example bwp-abc123456
```



