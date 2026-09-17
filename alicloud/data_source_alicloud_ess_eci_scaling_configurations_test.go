package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEssEciScalingConfigurationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_eci_scaling_configuration.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_eci_scaling_configuration.default.scaling_group_id}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_eci_scaling_configuration.default.scaling_configuration_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_eci_scaling_configuration.default.scaling_configuration_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_eci_scaling_configuration.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_eci_scaling_configuration.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_eci_scaling_configuration.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_eci_scaling_configuration.default.id}"]`,
			"name_regex":       `"${alicloud_ess_eci_scaling_configuration.default.scaling_configuration_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_eci_scaling_configuration.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_eci_scaling_configuration.default.id}_fake"]`,
			"name_regex":       `"${alicloud_ess_eci_scaling_configuration.default.scaling_configuration_name}"`,
		}),
	}

	var existEssEciScalingConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "1",
			"names.#":          "1",
			"names.0":          CHECKSET,
			"ids.0":            CHECKSET,
			"configurations.#": "1",
			// top-level scalar attributes
			"configurations.0.id":                               CHECKSET,
			"configurations.0.scaling_group_id":                 CHECKSET,
			"configurations.0.scaling_configuration_name":       CHECKSET,
			"configurations.0.description":                      "desc",
			"configurations.0.security_group_id":                CHECKSET,
			"configurations.0.container_group_name":             CHECKSET,
			"configurations.0.restart_policy":                   "restartPolicy",
			"configurations.0.cost_optimization":                "true",
			"configurations.0.instance_family_level":            "EntryLevel",
			"configurations.0.cpu_options_core":                 "1",
			"configurations.0.cpu_options_threads_per_core":     "2",
			"configurations.0.cpu":                              "2",
			"configurations.0.memory":                           "4",
			"configurations.0.resource_group_id":                "resourceGroupId",
			"configurations.0.dns_policy":                       "dnsPolicy",
			"configurations.0.enable_sls":                       "true",
			"configurations.0.image_snapshot_id":                "imageSnapshotId",
			"configurations.0.ram_role_name":                    "ramRoleName",
			"configurations.0.termination_grace_period_seconds": "60",
			"configurations.0.auto_match_image_cache":           "true",
			"configurations.0.ipv6_address_count":               "1",
			"configurations.0.active_deadline_seconds":          "60",
			"configurations.0.spot_strategy":                    "SpotWithPriceLimit",
			"configurations.0.spot_price_limit":                 "1.1",
			"configurations.0.auto_create_eip":                  "true",
			"configurations.0.eip_bandwidth":                    "1",
			"configurations.0.ephemeral_storage":                "1",
			"configurations.0.load_balancer_weight":             "1",
			"configurations.0.host_name":                        "hostname",
			"configurations.0.ingress_bandwidth":                "1",
			"configurations.0.egress_bandwidth":                 "1",
			"configurations.0.creation_time":                    CHECKSET,
			"configurations.0.lifecycle_state":                  CHECKSET,
			"configurations.0.tags.name":                        "tf-test",
			// nested blocks: presence + representative sub-fields
			"configurations.0.acr_registry_infos.#":               "1",
			"configurations.0.acr_registry_infos.0.instance_id":   CHECKSET,
			"configurations.0.acr_registry_infos.0.instance_name": "zzz",
			"configurations.0.acr_registry_infos.0.region_id":     "cn-hangzhou",
			"configurations.0.acr_registry_infos.0.domains.#":     "1",

			"configurations.0.image_registry_credentials.#":          "1",
			"configurations.0.image_registry_credentials.0.server":   "server",
			"configurations.0.image_registry_credentials.0.username": "username",

			"configurations.0.dns_config_options.#":       "1",
			"configurations.0.dns_config_options.0.name":  "test",
			"configurations.0.dns_config_options.0.value": "test",

			"configurations.0.security_context_sysctls.#":       "1",
			"configurations.0.security_context_sysctls.0.name":  "kernel.msgmax",
			"configurations.0.security_context_sysctls.0.value": "65536",

			"configurations.0.host_aliases.#":             "1",
			"configurations.0.host_aliases.0.ip":          "ip",
			"configurations.0.host_aliases.0.hostnames.#": "1",

			"configurations.0.volumes.#":                                                   "1",
			"configurations.0.volumes.0.name":                                              "name",
			"configurations.0.volumes.0.type":                                              "type",
			"configurations.0.volumes.0.disk_volume_disk_id":                               "disk_volume_disk_id",
			"configurations.0.volumes.0.disk_volume_fs_type":                               "disk_volume_fs_type",
			"configurations.0.volumes.0.disk_volume_disk_size":                             "1",
			"configurations.0.volumes.0.flex_volume_driver":                                "flex_volume_driver",
			"configurations.0.volumes.0.flex_volume_fs_type":                               "flex_volume_fs_type",
			"configurations.0.volumes.0.flex_volume_options":                               "flex_volume_options",
			"configurations.0.volumes.0.nfs_volume_path":                                   "nfs_volume_path",
			"configurations.0.volumes.0.nfs_volume_read_only":                              "true",
			"configurations.0.volumes.0.nfs_volume_server":                                 "nfs_volume_server",
			"configurations.0.volumes.0.host_path_volume_type":                             "Directory",
			"configurations.0.volumes.0.host_path_volume_path":                             "/etc/test1",
			"configurations.0.volumes.0.empty_dir_volume_medium":                           "memory",
			"configurations.0.volumes.0.empty_dir_volume_size_limit":                       "256 Gi",
			"configurations.0.volumes.0.config_file_volume_config_file_to_paths.#":         "1",
			"configurations.0.volumes.0.config_file_volume_config_file_to_paths.0.content": "content",
			"configurations.0.volumes.0.config_file_volume_config_file_to_paths.0.path":    "path",
			"configurations.0.volumes.0.config_file_volume_config_file_to_paths.0.mode":    CHECKSET,
			"configurations.0.volumes.0.config_file_volume_default_mode":                   CHECKSET,

			"configurations.0.init_containers.#":                                             "1",
			"configurations.0.init_containers.0.name":                                        "name",
			"configurations.0.init_containers.0.image":                                       CHECKSET,
			"configurations.0.init_containers.0.image_pull_policy":                           "policy",
			"configurations.0.init_containers.0.working_dir":                                 "workingDir",
			"configurations.0.init_containers.0.args.#":                                      "1",
			"configurations.0.init_containers.0.commands.#":                                  "1",
			"configurations.0.init_containers.0.cpu":                                         "1",
			"configurations.0.init_containers.0.gpu":                                         "1",
			"configurations.0.init_containers.0.memory":                                      "1",
			"configurations.0.init_containers.0.security_context_read_only_root_file_system": "true",
			"configurations.0.init_containers.0.security_context_run_as_user":                "1",
			"configurations.0.init_containers.0.security_context_capability_adds.#":          "1",
			"configurations.0.init_containers.0.security_context_capability_adds.0":          "adds",
			"configurations.0.init_containers.0.ports.#":                                     "1",
			"configurations.0.init_containers.0.environment_vars.#":                          "1",
			"configurations.0.init_containers.0.environment_vars.0.key":                      "key",
			"configurations.0.init_containers.0.environment_vars.0.value":                    "value",
			"configurations.0.init_containers.0.environment_vars.0.field_ref_field_path":     "path",
			"configurations.0.init_containers.0.volume_mounts.#":                             "1",
			"configurations.0.init_containers.0.volume_mounts.0.mount_path":                  "path",
			"configurations.0.init_containers.0.volume_mounts.0.name":                        "name",
			"configurations.0.init_containers.0.volume_mounts.0.read_only":                   "true",
			"configurations.0.init_containers.0.volume_mounts.0.mount_propagation":           "None",
			"configurations.0.init_containers.0.volume_mounts.0.sub_path":                    "data1/",

			"configurations.0.containers.#":                                             "1",
			"configurations.0.containers.0.name":                                        "name",
			"configurations.0.containers.0.image":                                       CHECKSET,
			"configurations.0.containers.0.image_pull_policy":                           "policy",
			"configurations.0.containers.0.working_dir":                                 "workingDir",
			"configurations.0.containers.0.args.#":                                      "1",
			"configurations.0.containers.0.cpu":                                         "1",
			"configurations.0.containers.0.gpu":                                         "1",
			"configurations.0.containers.0.memory":                                      "1",
			"configurations.0.containers.0.tty":                                         "true",
			"configurations.0.containers.0.stdin":                                       "true",
			"configurations.0.containers.0.security_context_read_only_root_file_system": "true",
			"configurations.0.containers.0.security_context_run_as_user":                "1",
			"configurations.0.containers.0.security_context_capability_adds.#":          "1",
			"configurations.0.containers.0.security_context_capability_adds.0":          "adds",
			"configurations.0.containers.0.lifecycle_pre_stop_handler_execs.#":          "1",
			"configurations.0.containers.0.lifecycle_pre_stop_handler_execs.0":          "echo 1",
			"configurations.0.containers.0.ports.#":                                     "1",
			"configurations.0.containers.0.environment_vars.#":                          "1",
			"configurations.0.containers.0.environment_vars.0.key":                      "key",
			"configurations.0.containers.0.environment_vars.0.value":                    "value",
			"configurations.0.containers.0.environment_vars.0.field_ref_field_path":     "path",
			"configurations.0.containers.0.volume_mounts.#":                             "1",
			"configurations.0.containers.0.volume_mounts.0.mount_path":                  "path",
			"configurations.0.containers.0.volume_mounts.0.name":                        "name",
			"configurations.0.containers.0.volume_mounts.0.read_only":                   "true",
			"configurations.0.containers.0.volume_mounts.0.mount_propagation":           "None",
			"configurations.0.containers.0.volume_mounts.0.sub_path":                    "data1/",
			"configurations.0.containers.0.liveness_probe_exec_commands.#":              "1",
			"configurations.0.containers.0.liveness_probe_period_seconds":               "1",
			"configurations.0.containers.0.liveness_probe_http_get_path":                "path",
			"configurations.0.containers.0.liveness_probe_failure_threshold":            "1",
			"configurations.0.containers.0.liveness_probe_initial_delay_seconds":        "1",
			"configurations.0.containers.0.liveness_probe_http_get_port":                "1",
			"configurations.0.containers.0.liveness_probe_http_get_scheme":              "HTTP",
			"configurations.0.containers.0.liveness_probe_tcp_socket_port":              "1",
			"configurations.0.containers.0.liveness_probe_timeout_seconds":              "1",
			"configurations.0.containers.0.liveness_probe_success_threshold":            "1",
			"configurations.0.containers.0.readiness_probe_exec_commands.#":             "1",
			"configurations.0.containers.0.readiness_probe_period_seconds":              "1",
			"configurations.0.containers.0.readiness_probe_http_get_path":               "path",
			"configurations.0.containers.0.readiness_probe_failure_threshold":           "1",
			"configurations.0.containers.0.readiness_probe_initial_delay_seconds":       "1",
			"configurations.0.containers.0.readiness_probe_http_get_port":               "1",
			"configurations.0.containers.0.readiness_probe_http_get_scheme":             "HTTP",
			"configurations.0.containers.0.readiness_probe_tcp_socket_port":             "1",
			"configurations.0.containers.0.readiness_probe_timeout_seconds":             "1",
			"configurations.0.containers.0.readiness_probe_success_threshold":           "1",
		}
	}

	var fakeEssEciScalingConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"configurations.#": "0",
			"ids.#":            "0",
			"names.#":          "0",
		}
	}

	var essEciScalingConfigurationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_eci_scaling_configurations.default",
		existMapFunc: existEssEciScalingConfigurationsMapFunc,
		fakeMapFunc:  fakeEssEciScalingConfigurationsMapFunc,
	}

	essEciScalingConfigurationsCheckInfo.dataSourceTestCheck(t, rand, scalingGroupIdConf, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tftestaccesseciscalingcfg%d"
}

data "alicloud_security_groups" "default" {
  name_regex     = "^tf_test_acc_alicloud_eci_container_group$"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
	group_type = "ECI"
}

locals {
    alicloud_security_group_id = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.ids.0 : concat(alicloud_security_group.default[*].id, [""])[0]
}

resource "alicloud_cr_ee_instance" "default" {
	payment_type   = "Subscription"
	period         = 1
	renewal_status = "ManualRenewal"
	instance_type  = "Economy"
	instance_name  = "testarc"
	image_scanner  = "DISABLE"
}

resource "alicloud_ess_eci_scaling_configuration" "default" {
	scaling_group_id                 = "${alicloud_ess_scaling_group.default.id}"
	scaling_configuration_name       = "${var.name}"
	description                      = "desc"
	security_group_id                = "${local.alicloud_security_group_id}"
	container_group_name             = "${var.name}"
	restart_policy                   = "restartPolicy"
	cost_optimization                = true
	instance_family_level            = "EntryLevel"
	cpu_options_core                 = 1
	cpu_options_threads_per_core     = 2
	cpu                              = 2
	memory                           = 4
	resource_group_id                = "resourceGroupId"
	dns_policy                       = "dnsPolicy"
	enable_sls                       = true
	image_snapshot_id                = "imageSnapshotId"
	ram_role_name                    = "ramRoleName"
	termination_grace_period_seconds = 60
	auto_match_image_cache           = true
	ipv6_address_count               = 1
	active_deadline_seconds          = 60
	spot_strategy                    = "SpotWithPriceLimit"
	spot_price_limit                 = 1.1
	auto_create_eip                  = true
	eip_bandwidth                    = 1
	ephemeral_storage                = 1
	load_balancer_weight             = 1
	host_name                        = "hostname"
	ingress_bandwidth                = 1
	egress_bandwidth                 = 1
	force_delete                     = true
	tags = {
		name = "tf-test"
	} 
	acr_registry_infos {
		domains       = ["test-registry-vpc.cn-hangzhou.cr.aliyuncs.com"]
		instance_id   = alicloud_cr_ee_instance.default.id
		region_id     = "cn-hangzhou"
		instance_name = "zzz"
	}
	image_registry_credentials {
		password = "password"
		server   = "server"
		username = "username"
	}
	dns_config_options {
		name  = "test"
		value = "test"
	}
	security_context_sysctls {
		name  = "kernel.msgmax"
		value = "65536"
	}
	host_aliases {
		hostnames = ["hostnames"]
		ip        = "ip"
	}
	volumes {
		name                            = "name"
		type                            = "type"
		disk_volume_disk_id             = "disk_volume_disk_id"
		disk_volume_fs_type             = "disk_volume_fs_type"
		disk_volume_disk_size           = 1
		flex_volume_driver              = "flex_volume_driver"
		flex_volume_fs_type             = "flex_volume_fs_type"
		flex_volume_options             = "flex_volume_options"
		nfs_volume_path                 = "nfs_volume_path"
		nfs_volume_read_only            = true
		nfs_volume_server               = "nfs_volume_server"
		host_path_volume_type           = "Directory"
		host_path_volume_path           = "/etc/test1"
		config_file_volume_default_mode = 777
		empty_dir_volume_medium         = "memory"
		empty_dir_volume_size_limit     = "256 Gi"
		config_file_volume_config_file_to_paths {
			content = "content"
			path    = "path"
			mode    = 777
		}
	}
	containers {
		name                                        = "name"
		image                                       = "registry-vpc.aliyuncs.com/eci_open/alpine:3.5"
		image_pull_policy                           = "policy"
		working_dir                                 = "workingDir"
		args                                        = ["arg"]
		cpu                                         = 1
		gpu                                         = 1
		memory                                      = 1
		tty                                         = true
		stdin                                       = true
		security_context_capability_adds            = ["adds"]
		security_context_read_only_root_file_system = true
		security_context_run_as_user                = 1
		lifecycle_pre_stop_handler_execs            = ["echo 1"]
		ports {
			protocol = "protocol"
			port     = 1
		}
		environment_vars {
			key                  = "key"
			value                = "value"
			field_ref_field_path = "path"
		}
		volume_mounts {
			mount_path        = "path"
			name              = "name"
			read_only         = true
			mount_propagation = "None"
			sub_path          = "data1/"
		}
		liveness_probe_exec_commands          = ["cmd"]
		liveness_probe_period_seconds         = 1
		liveness_probe_http_get_path          = "path"
		liveness_probe_failure_threshold      = 1
		liveness_probe_initial_delay_seconds  = 1
		liveness_probe_http_get_port          = 1
		liveness_probe_http_get_scheme        = "HTTP"
		liveness_probe_tcp_socket_port        = 1
		liveness_probe_timeout_seconds        = 1
		liveness_probe_success_threshold      = 1
		readiness_probe_exec_commands         = ["cmd"]
		readiness_probe_period_seconds        = 1
		readiness_probe_http_get_path         = "path"
		readiness_probe_failure_threshold     = 1
		readiness_probe_initial_delay_seconds = 1
		readiness_probe_http_get_port         = 1
		readiness_probe_http_get_scheme       = "HTTP"
		readiness_probe_tcp_socket_port       = 1
		readiness_probe_timeout_seconds       = 1
		readiness_probe_success_threshold     = 1
	}
	init_containers {
		name                                        = "name"
		image                                       = "registry-vpc.aliyuncs.com/eci_open/alpine:3.5"
		image_pull_policy                           = "policy"
		working_dir                                 = "workingDir"
		args                                        = ["arg"]
		commands                                    = ["cmd"]
		cpu                                         = 1
		gpu                                         = 1
		memory                                      = 1
		security_context_capability_adds            = ["adds"]
		security_context_read_only_root_file_system = true
		security_context_run_as_user                = 1
		ports {
			protocol = "protocol"
			port     = 1
		}
		environment_vars {
			key                  = "key"
			value                = "value"
			field_ref_field_path = "path"
		}
		volume_mounts {
			mount_path        = "path"
			name              = "name"
			read_only         = true
			mount_propagation = "None"
			sub_path          = "data1/"
		}
	}
}

data "alicloud_ess_eci_scaling_configurations" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

// TestAccAlicloudEssEciScalingConfigurationsDataSourceInstanceTypes covers the
// fields that belong to the instance-types specification mode (instance_types),
// which is mutually exclusive with the explicit cpu/memory mode exercised by
// TestAccAlicloudEssEciScalingConfigurationsDataSource.
func TestAccAliCloudEssEciScalingConfigurationsDataSourceInstanceTypes(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfigInstanceTypes(rand, map[string]string{
			"ids": `["${alicloud_ess_eci_scaling_configuration.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfigInstanceTypes(rand, map[string]string{
			"ids": `["${alicloud_ess_eci_scaling_configuration.default.id}_fake"]`,
		}),
	}

	var existEssEciScalingConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"configurations.#":                      "1",
			"configurations.0.id":                   CHECKSET,
			"configurations.0.scaling_group_id":     CHECKSET,
			"configurations.0.security_group_id":    CHECKSET,
			"configurations.0.container_group_name": CHECKSET,
			"configurations.0.creation_time":        CHECKSET,
			"configurations.0.instance_types.#":     "2",
			"configurations.0.containers.#":         "1",
		}
	}

	var fakeEssEciScalingConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"configurations.#": "0",
			"ids.#":            "0",
			"names.#":          "0",
		}
	}

	var essEciScalingConfigurationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_eci_scaling_configurations.default",
		existMapFunc: existEssEciScalingConfigurationsMapFunc,
		fakeMapFunc:  fakeEssEciScalingConfigurationsMapFunc,
	}

	essEciScalingConfigurationsCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func testAccCheckAlicloudEssEciScalingConfigurationsDataSourceConfigInstanceTypes(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tftestaccessecitypes%d"
}

data "alicloud_security_groups" "default" {
  name_regex     = "^tf_test_acc_alicloud_eci_container_group$"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 1
	scaling_group_name = "${var.name}"
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
	group_type = "ECI"
}

locals {
    alicloud_security_group_id = length(data.alicloud_security_groups.default.ids) > 0 ? data.alicloud_security_groups.default.ids.0 : concat(alicloud_security_group.default[*].id, [""])[0]
}

resource "alicloud_ess_eci_scaling_configuration" "default" {
	scaling_group_id           = "${alicloud_ess_scaling_group.default.id}"
	scaling_configuration_name = "${var.name}"
	container_group_name       = "${var.name}"
	security_group_id          = "${local.alicloud_security_group_id}"
	force_delete               = true
	instance_types             = ["${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.1.id}"]
	containers {
		name  = "name"
		image = "registry-vpc.aliyuncs.com/eci_open/alpine:3.5"
	}
}

data "alicloud_ess_eci_scaling_configurations" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
