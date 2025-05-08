---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_audit_policy"
sidebar_current: "docs-alicloud-resource-mongodb-audit-policy"
description: |-
  Provides a Alicloud MongoDB Audit Policy resource.
---

# alicloud_mongodb_audit_policy

Provides a MongoDB Audit Policy resource.

For information about MongoDB Audit Policy and how to use it, see [What is Audit Policy](https://www.alibabacloud.com/help/doc-detail/131941.html).

-> **NOTE:** Available since v1.148.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mongodb_audit_policy&exampleId=579e6708-068f-a8e0-b150-a675a47950d124d6f20b&activeTab=example&spm=docs.r.mongodb_audit_policy.0.579e670806&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_mongodb_zones" "default" {}
locals {
  index   = length(data.alicloud_mongodb_zones.default.zones) - 1
  zone_id = data.alicloud_mongodb_zones.default.zones[local.index].id
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = local.zone_id
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  vswitch_id          = alicloud_vswitch.default.id
  security_ip_list    = ["10.168.1.12", "100.69.7.112"]
  name                = var.name
  tags = {
    Created = "TF"
    For     = "example"
  }
}

resource "alicloud_mongodb_audit_policy" "default" {
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

## Timeouts

-> **NOTE:** Available since v1.161.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Audit Policy.
* `update` - (Defaults to 5 mins) Used when update the Audit Policy.

## Import

MongoDB Audit Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_audit_policy.example <db_instance_id>
```