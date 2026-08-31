data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_instance_type = "ecs.g7.xlarge"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.vpc_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = var.vswitch_cidr
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.vswitch_name
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = var.key_pair_name
}

resource "alicloud_security_group" "default" {
  name   = var.security_group_name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ram_role" "default" {
  name        = var.ram_name
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
  description = "this is a role test."
  force       = true
}

resource "alicloud_emrv2_cluster" "default" {
  payment_type    = "PayAsYouGo"
  cluster_type    = "DATALAKE"
  release_version = "EMR-5.10.0"
  cluster_name    = var.cluster_name
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