// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudGpdbDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbDBInstanceCreate,
		Read:   resourceAliCloudGpdbDBInstanceRead,
		Update: resourceAliCloudGpdbDBInstanceUpdate,
		Delete: resourceAliCloudGpdbDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"connection_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"HighAvailability", "Basic"}, false),
			},
			"db_instance_ip_array_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_ip_array_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_share_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"StorageElastic", "Serverless"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encryption_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"6.0", "7.0"}, false),
			},
			"idle_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"instance_network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_spec": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"master_cu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"modify_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"prod_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_management_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sample_data_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_ip_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"seg_disk_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"seg_node_num": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Calculate the number of nodes. Valid values:-The value range of the high-availability version of the storage elastic mode is 4 to 512, and the value must be a multiple of 4.-The value range of the basic version of the storage elastic mode is 2 to 512, and the value must be a multiple of 2.The-Serverless version has a value range of 2 to 512. The value must be a multiple of 2.>-this parameter must be passed in to create a storage elastic mode instance and a Serverless version instance.-During the public beta of the Serverless version (from 0101, 2022 to 0131, 2022), a maximum of 12 compute nodes can be created."),
			},
			"seg_storage_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"serverless_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Manual", "Auto"}, false),
			},
			"serverless_resource": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "Calculate resource thresholds. The value range is 8 to 32, the step size is 8, the unit is ACU. The default value is 32.> This parameter is required for only Serverless automatic scheduling mode instances."),
			},
			"ssl_enabled": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 2}),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 4000),
			},
			"tags": tagsSchema(),
			"upgrade_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"used_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vector_configuration_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudGpdbDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertGpdbDBInstancePayTypeRequest(v.(string))
	}
	request["EngineVersion"] = d.Get("engine_version")
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	request["Engine"] = "gpdb"
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}
	if v, ok := d.GetOk("instance_network_type"); ok {
		request["InstanceNetworkType"] = v
	}
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("seg_node_num"); ok {
		request["SegNodeNum"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("instance_spec"); ok {
		request["InstanceSpec"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("serverless_resource"); ok {
		request["ServerlessResource"] = v
	}
	if v, ok := d.GetOk("serverless_mode"); ok {
		request["ServerlessMode"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["DBInstanceDescription"] = v
	}
	if v, ok := d.GetOk("seg_disk_performance_level"); ok {
		request["SegDiskPerformanceLevel"] = v
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
	if v, ok := d.GetOk("seg_storage_type"); ok {
		request["SegStorageType"] = v
	}
	if v, ok := d.GetOk("db_instance_category"); ok {
		request["DBInstanceCategory"] = v
	}
	if v, ok := d.GetOk("master_cu"); ok {
		request["MasterCU"] = v
	}
	request["MasterNodeNum"] = "1"
	request["DBInstanceMode"] = d.Get("db_instance_mode")
	if v, ok := d.GetOk("ssl_enabled"); ok {
		request["EnableSSL"] = convertGpdbDBInstanceEnableSSLRequest(v.(int))
	}
	if v, ok := d.GetOk("idle_time"); ok {
		request["IdleTime"] = v
	}
	if v, ok := d.GetOk("storage_size"); ok {
		request["StorageSize"] = v
	}
	if v, ok := d.GetOk("sample_data_status"); ok {
		request["CreateSampleData"] = convertGpdbDBInstanceCreateSampleDataRequest(v.(string))
	}
	if v, ok := d.GetOk("prod_type"); ok {
		request["ProdType"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running", "IDLE"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGpdbDBInstanceUpdate(d, meta)
}

func resourceAliCloudGpdbDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbDBInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_instance DescribeGpdbDBInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["DBInstanceCategory"] != nil {
		d.Set("db_instance_category", objectRaw["DBInstanceCategory"])
	}
	if objectRaw["DBInstanceMode"] != nil {
		d.Set("db_instance_mode", objectRaw["DBInstanceMode"])
	}
	if objectRaw["DBInstanceDescription"] != nil {
		d.Set("description", objectRaw["DBInstanceDescription"])
	}
	if objectRaw["EncryptionKey"] != nil {
		d.Set("encryption_key", objectRaw["EncryptionKey"])
	}
	if objectRaw["EncryptionType"] != nil {
		d.Set("encryption_type", objectRaw["EncryptionType"])
	}
	if objectRaw["EngineVersion"] != nil {
		d.Set("engine_version", objectRaw["EngineVersion"])
	}
	if objectRaw["IdleTime"] != nil {
		d.Set("idle_time", objectRaw["IdleTime"])
	}
	if objectRaw["InstanceNetworkType"] != nil {
		d.Set("instance_network_type", objectRaw["InstanceNetworkType"])
	}
	if objectRaw["MaintainEndTime"] != nil {
		d.Set("maintain_end_time", objectRaw["MaintainEndTime"])
	}
	if objectRaw["MaintainStartTime"] != nil {
		d.Set("maintain_start_time", objectRaw["MaintainStartTime"])
	}
	if objectRaw["MasterCU"] != nil {
		d.Set("master_cu", objectRaw["MasterCU"])
	}
	if objectRaw["PayType"] != nil {
		d.Set("payment_type", convertGpdbDBInstanceItemsDBInstanceAttributePayTypeResponse(objectRaw["PayType"]))
	}
	if objectRaw["ProdType"] != nil {
		d.Set("prod_type", objectRaw["ProdType"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["SegDiskPerformanceLevel"] != nil {
		d.Set("seg_disk_performance_level", convertGpdbDBInstanceItemsDBInstanceAttributeSegDiskPerformanceLevelResponse(objectRaw["SegDiskPerformanceLevel"]))
	}
	if objectRaw["SegNodeNum"] != nil {
		d.Set("seg_node_num", objectRaw["SegNodeNum"])
	}
	if objectRaw["ServerlessMode"] != nil {
		d.Set("serverless_mode", objectRaw["ServerlessMode"])
	}
	if objectRaw["ServerlessResource"] != nil {
		d.Set("serverless_resource", objectRaw["ServerlessResource"])
	}
	if objectRaw["DBInstanceStatus"] != nil {
		d.Set("status", objectRaw["DBInstanceStatus"])
	}
	if objectRaw["StorageSize"] != nil {
		d.Set("storage_size", objectRaw["StorageSize"])
	}
	if objectRaw["VSwitchId"] != nil {
		d.Set("vswitch_id", objectRaw["VSwitchId"])
	}
	if objectRaw["VectorConfigurationStatus"] != nil {
		d.Set("vector_configuration_status", objectRaw["VectorConfigurationStatus"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}
	if objectRaw["ZoneId"] != nil {
		d.Set("zone_id", objectRaw["ZoneId"])
	}
	if objectRaw["StorageType"] != nil {
		d.Set("seg_storage_type", objectRaw["StorageType"])
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = gpdbServiceV2.DescribeDescribeDataReDistributeInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}

	objectRaw, err = gpdbServiceV2.DescribeDescribeDBInstanceSupportMaxPerformance(d.Id())
	if err != nil {
		return WrapError(err)
	}

	objectRaw, err = gpdbServiceV2.DescribeDescribeDBInstanceIPArrayList(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["DBInstanceIPArrayAttribute"] != nil {
		d.Set("db_instance_ip_array_attribute", objectRaw["DBInstanceIPArrayAttribute"])
	}
	if objectRaw["DBInstanceIPArrayName"] != nil {
		d.Set("db_instance_ip_array_name", objectRaw["DBInstanceIPArrayName"])
	}
	if objectRaw["SecurityIPList"] != nil {
		d.Set("security_ip_list", objectRaw["SecurityIPList"])
	}

	objectRaw, err = gpdbServiceV2.DescribeDescribeDBInstanceSSL(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["SSLEnabled"] != nil {
		d.Set("ssl_enabled", objectRaw["SSLEnabled"])
	}

	objectRaw, err = gpdbServiceV2.DescribeDescribeDBClusterNode(d.Id())
	if err != nil {
		return WrapError(err)
	}

	objectRaw, err = gpdbServiceV2.DescribeDescribeSampleData(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["SampleDataStatus"] != nil {
		d.Set("sample_data_status", objectRaw["SampleDataStatus"])
	}

	objectRaw, err = gpdbServiceV2.DescribeDescribeDataShareInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["DataShareStatus"] != nil {
		d.Set("data_share_status", objectRaw["DataShareStatus"])
	}

	objectRaw, err = gpdbServiceV2.DescribeDescribeDBResourceManagementMode(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["ResourceManagementMode"] != nil {
		d.Set("resource_management_mode", objectRaw["ResourceManagementMode"])
	}

	return nil
}

func resourceAliCloudGpdbDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyDBInstanceDescription"
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["DBInstanceDescription"] = d.Get("description")
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
	}
	update = false
	action = "ModifyDBInstanceMaintainTime"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	if d.HasChange("maintain_end_time") {
		update = true
	}
	request["EndTime"] = d.Get("maintain_end_time")
	if d.HasChange("maintain_start_time") {
		update = true
	}
	request["StartTime"] = d.Get("maintain_start_time")
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
	}
	update = false
	action = "ModifySecurityIps"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if v, ok := d.GetOk("modify_mode"); ok {
		request["ModifyMode"] = v
	}
	if d.HasChange("security_ip_list") {
		update = true
	}
	request["SecurityIPList"] = d.Get("security_ip_list")
	if d.HasChange("db_instance_ip_array_name") {
		update = true
		request["DBInstanceIPArrayName"] = d.Get("db_instance_ip_array_name")
	}

	if d.HasChange("db_instance_ip_array_attribute") {
		update = true
		request["DBInstanceIPArrayAttribute"] = d.Get("db_instance_ip_array_attribute")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
	}
	update = false
	action = "ModifyDBInstanceSSL"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("ssl_enabled") {
		update = true
	}
	request["SSLEnabled"] = d.Get("ssl_enabled")
	if v, ok := d.GetOk("connection_string"); ok {
		request["ConnectionString"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
		gpdbServiceV2 := GpdbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running", "IDLE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpgradeDBInstance"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()
	query["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if !d.IsNewResource() && d.HasChange("seg_node_num") {
		update = true
		request["SegNodeNum"] = d.Get("seg_node_num")
	}

	if v, ok := d.GetOk("instance_spec"); ok {
		request["InstanceSpec"] = v
	}
	if !d.IsNewResource() && d.HasChange("storage_size") {
		update = true
		request["StorageSize"] = d.Get("storage_size")
	}

	if v, ok := d.GetOk("upgrade_type"); ok {
		request["UpgradeType"] = v
	}
	if !d.IsNewResource() && d.HasChange("seg_disk_performance_level") {
		update = true
		request["SegDiskPerformanceLevel"] = d.Get("seg_disk_performance_level")
	}

	if v, ok := d.GetOk("seg_storage_type"); ok {
		request["SegStorageType"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing", "Throttling.User"}) || NeedRetry(err) {
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
		gpdbServiceV2 := GpdbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running", "IDLE"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyDBInstanceResourceGroup"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
	}
	update = false
	action = "ModifyVectorConfiguration"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("vector_configuration_status") {
		update = true
		request["VectorConfigurationStatus"] = d.Get("vector_configuration_status")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
		gpdbServiceV2 := GpdbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running", "IDLE"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyMasterSpec"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("master_cu") {
		update = true
	}
	request["MasterCU"] = d.Get("master_cu")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing", "Throttling.User"}) || NeedRetry(err) {
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
		gpdbServiceV2 := GpdbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running", "IDLE"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "SetDataShareInstance"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "InstanceList.0", d.Id())
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}
	query["RegionId"] = client.RegionId
	if d.HasChange("data_share_status") {
		update = true
	}
	request["OperationType"] = convertGpdbDBInstanceOperationTypeRequest(d.Get("data_share_status").(string))
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
		gpdbServiceV2 := GpdbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"opened", "closed"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyDBInstancePayType"
	conn, err = client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["PayType"] = convertGpdbDBInstancePayTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
	}

	if d.HasChange("status") {
		gpdbServiceV2 := GpdbServiceV2{client}
		object, err := gpdbServiceV2.DescribeGpdbDBInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["DBInstanceStatus"].(string) != target {
			if target == "STOPPED" {
				action := "PauseInstance"
				conn, err := client.NewGpdbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["DBInstanceId"] = d.Id()

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
				gpdbServiceV2 := GpdbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"STOPPED"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Running" {
				action := "ResumeInstance"
				conn, err := client.NewGpdbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["DBInstanceId"] = d.Id()

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
				gpdbServiceV2 := GpdbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}
	if d.HasChange("sample_data_status") {

		target := d.Get("sample_data_status").(string)
		if target == "loaded" {
			action := "CreateSampleData"
			conn, err := client.NewGpdbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["DBInstanceId"] = d.Id()

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
			gpdbServiceV2 := GpdbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"loaded"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
		if target == "unload" {
			action := "UnloadSampleData"
			conn, err := client.NewGpdbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["DBInstanceId"] = d.Id()

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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

		}
	}
	if d.HasChange("resource_management_mode") {

		target := d.Get("resource_management_mode").(string)
		if target == "resourceGroup" {
			action := "EnableDBResourceGroup"
			conn, err := client.NewGpdbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["DBInstanceId"] = d.Id()

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
			gpdbServiceV2 := GpdbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Running", "IDLE"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
		if target == "resourceQueue" {
			action := "DisableDBResourceGroup"
			conn, err := client.NewGpdbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["DBInstanceId"] = d.Id()

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
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
			gpdbServiceV2 := GpdbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Running", "IDLE"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
	}

	if d.HasChange("tags") {
		gpdbServiceV2 := GpdbServiceV2{client}
		if err := gpdbServiceV2.SetResourceTags(d, "instance"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudGpdbDBInstanceRead(d, meta)
}

func resourceAliCloudGpdbDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource alicloud_gpdb_instance which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewGpdbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["DBInstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-05-03"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, gpdbServiceV2.GpdbDBInstanceStateRefreshFunc(d.Id(), "DBInstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertGpdbDBInstanceItemsDBInstanceAttributePayTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Prepaid":
		return "Subscription"
	case "Postpaid":
		return "PayAsYouGo"
	}
	return source
}
func convertGpdbDBInstanceItemsDBInstanceAttributeSegDiskPerformanceLevelResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
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
func convertGpdbDBInstanceOperationTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "opened":
		return "add"
	case "closed":
		return "remove"
	}
	return source
}
func convertGpdbDBInstancePayTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}
func convertGpdbDBInstanceEnableSSLRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "1":
		return "true"
	case "0":
		return "false"
	}
	return source
}
func convertGpdbDBInstanceCreateSampleDataRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "loaded":
		return "true"
	case "unload":
		return "false"
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

func convertGpdbDbInstanceSSLEnabledResponse(source interface{}) interface{} {
	switch source {
	case false:
		return 0
	case true:
		return 1
	}

	return source
}
