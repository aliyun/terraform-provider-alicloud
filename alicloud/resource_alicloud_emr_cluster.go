package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEmrCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEmrClusterCreate,
		Read:   resourceAlicloudEmrClusterRead,
		Update: resourceAlicloudEmrClusterUpdate,
		Delete: resourceAlicloudEmrClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HADOOP", "KAFKA", "DRUID", "GATEWAY", "FLINK", "DATA_SCIENCE", "PRESTO", "SECURITY_CENTER", "DSW", "SHUFFLE_SERVICE", "EMR_STUDIO", "KUDU", "ZOOKEEPER"}, false),
			},
			"emr_ver": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PostPaid",
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("charge_type").(string) == "PostPaid"
				},
			},
			"tags": tagsSchema(),
			"host_group": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_group_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"MASTER", "CORE", "TASK", "GATEWAY"}, false),
						},
						"period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
						},
						"charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
						},
						"node_count": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"cloud", "cloud_efficiency", "cloud_ssd", "cloud_essd", "local_disk"}, false),
						},
						"disk_capacity": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk_count": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sys_disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"cloud", "cloud_efficiency", "cloud_ssd", "cloud_essd"}, false),
						},
						"sys_disk_capacity": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"gpu_driver": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_list": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable_graceful_decommission": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"decommission_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"bootstrap_action": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"arg": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"execution_target": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"execution_moment": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"BEFORE_INSTALL", "AFTER_STARTED"}, false),
						},
						"execution_fail_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"FAILED_BLOCKED", "FAILED_CONTINUE"}, false),
						},
					},
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_open_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"option_software_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"high_availability_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"use_local_metadb": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"ssh_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"master_pwd": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"eas_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"user_defined_emr_ecs_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deposit_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"FULLY_MANAGED", "HALF_MANAGED"}, false),
				Default:      "HALF_MANAGED",
			},
			"related_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEmrClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateClusterV2"
	request := make(map[string]interface{})
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}
	request = map[string]interface{}{
		"RegionId": client.RegionId,
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("emr_ver"); ok {
		request["EmrVer"] = v
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		request["ClusterType"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}

	if v, ok := d.GetOkExists("is_open_public_ip"); ok {
		request["IsOpenPublicIp"] = v
	}

	if v, ok := d.GetOk("user_defined_emr_ecs_role"); ok {
		request["UserDefinedEmrEcsRole"] = v
	}

	if v, ok := d.GetOkExists("ssh_enable"); ok {
		request["SshEnable"] = v
	}

	if v, ok := d.GetOk("master_pwd"); ok {
		request["MasterPwd"] = v
	}

	if v, ok := d.GetOk("charge_type"); ok {
		request["ChargeType"] = v
	}

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}

	if v, ok := d.GetOk("deposit_type"); ok {
		request["DepositType"] = v
	}

	if v, ok := d.GetOkExists("high_availability_enable"); ok {
		request["HighAvailabilityEnable"] = v
	}

	if v, ok := d.GetOkExists("eas_enable"); ok {
		request["EasEnable"] = v
	}

	if v, ok := d.GetOkExists("use_local_metadb"); ok {
		request["UseLocalMetaDb"] = v
	}

	if v, ok := d.GetOk("related_cluster_id"); ok {
		request["RelatedClusterId"] = v
	}

	if v, ok := d.GetOk("option_software_list"); ok {
		var softwareList []string
		for _, vv := range v.([]interface{}) {
			softwareList = append(softwareList, vv.(string))
		}
		request["OptionSoftWareList"] = softwareList
	}

	vpcService := VpcService{client}
	vswitchId := Trim(d.Get("vswitch_id").(string))
	request["NetType"] = "classic"

	if vswitchId != "" {
		request["VSwitchId"] = vswitchId
		request["NetType"] = "vpc"

		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return WrapError(err)
		}

		if request["ZoneId"] == "" {
			request["ZoneId"] = vsw.ZoneId
		} else if request["ZoneId"] != vsw.ZoneId {
			return WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request["ZoneId"]))
		}

		request["VpcId"] = vsw.VpcId
	}

	hostGroups := make([]map[string]interface{}, 0)
	if groups, ok := d.GetOk("host_group"); ok {
		nodeChecker := map[string]int{}
		for _, group := range groups.(*schema.Set).List() {
			kv := group.(map[string]interface{})
			hostGroup := map[string]interface{}{}

			if v, ok := kv["period"]; ok {
				hostGroup["Period"] = v
			}

			if v, ok := kv["sys_disk_capacity"]; ok {
				hostGroup["SysDiskCapacity"] = v
			}

			if v, ok := kv["disk_capacity"]; ok {
				hostGroup["DiskCapacity"] = v
			}

			if v, ok := kv["sys_disk_type"]; ok {
				hostGroup["SysDiskType"] = v
			}

			if v, ok := kv["disk_type"]; ok {
				hostGroup["DiskType"] = v
			}

			if v, ok := kv["host_group_name"]; ok {
				hostGroup["HostGroupName"] = v
			}

			if v, ok := kv["disk_count"]; ok {
				hostGroup["DiskCount"] = v
			}

			if v, ok := kv["auto_renew"]; ok {
				if v.(bool) == false {
					hostGroup["AutoRenew"] = "false"
				} else if v.(bool) == true {
					hostGroup["AutoRenew"] = "true"
				}
			}

			if v, ok := kv["gpu_driver"]; ok {
				hostGroup["GpuDriver"] = v
			}

			if v, ok := kv["node_count"]; ok {
				hostGroup["NodeCount"] = v
			}

			if v, ok := kv["instance_type"]; ok {
				hostGroup["InstanceType"] = v
			}

			if v, ok := kv["charge_type"]; ok {
				hostGroup["ChargeType"] = v
			}

			if v, ok := kv["host_group_type"]; ok {
				hostGroup["HostGroupType"] = v
				if nodeCount, exist := kv["node_count"]; exist {
					count, _ := strconv.Atoi(nodeCount.(string))
					nodeChecker[v.(string)] = count
				}
			}

			hostGroups = append(hostGroups, hostGroup)
		}
		// Gateway emr cluster do not need to check
		if request["ClusterType"] != "GATEWAY" {
			if nodeChecker["MASTER"] < 1 || nodeChecker["CORE"] < 2 {
				return WrapError(Error("%s emr cluster must contains 1 MASTER node and 2 CORE nodes.",
					request["ClusterType"]))
			}
			if taskNodeCount, exist := nodeChecker["TASK"]; exist && taskNodeCount < 1 {
				return WrapError(Error("%s emr cluster can not create with 0 Task node, must greater than 0.",
					request["ClusterType"]))
			}
			if ha, ok := d.GetOkExists("high_availability_enable"); ok && ha.(bool) && nodeChecker["MASTER"] < 2 {
				return WrapError(Error("High available %s emr cluster must contains 2 MASTER nodes",
					request["ClusterType"]))
			}
		}
	}
	request["HostGroup"] = hostGroups

	bootstrapActions := make([]map[string]interface{}, 0)
	if actions, ok := d.GetOk("bootstrap_action"); ok {
		for _, action := range actions.(*schema.Set).List() {
			kv := action.(map[string]interface{})
			bootstrapAction := map[string]interface{}{}

			if v, ok := kv["name"]; ok {
				bootstrapAction["Name"] = v
			}

			if v, ok := kv["path"]; ok {
				bootstrapAction["Path"] = v
			}

			if v, ok := kv["arg"]; ok {
				bootstrapAction["Arg"] = v
			}

			if v, ok := kv["execution_target"]; ok {
				bootstrapAction["ExecutionTarget"] = v
			}

			if v, ok := kv["execution_moment"]; ok {
				bootstrapAction["ExecutionMoment"] = v
			}

			if v, ok := kv["execution_fail_strategy"]; ok {
				bootstrapAction["ExecutionFailStrategy"] = v
			}

			bootstrapActions = append(bootstrapActions, bootstrapAction)
		}
	}
	request["BootstrapAction"] = bootstrapActions

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_emr_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ClusterId"]))

	d.Partial(true)
	emrService := EmrService{client}
	if err := emrService.setEmrClusterTags(d); err != nil {
		return WrapError(err)
	}
	d.Partial(false)

	stateConf := BuildStateConf([]string{"CREATING"}, []string{"IDLE"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, emrService.EmrClusterStateRefreshFunc(d.Id(), []string{"CREATE_FAILED"}))
	stateConf.PollInterval = 10 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEmrClusterRead(d, meta)
}

func resourceAlicloudEmrClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	emrService := EmrService{client}

	object, err := emrService.DescribeEmrCluster(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object["Name"])
	d.Set("charge_type", object["ChargeType"])
	d.Set("high_availability_enable", object["HighAvailabilityEnable"])
	d.Set("net_type", object["NetType"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("net_type", object["NetType"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("use_local_metadb", object["LocalMetaDb"])
	d.Set("deposit_type", object["DepositType"])
	d.Set("eas_enable", object["EasEnable"])
	d.Set("user_defined_emr_ecs_role", object["UserDefinedEmrEcsRole"])
	if v, ok := object["RelateClusterInfo"]; ok {
		d.Set("related_cluster_id", v.(map[string]interface{})["ClusterId"])
	}
	d.Set("zone_id", object["ZoneId"])
	if v, ok := object["SoftwareInfo"]; ok {
		d.Set("emr_ver", v.(map[string]interface{})["EmrVer"])
	}

	if v, ok := object["SoftwareInfo"]; ok {
		d.Set("cluster_type", v.(map[string]interface{})["ClusterType"])
	}
	tags, err := emrService.ListTagResources(d.Id(), "cluster")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(tags))

	return nil
}

func resourceAlicloudEmrClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	emrService := EmrService{client}
	d.Partial(true)
	if err := emrService.setEmrClusterTags(d); err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"Id":       d.Id(),
		"RegionId": client.RegionId,
	}
	if d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}
	if update {
		action := "ModifyClusterName"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("name")
	}

	if d.HasChange("host_group") {
		v1, v2 := d.GetChange("host_group")

		oldHostGroup := map[string]map[string]interface{}{}
		for _, v := range v1.(*schema.Set).List() {
			hostGroupName := v.(map[string]interface{})["host_group_name"].(string)
			oldHostGroup[hostGroupName] = v.(map[string]interface{})
		}

		newHostGroup := map[string]map[string]interface{}{}
		for _, v := range v2.(*schema.Set).List() {
			hostGroupName := v.(map[string]interface{})["host_group_name"].(string)
			newHostGroup[hostGroupName] = v.(map[string]interface{})
		}

		resizeRequest := map[string]interface{}{
			"ClusterId": d.Id(),
			"RegionId":  client.RegionId,
		}

		resizeHostGroups := make([]map[string]interface{}, 0)

		releaseRequest := map[string]interface{}{
			"ClusterId": d.Id(),
			"RegionId":  client.RegionId,
		}

		for k, v1 := range newHostGroup {
			if _, ok := oldHostGroup[k]; ok {
				newNodeCount, _ := strconv.Atoi(v1["node_count"].(string))
				listHostGroupRequest := map[string]interface{}{
					"ClusterId":     d.Id(),
					"HostGroupName": k,
					"RegionId":      client.RegionId,
				}
				action := "ListClusterHostGroup"
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, listHostGroupRequest, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, listHostGroupRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				resp, err := jsonpath.Get("$.HostGroupList.HostGroup", response)

				if err != nil {
					return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.HostGroupList.HostGroup", response)
				}

				if len(resp.([]interface{})) == 0 {
					continue
				}

				respHostGroupId := resp.([]interface{})[0].(map[string]interface{})["HostGroupId"]
				var hostGroupId string
				switch val := respHostGroupId.(type) {
				case string:
					hostGroupId = val
				case json.Number:
					hostGroupId = val.String()
				}
				oldNodeCount := formatInt(resp.([]interface{})[0].(map[string]interface{})["NodeCount"])

				// scala up
				if oldNodeCount < newNodeCount {
					count := newNodeCount - oldNodeCount
					resizeHostGroup := map[string]interface{}{}
					resizeHostGroup["ClusterId"] = d.Id()
					resizeHostGroup["HostGroupId"] = hostGroupId
					resizeHostGroup["NodeCount"] = strconv.Itoa(count)
					resizeHostGroup["InstanceType"] = v1["instance_type"]
					resizeHostGroup["HostGroupType"] = v1["host_group_type"]
					resizeHostGroup["HostGroupName"] = k
					resizeHostGroup["ChargeType"] = v1["charge_type"]
					resizeHostGroup["SysDiskType"] = v1["sys_disk_type"]
					resizeHostGroup["SysDiskCapacity"] = v1["sys_disk_capacity"]
					resizeHostGroup["DiskType"] = v1["disk_type"]
					resizeHostGroup["DiskCount"] = v1["disk_count"]
					resizeHostGroup["DiskCapacity"] = v1["disk_capacity"]

					resizeHostGroups = append(resizeHostGroups, resizeHostGroup)
				} else if oldNodeCount > newNodeCount { //scale down
					// EMR cluster can only scale down 'TASK' node group.
					if v1["host_group_type"].(string) != "TASK" {
						return WrapError(Error("EMR cluster can only scale down the node group type of [TASK]."))
					}

					// EMR cluster type of HADOOP support graceful decommission.
					clusterType, ok := d.GetOk("cluster_type")
					if ok && clusterType.(string) == "HADOOP" && v1["enable_graceful_decommission"] != nil {
						egd := v1["enable_graceful_decommission"].(bool)
						if egd {
							releaseRequest["EnableGracefulDecommission"] = egd
							releaseRequest["DecommissionTimeout"] = 3600
							if timeout, ok := v1["decommission_timeout"]; ok && timeout != 0 {
								releaseRequest["DecommissionTimeout"] = timeout.(int)
							}
						}
					}

					releaseRequest["HostGroupId"] = hostGroupId
					releaseRequest["InstanceIdList"] = v1["instance_list"]
					releaseRequest["ReleaseNumber"] = oldNodeCount - newNodeCount

					action := "ReleaseClusterHostGroup"
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					wait := incrementalWait(3*time.Second, 5*time.Second)
					err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, releaseRequest, &runtime)
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						return nil
					})
					addDebug(action, response, resizeRequest)
					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
					}
				}
			} else { // 'Task' HostGroupType may not exist when create emr_cluster
				clusterHostGroupRequest := map[string]interface{}{
					"ClusterId":     d.Id(),
					"HostGroupType": v1["host_group_type"],
					"HostGroupName": k,
					"RegionId":      client.RegionId,
				}

				action := "CreateClusterHostGroup"
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, clusterHostGroupRequest, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, clusterHostGroupRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

				listHostGroupRequest := map[string]interface{}{
					"ClusterId":     d.Id(),
					"HostGroupName": k,
					"RegionId":      client.RegionId,
				}
				action = "ListClusterHostGroup"
				runtime = util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait = incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, listHostGroupRequest, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, listHostGroupRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				resp, err := jsonpath.Get("$.HostGroupList.HostGroup", response)
				if err != nil {
					return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.HostGroupList.HostGroup", response)
				}

				if len(resp.([]interface{})) == 0 {
					continue
				}

				hostGroupId := resp.([]interface{})[0].(map[string]interface{})["HostGroupId"].(string)

				newNodeCount, _ := strconv.Atoi(v1["node_count"].(string))
				if newNodeCount <= 0 {
					return WrapError(Error("emr cluster can not resize with 0 Task node, must greater than 0."))
				}
				resizeHostGroup := map[string]interface{}{}
				resizeHostGroup["ClusterId"] = d.Id()
				resizeHostGroup["HostGroupId"] = hostGroupId
				resizeHostGroup["HostGroupName"] = k
				resizeHostGroup["NodeCount"] = strconv.Itoa(newNodeCount)
				resizeHostGroup["ChargeType"] = v1["charge_type"]
				resizeHostGroup["InstanceType"] = v1["instance_type"]
				resizeHostGroup["HostGroupType"] = v1["host_group_type"]
				resizeHostGroup["SysDiskType"] = v1["sys_disk_type"]
				resizeHostGroup["SysDiskCapacity"] = v1["sys_disk_capacity"]
				resizeHostGroup["DiskType"] = v1["disk_type"]
				resizeHostGroup["DiskCount"] = v1["disk_count"]
				resizeHostGroup["DiskCapacity"] = v1["disk_capacity"]

				resizeHostGroups = append(resizeHostGroups, resizeHostGroup)
			}
		}

		if len(resizeHostGroups) != 0 {
			resizeRequest["HostGroup"] = resizeHostGroups
			action := "ResizeClusterV2"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, resizeRequest, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, resizeRequest)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("host_group")
	}
	d.Partial(false)

	return nil
}

func resourceAlicloudEmrClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	emrService := EmrService{client}
	action := "ReleaseCluster"
	var response map[string]interface{}
	conn, err := client.NewEmrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Id":           d.Id(),
		"RegionId":     client.RegionId,
		"ForceRelease": requests.NewBoolean(true),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-08"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	stateConf := BuildStateConf([]string{"RELEASING"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, emrService.EmrClusterStateRefreshFunc(d.Id(), []string{"RELEASE_FAILED"}))
	stateConf.PollInterval = 5 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return WrapError(emrService.WaitForEmrCluster(d.Id(), Deleted, DefaultTimeoutMedium))
}
