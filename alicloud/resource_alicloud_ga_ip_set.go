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

func resourceAliCloudGaIpSet() *schema.Resource {
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
				ValidateFunc: StringInSlice([]string{"IPv4", "IPv6", "DUAL_STACK"}, false),
			},
			"isp_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"BGP", "BGP_PRO"}, false),
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
	gaService := GaService{client}
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
	request["ClientToken"] = buildClientToken("CreateIpSets")
	action := "CreateIpSets"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet", "NotExist.BasicBandwidthPackage", "NotSuitable.RegionSelection", "NotActive.Listener"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_ip_set", action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.IpSets", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	response = v.([]interface{})[0].(map[string]interface{})
	d.SetId(fmt.Sprint(response["IpSetId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaIpSetRead(d, meta)
}

func resourceAliCloudGaIpSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaIpSet(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_ip_set gaService.DescribeGaIpSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("accelerate_region_id", object["AccelerateRegionId"])
	d.Set("bandwidth", formatInt(object["Bandwidth"]))
	d.Set("ip_version", object["IpVersion"])
	d.Set("isp_type", object["IspType"])
	d.Set("ip_address_list", object["IpAddressList"])
	d.Set("status", object["State"])

	return nil
}

func resourceAliCloudGaIpSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("UpdateIpSet")
			response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet", "GreaterThanGa.IpSetBandwidth", "NotActive.Listener"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGaIpSetRead(d, meta)
}

func resourceAliCloudGaIpSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var err error

	request := map[string]interface{}{
		"IpSetIds.1": d.Id(),
		"RegionId":   client.RegionId,
	}

	action := "DeleteIpSets"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteIpSet")
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet", "NotActive.Listener"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.IpSets"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, gaService.GaIpSetStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
