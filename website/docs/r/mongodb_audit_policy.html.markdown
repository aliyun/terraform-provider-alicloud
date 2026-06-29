---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_audit_policy"
description: |-
  Provides a Alicloud Mongodb Audit Policy resource.
---

# alicloud_mongodb_audit_policy

Provides a Mongodb Audit Policy resource.



For information about Mongodb Audit Policy and how to use it, see [What is Audit Policy](https://www.alibabacloud.com/help/doc-detail/131941.html).

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
  audit_status   = "enable"
}
```

### Deleting `alicloud_mongodb_audit_policy` or removing it from your configuration

Terraform cannot destroy resource `alicloud_mongodb_audit_policy`. Terraform will remove this resource from the state file, however resources may remain.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_mongodb_audit_policy&spm=docs.r.mongodb_audit_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `audit_status` - (Required) Audit state. Valid values: `enable`, `disabled`. The audit policy cannot be created with `disabled` — the underlying API rejects it. Create the resource with `enable` and switch to `disabled` in a later apply.
* `db_instance_id` - (Required, ForceNew) Database Instance Id
* `storage_period` - (Optional, Int) Audit log retention duration, in days.
  - When `service_type` is `Standard`, the value range is 1 to 365 days. The default value is 30 days.
  - When `service_type` is `V2_Standard`, this is the cold storage duration and is required. Valid values: `30`, `180`, `365`, `1095`, `1825`.
* `service_type` - (Optional, Computed, Available since v1.284.0) The edition of the audit log. Valid values: `Standard`, `V2_Standard`. If omitted, the Provider sends `Standard`. In regions where only the V2 audit log is available, set this to `V2_Standard`. Changes to this field are ignored while `audit_status` is `disabled` — the server switches the edition internally when audit is off and restores the declared value on re-enable.
* `hot_storage_period` - (Optional, Int, Available since v1.284.0) The hot storage duration of the audit log, in days. The value range is 0 to 7. Only takes effect when `service_type` is `V2_Standard`.
* `filter` - (Optional, Available since v1.270.0) The type of logs collected by the audit log feature of the instance. Separate multiple types with commas (,). Valid values:
  - `admin`: O & M control operation.
  - `slow`: slow log.
  - `query`: the query operation.
  - `insert`: insert operation.
  - `update`: The update operation.
  - `delete`: deletes the operation.
  - `command`: Protocol command. For example, the aggregate aggregation method.
-> **NOTE:** `filter` only supports ApsaraDB for MongoDB replica set instances with `storage_type` of `local_ssd`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

-> **NOTE:** Available since v1.161.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Audit Policy.
* `update` - (Defaults to 15 mins) Used when update the Audit Policy.

## Import

Mongodb Audit Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_audit_policy.example <db_instance_id>
```
