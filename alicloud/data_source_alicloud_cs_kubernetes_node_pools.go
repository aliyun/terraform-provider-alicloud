package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudAckNodepools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudAckNodepoolRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nodepools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cis_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"compensate_with_on_demand": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cpu_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bursting_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kms_key_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"performance_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"encrypted": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_snapshot_policy_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mount_target": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_format": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"provisioned_iops": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_system": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"deployment_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desired_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"install_cloud_monitor": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kubelet_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu_manager_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allowed_unsafe_sysctls": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"topology_manager_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pod_pids_limit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpu_cfs_quota": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"serialize_image_pulls": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kube_api_burst": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_dns": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"system_reserved": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"feature_gates": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"registry_burst": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"read_only_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"registry_pull_qps": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reserved_memory": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"limits": {
													Type:     schema.TypeMap,
													Computed: true,
												},
												"numa_node": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"container_log_monitor_interval": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"container_log_max_workers": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"event_burst": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_gc_high_threshold_percent": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eviction_soft_grace_period": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"image_gc_low_threshold_percent": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"memory_manager_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpu_cfs_quota_period": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eviction_soft": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"event_record_qps": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kube_reserved": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"max_pods": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"container_log_max_files": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eviction_hard": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"tracing": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sampling_rate_per_million": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"endpoint": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"container_log_max_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kube_api_qps": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"login_as_non_root": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"management": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_upgrade_policy": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auto_upgrade_kubelet": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"auto_repair": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"auto_upgrade": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"surge_percentage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"surge": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"auto_vul_fix_policy": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"restart_node": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"vul_level": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"auto_repair_policy": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"restart_node": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"auto_vul_fix": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"max_unavailable": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"multi_az_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_name_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_pool_name": {
							Type:         schema.TypeString,
							ExactlyOneOf: []string{},
							Computed:     true,
						},
						"on_demand_base_capacity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"on_demand_percentage_above_base_capacity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"period_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pre_user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_pool_options": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private_pool_options_match_criteria": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_pool_options_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rds_instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eip_bandwidth": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_bond_eip": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"max_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"eip_internet_charge_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_hardening_os": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"soc_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"spot_instance_pools": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spot_instance_remedy": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"spot_price_limit": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"price_limit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_bursting_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"system_disk_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"system_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_encrypt_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_encrypted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"system_disk_kms_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_performance_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_provisioned_iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"system_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"system_disk_snapshot_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"taints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"effect": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tee_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tee_enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"unschedulable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudAckNodepoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	ClusterId := d.Get("cluster_id")
	action := fmt.Sprintf("/clusters/%s/nodepools", ClusterId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["ClusterId"] = d.Get("cluster_id")

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaGet("CS", "2015-12-15", action, query, nil, nil)

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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resp, _ := jsonpath.Get("$.nodepools[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			nodepool_id, _ := jsonpath.Get("$.nodepool_info.nodepool_id", item)
			if _, ok := idsMap[fmt.Sprint(nodepool_id)]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		kubernetes_configRawObj, _ := jsonpath.Get("$.kubernetes_config", objectRaw)
		kubernetes_configRaw := make(map[string]interface{})
		if kubernetes_configRawObj != nil {
			kubernetes_configRaw = kubernetes_configRawObj.(map[string]interface{})
		}
		mapping["cpu_policy"] = kubernetes_configRaw["cpu_policy"]
		mapping["install_cloud_monitor"] = kubernetes_configRaw["cms_enabled"]
		mapping["node_name_mode"] = kubernetes_configRaw["node_name_mode"]
		mapping["pre_user_data"] = kubernetes_configRaw["pre_user_data"]
		mapping["runtime_name"] = kubernetes_configRaw["runtime"]
		mapping["runtime_version"] = kubernetes_configRaw["runtime_version"]
		mapping["unschedulable"] = kubernetes_configRaw["unschedulable"]
		mapping["user_data"] = kubernetes_configRaw["user_data"]

		nodepool_infoRawObj, _ := jsonpath.Get("$.nodepool_info", objectRaw)
		nodepool_infoRaw := make(map[string]interface{})
		if nodepool_infoRawObj != nil {
			nodepool_infoRaw = nodepool_infoRawObj.(map[string]interface{})
		}
		mapping["node_pool_name"] = nodepool_infoRaw["name"]
		nodepool_id, _ := jsonpath.Get("$.nodepool_info.nodepool_id", objectRaw)
		mapping["node_pool_id"] = nodepool_id
		mapping["resource_group_id"] = nodepool_infoRaw["resource_group_id"]

		scaling_groupRawObj, _ := jsonpath.Get("$.scaling_group", objectRaw)
		scaling_groupRaw := make(map[string]interface{})
		if scaling_groupRawObj != nil {
			scaling_groupRaw = scaling_groupRawObj.(map[string]interface{})
		}
		mapping["auto_renew"] = scaling_groupRaw["auto_renew"]
		mapping["auto_renew_period"] = scaling_groupRaw["auto_renew_period"]
		mapping["cis_enabled"] = scaling_groupRaw["cis_enabled"]
		mapping["compensate_with_on_demand"] = scaling_groupRaw["compensate_with_on_demand"]
		mapping["deployment_set_id"] = scaling_groupRaw["deploymentset_id"]
		if v, ok := scaling_groupRaw["desired_size"].(json.Number); ok {
			mapping["desired_size"] = v.String()
		}

		mapping["image_id"] = scaling_groupRaw["image_id"]
		mapping["image_type"] = scaling_groupRaw["image_type"]
		mapping["instance_charge_type"] = scaling_groupRaw["instance_charge_type"]
		mapping["internet_charge_type"] = scaling_groupRaw["internet_charge_type"]
		mapping["internet_max_bandwidth_out"] = scaling_groupRaw["internet_max_bandwidth_out"]
		mapping["key_name"] = scaling_groupRaw["key_pair"]
		mapping["login_as_non_root"] = scaling_groupRaw["login_as_non_root"]
		mapping["multi_az_policy"] = scaling_groupRaw["multi_az_policy"]
		if v, ok := scaling_groupRaw["on_demand_base_capacity"].(json.Number); ok {
			mapping["on_demand_base_capacity"] = v.String()
		}

		if v, ok := scaling_groupRaw["on_demand_percentage_above_base_capacity"].(json.Number); ok {
			mapping["on_demand_percentage_above_base_capacity"] = v.String()
		}

		mapping["password"] = scaling_groupRaw["login_password"]
		mapping["period"] = scaling_groupRaw["period"]
		mapping["period_unit"] = scaling_groupRaw["period_unit"]
		mapping["platform"] = scaling_groupRaw["platform"]
		mapping["ram_role_name"] = scaling_groupRaw["ram_role_name"]
		mapping["scaling_group_id"] = scaling_groupRaw["scaling_group_id"]
		mapping["scaling_policy"] = scaling_groupRaw["scaling_policy"]
		mapping["security_group_id"] = scaling_groupRaw["security_group_id"]
		mapping["security_hardening_os"] = scaling_groupRaw["security_hardening_os"]
		mapping["soc_enabled"] = scaling_groupRaw["soc_enabled"]
		mapping["spot_instance_pools"] = scaling_groupRaw["spot_instance_pools"]
		mapping["spot_instance_remedy"] = scaling_groupRaw["spot_instance_remedy"]
		mapping["spot_strategy"] = scaling_groupRaw["spot_strategy"]
		mapping["system_disk_bursting_enabled"] = scaling_groupRaw["system_disk_bursting_enabled"]
		mapping["system_disk_category"] = scaling_groupRaw["system_disk_category"]
		mapping["system_disk_encrypt_algorithm"] = scaling_groupRaw["system_disk_encrypt_algorithm"]
		mapping["system_disk_encrypted"] = scaling_groupRaw["system_disk_encrypted"]
		mapping["system_disk_kms_key"] = scaling_groupRaw["system_disk_kms_key_id"]
		mapping["system_disk_performance_level"] = scaling_groupRaw["system_disk_performance_level"]
		mapping["system_disk_provisioned_iops"] = scaling_groupRaw["system_disk_provisioned_iops"]
		mapping["system_disk_size"] = scaling_groupRaw["system_disk_size"]
		mapping["system_disk_snapshot_policy_id"] = scaling_groupRaw["worker_system_disk_snapshot_policy_id"]

		data_disksRaw, _ := jsonpath.Get("$.scaling_group.data_disks", objectRaw)
		dataDisksMaps := make([]map[string]interface{}, 0)
		if data_disksRaw != nil {
			for _, data_disksChildRaw := range data_disksRaw.([]interface{}) {
				dataDisksMap := make(map[string]interface{})
				data_disksChildRaw := data_disksChildRaw.(map[string]interface{})
				if v, ok := data_disksChildRaw["auto_format"].(bool); ok {
					dataDisksMap["auto_format"] = strconv.FormatBool(v)
				}

				dataDisksMap["auto_snapshot_policy_id"] = data_disksChildRaw["auto_snapshot_policy_id"]
				dataDisksMap["bursting_enabled"] = data_disksChildRaw["bursting_enabled"]
				dataDisksMap["category"] = data_disksChildRaw["category"]
				dataDisksMap["device"] = data_disksChildRaw["device"]
				dataDisksMap["encrypted"] = data_disksChildRaw["encrypted"]
				dataDisksMap["file_system"] = data_disksChildRaw["file_system"]
				dataDisksMap["kms_key_id"] = data_disksChildRaw["kms_key_id"]
				dataDisksMap["mount_target"] = data_disksChildRaw["mount_target"]
				dataDisksMap["name"] = data_disksChildRaw["disk_name"]
				dataDisksMap["performance_level"] = data_disksChildRaw["performance_level"]
				dataDisksMap["provisioned_iops"] = data_disksChildRaw["provisioned_iops"]
				dataDisksMap["size"] = data_disksChildRaw["size"]
				dataDisksMap["snapshot_id"] = data_disksChildRaw["snapshot_id"]

				dataDisksMaps = append(dataDisksMaps, dataDisksMap)
			}
		}
		mapping["data_disks"] = dataDisksMaps
		instance_typesRaw, _ := jsonpath.Get("$.scaling_group.instance_types", objectRaw)
		mapping["instance_types"] = instance_typesRaw
		kubeletConfigurationMaps := make([]map[string]interface{}, 0)
		kubeletConfigurationMap := make(map[string]interface{})
		kubelet_configurationRawObj, _ := jsonpath.Get("$.node_config.kubelet_configuration", objectRaw)
		kubelet_configurationRaw := make(map[string]interface{})
		if kubelet_configurationRawObj != nil {
			kubelet_configurationRaw = kubelet_configurationRawObj.(map[string]interface{})
		}
		if len(kubelet_configurationRaw) > 0 {
			if v, ok := kubelet_configurationRaw["containerLogMaxFiles"].(json.Number); ok {
				kubeletConfigurationMap["container_log_max_files"] = v.String()
			}

			kubeletConfigurationMap["container_log_max_size"] = kubelet_configurationRaw["containerLogMaxSize"]
			if v, ok := kubelet_configurationRaw["containerLogMaxWorkers"].(json.Number); ok {
				kubeletConfigurationMap["container_log_max_workers"] = v.String()
			}

			kubeletConfigurationMap["container_log_monitor_interval"] = kubelet_configurationRaw["containerLogMonitorInterval"]
			if v, ok := kubelet_configurationRaw["cpuCFSQuota"].(bool); ok {
				kubeletConfigurationMap["cpu_cfs_quota"] = strconv.FormatBool(v)
			}

			kubeletConfigurationMap["cpu_cfs_quota_period"] = kubelet_configurationRaw["cpuCFSQuotaPeriod"]
			kubeletConfigurationMap["cpu_manager_policy"] = kubelet_configurationRaw["cpuManagerPolicy"]
			if v, ok := kubelet_configurationRaw["eventBurst"].(json.Number); ok {
				kubeletConfigurationMap["event_burst"] = v.String()
			}

			if v, ok := kubelet_configurationRaw["eventRecordQPS"].(json.Number); ok {
				kubeletConfigurationMap["event_record_qps"] = v.String()
			}

			kubeletConfigurationMap["eviction_hard"] = kubelet_configurationRaw["evictionHard"]
			kubeletConfigurationMap["eviction_soft"] = kubelet_configurationRaw["evictionSoft"]
			kubeletConfigurationMap["eviction_soft_grace_period"] = kubelet_configurationRaw["evictionSoftGracePeriod"]
			kubeletConfigurationMap["feature_gates"] = kubelet_configurationRaw["featureGates"]
			if v, ok := kubelet_configurationRaw["imageGCHighThresholdPercent"].(json.Number); ok {
				kubeletConfigurationMap["image_gc_high_threshold_percent"] = v.String()
			}

			if v, ok := kubelet_configurationRaw["imageGCLowThresholdPercent"].(json.Number); ok {
				kubeletConfigurationMap["image_gc_low_threshold_percent"] = v.String()
			}

			if v, ok := kubelet_configurationRaw["kubeAPIBurst"].(json.Number); ok {
				kubeletConfigurationMap["kube_api_burst"] = v.String()
			}

			if v, ok := kubelet_configurationRaw["kubeAPIQPS"].(json.Number); ok {
				kubeletConfigurationMap["kube_api_qps"] = v.String()
			}

			kubeletConfigurationMap["kube_reserved"] = kubelet_configurationRaw["kubeReserved"]
			if v, ok := kubelet_configurationRaw["maxPods"].(json.Number); ok {
				kubeletConfigurationMap["max_pods"] = v.String()
			}

			kubeletConfigurationMap["memory_manager_policy"] = kubelet_configurationRaw["memoryManagerPolicy"]
			if v, ok := kubelet_configurationRaw["podPidsLimit"].(json.Number); ok {
				kubeletConfigurationMap["pod_pids_limit"] = v.String()
			}

			if v, ok := kubelet_configurationRaw["readOnlyPort"].(json.Number); ok {
				kubeletConfigurationMap["read_only_port"] = v.String()
			}

			if v, ok := kubelet_configurationRaw["registryBurst"].(json.Number); ok {
				kubeletConfigurationMap["registry_burst"] = v.String()
			}

			if v, ok := kubelet_configurationRaw["registryPullQPS"].(json.Number); ok {
				kubeletConfigurationMap["registry_pull_qps"] = v.String()
			}

			if v, ok := kubelet_configurationRaw["serializeImagePulls"].(bool); ok {
				kubeletConfigurationMap["serialize_image_pulls"] = strconv.FormatBool(v)
			}

			kubeletConfigurationMap["system_reserved"] = kubelet_configurationRaw["systemReserved"]
			kubeletConfigurationMap["topology_manager_policy"] = kubelet_configurationRaw["topologyManagerPolicy"]

			allowedUnsafeSysctlsRaw, _ := jsonpath.Get("$.node_config.kubelet_configuration.allowedUnsafeSysctls", objectRaw)
			kubeletConfigurationMap["allowed_unsafe_sysctls"] = allowedUnsafeSysctlsRaw
			clusterDNSRaw, _ := jsonpath.Get("$.node_config.kubelet_configuration.clusterDNS", objectRaw)
			kubeletConfigurationMap["cluster_dns"] = clusterDNSRaw
			reservedMemoryRaw, _ := jsonpath.Get("$.node_config.kubelet_configuration.reservedMemory", objectRaw)
			reservedMemoryMaps := make([]map[string]interface{}, 0)
			if reservedMemoryRaw != nil {
				for _, reservedMemoryChildRaw := range reservedMemoryRaw.([]interface{}) {
					reservedMemoryMap := make(map[string]interface{})
					reservedMemoryChildRaw := reservedMemoryChildRaw.(map[string]interface{})
					reservedMemoryMap["limits"] = reservedMemoryChildRaw["limits"]
					reservedMemoryMap["numa_node"] = reservedMemoryChildRaw["numaNode"]

					reservedMemoryMaps = append(reservedMemoryMaps, reservedMemoryMap)
				}
			}
			kubeletConfigurationMap["reserved_memory"] = reservedMemoryMaps
			tracingMaps := make([]map[string]interface{}, 0)
			tracingMap := make(map[string]interface{})
			tracingRawObj, _ := jsonpath.Get("$.node_config.kubelet_configuration.tracing", objectRaw)
			tracingRaw := make(map[string]interface{})
			if tracingRawObj != nil {
				tracingRaw = tracingRawObj.(map[string]interface{})
			}
			if len(tracingRaw) > 0 {
				tracingMap["endpoint"] = tracingRaw["endpoint"]
				if v, ok := tracingRaw["samplingRatePerMillion"].(json.Number); ok {
					tracingMap["sampling_rate_per_million"] = v.String()
				}

				tracingMaps = append(tracingMaps, tracingMap)
			}
			kubeletConfigurationMap["tracing"] = tracingMaps
			kubeletConfigurationMaps = append(kubeletConfigurationMaps, kubeletConfigurationMap)
		}
		mapping["kubelet_configuration"] = kubeletConfigurationMaps
		labelsRaw, _ := jsonpath.Get("$.kubernetes_config.labels", objectRaw)
		labelsMaps := make([]map[string]interface{}, 0)
		if labelsRaw != nil {
			for _, labelsChildRaw := range labelsRaw.([]interface{}) {
				labelsMap := make(map[string]interface{})
				labelsChildRaw := labelsChildRaw.(map[string]interface{})
				labelsMap["key"] = labelsChildRaw["key"]
				labelsMap["value"] = labelsChildRaw["value"]

				labelsMaps = append(labelsMaps, labelsMap)
			}
		}
		mapping["labels"] = labelsMaps
		managementMaps := make([]map[string]interface{}, 0)
		managementMap := make(map[string]interface{})
		managementRaw := make(map[string]interface{})
		if objectRaw["management"] != nil {
			managementRaw = objectRaw["management"].(map[string]interface{})
		}
		if len(managementRaw) > 0 {
			managementMap["auto_repair"] = managementRaw["auto_repair"]
			managementMap["auto_upgrade"] = managementRaw["auto_upgrade"]
			managementMap["auto_vul_fix"] = managementRaw["auto_vul_fix"]
			managementMap["enable"] = managementRaw["enable"]

			upgrade_configRawObj, _ := jsonpath.Get("$.management.upgrade_config", objectRaw)
			upgrade_configRaw := make(map[string]interface{})
			if upgrade_configRawObj != nil {
				upgrade_configRaw = upgrade_configRawObj.(map[string]interface{})
			}
			if len(upgrade_configRaw) > 0 {
				managementMap["max_unavailable"] = upgrade_configRaw["max_unavailable"]
				managementMap["surge"] = upgrade_configRaw["surge"]
				managementMap["surge_percentage"] = upgrade_configRaw["surge_percentage"]
			}
			autoRepairPolicyMaps := make([]map[string]interface{}, 0)
			autoRepairPolicyMap := make(map[string]interface{})
			auto_repair_policyRaw := make(map[string]interface{})
			if managementRaw["auto_repair_policy"] != nil {
				auto_repair_policyRaw = managementRaw["auto_repair_policy"].(map[string]interface{})
			}
			if len(auto_repair_policyRaw) > 0 {
				autoRepairPolicyMap["restart_node"] = auto_repair_policyRaw["restart_node"]

				autoRepairPolicyMaps = append(autoRepairPolicyMaps, autoRepairPolicyMap)
			}
			managementMap["auto_repair_policy"] = autoRepairPolicyMaps
			autoUpgradePolicyMaps := make([]map[string]interface{}, 0)
			autoUpgradePolicyMap := make(map[string]interface{})
			auto_upgrade_policyRaw := make(map[string]interface{})
			if managementRaw["auto_upgrade_policy"] != nil {
				auto_upgrade_policyRaw = managementRaw["auto_upgrade_policy"].(map[string]interface{})
			}
			if len(auto_upgrade_policyRaw) > 0 {
				autoUpgradePolicyMap["auto_upgrade_kubelet"] = auto_upgrade_policyRaw["auto_upgrade_kubelet"]

				autoUpgradePolicyMaps = append(autoUpgradePolicyMaps, autoUpgradePolicyMap)
			}
			managementMap["auto_upgrade_policy"] = autoUpgradePolicyMaps
			autoVulFixPolicyMaps := make([]map[string]interface{}, 0)
			autoVulFixPolicyMap := make(map[string]interface{})
			auto_vul_fix_policyRaw := make(map[string]interface{})
			if managementRaw["auto_vul_fix_policy"] != nil {
				auto_vul_fix_policyRaw = managementRaw["auto_vul_fix_policy"].(map[string]interface{})
			}
			if len(auto_vul_fix_policyRaw) > 0 {
				autoVulFixPolicyMap["restart_node"] = auto_vul_fix_policyRaw["restart_node"]
				autoVulFixPolicyMap["vul_level"] = auto_vul_fix_policyRaw["vul_level"]

				autoVulFixPolicyMaps = append(autoVulFixPolicyMaps, autoVulFixPolicyMap)
			}
			managementMap["auto_vul_fix_policy"] = autoVulFixPolicyMaps
			managementMaps = append(managementMaps, managementMap)
		}
		mapping["management"] = managementMaps
		privatePoolOptionsMaps := make([]map[string]interface{}, 0)
		privatePoolOptionsMap := make(map[string]interface{})
		private_pool_optionsRawObj, _ := jsonpath.Get("$.scaling_group.private_pool_options", objectRaw)
		private_pool_optionsRaw := make(map[string]interface{})
		if private_pool_optionsRawObj != nil {
			private_pool_optionsRaw = private_pool_optionsRawObj.(map[string]interface{})
		}
		if len(private_pool_optionsRaw) > 0 {
			privatePoolOptionsMap["private_pool_options_id"] = private_pool_optionsRaw["id"]
			privatePoolOptionsMap["private_pool_options_match_criteria"] = private_pool_optionsRaw["match_criteria"]

			privatePoolOptionsMaps = append(privatePoolOptionsMaps, privatePoolOptionsMap)
		}
		mapping["private_pool_options"] = privatePoolOptionsMaps
		rds_instancesRaw, _ := jsonpath.Get("$.scaling_group.rds_instances", objectRaw)
		mapping["rds_instances"] = rds_instancesRaw
		scalingConfigMaps := make([]map[string]interface{}, 0)
		scalingConfigMap := make(map[string]interface{})
		auto_scalingRaw := make(map[string]interface{})
		if objectRaw["auto_scaling"] != nil {
			auto_scalingRaw = objectRaw["auto_scaling"].(map[string]interface{})
		}
		if len(auto_scalingRaw) > 0 {
			scalingConfigMap["eip_bandwidth"] = auto_scalingRaw["eip_bandwidth"]
			scalingConfigMap["eip_internet_charge_type"] = auto_scalingRaw["eip_internet_charge_type"]
			scalingConfigMap["enable"] = auto_scalingRaw["enable"]
			scalingConfigMap["is_bond_eip"] = auto_scalingRaw["is_bond_eip"]
			scalingConfigMap["max_size"] = auto_scalingRaw["max_instances"]
			scalingConfigMap["min_size"] = auto_scalingRaw["min_instances"]
			scalingConfigMap["type"] = auto_scalingRaw["type"]

			scalingConfigMaps = append(scalingConfigMaps, scalingConfigMap)
		}
		mapping["scaling_config"] = scalingConfigMaps
		security_group_idsRaw, _ := jsonpath.Get("$.scaling_group.security_group_ids", objectRaw)
		mapping["security_group_ids"] = security_group_idsRaw
		spot_price_limitRaw, _ := jsonpath.Get("$.scaling_group.spot_price_limit", objectRaw)
		spotPriceLimitMaps := make([]map[string]interface{}, 0)
		if spot_price_limitRaw != nil {
			for _, spot_price_limitChildRaw := range spot_price_limitRaw.([]interface{}) {
				spotPriceLimitMap := make(map[string]interface{})
				spot_price_limitChildRaw := spot_price_limitChildRaw.(map[string]interface{})
				spotPriceLimitMap["instance_type"] = spot_price_limitChildRaw["instance_type"]
				spotPriceLimitMap["price_limit"] = spot_price_limitChildRaw["price_limit"]

				spotPriceLimitMaps = append(spotPriceLimitMaps, spotPriceLimitMap)
			}
		}
		mapping["spot_price_limit"] = spotPriceLimitMaps
		system_disk_categoriesRaw, _ := jsonpath.Get("$.scaling_group.system_disk_categories", objectRaw)
		mapping["system_disk_categories"] = system_disk_categoriesRaw
		tagsMaps, _ := jsonpath.Get("$.scaling_group.tags", objectRaw)
		mapping["tags"] = tagsToMap(tagsMaps)
		taintsRaw, _ := jsonpath.Get("$.kubernetes_config.taints", objectRaw)
		taintsMaps := make([]map[string]interface{}, 0)
		if taintsRaw != nil {
			for _, taintsChildRaw := range taintsRaw.([]interface{}) {
				taintsMap := make(map[string]interface{})
				taintsChildRaw := taintsChildRaw.(map[string]interface{})
				taintsMap["effect"] = taintsChildRaw["effect"]
				taintsMap["key"] = taintsChildRaw["key"]
				taintsMap["value"] = taintsChildRaw["value"]

				taintsMaps = append(taintsMaps, taintsMap)
			}
		}
		mapping["taints"] = taintsMaps
		teeConfigMaps := make([]map[string]interface{}, 0)
		teeConfigMap := make(map[string]interface{})
		tee_configRaw := make(map[string]interface{})
		if objectRaw["tee_config"] != nil {
			tee_configRaw = objectRaw["tee_config"].(map[string]interface{})
		}
		if len(tee_configRaw) > 0 {
			teeConfigMap["tee_enable"] = tee_configRaw["tee_enable"]

			teeConfigMaps = append(teeConfigMaps, teeConfigMap)
		}
		mapping["tee_config"] = teeConfigMaps
		vswitch_idsRaw, _ := jsonpath.Get("$.scaling_group.vswitch_ids", objectRaw)
		mapping["vswitch_ids"] = vswitch_idsRaw

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["NodePoolName"])
		s = append(s, mapping)
	}

	d.Set("ids", ids)
	d.Set("names", names)
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("nodepools", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
