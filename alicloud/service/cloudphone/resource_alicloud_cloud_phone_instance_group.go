// Package cloudphone. This file is generated automatically. Please do not modify it manually, thank you!
package cloudphone

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceGroupCreate,
		Read:   resourceInstanceGroupRead,
		Update: resourceInstanceGroupUpdate,
		Delete: resourceInstanceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(9 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gpu_acceleration": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_group_spec": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"number_of_instances": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceInstanceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAndroidInstanceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["BizRegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("instance_group_name"); ok {
		request["InstanceGroupName"] = v
	}
	if v, ok := d.GetOk("office_site_id"); ok {
		request["OfficeSiteId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	request["ImageId"] = d.Get("image_id")
	if v, ok := d.GetOkExists("number_of_instances"); ok {
		request["NumberOfInstances"] = v
	}
	if v, ok := d.GetOkExists("amount"); ok {
		request["Amount"] = v
	}
	if v, ok := d.GetOk("charge_type"); ok {
		request["ChargeType"] = v
	}
	request["InstanceGroupSpec"] = d.Get("instance_group_spec")
	if v, ok := d.GetOkExists("gpu_acceleration"); ok {
		request["GpuAcceleration"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("policy_group_id"); ok {
		request["PolicyGroupId"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_phone_instance_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.InstanceGroupInfos[0].InstanceGroupId", response)
	d.SetId(fmt.Sprint(id))

	cloudPhoneService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, cloudPhoneService.InstanceGroupStateRefreshFunc(d.Id(), "InstanceGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceInstanceGroupRead(d, meta)
}

func resourceInstanceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudPhoneService := NewService(client)

	objectRaw, err := cloudPhoneService.DescribeInstanceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_phone_instance_group DescribeInstanceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("charge_type", objectRaw["ChargeType"])
	d.Set("image_id", objectRaw["ImageId"])
	d.Set("instance_group_name", objectRaw["InstanceGroupName"])
	d.Set("instance_group_spec", objectRaw["InstanceGroupSpec"])
	number_of_instancesVal, err := convert.FormatInt(objectRaw["NumberOfInstances"])
	if err != nil {
		return err
	}
	d.Set("number_of_instances", number_of_instancesVal)
	d.Set("status", objectRaw["InstanceGroupStatus"])

	return nil
}

func resourceInstanceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyAndroidInstanceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceGroupId"] = d.Id()

	if d.HasChange("instance_group_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_group_name"); ok || d.HasChange("instance_group_name") {
		request["NewInstanceGroupName"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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

	return resourceInstanceGroupRead(d, meta)
}

func resourceInstanceGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAndroidInstanceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceGroupIds.1"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)

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

	cloudPhoneService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudPhoneService.InstanceGroupStateRefreshFunc(d.Id(), "InstanceGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
