package alicloud

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/denverdino/aliyungo/common"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlikafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlikafkaInstanceCreate,
		Read:   resourceAlicloudAlikafkaInstanceRead,
		Update: resourceAlicloudAlikafkaInstanceUpdate,
		Delete: resourceAlicloudAlikafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"topic_quota": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute 'topic_quota' has been deprecated from 1.194.0 and it will be removed in the next future. Using new attribute 'partition_num' instead.",
			},
			"partition_num": {
				Type:          schema.TypeInt,
				Optional:      true,
				AtLeastOneOf:  []string{"partition_num", "topic_quota"},
				ConflictsWith: []string{"topic_quota"},
			},
			"disk_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"deploy_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 5}),
			},
			"io_max": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"paid_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
			},
			"spec_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "normal",
			},
			"eip_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(int) == 5
				},
			},
			"security_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: alikafkaInstanceConfigDiffSuppressFunc,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"end_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tagsSchema(),
			"selected_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudAlikafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	vpcService := VpcService{client}
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}

	// 1. Create order
	var createOrderAction string
	createOrderResponse := make(map[string]interface{})
	createOrderReq := make(map[string]interface{})
	createOrderReq["RegionId"] = client.RegionId
	if v, ok := d.GetOk("partition_num"); ok {
		createOrderReq["PartitionNum"] = v
	} else if v, ok := d.GetOk("topic_quota"); ok {
		createOrderReq["TopicQuota"] = v
	}
	createOrderReq["DiskType"] = d.Get("disk_type")
	createOrderReq["DiskSize"] = d.Get("disk_size")
	createOrderReq["DeployType"] = d.Get("deploy_type")
	createOrderReq["IoMax"] = d.Get("io_max")
	if v, ok := d.GetOk("spec_type"); ok {
		createOrderReq["SpecType"] = v
	}
	if v, ok := d.GetOkExists("eip_max"); ok {
		createOrderReq["EipMax"] = v
	}
	if v, ok := d.GetOk("paid_type"); ok {
		switch v.(string) {
		case "PostPaid":
			createOrderAction = "CreatePostPayOrder"
		case "PrePaid":
			createOrderAction = "CreatePrePayOrder"
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		createOrderResponse, err = conn.DoRequest(StringPointer(createOrderAction), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, createOrderReq, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL", "ONS_SYSTEM_ERROR"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(createOrderAction, createOrderResponse, createOrderReq)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", createOrderAction, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(createOrderResponse["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", createOrderAction, createOrderResponse))
	}

	alikafkaInstanceVO, err := alikafkaService.DescribeAliKafkaInstanceByOrderId(fmt.Sprint(createOrderResponse["OrderId"]), 60)
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprint(alikafkaInstanceVO["InstanceId"]))

	// 2. Start instance
	startInstanceAction := "StartInstance"
	startInstanceResponse := make(map[string]interface{})
	startInstanceReq := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		startInstanceReq["VpcId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		startInstanceReq["ZoneId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		startInstanceReq["VSwitchId"] = v
	}

	if (startInstanceReq["ZoneId"] == nil || startInstanceReq["VpcId"] == nil) && startInstanceReq["VSwitchId"] != nil {
		vsw, err := vpcService.DescribeVswitch(startInstanceReq["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if v, ok := startInstanceReq["VpcId"].(string); !ok || v == "" {
			startInstanceReq["VpcId"] = vsw["VpcId"]
		}
		if v, ok := startInstanceReq["ZoneId"].(string); !ok || v == "" {
			startInstanceReq["ZoneId"] = vsw["ZoneId"]
		}
	}

	startInstanceReq["RegionId"] = client.RegionId
	startInstanceReq["InstanceId"] = alikafkaInstanceVO["InstanceId"]
	if _, ok := d.GetOkExists("eip_max"); ok {
		startInstanceReq["DeployModule"] = "eip"
		startInstanceReq["IsEipInner"] = true
	}
	if v, ok := d.GetOk("name"); ok {
		startInstanceReq["Name"] = v
	}
	if v, ok := d.GetOk("security_group"); ok {
		startInstanceReq["SecurityGroup"] = v
	}
	if v, ok := d.GetOk("service_version"); ok {
		startInstanceReq["ServiceVersion"] = v
	}
	if v, ok := d.GetOk("config"); ok {
		startInstanceReq["Config"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		startInstanceReq["KMSKeyId"] = v
	}
	if v, ok := d.GetOk("selected_zones"); ok {
		startInstanceReq["SelectedZones"] = convertListToJsonString(v.([]interface{}))
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		startInstanceResponse, err = conn.DoRequest(StringPointer(startInstanceAction), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, startInstanceReq, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(startInstanceAction, startInstanceResponse, startInstanceReq)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", startInstanceAction, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(startInstanceResponse["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", startInstanceAction, startInstanceResponse))
	}

	// 3. wait until running
	stateConf := BuildStateConf([]string{}, []string{"5"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudAlikafkaInstanceUpdate(d, meta)
}

func resourceAlicloudAlikafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	object, err := alikafkaService.DescribeAliKafkaInstance(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object["Name"])
	d.Set("disk_type", object["DiskType"])
	d.Set("disk_size", object["DiskSize"])
	d.Set("deploy_type", object["DeployType"])
	d.Set("io_max", object["IoMax"])
	d.Set("eip_max", object["EipMax"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("paid_type", PostPaid)
	d.Set("spec_type", object["SpecType"])
	d.Set("security_group", object["SecurityGroup"])
	d.Set("end_point", object["EndPoint"])
	d.Set("status", object["ServiceStatus"])
	// object.UpgradeServiceDetailInfo.UpgradeServiceDetailInfoVO[0].Current2OpenSourceVersion can guaranteed not to be null
	d.Set("service_version", object["UpgradeServiceDetailInfo"].(map[string]interface{})["Current2OpenSourceVersion"])
	d.Set("config", object["AllConfig"])
	d.Set("kms_key_id", object["KmsKeyId"])
	if fmt.Sprint(object["PaidType"]) == "0" {
		d.Set("paid_type", PrePaid)
	}

	quota, err := alikafkaService.GetQuotaTip(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("topic_quota", quota["TopicQuota"])
	d.Set("partition_num", quota["PartitionNumOfBuy"])

	tags, err := alikafkaService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", alikafkaService.tagsToMap(tags))

	return nil
}

func resourceAlicloudAlikafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if err := alikafkaService.setInstanceTags(d, TagResourceInstance); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudAlikafkaInstanceRead(d, meta)
	}

	// Process change instance name.
	if d.HasChange("name") {
		action := "ModifyInstanceName"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("name"); ok {
			request["InstanceName"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("name")
	}

	// Process paid type change, note only support change from post to pre pay.
	if d.HasChange("paid_type") {
		o, n := d.GetChange("paid_type")
		oldPaidType := o.(string)
		newPaidType := n.(string)
		oldPaidTypeInt := 1
		newPaidTypeInt := 1
		if oldPaidType == string(PrePaid) {
			oldPaidTypeInt = 0
		}
		if newPaidType == string(PrePaid) {
			newPaidTypeInt = 0
		}
		if oldPaidTypeInt == 1 && newPaidTypeInt == 0 {
			action := "ConvertPostPayOrder"
			request := map[string]interface{}{
				"RegionId":   client.RegionId,
				"InstanceId": d.Id(),
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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

			err = alikafkaService.WaitForAliKafkaInstanceUpdated(d, strconv.Itoa(newPaidTypeInt), DefaultTimeoutMedium)
			if err != nil {
				return WrapError(err)
			}
		} else {
			return WrapError(errors.New("paid type only support change from post pay to pre pay"))
		}

		d.SetPartial("paid_type")
	}

	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}
	if d.HasChange("partition_num") {
		update = true
	}
	request["PartitionNum"] = d.Get("partition_num")
	if d.HasChange("topic_quota") {
		update = true
		request["TopicQuota"] = d.Get("topic_quota")
	}
	if d.HasChange("disk_size") {
		update = true
	}
	request["DiskSize"] = d.Get("disk_size")
	if d.HasChange("io_max") {
		update = true
	}
	request["IoMax"] = d.Get("io_max")
	if d.HasChange("spec_type") {
		update = true
	}
	request["SpecType"] = d.Get("spec_type")

	if d.HasChange("deploy_type") {
		update = true
	}
	if d.Get("deploy_type").(int) == 4 {
		request["EipModel"] = true
	} else {
		request["EipModel"] = false
	}
	if d.HasChange("eip_max") {
		update = true
	}
	request["EipMax"] = d.Get("eip_max").(int)

	if update {
		action := "UpgradePostPayOrder"
		if d.Get("paid_type").(string) == string(PrePaid) {
			action = "UpgradePrePayOrder"
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, raw, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_instances", action, AlibabaCloudSdkGoERROR)
		}

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("partition_num")
		d.SetPartial("topic_quota")
		d.SetPartial("disk_size")
		d.SetPartial("io_max")
		d.SetPartial("spec_type")
		d.SetPartial("eip_max")

		paidType := 1
		if d.Get("paid_type").(string) == string(PrePaid) {
			paidType = 0
		}
		err = alikafkaService.WaitForAliKafkaInstanceUpdated(d, strconv.Itoa(paidType), DefaultTimeoutMedium)
		if err != nil {
			return WrapError(err)
		}

		stateConf := BuildStateConf([]string{}, []string{"5"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("service_version") {
		action := "UpgradeInstanceVersion"
		request := map[string]interface{}{
			"InstanceId": d.Id(),
			"RegionId":   client.RegionId,
		}

		if v, ok := d.GetOk("service_version"); ok {
			request["TargetVersion"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				// means no need to update version
				if IsExpectedErrors(err, []string{"ONS_INIT_ENV_ERROR"}) {
					return nil
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// wait for upgrade task be invoke
		time.Sleep(60 * time.Second)
		// upgrade service may be last a long time
		stateConf := BuildStateConf([]string{}, []string{"5"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("service_version")
	}

	if d.HasChange("config") {
		action := "UpdateInstanceConfig"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("config"); ok {
			request["Config"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// wait for upgrade task be invoke
		stateConf := BuildStateConf([]string{}, []string{"5"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("config")
	}

	d.Partial(false)
	return resourceAlicloudAlikafkaInstanceRead(d, meta)
}

func resourceAlicloudAlikafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	action := "ReleaseInstance"
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId":          d.Id(),
		"RegionId":            client.RegionId,
		"ForceDeleteInstance": true,
	}

	// Pre paid instance can not be release.
	if d.Get("paid_type").(string) == string(PrePaid) {
		return nil
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	stateConf := BuildStateConf([]string{}, []string{"15"}, client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	action = "DeleteInstance"
	request = map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	stateConf = BuildStateConf([]string{}, []string{}, client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
