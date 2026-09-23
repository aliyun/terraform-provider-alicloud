package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb HadoopDataSource. >>> Resource test cases, automatically generated.
// Case Hadoop数据源_资源依赖_case_03 7604
func TestAccAliCloudGpdbHadoopDataSource_basic7604(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_hadoop_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbHadoopDataSourceMap7604)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbHadoopDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbhadoopdatasource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbHadoopDataSourceBasicDependence7604)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"hadoop_create_type":      "emr",
					"hdfs_conf":               "aaa",
					"yarn_conf":               "aaa",
					"data_source_name":        "${alicloud_gpdb_external_data_service.defaultyOxz1K.service_name}",
					"db_instance_id":          "${alicloud_gpdb_instance.defaultZoepvx.id}",
					"data_source_type":        "hive",
					"hadoop_hosts_address":    "aaa",
					"data_source_description": "test",
					"map_reduce_conf":         "aaa",
					"hadoop_core_conf":        "aaa",
					"hive_conf":               "aaa",
					"emr_instance_id":         "${alicloud_emrv2_cluster.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hadoop_create_type":      "emr",
						"hdfs_conf":               "aaa",
						"yarn_conf":               "aaa",
						"data_source_name":        CHECKSET,
						"db_instance_id":          CHECKSET,
						"data_source_type":        "hive",
						"hadoop_hosts_address":    "aaa",
						"data_source_description": "test",
						"map_reduce_conf":         "aaa",
						"hadoop_core_conf":        "aaa",
						"hive_conf":               "aaa",
						"emr_instance_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hadoop_create_type":      "selfCreate",
					"hdfs_conf":               "aaa2",
					"yarn_conf":               "aaa2",
					"data_source_type":        "hdfs",
					"hadoop_hosts_address":    "aaa2",
					"data_source_description": "aaa2",
					"map_reduce_conf":         "aaa2",
					"hadoop_core_conf":        "aaa2",
					"hive_conf":               "aaa2",
					"emr_instance_id":         "${alicloud_emrv2_cluster.default1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hadoop_create_type":      "selfCreate",
						"hdfs_conf":               "aaa2",
						"yarn_conf":               "aaa2",
						"data_source_type":        "hdfs",
						"hadoop_hosts_address":    "aaa2",
						"data_source_description": "aaa2",
						"map_reduce_conf":         "aaa2",
						"hadoop_core_conf":        "aaa2",
						"hive_conf":               "aaa2",
						"emr_instance_id":         CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudGpdbHadoopDataSourceMap7604 = map[string]string{
	"status":         CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
}

func AlicloudGpdbHadoopDataSourceBasicDependence7604(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "create_pxf_emr_instance_id" {
  default = "c-55fe9a8eed155351"
}

variable "modify_pxf_emr_instance_id" {
  default = "c-1b5a32a3422bafa4"
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

resource "alicloud_emrv2_cluster" "default1" {
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
  instance_spec              = "2C8G"
  description                = "创建依赖的Local实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = "cn-beijing-h"
  vswitch_id                 = data.alicloud_vswitches.default.ids[0]
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = data.alicloud_vpcs.default.ids.0
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
  db_instance_category       = "Basic"
}

resource "alicloud_gpdb_external_data_service" "defaultyOxz1K" {
  service_name        = var.name
  db_instance_id      = alicloud_gpdb_instance.defaultZoepvx.id
  service_description = "mytest"
  service_spec        = "8"
}


`, name)
}

// Test Gpdb HadoopDataSource. <<< Resource test cases, automatically generated.
