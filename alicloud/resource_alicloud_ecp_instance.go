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

func resourceAliCloudEcpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcpInstanceCreate,
		Read:   resourceAliCloudEcpInstanceRead,
		Update: resourceAliCloudEcpInstanceUpdate,
		Delete: resourceAliCloudEcpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"eip_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"resolution": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vnc_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Running", "Stopped"}, false),
			},
		},
	}
}

func resourceAliCloudEcpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudPhoneService := CloudphoneService{client}
	var response map[string]interface{}
	action := "RunInstances"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("RunInstances")
	request["Amount"] = 1
	request["InstanceType"] = d.Get("instance_type")
	request["ImageId"] = d.Get("image_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["SecurityGroupId"] = d.Get("security_group_id")

	if v, ok := d.GetOkExists("eip_bandwidth"); ok {
		request["EipBandwidth"] = v
	}

	if v, ok := d.GetOk("resolution"); ok {
		request["Resolution"] = v
	}

	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertEcpSyncPaymentTypeRequest(v.(string))
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}

	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecp_instance", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.InstanceIds.InstanceId", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ecp_instance")
	} else {
		instanceId := resp.([]interface{})[0]
		d.SetId(fmt.Sprint(instanceId))
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudPhoneService.EcpInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcpInstanceUpdate(d, meta)
}

func resourceAliCloudEcpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudPhoneService := CloudphoneService{client}

	object, err := cloudPhoneService.DescribeEcpInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecp_instance cloudPhoneService.DescribeEcpInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_type", object["InstanceType"])
	d.Set("image_id", object["ImageId"])
	d.Set("vswitch_id", object["VpcAttributes"].(map[string]interface{})["VSwitchId"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("resolution", object["Resolution"])
	d.Set("key_pair_name", object["KeyPairName"])
	d.Set("payment_type", convertEcpSyncPaymentTypeResponse(object["ChargeType"]))
	d.Set("instance_name", object["InstanceName"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])

	if eipAddress, ok := object["EipAddress"]; ok {
		eipAddressArg := eipAddress.(map[string]interface{})

		if eipBandwidth, ok := eipAddressArg["Bandwidth"]; ok {
			d.Set("eip_bandwidth", eipBandwidth)
		}
	}

	return nil
}

func resourceAliCloudEcpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudPhoneService := CloudphoneService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"InstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("key_pair_name") {
		update = true
	}
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	}

	if d.HasChange("vnc_password") {
		update = true

		if v, ok := d.GetOk("vnc_password"); ok {
			request["VncPassword"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "UpdateInstanceAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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

		d.SetPartial("key_pair_name")
		d.SetPartial("vnc_password")
		d.SetPartial("instance_name")
		d.SetPartial("description")
	}

	update = false
	updateResolutionReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"InstanceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resolution") {
		update = true
	}
	if v, ok := d.GetOk("resolution"); ok {
		updateResolutionReq["Resolution"] = v
	}

	if update {
		object, err := cloudPhoneService.DescribeEcpInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}

		if object["Status"].(string) != "Stopped" {
			err := cloudPhoneService.ModifyEcpInstanceStatus(d, "Stopped")
			if err != nil {
				return WrapError(err)
			}
		}

		action := "UpdateInstanceAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, updateResolutionReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateResolutionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if object["Status"].(string) != "Running" {
			err := cloudPhoneService.ModifyEcpInstanceStatus(d, "Running")
			if err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("resolution")
	}

	if d.HasChange("status") {
		object, err := cloudPhoneService.DescribeEcpInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Running" {
				err := cloudPhoneService.ModifyEcpInstanceStatus(d, "Running")
				if err != nil {
					return WrapError(err)
				}
			}

			if target == "Stopped" {
				err := cloudPhoneService.ModifyEcpInstanceStatus(d, "Stopped")
				if err != nil {
					return WrapError(err)
				}
			}
		}

		d.SetPartial("status")
	}

	d.Partial(false)

	return resourceAliCloudEcpInstanceRead(d, meta)
}

func resourceAliCloudEcpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("payment_type"); ok && fmt.Sprint(v) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resource alicloud_ecp_instance which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstances"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"InstanceId": []string{d.Id()},
	}

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"CloudPhoneInstances.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertEcpSyncPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}

	return source
}

func convertEcpSyncPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}

	return source
}
