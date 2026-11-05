// Package milvus. This file is generated automatically. Please do not modify it manually, thank you!
package milvus

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,
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
				ForceNew: true,
			},
			"components": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      componentsHash,
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
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_class": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"performance_level": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"disk_size_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validators.StringInSlice([]string{"Normal", "Large"}, false),
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
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
				ForceNew: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"load_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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
				ValidateFunc: validators.StringInSlice([]string{"month", "year"}, false),
			},
			"payment_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"promotion_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"tags": tags.TagsSchema(),
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

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/webapi/instance/create")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = helper.StringPointer(client.RegionId)
	query["clientToken"] = helper.StringPointer(helper.BuildClientToken(action))

	if v, ok := d.GetOkExists("load_replicas"); ok {
		request["loadReplicas"] = v
	}
	if v, ok := d.GetOk("components"); ok {
		componentsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["cuNum"] = dataLoopTmp["cu_num"]
			for _, dataDisks := range convert.ConvertToInterfaceArray(dataLoopTmp["data_disk"]) {
				dataDisksMap := map[string]interface{}{}
				dataDisksArg := dataDisks.(map[string]interface{})

				if size, ok := dataDisksArg["size"]; ok {
					dataDisksMap["size"] = size
				}

				if storageClass, ok := dataDisksArg["storage_class"]; ok {
					dataDisksMap["storageClass"] = storageClass
				}

				if enabled, ok := dataDisksArg["enabled"]; ok {
					dataDisksMap["enabled"] = enabled
				}

				if performanceLevel, ok := dataDisksArg["performance_level"]; ok {
					dataDisksMap["performanceLevel"] = performanceLevel
				}

				if len(dataDisksMap) > 0 {
					dataLoopMap["dataDisk"] = dataDisksMap
				}
			}
			dataLoopMap["type"] = dataLoopTmp["type"]
			dataLoopMap["replica"] = dataLoopTmp["replica"]
			dataLoopMap["diskSizeType"] = dataLoopTmp["disk_size_type"]
			dataLoopMap["cuType"] = dataLoopTmp["cu_type"]
			componentsMapsArray = append(componentsMapsArray, dataLoopMap)
		}
		request["components"] = componentsMapsArray
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
		for _, dataLoop2 := range convert.ConvertToInterfaceArray(v) {
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
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("milvus", "2023-10-12", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_milvus_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.instanceId", response)
	d.SetId(fmt.Sprint(id))

	milvusService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 4*time.Minute, milvusService.InstanceStateRefreshFunc(d.Id(), "status", []string{"creating_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	milvusService := NewService(client)

	objectRaw, err := milvusService.DescribeInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_milvus_instance DescribeInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
		for _, componentsChildRaw := range convert.ConvertToInterfaceArray(componentsRaw) {
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
				dataDiskMap["enabled"] = dataDiskRaw["enabled"]
				dataDiskMap["performance_level"] = dataDiskRaw["performanceLevel"]
				dataDiskMap["size"] = dataDiskRaw["size"]
				dataDiskMap["storage_class"] = dataDiskRaw["storageClass"]

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
		securityGroupIdsRaw = convert.ConvertToInterfaceArray(objectRaw["securityGroupIds"])
	}

	d.Set("security_group_ids", securityGroupIdsRaw)
	tagsMaps := objectRaw["tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))
	vSwitchIdsRaw := objectRaw["vSwitchIds"]
	vSwitchIdsMaps := make([]map[string]interface{}, 0)
	if vSwitchIdsRaw != nil {
		for _, vSwitchIdsChildRaw := range convert.ConvertToInterfaceArray(vSwitchIdsRaw) {
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

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
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
	query["RegionId"] = helper.StringPointer(client.RegionId)
	query["clientToken"] = helper.StringPointer(helper.BuildClientToken(action))
	if d.HasChange("components") {
		update = true
	}
	if v, ok := d.GetOk("components"); ok && d.HasChange("components") {
		componentsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["type"] = dataLoopTmp["type"]

			for _, dataDisks := range convert.ConvertToInterfaceArray(dataLoopTmp["data_disk"]) {
				dataDisksMap := map[string]interface{}{}
				dataDisksArg := dataDisks.(map[string]interface{})

				if size, ok := dataDisksArg["size"]; ok {
					dataDisksMap["size"] = size
				}

				if storageClass, ok := dataDisksArg["storage_class"]; ok {
					dataDisksMap["storageClass"] = storageClass
				}

				if enabled, ok := dataDisksArg["enabled"]; ok {
					dataDisksMap["enabled"] = enabled
				}

				if performanceLevel, ok := dataDisksArg["performance_level"]; ok {
					dataDisksMap["performanceLevel"] = performanceLevel
				}

				if len(dataDisksMap) > 0 {
					dataLoopMap["dataDisk"] = dataDisksMap
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("milvus", "2023-10-12", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		milvusService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, milvusService.InstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = fmt.Sprintf("/webapi/resourceGroup/change")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = helper.StringPointer(d.Id())
	query["RegionId"] = helper.StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["NewResourceGroupId"] = helper.StringPointer(v.(string))
	}

	query["ResourceType"] = helper.StringPointer("instance")
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("milvus", "2023-10-12", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		milvusService := NewService(client)
		if err := milvusService.SetResourceTags(d, "instance"); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := true
	if v, ok := d.GetOk("payment_type"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"Subscription"}) {
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
		query["instanceId"] = helper.StringPointer(d.Id())
		query["RegionId"] = helper.StringPointer(client.RegionId)

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RoaDelete("milvus", "2023-10-12", action, query, nil, nil, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		milvusService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 4*time.Minute, milvusService.InstanceStateRefreshFunc(d.Id(), "status", []string{"deleting_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	enableDelete = false
	if v, ok := d.GetOk("payment_type"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"Subscription"}) {
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
		request["InstanceId"] = helper.StringPointer(d.Id())

		request["clientToken"] = helper.BuildClientToken(action)

		request["ImmediatelyRelease"] = helper.StringPointer("1")
		var endpoint string
		request["ProductCode"] = helper.StringPointer("milvus")
		request["ProductType"] = helper.StringPointer("milvus_milvuspre_public_cn")
		if client.IsInternationalAccount() {
			request["ProductType"] = ""
		}

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{""}) {
					request["ProductCode"] = ""
					request["ProductType"] = ""
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		milvusService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 4*time.Minute, milvusService.InstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}
	return nil
}

// componentsHash computes the set identity hash for a Milvus instance
// component. It intentionally excludes data_disk (Optional+Computed, populated
// with server defaults by the API even when the user leaves it unset), pay_type
// (Computed only) and cu_type (Optional+Computed, defaulted to "general" by the
// API when the user leaves it unset) so that a component written in config
// without those fields still hashes identically to the same component read back
// from the API. Without this the TypeSet would show a perpetual diff, because
// the default schema.hashResource hashes Computed fields that are empty in
// config but populated in state. The remaining identity fields are sufficient
// to uniquely distinguish a component within one instance.
func componentsHash(i interface{}) int {
	m, ok := i.(map[string]interface{})
	if !ok {
		return 0
	}
	return schema.HashString(fmt.Sprintf("%v|%v|%v|%v",
		m["type"], m["disk_size_type"], m["replica"], m["cu_num"]))
}
