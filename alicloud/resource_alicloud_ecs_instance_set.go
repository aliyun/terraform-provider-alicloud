package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcsInstanceSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsInstanceSetCreate,
		Read:   resourceAlicloudEcsInstanceSetRead,
		Update: resourceAlicloudEcsInstanceSetUpdate,
		Delete: resourceAlicloudEcsInstanceSetDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hpc_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},
			"security_enhancement_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Deactive"}, false),
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"password_inherit": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"boot_check_os_with_assistant": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					diff := d.Get("instance_charge_type").(string) == "PrePaid"
					if diff {
						return diff
					}
					if old != "" && new != "" && strings.HasPrefix(new, strings.Trim(old, "Z")) {
						diff = true
					}
					return diff
				},
			},
			"data_disks": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(2, 128),
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"disk_category": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd", "cloud"}, false),
						},
						"disk_description": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(2, 256),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"performance_level": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
						},
					},
				},
			},
			"internet_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
				DiffSuppressFunc: ecsInternetDiffSuppressFunc,
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"system_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"system_disk_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"system_disk_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"system_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: ecsSystemDiskPerformanceLevelSuppressFunc,
				ValidateFunc:     validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"system_disk_auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchemaWithIgnore(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
			},
			"spot_price_limit": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("spot_strategy"); ok && v.(string) == "SpotWithPriceLimit" {
						return false
					}
					return true
				},
			},
			"dedicated_host_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"launch_template_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"launch_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"launch_template_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: ecsNotAutoRenewDiffSuppressFunc,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"network_interfaces": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"network_interface_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"primary_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"deployment_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"unique_suffix": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclude_instance_filter": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"InstanceId", "InstanceName"}, false),
						},
						"value": {
							MinItems: 1,
							MaxItems: 100,
							Required: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudEcsInstanceSetCreate(d *schema.ResourceData, meta interface{}) error {

	amount := 1
	if v, ok := d.GetOk("amount"); ok {
		amount = formatInt(v)
	}

	err, instanceIds := buildEcsInstanceSetRunInstanceRequest(d, meta, amount)
	if err != nil {
		return err
	}

	d.Set("instance_ids", instanceIds)
	d.SetId(encodeToBase64String(instanceIds))
	return resourceAlicloudEcsInstanceSetUpdate(d, meta)
}

func buildEcsInstanceSetRunInstanceRequest(d *schema.ResourceData, meta interface{}, amount int) (error, []string) {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "RunInstances"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err), nil
	}
	request["RegionId"] = client.RegionId
	request["Amount"] = amount
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("hpc_cluster_id"); ok {
		request["HpcClusterId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIds := make([]string, 0)
		for _, securityGroupId := range v.(*schema.Set).List() {
			securityGroupIds = append(securityGroupIds, fmt.Sprint(securityGroupId))
		}
		request["SecurityGroupIds"] = securityGroupIds
	}
	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		request["SecurityEnhancementStrategy"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("password_inherit"); ok {
		request["PasswordInherit"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	if v, ok := d.GetOk("auto_release_time"); ok {
		request["AutoReleaseTime"] = v
	}
	if v, ok := d.GetOk("data_disks"); ok {
		diskDataMaps := make([]map[string]interface{}, 0)
		for _, disk := range v.(*schema.Set).List() {
			diskArg := disk.(map[string]interface{})
			diskMap := map[string]interface{}{}
			diskMap["DiskName"] = diskArg["disk_name"]
			diskMap["Size"] = diskArg["disk_size"]
			diskMap["Category"] = diskArg["disk_category"]
			diskMap["Description"] = diskArg["disk_description"]
			diskMap["SnapshotId"] = diskArg["snapshot_id"]
			diskMap["Encrypted"] = diskArg["encrypted"]
			diskMap["KMSKeyId"] = diskArg["kms_key_id"]
			diskMap["AutoSnapshotPolicyId"] = diskArg["auto_snapshot_policy_id"]
			diskMap["PerformanceLevel"] = diskArg["performance_level"]

			diskDataMaps = append(diskDataMaps, diskMap)
		}

		request["DataDisk"] = diskDataMaps
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request["InternetMaxBandwidthOut"] = v
	}
	if v, ok := d.GetOk("system_disk_category"); ok {
		request["SystemDisk.Category"] = v
	}
	if v, ok := d.GetOk("system_disk_description"); ok {
		request["SystemDisk.Description"] = v
	}
	if v, ok := d.GetOk("system_disk_size"); ok {
		request["SystemDisk.Size"] = v
	}
	if v, ok := d.GetOk("system_disk_performance_level"); ok {
		request["SystemDisk.PerformanceLevel"] = v
	}
	if v, ok := d.GetOk("system_disk_auto_snapshot_policy_id"); ok {
		request["SystemDisk.AutoSnapshotPolicyId"] = v
	}
	if v, ok := d.GetOk("system_disk_name"); ok {
		request["SystemDisk.DiskName"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}

	if v, ok := d.GetOk("spot_strategy"); ok {
		request["SpotStrategy"] = v
	}
	if v, ok := d.GetOk("spot_price_limit"); ok {
		request["SpotPriceLimit"] = v
	}
	if v, ok := d.GetOk("dedicated_host_id"); ok {
		request["DedicatedHostId"] = v
	}
	if v, ok := d.GetOk("launch_template_name"); ok {
		request["LaunchTemplateName"] = v
	}
	if v, ok := d.GetOk("launch_template_id"); ok {
		request["LaunchTemplateId"] = v
	}
	if v, ok := d.GetOk("launch_template_version"); ok {
		request["LaunchTemplateVersion"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["InstanceChargeType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("unique_suffix"); ok {
		request["UniqueSuffix"] = v
	}
	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request["DeletionProtection"] = v
	}
	if v, ok := d.GetOk("deployment_set_id"); ok {
		request["DeploymentSetId"] = v
	}
	if v, ok := d.GetOk("network_interfaces"); ok {
		eniMaps := make([]map[string]interface{}, 0)
		for _, eni := range v.(*schema.Set).List() {
			eniArg := eni.(map[string]interface{})
			eniMap := map[string]interface{}{}
			eniMap["SecurityGroupId"] = eniArg["security_group_id"]
			eniMap["VSwitchId"] = eniArg["vswitch_id"]
			eniMap["Description"] = eniArg["description"]
			eniMap["NetworkInterfaceName"] = eniArg["network_interface_name"]
			eniMap["PrimaryIpAddress"] = eniArg["primary_ip_address"]
			eniMaps = append(eniMaps, eniMap)
		}

		request["NetworkInterface"] = eniMaps
	}

	request["ClientToken"] = buildClientToken(action)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "resource_alicloud_ecs_instance_set", action, AlibabaCloudSdkGoERROR), nil
	}

	resp, err := jsonpath.Get("$.InstanceIdSets.InstanceIdSet", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "resource_alicloud_ecs_instance_set", "$.InstanceIdSets.InstanceIdSet", response), nil
	}

	instanceIds := make([]string, 0)
	for _, v := range resp.([]interface{}) {
		instanceIds = append(instanceIds, fmt.Sprint(v))
	}

	ecsService := EcsService{client}
	var instanceHealthCheckFunc resource.StateRefreshFunc
	//Not check OK. GetOk will take "value" in judgement!
	if v, _ := d.GetOk("boot_check_os_with_assistant"); v != nil && v.(bool) == true {
		instanceHealthCheckFunc = ecsService.EcsInstanceSetStateRefreshFunc(encodeToBase64String(instanceIds), []string{"Stopping"})
	} else {
		//Default is false.
		instanceHealthCheckFunc = ecsService.EcsInstanceVmSetStateRefreshFuncWithoutOsCheck(encodeToBase64String(instanceIds), []string{"Stopping"})
	}
	stateConf := BuildStateConf([]string{"Pending", "Starting", "Stopped"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, instanceHealthCheckFunc)

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id()), nil
	}

	return nil, instanceIds
}

func resourceAlicloudEcsInstanceSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	objects, err := ecsService.DescribeEcsInstanceSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_instance ecsService.DescribeInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	amount := d.Get("amount")
	if len(objects) != formatInt(amount) {

		instanceIds := make([]string, 0)
		for _, object := range objects {
			instanceIds = append(instanceIds, fmt.Sprint(object["InstanceId"]))
		}
		d.Set("instance_ids", instanceIds)
		d.SetId(encodeToBase64String(instanceIds))
	}

	instance := objects[0]
	d.Set("resource_group_id", instance["ResourceGroupId"])
	d.Set("hpc_cluster_id", instance["HpcClusterId"])
	d.Set("description", instance["Description"])
	d.Set("image_id", instance["ImageId"])
	d.Set("instance_type", instance["InstanceType"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("instance_name", d.Get("instance_name"))
	d.Set("host_name", d.Get("host_name"))
	d.Set("auto_release_time", instance["AutoReleaseTime"])
	d.Set("internet_charge_type", instance["InternetChargeType"])
	d.Set("internet_max_bandwidth_out", instance["InternetMaxBandwidthOut"])
	d.Set("key_pair_name", instance["KeyPairName"])
	d.Set("spot_strategy", instance["SpotStrategy"])
	d.Set("spot_price_limit", instance["SpotPriceLimit"])
	d.Set("dedicated_host_id", instance["DedicatedHostId"])
	d.Set("instance_charge_type", instance["InstanceChargeType"])
	d.Set("deletion_protection", instance["DeletionProtection"])
	d.Set("deployment_set_id", instance["DeploymentSetId"])

	if v, ok := instance["SecurityGroupIds"].(map[string]interface{}); ok {
		securityGroupIds := make([]string, 0)
		for _, v := range v["SecurityGroupId"].([]interface{}) {
			securityGroupIds = append(securityGroupIds, fmt.Sprint(v))
		}
		d.Set("security_group_ids", securityGroupIds)
	}

	if v, ok := instance["VpcAttributes"].(map[string]interface{}); ok {
		d.Set("vswitch_id", fmt.Sprint(v["VSwitchId"]))
	}

	if v, ok := instance["Tags"]; ok {
		d.Set("tags", tagsToMap(v))
	}

	var disk ecs.Disk
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		disk, err = ecsService.DescribeInstanceSystemDisk(fmt.Sprint(instance["InstanceId"]), fmt.Sprint(instance["ResourceGroupId"]))
		if err != nil {
			if NotFoundError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapError(err)
	}

	d.Set("system_disk_category", disk.Category)
	d.Set("system_disk_size", disk.Size)
	d.Set("system_disk_name", disk.DiskName)
	d.Set("system_disk_description", disk.Description)
	d.Set("system_disk_auto_snapshot_policy_id", disk.AutoSnapshotPolicyId)
	d.Set("system_disk_performance_level", disk.PerformanceLevel)

	return nil
}

func resourceAlicloudEcsInstanceSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	d.Partial(false)

	if d.HasChange("tags") {
		instanceIds := make([]string, 0)
		if err := ecsService.SetInstanceSetResourceTags(d, "instance", instanceIds); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("exclude_instance_filter") {
		curInstanceIds := make([]string, 0)
		for _, v := range d.Get("instance_ids").([]interface{}) {
			curInstanceIds = append(curInstanceIds, fmt.Sprint(v))
		}

		oraw, nraw := d.GetChange("exclude_instance_filter")
		removed := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		created := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()

		objects, err := ecsService.DescribeEcsInstanceSet(d.Id())
		if err != nil {
			return WrapError(err)
		}

		if len(removed) > 0 {
			removeInstanceIds := make([]string, 0)

			removeItem := removed[0].(map[string]interface{})
			excludeType := fmt.Sprint(removeItem["key"])

			if strings.EqualFold(excludeType, "InstanceId") {
				instanceIds := removeItem["value"].([]interface{})
				for _, v := range instanceIds {
					instanceId := fmt.Sprint(v)
					for _, id := range curInstanceIds {
						if strings.EqualFold(instanceId, id) {
							removeInstanceIds = append(removeInstanceIds, instanceId)
							break
						}
					}
				}
			} else if strings.EqualFold(excludeType, "InstanceName") {
				instanceNames := removeItem["value"].([]interface{})
				for _, v := range instanceNames {
					instanceName := fmt.Sprint(v)
					for _, object := range objects {
						if strings.EqualFold(instanceName, fmt.Sprint(object["InstanceName"])) {
							removeInstanceIds = append(removeInstanceIds, fmt.Sprint(object["InstanceId"]))
							break
						}
					}
				}
			}

			if len(removeInstanceIds) > 0 {
				err = buildEcsInstanceSetDeleteInstancesRequest(d, meta, removeInstanceIds)
				if err != nil {
					return WrapError(err)
				}

				tmpInstanceIds := make([]string, 0)
				for _, cur := range curInstanceIds {
					flag := true
					for _, remove := range removeInstanceIds {
						if strings.EqualFold(cur, remove) {
							flag = false
							break
						}
					}

					if flag {
						tmpInstanceIds = append(tmpInstanceIds, cur)
					}
				}
				curInstanceIds = tmpInstanceIds
			}
		}

		if len(created) > 0 {
			addInstanceNames := make([]string, 0)
			instanceNamePrefix := d.Get("instance_name").(string)

			addItem := created[0].(map[string]interface{})
			excludeType := fmt.Sprint(addItem["key"])

			if strings.EqualFold(excludeType, "InstanceId") {
				instanceIds := addItem["value"].([]interface{})
				for _, v := range instanceIds {
					addInstanceNames = append(addInstanceNames, fmt.Sprintf("%s-%s", instanceNamePrefix, fmt.Sprint(v)))
				}
			} else if strings.EqualFold(excludeType, "InstanceName") {
				instanceNames := addItem["value"].([]interface{})
				for _, v := range instanceNames {
					addInstanceNames = append(addInstanceNames, fmt.Sprint(v))
				}
			}

			amount := len(addInstanceNames)
			if amount > 0 {
				err, instanceIds := buildEcsInstanceSetRunInstanceRequest(d, meta, amount)
				if err != nil {
					return err
				}

				curInstanceIds = append(curInstanceIds, instanceIds...)
				for index, name := range addInstanceNames {
					err = modifyEcsInstanceRequest(d, meta, instanceIds[index], name)
					if err != nil {
						return err
					}
				}
			}
		}

		// 更新InstanceIds
		d.Set("instance_ids", curInstanceIds)
		d.SetId(encodeToBase64String(curInstanceIds))
	}

	return resourceAlicloudEcsInstanceSetRead(d, meta)
}

func resourceAlicloudEcsInstanceSetDelete(d *schema.ResourceData, meta interface{}) error {

	instanceIds, err := decodeFromBase64String(d.Id())
	if err != nil {
		return WrapError(err)
	}

	err = buildEcsInstanceSetDeleteInstancesRequest(d, meta, instanceIds)
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func buildEcsInstanceSetDeleteInstancesRequest(d *schema.ResourceData, meta interface{}, instanceIds []string) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "DeleteInstances"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	for i, instanceId := range instanceIds {
		request[fmt.Sprintf("InstanceId.%d", i+1)] = instanceId
	}

	request["Force"] = true
	request["RegionId"] = client.RegionId
	request["ClientToken"] = fmt.Sprint(strings.Trim(uuid.New().String(), "-"))[1:30]
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "DependencyViolation.RouteEntry", "IncorrectInstanceStatus.Initializing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceIds.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "resource_alicloud_ecs_instance_set", action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func modifyEcsInstanceRequest(d *schema.ResourceData, meta interface{}, instanceId, instanceName string) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyInstanceAttribute"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	request["InstanceId"] = instanceId
	request["InstanceName"] = instanceName
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapError(err)
	}

	return nil
}
