---
subcategory: "Open Search"
layout: "alicloud"
page_title: "Alicloud: alicloud_open_search_app_group"
sidebar_current: "docs-alicloud-resource-open-search-app-group"
description: |-
  Provides a Alicloud Open Search App Group resource.
---

# alicloud\_open\_search\_app\_group

Provides a Open Search App Group resource.

For information about Open Search App Group and how to use it, see [What is App Group](https://www.aliyun.com/product/opensearch).

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "name"
}
resource "alicloud_open_search_app_group" "default" {
  app_group_name = var.name
  payment_type   = "PayAsYouGo"
  type           = "standard"
  quota {
    doc_size         = 1
    compute_resource = 20
    spec             = "opensearch.share.common"
  }
}

```

## Argument Reference

The resource does not support any argument.

* `app_group_name` - (Required,ForceNew) Application Group Name.
* `type` - (Required,ForceNew) Application type. Valid Values: `standard`, `enhanced`.
* `payment_type` - (Required,ForceNew) The billing method of the resource. Valid values: `Subscription` and `PayAsYouGo`.
* `quota` - (Required,ForceNew) Quota information.  The details see Block quota.
* `order` - (Optional,ForceNew) Order cycle information. The details see Block order.
* `order_type` - (Optional) Order change type. Valid values: `UPGRADE` and `DOWNGRADE`.
* `charge_way` - (Optional) Billing model. Valid values:`compute_resource` and `qps`.
* `description` - (Optional) The description of the resource.
* `current_version` - (Optional) The version of Application Group Name.

#### quota
The quota supports the following:

* `doc_size` - (Required) Storage Size. Unit: GB.
* `compute_resource` - (Required) Computing resources. Unit: LCU.
* `qps` - (Required) Search request. Unit: times/second. 
* `spec` - (Required) Specification. Valid values: 
  * `opensearch.share.junior`: Entry-level.
  * `opensearch.share.common`: Shared universal.
  * `opensearch.share.compute`: Shared computing.
  * `opensearch.share.storage`: Shared storage type.
  * `opensearch.private.common`: Exclusive universal type.
  * `opensearch.private.compute`: Exclusive computing type.
  * `opensearch.private.storage`: Exclusive storage type

#### order
The quota supports the following:
* `duration` - (Optional) Order cycle. The minimum value is not less than 0.
* `pricing_cycle` - (Optional) Order cycle unit. Valid values: `Year` and `Month`.
* `auto_renew` - (Optional) Whether to renew automatically. It only takes effect when the parameter payment_type takes the value `Subscription`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of App Group. It is the same as `app_group_name`.
* `status` - The status of the resource. Valid values: `producing`,`review_pending`,`config_pending`,`normal`,`frozen`.
* `instance_id` - The instance id.

## Import

Open Search App Group can be imported using the id, e.g.

```
$ terraform import alicloud_open_search_app_group.example <id>
```
