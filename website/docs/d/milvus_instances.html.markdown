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

provider "alicloud" {
  region = ""
}


resource "alicloud_milvus_instance" "default" {
  ai_function = false
  zone_id     = "cn-hangzhou-j"
  vswitch_ids {
    zone_id = "cn-hangzhou-j"
    vsw_id  = "vsw-bp1pommb2vygb0kzvf8i6"
  }
  vswitch_ids {
    zone_id = "cn-hangzhou-k"
    vsw_id  = "vsw-bp1tomony773mb6nlabw9"
  }
  encrypted             = false
  auto_renew            = false
  payment_duration_unit = "month"
  auto_pay              = true
  load_replicas         = "2"
  payment_duration      = "1"
  db_admin_password     = "@1234Test"
  instance_name         = "tf-parity-sub-month-{{function.randomIntString(100000,999999)}}"
  components {
    cu_type = "general"
    type    = "streaming"
    cu_num  = "4"
    replica = "2"
  }
  components {
    cu_type = "general"
    type    = "data"
    cu_num  = "4"
    replica = "1"
  }
  components {
    cu_type        = "general"
    type           = "query"
    cu_num         = "16"
    disk_size_type = "Normal"
    replica        = "2"
  }
  components {
    cu_type = "general"
    type    = "proxy"
    cu_num  = "2"
    replica = "2"
  }
  components {
    cu_type = "general"
    type    = "mix_coordinator"
    cu_num  = "4"
    replica = "2"
  }
  db_version          = "2.6"
  vpc_id              = "vpc-bp168d0ay5yft9aira762"
  is_multi_az_storage = true
  payment_type        = "Subscription"
  ha                  = true
  multi_zone_mode     = "single"
  auto_backup         = true
  promotion_no        = "youhuiquan_promotion_option_id_for_blank"
}

data "alicloud_milvus_instances" "default" {
  ids               = ["${alicloud_milvus_instance.default.id}"]
  name_regex        = alicloud_milvus_instance.default.instance_name
  instance_name     = "tf-parity-sub-month-{{function.randomIntString(100000,999999)}}"
  resource_group_id = ""
}

output "alicloud_milvus_instance_example_id" {
  value = data.alicloud_milvus_instances.default.instances.0.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Optional) Instance ID
* `instance_name` - (Optional) Instance name. The length is limited to 1-64 characters and can only contain Chinese, letters, numbers,-,_
* `resource_group_id` - (Optional) Resource Group ID
* `tags` - (Optional) A map of tags to filter instances by.
* `ids` - (Optional) A list of Instance IDs.
* `name_regex` - (Optional) A regex string to filter results by instance name.
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
