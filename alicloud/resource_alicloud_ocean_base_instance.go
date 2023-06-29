package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOceanBaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOceanBaseInstanceCreate,
		Read:   resourceAlicloudOceanBaseInstanceRead,
		Update: resourceAlicloudOceanBaseInstanceUpdate,
		Delete: resourceAlicloudOceanBaseInstanceDelete,
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
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				Type:             schema.TypeBool,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Optional:         true,
				Type:             schema.TypeInt,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"backup_retain_mode": {
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"delete_all", "receive_all", "receive_last"}, false),
				Type:         schema.TypeString,
			},
			"commodity_code": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"cpu": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"disk_size": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: IntBetween(100, 10000),
			},
			"instance_class": {
				Required:     true,
				ValidateFunc: StringInSlice([]string{"14C70GB", "30C180GB", "62C400GB", "8C32GB", "16C70GB", "24C120GB", "32C160GB", "64C380GB", "20C32GB", "40C64GB", "4C16GB"}, false),
				Type:         schema.TypeString,
			},
			"instance_name": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"node_num": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"payment_type": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Type:         schema.TypeString,
			},
			"period": {
				Optional:         true,
				Type:             schema.TypeInt,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"period_unit": {
				Optional:         true,
				Type:             schema.TypeString,
				ValidateFunc:     StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"resource_group_id": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"series": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"normal", "normal_ssd", "history"}, false),
				Type:         schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"zones": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudOceanBaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oceanBaseProService := OceanBaseProService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	conn, err := client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("disk_size"); ok {
		request["DiskSize"] = v
	}
	if v, ok := d.GetOk("instance_class"); ok {
		request["InstanceClass"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertOceanBaseInstancePaymentTypeRequest(v)
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("series"); ok {
		request["Series"] = v
	}
	if v, ok := d.GetOk("zones"); ok {
		request["Zones"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	var response map[string]interface{}
	action := "CreateInstance"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ocean_base_instance", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.Data.InstanceId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ocean_base_instance")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, oceanBaseProService.OceanBaseInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudOceanBaseInstanceUpdate(d, meta)
}

func resourceAlicloudOceanBaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oceanBaseProService := OceanBaseProService{client}

	object, err := oceanBaseProService.DescribeOceanBaseInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ocean_base_instance oceanBaseProService.DescribeOceanBaseInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", object["CreateTime"])
	if v, ok := object["InstanceClass"]; ok && fmt.Sprint(v) != "" {
		d.Set("instance_class", fmt.Sprint(v, "B"))
	}
	d.Set("instance_name", object["InstanceName"])
	d.Set("node_num", object["NodeNum"])
	d.Set("payment_type", convertOceanBaseInstancePaymentTypeResponse(object["PayType"]))
	d.Set("series", convertOceanBaseInstanceSeriesResponse(object["Series"]))
	d.Set("status", object["Status"])
	zones, _ := jsonpath.Get("$.Zones", object)
	d.Set("zones", zones)
	d.Set("auto_renew", object["AutoRenewal"])

	object, err = oceanBaseProService.DescribeInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("commodity_code", object["CommodityCode"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("disk_size", formatInt(object["DiskSize"]))
	d.Set("cpu", formatInt(object["Cpu"]))

	return nil
}

func resourceAlicloudOceanBaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	conn, err := client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}
	oceanBaseProService := OceanBaseProService{client}
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
	}
	request["InstanceName"] = d.Get("instance_name")

	if update {
		action := "ModifyInstanceName"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("instance_name")
	}

	update = false
	request = map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	if d.HasChange("node_num") {
		update = true
	}
	request["NodeNum"] = d.Get("node_num")

	if update {
		action := "ModifyInstanceNodeNum"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Instance.Order.CreateFailed"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, oceanBaseProService.OceanBaseInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("node_num")
	}

	update = false
	request = map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	if !d.IsNewResource() && d.HasChange("disk_size") {
		update = true
		if v, ok := d.GetOk("disk_size"); ok {
			request["DiskSize"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("instance_class") {
		update = true
		if v, ok := d.GetOk("instance_class"); ok {
			request["InstanceClass"] = v
		}
	}

	if update {
		action := "ModifyInstanceSpec"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Instance.Order.CreateFailed"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, oceanBaseProService.OceanBaseInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("disk_size")
		d.SetPartial("instance_class")
	}

	d.Partial(false)
	return resourceAlicloudOceanBaseInstanceRead(d, meta)
}

func resourceAlicloudOceanBaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oceanBaseProService := OceanBaseProService{client}
	conn, err := client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceIds": fmt.Sprintf("[\"%s\"]", d.Id()),
		"RegionId":    client.RegionId,
	}

	if v, ok := d.GetOk("backup_retain_mode"); ok {
		request["BackupRetainMode"] = v
	}

	action := "DeleteInstances"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, oceanBaseProService.OceanBaseInstanceStateRefreshFunc(d.Id(), []string{}))
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

func convertOceanBaseInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	case "Subscription":
		return "PrePay"
	}
	return source
}

func convertOceanBaseInstanceSeriesResponse(source interface{}) interface{} {
	switch source {
	case "NORMAL":
		return "normal"
	case "BASIC":
		return "basic"
	}
	return source
}
