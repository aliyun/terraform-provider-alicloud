---
subcategory: "Milvus"
layout: "alicloud"
page_title: "Alicloud: alicloud_milvus_instance"
description: |-
  Provides a Alicloud Milvus Instance resource.
---

# alicloud_milvus_instance

Provides a Milvus Instance resource.



For information about Milvus Instance and how to use it, see [What is Instance](https://next.api.alibabacloud.com/document/milvus/2023-10-12/CreateInstance).

-> **NOTE:** Available since v1.264.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_milvus_instance&exampleId=74471fee-a02b-33d5-e171-b47ee729834fa8fa9bad&activeTab=example&spm=docs.r.milvus_instance.0.74471feea0&intl_lang=EN_US" target="_blank">
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

variable "region_id" {
  default = "cn-hangzhou"
}

variable "zone_id" {
  default = "cn-hangzhou-j"
}

resource "alicloud_vpc" "defaultILXuit" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultN80M7S" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "milvus-example"
}


resource "alicloud_milvus_instance" "default" {
  zone_id = var.zone_id
  vswitch_ids {
    vsw_id  = alicloud_vswitch.defaultN80M7S.id
    zone_id = alicloud_vswitch.defaultN80M7S.zone_id
  }
  db_admin_password = "Test123456@"
  components {
    type    = "standalone"
    cu_num  = "8"
    replica = "1"
    cu_type = "general"
  }
  instance_name         = "镇远测试包年包月"
  db_version            = "2.4"
  vpc_id                = alicloud_vpc.defaultILXuit.id
  ha                    = false
  payment_type          = "Subscription"
  multi_zone_mode       = "Single"
  payment_duration_unit = "year"
  payment_duration      = "1"
}
```

### Deleting `alicloud_milvus_instance` or removing it from your configuration

The `alicloud_milvus_instance` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `auto_backup` - (Optional, Computed) Whether to enable automatic backup
* `components` - (Optional, Set) Instance component information. Includes Starter Edition/Standard Edition.
  - Starter version: Array including standalone
  - Standard Edition: The configuration is different according to the 2.5 version and 2.6 version.
2.5: proxy ,mix_coordinator,data,query,index
2.6 need to configure: proxy,mix_coordinator,data,query,streaming See [`components`](#components) below.
* `configuration` - (Optional) User-defined configuration
* `db_admin_password` - (Optional) DB administrator password, which can be used to log in to attu.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `db_version` - (Required, ForceNew) Milvus kernel version. Supported versions: 2.4, 2.5, 2.6.
* `encrypted` - (Optional, ForceNew) Whether to use kms encryption. After enabling, you need to configure KmsKeyId. The default is false.
* `ha` - (Optional) Whether to enable multiple copies of data
* `instance_name` - (Required) Instance name. The length is limited to 1-64 characters and can only contain Chinese, letters, numbers,-,_
* `kms_key_id` - (Optional, ForceNew) Kms Key encryption id, need to be encrypted set to true.
* `multi_zone_mode` - (Optional, ForceNew) Availability Zone mode. The default Single.
  - Single: Single zone.
  - Two: Dual Availability Zones.
* `payment_duration` - (Optional, Int) Instance Payment Duration

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `payment_duration_unit` - (Optional) Paid unit , Enumeration value:
  - Month: Month
  - Year: Year

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `payment_type` - (Required, ForceNew) Payment Type ,Enumeration value:
  - PayAsYouGo: Pay by volume
  - Subscription: Package year package month
* `resource_group_id` - (Optional, Computed) Resource Group ID
* `tags` - (Optional, Map) User Defined Label
* `vswitch_ids` - (Optional, ForceNew, List) Switch list, configure the switch and zone. See [`vswitch_ids`](#vswitch_ids) below.
* `vpc_id` - (Required, ForceNew) The VPC network ID. vpc-xxx.
* `zone_id` - (Optional, ForceNew) The zone id. When multi-zone is enabled, it represents the primary zone.

### `components`

The components supports the following:
* `cu_num` - (Required, Int) The number of CU. For example: 4
* `cu_type` - (Optional, ForceNew, Computed) The calculation type. The default value is general, and the ram type needs to be opened with a work order.
  - general: Generic
  - ram: Capacity
* `disk_size_type` - (Optional, ForceNew, Computed) Default Normal. The Query Node is configured with the capacity type, performance type, and capacity type Large, and the rest are configured with Normal.
* `replica` - (Required, Int) The number of component replicas. The number of highly available replicas must be greater than or equal to 2.
* `type` - (Required) The component type. Different types need to be configured according to different versions.
  - Starter version: Array including standalone
  - Standard Edition: The configuration is different according to the 2.5 version and 2.6 version.
2.5: proxy ,mix_coordinator,data,query,index
2.6 need to configure: proxy,mix_coordinator,data,query,streaming

### `vswitch_ids`

The vswitch_ids supports the following:
* `vsw_id` - (Optional, ForceNew) VSwitch id, which must correspond to the zone id.
* `zone_id` - (Optional, ForceNew) The availability zone must correspond to the vswId.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Instance creation time.
* `region_id` - regionId. For example: cn-hangzhou
* `status` - Instance status. Value range:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 14 mins) Used when create the Instance.
* `delete` - (Defaults to 14 mins) Used when delete the Instance.
* `update` - (Defaults to 21 mins) Used when update the Instance.

## Import

Milvus Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_milvus_instance.example <id>
```