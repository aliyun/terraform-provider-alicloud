package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEmrClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}_fake"]`,
			"enable_details": "true",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"name_regex":     `"${alicloud_emr_cluster.default.name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"name_regex":     `"${alicloud_emr_cluster.default.name}_fake"`,
			"enable_details": "true",
		}),
	}

	clusterNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"cluster_name":   `"${alicloud_emr_cluster.default.name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"cluster_name":   `"${alicloud_emr_cluster.default.name}_fake"`,
			"enable_details": "true",
		}),
	}

	clusterTypeListConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"cluster_type_list": `["HADOOP","KAFKA"]`,
			"enable_details":    "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"cluster_type_list": `["ZOOKPEER","CLICKHOSUE"]`,
			"enable_details":    "true",
		}),
	}

	createTypeListConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"create_type":    `"MANUAL"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"create_type":    `"ON-DEMAND"`,
			"enable_details": "true",
		}),
	}

	defaultStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"default_status": `"true"`,
			"enable_details": "true",
		}),
	}

	depositTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"deposit_type":   `"HALF_MANAGED"`,
			"enable_details": "true",
		}),
	}

	isDescConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"is_desc":        `"false"`,
			"enable_details": "true",
		}),
	}

	machineTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"machine_type":   `"ECS"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"machine_type":   `"DOCKER"`,
			"enable_details": "true",
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"enable_details":    "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
			"enable_details":    "true",
		}),
	}

	statusListConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"status_list":    `["CREATING","IDLE"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"status_list":    `["ABNORMAL","RELEASING"]`,
			"enable_details": "true",
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"vpc_id":         `"${data.alicloud_vpcs.default.ids.0}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"vpc_id":         `"${data.alicloud_vpcs.default.ids.0}_fake"`,
			"enable_details": "true",
		}),
	}

	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"cluster_name":   `"${alicloud_emr_cluster.default.name}"`,
			"page_number":    `1`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"cluster_name":   `"${alicloud_emr_cluster.default.name}"`,
			"page_number":    `2`,
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"name_regex":        `"${alicloud_emr_cluster.default.name}"`,
			"cluster_type_list": `["HADOOP","KAFKA"]`,
			"create_type":       `"MANUAL"`,
			"deposit_type":      `"HALF_MANAGED"`,
			"machine_type":      `"ECS"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"status_list":       `["CREATING","IDLE"]`,
			"default_status":    `"true"`,
			"is_desc":           `"false"`,
			"vpc_id":            `"${data.alicloud_vpcs.default.ids.0}"`,
			"enable_details":    "true",
			"page_number":       `1`,
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}_fake"]`,
			"name_regex":        `"${alicloud_emr_cluster.default.name}_fake"`,
			"cluster_type_list": `["ZOOKPEER","CLICKHOSUE"]`,
			"create_type":       `"ON-DEMAND"`,
			"machine_type":      `"DOCKER"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
			"status_list":       `["ABNORMAL","RELEASING"]`,
			"vpc_id":            `"${data.alicloud_vpcs.default.ids.0}_fake"`,
			"enable_details":    "true",
			"page_number":       `2`,
		}),
	}

	var existAlicloudEmrClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"clusters.#":                      "1",
			"ids.0":                           CHECKSET,
			"names.#":                         "1",
			"total_count":                     CHECKSET,
			"clusters.0.id":                   CHECKSET,
			"clusters.0.auto_scaling_allowed": CHECKSET,
			"clusters.0.auto_scaling_by_load_allowed":                   CHECKSET,
			"clusters.0.auto_scaling_enable":                            CHECKSET,
			"clusters.0.auto_scaling_spot_with_limit_allowed":           CHECKSET,
			"clusters.0.bootstrap_action_list.#":                        "0",
			"clusters.0.bootstrap_failed":                               CHECKSET,
			"clusters.0.cluster_id":                                     CHECKSET,
			"clusters.0.cluster_name":                                   fmt.Sprintf("tf-testAccClusters-%d", rand),
			"clusters.0.create_resource":                                "ECM_EMR",
			"clusters.0.create_time":                                    CHECKSET,
			"clusters.0.create_type":                                    "MANUAL",
			"clusters.0.deposit_type":                                   "HALF_MANAGED",
			"clusters.0.eas_enable":                                     CHECKSET,
			"clusters.0.expired_time":                                   CHECKSET,
			"clusters.0.extra_info":                                     CHECKSET,
			"clusters.0.high_availability_enable":                       CHECKSET,
			"clusters.0.host_group_list.#":                              "3",
			"clusters.0.host_group_list.0.band_width":                   "",
			"clusters.0.host_group_list.0.charge_type":                  CHECKSET,
			"clusters.0.host_group_list.0.cpu_core":                     CHECKSET,
			"clusters.0.host_group_list.0.disk_capacity":                CHECKSET,
			"clusters.0.host_group_list.0.disk_count":                   CHECKSET,
			"clusters.0.host_group_list.0.disk_type":                    CHECKSET,
			"clusters.0.host_group_list.0.host_group_change_type":       "",
			"clusters.0.host_group_list.0.host_group_id":                CHECKSET,
			"clusters.0.host_group_list.0.host_group_name":              CHECKSET,
			"clusters.0.host_group_list.0.host_group_type":              CHECKSET,
			"clusters.0.host_group_list.0.instance_type":                CHECKSET,
			"clusters.0.host_group_list.0.memory_capacity":              CHECKSET,
			"clusters.0.host_group_list.0.node_count":                   CHECKSET,
			"clusters.0.host_group_list.0.nodes.#":                      CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.create_time":          CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.disk_infos.#":         CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.disk_infos.0.device":  CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.disk_infos.0.disk_id": CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.disk_infos.0.size":    CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.disk_infos.0.type":    CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.emr_expired_time":     CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.expired_time":         CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.inner_ip":             CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.instance_id":          CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.pub_ip":               CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.status":               CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.support_ipv6":         CHECKSET,
			"clusters.0.host_group_list.0.nodes.0.zone_id":              CHECKSET,
			"clusters.0.host_group_list.0.period":                       "",
			"clusters.0.image_id":                                       CHECKSET,
			"clusters.0.local_meta_db":                                  CHECKSET,
			"clusters.0.machine_type":                                   "ECS",
			"clusters.0.meta_store_type":                                CHECKSET,
			"clusters.0.net_type":                                       CHECKSET,
			"clusters.0.payment_type":                                   CHECKSET,
			"clusters.0.period":                                         CHECKSET,
			"clusters.0.resize_disk_enable":                             CHECKSET,
			"clusters.0.running_time":                                   CHECKSET,
			"clusters.0.security_group_id":                              CHECKSET,
			"clusters.0.security_group_name":                            CHECKSET,
			"clusters.0.software_info.#":                                "1",
			"clusters.0.software_info.0.cluster_type":                   CHECKSET,
			"clusters.0.software_info.0.emr_ver":                        CHECKSET,
			"clusters.0.software_info.0.softwares.#":                    CHECKSET,
			"clusters.0.software_info.0.softwares.0.display_name":       CHECKSET,
			"clusters.0.software_info.0.softwares.0.name":               CHECKSET,
			"clusters.0.software_info.0.softwares.0.only_display":       CHECKSET,
			"clusters.0.software_info.0.softwares.0.start_tpe":          CHECKSET,
			"clusters.0.software_info.0.softwares.0.version":            CHECKSET,
			"clusters.0.start_time":                                     CHECKSET,
			"clusters.0.status":                                         "IDLE",
			"clusters.0.stop_time":                                      "",
			"clusters.0.tags.%":                                         "2",
			"clusters.0.tags.Created":                                   "TF",
			"clusters.0.tags.For":                                       "acceptance test",
			"clusters.0.user_defined_emr_ecs_role":                      CHECKSET,
			"clusters.0.user_id":                                        CHECKSET,
			"clusters.0.vswitch_id":                                     CHECKSET,
			"clusters.0.vpc_id":                                         CHECKSET,
			"clusters.0.zone_id":                                        CHECKSET,
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
		resourceId:   "data.alicloud_emr_clusters.default",
		existMapFunc: existAlicloudEmrClustersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEmrClustersDataSourceNameMapFunc,
	}
	alicloudEmrClustersCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, clusterNameConf, clusterTypeListConf, createTypeListConf, defaultStatusConf, depositTypeConf, isDescConf, machineTypeConf, resourceGroupIdConf, statusListConf, vpcIdConf, allConf, pagingConf)
}

func TestAccAlicloudEmrClustersDataSource1(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"enable_details": "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}_fake"]`,
			"enable_details": "false",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"name_regex":     `"${alicloud_emr_cluster.default.name}"`,
			"enable_details": "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"name_regex":     `"${alicloud_emr_cluster.default.name}_fake"`,
			"enable_details": "false",
		}),
	}

	clusterNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"cluster_name":   `"${alicloud_emr_cluster.default.name}"`,
			"enable_details": "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"cluster_name":   `"${alicloud_emr_cluster.default.name}_fake"`,
			"enable_details": "false",
		}),
	}

	clusterTypeListConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"cluster_type_list": `["HADOOP","KAFKA"]`,
			"enable_details":    "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"cluster_type_list": `["ZOOKPEER","CLICKHOSUE"]`,
			"enable_details":    "false",
		}),
	}

	createTypeListConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"create_type":    `"MANUAL"`,
			"enable_details": "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"create_type":    `"ON-DEMAND"`,
			"enable_details": "false",
		}),
	}

	defaultStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"default_status": `"true"`,
			"enable_details": "false",
		}),
	}

	depositTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"deposit_type":   `"HALF_MANAGED"`,
			"enable_details": "false",
		}),
	}

	isDescConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"is_desc":        `"false"`,
			"enable_details": "false",
		}),
	}

	machineTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"machine_type":   `"ECS"`,
			"enable_details": "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"machine_type":   `"DOCKER"`,
			"enable_details": "false",
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"enable_details":    "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
			"enable_details":    "false",
		}),
	}

	statusListConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"status_list":    `["CREATING","IDLE"]`,
			"enable_details": "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"status_list":    `["ABNORMAL","RELEASING"]`,
			"enable_details": "false",
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"vpc_id":         `"${data.alicloud_vpcs.default.ids.0}"`,
			"enable_details": "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_emr_cluster.default.id}"]`,
			"vpc_id":         `"${data.alicloud_vpcs.default.ids.0}_fake"`,
			"enable_details": "false",
		}),
	}

	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"cluster_name":   `"${alicloud_emr_cluster.default.name}"`,
			"page_number":    `1`,
			"enable_details": "false",
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"cluster_name":   `"${alicloud_emr_cluster.default.name}"`,
			"page_number":    `2`,
			"enable_details": "false",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}"]`,
			"name_regex":        `"${alicloud_emr_cluster.default.name}"`,
			"cluster_type_list": `["HADOOP","KAFKA"]`,
			"create_type":       `"MANUAL"`,
			"deposit_type":      `"HALF_MANAGED"`,
			"machine_type":      `"ECS"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"status_list":       `["CREATING","IDLE"]`,
			"default_status":    `"true"`,
			"is_desc":           `"false"`,
			"vpc_id":            `"${data.alicloud_vpcs.default.ids.0}"`,
			"enable_details":    "false",
			"page_number":       `1`,
		}),
		fakeConfig: testAccCheckAlicloudEmrClustersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_emr_cluster.default.id}_fake"]`,
			"name_regex":        `"${alicloud_emr_cluster.default.name}_fake"`,
			"cluster_type_list": `["ZOOKPEER","CLICKHOSUE"]`,
			"create_type":       `"ON-DEMAND"`,
			"machine_type":      `"DOCKER"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
			"status_list":       `["ABNORMAL","RELEASING"]`,
			"vpc_id":            `"${data.alicloud_vpcs.default.ids.0}_fake"`,
			"enable_details":    "false",
			"page_number":       `2`,
		}),
	}

	var existAlicloudEmrClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"clusters.#":                       "1",
			"ids.0":                            CHECKSET,
			"names.#":                          "1",
			"clusters.0.id":                    CHECKSET,
			"clusters.0.cluster_id":            CHECKSET,
			"clusters.0.cluster_name":          fmt.Sprintf("tf-testAccClusters-%d", rand),
			"clusters.0.create_resource":       "ECM_EMR",
			"clusters.0.create_time":           CHECKSET,
			"clusters.0.deposit_type":          "HALF_MANAGED",
			"clusters.0.expired_time":          CHECKSET,
			"clusters.0.machine_type":          "ECS",
			"clusters.0.meta_store_type":       CHECKSET,
			"clusters.0.payment_type":          CHECKSET,
			"clusters.0.period":                CHECKSET,
			"clusters.0.running_time":          CHECKSET,
			"clusters.0.status":                "IDLE",
			"clusters.0.type":                  CHECKSET,
			"clusters.0.has_uncompleted_order": CHECKSET,
			"clusters.0.tags.%":                "2",
			"clusters.0.tags.Created":          "TF",
			"clusters.0.tags.For":              "acceptance test",
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
		resourceId:   "data.alicloud_emr_clusters.default",
		existMapFunc: existAlicloudEmrClustersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEmrClustersDataSourceNameMapFunc,
	}
	alicloudEmrClustersCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, clusterNameConf, clusterTypeListConf, createTypeListConf, defaultStatusConf, depositTypeConf, isDescConf, machineTypeConf, resourceGroupIdConf, statusListConf, vpcIdConf, pagingConf, allConf)
}

func testAccCheckAlicloudEmrClustersDataSourceName(rand int, attrMap map[string]string) string {
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

data "alicloud_emr_main_versions" "default" {}

data "alicloud_emr_instance_types" "default" {
  destination_resource  = "InstanceType"
  cluster_type          = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  support_local_storage = false
  instance_charge_type  = "PostPaid"
  support_node_type     = ["MASTER", "CORE", "TASK"]
}

data "alicloud_emr_disk_types" "data_disk" {
  destination_resource = "DataDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
  destination_resource = "SystemDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = data.alicloud_vpcs.default.ids.0
}


data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
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

resource "alicloud_emr_cluster" "default" {
  name = var.name

  emr_ver = data.alicloud_emr_main_versions.default.main_versions.0.emr_version

  cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0

  host_group {
    host_group_name   = "master_group"
    host_group_type   = "MASTER"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "1"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "core_group"
    host_group_type   = "CORE"
    node_count        = "3"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "task_group"
    host_group_type   = "TASK"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  high_availability_enable  = true
  zone_id                   = data.alicloud_emr_instance_types.default.types.0.zone_id
  security_group_id         = alicloud_security_group.default.id
  is_open_public_ip         = true
  charge_type               = "PostPaid"
  vswitch_id                = data.alicloud_vswitches.default.ids.0
  user_defined_emr_ecs_role = alicloud_ram_role.default.name
  ssh_enable                = true
  master_pwd                = "ABCtest1234!"
  tags = {
      Created = "TF"
      For    = "acceptance test"
  }
}

data "alicloud_emr_clusters" "default" {   
   %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
