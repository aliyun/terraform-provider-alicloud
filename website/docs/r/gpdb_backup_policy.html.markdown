---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_backup_policy"
description: |-
  Provides a Alicloud GPDB Backup Policy resource.
---

# alicloud_gpdb_backup_policy

Provides a GPDB Backup Policy resource. Describe the instance backup strategy.

For information about GPDB Backup Policy and how to use it, see [What is Backup Policy](https://www.alibabacloud.com/help/en/analyticdb-for-postgresql/latest/api-gpdb-2016-05-03-modifybackuppolicy).

-> **NOTE:** Available since v1.211.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_backup_policy&exampleId=fc78f005-be65-a4d8-4f9e-5c8db3140e1195ea21ee&activeTab=example&spm=docs.r.gpdb_backup_policy.0.fc78f005be&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_gpdb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.ids.0
  vswitch_name = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = local.vswitch_id
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_gpdb_backup_policy" "default" {
  db_instance_id          = alicloud_gpdb_instance.default.id
  recovery_point_period   = "1"
  enable_recovery_point   = "true"
  preferred_backup_period = "Wednesday"
  preferred_backup_time   = "15:00Z-16:00Z"
  backup_retention_period = "7"
}
```

### Deleting `alicloud_gpdb_backup_policy` or removing it from your configuration

Terraform cannot destroy resource `alicloud_gpdb_backup_policy`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `backup_retention_period` - (Optional) Data backup retention days.
* `db_instance_id` - (Required, ForceNew) The instance ID.
-> **NOTE:**  You can call the [DescribeDBInstances](~~ 86911 ~~) operation to view the details of all AnalyticDB PostgreSQL instances in the target region, including the instance ID.
* `enable_recovery_point` - (Optional) Whether to enable automatic recovery points. Value Description:
  - **true**: enabled.
  - **false**: Closed.
* `preferred_backup_period` - (Required) Data backup cycle. Separate multiple values with commas (,). Value Description:
  - **Monday**: Monday.
  - **Tuesday**: Tuesday.
  - **Wednesday**: Wednesday.
  - **Thursday**: Thursday.
  - **Friday**: Friday.
  - **Saturday**: Saturday.
  - **Sunday**: Sunday.
* `preferred_backup_time` - (Required) Data backup time. Format: HH:mmZ-HH:mmZ(UTC time).
* `recovery_point_period` - (Optional) Recovery point frequency. Value Description:
  - **1**: Hourly.
  - **2**: Every two hours.
  - **4**: Every four hours.
  - **8**: Every eight hours.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Backup Policy.
* `update` - (Defaults to 5 mins) Used when update the Backup Policy.

## Import

GPDB Backup Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_backup_policy.example <id>
```