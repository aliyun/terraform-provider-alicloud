// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEipSegmentAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEipSegmentAddressCreate,
		Read:   resourceAliCloudEipSegmentAddressRead,
		Update: resourceAliCloudEipSegmentAddressUpdate,
		Delete: resourceAliCloudEipSegmentAddressDelete,
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_mask": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"isp": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"BGP", "BGP_PRO", "ChinaTelecom", "ChinaUnicom", "ChinaMobile", "ChinaTelecom_L2", "ChinaUnicom_L2", "ChinaMobile_L2", "BGP_FinanceCloud"}, false),
			},
			"netmode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"public"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"segment_address_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEipSegmentAddressCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AllocateEipSegmentAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("isp"); ok {
		request["Isp"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	request["EipMask"] = d.Get("eip_mask")
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("netmode"); ok {
		request["Netmode"] = v
	}
	if v, ok := d.GetOk("zone"); ok {
		request["Zone"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.LastTokenProcessing", "IncorrectStatus", "SystemBusy", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eip_segment_address", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EipSegmentInstanceId"]))

	eipServiceV2 := EipServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Allocated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, eipServiceV2.EipSegmentAddressStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEipSegmentAddressRead(d, meta)
}

func resourceAliCloudEipSegmentAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eipServiceV2 := EipServiceV2{client}

	objectRaw, err := eipServiceV2.DescribeEipSegmentAddress(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eip_segment_address DescribeEipSegmentAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["Name"] != nil {
		d.Set("segment_address_name", objectRaw["Name"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["Zone"] != nil {
		d.Set("zone", objectRaw["Zone"])
	}

	return nil
}

func resourceAliCloudEipSegmentAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Segment Address.")
	return nil
}

func resourceAliCloudEipSegmentAddressDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseEipSegmentAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["SegmentInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.LastTokenProcessing", "LastTokenProcessing", "IncorrectStatus", "IncorrectEipStatus", "OperationFailed.EipStatusInvalid", "OperationConflict", "TaskConflict.AssociateGlobalAccelerationInstance", "SystemBusy", "ServiceUnavailable"}) || NeedRetry(err) {
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

	eipServiceV2 := EipServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, eipServiceV2.EipSegmentAddressStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
