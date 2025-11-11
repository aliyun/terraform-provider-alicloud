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
			"components": {
				Type:     schema.TypeSet,
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
						"disk_size_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Normal", "Large"}, false),
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
			"ha": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"multi_zone_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	if v, ok := d.GetOk("components"); ok {
		componentsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["cuNum"] = dataLoopTmp["cu_num"]
			dataLoopMap["type"] = dataLoopTmp["type"]
			dataLoopMap["replica"] = dataLoopTmp["replica"]
			dataLoopMap["diskSizeType"] = dataLoopTmp["disk_size_type"]
			dataLoopMap["cuType"] = dataLoopTmp["cu_type"]
			componentsMapsArray = append(componentsMapsArray, dataLoopMap)
		}
		request["components"] = componentsMapsArray
	}

	request["paymentType"] = d.Get("payment_type")
	if v, ok := d.GetOk("configuration"); ok {
		request["configuration"] = v
	}
	if v, ok := d.GetOk("vswitch_ids"); ok {
		vSwitchIdsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["zoneId"] = dataLoop1Tmp["zone_id"]
			dataLoop1Map["vswId"] = dataLoop1Tmp["vsw_id"]
			vSwitchIdsMapsArray = append(vSwitchIdsMapsArray, dataLoop1Map)
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
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	request["instanceName"] = d.Get("instance_name")
	if v, ok := d.GetOkExists("encrypted"); ok {
		request["encrypted"] = v
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
	d.Set("ha", objectRaw["ha"])
	d.Set("instance_name", objectRaw["instanceName"])
	d.Set("kms_key_id", objectRaw["kmsKeyId"])
	d.Set("multi_zone_mode", objectRaw["multiZoneMode"])
	d.Set("payment_type", objectRaw["paymentType"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
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
			componentsMap["replica"] = componentsChildRaw["replica"]
			componentsMap["type"] = componentsChildRaw["type"]

			componentsMaps = append(componentsMaps, componentsMap)
		}
	}
	if err := d.Set("components", componentsMaps); err != nil {
		return err
	}
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
	d.Partial(true)

	var err error
	action := fmt.Sprintf("/webapi/instance/update")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instanceId"] = d.Id()
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("components") {
		update = true
	}
	if v, ok := d.GetOk("components"); ok && d.HasChange("components") {
		componentsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["type"] = dataLoopTmp["type"]
			dataLoopMap["cuNum"] = dataLoopTmp["cu_num"]
			dataLoopMap["replica"] = dataLoopTmp["replica"]
			componentsMapsArray = append(componentsMapsArray, dataLoopMap)
		}
		request["components"] = componentsMapsArray
	}

	if d.HasChange("ha") {
		update = true
	}
	if v, ok := d.GetOk("ha"); ok && d.HasChange("ha") {
		request["ha"] = v
	}
	if d.HasChange("instance_name") {
		update = true
	}
	request["instanceName"] = d.Get("instance_name")
	if d.HasChange("auto_backup") {
		update = true
	}
	if v, ok := d.GetOk("auto_backup"); ok && d.HasChange("auto_backup") {
		request["autoBackup"] = v
	}
	if d.HasChange("configuration") {
		update = true
	}
	if v, ok := d.GetOk("configuration"); ok && d.HasChange("configuration") {
		request["configuration"] = v
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
	d.Partial(false)
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
