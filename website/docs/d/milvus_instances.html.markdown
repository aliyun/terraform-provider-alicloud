---
subcategory: "Milvus"
layout: "alicloud"
page_title: "Alicloud: alicloud_milvus_instances"
sidebar_current: "docs-alicloud-datasource-milvus-instances"
description: |-
  Provides a list of Milvus Instance owned by an Alibaba Cloud account.
---

# alicloud_milvus_instances

This data source provides Milvus Instance available to the user.[What is Instance](https://next.api.alibabacloud.com/document/milvus/2023-10-12/CreateInstance)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "zone_id" {
  default = "cn-hangzhou-j"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "terraform-example"
}

resource "alicloud_milvus_instance" "default" {
  zone_id = var.zone_id
  vswitch_ids {
    vsw_id  = alicloud_vswitch.default.id
    zone_id = alicloud_vswitch.default.zone_id
  }
  db_admin_password = "TerraformExample123!"
  components {
    type           = "data"
    cu_num         = 2
    replica        = 1
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  components {
    type           = "index"
    cu_num         = 4
    replica        = 2
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  components {
    type           = "query"
    cu_num         = 4
    replica        = 2
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  components {
    type           = "proxy"
    cu_num         = 2
    replica        = 2
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  components {
    type           = "mix_coordinator"
    cu_num         = 4
    replica        = 2
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  instance_name         = var.name
  db_version            = "2.4"
  vpc_id                = alicloud_vpc.default.id
  ha                    = false
  payment_type          = "Subscription"
  multi_zone_mode       = "Single"
  payment_duration_unit = "year"
  payment_duration      = 1
  auto_pay              = true
}

data "alicloud_milvus_instances" "default" {
  ids           = [alicloud_milvus_instance.default.id]
  instance_name = var.name
}

output "first_instance_id" {
  value = data.alicloud_milvus_instances.default.instances.0.instance_id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (ForceNew, Optional) Instance ID
* `instance_name` - (ForceNew, Optional) Instance name. The length is limited to 1-64 characters and can only contain Chinese, letters, numbers,-,_
* `resource_group_id` - (ForceNew, Optional) Resource Group ID
* `tags` - (Optional) User Defined Label
* `ids` - (Optional, Computed) A list of Instance IDs.
* `name_regex` - (Optional) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Instance IDs.
* `names` - A list of name of Instances.
* `instances` - A list of Instance Entries. Each element contains the following attributes:
  * `auto_backup` - Whether to enable automatic backup.
  * `components` - Instance component information.
    * `cu_num` - The number of CU.
    * `cu_type` - The calculation type.
    * `data_disk` - The QueryNode data disk configuration.
      * `enabled` - Whether to enable the QueryNode data disk.
      * `performance_level` - The ESSD performance level.
      * `size` - The data disk size in GiB.
      * `storage_class` - The data disk StorageClass.
    * `disk_size_type` - Default Normal.
    * `pay_type` - The payment type of the component.
    * `replica` - The number of component replicas.
    * `type` - The component type.
  * `configuration` - User-defined configuration.
  * `create_time` - Instance creation time.
  * `db_version` - Milvus kernel version.
  * `encrypted` - Whether to use kms encryption.
  * `expire_time` - The expiration time of the instance, which is returned by the package year and month cluster.
  * `ha` - Whether to enable multiple copies of data.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `kms_key_id` - Kms Key encryption id, need to be encrypted set to true.
  * `multi_zone_mode` - Availability Zone mode.
  * `order_id` - Alibaba Cloud Order Number.
  * `payment_type` - Payment Type.
  * `region_id` - regionId.
  * `resource_group_id` - Resource Group ID.
  * `running_time` - Instance running time.
  * `security_group_ids` - Configured Security Group id.
  * `status` - Instance status.
  * `tags` - User Defined Label.
  * `vswitch_ids` - Switch list, configure the switch and zone.
    * `vsw_id` - VSwitch id, which must correspond to the zone id.
    * `zone_id` - The availability zone must correspond to the vswId.
  * `vpc_id` - The VPC network ID.
  * `zone_id` - The zone id.
  * `id` - The ID of the resource supplied above.
