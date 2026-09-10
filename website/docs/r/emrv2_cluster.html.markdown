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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_emrv2_cluster&exampleId=5de14baf-18bd-38b6-f72b-9c99227d119a26075001&activeTab=example&spm=docs.r.emrv2_cluster.0.5de14baf18&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_kms_keys" "default" {
  status = "Enabled"
}

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

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ram_role" "default" {
  name        = var.name
  document    = <<EOF
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
  node_groups {
    vswitch_ids = [
      "${alicloud_vswitch.default.id}"
    ]
    instance_types = [
      "ecs.g7.xlarge"
    ]
    node_count           = "1"
    spot_instance_remedy = "false"
    data_disks {
      count             = "3"
      category          = "cloud_essd"
      size              = "80"
      performance_level = "PL0"
    }

    node_group_name   = "emr-master"
    payment_type      = "PayAsYouGo"
    with_public_ip    = "false"
    graceful_shutdown = "false"
    system_disk {
      category          = "cloud_essd"
      size              = "80"
      performance_level = "PL0"
      count             = "1"
    }

    node_group_type = "MASTER"
  }
  node_groups {
    spot_instance_remedy = "false"
    node_group_type      = "CORE"
    vswitch_ids = [
      "${alicloud_vswitch.default.id}"
    ]
    node_count        = "2"
    graceful_shutdown = "false"
    system_disk {
      performance_level = "PL0"
      count             = "1"
      category          = "cloud_essd"
      size              = "80"
    }

    data_disks {
      count             = "3"
      performance_level = "PL0"
      category          = "cloud_essd"
      size              = "80"
    }

    node_group_name = "emr-core"
    payment_type    = "PayAsYouGo"
    instance_types = [
      "ecs.g7.xlarge"
    ]
    with_public_ip = "false"
  }

  deploy_mode = "NORMAL"
  tags = {
    Created = "TF"
    For     = "example"
  }
  release_version = "EMR-5.10.0"
  applications = [
    "HADOOP-COMMON",
    "HDFS",
    "YARN"
  ]
  node_attributes {
    zone_id              = data.alicloud_zones.default.zones.0.id
    key_pair_name        = alicloud_ecs_key_pair.default.id
    data_disk_encrypted  = "true"
    data_disk_kms_key_id = data.alicloud_kms_keys.default.ids.0
    vpc_id               = alicloud_vpc.default.id
    ram_role             = alicloud_ram_role.default.name
    security_group_id    = alicloud_security_group.default.id
  }

  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  cluster_name      = var.name
  payment_type      = "PayAsYouGo"
  cluster_type      = "DATAFLOW"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_emrv2_cluster&spm=docs.r.emrv2_cluster.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `resource_group_id` - (Optional) The Id of resource group which the emr-cluster belongs.
* `payment_type` - (Optional) Payment Type for this cluster. Supported value: PayAsYouGo or Subscription. **NOTE:** From version 1.227.0, `payment_type` can be modified.
* `subscription_config` - (Optional) The detail configuration of subscription payment type. See [`subscription_config`](#subscription_config) below.
* `cluster_type` - (Required, ForceNew) EMR Cluster Type, e.g. DATALAKE, OLAP, DATAFLOW, DATASERVING, CUSTOM etc. You can find all valid EMR cluster type in emr web console.
* `release_version` - (Required, ForceNew) EMR Version, e.g. EMR-5.10.0. You can find the all valid EMR Version in emr web console.
* `cluster_name` - (Required) The name of emr cluster. The name length must be less than 64. Supported characters: chinese character, english character, number, "-", "_".
* `deploy_mode` - (Optional, ForceNew) The deploy mode of EMR cluster. Supported value: NORMAL or HA.
* `log_collect_strategy` - (Optional, Available since v1.219.0) The log collect strategy of EMR cluster. 
* `deletion_protection` - (Optional, Available since v1.236.0) The deletion protection of EMR cluster.
* `security_mode` - (Optional) The security mode of EMR cluster. Supported value: NORMAL or KERBEROS.
* `applications` - (Required, ForceNew) The applications of EMR cluster to be installed, e.g. HADOOP-COMMON, HDFS, YARN, HIVE, SPARK2, SPARK3, ZOOKEEPER etc. You can find all valid applications in emr web console.
* `application_configs` - (Optional) The application configurations of EMR cluster. See [`application_configs`](#application_configs) below.
* `node_attributes` - (Required, ForceNew) The node attributes of ecs instances which the emr-cluster belongs. See [`node_attributes`](#node_attributes) below.
* `node_groups` - (Required) Groups of node, You can specify MASTER as a group, CORE as a group (just like the above example). See [`node_groups`](#node_groups) below. **NOTE:** Since version 1.227.0, the type of `node_groups` changed from Set to List.
* `bootstrap_scripts` (Optional) The bootstrap scripts to be effected when creating emr-cluster or resize emr-cluster, if priority is not specified, the scripts will execute in the declared order. See [`bootstrap_scripts`](#bootstrap_scripts) below.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### `subscription_config`

The `subscription_config` block supports the following:

* `payment_duration_unit` - (Required) If paymentType is Subscription, this should be specified. Supported value: Month or Year.
* `payment_duration` - (Required) If paymentType is Subscription, this should be specified. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36、48.
* `auto_renew` - (Optional) Auto renew for prepaid, ’true’ or ‘false’ . Default value: false.
* `auto_pay_order` - (Optional, Available since v1.219.0) Auto pay order for payment type of subscription, ’true’ or ‘false’ .  Default value is ’true’.
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

* `vpc_id` - (Required, ForceNew) Used to retrieve instances belong to specified VPC.
* `zone_id` - (Required, ForceNew) Zone ID, e.g. cn-hangzhou-i
* `security_group_id` - (Required, ForceNew) Security Group ID for Cluster.
* `ram_role` - (Required, ForceNew) Alicloud EMR uses roles to perform actions on your behalf when provisioning cluster resources, running applications, dynamically scaling resources. EMR uses the following roles when interacting with other Alicloud services. Default value is AliyunEmrEcsDefaultRole.
* `key_pair_name` - (Required, ForceNew) The name of the key pair.
* `data_disk_encrypted` - (Optional, ForceNew, Available since v1.204.0) Whether to enable data disk encryption.
* `data_disk_kms_key_id` - (Optional, ForceNew, Available since v1.204.0) The kms key id used to encrypt the data disk. It takes effect when data_disk_encrypted is true.
* `system_disk_encrypted` - (Optional, ForceNew, Available since v1.242.0) Whether to enable system disk encryption.
* `system_disk_kms_key_id` - (Optional, ForceNew, Available since v1.242.0) The kms key id used to encrypt the system disk. It takes effect when system_disk_encrypted is true.

### `node_groups`

The node_groups mapping supports the following: 

* `node_group_type` - (Required) The node group type of emr cluster, supported value: MASTER, CORE or TASK. Node group type of GATEWAY is available since v1.219.0. Node group type of MASTER-EXTEND is available since v1.243.0.
* `node_group_name` - (Required) The node group name of emr cluster.
* `payment_type` - (Optional) Payment Type for this cluster. Supported value: PayAsYouGo or Subscription.
* `subscription_config` - (Optional) The detail configuration of subscription payment type. See [`subscription_config`](#node_groups-subscription_config) below.
* `private_pool_options` - (Optional, Available since v1.253.0) The node group specific private pool resources. See [`private_pool_options`](#node_groups-private_pool_options) below.
* `spot_bid_prices` - (Optional) The spot bid prices of a PayAsYouGo instance. See [`spot_bid_prices`](#node_groups-spot_bid_prices) below.
* `vswitch_ids` - (Optional) Global vSwitch ids, you can also specify it in node group. **NOTE:** From version 1.236.0, `vswitch_ids` can be modified.
* `with_public_ip` - (Optional) Whether the node has a public IP address enabled. **NOTE:** From version 1.236.0, `with_public_ip` can be modified.
* `additional_security_group_ids` - (Optional) Additional security Group IDS for Cluster, you can also specify this key for each node group. **NOTE:** From version 1.236.0, `additional_security_group_ids` can be modified.
* `instance_types` - (Required) Host Ecs instance types. **NOTE:** From version 1.236.0, `instance_types` can be modified.
* `node_count` - (Required) Host Ecs number in this node group.
* `system_disk` - (Required) Host Ecs system disk information in this node group. See [`system_disk`](#node_groups-system_disk) below.
* `data_disks` - (Required) Host Ecs data disks information in this node group. See [`data_disks`](#node_groups-data_disks) below.
* `graceful_shutdown` - (Optional) Enable emr cluster of task node graceful decommission, ’true’ or ‘false’ .
* `spot_instance_remedy` - (Optional) Whether to replace spot instances with newly created spot/onDemand instance when receive a spot recycling message.
* `spot_strategy` - (Optional, Available since v1.236.0) The spot strategy configuration of emr cluster. Valid values: `NoSpot`, `SpotWithPriceLimit`, `SpotAsPriceGo`.
* `cost_optimized_config` - (Optional) The detail cost optimized configuration of emr cluster. See [`cost_optimized_config`](#node_groups-cost_optimized_config) below. **NOTE:** From version 1.236.0, `cost_optimized_config` can be modified.
* `deployment_set_strategy` - (Optional, Available since v1.219.0) Deployment set strategy for this cluster node group. Supported value: NONE, CLUSTER or NODE_GROUP. **NOTE:** From version 1.236.0, `deployment_set_strategy` can be modified.
* `auto_scaling_policy` - (Optional, Available since v1.227.0) The node group auto scaling policy for emr cluster. See [`auto_scaling_policy`](#node_groups-auto_scaling_policy) below.
* `ack_config` - (Optional, Available since v1.236.0) The node group of ack configuration for emr cluster to deploying on kubernetes. See [`ack_config`](#node_groups-ack_config) below.
* `node_resize_strategy` - (Optional, Available since v1.219.0) Node resize strategy for this cluster node group. Supported value: PRIORITY, COST_OPTIMIZED.

### `node_groups-private_pool_options`

The private_pool_options mapping supports the following:

* `private_pool_ids` - (Optional) The node group specific private pool resource ids.
* `match_criteria` - (Optional) The node group specific private pool resource match criteria. Valid values: `Open`, `Target`, `None`.

### `node_groups-subscription_config`

The subscription_config mapping supports the following: 

* `payment_duration_unit` - (Required) If paymentType is Subscription, this should be specified. Supported value: Month or Year.
* `payment_duration` - (Required) If paymentType is Subscription, this should be specified. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36、48.
* `auto_renew` - (Optional) Auto renew for prepaid, ’true’ or ‘false’ . Default value: false.
* `auto_pay_order` - (Optional, Available since v1.219.0) Auto pay order for payment type of subscription, ’true’ or ‘false’ . Default value is ’true’.
* `auto_renew_duration_unit` - (Optional) If paymentType is Subscription, this should be specified. Supported value: Month or Year.
* `auto_renew_duration` - (Optional) If paymentType is Subscription, this should be specified. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36、48. 

### `node_groups-spot_bid_prices`

The spot_bid_prices mapping supports the following: 

* `instance_type` - (Required) Host Ecs instance type.
* `bid_price` - (Required) The spot bid price of a PayAsYouGo instance.

### `node_groups-system_disk`

The system_disk mapping supports the following: 

* `category` - (Required) The type of the data disk. Valid values: `cloud_efficiency`, `cloud_essd`, `cloud_ssd`. **NOTE:** Since version v1.230.0, the category `cloud_ssd` is available.
* `size` - (Required)The size of a data disk, at least 40. Unit: GiB.
* `performance_level` - (Optional) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity.
* `count` - (Optional) The count of a data disk.

### `node_groups-data_disks`

The data_disks mapping supports the following: 

* `category` - (Required) The type of the data disk. Valid values: `cloud_efficiency`, `cloud_essd`, `cloud`, `local_hdd_pro`, `local_disk`, `local_ssd_pro`. **NOTE:** Since version v1.230.0, the categories `cloud`, `local_hdd_pro`, `local_disk`, `local_ssd_pro` are available.
* `size` - (Required)The size of a data disk, at least 40. Unit: GiB.
* `performance_level` - (Optional) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity.
* `count` - (Optional) The count of a data disk.

### `node_groups-cost_optimized_config`

The cost_optimized_config mapping supports the following: 

* `on_demand_base_capacity` - (Required) The cost optimized configuration which on demand based capacity.
* `on_demand_percentage_above_base_capacity` - (Required) The cost optimized configuration which on demand percentage above based capacity.
* `spot_instance_pools` - (Required) The cost optimized configuration with spot instance pools.

### `node_groups-auto_scaling_policy`

The auto_scaling_policy mapping supports the following:

* `scaling_rules` - (Optional) The scaling rules of auto scaling policy. See [`scaling_rules`](#node_groups-auto_scaling_policy-scaling_rules) below.
* `constraints` - (Optional) The constraints of auto scaling policy. See [`constraints`](#node_groups-auto_scaling_policy-constraints) below.

### `node_groups-ack_config`

The ack_config mapping supports the following:

* `ack_instance_id` - (Required) The ack cluster instance id.
* `node_selectors` - (Optional) The ack cluster node selectors for job pods scheduling. See [`node_selectors`](#node_groups-ack_config-node_selectors) below.
* `tolerations` - (Optional) The ack cluster tolerations. See [`tolerations`](#node_groups-ack_config-tolerations) below.
* `namespace` - (Required) The ack cluster namespace.
* `request_cpu` - (Required) The job pod resource of request cpu.
* `request_memory` - (Required) The job pod resource of request memory.
* `limit_cpu` - (Required) The job pod resource of limit cpu.
* `limit_memory` - (Required) The job pod resource of limit memory.
* `custom_labels` - (Optional) The ack cluster custom labels. See [`custom_labels`](#node_groups-ack_config-custom_labels) below.
* `custom_annotations` - (Optional) The ack cluster custom annotations. See [`custom_annotations`](#node_groups-ack_config-custom_annotations) below.
* `pvcs` - (Optional) The ack cluster persistent volume claim. See [`pvcs`](#node_groups-ack_config-pvcs) below.
* `volumes` - (Optional) The ack cluster volumes. See [`volumes`](#node_groups-ack_config-volumes) below.
* `volume_mounts` - (Optional) The ack cluster volume mounts. See [`volume_mounts`](#node_groups-ack_config-volume_mounts) below.
* `pre_start_command` - (Optional) The job pod pre start command.
* `pod_affinity` - (Optional) The job pod affinity.
* `pod_anti_affinity` - (Optional) The job pod anti-affinity.
* `node_affinity` - (Optional) The ack cluster node affinity.

### `node_groups-ack_config-node_selectors`

The node_selectors mapping supports the following:

* `key` - (Required) The key of ack cluster node selector.
* `value` - (Optional) The value of ack cluster node selector.

### `node_groups-ack_config-tolerations`

The tolerations mapping supports the following:

* `key` - (Optional) The key of ack cluster tolerations.
* `value` - (Optional) The value of ack cluster tolerations.
* `operator` - (Optional) The operator of ack cluster tolerations.
* `effect` - (Optional) The effect of ack cluster tolerations.

### `node_groups-ack_config-custom_labels`

The custom_labels mapping supports the following:

* `key` - (Required) The key of ack cluster custom labels.
* `value` - (Optional) The value of ack cluster custom labels.

### `node_groups-ack_config-custom_annotations`

The custom_annotations mapping supports the following:

* `key` - (Required) The key of ack cluster custom labels.
* `value` - (Optional) The value of ack cluster custom labels.

### `node_groups-ack_config-pvcs`

The pvcs mapping supports the following:

* `name` - (Required) The name of ack cluster job pod persistent volume claim.
* `path` - (Required) The path of ack cluster job pod persistent volume claim.
* `data_disk_storage_class` - (Required) The ack cluster job pod data disk storage class of persistent volume claim.
* `data_disk_size` - (Required) The ack cluster job pod data disk size of persistent volume claim.

### `node_groups-ack_config-volumes`

The volumes mapping supports the following:

* `name` - (Required) The name of ack cluster job pod volumes.
* `path` - (Required) The path of ack cluster job pod volumes.
* `type` - (Required) The ack cluster job pod volumes type.

### `node_groups-ack_config-volume_mounts`

The volume_mounts mapping supports the following:

* `name` - (Required) The name of ack cluster job pod volume mounts.
* `path` - (Required) The path of ack cluster job pod volume mounts.

### `node_groups-auto_scaling_policy-scaling_rules`

The scaling_rules mapping supports the following:

* `rule_name` - (Required) The rule name of auto scaling policy.
* `trigger_type` - (Required) The trigger type of auto scaling policy. Valid values: `TIME_TRIGGER` and `METRICS_TRIGGER`.
* `activity_type` - (Required) The activity type of auto scaling policy. Valid values: `SCALE_OUT` and `SCALE_IN`.
* `adjustment_type` - (Optional) The adjustment type of auto scaling policy. Valid values: `CHANGE_IN_CAPACITY` and `EXACT_CAPACITY`.
* `adjustment_value` - (Required) The adjustment value of auto scaling policy. The value should between 1 and 5000.
* `min_adjustment_value` - (Optional) The minimum adjustment value of auto scaling policy.
* `time_trigger` - (Optional) The trigger time of scaling rules for emr node group auto scaling policy. See [`time_trigger`](#node_groups-auto_scaling_policy-scaling_rules-time_trigger) below.
* `metrics_trigger` - (Optional) The trigger metrics of scaling rules for emr node group auto scaling policy. See [`metrics_trigger`](#node_groups-auto_scaling_policy-scaling_rules-metrics_trigger) below.

### `node_groups-auto_scaling_policy-scaling_rules-time_trigger`

The time_trigger mapping supports the following:

* `launch_time` - (Required) The launch time for this scaling rule specific time trigger.
* `start_time` - (Optional) The start time for this scaling rule specific time trigger.
* `end_time` - (Optional) The end time for this scaling rule specific time trigger.
* `launch_expiration_time` - (Optional) The launch expiration time for this scaling rule specific time trigger. The value should between 0 and 3600.
* `recurrence_type` - (Optional) The recurrence type for this scaling rule specific time trigger. Valid values: `MINUTELY`, `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`.
* `recurrence_value` - (Optional) The recurrence value for this scaling rule specific time trigger.

### `node_groups-auto_scaling_policy-scaling_rules-metrics_trigger`

The metrics_trigger mapping supports the following:

* `time_window` - (Required) The time window for this scaling rule specific metrics trigger.
* `evaluation_count` - (Required) The evaluation count for this scaling rule specific metrics trigger.
* `cool_down_interval` - (Optional) The time of cool down interval for this scaling rule specific metrics trigger.
* `condition_logic_operator` - (Optional) The condition logic operator for this scaling rule specific metrics trigger. Valid values: `And` and `Or`.
* `time_constraints` - (Optional) The time constraints for this scaling rule specific metrics trigger. See [`time_constraints`](#node_groups-auto_scaling_policy-scaling_rules-metrics_trigger-time_constraints) below.
* `conditions` - (Optional) The conditions for this scaling rule specific metrics trigger. See [`conditions`](#node_groups-auto_scaling_policy-scaling_rules-metrics_trigger-conditions) below.

### `node_groups-auto_scaling_policy-scaling_rules-metrics_trigger-time_constraints`

The time_constraints mapping supports the following:

* `start_time` - (Optional) The start time for this scaling rule specific metrics trigger.
* `end_time` - (Optional) The end time for this scaling rule specific metrics trigger.

### `node_groups-auto_scaling_policy-scaling_rules-metrics_trigger-conditions`

The conditions mapping supports the following:

* `metric_name` - (Required) The metric name for this scaling rule specific metrics trigger.
* `statistics` - (Required) The statistics for this scaling rule specific metrics trigger.
* `comparison_operator` - (Required) The comparison operator for this scaling rule specific metrics trigger. Invalid values: `EQ`, `NE`, `GT`, `LT`, `GE`, `LE`.
* `threshold` - (Required) The threshold for this scaling rule specific metrics trigger.
* `tags` - (Optional) The tags for this scaling rule specific metrics trigger. See [`tags`](#node_groups-auto_scaling_policy-scaling_rules-metrics_trigger-conditions-tags) below.

### `node_groups-auto_scaling_policy-scaling_rules-metrics_trigger-conditions-tags`

The tags mapping supports the following:

* `key` - (Required) The tag key for this scaling rule specific metrics trigger.
* `value` - (Optional) The tag value for this scaling rule specific metrics trigger.

### `node_groups-auto_scaling_policy-constraints`

The constraints supports the following:

* `max_capacity` - (Optional) The maximum capacity of constraints for emr node group auto scaling policy.
* `min_capacity` - (Optional) The minimum capacity of constraints for emr node group auto scaling policy.

### `bootstrap_scripts`

The bootstrap_scripts mapping supports the following: 

* `script_name` - (Required) The bootstrap script name.
* `script_path` - (Required) The bootstrap script path, e.g. "oss://bucket/path".
* `script_args` - (Required) The bootstrap script args, e.g. "--a=b".
* `priority` - (Deprecated since v1.227.0) The bootstrap scripts priority.
* `execution_moment` - (Required) The bootstrap scripts execution moment, ’BEFORE_INSTALL’, ‘AFTER_STARTED’ or ‘BEFORE_START’. The execution moment of BEFORE_START is available since v1.243.0.
* `execution_fail_strategy` - (Required) The bootstrap scripts execution fail strategy, ’FAILED_BLOCK’ or ‘FAILED_CONTINUE’ .
* `node_selector` - (Required) The bootstrap scripts execution target. See [`node_selector`](#bootstrap_scripts-node_selector) below.

### `bootstrap_scripts-node_selector`

The node_selector mapping supports the following: 

* `node_select_type` - (Required) The bootstrap scripts execution target node select type. Supported value: NODE, NODEGROUP or CLUSTER.
* `node_names` - (Optional) The bootstrap scripts execution target node names.
* `node_group_id` - (Deprecated since v1.227.0) It has been deprecated from version 1.227.0 and new field 'node_group_ids' replaces it.
* `node_group_ids` - (Optional, Available since v1.227.0) The bootstrap scripts execution target node group ids.
* `node_group_types` - (Optional) The bootstrap scripts execution target node group types.
* `node_group_name` - (Deprecated since v1.227.0) It has been deprecated from version 1.227.0 and new field 'node_group_names' replaces it.
* `node_group_names` - (Optional, Available since v1.227.0) The bootstrap scripts execution target node group names.

## Attributes Reference

The following attributes are exported:

* `id` - The emr cluster ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the cluster (until it reaches the initial `RUNNING` status). 
* `delete` - (Defaults to 5 mins) Used when terminating the cluster.

## Import

Aliclioud E-MapReduce cluster can be imported using the id e.g.

```shell
$ terraform import alicloud_emrv2_cluster.default <id>
```