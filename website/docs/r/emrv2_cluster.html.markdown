---
subcategory: "E-MapReduce (EMR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_emrv2_cluster"
sidebar_current: "docs-alicloud-resource-emr-cluster-new"
description: |-
  Provides a EMR Cluster resource.
---

# alicloud_emrv2_cluster

Provides a EMR cluster resource. This resource is based on EMR's new version OpenAPI.

For information about EMR New and how to use it, see [Add a domain](https://www.alibabacloud.com/help/doc-detail/28068.htm).

-> **NOTE:** Available since v1.199.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_zones" "default" {
  available_instance_type = "ecs.g7.xlarge"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}
resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = var.name
}
resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ram_role" "default" {
  name     = "tfexampleroleemrv2"
  document = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF

  description = "this is a role example."
  force       = true
}

resource "alicloud_emrv2_cluster" "default" {
  payment_type    = "PayAsYouGo"
  cluster_type    = "DATALAKE"
  release_version = "EMR-5.10.0"
  cluster_name    = var.name
  deploy_mode     = "NORMAL"
  security_mode   = "NORMAL"

  applications = ["HADOOP-COMMON", "HDFS", "YARN", "HIVE", "SPARK3", "TEZ"]

  application_configs {
    application_name  = "HIVE"
    config_file_name  = "hivemetastore-site.xml"
    config_item_key   = "hive.metastore.type"
    config_item_value = "DLF"
    config_scope      = "CLUSTER"
  }

  application_configs {
    application_name  = "SPARK3"
    config_file_name  = "hive-site.xml"
    config_item_key   = "hive.metastore.type"
    config_item_value = "DLF"
    config_scope      = "CLUSTER"
  }

  node_attributes {
    ram_role          = alicloud_ram_role.default.name
    security_group_id = alicloud_security_group.default.id
    vpc_id            = alicloud_vpc.default.id
    zone_id           = data.alicloud_zones.default.zones.0.id
    key_pair_name     = alicloud_ecs_key_pair.default.id
  }

  tags = {
    created = "tf"
  }

  node_groups {
    node_group_type = "MASTER"
    node_group_name = "emr-master"
    payment_type    = "PayAsYouGo"
    vswitch_ids     = [alicloud_vswitch.default.id]
    with_public_ip  = false
    instance_types  = ["ecs.g7.xlarge"]
    node_count      = 1

    system_disk {
      category = "cloud_essd"
      size     = 80
      count    = 1
    }

    data_disks {
      category = "cloud_essd"
      size     = 80
      count    = 3
    }
  }

  node_groups {
    node_group_type = "CORE"
    node_group_name = "emr-core"
    payment_type    = "PayAsYouGo"
    vswitch_ids     = [alicloud_vswitch.default.id]
    with_public_ip  = false
    instance_types  = ["ecs.g7.xlarge"]
    node_count      = 3

    system_disk {
      category = "cloud_essd"
      size     = 80
      count    = 1
    }

    data_disks {
      category = "cloud_essd"
      size     = 80
      count    = 3
    }
  }

  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```
## Argument Reference

The following arguments are supported:

* `resource_group_id` - (Optional) The Id of resource group which the emr-cluster belongs.
* `payment_type` - (Optional, ForceNew) Payment Type for this cluster. Supported value: PayAsYouGo or Subscription.
* `subscription_config` - (Optional) The detail configuration of subscription payment type. See [`subscription_config`](#subscription_config) below.
* `cluster_type` - (Required, ForceNew) EMR Cluster Type, e.g. DATALAKE, OLAP, DATAFLOW, DATASERVING, CUSTOM etc. You can find all valid EMR cluster type in emr web console.
* `release_version` - (Required, ForceNew) EMR Version, e.g. EMR-5.10.0. You can find the all valid EMR Version in emr web console.
* `cluster_name` - (Required) The name of emr cluster. The name length must be less than 64. Supported characters: chinese character, english character, number, "-", "_".
* `deploy_mode` - (Optional, ForceNew) The deploy mode of EMR cluster. Supported value: NORMAL or HA.
* `security_mode` - (Optional) The security mode of EMR cluster. Supported value: NORMAL or KERBEROS.
* `applications` - (Required, ForceNew) The applications of EMR cluster to be installed, e.g. HADOOP-COMMON, HDFS, YARN, HIVE, SPARK2, SPARK3, ZOOKEEPER etc. You can find all valid applications in emr web console.
* `application_configs` - (Optional) The application configurations of EMR cluster. See [`application_configs`](#application_configs) below.
* `node_attributes` - (Required, ForceNew) The node attributes of ecs instances which the emr-cluster belongs. See [`node_attributes`](#node_attributes) below.
* `node_groups` - (Required) Groups of node, You can specify MASTER as a group, CORE as a group (just like the above example). See [`node_groups`](#node_groups) below.
* `bootstrap_scripts` (Optional) The bootstrap scripts to be effected when creating emr-cluster or resize emr-cluster. See [`bootstrap_scripts`](#bootstrap_scripts) below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### `subscription_config`

The `subscription_config` block supports the following:

* `payment_duration_unit` - (Required) If paymentType is Subscription, this should be specified. Supported value: Month or Year.
* `payment_duration` - (Required) If paymentType is Subscription, this should be specified. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36、48.
* `auto_renew` - (Optional) Auto renew for prepaid, ’true’ or ‘false’ . Default value: false.
* `auto_renew_duration_unit` - (Optional) If paymentType is Subscription, this should be specified. Supported value: Month or Year.
* `auto_renew_duration` - (Optional) If paymentType is Subscription, this should be specified. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36、48. 

### `application_configs`

The `application_configs` block supports the following:

* `application_name` - (Required) The application name of EMR cluster which has installed.
* `config_file_name` - (Required) The configuration file name of application installed.
* `config_item_key` - (Required) The configuration item key of application installed.
* `config_item_value` - (Required) The configuration item value of application installed.
* `config_scope` - (Optional) The configuration scope of emr cluster. Supported value: CLUSTER or NODEGROUP.
* `config_description` - (Optional) The configuration description of application installed.
* `node_group_name` - (Optional) The configuration effected which node group name of emr cluster.
* `node_group_id` - (Optional) The configuration effected which node group id of emr cluster.

### `node_attributes`

The `node_attributes` block supports the following:

* `vpc_id` - (Required) Used to retrieve instances belong to specified VPC.
* `zone_id` - (Required) Zone ID, e.g. cn-hangzhou-i
* `security_group_id` - (Required) Security Group ID for Cluster.
* `ram_role` - (Required) Alicloud EMR uses roles to perform actions on your behalf when provisioning cluster resources, running applications, dynamically scaling resources. EMR uses the following roles when interacting with other Alicloud services. Default value is AliyunEmrEcsDefaultRole.
* `key_pair_name` - (Required) The name of the key pair.
* `data_disk_encrypted` - (Optional, ForceNew, Available since v1.204.0) Whether to enable data disk encryption.
* `data_disk_kms_key_id` - (Optional, ForceNew, Available since v1.204.0) The kms key id used to encrypt the data disk. It takes effect when data_disk_encrypted is true.

### `node_groups`

The node_groups mapping supports the following: 

* `node_group_type` - (Required) The node group type of emr cluster, supported value: MASTER, CORE or TASK.
* `node_group_name` - (Required) The node group name of emr cluster.
* `payment_type` - (Optional) Payment Type for this cluster. Supported value: PayAsYouGo or Subscription.
* `subscription_config` - (Optional) The detail configuration of subscription payment type. See [`subscription_config`](#node_groups-subscription_config) below.
* `spot_bid_prices` - (Optional) The spot bid prices of a PayAsYouGo instance. See [`spot_bid_prices`](#node_groups-spot_bid_prices) below.
* `vswitch_ids` - (Optional) Global vSwitch ids, you can also specify it in node group.
* `with_public_ip` - (Optional) Whether the node has a public IP address enabled.
* `additional_security_group_ids` - (Optional) Additional security Group IDS for Cluster, you can also specify this key for each node group.
* `instance_types` - (Required) Host Ecs instance types.
* `node_count` - (Required) Host Ecs number in this node group.
* `system_disk` - (Required) Host Ecs system disk information in this node group. See [`system_disk`](#node_groups-system_disk) below.
* `data_disks` - (Required) Host Ecs data disks information in this node group. See [`data_disks`](#node_groups-data_disks) below.
* `graceful_shutdown` - (Optional) Enable emr cluster of task node graceful decommission, ’true’ or ‘false’ .
* `spot_instance_remedy` - (Optional) Whether to replace spot instances with newly created spot/onDemand instance when receive a spot recycling message.
* `cost_optimized_config` - (Optional) The detail cost optimized configuration of emr cluster. See [`cost_optimized_config`](#node_groups-cost_optimized_config) below.

### `node_groups-subscription_config`

The subscription_config mapping supports the following: 

* `payment_duration_unit` - (Required) If paymentType is Subscription, this should be specified. Supported value: Month or Year.
* `payment_duration` - (Required) If paymentType is Subscription, this should be specified. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36、48.
* `auto_renew` - (Optional) Auto renew for prepaid, ’true’ or ‘false’ . Default value: false.
* `auto_renew_duration_unit` - (Optional) If paymentType is Subscription, this should be specified. Supported value: Month or Year.
* `auto_renew_duration` - (Optional) If paymentType is Subscription, this should be specified. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36、48. 

### `node_groups-spot_bid_prices`

The spot_bid_prices mapping supports the following: 

* `instance_type` - (Required) Host Ecs instance type.
* `bid_price` - (Required) The spot bid price of a PayAsYouGo instance.

### `node_groups-system_disk`

The system_disk mapping supports the following: 

* `category` - (Required) The type of the data disk. Valid values: `cloud_efficiency` and `cloud_essd`.
* `size` - (Required)The size of a data disk, at least 40. Unit: GiB.
* `performance_level` - (Optional) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity.
* `count` - (Optional) The count of a data disk.

### `node_groups-data_disks`

The data_disks mapping supports the following: 

* `category` - (Required) The type of the data disk. Valid values: `cloud_efficiency` and `cloud_essd`.
* `size` - (Required)The size of a data disk, at least 40. Unit: GiB.
* `performance_level` - (Optional) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity.
* `count` - (Optional) The count of a data disk.

### `node_groups-cost_optimized_config`

The cost_optimized_config mapping supports the following: 

* `on_demand_base_capacity` - (Required) The cost optimized configuration which on demand based capacity.
* `on_demand_percentage_above_base_capacity` - (Required) The cost optimized configuration which on demand percentage above based capacity.
* `spot_instance_pools` - (Required) The cost optimized configuration with spot instance pools.

### `bootstrap_scripts`

The bootstrap_scripts mapping supports the following: 

* `script_name` - (Required) The bootstrap script name.
* `script_path` - (Required) The bootstrap script path, e.g. "oss://bucket/path".
* `script_args` - (Required) The bootstrap script args, e.g. "--a=b".
* `priority` - (Optional) The bootstrap scripts priority.
* `execution_moment` - (Required) The bootstrap scripts execution moment, ’BEFORE_INSTALL’ or ‘AFTER_STARTED’ .
* `execution_fail_strategy` - (Required) The bootstrap scripts execution fail strategy, ’FAILED_BLOCKED’ or ‘FAILED_CONTINUE’ .
* `node_selector` - (Required) The bootstrap scripts execution target. See [`node_selector`](#bootstrap_scripts-node_selector) below.

### `bootstrap_scripts-node_selector`

The node_selector mapping supports the following: 

* `node_select_type` - (Required) The bootstrap scripts execution target node select type. Supported value: NODE, NODEGROUP or CLUSTER.
* `node_names` - (Optional) The bootstrap scripts execution target node names.
* `node_group_id` - (Optional) The bootstrap scripts execution target node group id.
* `node_group_types` - (Optional) The bootstrap scripts execution target node group types.
* `node_group_name` - (Optional) The bootstrap scripts execution target node group name.

## Attributes Reference

The following attributes are exported:

* `id` - The emr cluster ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the cluster (until it reaches the initial `RUNNING` status). 
* `delete` - (Defaults to 5 mins) Used when terminating the cluster.

## Import

Aliclioud E-MapReduce cluster can be imported using the id e.g.

```shell
$ terraform import alicloud_emrv2_cluster.default <id>
```