package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEhpcCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEhpcClusterCreate,
		Read:   resourceAlicloudEhpcClusterRead,
		Update: resourceAlicloudEhpcClusterUpdate,
		Delete: resourceAlicloudEhpcClusterDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"nis", "ldap"}, false),
			},
			"additional_volumes": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_queue": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"local_directory": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"location": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"remote_directory": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"roles": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 8,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"volume_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"volume_mount_option": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"volume_mountpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"volume_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"nfs", "smb"}, false),
						},
						"volume_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"application": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"client_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 64),
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"compute_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 99),
			},
			"compute_enable_ht": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"compute_instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"compute_spot_price_limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compute_spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NoSpot", "SpotWithPriceLimit", "SpotAsPriceGo"}, false),
			},
			"deploy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Simple", "Tiny"}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("account_type"); ok && v.(string) == "ldap" {
						return false
					}
					return true
				},
			},
			"ecs_charge_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
			},
			"manager_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: intBetween(1, 2),
			},
			"ehpc_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ha_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_owner_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"marketplace", "others", "self", "system"}, false),
			},
			"input_file_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_compute_ess": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"job_queue": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_pair_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"key_pair_name", "password"},
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"login_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1}),
			},
			"login_instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"os_tag": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
				Optional:     true,
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Year", "Month", "Hour", "Week"}, false),
			},
			"plugin": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("scheduler_type"); ok && v.(string) == "custom" {
						return false
					}
					return true
				},
			},
			"scheduler_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"pbs", "slurm", "opengridscheduler", "deadline"}, false),
			},
			"post_install_script": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"args": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ram_node_types": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"release_instance": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"remote_directory": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"remote_vis_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scc_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_disk_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(40, 500),
			},
			"system_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud"}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_mount_option": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_mountpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"nfs", "smb"}, false),
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"without_agent": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"without_elastic_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEhpcClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateCluster"
	request := make(map[string]interface{})
	conn, err := client.NewEhsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("account_type"); ok {
		request["AccountType"] = v
	}
	additionalVolumesArr := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("additional_volumes"); ok {
		for _, additionalVolumes := range v.(*schema.Set).List() {
			additionalVolumesArg := additionalVolumes.(map[string]interface{})
			mapping := map[string]interface{}{
				"JobQueue":          additionalVolumesArg["job_queue"],
				"VolumeId":          additionalVolumesArg["volume_id"],
				"VolumeMountpoint":  additionalVolumesArg["volume_mountpoint"],
				"VolumeMountOption": additionalVolumesArg["volume_mount_option"],
				"RemoteDirectory":   additionalVolumesArg["remote_directory"],
				"VolumeType":        additionalVolumesArg["volume_type"],
				"LocalDirectory":    additionalVolumesArg["local_directory"],
				"VolumeProtocol":    additionalVolumesArg["volume_protocol"],
				"Location":          additionalVolumesArg["location"],
			}
			rolesArr := make([]map[string]interface{}, 0)
			if v, ok := additionalVolumesArg["roles"]; ok {
				for _, roles := range v.(*schema.Set).List() {
					rolesArg := roles.(map[string]interface{})
					rolesArr = append(rolesArr, map[string]interface{}{
						"Name": rolesArg["name"],
					})
				}
				mapping["Roles"] = rolesArr
			}

			additionalVolumesArr = append(additionalVolumesArr, mapping)
		}
		request["AdditionalVolumes"] = additionalVolumesArr
	}
	applicationArr := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("application"); ok {
		for _, application := range v.(*schema.Set).List() {
			applicationArg := application.(map[string]interface{})
			applicationArr = append(applicationArr, map[string]interface{}{
				"Tag": applicationArg["tag"],
			})
		}
		request["Application"] = applicationArr
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("client_version"); ok {
		request["ClientVersion"] = v
	}
	request["Name"] = d.Get("cluster_name")
	if v, ok := d.GetOk("cluster_version"); ok {
		request["ClusterVersion"] = v
	}
	request["EcsOrder.Compute.Count"] = d.Get("compute_count")
	if v, ok := d.GetOkExists("compute_enable_ht"); ok {
		request["ComputeEnableHt"] = v
	}
	request["EcsOrder.Compute.InstanceType"] = d.Get("compute_instance_type")
	if v, ok := d.GetOk("compute_spot_price_limit"); ok {
		request["ComputeSpotPriceLimit"] = v
	}
	if v, ok := d.GetOk("compute_spot_strategy"); ok {
		request["ComputeSpotStrategy"] = v
	}
	if v, ok := d.GetOk("deploy_mode"); ok {
		request["DeployMode"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	if v, ok := d.GetOk("ecs_charge_type"); ok {
		request["EcsChargeType"] = v
	}
	if v, ok := d.GetOk("manager_count"); ok {
		request["EcsOrder.Manager.Count"] = v
	}
	if v, ok := d.GetOk("ehpc_version"); ok {
		request["EhpcVersion"] = v
	}
	if v, ok := d.GetOkExists("ha_enable"); ok {
		request["HaEnable"] = v
	}
	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}
	if v, ok := d.GetOk("image_owner_alias"); ok {
		request["ImageOwnerAlias"] = v
	}
	if v, ok := d.GetOk("input_file_url"); ok {
		request["InputFileUrl"] = v
	}
	if v, ok := d.GetOkExists("is_compute_ess"); ok {
		request["IsComputeEss"] = v
	}
	if v, ok := d.GetOk("job_queue"); ok {
		request["JobQueue"] = v
	}
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}
	request["EcsOrder.Login.Count"] = d.Get("login_count")
	request["EcsOrder.Login.InstanceType"] = d.Get("login_instance_type")
	request["EcsOrder.Manager.InstanceType"] = d.Get("manager_instance_type")
	request["OsTag"] = d.Get("os_tag")
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("plugin"); ok {
		request["Plugin"] = v
	}
	postInstallScriptArr := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("post_install_script"); ok {
		for _, postInstallScript := range v.(*schema.Set).List() {
			postInstallScriptArg := postInstallScript.(map[string]interface{})
			postInstallScriptArr = append(postInstallScriptArr, map[string]interface{}{
				"Args": postInstallScriptArg["args"],
				"Url":  postInstallScriptArg["url"],
			})
		}
		request["PostInstallScript"] = postInstallScriptArr
	}
	if v, ok := d.GetOk("ram_node_types"); ok {
		request["RamNodeTypes"] = convertArrayToString(v, ",")
	}
	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}
	if v, ok := d.GetOk("remote_directory"); ok {
		request["RemoteDirectory"] = v
	}
	if v, ok := d.GetOkExists("remote_vis_enable"); ok {
		request["RemoteVisEnable"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("scc_cluster_id"); ok {
		request["SccClusterId"] = v
	}
	if v, ok := d.GetOk("scheduler_type"); ok {
		request["SchedulerType"] = v
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOk("security_group_name"); ok {
		request["SecurityGroupName"] = v
	}
	if v, ok := d.GetOk("system_disk_level"); ok {
		request["SystemDiskLevel"] = v
	}
	if v, ok := d.GetOk("system_disk_size"); ok {
		request["SystemDiskSize"] = v
	}
	if v, ok := d.GetOk("system_disk_type"); ok {
		request["SystemDiskType"] = v
	}
	if v, ok := d.GetOk("volume_id"); ok {
		request["VolumeId"] = v
	}
	if v, ok := d.GetOk("volume_mount_option"); ok {
		request["VolumeMountOption"] = v
	}
	if v, ok := d.GetOk("volume_mountpoint"); ok {
		request["VolumeMountpoint"] = v
	}
	if v, ok := d.GetOk("volume_protocol"); ok {
		request["VolumeProtocol"] = v
	}
	if v, ok := d.GetOk("volume_type"); ok {
		request["VolumeType"] = v
	}
	if v, ok := d.GetOkExists("without_agent"); ok {
		request["WithoutAgent"] = v
	}
	if v, ok := d.GetOkExists("without_elastic_ip"); ok {
		request["WithoutElasticIp"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	request["ClientToken"] = buildClientToken("CreateCluster")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ehpc_cluster", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["ClusterId"]))

	ehpcService := EhpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ehpcService.EhpcClusterStateRefreshFunc(d.Id(), []string{"exception"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEhpcClusterRead(d, meta)
}
func resourceAlicloudEhpcClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ehpcService := EhpcService{client}
	object, err := ehpcService.DescribeEhpcCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ehpc_cluster ehpcService.DescribeEhpcCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("account_type", object["AccountType"])
	d.Set("client_version", object["ClientVersion"])
	d.Set("cluster_name", object["Name"])
	d.Set("deploy_mode", object["DeployMode"])
	d.Set("description", object["Description"])
	d.Set("ha_enable", object["HaEnable"])
	d.Set("image_id", object["ImageId"])
	d.Set("image_owner_alias", object["ImageOwnerAlias"])
	if ecsInfo, ok := object["EcsInfo"]; ok {
		if compute, ok := ecsInfo.(map[string]interface{})["Compute"]; ok {
			d.Set("compute_count", compute.(map[string]interface{})["Count"])
			d.Set("compute_instance_type", compute.(map[string]interface{})["InstanceType"])
		}
		if login, ok := ecsInfo.(map[string]interface{})["Login"]; ok {
			d.Set("login_count", login.(map[string]interface{})["Count"])
			d.Set("login_instance_type", login.(map[string]interface{})["InstanceType"])
		}
		if manager, ok := ecsInfo.(map[string]interface{})["Manager"]; ok {
			d.Set("manager_instance_type", manager.(map[string]interface{})["InstanceType"])
			d.Set("manager_count", manager.(map[string]interface{})["Count"])
		}
	}
	postInstallScriptsList := make([]map[string]interface{}, 0)
	if postInstallScripts, ok := object["PostInstallScripts"]; ok {
		if postInstallScriptInfo, ok := postInstallScripts.(map[string]interface{})["PostInstallScriptInfo"]; ok {
			for _, v := range postInstallScriptInfo.([]interface{}) {
				postInstallScriptInfoArg := v.(map[string]interface{})
				postInstallScriptsList = append(postInstallScriptsList, map[string]interface{}{
					"url":  postInstallScriptInfoArg["Url"],
					"args": postInstallScriptInfoArg["Args"],
				})
			}
			d.Set("post_install_script", postInstallScriptsList)
		}
	}
	applicationList := make([]map[string]interface{}, 0)
	if applications, ok := object["Applications"]; ok {
		if applicationInfo, ok := applications.(map[string]interface{})["ApplicationInfo"]; ok {
			for _, v := range applicationInfo.([]interface{}) {
				applicationArg := v.(map[string]interface{})
				if vv, ok := applicationArg["Tag"]; ok {
					applicationList = append(applicationList, map[string]interface{}{
						"tag": vv,
					})
				}
			}
			d.Set("application", applicationList)
		}
	}

	d.Set("os_tag", object["OsTag"])
	d.Set("key_pair_name", object["KeyPairName"])
	d.Set("ecs_charge_type", object["EcsChargeType"])
	d.Set("remote_directory", object["RemoteDirectory"])
	d.Set("scc_cluster_id", object["SccClusterId"])
	d.Set("scheduler_type", object["SchedulerType"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("status", object["Status"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("volume_id", object["VolumeId"])
	d.Set("volume_mountpoint", object["VolumeMountpoint"])
	d.Set("volume_protocol", object["VolumeProtocol"])
	d.Set("volume_type", object["VolumeType"])
	return nil
}
func resourceAlicloudEhpcClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEhsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ClusterId": d.Id(),
	}
	if d.HasChange("cluster_name") {
		update = true
		request["Name"] = d.Get("cluster_name")
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("image_id") {
		update = true
		if v, ok := d.GetOk("image_id"); ok {
			request["ImageId"] = v
		}
	}
	if d.HasChange("image_owner_alias") {
		update = true
		if v, ok := d.GetOk("image_owner_alias"); ok {
			request["ImageOwnerAlias"] = v
		}
	}
	if update {
		action := "ModifyClusterAttributes"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		ehpcService := EhpcService{client}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ehpcService.EhpcClusterStateRefreshFunc(d.Id(), []string{"exception"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudEhpcClusterRead(d, meta)
}
func resourceAlicloudEhpcClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCluster"
	var response map[string]interface{}
	conn, err := client.NewEhsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ClusterId": d.Id(),
	}

	if v, ok := d.GetOkExists("release_instance"); ok {
		request["ReleaseInstance"] = fmt.Sprintf("%v", v.(bool))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"ClusterNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	ehpcService := EhpcService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ehpcService.EhpcClusterStateRefreshFunc(d.Id(), []string{"exception"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
