package alicloud

import (
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEmrCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEmrClusterCreate,
		Read:   resourceAlicloudEmrClusterRead,
		Update: resourceAlicloudEmrClusterUpdate,
		Delete: resourceAlicloudEmrClusterDelete,

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ValidateFunc: validateAllowedStringValue([]string{string(PrePaid), string(PostPaid)}),
			},
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
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validateAllowedStringValue([]string{HOST_GROUP_TYPE_MASTER,
								HOST_GROUP_TYPE_CORE, HOST_GROUP_TYPE_TASK, HOST_GROUP_TYPE_GATEWAY}),
						},
						"period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
						},
						"charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue([]string{string(PrePaid), string(PostPaid)}),
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
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validateAllowedStringValue([]string{string(DiskCloud), string(DiskCloudEfficiency),
								string(DiskCloudSSD), string(DiskCloudESSD), string(DiskLocalDisk)}),
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
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validateAllowedStringValue([]string{string(DiskCloud), string(DiskCloudEfficiency),
								string(DiskCloudSSD), string(DiskCloudESSD)}),
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
				ValidateFunc: validateAllowedStringValue([]string{DEPOSIT_TYPE_FULLY_MANAGED, DEPOSIT_TYPE_HALF_MANAGED}),
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

	request := emr.CreateCreateClusterV2Request()
	if name, ok := d.GetOk("name"); ok {
		request.Name = name.(string)
	}

	if emrVer, ok := d.GetOk("emr_ver"); ok {
		request.EmrVer = emrVer.(string)
	}

	if clusterType, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = clusterType.(string)
	}

	if zoneId, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = zoneId.(string)
	}

	if securityGroupId, ok := d.GetOk("security_group_id"); ok {
		request.SecurityGroupId = securityGroupId.(string)
	}

	if publicIp, ok := d.GetOk("is_open_public_ip"); ok {
		request.IsOpenPublicIp = requests.NewBoolean(publicIp.(bool))
	}

	if userDefinedRole, ok := d.GetOk("user_defined_emr_ecs_role"); ok {
		request.UserDefinedEmrEcsRole = userDefinedRole.(string)
	}

	if sshEnable, ok := d.GetOk("ssh_enable"); ok {
		request.SshEnable = requests.NewBoolean(sshEnable.(bool))
	}

	if masterPwd, ok := d.GetOk("master_pwd"); ok {
		request.MasterPwd = masterPwd.(string)
	}

	if chargeType, ok := d.GetOk("charge_type"); ok {
		request.ChargeType = chargeType.(string)
	}

	if keyPairName, ok := d.GetOk("key_pair_name"); ok {
		request.KeyPairName = keyPairName.(string)
	}

	if depositType, ok := d.GetOk("deposit_type"); ok {
		request.DepositType = depositType.(string)
	}

	if ha, ok := d.GetOk("high_availability_enable"); ok {
		request.HighAvailabilityEnable = requests.NewBoolean(ha.(bool))
	}

	if eas, ok := d.GetOk("eas_enable"); ok {
		request.EasEnable = requests.NewBoolean(eas.(bool))
	}

	if localMeta, ok := d.GetOk("use_local_metadb"); ok {
		request.UseLocalMetaDb = requests.NewBoolean(localMeta.(bool))
	}

	if relatedClusterId, ok := d.GetOk("related_cluster_id"); ok {
		request.RelatedClusterId = relatedClusterId.(string)
	}

	if optionSoftware, ok := d.GetOk("option_software_list"); ok {
		var softwareList []string
		for _, vv := range optionSoftware.([]interface{}) {
			softwareList = append(softwareList, vv.(string))
		}
		request.OptionSoftWareList = &softwareList
	}

	vpcService := VpcService{client}
	vswitchId := Trim(d.Get("vswitch_id").(string))
	request.NetType = "classic"

	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.NetType = "vpc"

		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if request.ZoneId != vsw.ZoneId {
			return WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}

		request.VpcId = vsw.VpcId
	}

	var hostGroups []emr.CreateClusterV2HostGroup
	if groups, ok := d.GetOk("host_group"); ok {
		for _, group := range groups.(*schema.Set).List() {
			kv := group.(map[string]interface{})
			hostGroup := emr.CreateClusterV2HostGroup{}

			if v, ok := kv["period"]; ok {
				hostGroup.Period = strconv.Itoa(v.(int))
			}

			if v, ok := kv["sys_disk_capacity"]; ok {
				hostGroup.SysDiskCapacity = v.(string)
			}

			if v, ok := kv["disk_capacity"]; ok {
				hostGroup.DiskCapacity = v.(string)
			}

			if v, ok := kv["sys_disk_type"]; ok {
				hostGroup.SysDiskType = strings.ToUpper(v.(string))
			}

			if v, ok := kv["disk_type"]; ok {
				hostGroup.DiskType = strings.ToUpper(v.(string))
			}

			if v, ok := kv["host_group_name"]; ok {
				hostGroup.HostGroupName = v.(string)
			}

			if v, ok := kv["disk_count"]; ok {
				hostGroup.DiskCount = v.(string)
			}

			if v, ok := kv["auto_renew"]; ok {
				if v.(bool) == false {
					hostGroup.AutoRenew = "false"
				} else if v.(bool) == true {
					hostGroup.AutoRenew = "true"
				}
			}

			if v, ok := kv["gpu_driver"]; ok {
				hostGroup.GpuDriver = v.(string)
			}

			if v, ok := kv["node_count"]; ok {
				hostGroup.NodeCount = v.(string)
			}

			if v, ok := kv["instance_type"]; ok {
				hostGroup.InstanceType = v.(string)
			}

			if v, ok := kv["charge_type"]; ok {
				hostGroup.ChargeType = v.(string)
			}

			if v, ok := kv["host_group_type"]; ok {
				hostGroup.HostGroupType = v.(string)
			}

			hostGroups = append(hostGroups, hostGroup)
		}
	}
	request.HostGroup = &hostGroups

	raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.CreateClusterV2(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_emr_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*emr.CreateClusterV2Response)
	d.SetId(response.ClusterId)

	emrService := EmrService{client}

	stateConf := BuildStateConf([]string{"CREATING"}, []string{"IDLE"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, emrService.EmrClusterStateRefreshFunc(d.Id(), []string{"CREATE_FAILED"}))
	stateConf.PollInterval = 5 * time.Second
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

	d.Set("name", object.ClusterInfo.Name)
	d.Set("charge_type", object.ClusterInfo.ChargeType)

	d.Set("high_availability_enable", object.ClusterInfo.HighAvailabilityEnable)
	d.Set("net_type", object.ClusterInfo.NetType)
	d.Set("security_group_id", object.ClusterInfo.SecurityGroupId)
	d.Set("net_type", object.ClusterInfo.NetType)
	d.Set("vpc_id", object.ClusterInfo.VpcId)
	d.Set("vswitch_id", object.ClusterInfo.VSwitchId)
	d.Set("use_local_metadb", object.ClusterInfo.UseLocalMetaDb)
	d.Set("deposit_type", object.ClusterInfo.DepositType)
	d.Set("eas_enable", object.ClusterInfo.EasEnable)
	d.Set("user_defined_emr_ecs_role", object.ClusterInfo.UserDefinedEmrEcsRole)
	d.Set("related_cluster_id", object.ClusterInfo.RelateClusterInfo.ClusterId)

	return nil
}

func resourceAlicloudEmrClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	if d.HasChange("name") {
		request := emr.CreateModifyClusterNameRequest()
		request.Name = d.Get("name").(string)
		request.Id = d.Id()
		raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
			return emrClient.ModifyClusterName(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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

		resizeRequest := emr.CreateResizeClusterV2Request()
		resizeRequest.ClusterId = d.Id()

		var resizeHostGroups []emr.ResizeClusterV2HostGroup

		releaseRequest := emr.CreateReleaseClusterHostGroupRequest()
		releaseRequest.ClusterId = d.Id()

		for k, v1 := range newHostGroup {
			if _, ok := oldHostGroup[k]; ok {
				newNodeCount, _ := strconv.Atoi(v1["node_count"].(string))

				listHostGroupRequest := emr.CreateListClusterHostGroupRequest()
				listHostGroupRequest.ClusterId = d.Id()
				listHostGroupRequest.HostGroupName = k

				raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
					return emrClient.ListClusterHostGroup(listHostGroupRequest)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "alicloud_emr_cluster", listHostGroupRequest.GetActionName(), AlibabaCloudSdkGoERROR)
				}

				resp := raw.(*emr.ListClusterHostGroupResponse)
				hostGroupId := resp.HostGroupList.HostGroup[0].HostGroupId
				oldNodeCount := resp.HostGroupList.HostGroup[0].NodeCount

				// scala up
				if oldNodeCount < newNodeCount {
					count := newNodeCount - oldNodeCount
					resizeHostGroup := emr.ResizeClusterV2HostGroup{}
					resizeHostGroup.ClusterId = d.Id()
					resizeHostGroup.HostGroupId = hostGroupId
					resizeHostGroup.NodeCount = strconv.Itoa(count)
					resizeHostGroup.InstanceType = v1["instance_type"].(string)
					resizeHostGroup.HostGroupType = v1["host_group_type"].(string)

					resizeHostGroups = append(resizeHostGroups, resizeHostGroup)
				} else if oldNodeCount > newNodeCount { //scale down
					releaseRequest.HostGroupId = hostGroupId
					releaseRequest.InstanceIdList = v1["instance_list"].(string)

					_, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
						return emrClient.ReleaseClusterHostGroup(releaseRequest)
					})

					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, "alicloud_emr_cluster", releaseRequest.GetActionName(), AlibabaCloudSdkGoERROR)
					}
				}
			}
		}

		if len(resizeHostGroups) != 0 {
			resizeRequest.HostGroup = &resizeHostGroups

			_, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
				return emrClient.ResizeClusterV2(resizeRequest)
			})

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "alicloud_emr_cluster", resizeRequest.GetActionName(), AlibabaCloudSdkGoERROR)
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
	request := emr.CreateReleaseClusterRequest()

	request.Id = d.Id()
	request.ForceRelease = requests.NewBoolean(true)

	raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.ReleaseCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{"RELEASING"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, emrService.EmrClusterStateRefreshFunc(d.Id(), []string{"RELEASE_FAILED"}))
	stateConf.PollInterval = 5 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return WrapError(emrService.WaitForEmrCluster(d.Id(), Deleted, DefaultTimeoutMedium))
}
