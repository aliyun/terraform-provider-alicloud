package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEmrV2ClusterInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClusterInstancesDataSourceName(rand, map[string]string{
			"cluster_id":       `"${alicloud_emrv2_cluster.default.id}"`,
			"node_group_names": `["${alicloud_emrv2_cluster.default.node_groups.0.node_group_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClusterInstancesDataSourceName(rand, map[string]string{
			"cluster_id":       `"${alicloud_emrv2_cluster.default.id}_fake"`,
			"node_group_names": `["${alicloud_emrv2_cluster.default.node_groups.0.node_group_name}_fake"]`,
		}),
	}

	clusterStatesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClusterInstancesDataSourceName(rand, map[string]string{
			"cluster_id":      `"${alicloud_emrv2_cluster.default.id}"`,
			"instance_states": `["Running"]`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClusterInstancesDataSourceName(rand, map[string]string{
			"cluster_id":      `"${alicloud_emrv2_cluster.default.id}_fake"`,
			"instance_states": `["Stopped"]`,
		}),
	}

	var existAlicloudEmrClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"instances.#":                 "1",
			"names.#":                     "1",
			"total_count":                 "1",
			"instances.0.instance_state":  "Running",
			"instances.0.create_time":     CHECKSET,
			"instances.0.instance_id":     CHECKSET,
			"instances.0.instance_name":   CHECKSET,
			"instances.0.node_group_id":   CHECKSET,
			"instances.0.node_group_type": "MASTER",
		}
	}
	var fakeAlicloudEmrClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
			"total_count": "0",
		}
	}
	var alicloudEmrClustersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_emrv2_cluster_instances.default",
		existMapFunc: existAlicloudEmrClustersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEmrClustersDataSourceNameMapFunc,
	}
	alicloudEmrClustersCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, clusterStatesConf)
}

func testAccCheckAlicloudEmrV2ClusterInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {  
   default = "tf-testAccClusters-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

data "alicloud_zones" "default" {
	available_instance_type = "ecs.g7.xlarge"
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
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
  description = "this is a role test."
  force       = true
}

resource "alicloud_emrv2_cluster" "default" {
  cluster_name = var.name
  cluster_type = "OLAP"
  release_version = "EMR-5.10.0"
  payment_type = "PayAsYouGo"
  deploy_mode = "NORMAL"
  security_mode = "NORMAL"
  resource_group_id = "${data.alicloud_resource_manager_resource_groups.default.ids.0}"
  applications = ["ZOOKEEPER"]

  node_attributes {
    ram_role = "${alicloud_ram_role.default.name}"
    security_group_id = "${alicloud_security_group.default.id}"
    vpc_id = "${alicloud_vpc.default.id}"
    zone_id = "${data.alicloud_zones.default.zones.0.id}"
    key_pair_name = "${alicloud_ecs_key_pair.default.id}"
  }
  
  node_groups {
    node_group_type = "MASTER"
    node_group_name = "emr-master"
    payment_type = "PayAsYouGo"
    vswitch_ids = ["${alicloud_vswitch.default.id}"]
    with_public_ip = false
    instance_types = ["ecs.g7.xlarge"]
    node_count = 1
    system_disk {
      category = "cloud_essd"
      size = 80
      count = 1
    }
    data_disks {
      category = "cloud_essd"
      size = 80
      count = 1
    }
  }

  tags = {
      Created = "TF"
  }
}

data "alicloud_emrv2_cluster_instances" "default" {   
   %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
