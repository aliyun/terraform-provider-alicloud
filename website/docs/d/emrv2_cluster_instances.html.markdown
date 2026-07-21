---
subcategory: "E-MapReduce (EMR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_emrv2_cluster_instances"
sidebar_current: "docs-alicloud-datasource-emr-cluster-instances"
description: |-
  Provides a list of Emr Cluster ecs instances to the user based on EMR's new version OpenAPI.
---

# alicloud_emrv2_cluster_instances

This data source provides the Emr Cluster ecs instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_instance_type = "ecs.g7.xlarge"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "TF-VPC"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = "TF_VSwitch"
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = "terraform-kp"
}

resource "alicloud_security_group" "default" {
  name   = "TF_SECURITY_GROUP"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ram_role" "default" {
  name     = "emrtf"
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

  description = "this is a role test."
  force       = true
}

resource "alicloud_emrv2_cluster" "default" {
  payment_type    = "PayAsYouGo"
  cluster_type    = "DATALAKE"
  release_version = "EMR-5.10.0"
  cluster_name    = "terraform-emr-cluster-v2"
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

data "alicloud_emrv2_cluster_instances" "ids" {}
output "emrv2_cluster_instances_id_1" {
  value = data.alicloud_emrv2_cluster_instances.ids.instances.0.instance_id
}

data "alicloud_emrv2_cluster_instances" "nodeGroupNames" {
  node_group_names = ["emr-core"]
}
output "emrv2_cluster_instances_id_2" {
  value = data.alicloud_emrv2_cluster_instances.nodeGroupNames.instances.0.instance_id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew)  A list of Cluster IDs.
* `cluster_id` - (Optional, ForceNew) The emr cluster ID.
* `node_group_ids` - (Optional, ForceNew) The cluster node group ids.
* `node_group_names` - (Optional, ForceNew) The cluster node group names.
* `instance_states` - (Optional, ForceNew) The cluster ecs instance states.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `next_token` - (Optional) The next token is used to list emr cluster ecs instances for next page.
* `max_results` - (Optional) The max results is used to list emr cluster ecs instances for next page.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Cluster ecs instance names.
* `ids` - A list of Cluster ecs instance IDS.
* `instances` - A list of Emr Cluster ecs instances. Each element contains the following attributes:
	* `instance_id` - The emr cluster ecs instance ID.
	* `instance_name` - The emr cluster ecs instance name.
	* `instance_type` - The emr cluster ecs instance type.
	* `instance_state` - The emr cluster ecs instance state.
	* `node_group_id` - The emr cluster node group ID.
	* `node_group_type` - The emr cluster node group type.
	* `zone_id` - The emr cluster node group zone ID.
	* `public_ip` - The emr cluster ecs instance public ip.
	* `private_ip` - The emr cluster ecs instance private ip.
	* `auto_renew` - The emr cluster node group whether auto renew when payment type is 'Subscription'.
	* `auto_renew_duration_unit` - The emr cluster node group auto renew duration unit when payment type is 'Subscription'.
	* `auto_renew_duration` - The emr cluster node group auto renew duration when payment type is 'Subscription'.
	* `create_time` - The creation time of the resource.
	* `expire_time` - The expire time of the resource.
* `total_count` - The total count of list emr cluster ecs instances.