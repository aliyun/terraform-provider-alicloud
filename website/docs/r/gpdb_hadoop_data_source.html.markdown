---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_hadoop_data_source"
description: |-
  Provides a Alicloud GPDB Hadoop Data Source resource.
---

# alicloud_gpdb_hadoop_data_source

Provides a GPDB Hadoop Data Source resource.

Hadoop DataSource Config.

For information about GPDB Hadoop Data Source and how to use it, see [What is Hadoop Data Source](https://www.alibabacloud.com/help/en/analyticdb/analyticdb-for-postgresql/developer-reference/api-gpdb-2016-05-03-createhadoopdatasource).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_hadoop_data_source&exampleId=272afeb1-2bde-883a-9063-f94812d51a2d804ac0d4&activeTab=example&spm=docs.r.gpdb_hadoop_data_source.0.272afeb12b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-h"
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
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

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_kms_keys" "default" {
  status = "Enabled"
}

resource "alicloud_emrv2_cluster" "default" {
  node_groups {
    vswitch_ids = [
      data.alicloud_vswitches.default.ids[0]
    ]
    instance_types = [
      "ecs.g6.xlarge"
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
      data.alicloud_vswitches.default.ids[0]
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
      "ecs.g6.xlarge"
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
    zone_id              = "cn-beijing-h"
    key_pair_name        = alicloud_ecs_key_pair.default.id
    data_disk_encrypted  = "true"
    data_disk_kms_key_id = data.alicloud_kms_keys.default.ids.0
    vpc_id               = data.alicloud_vpcs.default.ids.0
    ram_role             = alicloud_ram_role.default.name
    security_group_id    = alicloud_security_group.default.id
  }

  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  cluster_name      = var.name
  payment_type      = "PayAsYouGo"
  cluster_type      = "DATAFLOW"
}

resource "alicloud_gpdb_instance" "defaultZoepvx" {
  instance_spec         = "2C8G"
  description           = var.name
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = "cn-beijing-h"
  vswitch_id            = data.alicloud_vswitches.default.ids[0]
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = data.alicloud_vpcs.default.ids.0
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
  db_instance_category  = "Basic"
}

resource "alicloud_gpdb_external_data_service" "defaultyOxz1K" {
  service_name        = var.name
  db_instance_id      = alicloud_gpdb_instance.defaultZoepvx.id
  service_description = var.name
  service_spec        = "8"
}

resource "alicloud_gpdb_hadoop_data_source" "default" {
  hdfs_conf               = "aaa"
  data_source_name        = alicloud_gpdb_external_data_service.defaultyOxz1K.service_name
  yarn_conf               = "aaa"
  hive_conf               = "aaa"
  hadoop_create_type      = "emr"
  data_source_description = var.name
  map_reduce_conf         = "aaa"
  data_source_type        = "hive"
  hadoop_core_conf        = "aaa"
  emr_instance_id         = alicloud_emrv2_cluster.default.id
  db_instance_id          = alicloud_gpdb_instance.defaultZoepvx.id
  hadoop_hosts_address    = "aaa"
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The instance ID.

* `data_source_description` - (Optional) Data Source Description
* `data_source_name` - (Optional, ForceNew) Data Source Name
* `data_source_type` - (Optional) The type of the data source. Valid values:

  *   mysql
  - postgresql

  *   hdfs
  - hive
* `emr_instance_id` - (Optional) The ID of the Emr instance.
* `hadoop_core_conf` - (Optional) The string that specifies the content of the Hadoop core-site.xml file.

* `hadoop_create_type` - (Optional) The type of the external service. Valid values:
  - emr: E-MapReduce (EMR) Hadoop cluster.
  - selfCreate: self-managed Hadoop cluster.

* `hadoop_hosts_address` - (Optional) The IP address and hostname of the Hadoop cluster (data source) in the /etc/hosts file.

* `hdfs_conf` - (Optional) The string that specifies the content of the Hadoop hdfs-site.xml file. This parameter must be specified when DataSourceType is set to HDFS.

* `hive_conf` - (Optional) The string that specifies the content of the Hadoop hive-site.xml file. This parameter must be specified when DataSourceType is set to Hive.

* `map_reduce_conf` - (Optional) The content of the Hadoop mapred-site.xml file. This parameter must be specified when DataSourceType is set to HDFS.

* `yarn_conf` - (Optional) The string that specifies the content of the Hadoop yarn-site.xml file. This parameter must be specified when DataSourceType is set to HDFS.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<data_source_id>`.
* `create_time` - Creation time
* `data_source_id` - The data source ID.

* `status` - Data Source Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Hadoop Data Source.
* `delete` - (Defaults to 5 mins) Used when delete the Hadoop Data Source.
* `update` - (Defaults to 5 mins) Used when update the Hadoop Data Source.

## Import

GPDB Hadoop Data Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_hadoop_data_source.example <db_instance_id>:<data_source_id>
```