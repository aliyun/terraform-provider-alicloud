---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_control_policy_order"
sidebar_current: "docs-alicloud-resource-cloud-firewall-control-policy_order"
description: |-
  Provides a Alicloud Cloud Firewall Control Policy Order resource.
---

# alicloud\_cloud\_firewall\_control\_policy\_order

Provides a Cloud Firewall Control Policy resource.

For information about Cloud Firewall Control Policy Order and how to use it, see [What is Control Policy Order](https://www.alibabacloud.com/help/doc-detail/138867.htm).

-> **NOTE:** Available in v1.130.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_control_policy" "example1" {
  application_name = "ANY"
  acl_action       = "accept"
  description      = "example"
  destination_type = "net"
  destination      = "100.1.1.0/24"
  direction        = "out"
  proto            = "ANY"
  source           = "1.2.3.0/24"
  source_type      = "net"
}

resource "alicloud_cloud_firewall_control_policy_order" "example2" {
  acl_uuid  = alicloud_cloud_firewall_control_policy.example1.acl_uuid
  direction = alicloud_cloud_firewall_control_policy.example1.direction
  order     = 1
}

```

## Argument Reference

The following arguments are supported:


* `acl_uuid` - (Required) The unique ID of the access control policy.
* `direction` - (Required) Direction. Valid values: `in`, `out`.
* `order` - (Optional) The priority of the access control policy. The priority value starts from 1. A small priority value indicates a high priority. **NOTE:** The value of -1 indicates the lowest priority.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Control Policy Order. The value formats as `<acl_uuid>:<direction>`.

## Import

Cloud Firewall Control Policy Order can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_firewall_control_policy_order.example <acl_uuid>:<direction>
```
