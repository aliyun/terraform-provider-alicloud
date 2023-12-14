// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOceanBaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOceanBaseInstanceCreate,
		Read:   resourceAliCloudOceanBaseInstanceRead,
		Update: resourceAliCloudOceanBaseInstanceUpdate,
		Delete: resourceAliCloudOceanBaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"backup_retain_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"receive_all", "delete_all", "receive_last"}, false),
			},
			"commodity_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(100, 10000),
			},
			"disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cloud_essd_pl1", "cloud_essd_pl0"}, false),
			},
			"instance_class": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"8C32GB", "14C70GB", "30C180GB", "62C400GB", "16C70GB", "24C120GB", "32C160GB", "64C380GB", "20C32GB", "40C64GB", "4C16GB", "32C180GB", "64C400GB", "16C32GB", "32C70GB", "64C180GB", "104C600GB"}, false),
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ob_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, 60}),
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc:     StringInSlice([]string{"Month", "Hour"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"series": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"normal", "normal_ssd", "history"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zones": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudOceanBaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["DiskSize"] = d.Get("disk_size")
	request["Series"] = d.Get("series")
	request["InstanceClass"] = d.Get("instance_class")
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ChargeType"] = convertOceanBaseChargeTypeRequest(d.Get("payment_type").(string))
	jsonPathResult10, err := jsonpath.Get("$", d.Get("zones"))
	if err != nil {
		return WrapError(err)
	}
	request["Zones"] = convertListToCommaSeparate(jsonPathResult10.(*schema.Set).List())

	if v, ok := d.GetOk("ob_version"); ok {
		request["ObVersion"] = v
	}
	if v, ok := d.GetOk("disk_type"); ok {
		request["DiskType"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), query, request, &runtime)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ocean_base_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	oceanBaseServiceV2 := OceanBaseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutCreate), 20*time.Minute, oceanBaseServiceV2.OceanBaseInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudOceanBaseInstanceUpdate(d, meta)
}

func resourceAliCloudOceanBaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oceanBaseServiceV2 := OceanBaseServiceV2{client}

	objectRaw, err := oceanBaseServiceV2.DescribeOceanBaseInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ocean_base_instance DescribeOceanBaseInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("disk_type", objectRaw["DiskType"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("node_num", objectRaw["NodeNum"])
	d.Set("ob_version", objectRaw["Version"])
	d.Set("payment_type", convertOceanBaseInstancePayTypeResponse(objectRaw["PayType"]))
	d.Set("series", convertOceanBaseInstanceSeriesResponse(objectRaw["Series"]))
	d.Set("status", objectRaw["Status"])
	d.Set("auto_renew", objectRaw["AutoRenewal"])

	zones1Raw := make([]interface{}, 0)
	if objectRaw["Zones"] != nil {
		zones1Raw = objectRaw["Zones"].([]interface{})
	}

	d.Set("zones", zones1Raw)

	e := jsonata.MustCompile("$.InstanceClass & 'B'")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("instance_class", evaluation)
	objectRaw, err = oceanBaseServiceV2.DescribeDescribeInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("commodity_code", objectRaw["CommodityCode"])
	d.Set("cpu", formatInt(objectRaw["Cpu"]))
	d.Set("disk_size", objectRaw["DiskSize"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])

	return nil
}

func resourceAliCloudOceanBaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyInstanceName"
	conn, err := client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		request["InstanceName"] = d.Get("instance_name")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), query, request, &runtime)

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
		d.SetPartial("instance_name")
	}
	update = false
	action = "ModifyInstanceNodeNum"
	conn, err = client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	if d.HasChange("node_num") {
		update = true
		request["NodeNum"] = d.Get("node_num")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), query, request, &runtime)

			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.Order.CreateFailed"}) || NeedRetry(err) {
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
		oceanBaseServiceV2 := OceanBaseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, oceanBaseServiceV2.OceanBaseInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("node_num")
	}
	update = false
	action = "ModifyInstanceSpec"
	conn, err = client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_class") {
		update = true
	}
	request["InstanceClass"] = d.Get("instance_class")
	if !d.IsNewResource() && d.HasChange("disk_size") {
		update = true
	}
	request["DiskSize"] = d.Get("disk_size")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), query, request, &runtime)

			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.Order.CreateFailed"}) || NeedRetry(err) {
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
		oceanBaseServiceV2 := OceanBaseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, oceanBaseServiceV2.OceanBaseInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_class")
		d.SetPartial("disk_size")
	}

	d.Partial(false)
	return resourceAliCloudOceanBaseInstanceRead(d, meta)
}

func resourceAliCloudOceanBaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstances"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceIds"] = fmt.Sprintf("[\"%s\"]", d.Id())

	if v, ok := d.GetOk("backup_retain_mode"); ok {
		request["BackupRetainMode"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), query, request, &runtime)

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
		if IsExpectedErrors(err, []string{"UnknownError", "IllegalOperation.Resource"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	oceanBaseServiceV2 := OceanBaseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 20*time.Second, oceanBaseServiceV2.OceanBaseInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertOceanBaseInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}

func convertOceanBaseInstancePayTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	case "PrePay":
		return "Subscription"
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}
func convertOceanBaseInstanceSeriesResponse(source interface{}) interface{} {
	switch source {
	case "NORMAL":
		return "normal"
	case "NORMAL_SSD":
		return "normal_ssd"
	case "HISTORY":
		return "history"
	}
	return source
}
func convertOceanBaseChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	case "Subscription":
		return "PrePay"
	}
	return source
}
