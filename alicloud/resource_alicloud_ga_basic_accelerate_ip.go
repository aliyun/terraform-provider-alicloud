package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaBasicAccelerateIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaBasicAccelerateIpCreate,
		Read:   resourceAliCloudGaBasicAccelerateIpRead,
		Delete: resourceAliCloudGaBasicAccelerateIpDelete,
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
			"ip_set_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accelerate_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaBasicAccelerateIpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateBasicAccelerateIp"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateBasicAccelerateIp")
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["IpSetId"] = d.Get("ip_set_id")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "NotActive.IpSet"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_basic_accelerate_ip", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AccelerateIpId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaBasicAccelerateIpStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaBasicAccelerateIpRead(d, meta)
}

func resourceAliCloudGaBasicAccelerateIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaBasicAccelerateIp(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("ip_set_id", object["IpSetId"])
	d.Set("accelerate_ip_address", object["AccelerateIpAddress"])
	d.Set("status", object["State"])

	return nil
}

func resourceAliCloudGaBasicAccelerateIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteBasicAccelerateIp"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"ClientToken":    buildClientToken("DeleteBasicAccelerateIp"),
		"AccelerateIpId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "NotActive.IpSet"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.AccelerateIp"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaBasicAccelerateIpStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
