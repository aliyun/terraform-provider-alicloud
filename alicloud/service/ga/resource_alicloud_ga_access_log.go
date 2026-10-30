package ga

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAlicloudGaAccessLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaAccessLogCreate,
		Read:   resourceAlicloudGaAccessLogRead,
		Delete: resourceAlicloudGaAccessLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sls_project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sls_log_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sls_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaAccessLogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "AttachLogStoreToEndpointGroup"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("AttachLogStoreToEndpointGroup")
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")
	request["EndpointGroupIds.1"] = d.Get("endpoint_group_id")
	request["SlsProjectName"] = d.Get("sls_project_name")
	request["SlsLogStoreName"] = d.Get("sls_log_store_name")
	request["SlsRegionId"] = d.Get("sls_region_id")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_access_log", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["AcceleratorId"], request["ListenerId"], d.Get("endpoint_group_id")))

	stateConf := helper.BuildStateConf([]string{}, []string{"on"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, gaService.GaAccessLogStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAlicloudGaAccessLogRead(d, meta)
}

func resourceAlicloudGaAccessLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	object, err := gaService.DescribeGaAccessLog(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("endpoint_group_id", object["EndpointGroupId"])
	d.Set("sls_project_name", object["SlsProjectName"])
	d.Set("sls_log_store_name", object["SlsLogStoreName"])
	d.Set("sls_region_id", object["SlsRegionId"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudGaAccessLogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "DetachLogStoreFromEndpointGroup"

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"AcceleratorId":    parts[0],
		"ListenerId":       parts[1],
		"EndpointGroupIds": []string{parts[2]},
	}
	request["ClientToken"] = helper.BuildClientToken("DetachLogStoreFromEndpointGroup")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotExist.Accelerator", "NotExist.EndPointGroup"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, gaService.GaAccessLogStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
