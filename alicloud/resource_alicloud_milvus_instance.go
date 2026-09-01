// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMilvusInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMilvusInstanceCreate,
		Read:   resourceAliCloudMilvusInstanceRead,
		Update: resourceAliCloudMilvusInstanceUpdate,
		Delete: resourceAliCloudMilvusInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(14 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(14 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_backup": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backup_restore_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_cluster_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backup_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"components": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cu_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cu_num": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"data_disk": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"performance_level": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"disk_size_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"Normal", "Large"}, false),
						},
						"pay_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replica": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"configuration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_admin_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ha": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_multi_az_storage": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"load_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"multi_zone_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"payment_duration_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"month", "year"}, false),
			},
			"payment_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"promotion_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"running_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vsw_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudMilvusInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/webapi/instance/create")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)
	query["clientToken"] = StringPointer(buildClientToken(action))

	if v, ok := d.GetOkExists("load_replicas"); ok {
		request["loadReplicas"] = v
	}
	if v, ok := d.GetOk("components"); ok {
		componentsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["cuNum"] = dataLoopTmp["cu_num"]
			localData1 := make(map[string]interface{})
			size1, _ := jsonpath.Get("$[0].size", dataLoopTmp["data_disk"])
			if size1 != nil && size1 != "" {
				localData1["size"] = size1
			}
			storageClass1, _ := jsonpath.Get("$[0].storage_class", dataLoopTmp["data_disk"])
			if storageClass1 != nil && storageClass1 != "" {
				localData1["storageClass"] = storageClass1
			}
			enabled1, _ := jsonpath.Get("$[0].enabled", dataLoopTmp["data_disk"])
			if enabled1 != nil && enabled1 != "" {
				localData1["enabled"] = enabled1
			}
			performanceLevel1, _ := jsonpath.Get("$[0].performance_level", dataLoopTmp["data_disk"])
			if performanceLevel1 != nil && performanceLevel1 != "" {
				localData1["performanceLevel"] = performanceLevel1
			}
			if len(localData1) > 0 {
				dataLoopMap["dataDisk"] = localData1
			}
			dataLoopMap["type"] = dataLoopTmp["type"]
			dataLoopMap["replica"] = dataLoopTmp["replica"]
			dataLoopMap["diskSizeType"] = dataLoopTmp["disk_size_type"]
			dataLoopMap["cuType"] = dataLoopTmp["cu_type"]
			componentsMapsArray = append(componentsMapsArray, dataLoopMap)
		}
		request["components"] = componentsMapsArray
	}

	backupRestoreInfo := make(map[string]interface{})

	if v := d.Get("backup_restore_info"); !IsNil(v) {
		backupId1, _ := jsonpath.Get("$[0].backup_id", v)
		if backupId1 != nil && backupId1 != "" {
			backupRestoreInfo["backupId"] = backupId1
		}
		backupName1, _ := jsonpath.Get("$[0].backup_name", v)
		if backupName1 != nil && backupName1 != "" {
			backupRestoreInfo["backupName"] = backupName1
		}
		sourceClusterId1, _ := jsonpath.Get("$[0].source_cluster_id", v)
		if sourceClusterId1 != nil && sourceClusterId1 != "" {
			backupRestoreInfo["sourceClusterId"] = sourceClusterId1
		}

		request["backupRestoreInfo"] = backupRestoreInfo
	}

	request["paymentType"] = d.Get("payment_type")
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["autoPay"] = v
	}
	if v, ok := d.GetOk("promotion_no"); ok {
		request["promotionNo"] = v
	}
	if v, ok := d.GetOk("configuration"); ok {
		request["configuration"] = v
	}
	if v, ok := d.GetOk("vswitch_ids"); ok {
		vSwitchIdsMapsArray := make([]interface{}, 0)
		for _, dataLoop2 := range convertToInterfaceArray(v) {
			dataLoop2Tmp := dataLoop2.(map[string]interface{})
			dataLoop2Map := make(map[string]interface{})
			dataLoop2Map["zoneId"] = dataLoop2Tmp["zone_id"]
			dataLoop2Map["vswId"] = dataLoop2Tmp["vsw_id"]
			vSwitchIdsMapsArray = append(vSwitchIdsMapsArray, dataLoop2Map)
		}
		request["vSwitchIds"] = vSwitchIdsMapsArray
	}

	if v, ok := d.GetOk("multi_zone_mode"); ok {
		request["multiZoneMode"] = v
	}
	if v, ok := d.GetOk("payment_duration_unit"); ok {
		request["paymentDurationUnit"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request["kmsKeyId"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["autoRenew"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	request["instanceName"] = d.Get("instance_name")
	if v, ok := d.GetOkExists("encrypted"); ok {
		request["encrypted"] = v
	}
	if v, ok := d.GetOkExists("is_multi_az_storage"); ok {
		request["isMultiAzStorage"] = v
	}
	if v, ok := d.GetOkExists("ha"); ok {
		request["ha"] = v
	}
	request["vpcId"] = d.Get("vpc_id")
	request["dbVersion"] = d.Get("db_version")
	if v, ok := d.GetOkExists("auto_backup"); ok {
		request["autoBackup"] = v
	}
	request["zoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("db_admin_password"); ok {
		request["dbAdminPassword"] = v
	}
	if v, ok := d.GetOkExists("payment_duration"); ok {
		request["paymentDuration"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("milvus", "2023-10-12", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_milvus_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.instanceId", response)
	d.SetId(fmt.Sprint(id))

	milvusServiceV2 := MilvusServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 4*time.Minute, milvusServiceV2.MilvusInstanceStateRefreshFunc(d.Id(), "status", []string{"creating_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudMilvusInstanceRead(d, meta)
}

func resourceAliCloudMilvusInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	milvusServiceV2 := MilvusServiceV2{client}

	objectRaw, err := milvusServiceV2.DescribeMilvusInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_milvus_instance DescribeMilvusInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_backup", objectRaw["autoBackup"])
	d.Set("configuration", objectRaw["configuration"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("db_version", objectRaw["dbVersion"])
	d.Set("encrypted", objectRaw["encrypted"])
	d.Set("expire_time", objectRaw["expireTime"])
	d.Set("ha", objectRaw["ha"])
	d.Set("instance_name", objectRaw["instanceName"])
	d.Set("kms_key_id", objectRaw["kmsKeyId"])
	d.Set("multi_zone_mode", objectRaw["multiZoneMode"])
	d.Set("order_id", objectRaw["orderId"])
	d.Set("payment_type", objectRaw["paymentType"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("running_time", objectRaw["runningTime"])
	d.Set("status", objectRaw["status"])
	d.Set("vpc_id", objectRaw["vpcId"])
	d.Set("zone_id", objectRaw["zoneId"])

	componentsRaw := objectRaw["components"]
	componentsMaps := make([]map[string]interface{}, 0)
	if componentsRaw != nil {
		for _, componentsChildRaw := range convertToInterfaceArray(componentsRaw) {
			componentsMap := make(map[string]interface{})
			componentsChildRaw := componentsChildRaw.(map[string]interface{})
			componentsMap["cu_num"] = componentsChildRaw["cuNum"]
			componentsMap["cu_type"] = componentsChildRaw["cuType"]
			componentsMap["disk_size_type"] = componentsChildRaw["diskSizeType"]
			componentsMap["pay_type"] = componentsChildRaw["payType"]
			componentsMap["replica"] = componentsChildRaw["replica"]
			componentsMap["type"] = componentsChildRaw["type"]

			dataDiskMaps := make([]map[string]interface{}, 0)
			dataDiskMap := make(map[string]interface{})
			dataDiskRaw := make(map[string]interface{})
			if componentsChildRaw["dataDisk"] != nil {
				dataDiskRaw = componentsChildRaw["dataDisk"].(map[string]interface{})
			}
			if len(dataDiskRaw) > 0 {
				dataDiskMap["enabled"] = dataDiskRaw["Enabled"]
				dataDiskMap["performance_level"] = dataDiskRaw["PerformanceLevel"]
				dataDiskMap["size"] = dataDiskRaw["Size"]
				dataDiskMap["storage_class"] = dataDiskRaw["StorageClass"]

				dataDiskMaps = append(dataDiskMaps, dataDiskMap)
			}
			componentsMap["data_disk"] = dataDiskMaps
			componentsMaps = append(componentsMaps, componentsMap)
		}
	}
	if err := d.Set("components", componentsMaps); err != nil {
		return err
	}
	securityGroupIdsRaw := make([]interface{}, 0)
	if objectRaw["securityGroupIds"] != nil {
		securityGroupIdsRaw = convertToInterfaceArray(objectRaw["securityGroupIds"])
	}

	d.Set("security_group_ids", securityGroupIdsRaw)
	tagsMaps := objectRaw["tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	vSwitchIdsRaw := objectRaw["vSwitchIds"]
	vSwitchIdsMaps := make([]map[string]interface{}, 0)
	if vSwitchIdsRaw != nil {
		for _, vSwitchIdsChildRaw := range convertToInterfaceArray(vSwitchIdsRaw) {
			vSwitchIdsMap := make(map[string]interface{})
			vSwitchIdsChildRaw := vSwitchIdsChildRaw.(map[string]interface{})
			vSwitchIdsMap["vsw_id"] = vSwitchIdsChildRaw["vswId"]
			vSwitchIdsMap["zone_id"] = vSwitchIdsChildRaw["zoneId"]

			vSwitchIdsMaps = append(vSwitchIdsMaps, vSwitchIdsMap)
		}
	}
	if err := d.Set("vswitch_ids", vSwitchIdsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudMilvusInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/webapi/instance/update")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()
	query["RegionId"] = StringPointer(client.RegionId)
	query["clientToken"] = StringPointer(buildClientToken(action))
	if d.HasChange("components") {
		update = true
	}
	if v, ok := d.GetOk("components"); ok && d.HasChange("components") {
		componentsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["type"] = dataLoopTmp["type"]
			if !IsNil(dataLoopTmp["data_disk"]) {
				localData1 := make(map[string]interface{})
				storageClass1, _ := jsonpath.Get("$[0].storage_class", dataLoopTmp["data_disk"])
				if storageClass1 != nil && storageClass1 != "" {
					localData1["storageClass"] = storageClass1
				}
				enabled1, _ := jsonpath.Get("$[0].enabled", dataLoopTmp["data_disk"])
				if enabled1 != nil && enabled1 != "" {
					localData1["enabled"] = enabled1
				}
				performanceLevel1, _ := jsonpath.Get("$[0].performance_level", dataLoopTmp["data_disk"])
				if performanceLevel1 != nil && performanceLevel1 != "" {
					localData1["performanceLevel"] = performanceLevel1
				}
				size1, _ := jsonpath.Get("$[0].size", dataLoopTmp["data_disk"])
				if size1 != nil && size1 != "" {
					localData1["size"] = size1
				}
				if len(localData1) > 0 {
					dataLoopMap["dataDisk"] = localData1
				}
			}
			dataLoopMap["cuNum"] = dataLoopTmp["cu_num"]
			dataLoopMap["replica"] = dataLoopTmp["replica"]
			dataLoopMap["cuType"] = dataLoopTmp["cu_type"]
			componentsMapsArray = append(componentsMapsArray, dataLoopMap)
		}
		request["components"] = componentsMapsArray
	}

	if d.HasChange("instance_name") {
		update = true
	}
	request["instanceName"] = d.Get("instance_name")
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["autoPay"] = v
	}
	if d.HasChange("configuration") {
		update = true
	}
	if v, ok := d.GetOk("configuration"); ok && d.HasChange("configuration") {
		request["configuration"] = v
	}
	if d.HasChange("ha") {
		update = true
	}
	if v, ok := d.GetOk("ha"); ok && d.HasChange("ha") {
		request["ha"] = v
	}
	if d.HasChange("auto_backup") {
		update = true
	}
	if v, ok := d.GetOk("auto_backup"); ok && d.HasChange("auto_backup") {
		request["autoBackup"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("milvus", "2023-10-12", action, query, nil, body, true)
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
		milvusServiceV2 := MilvusServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, milvusServiceV2.MilvusInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = fmt.Sprintf("/webapi/resourceGroup/change")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["NewResourceGroupId"] = StringPointer(v.(string))
	}

	query["ResourceType"] = StringPointer("instance")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("milvus", "2023-10-12", action, query, nil, body, true)
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

	if d.HasChange("tags") {
		milvusServiceV2 := MilvusServiceV2{client}
		if err := milvusServiceV2.SetResourceTags(d, "instance"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudMilvusInstanceRead(d, meta)
}

func resourceAliCloudMilvusInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := true
	if v, ok := d.GetOk("payment_type"); ok {
		if InArray(fmt.Sprint(v), []string{"Subscription"}) {
			enableDelete = false
			log.Printf("[WARN] Cannot destroy resource alicloud_milvus_instance which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		action := fmt.Sprintf("/webapi/instance/delete")
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		var err error
		request = make(map[string]interface{})
		query["instanceId"] = StringPointer(d.Id())
		query["RegionId"] = StringPointer(client.RegionId)

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RoaDelete("milvus", "2023-10-12", action, query, nil, nil, true)
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
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		milvusServiceV2 := MilvusServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 4*time.Minute, milvusServiceV2.MilvusInstanceStateRefreshFunc(d.Id(), "status", []string{"deleting_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	enableDelete = false
	if v, ok := d.GetOk("payment_type"); ok {
		if InArray(fmt.Sprint(v), []string{"Subscription"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		action := "RefundInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		var err error
		request = make(map[string]interface{})
		request["InstanceId"] = StringPointer(d.Id())

		request["clientToken"] = buildClientToken(action)

		request["ImmediatelyRelease"] = StringPointer("1")
		var endpoint string
		request["ProductCode"] = StringPointer("milvus")
		request["ProductType"] = StringPointer("milvus_milvuspre_public_cn")
		if client.IsInternationalAccount() {
			request["ProductType"] = ""
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{""}) {
					request["ProductCode"] = ""
					request["ProductType"] = ""
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		milvusServiceV2 := MilvusServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 4*time.Minute, milvusServiceV2.MilvusInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	return nil
}
