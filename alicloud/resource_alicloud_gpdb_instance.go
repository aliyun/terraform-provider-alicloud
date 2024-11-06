package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGpdbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbDbInstanceCreate,
		Read:   resourceAliCloudGpdbDbInstanceRead,
		Update: resourceAliCloudGpdbDbInstanceUpdate,
		Delete: resourceAliCloudGpdbDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_instance_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"StorageElastic", "Serverless", "Classic"}, false),
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Basic", "HighAvailability"}, false),
			},
			"instance_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"2C16G", "4C32G", "16C128G", "2C8G", "4C16G", "8C32G", "16C64G"}, false),
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"resource_management_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"resourceGroup", "resourceQueue"}, false),
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"VPC"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"instance_group_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"master_cu": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{2, 4, 8, 16, 32}),
			},
			"seg_node_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 512),
			},
			"seg_storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"seg_disk_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"pl0", "pl1", "pl2"}, false),
			},
			"create_sample_data": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ssl_enabled": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
			},
			"encryption_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"CloudDisk"}, false),
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vector_configuration_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"serverless_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Manual", "Auto"}, false),
			},
			"prod_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"data_share_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"opened", "closed"}, false),
			},
			"used_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"ip_whitelist": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_ip_list"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_group_attribute": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_group_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old != "" && new != "" && old != new {
									oldParts := strings.Split(old, ",")
									sort.Strings(oldParts)
									newParts := strings.Split(new, ",")
									sort.Strings(newParts)
									return reflect.DeepEqual(newParts, oldParts)
								}
								return false
							},
						},
					},
				},
			},
			"parameters": {
				Type:     schema.TypeSet,
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"default_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_changeable_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"force_restart_instance": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"optional_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"security_ip_list": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"ip_whitelist"},
				Deprecated:    "Field 'security_ip_list' has been deprecated from version 1.187.0. Use 'ip_whitelist' instead.",
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  StringInSlice([]string{"Prepaid", "Postpaid"}, false),
				ConflictsWith: []string{"payment_type"},
				Deprecated:    "Field `instance_charge_type` has been deprecated from version 1.187.0. Use `payment_type` instead.",
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"zone_id"},
				Deprecated:    "Field 'availability_zone' has been deprecated from version 1.187.0. Use 'zone_id' instead.",
			},
			"master_node_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: IntInSlice([]int{1, 2}),
				Deprecated:   "Field `master_node_num` has been deprecated from provider version 1.213.0.",
			},
			"private_ip_address": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field `private_ip_address` has been deprecated from provider version 1.213.0.",
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGpdbDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	var response map[string]interface{}
	action := "CreateDBInstance"
	request := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateDBInstance")
	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["DBInstanceMode"] = d.Get("db_instance_mode")
	request["SecurityIPList"] = LOCAL_HOST_IP

	if v, ok := d.GetOk("db_instance_class"); ok {
		request["DBInstanceClass"] = v
	}

	if v, ok := d.GetOk("db_instance_category"); ok {
		request["DBInstanceCategory"] = v
	}

	if v, ok := d.GetOk("instance_spec"); ok {
		request["InstanceSpec"] = v
	}

	if v, ok := d.GetOk("storage_size"); ok {
		request["StorageSize"] = v
	}

	if v, ok := d.GetOk("instance_network_type"); ok {
		request["InstanceNetworkType"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	} else if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	}

	if v, ok := d.GetOk("instance_group_count"); ok {
		request["DBInstanceGroupCount"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertGpdbDbInstancePaymentTypeRequest(v.(string))
	} else if v, ok := d.GetOk("instance_charge_type"); ok {
		request["PayType"] = v
	}

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("master_cu"); ok {
		request["MasterCU"] = v
	}

	if v, ok := d.GetOk("seg_node_num"); ok {
		request["SegNodeNum"] = v
	}

	if v, ok := d.GetOk("seg_storage_type"); ok {
		request["SegStorageType"] = v
	}

	if v, ok := d.GetOk("seg_disk_performance_level"); ok {
		request["SegDiskPerformanceLevel"] = v
	}

	if v, ok := d.GetOkExists("create_sample_data"); ok {
		request["CreateSampleData"] = v
	}

	if v, ok := d.GetOk("encryption_type"); ok {
		request["EncryptionType"] = v
	}

	if v, ok := d.GetOk("encryption_key"); ok {
		request["EncryptionKey"] = v
	}

	if v, ok := d.GetOk("vector_configuration_status"); ok {
		request["VectorConfigurationStatus"] = v
	}

	if v, ok := d.GetOk("serverless_mode"); ok {
		request["ServerlessMode"] = v
	}

	if v, ok := d.GetOk("prod_type"); ok {
		request["ProdType"] = v
	}

	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["DBInstanceDescription"] = v
	}

	if request["VpcId"] == nil && request["VSwitchId"] != nil {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(request["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if v, ok := request["VPCId"].(string); !ok || v == "" {
			request["VPCId"] = vsw["VpcId"]
		}
	}

	if v, ok := d.GetOk("security_ip_list"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			request["SecurityIpList"] = strings.Join(expandStringList(v.(*schema.Set).List())[:], COMMA_SEPARATED)
		}
	}

	if v, ok := d.GetOk("master_node_num"); ok {
		request["MasterNodeNum"] = v
	}

	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}

	if v, ok := d.GetOkExists("ssl_enabled"); ok {
		request["EnableSSL"] = convertGpdbDbInstanceSSLEnabledRequest(v)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	if v, ok := d.GetOkExists("ssl_enabled"); ok {
		sslEnabledStateConf := BuildStateConf([]string{}, []string{fmt.Sprint(v)}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbService.DBInstanceSSLStateRefreshFunc(d, []string{}))
		if _, err := sslEnabledStateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGpdbDbInstanceUpdate(d, meta)
}

func resourceAliCloudGpdbDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbDbInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_db_instance gpdbService.DescribeGpdbDbInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("db_instance_category", object["DBInstanceCategory"])
	d.Set("db_instance_mode", object["DBInstanceMode"])
	d.Set("instance_network_type", object["InstanceNetworkType"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("availability_zone", object["ZoneId"])
	d.Set("payment_type", convertGpdbDbInstancePaymentTypeResponse(object["PayType"]))
	d.Set("instance_charge_type", object["PayType"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("master_cu", formatInt(object["MasterCU"]))
	d.Set("seg_node_num", formatInt(object["SegNodeNum"]))
	d.Set("seg_storage_type", object["StorageType"])
	d.Set("seg_disk_performance_level", convertGpdbDbInstanceSegDiskPerformanceLevelResponse(object["SegDiskPerformanceLevel"]))
	d.Set("encryption_type", object["EncryptionType"])
	d.Set("encryption_key", object["EncryptionKey"])
	d.Set("vector_configuration_status", object["VectorConfigurationStatus"])
	d.Set("maintain_start_time", object["MaintainStartTime"])
	d.Set("maintain_end_time", object["MaintainEndTime"])
	d.Set("serverless_mode", object["ServerlessMode"])
	d.Set("prod_type", object["ProdType"])
	d.Set("description", object["DBInstanceDescription"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("port", object["Port"])
	d.Set("master_node_num", formatInt(object["MasterNodeNum"]))
	d.Set("status", object["DBInstanceStatus"])

	if v, ok := object["SegmentCounts"]; ok && fmt.Sprint(v) != "0" {
		d.Set("node_num", formatInt(v))
	}

	if v, ok := object["StorageSize"]; ok && fmt.Sprint(v) != "0" {
		d.Set("storage_size", formatInt(v))
	}

	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}

	describeDBInstanceIPArrayListObject, err := gpdbService.DescribeDBInstanceIPArrayList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if iPWhitelistMap, ok := describeDBInstanceIPArrayListObject["Items"].(map[string]interface{}); ok && iPWhitelistMap != nil {
		if dBInstanceIPArrayList, ok := iPWhitelistMap["DBInstanceIPArray"]; ok && dBInstanceIPArrayList != nil {
			iPWhitelistMaps := make([]map[string]interface{}, 0)
			for _, dBInstanceIPArrayListItem := range dBInstanceIPArrayList.([]interface{}) {
				if dBInstanceIPArrayListItemMap, ok := dBInstanceIPArrayListItem.(map[string]interface{}); ok {
					if dBInstanceIPArrayListItem.(map[string]interface{})["DBInstanceIPArrayAttribute"] == "hidden" {
						continue
					}
					dBInstanceIPArrayListMap := map[string]interface{}{}
					dBInstanceIPArrayListMap["ip_group_attribute"] = dBInstanceIPArrayListItemMap["DBInstanceIPArrayAttribute"]
					dBInstanceIPArrayListMap["ip_group_name"] = dBInstanceIPArrayListItemMap["DBInstanceIPArrayName"]
					dBInstanceIPArrayListMap["security_ip_list"] = dBInstanceIPArrayListItemMap["SecurityIPList"]
					iPWhitelistMaps = append(iPWhitelistMaps, dBInstanceIPArrayListMap)
				}
			}
			d.Set("ip_whitelist", iPWhitelistMaps)
		}
	}

	describeDBInstanceSSLObject, err := gpdbService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if v, ok := describeDBInstanceSSLObject["SSLEnabled"]; ok && strconv.FormatBool(v.(bool)) != "" {
		d.Set("ssl_enabled", convertGpdbDbInstanceSSLEnabledResponse(v))
	}

	resourceManagementModeObject, err := gpdbService.DescribeDBResourceManagementMode(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if v, ok := resourceManagementModeObject["ResourceManagementMode"]; ok {
		d.Set("resource_management_mode", v)
	}

	dataShareStatusObject, err := gpdbService.DescribeGpdbDbInstanceDataShareStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if dataShareStatus, ok := dataShareStatusObject["DataShareStatus"]; ok {
		d.Set("data_share_status", dataShareStatus)
	}

	parameterObject, err := gpdbService.DescribeParameters(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if parameterList, ok := parameterObject["Parameters"].([]interface{}); ok && parameterList != nil {
		parameterListMaps := make([]map[string]interface{}, 0)
		for _, parameterListItem := range parameterList {
			if parameterListValueItemMap, ok := parameterListItem.(map[string]interface{}); ok {
				parameterListMap := map[string]interface{}{}
				parameterListMap["name"] = parameterListValueItemMap["ParameterName"]
				parameterListMap["value"] = parameterListValueItemMap["CurrentValue"]
				parameterListMap["default_value"] = parameterListValueItemMap["ParameterValue"]
				parameterListMap["is_changeable_config"] = parameterListValueItemMap["IsChangeableConfig"]
				parameterListMap["force_restart_instance"] = parameterListValueItemMap["ForceRestartInstance"]
				parameterListMap["optional_range"] = parameterListValueItemMap["OptionalRange"]
				parameterListMap["parameter_description"] = parameterListValueItemMap["ParameterDescription"]
				parameterListMaps = append(parameterListMaps, parameterListMap)
			}
		}
		d.Set("parameters", parameterListMaps)
	}

	return nil
}

func resourceAliCloudGpdbDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := make(map[string]interface{})
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)

	d.Partial(true)

	if d.HasChange("tags") {
		if err := gpdbService.SetResourceTags(d, "ALIYUN::GPDB::INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["DBInstanceDescription"] = v
	}
	if update {
		action := "ModifyDBInstanceDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		d.SetPartial("description")
	}

	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["NewResourceGroupId"] = v
	}
	if update {
		action := "ModifyDBInstanceResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		d.SetPartial("resource_group_id")
	}

	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("maintain_end_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_end_time"); ok {
		request["EndTime"] = v
	}
	if !d.IsNewResource() && d.HasChange("maintain_start_time") {
		update = true
	}
	if v, ok := d.GetOk("maintain_start_time"); ok {
		request["StartTime"] = v
	}
	if update {
		action := "ModifyDBInstanceMaintainTime"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		d.SetPartial("maintain_end_time")
		d.SetPartial("maintain_start_time")
	}

	if d.HasChange("ip_whitelist") {
		request = map[string]interface{}{
			"DBInstanceId": d.Id(),
		}
		o, n := d.GetChange("ip_whitelist")
		newGroupKeys := make(map[string]struct{})
		if n != nil {
			for _, iPWhitelist := range n.([]interface{}) {
				iPWhitelistArg := iPWhitelist.(map[string]interface{})
				request["DBInstanceIPArrayAttribute"] = iPWhitelistArg["ip_group_attribute"]
				request["DBInstanceIPArrayName"] = iPWhitelistArg["ip_group_name"]
				request["SecurityIPList"] = iPWhitelistArg["security_ip_list"]
				request["ModifyMode"] = 0
				ipGroupName := "default"
				if fmt.Sprint(iPWhitelistArg["ip_group_name"]) != "" {
					ipGroupName = fmt.Sprint(iPWhitelistArg["ip_group_name"])
				}
				newGroupKeys[ipGroupName] = struct{}{}

				action := "ModifySecurityIps"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
			}
		}
		if o != nil {
			for _, iPWhitelist := range o.([]interface{}) {
				iPWhitelistArg := iPWhitelist.(map[string]interface{})
				request["DBInstanceIPArrayAttribute"] = iPWhitelistArg["ip_group_attribute"]
				request["DBInstanceIPArrayName"] = iPWhitelistArg["ip_group_name"]
				request["SecurityIPList"] = iPWhitelistArg["security_ip_list"]
				request["ModifyMode"] = 2
				ipGroupName := "default"
				if fmt.Sprint(iPWhitelistArg["ip_group_name"]) != "" {
					ipGroupName = fmt.Sprint(iPWhitelistArg["ip_group_name"])
				}
				if _, ok := newGroupKeys[ipGroupName]; ok {
					continue
				}
				action := "ModifySecurityIps"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
			}
		}
		d.SetPartial("ip_whitelist")
	}

	if !d.IsNewResource() && d.HasChange("security_ip_list") {
		request = map[string]interface{}{
			"DBInstanceId": d.Id(),
		}
		if v, ok := d.GetOk("security_ip_list"); ok {
			if len(v.(*schema.Set).List()) > 0 {
				request["SecurityIpList"] = strings.Join(expandStringList(v.(*schema.Set).List())[:], COMMA_SEPARATED)
			}
			action := "ModifySecurityIps"
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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

			d.SetPartial("security_ip_list")
		}
	}

	update = false
	request = map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("seg_node_num") {
		update = true

		if v, ok := d.GetOkExists("seg_node_num"); ok {
			request["UpgradeType"] = 0
			request["SegNodeNum"] = v
		}
	}

	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("seg_node_num")
	}

	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("master_node_num") {
		update = true
		if v, ok := d.GetOk("master_node_num"); ok {
			request["UpgradeType"] = 2
			request["MasterNodeNum"] = v
		}
	}

	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("master_node_num")
	}

	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("instance_spec") {
		update = true
		if v, ok := d.GetOk("instance_spec"); ok {
			request["UpgradeType"] = 1
			request["InstanceSpec"] = v
		}
	}
	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_spec")
	}

	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("storage_size") {
		update = true
		if v, ok := d.GetOk("storage_size"); ok {
			request["UpgradeType"] = 1
			request["StorageSize"] = v
		}
	}

	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing", "InternalError"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("storage_size")
	}

	update = false
	modifySegDiskPerformanceLevelReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Id(),
		"UpgradeType":  3,
	}

	if v, ok := d.GetOk("seg_storage_type"); ok {
		modifySegDiskPerformanceLevelReq["SegStorageType"] = v
	}

	if !d.IsNewResource() && d.HasChange("seg_disk_performance_level") {
		update = true
	}
	if v, ok := d.GetOk("seg_disk_performance_level"); ok {
		modifySegDiskPerformanceLevelReq["SegDiskPerformanceLevel"] = v
	}

	if update {
		action := "UpgradeDBInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, modifySegDiskPerformanceLevelReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifySegDiskPerformanceLevelReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("seg_disk_performance_level")
	}

	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("ssl_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("ssl_enabled"); ok {
		request["SSLEnabled"] = v
	}

	if update {
		action := "ModifyDBInstanceSSL"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		sslEnabledStateConf := BuildStateConf([]string{}, []string{fmt.Sprint(request["SSLEnabled"])}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbService.DBInstanceSSLStateRefreshFunc(d, []string{}))
		if _, err := sslEnabledStateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("ssl_enabled")
	}

	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("vector_configuration_status") {
		update = true
	}
	if v, ok := d.GetOk("vector_configuration_status"); ok {
		request["VectorConfigurationStatus"] = v
	}

	if update {
		action := "ModifyVectorConfiguration"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("vector_configuration_status")
	}

	if d.HasChange("parameters") {
		action := "ModifyParameters"
		request, err = getModifyParametersRequest(d)
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("parameters")
	}

	update = false
	modifyMasterSpec := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("master_cu") {
		update = true
	}
	if v, ok := d.GetOk("master_cu"); ok {
		modifyMasterSpec["MasterCU"] = v
	}

	if update {
		action := "ModifyMasterSpec"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, modifyMasterSpec, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyMasterSpec)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("master_cu")
	}

	update = false
	request = map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	action := "EnableDBResourceGroup"
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_management_mode"); ok && d.HasChange("resource_management_mode") {
		update = true
		if v == "resourceGroup" {
			action = "EnableDBResourceGroup"
		} else {
			action = "DisableDBResourceGroup"
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{"Running", "IDLE"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	dataShareStatus := ""
	setDataShareInstanceReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"InstanceList": "[\"" + d.Id() + "\"]",
	}

	if d.HasChange("data_share_status") {
		update = true
	}
	if v, ok := d.GetOk("data_share_status"); ok {
		setDataShareInstanceReq["OperationType"] = convertGpdbDbInstanceDataShareStatusRequest(v.(string))
		dataShareStatus = v.(string)
	}

	if update {
		action := "SetDataShareInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, setDataShareInstanceReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setDataShareInstanceReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, gpdbService.GpdbDbInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		dataShareStatusStateConf := BuildStateConf([]string{}, []string{dataShareStatus}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbService.GpdbDbInstanceDataShareStatusStateRefreshFunc(d.Id(), []string{}))
		if _, err := dataShareStatusStateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("data_share_status")
	}

	d.Partial(false)

	return resourceAliCloudGpdbDbInstanceRead(d, meta)
}

func resourceAliCloudGpdbDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var response map[string]interface{}
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resourceGpdbDbInstance. Because payment_type = 'Subscription'. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	request["ClientToken"] = buildClientToken("DeleteDBInstance")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", "OperationDenied.DBInstancePayType"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func getModifyParametersRequest(d *schema.ResourceData) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}
	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	config := make(map[string]string)
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
	}
	configJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	request["Parameters"] = string(configJson)
	return request, nil
}

func convertGpdbDbInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "Prepaid"
	case "PayAsYouGo":
		return "Postpaid"
	}

	return source
}

func convertGpdbDbInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Prepaid":
		return "Subscription"
	case "Postpaid":
		return "PayAsYouGo"
	}

	return source
}

func convertGpdbDbInstanceSSLEnabledRequest(source interface{}) interface{} {
	switch source {
	case 0:
		return false
	case 1:
		return true
	}

	return source
}

func convertGpdbDbInstanceSSLEnabledResponse(source interface{}) interface{} {
	switch source {
	case false:
		return 0
	case true:
		return 1
	}

	return source
}

func convertGpdbDbInstanceDataShareStatusRequest(source interface{}) interface{} {
	switch source {
	case "opened":
		return "add"
	case "closed":
		return "remove"
	}

	return source
}

func convertGpdbDbInstanceSegDiskPerformanceLevelResponse(source interface{}) interface{} {
	switch source {
	case "PL0":
		return "pl0"
	case "PL1":
		return "pl1"
	case "PL2":
		return "pl2"
	}

	return source
}
