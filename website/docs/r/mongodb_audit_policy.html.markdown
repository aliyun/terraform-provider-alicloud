---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_audit_policy"
sidebar_current: "docs-alicloud-resource-mongodb-audit-policy"
description: |-
  Provides a Alicloud MongoDB Audit Policy resource.
---

# alicloud\_mongodb\_audit\_policy

Provides a MongoDB Audit Policy resource.

For information about MongoDB Audit Policy and how to use it, see [What is Audit Policy](https://www.alibabacloud.com/help/doc-detail/131941.html).

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_name = "subnet-for-local-test"
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = "example_value"
  vswitch_id          = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_mongodb_audit_policy" "example" {
  db_instance_id = alicloud_mongodb_instance.default.id
  audit_status   = "disabled"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the instance.
* `audit_status` - (Required) The status of the audit log. Valid values: `disabled`, `enable`.
* `storage_period` - (Optional) The retention period of audit logs. Valid values: `1` to `30`. Default value: `30`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Audit Policy. Its value is same as `db_instance_id`.

### Timeouts

-> **NOTE:** Available in 1.161.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Audit Policy.
* `update` - (Defaults to 5 mins) Used when update the Audit Policy.

## Import

MongoDB Audit Policy can be imported using the id, e.g.

```
$ terraform import alicloud_mongodb_audit_policy.example <db_instance_id>
```