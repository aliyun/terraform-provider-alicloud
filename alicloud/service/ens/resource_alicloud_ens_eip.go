// Package ens. This file is generated automatically. Please do not modify it manually, thank you!
package ens

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEnsEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnsEipCreate,
		Read:   resourceEnsEipRead,
		Update: resourceEnsEipUpdate,
		Delete: resourceEnsEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eip_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"95BandwidthByMonth"}, false),
			},
			"isp": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEnsEipCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEipInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = helper.BuildClientToken(action)

	request["EnsRegionId"] = d.Get("ens_region_id")
	if v, ok := d.GetOkExists("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	request["InternetChargeType"] = d.Get("internet_charge_type")
	if v, ok := d.GetOk("isp"); ok {
		request["Isp"] = v
	}
	request["InstanceChargeType"] = convertEnsInstanceInstanceChargeTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("eip_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ens_eip", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AllocationId"]))

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, ensServiceV2.EnsEipStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceEnsEipRead(d, meta)
}

func resourceEnsEipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := NewService(client)

	objectRaw, err := ensServiceV2.DescribeEnsEip(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_eip DescribeEnsEip Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("create_time", objectRaw["AllocationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("eip_name", objectRaw["Name"])
	d.Set("ens_region_id", objectRaw["EnsRegionId"])
	d.Set("internet_charge_type", objectRaw["InternetChargeType"])
	d.Set("isp", objectRaw["Isp"])
	d.Set("payment_type", convertEnsEipAddressesEipAddressChargeTypeResponse(objectRaw["ChargeType"]))
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceEnsEipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyEnsEipAddressAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AllocationId"] = d.Id()

	if d.HasChange("eip_name") {
		update = true
	}
	if v, ok := d.GetOk("eip_name"); ok {
		request["Name"] = v
	}
	if d.HasChange("bandwidth") {
		update = true

		if v, ok := d.GetOkExists("bandwidth"); ok {
			request["Bandwidth"] = v
		}
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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

	return resourceEnsEipRead(d, meta)
}

func resourceEnsEipDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEip"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

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

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, ensServiceV2.EnsEipStateRefreshFunc(d.Id(), "$.AllocationId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func convertEnsEipAddressesEipAddressChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
