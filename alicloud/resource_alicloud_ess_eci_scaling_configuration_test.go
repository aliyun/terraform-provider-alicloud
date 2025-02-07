package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ess_eci_scaling_configuration",
		&resource.Sweeper{
			Name: "alicloud_ess_eci_scaling_configuration",
			F:    testSweepEciScalingConfiguration,
		})
}

func testSweepEciScalingConfiguration(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc`",
	}
	var response map[string]interface{}
	action := "DescribeEciScalingConfigurations"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	response, err = client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_eci_scaling_configuration", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	resp, err := jsonpath.Get("$.ScalingConfigurations", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ScalingConfigurations", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["ScalingConfigurationName"].(string)
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Eci configuration: %s ", name)
			continue
		}
		log.Printf("[INFO] Delete Eci configuration: %s ", name)
		action := "DeleteEciScalingConfiguration"
		request := map[string]interface{}{
			"ScalingConfigurationId": item["ScalingConfigurationId"],
		}
		request["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*10, func() *resource.RetryError {
			response, err = client.RpcPost("Ess", "2014-08-28", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Eci Scaling Configuration (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAliCloudEssEciScalingConfigurationBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ess_eci_scaling_configuration.default"
	ra := resourceAttrInit(resourceId, AlicloudEssEciScalingConfigurationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEssEciScalingConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-test-acc-alicloud-eci-container-group%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssEciScalingConfiguration)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                 "${alicloud_ess_scaling_group.default.id}",
					"scaling_configuration_name":       name,
					"description":                      "desc",
					"security_group_id":                "${local.alicloud_security_group_id}",
					"container_group_name":             name,
					"restart_policy":                   "restartPolicy",
					"cpu":                              "2",
					"memory":                           "4",
					"resource_group_id":                "resourceGroupId",
					"dns_policy":                       "dnsPolicy",
					"enable_sls":                       "true",
					"image_snapshot_id":                "imageSnapshotId",
					"ram_role_name":                    "ramRoleName",
					"termination_grace_period_seconds": "60",
					"auto_match_image_cache":           "true",
					"ipv6_address_count":               "1",
					"active_deadline_seconds":          "60",
					"spot_strategy":                    "SpotWithPriceLimit",
					"spot_price_limit":                 "1.1",
					"auto_create_eip":                  "true",
					"eip_bandwidth":                    "1",
					"ephemeral_storage":                "1",
					"load_balancer_weight":             "1",
					"host_name":                        "hostname",
					"ingress_bandwidth":                "1",
					"egress_bandwidth":                 "1",
					"tags": map[string]string{
						"name": "tf-test",
					},
					"acr_registry_infos": []map[string]interface{}{
						{
							"domains":       []string{"test-registry-vpc.cn-hangzhou.cr.aliyuncs.com"},
							"instance_id":   "cri-47rme9691uiowvfv",
							"region_id":     "cn-hangzhou",
							"instance_name": "zzz",
						},
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "password",
							"server":   "server",
							"username": "username",
						},
					},
					"containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"lifecycle_pre_stop_handler_execs":            []string{"echo 1"},
							"security_context_read_only_root_file_system": "true",
							"security_context_run_as_user":                "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "path",
									"name":       "name",
									"read_only":  "true",
								},
							},
							"liveness_probe_exec_commands":          []string{"cmd"},
							"liveness_probe_period_seconds":         "1",
							"liveness_probe_http_get_path":          "path",
							"liveness_probe_failure_threshold":      "1",
							"liveness_probe_initial_delay_seconds":  "1",
							"liveness_probe_http_get_port":          "1",
							"liveness_probe_http_get_scheme":        "HTTP",
							"liveness_probe_tcp_socket_port":        "1",
							"liveness_probe_timeout_seconds":        "1",
							"readiness_probe_exec_commands":         []string{"cmd"},
							"readiness_probe_period_seconds":        "1",
							"readiness_probe_http_get_path":         "path",
							"readiness_probe_failure_threshold":     "1",
							"readiness_probe_initial_delay_seconds": "1",
							"readiness_probe_http_get_port":         "1",
							"readiness_probe_http_get_scheme":       "HTTP",
							"readiness_probe_tcp_socket_port":       "1",
							"readiness_probe_timeout_seconds":       "1",
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"security_context_read_only_root_file_system": "true",
							"security_context_run_as_user":                "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "path",
									"name":       "name",
									"read_only":  "true",
								},
							},
							"commands": []string{"cmd"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content",
									"path":    "path",
								},
							},
							"disk_volume_disk_id":   "disk_volume_disk_id",
							"disk_volume_fs_type":   "disk_volume_fs_type",
							"disk_volume_disk_size": "1",
							"flex_volume_driver":    "flex_volume_driver",
							"flex_volume_fs_type":   "flex_volume_fs_type",
							"flex_volume_options":   "flex_volume_options",
							"nfs_volume_path":       "nfs_volume_path",
							"nfs_volume_read_only":  "true",
							"nfs_volume_server":     "nfs_volume_server",
							"name":                  "name",
							"type":                  "type",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames"},
							"ip":        "ip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":                 CHECKSET,
						"scaling_configuration_name":       name,
						"description":                      "desc",
						"security_group_id":                CHECKSET,
						"container_group_name":             name,
						"restart_policy":                   "restartPolicy",
						"cpu":                              "2",
						"memory":                           "4",
						"resource_group_id":                "resourceGroupId",
						"dns_policy":                       "dnsPolicy",
						"enable_sls":                       "true",
						"image_snapshot_id":                "imageSnapshotId",
						"ram_role_name":                    "ramRoleName",
						"termination_grace_period_seconds": "60",
						"auto_match_image_cache":           "true",
						"ipv6_address_count":               "1",
						"active_deadline_seconds":          "60",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit":                 "1.1",
						"auto_create_eip":                  "true",
						"host_name":                        "hostname",
						"ingress_bandwidth":                "1",
						"egress_bandwidth":                 "1",
						"ephemeral_storage":                "1",
						"load_balancer_weight":             "1",
						"tags.name":                        "tf-test",
						"image_registry_credentials.#":     "1",
						"containers.#":                     "1",
						"containers.0.security_context_read_only_root_file_system":      "true",
						"containers.0.security_context_run_as_user":                     "1",
						"containers.0.security_context_capability_adds.#":               "1",
						"containers.0.security_context_capability_adds.0":               "adds",
						"containers.0.lifecycle_pre_stop_handler_execs.#":               "1",
						"containers.0.lifecycle_pre_stop_handler_execs.0":               "echo 1",
						"containers.0.environment_vars.#":                               "1",
						"init_containers.#":                                             "1",
						"init_containers.0.security_context_run_as_user":                "1",
						"init_containers.0.security_context_read_only_root_file_system": "true",
						"init_containers.0.security_context_capability_adds.#":          "1",
						"init_containers.0.security_context_capability_adds.0":          "adds",
						"volumes.#":            "1",
						"host_aliases.#":       "1",
						"acr_registry_infos.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_configuration_name":       "newName",
					"active":                           "true",
					"force_delete":                     "true",
					"description":                      "newDesc",
					"security_group_id":                "${local.alicloud_security_group_id1}",
					"container_group_name":             "1new-name",
					"restart_policy":                   "newPolicy",
					"cpu":                              "4",
					"memory":                           "8",
					"resource_group_id":                "newGroupId",
					"dns_policy":                       "newDnsPolicy",
					"enable_sls":                       "false",
					"image_snapshot_id":                "imageSnapshotId2",
					"ram_role_name":                    "newRoleName",
					"termination_grace_period_seconds": "120",
					"auto_match_image_cache":           "false",
					"ipv6_address_count":               "2",
					"active_deadline_seconds":          "120",
					"spot_strategy":                    "SpotAsPriceGo",
					"spot_price_limit":                 "1.2",
					"auto_create_eip":                  "false",
					"eip_bandwidth":                    "3",
					"host_name":                        "newHostName",
					"ingress_bandwidth":                "2",
					"egress_bandwidth":                 "2",
					"ephemeral_storage":                "2",
					"load_balancer_weight":             "2",
					"tags": map[string]string{
						"name": "tf-test2",
					},
					"acr_registry_infos": []map[string]interface{}{
						{
							"domains":       []string{"test-registry-vpc.cn-hangzhou.cr.aliyuncs.com2"},
							"instance_id":   "cri-47rme9691uiowvfv2",
							"region_id":     "cn-beijing",
							"instance_name": "zzz2",
						},
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "newPassword",
							"server":   "newServer",
							"username": "newUserName",
						},
					},
					"containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds2"},
							"lifecycle_pre_stop_handler_execs":            []string{"echo 2"},
							"security_context_read_only_root_file_system": "false",
							"security_context_run_as_user":                "2",
							"ports": []map[string]interface{}{
								{
									"protocol": "newProtocol",
									"port":     "2",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "newKey",
									"value":                "newValue",
									"field_ref_field_path": "newPath",
								},
							},
							"working_dir":       "newWorkingDir",
							"args":              []string{"arg2"},
							"cpu":               "2",
							"gpu":               "2",
							"memory":            "2",
							"name":              "newName",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.6",
							"image_pull_policy": "newPolicy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "newPath",
									"name":       "newName",
									"read_only":  "false",
								},
							},
							"commands":                              []string{"cmd2"},
							"liveness_probe_exec_commands":          []string{"cmd2"},
							"liveness_probe_period_seconds":         "2",
							"liveness_probe_http_get_path":          "path2",
							"liveness_probe_failure_threshold":      "2",
							"liveness_probe_initial_delay_seconds":  "2",
							"liveness_probe_http_get_port":          "2",
							"liveness_probe_http_get_scheme":        "HTTPS",
							"liveness_probe_tcp_socket_port":        "2",
							"liveness_probe_success_threshold":      "1",
							"liveness_probe_timeout_seconds":        "2",
							"readiness_probe_exec_commands":         []string{"cmd2"},
							"readiness_probe_period_seconds":        "2",
							"readiness_probe_http_get_path":         "path2",
							"readiness_probe_failure_threshold":     "2",
							"readiness_probe_initial_delay_seconds": "2",
							"readiness_probe_http_get_port":         "2",
							"readiness_probe_http_get_scheme":       "HTTPS",
							"readiness_probe_tcp_socket_port":       "2",
							"readiness_probe_success_threshold":     "1",
							"readiness_probe_timeout_seconds":       "2",
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds2"},
							"security_context_read_only_root_file_system": "false",
							"security_context_run_as_user":                "2",
							"ports": []map[string]interface{}{
								{
									"protocol": "newProtocol",
									"port":     "2",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "newKey",
									"value":                "newValue",
									"field_ref_field_path": "newPath",
								},
							},
							"working_dir":       "newWorkingDir",
							"args":              []string{"arg2"},
							"cpu":               "2",
							"gpu":               "2",
							"memory":            "2",
							"name":              "newName",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.6",
							"image_pull_policy": "newPolicy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "newPath",
									"name":       "newName",
									"read_only":  "false",
								},
							},
							"commands": []string{"cmd2"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content2",
									"path":    "path2",
								},
							},
							"disk_volume_disk_id":   "disk_volume_disk_id2",
							"disk_volume_fs_type":   "disk_volume_fs_type2",
							"disk_volume_disk_size": "2",
							"flex_volume_driver":    "flex_volume_driver2",
							"flex_volume_fs_type":   "flex_volume_fs_type2",
							"flex_volume_options":   "flex_volume_options2",
							"nfs_volume_path":       "nfs_volume_path2",
							"nfs_volume_read_only":  "false",
							"nfs_volume_server":     "nfs_volume_server2",
							"name":                  "name2",
							"type":                  "type2",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames2"},
							"ip":        "ip2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_configuration_name":       "newName",
						"active":                           "true",
						"force_delete":                     "true",
						"description":                      "newDesc",
						"security_group_id":                CHECKSET,
						"container_group_name":             "1new-name",
						"restart_policy":                   "newPolicy",
						"cpu":                              "4",
						"memory":                           "8",
						"resource_group_id":                "newGroupId",
						"dns_policy":                       "newDnsPolicy",
						"enable_sls":                       "false",
						"image_snapshot_id":                "imageSnapshotId2",
						"ram_role_name":                    "newRoleName",
						"termination_grace_period_seconds": "120",
						"auto_match_image_cache":           "false",
						"ipv6_address_count":               "2",
						"active_deadline_seconds":          "120",
						"spot_strategy":                    "SpotAsPriceGo",
						"spot_price_limit":                 "1.2",
						"auto_create_eip":                  "false",
						"eip_bandwidth":                    "3",
						"host_name":                        "newHostName",
						"ingress_bandwidth":                "2",
						"egress_bandwidth":                 "2",
						"ephemeral_storage":                "2",
						"load_balancer_weight":             "2",
						"tags.name":                        "tf-test2",
						"image_registry_credentials.#":     "1",
						"containers.#":                     "1",
						"containers.0.security_context_read_only_root_file_system":      "false",
						"containers.0.security_context_run_as_user":                     "2",
						"containers.0.security_context_capability_adds.#":               "1",
						"containers.0.security_context_capability_adds.0":               "adds2",
						"containers.0.lifecycle_pre_stop_handler_execs.#":               "1",
						"containers.0.lifecycle_pre_stop_handler_execs.0":               "echo 2",
						"containers.0.environment_vars.#":                               "1",
						"containers.0.image":                                            "registry-vpc.aliyuncs.com/eci_open/alpine:3.6",
						"containers.0.liveness_probe_exec_commands.0":                   "cmd2",
						"containers.0.liveness_probe_period_seconds":                    "2",
						"containers.0.liveness_probe_http_get_path":                     "path2",
						"containers.0.liveness_probe_failure_threshold":                 "2",
						"containers.0.liveness_probe_initial_delay_seconds":             "2",
						"containers.0.liveness_probe_http_get_port":                     "2",
						"containers.0.liveness_probe_http_get_scheme":                   "HTTPS",
						"containers.0.liveness_probe_tcp_socket_port":                   "2",
						"containers.0.liveness_probe_success_threshold":                 "1",
						"containers.0.liveness_probe_timeout_seconds":                   "2",
						"containers.0.readiness_probe_exec_commands.0":                  "cmd2",
						"containers.0.readiness_probe_period_seconds":                   "2",
						"containers.0.readiness_probe_http_get_path":                    "path2",
						"containers.0.readiness_probe_failure_threshold":                "2",
						"containers.0.readiness_probe_initial_delay_seconds":            "2",
						"containers.0.readiness_probe_http_get_port":                    "2",
						"containers.0.readiness_probe_http_get_scheme":                  "HTTPS",
						"containers.0.readiness_probe_tcp_socket_port":                  "2",
						"containers.0.readiness_probe_success_threshold":                "1",
						"containers.0.readiness_probe_timeout_seconds":                  "2",
						"init_containers.#":                                             "1",
						"init_containers.0.image":                                       "registry-vpc.aliyuncs.com/eci_open/alpine:3.6",
						"init_containers.0.security_context_read_only_root_file_system": "false",
						"init_containers.0.security_context_run_as_user":                "2",
						"init_containers.0.security_context_capability_adds.#":          "1",
						"init_containers.0.security_context_capability_adds.0":          "adds2",
						"init_containers.0.environment_vars.#":                          "1",
						"volumes.#":                                                     "1",
						"host_aliases.#":                                                "1",
						"acr_registry_infos.#":                                          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.1.id}", "${data.alicloud_instance_types.default.instance_types.2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{"${data.alicloud_instance_types.default.instance_types.1.id}", "${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_types.#": "3",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssEciScalingConfiguration_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ess_eci_scaling_configuration.default"
	ra := resourceAttrInit(resourceId, AlicloudEssEciScalingConfigurationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEssEciScalingConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-test-acc-alicloud-eci-container-group%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssEciScalingConfiguration)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                 "${alicloud_ess_scaling_group.default.id}",
					"scaling_configuration_name":       name,
					"description":                      "desc",
					"security_group_id":                "${local.alicloud_security_group_id}",
					"container_group_name":             name,
					"restart_policy":                   "restartPolicy",
					"cpu_options_core":                 "1",
					"cpu_options_threads_per_core":     "2",
					"cpu":                              "2",
					"memory":                           "4",
					"resource_group_id":                "resourceGroupId",
					"dns_policy":                       "dnsPolicy",
					"enable_sls":                       "true",
					"image_snapshot_id":                "imageSnapshotId",
					"ram_role_name":                    "ramRoleName",
					"termination_grace_period_seconds": "60",
					"auto_match_image_cache":           "true",
					"ipv6_address_count":               "1",
					"active_deadline_seconds":          "60",
					"spot_strategy":                    "SpotWithPriceLimit",
					"spot_price_limit":                 "1.1",
					"auto_create_eip":                  "true",
					"eip_bandwidth":                    "1",
					"ephemeral_storage":                "1",
					"load_balancer_weight":             "1",
					"host_name":                        "hostname",
					"ingress_bandwidth":                "1",
					"egress_bandwidth":                 "1",
					"tags": map[string]string{
						"name": "tf-test",
					},
					"acr_registry_infos": []map[string]interface{}{
						{
							"domains":       []string{"test-registry-vpc.cn-hangzhou.cr.aliyuncs.com"},
							"instance_id":   "cri-47rme9691uiowvfv",
							"region_id":     "cn-hangzhou",
							"instance_name": "zzz",
						},
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "password",
							"server":   "server",
							"username": "username",
						},
					},
					"containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"lifecycle_pre_stop_handler_execs":            []string{"echo 1"},
							"security_context_read_only_root_file_system": "true",
							"security_context_run_as_user":                "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "path",
									"name":       "name",
									"read_only":  "true",
								},
							},
							"liveness_probe_exec_commands":          []string{"cmd"},
							"liveness_probe_period_seconds":         "1",
							"liveness_probe_http_get_path":          "path",
							"liveness_probe_failure_threshold":      "1",
							"liveness_probe_initial_delay_seconds":  "1",
							"liveness_probe_http_get_port":          "1",
							"liveness_probe_http_get_scheme":        "HTTP",
							"liveness_probe_tcp_socket_port":        "1",
							"liveness_probe_timeout_seconds":        "1",
							"readiness_probe_exec_commands":         []string{"cmd"},
							"readiness_probe_period_seconds":        "1",
							"readiness_probe_http_get_path":         "path",
							"readiness_probe_failure_threshold":     "1",
							"readiness_probe_initial_delay_seconds": "1",
							"readiness_probe_http_get_port":         "1",
							"readiness_probe_http_get_scheme":       "HTTP",
							"readiness_probe_tcp_socket_port":       "1",
							"readiness_probe_timeout_seconds":       "1",
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"security_context_read_only_root_file_system": "true",
							"security_context_run_as_user":                "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "path",
									"name":       "name",
									"read_only":  "true",
								},
							},
							"commands": []string{"cmd"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content",
									"path":    "path",
								},
							},
							"disk_volume_disk_id":   "disk_volume_disk_id",
							"disk_volume_fs_type":   "disk_volume_fs_type",
							"disk_volume_disk_size": "1",
							"flex_volume_driver":    "flex_volume_driver",
							"flex_volume_fs_type":   "flex_volume_fs_type",
							"flex_volume_options":   "flex_volume_options",
							"nfs_volume_path":       "nfs_volume_path",
							"nfs_volume_read_only":  "true",
							"nfs_volume_server":     "nfs_volume_server",
							"name":                  "name",
							"type":                  "type",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames"},
							"ip":        "ip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":                 CHECKSET,
						"scaling_configuration_name":       name,
						"description":                      "desc",
						"security_group_id":                CHECKSET,
						"container_group_name":             name,
						"restart_policy":                   "restartPolicy",
						"cpu_options_core":                 "1",
						"cpu_options_threads_per_core":     "2",
						"cpu":                              "2",
						"memory":                           "4",
						"resource_group_id":                "resourceGroupId",
						"dns_policy":                       "dnsPolicy",
						"enable_sls":                       "true",
						"image_snapshot_id":                "imageSnapshotId",
						"ram_role_name":                    "ramRoleName",
						"termination_grace_period_seconds": "60",
						"auto_match_image_cache":           "true",
						"ipv6_address_count":               "1",
						"active_deadline_seconds":          "60",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit":                 "1.1",
						"auto_create_eip":                  "true",
						"host_name":                        "hostname",
						"ingress_bandwidth":                "1",
						"egress_bandwidth":                 "1",
						"ephemeral_storage":                "1",
						"load_balancer_weight":             "1",
						"tags.name":                        "tf-test",
						"image_registry_credentials.#":     "1",
						"containers.#":                     "1",
						"containers.0.security_context_read_only_root_file_system":      "true",
						"containers.0.security_context_run_as_user":                     "1",
						"containers.0.security_context_capability_adds.#":               "1",
						"containers.0.security_context_capability_adds.0":               "adds",
						"containers.0.lifecycle_pre_stop_handler_execs.#":               "1",
						"containers.0.lifecycle_pre_stop_handler_execs.0":               "echo 1",
						"containers.0.environment_vars.#":                               "1",
						"init_containers.#":                                             "1",
						"init_containers.0.security_context_run_as_user":                "1",
						"init_containers.0.security_context_read_only_root_file_system": "true",
						"init_containers.0.security_context_capability_adds.#":          "1",
						"init_containers.0.security_context_capability_adds.0":          "adds",
						"volumes.#":            "1",
						"host_aliases.#":       "1",
						"acr_registry_infos.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_configuration_name":       "newName",
					"active":                           "true",
					"force_delete":                     "true",
					"cpu_options_core":                 "2",
					"cpu_options_threads_per_core":     "1",
					"description":                      "newDesc",
					"security_group_id":                "${local.alicloud_security_group_id1}",
					"container_group_name":             "1new-name",
					"restart_policy":                   "newPolicy",
					"cpu":                              "4",
					"memory":                           "8",
					"resource_group_id":                "newGroupId",
					"dns_policy":                       "newDnsPolicy",
					"enable_sls":                       "false",
					"image_snapshot_id":                "imageSnapshotId2",
					"ram_role_name":                    "newRoleName",
					"termination_grace_period_seconds": "120",
					"auto_match_image_cache":           "false",
					"ipv6_address_count":               "2",
					"active_deadline_seconds":          "120",
					"spot_strategy":                    "SpotAsPriceGo",
					"spot_price_limit":                 "1.2",
					"auto_create_eip":                  "false",
					"eip_bandwidth":                    "3",
					"host_name":                        "newHostName",
					"ingress_bandwidth":                "2",
					"egress_bandwidth":                 "2",
					"ephemeral_storage":                "2",
					"load_balancer_weight":             "2",
					"tags": map[string]string{
						"name": "tf-test2",
					},
					"acr_registry_infos": []map[string]interface{}{
						{
							"domains":       []string{"test-registry-vpc.cn-hangzhou.cr.aliyuncs.com2"},
							"instance_id":   "cri-47rme9691uiowvfv2",
							"region_id":     "cn-beijing",
							"instance_name": "zzz2",
						},
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "newPassword",
							"server":   "newServer",
							"username": "newUserName",
						},
					},
					"containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds2"},
							"lifecycle_pre_stop_handler_execs":            []string{"echo 2"},
							"security_context_read_only_root_file_system": "false",
							"security_context_run_as_user":                "2",
							"ports": []map[string]interface{}{
								{
									"protocol": "newProtocol",
									"port":     "2",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "newKey",
									"value":                "newValue",
									"field_ref_field_path": "newPath",
								},
							},
							"working_dir":       "newWorkingDir",
							"args":              []string{"arg2"},
							"cpu":               "2",
							"gpu":               "2",
							"memory":            "2",
							"name":              "newName",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.6",
							"image_pull_policy": "newPolicy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "newPath",
									"name":       "newName",
									"read_only":  "false",
								},
							},
							"commands":                              []string{"cmd2"},
							"liveness_probe_exec_commands":          []string{"cmd2"},
							"liveness_probe_period_seconds":         "2",
							"liveness_probe_http_get_path":          "path2",
							"liveness_probe_failure_threshold":      "2",
							"liveness_probe_initial_delay_seconds":  "2",
							"liveness_probe_http_get_port":          "2",
							"liveness_probe_http_get_scheme":        "HTTPS",
							"liveness_probe_tcp_socket_port":        "2",
							"liveness_probe_success_threshold":      "1",
							"liveness_probe_timeout_seconds":        "2",
							"readiness_probe_exec_commands":         []string{"cmd2"},
							"readiness_probe_period_seconds":        "2",
							"readiness_probe_http_get_path":         "path2",
							"readiness_probe_failure_threshold":     "2",
							"readiness_probe_initial_delay_seconds": "2",
							"readiness_probe_http_get_port":         "2",
							"readiness_probe_http_get_scheme":       "HTTPS",
							"readiness_probe_tcp_socket_port":       "2",
							"readiness_probe_success_threshold":     "1",
							"readiness_probe_timeout_seconds":       "2",
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds2"},
							"security_context_read_only_root_file_system": "false",
							"security_context_run_as_user":                "2",
							"ports": []map[string]interface{}{
								{
									"protocol": "newProtocol",
									"port":     "2",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "newKey",
									"value":                "newValue",
									"field_ref_field_path": "newPath",
								},
							},
							"working_dir":       "newWorkingDir",
							"args":              []string{"arg2"},
							"cpu":               "2",
							"gpu":               "2",
							"memory":            "2",
							"name":              "newName",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.6",
							"image_pull_policy": "newPolicy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "newPath",
									"name":       "newName",
									"read_only":  "false",
								},
							},
							"commands": []string{"cmd2"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content2",
									"path":    "path2",
								},
							},
							"disk_volume_disk_id":   "disk_volume_disk_id2",
							"disk_volume_fs_type":   "disk_volume_fs_type2",
							"disk_volume_disk_size": "2",
							"flex_volume_driver":    "flex_volume_driver2",
							"flex_volume_fs_type":   "flex_volume_fs_type2",
							"flex_volume_options":   "flex_volume_options2",
							"nfs_volume_path":       "nfs_volume_path2",
							"nfs_volume_read_only":  "false",
							"nfs_volume_server":     "nfs_volume_server2",
							"name":                  "name2",
							"type":                  "type2",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames2"},
							"ip":        "ip2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_configuration_name":       "newName",
						"active":                           "true",
						"force_delete":                     "true",
						"cpu_options_core":                 "2",
						"cpu_options_threads_per_core":     "1",
						"description":                      "newDesc",
						"security_group_id":                CHECKSET,
						"container_group_name":             "1new-name",
						"restart_policy":                   "newPolicy",
						"cpu":                              "4",
						"memory":                           "8",
						"resource_group_id":                "newGroupId",
						"dns_policy":                       "newDnsPolicy",
						"enable_sls":                       "false",
						"image_snapshot_id":                "imageSnapshotId2",
						"ram_role_name":                    "newRoleName",
						"termination_grace_period_seconds": "120",
						"auto_match_image_cache":           "false",
						"ipv6_address_count":               "2",
						"active_deadline_seconds":          "120",
						"spot_strategy":                    "SpotAsPriceGo",
						"spot_price_limit":                 "1.2",
						"auto_create_eip":                  "false",
						"eip_bandwidth":                    "3",
						"host_name":                        "newHostName",
						"ingress_bandwidth":                "2",
						"egress_bandwidth":                 "2",
						"ephemeral_storage":                "2",
						"load_balancer_weight":             "2",
						"tags.name":                        "tf-test2",
						"image_registry_credentials.#":     "1",
						"containers.#":                     "1",
						"containers.0.security_context_read_only_root_file_system":      "false",
						"containers.0.security_context_run_as_user":                     "2",
						"containers.0.security_context_capability_adds.#":               "1",
						"containers.0.security_context_capability_adds.0":               "adds2",
						"containers.0.lifecycle_pre_stop_handler_execs.#":               "1",
						"containers.0.lifecycle_pre_stop_handler_execs.0":               "echo 2",
						"containers.0.environment_vars.#":                               "1",
						"containers.0.image":                                            "registry-vpc.aliyuncs.com/eci_open/alpine:3.6",
						"containers.0.liveness_probe_exec_commands.0":                   "cmd2",
						"containers.0.liveness_probe_period_seconds":                    "2",
						"containers.0.liveness_probe_http_get_path":                     "path2",
						"containers.0.liveness_probe_failure_threshold":                 "2",
						"containers.0.liveness_probe_initial_delay_seconds":             "2",
						"containers.0.liveness_probe_http_get_port":                     "2",
						"containers.0.liveness_probe_http_get_scheme":                   "HTTPS",
						"containers.0.liveness_probe_tcp_socket_port":                   "2",
						"containers.0.liveness_probe_success_threshold":                 "1",
						"containers.0.liveness_probe_timeout_seconds":                   "2",
						"containers.0.readiness_probe_exec_commands.0":                  "cmd2",
						"containers.0.readiness_probe_period_seconds":                   "2",
						"containers.0.readiness_probe_http_get_path":                    "path2",
						"containers.0.readiness_probe_failure_threshold":                "2",
						"containers.0.readiness_probe_initial_delay_seconds":            "2",
						"containers.0.readiness_probe_http_get_port":                    "2",
						"containers.0.readiness_probe_http_get_scheme":                  "HTTPS",
						"containers.0.readiness_probe_tcp_socket_port":                  "2",
						"containers.0.readiness_probe_success_threshold":                "1",
						"containers.0.readiness_probe_timeout_seconds":                  "2",
						"init_containers.#":                                             "1",
						"init_containers.0.image":                                       "registry-vpc.aliyuncs.com/eci_open/alpine:3.6",
						"init_containers.0.security_context_read_only_root_file_system": "false",
						"init_containers.0.security_context_run_as_user":                "2",
						"init_containers.0.security_context_capability_adds.#":          "1",
						"init_containers.0.security_context_capability_adds.0":          "adds2",
						"init_containers.0.environment_vars.#":                          "1",
						"volumes.#":                                                     "1",
						"host_aliases.#":                                                "1",
						"acr_registry_infos.#":                                          "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssEciScalingConfiguration_supply(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ess_eci_scaling_configuration.default"
	ra := resourceAttrInit(resourceId, AlicloudEssEciScalingConfigurationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEssEciScalingConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-test-acc-alicloud-eci-container-group%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssEciScalingConfiguration)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                 "${alicloud_ess_scaling_group.default.id}",
					"scaling_configuration_name":       name,
					"description":                      "desc",
					"security_group_id":                "${local.alicloud_security_group_id}",
					"container_group_name":             name,
					"restart_policy":                   "restartPolicy",
					"cost_optimization":                "false",
					"cpu_options_core":                 "1",
					"cpu_options_threads_per_core":     "2",
					"cpu":                              "2",
					"memory":                           "4",
					"resource_group_id":                "resourceGroupId",
					"dns_policy":                       "dnsPolicy",
					"enable_sls":                       "true",
					"image_snapshot_id":                "imageSnapshotId",
					"ram_role_name":                    "ramRoleName",
					"termination_grace_period_seconds": "60",
					"auto_match_image_cache":           "true",
					"ipv6_address_count":               "1",
					"active_deadline_seconds":          "60",
					"spot_strategy":                    "SpotWithPriceLimit",
					"spot_price_limit":                 "1.1",
					"auto_create_eip":                  "true",
					"eip_bandwidth":                    "1",
					"ephemeral_storage":                "1",
					"load_balancer_weight":             "1",
					"host_name":                        "hostname",
					"ingress_bandwidth":                "1",
					"egress_bandwidth":                 "1",
					"tags": map[string]string{
						"name": "tf-test",
					},
					"acr_registry_infos": []map[string]interface{}{
						{
							"domains":       []string{"test-registry-vpc.cn-hangzhou.cr.aliyuncs.com"},
							"instance_id":   "cri-47rme9691uiowvfv",
							"region_id":     "cn-hangzhou",
							"instance_name": "zzz",
						},
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "password",
							"server":   "server",
							"username": "username",
						},
					},
					"dns_config_options": []map[string]interface{}{
						{
							"name":  "test",
							"value": "test",
						},
					},
					"security_context_sysctls": []map[string]interface{}{
						{
							"name":  "kernel.msgmax",
							"value": "65536",
						},
					},
					"containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"lifecycle_pre_stop_handler_execs":            []string{"echo 1"},
							"security_context_read_only_root_file_system": "true",
							"tty":                          "true",
							"stdin":                        "true",
							"security_context_run_as_user": "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "path",
									"name":              "name",
									"read_only":         "true",
									"mount_propagation": "None",
									"sub_path":          "data1/",
								},
							},
							"liveness_probe_exec_commands":          []string{"cmd"},
							"liveness_probe_period_seconds":         "1",
							"liveness_probe_http_get_path":          "path",
							"liveness_probe_failure_threshold":      "1",
							"liveness_probe_initial_delay_seconds":  "1",
							"liveness_probe_http_get_port":          "1",
							"liveness_probe_http_get_scheme":        "HTTP",
							"liveness_probe_tcp_socket_port":        "1",
							"liveness_probe_timeout_seconds":        "1",
							"readiness_probe_exec_commands":         []string{"cmd"},
							"readiness_probe_period_seconds":        "1",
							"readiness_probe_http_get_path":         "path",
							"readiness_probe_failure_threshold":     "1",
							"readiness_probe_initial_delay_seconds": "1",
							"readiness_probe_http_get_port":         "1",
							"readiness_probe_http_get_scheme":       "HTTP",
							"readiness_probe_tcp_socket_port":       "1",
							"readiness_probe_timeout_seconds":       "1",
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"security_context_read_only_root_file_system": "true",
							"security_context_run_as_user":                "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "path",
									"name":              "name",
									"read_only":         "true",
									"mount_propagation": "None",
									"sub_path":          "data1/",
								},
							},
							"commands": []string{"cmd"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content",
									"path":    "path",
									"mode":    "0777",
								},
							},
							"disk_volume_disk_id":             "disk_volume_disk_id",
							"disk_volume_fs_type":             "disk_volume_fs_type",
							"disk_volume_disk_size":           "1",
							"flex_volume_driver":              "flex_volume_driver",
							"flex_volume_fs_type":             "flex_volume_fs_type",
							"flex_volume_options":             "flex_volume_options",
							"nfs_volume_path":                 "nfs_volume_path",
							"nfs_volume_read_only":            "true",
							"nfs_volume_server":               "nfs_volume_server",
							"name":                            "name",
							"type":                            "type",
							"host_path_volume_type":           "Directory",
							"host_path_volume_path":           "/etc/test1",
							"config_file_volume_default_mode": "0777",
							"empty_dir_volume_medium":         "memory",
							"empty_dir_volume_size_limit":     "256 Gi",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames"},
							"ip":        "ip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":                 CHECKSET,
						"scaling_configuration_name":       name,
						"description":                      "desc",
						"security_group_id":                CHECKSET,
						"container_group_name":             name,
						"restart_policy":                   "restartPolicy",
						"cpu_options_core":                 "1",
						"cost_optimization":                "false",
						"cpu_options_threads_per_core":     "2",
						"cpu":                              "2",
						"memory":                           "4",
						"resource_group_id":                "resourceGroupId",
						"dns_policy":                       "dnsPolicy",
						"enable_sls":                       "true",
						"image_snapshot_id":                "imageSnapshotId",
						"ram_role_name":                    "ramRoleName",
						"termination_grace_period_seconds": "60",
						"auto_match_image_cache":           "true",
						"ipv6_address_count":               "1",
						"active_deadline_seconds":          "60",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit":                 "1.1",
						"auto_create_eip":                  "true",
						"host_name":                        "hostname",
						"ingress_bandwidth":                "1",
						"egress_bandwidth":                 "1",
						"ephemeral_storage":                "1",
						"load_balancer_weight":             "1",
						"tags.name":                        "tf-test",
						"image_registry_credentials.#":     "1",
						"dns_config_options.#":             "1",
						"security_context_sysctls.#":       "1",
						"containers.#":                     "1",
						"containers.0.security_context_read_only_root_file_system": "true",
						"containers.0.tty":                                              "true",
						"containers.0.stdin":                                            "true",
						"containers.0.security_context_run_as_user":                     "1",
						"containers.0.security_context_capability_adds.#":               "1",
						"containers.0.security_context_capability_adds.0":               "adds",
						"containers.0.lifecycle_pre_stop_handler_execs.#":               "1",
						"containers.0.lifecycle_pre_stop_handler_execs.0":               "echo 1",
						"containers.0.environment_vars.#":                               "1",
						"containers.0.volume_mounts.#":                                  "1",
						"init_containers.#":                                             "1",
						"init_containers.0.security_context_run_as_user":                "1",
						"init_containers.0.security_context_read_only_root_file_system": "true",
						"init_containers.0.security_context_capability_adds.#":          "1",
						"init_containers.0.security_context_capability_adds.0":          "adds",
						"init_containers.0.volume_mounts.#":                             "1",
						"volumes.#":                                                     "1",
						"host_aliases.#":                                                "1",
						"acr_registry_infos.#":                                          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                 "${alicloud_ess_scaling_group.default.id}",
					"scaling_configuration_name":       name,
					"description":                      "desc",
					"cost_optimization":                "true",
					"instance_family_level":            "EntryLevel",
					"security_group_id":                "${local.alicloud_security_group_id}",
					"container_group_name":             name,
					"restart_policy":                   "restartPolicy",
					"cpu_options_core":                 "1",
					"cpu_options_threads_per_core":     "2",
					"cpu":                              "2",
					"memory":                           "4",
					"resource_group_id":                "resourceGroupId",
					"dns_policy":                       "dnsPolicy",
					"enable_sls":                       "true",
					"image_snapshot_id":                "imageSnapshotId",
					"ram_role_name":                    "ramRoleName",
					"termination_grace_period_seconds": "60",
					"auto_match_image_cache":           "true",
					"ipv6_address_count":               "1",
					"active_deadline_seconds":          "60",
					"spot_strategy":                    "SpotWithPriceLimit",
					"spot_price_limit":                 "1.1",
					"auto_create_eip":                  "true",
					"eip_bandwidth":                    "1",
					"ephemeral_storage":                "1",
					"load_balancer_weight":             "1",
					"host_name":                        "hostname",
					"ingress_bandwidth":                "1",
					"egress_bandwidth":                 "1",
					"tags": map[string]string{
						"name": "tf-test",
					},
					"acr_registry_infos": []map[string]interface{}{
						{
							"domains":       []string{"test-registry-vpc.cn-hangzhou.cr.aliyuncs.com"},
							"instance_id":   "cri-47rme9691uiowvfv",
							"region_id":     "cn-hangzhou",
							"instance_name": "zzz",
						},
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "password",
							"server":   "server",
							"username": "username",
						},
					},
					"dns_config_options": []map[string]interface{}{
						{
							"name":  "test",
							"value": "test",
						},
						{
							"name":  "test1",
							"value": "test1",
						},
					},
					"security_context_sysctls": []map[string]interface{}{
						{
							"name":  "kernel.msgmax",
							"value": "65536",
						},
						{
							"name":  "kernel.msgmin",
							"value": "65535",
						},
					},
					"containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"lifecycle_pre_stop_handler_execs":            []string{"echo 1"},
							"security_context_read_only_root_file_system": "true",
							"tty":                          "false",
							"stdin":                        "false",
							"security_context_run_as_user": "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "path2",
									"name":              "name2",
									"read_only":         "true",
									"mount_propagation": "None",
									"sub_path":          "data3/",
								},
								{
									"mount_path":        "path1",
									"name":              "name1",
									"read_only":         "false",
									"mount_propagation": "HostToCotainer",
									"sub_path":          "data2/",
								},
							},
							"liveness_probe_exec_commands":          []string{"cmd"},
							"liveness_probe_period_seconds":         "1",
							"liveness_probe_http_get_path":          "path",
							"liveness_probe_failure_threshold":      "1",
							"liveness_probe_initial_delay_seconds":  "1",
							"liveness_probe_http_get_port":          "1",
							"liveness_probe_http_get_scheme":        "HTTP",
							"liveness_probe_tcp_socket_port":        "1",
							"liveness_probe_timeout_seconds":        "1",
							"readiness_probe_exec_commands":         []string{"cmd"},
							"readiness_probe_period_seconds":        "1",
							"readiness_probe_http_get_path":         "path",
							"readiness_probe_failure_threshold":     "1",
							"readiness_probe_initial_delay_seconds": "1",
							"readiness_probe_http_get_port":         "1",
							"readiness_probe_http_get_scheme":       "HTTP",
							"readiness_probe_tcp_socket_port":       "1",
							"readiness_probe_timeout_seconds":       "1",
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"security_context_read_only_root_file_system": "true",
							"security_context_run_as_user":                "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "path2",
									"name":              "name2",
									"read_only":         "true",
									"mount_propagation": "None",
									"sub_path":          "data3/",
								},
								{
									"mount_path":        "path1",
									"name":              "name1",
									"read_only":         "false",
									"mount_propagation": "HostToCotainer",
									"sub_path":          "data2/",
								},
							},
							"commands": []string{"cmd"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content",
									"path":    "path",
									"mode":    "0644",
								},
							},
							"disk_volume_disk_id":             "disk_volume_disk_id",
							"disk_volume_fs_type":             "disk_volume_fs_type",
							"disk_volume_disk_size":           "1",
							"flex_volume_driver":              "flex_volume_driver",
							"flex_volume_fs_type":             "flex_volume_fs_type",
							"flex_volume_options":             "flex_volume_options",
							"nfs_volume_path":                 "nfs_volume_path",
							"nfs_volume_read_only":            "true",
							"nfs_volume_server":               "nfs_volume_server",
							"name":                            "name",
							"type":                            "type",
							"host_path_volume_type":           "File",
							"host_path_volume_path":           "/etc/test",
							"config_file_volume_default_mode": "0644",
							"empty_dir_volume_medium":         "",
							"empty_dir_volume_size_limit":     "256 Mi",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames"},
							"ip":        "ip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":                 CHECKSET,
						"scaling_configuration_name":       name,
						"description":                      "desc",
						"security_group_id":                CHECKSET,
						"container_group_name":             name,
						"cost_optimization":                "true",
						"instance_family_level":            "EntryLevel",
						"restart_policy":                   "restartPolicy",
						"cpu_options_core":                 "1",
						"cpu_options_threads_per_core":     "2",
						"cpu":                              "2",
						"memory":                           "4",
						"resource_group_id":                "resourceGroupId",
						"dns_policy":                       "dnsPolicy",
						"enable_sls":                       "true",
						"image_snapshot_id":                "imageSnapshotId",
						"ram_role_name":                    "ramRoleName",
						"termination_grace_period_seconds": "60",
						"auto_match_image_cache":           "true",
						"ipv6_address_count":               "1",
						"active_deadline_seconds":          "60",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit":                 "1.1",
						"auto_create_eip":                  "true",
						"host_name":                        "hostname",
						"ingress_bandwidth":                "1",
						"egress_bandwidth":                 "1",
						"ephemeral_storage":                "1",
						"load_balancer_weight":             "1",
						"tags.name":                        "tf-test",
						"image_registry_credentials.#":     "1",
						"dns_config_options.#":             "2",
						"security_context_sysctls.#":       "2",
						"containers.#":                     "1",
						"containers.0.security_context_read_only_root_file_system": "true",
						"containers.0.tty":                                              "false",
						"containers.0.stdin":                                            "false",
						"containers.0.security_context_run_as_user":                     "1",
						"containers.0.security_context_capability_adds.#":               "1",
						"containers.0.security_context_capability_adds.0":               "adds",
						"containers.0.lifecycle_pre_stop_handler_execs.#":               "1",
						"containers.0.lifecycle_pre_stop_handler_execs.0":               "echo 1",
						"containers.0.environment_vars.#":                               "1",
						"containers.0.volume_mounts.#":                                  "2",
						"init_containers.#":                                             "1",
						"init_containers.0.security_context_run_as_user":                "1",
						"init_containers.0.security_context_read_only_root_file_system": "true",
						"init_containers.0.security_context_capability_adds.#":          "1",
						"init_containers.0.security_context_capability_adds.0":          "adds",
						"init_containers.0.volume_mounts.#":                             "2",
						"volumes.#":                                                     "1",
						"host_aliases.#":                                                "1",
						"acr_registry_infos.#":                                          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                 "${alicloud_ess_scaling_group.default.id}",
					"scaling_configuration_name":       name,
					"description":                      "desc",
					"cost_optimization":                "true",
					"instance_family_level":            "EnterpriseLevel",
					"security_group_id":                "${local.alicloud_security_group_id}",
					"container_group_name":             name,
					"restart_policy":                   "restartPolicy",
					"cpu_options_core":                 "1",
					"cpu_options_threads_per_core":     "2",
					"cpu":                              "2",
					"memory":                           "4",
					"resource_group_id":                "resourceGroupId",
					"dns_policy":                       "dnsPolicy",
					"enable_sls":                       "true",
					"image_snapshot_id":                "imageSnapshotId",
					"ram_role_name":                    "ramRoleName",
					"termination_grace_period_seconds": "60",
					"auto_match_image_cache":           "true",
					"ipv6_address_count":               "1",
					"active_deadline_seconds":          "60",
					"spot_strategy":                    "SpotWithPriceLimit",
					"spot_price_limit":                 "1.1",
					"auto_create_eip":                  "true",
					"eip_bandwidth":                    "1",
					"ephemeral_storage":                "1",
					"load_balancer_weight":             "1",
					"host_name":                        "hostname",
					"ingress_bandwidth":                "1",
					"egress_bandwidth":                 "1",
					"tags": map[string]string{
						"name": "tf-test",
					},
					"acr_registry_infos": []map[string]interface{}{
						{
							"domains":       []string{"test-registry-vpc.cn-hangzhou.cr.aliyuncs.com"},
							"instance_id":   "cri-47rme9691uiowvfv",
							"region_id":     "cn-hangzhou",
							"instance_name": "zzz",
						},
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "password",
							"server":   "server",
							"username": "username",
						},
					},
					"dns_config_options": []map[string]interface{}{
						{
							"name":  "test",
							"value": "test",
						},
						{
							"name":  "test1",
							"value": "test1",
						},
					},
					"security_context_sysctls": []map[string]interface{}{
						{
							"name":  "kernel.msgmax",
							"value": "65536",
						},
						{
							"name":  "kernel.msgmin",
							"value": "65535",
						},
					},
					"containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"lifecycle_pre_stop_handler_execs":            []string{"echo 1"},
							"security_context_read_only_root_file_system": "true",
							"security_context_run_as_user":                "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "path2",
									"name":              "name2",
									"read_only":         "true",
									"mount_propagation": "None",
									"sub_path":          "data3/",
								},
								{
									"mount_path":        "path1",
									"name":              "name1",
									"read_only":         "false",
									"mount_propagation": "HostToCotainer",
									"sub_path":          "data2/",
								},
							},
							"liveness_probe_exec_commands":          []string{"cmd"},
							"liveness_probe_period_seconds":         "1",
							"liveness_probe_http_get_path":          "path",
							"liveness_probe_failure_threshold":      "1",
							"liveness_probe_initial_delay_seconds":  "1",
							"liveness_probe_http_get_port":          "1",
							"liveness_probe_http_get_scheme":        "HTTP",
							"liveness_probe_tcp_socket_port":        "1",
							"liveness_probe_timeout_seconds":        "1",
							"readiness_probe_exec_commands":         []string{"cmd"},
							"readiness_probe_period_seconds":        "1",
							"readiness_probe_http_get_path":         "path",
							"readiness_probe_failure_threshold":     "1",
							"readiness_probe_initial_delay_seconds": "1",
							"readiness_probe_http_get_port":         "1",
							"readiness_probe_http_get_scheme":       "HTTP",
							"readiness_probe_tcp_socket_port":       "1",
							"readiness_probe_timeout_seconds":       "1",
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"security_context_capability_adds":            []string{"adds"},
							"security_context_read_only_root_file_system": "true",
							"security_context_run_as_user":                "1",
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":                  "key",
									"value":                "value",
									"field_ref_field_path": "path",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path":        "path2",
									"name":              "name2",
									"read_only":         "true",
									"mount_propagation": "None",
									"sub_path":          "data3/",
								},
								{
									"mount_path":        "path1",
									"name":              "name1",
									"read_only":         "false",
									"mount_propagation": "HostToCotainer",
									"sub_path":          "data2/",
								},
							},
							"commands": []string{"cmd"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content",
									"path":    "path",
									"mode":    "0644",
								},
							},
							"disk_volume_disk_id":             "disk_volume_disk_id",
							"disk_volume_fs_type":             "disk_volume_fs_type",
							"disk_volume_disk_size":           "1",
							"flex_volume_driver":              "flex_volume_driver",
							"flex_volume_fs_type":             "flex_volume_fs_type",
							"flex_volume_options":             "flex_volume_options",
							"nfs_volume_path":                 "nfs_volume_path",
							"nfs_volume_read_only":            "true",
							"nfs_volume_server":               "nfs_volume_server",
							"name":                            "name",
							"type":                            "type",
							"host_path_volume_type":           "File",
							"host_path_volume_path":           "/etc/test",
							"config_file_volume_default_mode": "0644",
							"empty_dir_volume_medium":         "",
							"empty_dir_volume_size_limit":     "256 Mi",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames"},
							"ip":        "ip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":                 CHECKSET,
						"scaling_configuration_name":       name,
						"description":                      "desc",
						"security_group_id":                CHECKSET,
						"container_group_name":             name,
						"cost_optimization":                "true",
						"instance_family_level":            "EnterpriseLevel",
						"restart_policy":                   "restartPolicy",
						"cpu_options_core":                 "1",
						"cpu_options_threads_per_core":     "2",
						"cpu":                              "2",
						"memory":                           "4",
						"resource_group_id":                "resourceGroupId",
						"dns_policy":                       "dnsPolicy",
						"enable_sls":                       "true",
						"image_snapshot_id":                "imageSnapshotId",
						"ram_role_name":                    "ramRoleName",
						"termination_grace_period_seconds": "60",
						"auto_match_image_cache":           "true",
						"ipv6_address_count":               "1",
						"active_deadline_seconds":          "60",
						"spot_strategy":                    "SpotWithPriceLimit",
						"spot_price_limit":                 "1.1",
						"auto_create_eip":                  "true",
						"host_name":                        "hostname",
						"ingress_bandwidth":                "1",
						"egress_bandwidth":                 "1",
						"ephemeral_storage":                "1",
						"load_balancer_weight":             "1",
						"tags.name":                        "tf-test",
						"image_registry_credentials.#":     "1",
						"dns_config_options.#":             "2",
						"security_context_sysctls.#":       "2",
						"containers.#":                     "1",
						"containers.0.security_context_read_only_root_file_system":      "true",
						"containers.0.security_context_run_as_user":                     "1",
						"containers.0.security_context_capability_adds.#":               "1",
						"containers.0.security_context_capability_adds.0":               "adds",
						"containers.0.lifecycle_pre_stop_handler_execs.#":               "1",
						"containers.0.lifecycle_pre_stop_handler_execs.0":               "echo 1",
						"containers.0.environment_vars.#":                               "1",
						"containers.0.volume_mounts.#":                                  "2",
						"init_containers.#":                                             "1",
						"init_containers.0.security_context_run_as_user":                "1",
						"init_containers.0.security_context_read_only_root_file_system": "true",
						"init_containers.0.security_context_capability_adds.#":          "1",
						"init_containers.0.security_context_capability_adds.0":          "adds",
						"init_containers.0.volume_mounts.#":                             "2",
						"volumes.#":                                                     "1",
						"host_aliases.#":                                                "1",
						"acr_registry_infos.#":                                          "1",
					}),
				),
			},
		},
	})
}

var AlicloudEssEciScalingConfigurationMap = map[string]string{}

func resourceEssEciScalingConfiguration(name string) string {
	return fmt.Sprintf(`
	%s

	variable "name" {
		default = "%s"
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
        alicloud_security_group_id1 = length(data.alicloud_security_groups.default.ids) > 1 ? data.alicloud_security_groups.default.ids.1 : concat(alicloud_security_group.default[*].id, [""])[0]
 
	}
    
	`, EcsInstanceCommonTestCase, name)
}
