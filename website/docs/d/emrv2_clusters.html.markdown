---
subcategory: "E-MapReduce (EMR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_emrv2_clusters"
sidebar_current: "docs-alicloud-datasource-emrv2-clusters"
description: |-
  Provides a list of Emr Clusters to the user based on EMR's new version OpenAPI.
---

# alicloud_emrv2_clusters

This data source provides the Emr Clusters of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.199.0.

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

data "alicloud_emrv2_clusters" "ids" {}
output "emrv2_clusters_id_1" {
  value = data.alicloud_emrv2_clusters.ids.clusters.0.id
}

data "alicloud_emrv2_clusters" "nameRegex" {
  name_regex = alicloud_emr_cluster.default.name
}
output "emrv2_clusters_id_2" {
  value = data.alicloud_emrv2_clusters.nameRegex.clusters.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Cluster IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Cluster name.
* `cluster_name` - (Optional, ForceNew) The cluster name.
* `resource_group_id` - (Optional, ForceNew) The Resource Group ID.
* `cluster_types` - (Optional, ForceNew) The cluster types.
* `cluster_states` - (Optional, ForceNew) The cluster states.
* `payment_types` - (Optional, ForceNew) The cluster payment types.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `next_token` - (Optional) The next token is used to list clusters for next page.
* `max_results` - (Optional) The max results is used to list clusters for next page.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Cluster names.
* `ids` - A list of Cluster IDS.
* `clusters` - A list of Emr Clusters. Each element contains the following attributes:
	* `cluster_id` - The first ID of the resource.
	* `cluster_name` - The name of the emr cluster.
	* `cluster_type` - The type of the emr cluster.
	* `cluster_state` - The state of the emr cluster.
	* `payment_type` - The payment type of the emr cluster.
	* `create_time` - The creation time of the resource.
	* `ready_time` - The ready time of the resource.
	* `expire_time` - The expire time of the resource.
	* `end_time` - The end time of the resource.
	* `release_version` - The release version of the resource.
	* `resource_group_id` - The resource group id of the resource.
	* `state_change_reason` - The cluster state change reason.
	* `tags` - A mapping of tags to assign to the resource.
	* `emr_default_role` - The ecs default role belongs to this emr cluster.
* `total_count` - The total count of list clusters.