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

func resourceAliCloudCloudPhoneInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudPhoneInstanceGroupCreate,
		Read:   resourceAliCloudCloudPhoneInstanceGroupRead,
		Update: resourceAliCloudCloudPhoneInstanceGroupUpdate,
		Delete: resourceAliCloudCloudPhoneInstanceGroupDelete,
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

func resourceAliCloudCloudPhoneInstanceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAndroidInstanceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["BizRegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_phone_instance_group", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.InstanceGroupInfos[0].InstanceGroupId", response)
	d.SetId(fmt.Sprint(id))

	cloudPhoneServiceV2 := CloudPhoneServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, cloudPhoneServiceV2.CloudPhoneInstanceGroupStateRefreshFunc(d.Id(), "InstanceGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudPhoneInstanceGroupRead(d, meta)
}

func resourceAliCloudCloudPhoneInstanceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudPhoneServiceV2 := CloudPhoneServiceV2{client}

	objectRaw, err := cloudPhoneServiceV2.DescribeCloudPhoneInstanceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_phone_instance_group DescribeCloudPhoneInstanceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("charge_type", objectRaw["ChargeType"])
	d.Set("image_id", objectRaw["ImageId"])
	d.Set("instance_group_name", objectRaw["InstanceGroupName"])
	d.Set("instance_group_spec", objectRaw["InstanceGroupSpec"])
	d.Set("number_of_instances", formatInt(objectRaw["NumberOfInstances"]))
	d.Set("status", objectRaw["InstanceGroupStatus"])

	return nil
}

func resourceAliCloudCloudPhoneInstanceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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

	return resourceAliCloudCloudPhoneInstanceGroupRead(d, meta)
}

func resourceAliCloudCloudPhoneInstanceGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAndroidInstanceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceGroupIds.1"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)

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

	cloudPhoneServiceV2 := CloudPhoneServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudPhoneServiceV2.CloudPhoneInstanceGroupStateRefreshFunc(d.Id(), "InstanceGroupStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
