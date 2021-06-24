---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_schema_group"
sidebar_current: "docs-alicloud-resource-event-bridge-schema-group"
description: |-
  Provides a Alicloud Event Bridge Schema Group resource.
---

# alicloud\_event\_bridge\_schema\_group

Provides a Event Bridge Schema Group resource.

For information about Event Bridge Schema Group and how to use it, see [What is Schema Group](https://www.alibabacloud.com/help/en/doc-detail/214326.htm).

-> **NOTE:** Available in v1.126.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_event_bridge_schema_group" "example" {
  group_id = "tf-testacc1234"
  format = "OPEN_API_3_0"
}

```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) The first ID of the resource.
* `description` - (Optional) - the description of group.
* `format` - (Optional, ForceNew) The format of schemas. Valid values: `AVRO`: Avro Schema Format, `JSON_SCHEMA_DRAFT_4`: Json Schema Format, `OPEN_API_3_0`: Open Api 3.0 Schema Format, which is a subset of Json Schema.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Schema Group. Value as `group_id`.

## Import

Event Bridge Schema Group can be imported using the id, e.g.

```
$ terraform import alicloud_event_bridge_schema_group.example <group_id>
```