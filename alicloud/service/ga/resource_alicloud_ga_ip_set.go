package ga

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudGaIpSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaIpSetCreate,
		Read:   resourceAliCloudGaIpSetRead,
		Update: resourceAliCloudGaIpSetUpdate,
		Delete: resourceAliCloudGaIpSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accelerate_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"IPv4", "IPv6", "DUAL_STACK"}, false),
			},
			"isp_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"BGP", "BGP_PRO"}, false),
			},
			"ip_address_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaIpSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	var err error

	request["AcceleratorId"] = d.Get("accelerator_id")
	request["AccelerateRegion.1.AccelerateRegionId"] = d.Get("accelerate_region_id")

	if v, ok := d.GetOk("bandwidth"); ok {
		request["AccelerateRegion.1.Bandwidth"] = v
	}

	if v, ok := d.GetOk("ip_version"); ok {
		request["AccelerateRegion.1.IpVersion"] = v
	}

	if v, ok := d.GetOk("isp_type"); ok {
		request["AccelerateRegion.1.IspType"] = v
	}

	var response map[string]interface{}
	request["ClientToken"] = helper.BuildClientToken("CreateIpSets")
	action := "CreateIpSets"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet", "NotExist.BasicBandwidthPackage", "NotSuitable.RegionSelection", "NotActive.Listener"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_ip_set", action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.IpSets", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	response = v.([]interface{})[0].(map[string]interface{})
	d.SetId(fmt.Sprint(response["IpSetId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudGaIpSetRead(d, meta)
}

func resourceAliCloudGaIpSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaIpSet(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_ip_set gaService.DescribeGaIpSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("accelerate_region_id", object["AccelerateRegionId"])
	bandwidthVal, err := convert.FormatInt(object["Bandwidth"])
	if err != nil {
		return err
	}
	d.Set("bandwidth", bandwidthVal)
	d.Set("ip_version", object["IpVersion"])
	d.Set("isp_type", object["IspType"])
	d.Set("ip_address_list", object["IpAddressList"])
	d.Set("status", object["State"])

	return nil
}

func resourceAliCloudGaIpSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var err error

	update := false
	request := map[string]interface{}{
		"IpSetId":  d.Id(),
		"RegionId": client.RegionId,
	}

	if d.HasChange("bandwidth") {
		update = true
	}
	request["Bandwidth"] = d.Get("bandwidth")

	if update {
		action := "UpdateIpSet"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			request["ClientToken"] = helper.BuildClientToken("UpdateIpSet")
			response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet", "GreaterThanGa.IpSetBandwidth", "NotActive.Listener"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudGaIpSetRead(d, meta)
}

func resourceAliCloudGaIpSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var err error

	request := map[string]interface{}{
		"IpSetIds.1": d.Id(),
		"RegionId":   client.RegionId,
	}

	action := "DeleteIpSets"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		request["ClientToken"] = helper.BuildClientToken("DeleteIpSet")
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet", "NotActive.Listener"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotExist.IpSets"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
