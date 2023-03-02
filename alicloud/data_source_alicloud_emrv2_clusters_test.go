package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEmrV2ClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_emrv2_cluster.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_emrv2_cluster.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_emrv2_cluster.default.id}"]`,
			"name_regex": `"${alicloud_emrv2_cluster.default.cluster_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_emrv2_cluster.default.id}"]`,
			"name_regex": `"${alicloud_emrv2_cluster.default.cluster_name}_fake"`,
		}),
	}

	clusterNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"cluster_name": `"${alicloud_emrv2_cluster.default.cluster_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"cluster_name": `"${alicloud_emrv2_cluster.default.cluster_name}_fake"`,
		}),
	}

	clusterTypesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_emrv2_cluster.default.id}"]`,
			"cluster_types": `["DATALAKE","OLAP"]`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_emrv2_cluster.default.id}"]`,
			"cluster_types": `["DATALEKA","OALP"]`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emrv2_cluster.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emrv2_cluster.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	clusterStatesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emrv2_cluster.default.id}"]`,
			"cluster_states": `["STARTING","RUNNING"]`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emrv2_cluster.default.id}"]`,
			"cluster_states": `["TERMINATING"]`,
		}),
	}

	existTagBytes, _ := json.Marshal(map[string]interface{}{"Created": "TF"})
	fakeTagBytes, _ := json.Marshal(map[string]interface{}{"CreatedFake": "TFFake"})
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_emrv2_cluster.default.id}"]`,
			"tags": string(existTagBytes),
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_emrv2_cluster.default.id}"]`,
			"tags": string(fakeTagBytes),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emrv2_cluster.default.id}"]`,
			"name_regex":        `"${alicloud_emrv2_cluster.default.cluster_name}"`,
			"cluster_types":     `["DATALAKE","OLAP"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"cluster_states":    `["STARTING","RUNNING"]`,
		}),
		fakeConfig: testAccCheckAlicloudEmrV2ClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emrv2_cluster.default.id}_fake"]`,
			"name_regex":        `"${alicloud_emrv2_cluster.default.cluster_name}_fake"`,
			"cluster_types":     `["DATALEKA","OALP"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
			"cluster_states":    `["TERMINATING"]`,
		}),
	}

	var existAlicloudEmrClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"clusters.#":                             "1",
			"ids.0":                                  CHECKSET,
			"names.#":                                "1",
			"total_count":                            CHECKSET,
			"clusters.0.cluster_id":                  CHECKSET,
			"clusters.0.cluster_name":                fmt.Sprintf("tf-testAccClusters-%d", rand),
			"clusters.0.cluster_type":                "OLAP",
			"clusters.0.cluster_state":               "RUNNING",
			"clusters.0.payment_type":                "PayAsYouGo",
			"clusters.0.create_time":                 CHECKSET,
			"clusters.0.ready_time":                  CHECKSET,
			"clusters.0.expire_time":                 CHECKSET,
			"clusters.0.end_time":                    CHECKSET,
			"clusters.0.release_version":             "EMR-5.10.0",
			"clusters.0.resource_group_id":           CHECKSET,
			"clusters.0.tags.#":                      "1",
			"clusters.0.tags.0.key":                  "Created",
			"clusters.0.tags.0.value":                "TF",
			"clusters.0.emr_default_role":            CHECKSET,
			"clusters.0.state_change_reason.Code":    CHECKSET,
			"clusters.0.state_change_reason.Message": CHECKSET,
		}
	}
	var fakeAlicloudEmrClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}
	var alicloudEmrClustersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_emrv2_clusters.default",
		existMapFunc: existAlicloudEmrClustersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEmrClustersDataSourceNameMapFunc,
	}
	alicloudEmrClustersCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, clusterNameConf, clusterTypesConf, resourceGroupIdConf, clusterStatesConf, tagsConf, allConf)
}

func testAccCheckAlicloudEmrV2ClustersDataSourceName(rand int, attrMap map[string]string) string {
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

data "alicloud_emrv2_clusters" "default" {   
   %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
