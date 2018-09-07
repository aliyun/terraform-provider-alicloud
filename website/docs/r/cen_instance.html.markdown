---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance"
sidebar_current: "docs-alicloud-resource-cen-instance"
description: |-
  Provides a Alicloud CEN instance resource.
---

# alicloud\_cen_instance

Provides a CEN instance resource. Cloud Enterprise Network (CEN) is a service that allows you to create a global network for rapidly building a distributed business system with a hybrid cloud computing solution. CEN enables you to build a secure, private, and enterprise-class interconnected network between VPCs in different regions and your local data centers. CEN provides enterprise-class scalability that automatically responds to your dynamic computing requirements.

For information about CEN and how to use it, see [What is Cloud Enterprise Network](https://www.alibabacloud.com/help/doc-detail/59870.htm).

## Example Usage

Basic Usage

```
resource "alicloud_cen_instance" "cen" {
	name = "tf_test_foo"
	description = "an example for cen"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the CEN instance. Defaults to null.
* `description` - (Optional) The description of the CEN instance. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CEN instance.
* `name` - The name of the CEN instance.
* `description` - The description of the CEN instance.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_instance.example cen-abc123456
```

