---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_control_policy_order"
sidebar_current: "docs-alicloud-resource-cloud-firewall-control-policy-order"
description: |-
  Provides a Alicloud Cloud Firewall Control Policy Order resource.
---

# alicloud_cloud_firewall_control_policy_order

Provides a Cloud Firewall Control Policy Order resource.

For information about Cloud Firewall Control Policy Order and how to use it, see [What is Control Policy Order](https://www.alibabacloud.com/help/doc-detail/138867.htm).

-> **NOTE:** Available since v1.130.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_control_policy_order&exampleId=5a361d38-71eb-6652-01ae-b979cb63ce2a68cd5f16&activeTab=example&spm=docs.r.cloud_firewall_control_policy_order.0.5a361d3871&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cloud_firewall_control_policy" "default" {
  direction        = "in"
  application_name = "ANY"
  description      = var.name
  acl_action       = "accept"
  source           = "127.0.0.1/32"
  source_type      = "net"
  destination      = "127.0.0.2/32"
  destination_type = "net"
  proto            = "ANY"
}

resource "alicloud_cloud_firewall_control_policy_order" "default" {
  acl_uuid  = alicloud_cloud_firewall_control_policy.default.acl_uuid
  direction = alicloud_cloud_firewall_control_policy.default.direction
  order     = 1
}
```

## Argument Reference

The following arguments are supported:

* `acl_uuid` - (Required, ForceNew) The unique ID of the access control policy.
* `direction` - (Required, ForceNew) The direction of the traffic to which the access control policy applies. Valid values: `in`, `out`.
* `order` - (Required, Int) The priority of the access control policy. The priority value starts from 1. A small priority value indicates a high priority. **NOTE:** The value of `-1` indicates the lowest priority.
-> **NOTE:** From version 1.227.1, `order` must be set.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Control Policy Order. It formats as `<acl_uuid>:<direction>`.

## Import

Cloud Firewall Control Policy Order can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_control_policy_order.example <acl_uuid>:<direction>
```
