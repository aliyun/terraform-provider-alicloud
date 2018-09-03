---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen"
sidebar_current: "docs-alicloud-resource-cen"
description: |-
  Provides a Alicloud CEN resource.
---

# alicloud\_cen

Provides a CEN resource.

## Example Usage

Basic Usage

```
resource "alicloud_cen" "cen" {
	name = "tf_test_foo"
	description = "an example for cen"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the CEN. Defaults to null.
* `description` - (Optional) The CEN description. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CEN.
* `name` - The name of the CEN.
* `description` - The description of the CEN.

## Import

CEN can be imported using the id, e.g.

```
$ terraform import alicloud_cen.example cen-abc123456
```

