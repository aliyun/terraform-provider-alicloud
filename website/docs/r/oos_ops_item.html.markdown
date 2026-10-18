---
subcategory: "OOS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_ops_item"
sidebar_current: "docs-alicloud-resource-oos-ops-item"
description: |-
  Provides an OOS Ops Item resource.
---

# alicloud_oos_ops_item

Provides an OOS Ops Item resource. An Ops Item tracks an operational event across resources.

For information about OOS Ops Item and how to use it, see [What is Ops Item](https://www.alibabacloud.com/help/en/oos/developer-reference/api-oos-2019-06-01-createopsitem).

-> **NOTE:** Available since v1.295.0.

-> **NOTE:** OOS exposes no hard-delete API for Ops Items. `terraform destroy` transitions the Ops Item to the `Resolved` status rather than removing it from OOS.

## Example Usage

Basic Usage

```terraform
resource "alicloud_oos_ops_item" "example" {
  title       = "terraform-example-ops-item"
  severity    = "High"
  category    = "Security"
  source      = "terraform"
  priority    = 2
  description = "An example Ops Item managed by Terraform"
  status      = "InProgress"
  tags = {
    Env  = "terraform"
    Name = "ops-item-example"
  }
  resources = [
    "acs:oos:cn-hangzhou:1234567890123456:application/example/applicationgroup/default",
  ]
  solutions {
    priority    = 3
    type        = "SingleAZEcs"
    url         = "https://example.com"
    description = "Create ECS instances across zones"
  }
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The title of the Ops Item.
* `severity` - (Required) The severity of the Ops Item. Valid values: `Critical`, `High`, `Medium`, `Low`.
* `category` - (Optional) The category of the Ops Item. Valid values: `Availability`, `Cost`, `Performance`, `Recovery`, `Security`.
* `source` - (Required) The source of the Ops Item. The CreateOpsItem API returns 400 MissingSource when it is omitted.
* `priority` - (Required) The priority of the Ops Item. The CreateOpsItem API returns 400 InvalidParameter when it is omitted.
* `description` - (Optional) The description of the Ops Item.
* `status` - (Optional) The status of the Ops Item. OOS creates Ops Items as `Open`; the status is applied on update, not at creation. Valid values: `Open`, `InProgress`, `Resolved`. Destroying the resource sets `Resolved`.
* `dedup_string` - (Optional, ForceNew) The deduplication string of the Ops Item. Changing it forces replacement. OOS does not return this field on read, so the configured value is preserved in state.
* `resource_group_id` - (Optional) The ID of the resource group.
* `tags` - (Optional) The tags of the Ops Item.
* `resources` - (Optional) The list of associated resource ARNs.
* `solutions` - (Optional) The solutions of the Ops Item. See [solutions](#solutions-block) below.

<a name="solutions-block"></a>The `solutions` block supports:

* `priority` - (Optional) The priority of the solution.
* `type` - (Optional) The type of the solution.
* `url` - (Optional) The URL of the solution.
* `description` - (Optional) The description of the solution.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `ops_item_id`.
* `ops_item_id` - The ID of the Ops Item.
* `create_date` - The creation time of the Ops Item.
* `update_date` - The last update time of the Ops Item.
* `create_by` - The creator of the Ops Item.
* `last_modified_by` - The last modifier of the Ops Item.
* `attributes` - The custom attributes of the Ops Item.

## Import

OOS Ops Item can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_ops_item.example <ops_item_id>
```
